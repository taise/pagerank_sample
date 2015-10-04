package graph

import (
	"fmt"
	"math/rand"
)

/* 有向グラフのリンク
 *    ex) map[int][]int {fromNodeId: {toNodeIds}}
 *
 * links := Links{
 * 	1: {2, 3},
 * 	2: {3},
 * 	3: {1, 4},
 * 	4: {1, 2, 3, 4},
 * }
 */

type Links map[int][]int

func (self Links) Print() {
	for k, links := range self {
		fmt.Println("key: ", k, "links: ", links)
	}
}

func (self Links) ToNodes() Nodes {
	const v = float64(1.0)
	nodes := Nodes{}

	for fromId, toIds := range self {
		nodes[fromId] = v
		for _, id := range toIds {
			nodes[id] = v
		}
	}
	return nodes
}

func shuffle(a []int) {
	for i := range a {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
}

func GenerateLinks(linkSize int) Links {
	links := Links{}

	for i := 0; i < linkSize; i++ {
		out := []int{}

		linkIds := []int{}
		for i := 0; i < linkSize; i++ {
			linkIds = append(linkIds, i)
		}
		shuffle(linkIds)

		outSize := rand.Intn(linkSize / 10)
		for j := 0; j < outSize; j++ {
			out = append(out, linkIds[j])
		}
		links[i] = out
	}
	return links
}
