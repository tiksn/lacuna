package main

import (
	"os"

	"github.com/docker/distribution/reference"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

func setBaseImageVersion(infile string, imageName string, versionNumber string, outfile string) {

	var reader, err = os.Open(infile)
	if err != nil {
		panic(err)
	}

	var result, err2 = parser.Parse(reader)
	if err2 != nil {
		panic(err2)
	}

	setBaseImageAstVersion(result.AST, imageName, versionNumber)

	writer, err3 := os.Create(outfile)
	if err3 != nil {
		panic(err3)
	}

	defer writer.Close()

	writer.WriteString(result.AST.Dump())
	writer.Sync()
}

func setBaseImageAstVersion(node *parser.Node, imageName string, versionNumber string) {

	for _, c := range node.Children {
		setBaseImageAstVersion(c, imageName, versionNumber)
	}

	if node.Value == "from" {
		var imageTagNode = node.Next
		var r, err = reference.Parse(imageTagNode.Value)
		if err != nil {
			panic(err)
		}

		if nt, isTagged := r.(reference.NamedTagged); isTagged {
			//TODO: Check old tag log.Fatalln(nt.Tag())
			setBaseImageTagVersion(imageTagNode, nt, versionNumber)
		} else {
			panic("Can't extract tags.")
		}
	}

}

func setBaseImageTagVersion(imageTagNode *parser.Node, fromNode reference.NamedTagged, versionNumber string) {
	var changedReference, err = reference.WithTag(fromNode, versionNumber)
	if err != nil {
		panic(err)
	}
	imageTagNode.Value = changedReference.String()
}
