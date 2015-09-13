package main

import (
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
 *
 */

func main() {
	// rankの初期値
	const v = float64(1.0)

	// 有向グラフのリンク
	//   ex) map[int][]int {from_node_id: {to_node_ids}}
	links := map[int][]int{
		1: {2, 3, 4},
		2: {1, 4, 5},
		3: {1, 4},
		4: {2, 3},
		5: {2, 3, 6},
		6: {1, 4, 7},
	}
	printLinks(links)

	// 有向グラフのノード
	//	 ex) map[int]float64 {node_id: rank}
	nodes := generateNodeFrom(links, v)

	// 50 step以内にグラフのrankが収束するはず
	for step := 1; step < 50; step++ {
		fmt.Printf("\n===== step %v =====\n", step)
		nodes = updateRank(nodes, links)
	}
}

func Round(f float64) float64 {
	const place = 15
	shift := math.Pow(10, float64(place))
	return math.Floor(f*shift+.5) / shift
}

// リンクからノードを生成する
func generateNodeFrom(links map[int][]int, v float64) map[int]float64 {
	nodes := map[int]float64{}

	for from_id, to_ids := range links {
		nodes[from_id] = v
		for _, id := range to_ids {
			nodes[id] = v
		}
	}
	return nodes
}

func copyNodeKey(nodes map[int]float64) map[int]float64 {
	newnodes := map[int]float64{}
	for k, _ := range nodes {
		newnodes[k] = float64(0)
	}
	return newnodes
}

func printLinks(links map[int][]int) {
	for k, links := range links {
		fmt.Println("key: ", k, "links: ", links)
	}
}

func printRank(nodes map[int]float64) {
	for id, rank := range nodes {
		fmt.Println("id: ", id, ",rank: ", rank)
	}
}

// p: ノードから他のノードへの遷移確率
//    リンクの重みがなければ、1 / 他のノードへのリンク数
func p(link []int) float64 {
	return 1 / float64(len(link))
}

/*
 * ## 更新手順
 *
 * 1. ノードを1つ取り出す
 * 2. 各ノードから他のノードへのリンクを元にpを算出する
 * 3. ノードを持つrankとpを掛け、他のノードにrankを配分する
 * 4. 他のノードに配分していないノードがあれば 1. に戻る
 *
 */
func updateRank(nodes map[int]float64, links map[int][]int) map[int]float64 {
	nextNodes := copyNodeKey(nodes)
	for id, rank := range nodes {
		shareRank := Round(p(links[id]) * rank)

		for _, target_id := range links[id] {
			nextNodes[target_id] += shareRank
		}
	}

	printRank(nextNodes)
	return nextNodes
}
