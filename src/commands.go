package main

import (
	"fmt"
	"github.com/docker/distribution/reference"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"os"
)

func NodeToString(node *parser.Node) string {
	str := ""
	str += node.Value

	if len(node.Flags) > 0 {
		str += fmt.Sprintf(" %q", node.Flags)
	}

	for _, n := range node.Children {
		str += NodeToString(n) + "\n"
	}

	for n := node.Next; n != nil; n = n.Next {
		if len(n.Children) > 0 {
			str += " " + NodeToString(n)
		} else {
			str += " " + n.Value
		}
	}

	return str
}

func setBaseImageVersion(node *parser.Node, imageName string, versionNumber string, outfile string) {

	setBaseImageAstVersion(node, imageName, versionNumber)

	writer, err3 := os.Create(outfile)
	if err3 != nil {
		panic(err3)
	}

	defer writer.Close()

	writer.WriteString(NodeToString(node))
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
