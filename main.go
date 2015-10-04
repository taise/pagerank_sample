package main

import (
	"./graph"
	"./pagerank"
)

func main() {
	links := graph.GenerateLinks(20)
	links.Print()

	pagerank.Run(links)
}
