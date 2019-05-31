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
				setBaseImageVersion(c.Args().Get(0), c.Args().Get(1), c.Args().Get(2))
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