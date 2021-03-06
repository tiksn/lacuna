package main

import (
	"os"

	"github.com/urfave/cli"
)

var app = cli.NewApp()

func info() {
	app.Name = "lacuna"
	app.Usage = "Docker base image version setter CLI"
	app.Author = "Tigran TIKSN Torosyan"
	app.Version = "1.0.0"
}

func commands() {
	app.Commands = []cli.Command{
		{
			Name:    "format",
			Aliases: []string{"f"},
			Usage:   "Format Docker file",
			Action: func(c *cli.Context) {
				var rootNode = readInputFile(c)
				writeOutputFile(c, rootNode)
			},
			Flags: []cli.Flag{
				cli.StringFlag{Name: "input"},
				cli.StringFlag{Name: "output"},
			},
		},
		{
			Name:    "set-base-image-tag",
			Aliases: []string{"s"},
			Usage:   "Sets docker file's base image tag",
			Action: func(c *cli.Context) {
				var versionNumber = c.String("version")
				var imageName = c.String("image")
				var rootNode = readInputFile(c)
				setBaseImageVersion(rootNode, imageName, versionNumber)
				writeOutputFile(c, rootNode)
			},
			Flags: []cli.Flag{
				cli.StringFlag{Name: "input"},
				cli.StringFlag{Name: "output"},
				cli.StringFlag{Name: "version"},
				cli.StringFlag{Name: "image"},
			},
		},
	}
}

func main() {
	info()
	commands()
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
