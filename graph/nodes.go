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

func (self Nodes) SplitByNodeIds(d int) []Nodes {
	nodeSlice := []Nodes{}
	for i := 0; i < d; i++ {
		nodeSlice = append(nodeSlice, Nodes{})
	}

	for k, v := range self {
		fmt.Printf("i: %v, k: %v, v: %v\n", k%d, k, v)
		nodeSlice[k%d][k] = v
	}
	return nodeSlice
}
