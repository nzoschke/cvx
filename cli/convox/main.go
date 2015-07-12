package convox

import (
	"os"

	"github.com/nzoschke/convox/Godeps/_workspace/src/github.com/codegangsta/cli"
)

func Run() error {
	app := cli.NewApp()
	app.Name = "convox"
	return app.Run(os.Args)
}
