package main

import (
	"bytes"
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
	convox "github.com/nzoschke/cvx/convox"

	"github.com/nzoschke/cvx/Godeps/_workspace/src/github.com/aws/aws-sdk-go/aws"
	"github.com/nzoschke/cvx/Godeps/_workspace/src/github.com/aws/aws-sdk-go/internal/protocol/xml/xmlutil"
	"github.com/nzoschke/cvx/Godeps/_workspace/src/github.com/aws/aws-sdk-go/service/cloudformation"
)

type Case struct {
	got, want interface{}
}

type Out struct {
	stdout, stderr string
}

type Cases []Case

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

	help := `NAME:
   convox - A new cli application

USAGE:
   convox [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
   apps		
   help, h	Shows a list of commands or help for one command
   
GLOBAL OPTIONS:
   --help, -h		show help
   --version, -v	print the version
   
`

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
		{Run([]string{"convox", "help"}), Out{help, ""}},
		{Run([]string{"convox", "apps"}), Out{text, ""}},
		{Run([]string{"convox", "apps", "--output", "json"}), Out{json, ""}},
	}

	assert(t, cases)
}

func assert(t *testing.T, cases Cases) {
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("got `%q` want `%q`", c.got, c.want)
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

		t := reflect.TypeOf(output).Name()
		t = strings.Replace(t, "Output", "", 1)

		w.Header().Set("Content-Type", "text/xml")
		w.Header().Set("X-Amzn-Requestid", "b123290e-28ae-11e5-b834-6f3c1afbf01a")

		w.Write([]byte(fmt.Sprintf("<%sResponse><%sResult>%s</%sResult><ResponseMetadata><RequestId>b123290e-28ae-11e5-b834-6f3c1afbf01a</RequestId></ResponseMetadata></%sResponse>", t, t, b.String(), t, t)))
	}))

	aws.DefaultConfig.Endpoint = s.URL

	return s
}

func NewApiServer() *httptest.Server {
	s := httptest.NewServer(api.Handler())

	convox.DefaultConfig.Endpoint = s.URL

	return s
}

func Run(args []string) Out {
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

	// Run CLI app
	os.Args = args
	convox.Run()

	// restore stderr, stdout
	ew.Close()
	os.Stderr = oldErr
	err := <-errC

	ow.Close()
	os.Stdout = oldOut
	out := <-outC

	return Out{out, err}
}
