package main

import (
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Server-test"
	app.Usage = "Multi File Transfer Server"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		ServerCommand(),
	}
	app.Run(os.Args)
}
