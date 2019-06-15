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

		if nt, isTagged := r.(reference.NamedTagged); isTagged {
			//TODO: Check old tag log.Fatalln(nt.Tag())
			setBaseImageTagVersion(imageTagNode, nt, versionNumber)
			log.Fatalln(nt)
		}

		log.Fatalln("Can't extract tags.")
	}

}

func setBaseImageTagVersion(imageTagNode *parser.Node, fromNode reference.NamedTagged, versionNumber string) {
	var changedReference, err = reference.WithTag(fromNode, versionNumber)
	if err != nil {
		log.Fatalln(err)
	}
	imageTagNode.Value = changedReference.String()
	log.Fatalln(changedReference)
}