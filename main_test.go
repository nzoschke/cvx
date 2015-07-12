package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/nzoschke/convox/api"
	"github.com/nzoschke/convox/cli/convox"

	"github.com/nzoschke/convox/Godeps/_workspace/src/github.com/aws/aws-sdk-go/aws"
	"github.com/nzoschke/convox/Godeps/_workspace/src/github.com/aws/aws-sdk-go/internal/protocol/xml/xmlutil"
	"github.com/nzoschke/convox/Godeps/_workspace/src/github.com/aws/aws-sdk-go/service/cloudformation"
)

type Case struct {
	got, want interface{}
}

type Cases []Case

func TestHttpServer(t *testing.T) {
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

	api := httptest.NewServer(api.Handler())
	defer api.Close()

	res, err := http.Get(api.URL + "/apps")
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		log.Fatal(err)
	}

	cases := Cases{
		{res.StatusCode, 200},
		{string(body), `[{"Name":"app1","Status":"","Tags":{"Type":"app"}},{"Name":"app2","Status":"","Tags":{"Type":"app"}}]`},
		// {string(body), `[{"Name":"s-150531195417","Status":"","Tags":{}},{"Name":"s-150531193549","Status":"","Tags":{}},{"Name":"s-1505311901","Status":"","Tags":{}},{"Name":"s","Status":"","Tags":{}},{"Name":"staging","Status":"","Tags":{}}]`},
	}

	assert(t, cases)
}

func TestApps(t *testing.T) {
	help := `
usage: convox [--version] [--help] [--app=<name>]

Commands:
  apps  List, create, or delete apps
`

	out := `
 * app1
  app2
`

	cases := Cases{
		{convox.Run([]string{"apps", "help"}), help},
		{convox.Run([]string{"apps"}), out},
	}

	assert(t, cases)
}

func assert(t *testing.T, cases Cases) {
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("got `%v` want `%v`", c.got, c.want)
		}
	}
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
