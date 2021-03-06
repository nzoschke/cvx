package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/nzoschke/cvx/api"
	"github.com/nzoschke/cvx/cli"

	"github.com/nzoschke/cvx/Godeps/_workspace/src/github.com/aws/aws-sdk-go/aws"
	"github.com/nzoschke/cvx/Godeps/_workspace/src/github.com/aws/aws-sdk-go/internal/protocol/xml/xmlutil"
	"github.com/nzoschke/cvx/Godeps/_workspace/src/github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/nzoschke/cvx/Godeps/_workspace/src/github.com/aws/aws-sdk-go/service/dynamodb"
)

type Case struct {
	got, want interface{}
}

type Cases []Case

func TestCLIHelp(t *testing.T) {
	help := `NAME:
   convox - A new cli application

USAGE:
   convox [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
   apps		
   builds	
   stacks	
   help, h	Shows a list of commands or help for one command
   
GLOBAL OPTIONS:
   --help, -h		show help
   --version, -v	print the version
   
`

	cases := Cases{
		{Run([]string{"convox"}), help},
		{Run([]string{"convox", "help"}), help},
	}

	assert(t, cases)
}

func TestApps(t *testing.T) {
	awsServer := NewAwsServer(cloudformation.DescribeStacksOutput{
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
	})
	defer awsServer.Close()

	apiServer := NewApiServer()
	defer apiServer.Close()

	json := `[
  {
    "Name": "app1",
    "Status": "",
    "Tags": {
      "Type": "app"
    }
  },
  {
    "Name": "app2",
    "Status": "",
    "Tags": {
      "Type": "app"
    }
  }
]
`

	text := `app1
app2
`

	cases := Cases{
		{Run([]string{"convox", "apps"}), text},
		{Run([]string{"convox", "apps", "--output", "json"}), json},
	}

	assert(t, cases)
}

func TestBuilds(t *testing.T) {
	aws.DefaultConfig.Endpoint = ""

	awsServer := NewAwsServer(dynamodb.QueryOutput{
		Count: aws.Long(2),
		Items: []map[string]*dynamodb.AttributeValue{
			{
				"id":     {S: aws.String("BOBSFPGWQBY")},
				"app":    {S: aws.String("app1")},
				"status": {S: aws.String("complete")},
			},
			{
				"id":     {S: aws.String("BFEOTKNIURY")},
				"app":    {S: aws.String("app1")},
				"status": {S: aws.String("failed")},
			},
		},
	})
	defer awsServer.Close()

	apiServer := NewApiServer()
	defer apiServer.Close()

	json := `[
  {
    "Id": "BOBSFPGWQBY",
    "App": "app1",
    "Logs": "",
    "Release": "",
    "Status": "complete",
    "Started": "0001-01-01T00:00:00Z",
    "Ended": "0001-01-01T00:00:00Z"
  },
  {
    "Id": "BFEOTKNIURY",
    "App": "app1",
    "Logs": "",
    "Release": "",
    "Status": "failed",
    "Started": "0001-01-01T00:00:00Z",
    "Ended": "0001-01-01T00:00:00Z"
  }
]
`

	text := `app1 BOBSFPGWQBY complete
app1 BFEOTKNIURY failed
`

	cases := Cases{
		{Run([]string{"convox", "builds"}), text},
		{Run([]string{"convox", "builds", "--output", "json"}), json},
	}

	assert(t, cases)
}

func TestStacks(t *testing.T) {
	awsServer := NewAwsServer(cloudformation.DescribeStacksOutput{
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
	})
	defer awsServer.Close()

	apiServer := NewApiServer()
	defer apiServer.Close()

	json := "[\n  {\n    \"Capabilities\": null,\n    \"CreationTime\": null,\n    \"Description\": null,\n    \"DisableRollback\": null,\n    \"LastUpdatedTime\": null,\n    \"NotificationARNs\": null,\n    \"Outputs\": null,\n    \"Parameters\": null,\n    \"StackID\": \"arn:aws:cloudformation:us-east-1:901416387788:stack/app1/a9196ca0-24e3-11e5-a58b-500150b34c7c\",\n    \"StackName\": \"app1\",\n    \"StackStatus\": null,\n    \"StackStatusReason\": null,\n    \"Tags\": [\n      {\n        \"Key\": \"Type\",\n        \"Value\": \"app\"\n      }\n    ],\n    \"TimeoutInMinutes\": null\n  },\n  {\n    \"Capabilities\": null,\n    \"CreationTime\": null,\n    \"Description\": null,\n    \"DisableRollback\": null,\n    \"LastUpdatedTime\": null,\n    \"NotificationARNs\": null,\n    \"Outputs\": null,\n    \"Parameters\": null,\n    \"StackID\": \"arn:aws:cloudformation:us-east-1:901416387788:stack/app2/185779b0-1632-11e5-98be-50d501114c2c\",\n    \"StackName\": \"app2\",\n    \"StackStatus\": null,\n    \"StackStatusReason\": null,\n    \"Tags\": [\n      {\n        \"Key\": \"Type\",\n        \"Value\": \"app\"\n      }\n    ],\n    \"TimeoutInMinutes\": null\n  }\n]\n"
	cases := Cases{
		{Run([]string{"convox", "stacks", "--output", "json", "--verbose"}), json},
	}

	assert(t, cases)
}

func assert(t *testing.T, cases Cases) {
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("\n%q\n%q", c.got, c.want)
		}
	}
}

func Get(t *testing.T, url string) string {
	res, err := http.Get(url)

	if err != nil {
		t.Error(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		t.Error(err)
	}

	return string(body)
}

func NewAwsServer(output interface{}) *httptest.Server {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var b bytes.Buffer
		enc := xml.NewEncoder(&b)
		xmlutil.BuildXML(output, enc)

		t := reflect.TypeOf(output)

		p := t.PkgPath()
		parts := strings.Split(p, "/")
		service := parts[len(parts)-1]

		switch service {
		case "dynamodb", "ecs": // jsonrpc services
			b, err := json.Marshal(output)

			if err != nil {
				fmt.Println("error:", err)
			}

			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Amzn-Requestid", "24b57e1a-8Bcb-4655-aea5-d00864d514ac")

			w.Write(b)
		case "cloudformation": // XML services
			n := t.Name()
			n = strings.Replace(n, "Output", "", 1)

			w.Header().Set("Content-Type", "text/xml")
			w.Header().Set("X-Amzn-Requestid", "b123290e-28ae-11e5-b834-6f3c1afbf01a")

			w.Write([]byte(fmt.Sprintf("<%sResponse><%sResult>%s</%sResult><ResponseMetadata><RequestId>b123290e-28ae-11e5-b834-6f3c1afbf01a</RequestId></ResponseMetadata></%sResponse>", n, n, b.String(), n, n)))
		}
	}))

	aws.DefaultConfig.Endpoint = s.URL

	return s
}

func NewApiServer() *httptest.Server {
	s := httptest.NewServer(api.Handler())

	cli.DefaultConfig.Endpoint = s.URL

	return s
}

func Run(args []string) string {
	// Capture stdout and stderr to strings via Pipes
	oldErr := os.Stderr
	oldOut := os.Stdout

	er, ew, _ := os.Pipe()
	or, ow, _ := os.Pipe()

	os.Stderr = ew
	os.Stdout = ow

	errC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, er)
		errC <- buf.String()
	}()

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, or)
		outC <- buf.String()
	}()

	os.Args = args
	cli.Run()

	// restore stderr, stdout
	ew.Close()
	os.Stderr = oldErr
	<-errC

	ow.Close()
	os.Stdout = oldOut
	out := <-outC

	return string(out)
}
