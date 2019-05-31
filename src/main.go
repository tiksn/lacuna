package main

import (
	"fmt"
	"github.com/urfave/cli"
	"log"
	"os"
	// "github.com/moby/buildkit/frontend/dockerfile/parser"
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
			Name:    "set-base-image-tag",
			Aliases: []string{"s"},
			Usage:   "Sets base image tag",
			Action: func(c *cli.Context) {
				var inputFile = c.String("input")
				var versionNumber = c.String("input")
				var outputFile = c.String("output")
				setBaseImageVersion(inputFile, versionNumber, outputFile)
			},
			Flags: []cli.Flag{
				cli.StringFlag{Name: "input"},
				cli.StringFlag{Name: "version"},
				cli.StringFlag{Name: "output"},
			},
		},
	}
}

func main() {
	info()
	commands()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func setBaseImageVersion(infile string, version string, outfile string) {
	fmt.Printf("not implemented yet. Versions is %q %q %q", infile, version, outfile)
}
