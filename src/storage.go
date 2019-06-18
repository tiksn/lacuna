package main

import (
	"os"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/urfave/cli"
)

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
