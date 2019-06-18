package main

import (
	"os"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/urfave/cli"
)

var app = cli.NewApp()

func info() {
	app.Name = "lacuna"
	app.Usage = "Docker base image version setter CLI"
	app.Author = "Tigran TIKSN Torosyan"
	app.Version = "1.0.0"
}

func readInputFile(c *cli.Context) *parser.Node {
	var inputFile = c.String("input")

	var reader, err = os.Open(inputFile)
	if err != nil {
		panic(err)
	}

	var result, err2 = parser.Parse(reader)
	if err2 != nil {
		panic(err2)
	}

	return result.AST
}

func commands() {
	app.Commands = []cli.Command{
		{
			Name:    "set-base-image-tag",
			Aliases: []string{"s"},
			Usage:   "Sets base image tag",
			Action: func(c *cli.Context) {
				var versionNumber = c.String("version")
				var outputFile = c.String("output")
				var imageName = c.String("image")
				var rootNode = readInputFile(c)
				setBaseImageVersion(rootNode, imageName, versionNumber, outputFile)
			},
			Flags: []cli.Flag{
				cli.StringFlag{Name: "input"},
				cli.StringFlag{Name: "version"},
				cli.StringFlag{Name: "output"},
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
