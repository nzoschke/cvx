package main

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "net/http/httptest"
  "testing"

  "github.com/aws/aws-sdk-go/aws"
  "github.com/nzoschke/convox/Godeps/_workspace/src/github.com/aws/aws-sdk-go/service/cloudformation"
)

func TestApps(t *testing.T) {
  ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    out := cloudformation.DescribeStacksOutput{
      Stacks: []*cloudformation.Stack{
        {
          StackID:   aws.String("arn:aws:cloudformation:us-east-1:901416387788:stack/app1/a9196ca0-24e3-11e5-a58b-500150b34c7c"),
          StackName: aws.String("app1"),
          Tags: []*cloudformation.Tag{
            {
              Key:   aws.String("Type"),
              Value: aws.String("app"),
            },
          },
        },
        {
          StackID:   aws.String("arn:aws:cloudformation:us-east-1:901416387788:stack/app2/185779b0-1632-11e5-98be-50d501114c2c"),
          StackName: aws.String("app2"),
          Tags: []*cloudformation.Tag{
            {
              Key:   aws.String("Type"),
              Value: aws.String("app"),
            },
          },
        },
      },
    }

    b, err := json.Marshal(out)

    if err != nil {
      t.Error(err)
    }

    w.Write(b)
  }))
  defer ts.Close()

  res, err := http.Get(ts.URL)

  body, err := ioutil.ReadAll(res.Body)
  res.Body.Close()
  if err != nil {
    log.Fatal(err)
  }

  fmt.Printf("%+v\n", string(body))

  // GET /apps
  // ['convox', 'apps']

  // cloudformation.DescribeStacksOutput
  // aws_output = `{"Stacks": []}`
}
