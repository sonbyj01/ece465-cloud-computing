package main

import (
	"fmt"
	"proj1/graph"
	"sync"
)

//func printGraph(g graph.Graph) {
//	fmt.Println("===============")
//
//	for index, node := range g.Nodes {
//		fmt.Println("-- DEBUGGING-- ", "Index: ", index,
//			"Node Value: ", node.Value, "Edge Array: ", node.Adj)
//	}
//
//	fmt.Println("===============")
//}
//
//func printNodes(n []graph.Node) {
//	fmt.Println("===============")
//
//	for index, node := range n {
//		fmt.Println("-- DEBUGGING-- ", "Index: ", index,
//			"Node Value: ", node.Value, "Edge Array: ", node.Adj)
//	}
//
//	fmt.Println("===============")
//}

//func speculativeColoring(u *graph.Graph) {
//	var wg sync.WaitGroup
//
//	for index, v := range u.Nodes {
//		wg.Add(1)
//
//		go func() {
//			defer wg.Done()
//			colors := make([]bool, 20)
//
//			// marks color w as forbidden to v
//			for _, w := range v.Adj {
//				colors[w.Value] = true
//			}
//
//			// assigns smallest available value to color v
//			for minColorValue, boolValue := range colors {
//				if !boolValue {
//					u.Nodes[index].Value = minColorValue
//					break
//				}
//			}
//		}()
//	}
//	wg.Wait()
//}
//
//func conflictResolution(u graph.Graph) graph.Graph {
//	var wg sync.WaitGroup
//	var r graph.Graph
//
//	for index, v := range u.Nodes {
//		wg.Add(1)
//
//		go func() {
//			defer wg.Done()
//
//			for _, w := range v.Adj {
//				//fmt.Println("v", v.Value)
//				//fmt.Println("w", w.Value)
//				if v.Value == w.Value {
//					fmt.Println("conflict")
//					if v.Index > w.Index {
//						fmt.Println("appended")
//						r.Nodes = append(r.Nodes, u.Nodes[index])
//					}
//				}
//			}
//		}()
//	}
//	wg.Wait()
//	return r
//}

func assignColors(G graph.Graph, C []int, Conf []graph.Node) []int {
	var wg sync.WaitGroup

	for index, v := range Conf {
		wg.Add(1)
		fmt.Println(index, v)

		go func() {
			defer wg.Done()
			Forbidden := make([]bool, len(Conf))
			fmt.Println("before: ", v.Index, Forbidden)
			for _, u := range Conf[index].Adj {
				Forbidden[u.Value] = true
			}
			fmt.Println("after: ", v.Index, Forbidden)

			for minColorVal, boolVal := range Forbidden {
				if !boolVal {
					C[v.Index] = minColorVal
					G.Nodes[v.Index].Value = minColorVal
					Conf[v.Index].Value = minColorVal
					break
				}
			}
		}()
	}
	wg.Wait()
	return C
}

func detectConflicts(G graph.Graph, C []int, Conf []graph.Node) []graph.Node {
	var NewConf []graph.Node
	ch := make(chan []graph.Node)

	for _, v := range Conf {
		go func() {
			for _, u := range v.Adj {
				if C[v.Index] == C[u.Index] && u.Index < v.Index {
					ch <- append(NewConf, v)
				}
			}
		}()
	}
	NewConfP := <-ch
	return NewConfP
}

// C is a color array associated with the respective node index
// Conf is all the nodes of graph G
func main() {
	G := graph.NewCompleteGraph(4)
	C := make([]int, len(G.Nodes))
	Conf := G.Nodes

	for len(Conf) != 0 {
		//fmt.Println(C)
		//printNodes(Conf)
		C = assignColors(G, C, Conf)
		//fmt.Println(C)
		//printNodes(Conf)
		break
		//Conf = detectConflicts(G, C, Conf)
		//time.Sleep(10 * time.Millisecond)
	}
}
