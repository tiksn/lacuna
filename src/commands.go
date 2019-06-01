package main

import (
	"log"
	"os"

	"github.com/docker/distribution/reference"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

func setBaseImageVersion(infile string, imageName string, versionNumber string, outfile string) {

	var reader, err = os.Open(infile)
	if err != nil {
		log.Fatalln(err)
	}

	var result, err2 = parser.Parse(reader)
	if err2 != nil {
		log.Fatalln(err2)
	}

	setBaseImageAstVersion(result.AST, imageName, versionNumber)
}

func setBaseImageAstVersion(node *parser.Node, imageName string, versionNumber string) {

	for _, c := range node.Children {
		setBaseImageAstVersion(c, imageName, versionNumber)
	}

	if node.Value == "from" {
		var imageTagNode = node.Next
		var r, err = reference.Parse(imageTagNode.Value)
		if err != nil {
			log.Fatalln(err)
		}

		log.Fatalln(r.String())
	}

}
