package api

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/nzoschke/cvx/Godeps/_workspace/src/github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/nzoschke/cvx/Godeps/_workspace/src/github.com/aws/aws-sdk-go/aws"
	"github.com/nzoschke/cvx/Godeps/_workspace/src/github.com/aws/aws-sdk-go/service/cloudformation"

	"github.com/nzoschke/cvx/Godeps/_workspace/src/github.com/codegangsta/negroni"
)

type App struct {
	Name string

	Status string

	// Outputs    map[string]string
	// Parameters map[string]string
	Tags map[string]string
}

type Apps []App

type Build struct {
	Id string

	App string

	Logs    string
	Release string
	Status  string

	Started time.Time
	Ended   time.Time

	kinesis string
}

type Builds []Build

var SortableTime = "20060102.150405.000000000"

func Run() {
	mux := Handler()
	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(":3000")
}

func Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/apps", func(w http.ResponseWriter, req *http.Request) {
		svc := cloudformation.New(&aws.Config{
			Region:      "us-east-1",
			Logger:      os.Stdout,
			LogLevel:    0,
			LogHTTPBody: true,
		})

		res, err := svc.DescribeStacks(&cloudformation.DescribeStacksInput{})

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		apps := make(Apps, len(res.Stacks))

		for _, stack := range res.Stacks {
			tags := make(map[string]string)

			for _, tag := range stack.Tags {
				tags[*tag.Key] = *tag.Value
			}

			if tags["Type"] == "app" {
				app := App{
					Name: *stack.StackName,
					Tags: tags,
				}

				apps = append(apps, app)
			}
		}

		b, err := json.MarshalIndent(apps, "", "  ")

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Write(b)
	})

	mux.HandleFunc("/builds", func(w http.ResponseWriter, req *http.Request) {
		svc := dynamodb.New(&aws.Config{
			Region:      "us-east-1",
			Logger:      os.Stdout,
			LogLevel:    0,
			LogHTTPBody: true,
		})

		res, err := svc.Query(&dynamodb.QueryInput{
			KeyConditions: map[string]*dynamodb.Condition{
				"app": &dynamodb.Condition{
					AttributeValueList: []*dynamodb.AttributeValue{&dynamodb.AttributeValue{S: aws.String("lugg-api")}},
					ComparisonOperator: aws.String("EQ"),
				},
			},
			IndexName:        aws.String("app.created"),
			Limit:            aws.Long(10),
			ScanIndexForward: aws.Boolean(false),
			TableName:        aws.String("convox-builds"),
		})

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		builds := make(Builds, len(res.Items))

		for _, item := range res.Items {
			started, _ := time.Parse(SortableTime, coalesce(item["created"], ""))
			ended, _ := time.Parse(SortableTime, coalesce(item["ended"], ""))

			build := Build{
				Id:      coalesce(item["id"], ""),
				App:     coalesce(item["app"], ""),
				Logs:    coalesce(item["logs"], ""),
				Release: coalesce(item["release"], ""),
				Status:  coalesce(item["status"], ""),
				Started: started,
				Ended:   ended,
			}

			builds = append(builds, build)
		}

		b, err := json.MarshalIndent(builds, "", "  ")

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Write(b)
	})

	return mux
}

func coalesce(s *dynamodb.AttributeValue, def string) string {
	if s != nil {
		return *s.S
	} else {
		return def
	}
}