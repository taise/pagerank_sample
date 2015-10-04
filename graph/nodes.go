package graph

import "fmt"

type Nodes map[int]float64

func (self Nodes) CopyKey() Nodes {
	newnodes := Nodes{}
	for k, _ := range self {
		newnodes[k] = float64(0)
	}
	return newnodes
}

func (self Nodes) Print() {
	for id, rank := range self {
		fmt.Println("id: ", id, ",rank: ", rank)
	}
}
