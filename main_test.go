package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nzoschke/convox/api"
	"github.com/nzoschke/convox/cli/convox"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/nzoschke/convox/Godeps/_workspace/src/github.com/aws/aws-sdk-go/service/cloudformation"
)

type Case struct {
	got, want interface{}
}

type Cases []Case

func TestHttp(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "something failed", http.StatusInternalServerError)
	}

	req, err := http.NewRequest("GET", "http://example.com/foo", nil)
	if err != nil {
		log.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler(w, req)

	cases := Cases{
		{w.Code, 500},
		{w.Body.String(), "something failed\n"},
	}

	assert(t, cases)
}

func TestHttpServer(t *testing.T) {
	svcOut := cloudformation.DescribeStacksOutput{
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

	ts := httptest.NewServer(api.TestHandler(svcOut))
	defer ts.Close()

	res, err := http.Get(ts.URL + "/apps")
	if err != nil {
		log.Fatal(err)
	}

	greeting, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		log.Fatal(err)
	}

	cases := Cases{
		{res.StatusCode, 200},
		{string(greeting), `[{"Name":"app1","Status":"","Tags":{"Type":"app"}},{"Name":"app2","Status":"","Tags":{"Type":"app"}}]`},
	}

	assert(t, cases)

}

func TestApps(t *testing.T) {
	svcOut := cloudformation.DescribeStacksOutput{
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

	api.Set("/apps", svcOut)

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
