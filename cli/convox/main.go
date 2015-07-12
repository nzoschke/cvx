package convox

import (
	"os"

	"github.com/nzoschke/convox/Godeps/_workspace/src/github.com/codegangsta/cli"
)

func main() {
	cli.NewApp().Run(os.Args)
}

func Run(args []string) string {
	if len(args) == 2 {
		return `
usage: convox [--version] [--help] [--app=<name>]

Commands:
  apps  List, create, or delete apps
`
	} else {
		return `
 * app1
  app2
`
	}
}
