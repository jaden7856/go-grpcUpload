package main

import (
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Server"
	app.Usage = "Multi File Transferred Server"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		ServerCommand(),
	}
	err := app.Run(os.Args)
	if err != nil {
		return
	}
}
