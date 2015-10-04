package graph

type AdjacencyList map[int]map[int]float64

func BuildAdjacencyList(links Links, nodes Nodes) AdjacencyList {
	adjacencyList := AdjacencyList{}

	// 有向グラフの場合は対称行列ではない
	// nodeIdの隣接行列を0で埋める
	for i, _ := range nodes {
		tmpMap := map[int]float64{}
		for j, _ := range nodes {
			tmpMap[j] = 0
		}
		adjacencyList[i] = tmpMap
	}
	return adjacencyList
}
