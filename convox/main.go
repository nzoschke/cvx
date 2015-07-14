package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/nzoschke/cvx/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/nzoschke/cvx/api"
)

var DefaultConfig = &Config{
	Endpoint: "http://localhost:3000",
}

type Config struct {
	Endpoint string
}

func main() {
	Run()
}

func Run() error {
	app := cli.NewApp()
	app.Name = "convox"

	app.Commands = []cli.Command{
		{
			Name:   "apps",
			Action: Apps,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "output",
					Usage: "output 'text' or 'json'",
					Value: "text",
				},
			},
		},
		{
			Name:   "builds",
			Action: Builds,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "output",
					Usage: "output 'text' or 'json'",
					Value: "text",
				},
			},
		},
		{
			Name:   "stacks",
			Action: Stacks,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "output",
					Usage: "output 'text' or 'json'",
					Value: "text",
				},
			},
		},
	}
	return app.Run(os.Args)
}

func Apps(c *cli.Context) {
	res, err := http.Get(DefaultConfig.Endpoint + "/apps")

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	var apps *api.Apps
	err = json.Unmarshal(body, &apps)

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	switch c.String("output") {
	case "json":
		fmt.Printf("%s\n", body)
	case "text":
		for _, app := range *apps {
			fmt.Printf("%s\n", app.Name)
		}
	}
}

func Builds(c *cli.Context) {
	res, err := http.Get(DefaultConfig.Endpoint + "/builds")

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	var builds *api.Builds
	err = json.Unmarshal(body, &builds)

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	switch c.String("output") {
	case "json":
		fmt.Printf("%s\n", body)
	case "text":
		for _, build := range *builds {
			fmt.Printf("%s\n", build.Id)
		}
	}
}

func Stacks(c *cli.Context) {
	res, err := http.Get(DefaultConfig.Endpoint + "/stacks")

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	fmt.Printf("%s\n", body)
}
