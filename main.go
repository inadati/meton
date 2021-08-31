package main

import (
	"os"

	"github.com/mattn/go-colorable"
	"github.com/meton888/meton/cmd"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	logrus.SetOutput(colorable.NewColorableStdout())

	if err := exec(); err != nil {
		logrus.Fatal(err)
	}
}

func exec() error {
	app := cli.NewApp()
	app.Name = "meton"
	app.Usage = "Very easy to build mesos and marathon clusters."
	app.HideVersion = true
	app.Authors = []*cli.Author{
		&cli.Author{
			Name:  "inadati",
			Email: "dr.inadati@gmail.com",
		},
	}
	app.Commands = []*cli.Command{
		cmd.UpCommand,
		cmd.DownCommand,
		cmd.TunnelCommand,
	}

	return app.Run(os.Args)
}
