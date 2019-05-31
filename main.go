package main

import (
	"os"
	"log"
	"github.com/urfave/cli"
	// "github.com/moby/buildkit/frontend/dockerfile/parser"
)

var app = cli.NewApp()

func main() {
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
