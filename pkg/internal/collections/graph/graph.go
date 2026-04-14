package graph

import (
	"github.com/4strodev/wiring_graphs/pkg/internal/collections/set"
)

type Graph[T any] struct {
	nodes set.Set[*Node[T]]
}

func NewGraph[T any]() Graph[T] {
	return Graph[T]{
		nodes: set.New[*Node[T]](),
	}
}

func (g *Graph[T]) Add(node *Node[T]) {
	g.nodes.Add(node)
}

func (g Graph[T]) Connect(originNode, destinationNode *Node[T], dir connectionDirection) {
	g.Add(originNode)
	g.Add(destinationNode)

	originNode.Connect(destinationNode, dir)
}

func (g Graph[T]) GetNodes() set.Set[*Node[T]] {
	return g.nodes
}

func (g Graph[T]) GetRootNodes() []*Node[T] {
	var rootNodes = []*Node[T]{}
	for node := range g.nodes {
		if !node.HasIncomingNodes() {
			rootNodes = append(rootNodes, node)
		}
	}

	return rootNodes
}

// DetectCircularRelations using DFS detects circular relations between nodes on this graph
func (g Graph[T]) DetectCircularRelations() ([]*Node[T], bool) {
	for node := range g.nodes {
		if !node.HasConnections() {
			continue
		}

		state := map[*Node[T]]dfsStates{} // 0: Unvisited, 1: Visiting, 2: Visited
		parent := map[*Node[T]]*Node[T]{}
		var cycle []*Node[T]
		found := dfsVisit(node, state, parent, &cycle)
		if found {
			return cycle, found
		}
	}

	return nil, false
}
