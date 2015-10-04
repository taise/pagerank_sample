package main

import (
	"./graph"
	"fmt"
	"math"
)

/*
 * ## グラフの基礎知識
 *
 * - ノード(node): ネットワークの頂点
 * - リンク(link): ネットワークの頂点と頂点を結ぶ線
 * - 有向グラフ(directed graph): リンクのつながりに方向があるグラフ
 *     - ex) Twitterのフォロー関係
 * - 無向グラフ(undirected graph): リンクのつながりに方向がないグラフ
 *     - ex) Facebookの友人関係
 * - グラフ構造の値: 各ノードとリンクは値を持つことができる
 *     - ex) 地点Aの人口
 *     - ex) 地点Aから地点Bまでの距離
 *
 *
 * ex) 3つのノードを持つ有向グラフ
 *
 *  <------------,
 * 1<----,       |
 * └─ -> 2 <--> 3
 *
 * - nodes: {1, 2, 3}
 * - links: { 1: {2}, 2: {1, 3}, 3: {1, 2}}
 * - adjacencyList(matrix):
 *  {1: {1: 0,   2: 1,   3: 0},
 *   2: {1: 0.5, 2: 0,   3: 0.5}
 *   3: {1: 0.5, 2: 0.5, 3: 0}}
 *
 * |  |1  |2  |3  |
 * |:-|--:|--:|--:|
 * |1 |  0|  1|  0|
 * |2 |1/2|  0|1/2|
 * |3 |1/2|1/2|  0|
 *
 */

func main() {
	links := graph.GenerateLinks(20)
	links.Print()

	// 有向グラフのノードとそのRank
	//	 ex) map[int]float64 {nodeId: rank}
	nodes := links.ToNodes()

	// 各ノードからノードへの遷移確率をもつ隣接行列G
	G := meyersAdjacencyList(links, nodes)

	// 50 step以内にグラフのrankが収束するはず
	for step := 1; step <= 50; step++ {
		fmt.Printf("\n===== step %v =====\n", step)
		nodes = updateRank(nodes, G)
		printRank(nodes)
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
func probabilityAdjacencyList(links graph.Links, nodes graph.Nodes) graph.AdjacencyList {
	matrix := graph.BuildAdjacencyList(links, nodes)

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
func meyersAdjacencyList(links graph.Links, nodes graph.Nodes) graph.AdjacencyList {
	const d = float64(0.85)

	S := probabilityAdjacencyList(links, nodes)
	G := graph.BuildAdjacencyList(links, nodes)
	n := float64(len(G))

	for i, cols := range S {
		for j, _ := range cols {
			s := S[i][j]
			G[i][j] = d*(s) + (1.0-d)*(1.0/n)
		}
	}
	return G
}

// ノードからノードへ遷移確率pの分のrankをそれぞれ配分する
func updateRank(nodes graph.Nodes, G graph.AdjacencyList) graph.Nodes {
	nextNodes := nodes.CopyKey()

	for i, rank := range nodes {
		for j, p := range G[i] {
			nextNodes[j] += p * rank
		}
	}

	return nextNodes
}

func printRank(nodes graph.Nodes) {
	for id, rank := range nodes {
		fmt.Println("id: ", id, ",rank: ", rank)
	}
}
