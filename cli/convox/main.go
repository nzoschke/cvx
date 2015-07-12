package convox

import (
	"fmt"
	"os"

	"github.com/nzoschke/convox/Godeps/_workspace/src/github.com/codegangsta/cli"
)

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
		},
	}
	return app.Run(os.Args)
}

func cmdApps(c *cli.Context) {
	fmt.Printf("app1\napp2\n")
	// data, err := ConvoxGet("/apps")

	// if err != nil {
	// 	stdcli.Error(err)
	// 	return
	// }

	// var apps *Apps
	// err = json.Unmarshal(data, &apps)

	// if err != nil {
	// 	stdcli.Error(err)
	// 	return
	// }

	// for _, app := range *apps {
	// 	fmt.Printf("%s\n", app.Name)
	// }
}
