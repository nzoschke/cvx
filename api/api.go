package api

import (
	"encoding/json"
	"net/http"
	"os"

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

		apps := make(Apps, 0)

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

	return mux
}
