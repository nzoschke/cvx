package main

import (
	"encoding/json"
	"net/http"

	"github.com/nzoschke/convox/Godeps/_workspace/src/github.com/aws/aws-sdk-go/aws"
	"github.com/nzoschke/convox/Godeps/_workspace/src/github.com/aws/aws-sdk-go/service/cloudformation"

	"github.com/nzoschke/convox/Godeps/_workspace/src/github.com/codegangsta/negroni"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/apps", func(w http.ResponseWriter, req *http.Request) {
		svc := cloudformation.New(&aws.Config{
			Region: "us-west-2",
		})

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

	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(":3000")
}
