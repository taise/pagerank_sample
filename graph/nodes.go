package graph

type Nodes map[int]float64

func (self Nodes) CopyKey() Nodes {
	newnodes := Nodes{}
	for k, _ := range self {
		newnodes[k] = float64(0)
	}
	return newnodes
}
