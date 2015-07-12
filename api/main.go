package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/nzoschke/convox/Godeps/_workspace/src/github.com/aws/aws-sdk-go/aws"
	"github.com/nzoschke/convox/Godeps/_workspace/src/github.com/aws/aws-sdk-go/service/cloudformation"

	"github.com/nzoschke/convox/Godeps/_workspace/src/github.com/codegangsta/negroni"
)

type App struct {
	Name string

	Status string

	// Outputs    map[string]string
	// Parameters map[string]string
	Tags map[string]string
}

type Apps []App

func main() {
	mux := Handler()
	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(":3000")
}

func TestHandler(awsEndpoint string) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/apps", func(w http.ResponseWriter, req *http.Request) {
		svc := cloudformation.New(&aws.Config{
			Region:      "us-west-2",
			Logger:      os.Stdout,
			LogLevel:    0,
			LogHTTPBody: true,
		})

		svc.Endpoint = awsEndpoint

		res, err := svc.DescribeStacks(&cloudformation.DescribeStacksInput{})

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// res, ok := awsOutput.(cloudformation.DescribeStacksOutput)

		// if !ok {
		// 	http.Error(w, "error converting to cloudformation.DescribeStacksOutput", 500)
		// 	return
		// }

		// fmt.Printf("%+v\n", res.Stacks)
		apps := make(Apps, 0)

		for _, stack := range res.Stacks {
			tags := make(map[string]string)

			for _, tag := range stack.Tags {
				tags[*tag.Key] = *tag.Value
			}

			// if tags["Type"] == "app" {
			app := App{
				Name: *stack.StackName,
				Tags: tags,
			}

			apps = append(apps, app)
			// }
		}

		b, err := json.Marshal(apps)

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Write(b)
	})

	return mux
}

func Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/apps", func(w http.ResponseWriter, req *http.Request) {
		svc := cloudformation.New(&aws.Config{
			Region: "us-west-2",
		})

		fmt.Printf("%+v\n", svc.Endpoint)

		res, err := svc.DescribeStacks(&cloudformation.DescribeStacksInput{})

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		b, err := json.Marshal(res)

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Write(b)
	})

	return mux
}

func Set(url string, svcOut interface{}) string {
	// mux := http.NewServeMux()
	// mux.HandleFunc(url, func(w http.ResponseWriter, req *http.Request) {
	// })
	// b, err := json.Marshal(res)
	return ""
}
