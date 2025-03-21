package main

import (
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Upload Server"
	app.Usage = "Multi File Transferred Server"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		ServerCommand(),
	}
	if err := app.Run(os.Args); err != nil {
		return
	}
}
