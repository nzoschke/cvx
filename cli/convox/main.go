package convox

import (
	"os"

	"github.com/nzoschke/convox/Godeps/_workspace/src/github.com/codegangsta/cli"
)

func Run() error {
	app := cli.NewApp()
	app.Name = "convox"

	app.Commands = []cli.Command{
		{
			Name: "apps",
			Action: func(c *cli.Context) {
				println("listing apps")
			},
		},
	}
	return app.Run(os.Args)
}
