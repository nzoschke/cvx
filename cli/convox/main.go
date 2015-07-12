package convox

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/nzoschke/convox/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/nzoschke/convox/api"
)

var DefaultConfig = &Config{
	Endpoint: "",
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
			Action: cmdApps,
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

func cmdApps(c *cli.Context) {
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

	// switch
	if c.String("output") == "json" {
		fmt.Printf("%s\n", body)
		return
	}

	for _, app := range *apps {
		fmt.Printf("%s\n", app.Name)
	}
}
