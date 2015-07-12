package main

import (
	"os"

	"github.com/nzoschke/convox/Godeps/_workspace/src/github.com/codegangsta/cli"
)

func main() {
	cli.NewApp().Run(os.Args)
}
