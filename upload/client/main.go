package main

import (
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Upload Client"
	app.Usage = "Multi File Transferred"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		uploadCommand(),
	}
	if err := app.Run(os.Args); err != nil {
		return
	}
}
