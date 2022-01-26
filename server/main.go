package main

import (
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Server"
	app.Usage = "Multi File Transferer Server"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		ServerCommand(),
	}
	app.Run(os.Args)
}
