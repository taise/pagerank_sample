package pagerank

import (
	. "../graph"
	"fmt"
	"math"
)

func Run(links Links) {
	// 並列数
	maxRoutineSize := 4

	// 有向グラフのノードとそのRank
	//	 ex) map[int]float64 {nodeId: rank}
	nodes := links.ToNodes()

	// 各ノードからノードへの遷移確率をもつ隣接行列G
	G := meyersAdjacencyList(links, nodes)

	// 50 step以内にグラフのrankが収束するはず
	// for step := 1; step <= 50; step++ {
	for step := 1; step <= 2; step++ {
		fmt.Printf("\n===== step %v =====\n", step)
		nodes = updateRank(nodes, G, maxRoutineSize)
		nodes.Print()
	}
}

func Round(f float64) float64 {
	const place = 15
	shift := math.Pow(10, float64(place))
	return math.Floor(f*shift+.5) / shift
}

/*
 * 遷移確率の隣接行列S
 *   リンクがないノードjへの遷移確率は0とする
 *   ai = {0: リンクあり, 1: リンクなし}
 *
 *   Niはノードiのリンク数
 *   nはノード数
 *
 *   あるノードiからあるノードjへの遷移確率pij
 *   pij = (1/Ni) + ai(1/n)
 *
 * |  |1  |2  |3  |4  |5  |6  |
 * |:-|--:|--:|--:|--:|--:|--:|
 * |1 |  0|1/2|1/2|  0|  0|  0|
 * |2 |  0|  0|  1|  0|  0|  0|
 * |3 |1/3|  0|  0|1/3|1/3|  0|
 * |4 |1/6|1/6|1/6|1/6|1/6|1/6|
 * |5 |  0|  0|  0|  0|  0|  1|
 * |6 |  0|  0|  0|  0|  1|  0|
 */
func probabilityAdjacencyList(links Links, nodes Nodes) AdjacencyList {
	matrix := BuildAdjacencyList(links, nodes)

	for row, link := range links {
		ni := float64(len(link))
		for _, col := range link {
			matrix[row][col] = 1.0 / ni
		}
	}
	return matrix
}

/*
 * Meyer's Random Surfer Model
 *   damping factorを考慮した遷移確率
 *   d(damping factor)は、リンクがなくても一定確率で別ノードに遷移する確率
 *
 *   nはノード数
 *
 *   あるノードiからノードjへの遷移確率pij
 *   pij = d * Sij + (1 - d) * (1 / n)
 *
 * 以下は d = 0.9の場合の隣接行列
 * |  |1    |2    |3    |4    |5    |6    |
 * |:-|----:|----:|----:|----:|----:|----:|
 * |1 | 1/60|28/60|28/60| 1/60| 1/60| 1/60|
 * |2 | 1/60| 1/60|55/60| 1/60| 1/60| 1/60|
 * |3 |19/60| 1/60| 1/60|19/60|19/60| 1/60|
 * |4 |10/60|10/60|10/60|10/60|10/60|10/60|
 * |5 | 1/60| 1/60| 1/60| 1/60| 1/60|55/60|
 * |6 | 1/60| 1/60| 1/60| 1/60|55/60| 1/60|
 */
func meyersAdjacencyList(links Links, nodes Nodes) AdjacencyList {
	const d = float64(0.85)

	S := probabilityAdjacencyList(links, nodes)
	G := BuildAdjacencyList(links, nodes)
	n := float64(len(G))

	for i, cols := range S {
		// - リンクのあるノード
		//     - 確率`d`で、リンク先のノードに遷移する
		//     - 確率`1 - d`で、全ノードのなかからランダムで遷移する
		// - リンクのないノード(dangling node)
		//     - 確率`1 - d`で、全ノードのなかからランダムで遷移する
		p := d
		if links.IsDanglingNode(i) {
			p = 0
		}
		for j, _ := range cols {
			s := S[i][j]
			G[i][j] = p*(s) + (1.0-p)*(1.0/n)
		}
	}
	return G
}

// ノードからノードへ遷移確率pの分のrankをそれぞれ配分する
func divideRank(nodes Nodes, G AdjacencyList) MapResult {
	mapResult := MapResult{}

	for i, rank := range nodes {
		for j, p := range G[i] {
			mapResult[j] += p * rank
		}
	}
	return mapResult
}

// 遷移確率GとNodesのランク値を使って、次の時点のNodesのランク値を計算する
func updateRank(nodes Nodes, G AdjacencyList, maxRoutineSize int) Nodes {
	nextNodes := nodes.CopyKey()

	nodeMap := nodes.SplitByNodeIds(4)

	for _, nodes := range nodeMap {
		mapResult := divideRank(nodes, G)
		fmt.Print(mapResult)
	}

	//for i, rank := range nodes {
	//	for j, p := range G[i] {
	//		nextNodes[j] += p * rank
	//	}
	//}

	return nextNodes
}
