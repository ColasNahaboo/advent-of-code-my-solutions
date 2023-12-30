// Copyright 2013 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package astar implements the A* search algorithm for finding least-cost paths.

// Colas: This is a modified version of https://github.com/fzipp/astar
// modified to not requiring playing with types and methods, and working
// with graphs that are dynamically lazy-created on demand, useful to explore
// solutions in potentially huge spaces, and the goal to reach can be many states
// satisfying an end condition.
// Typically, finding a path in all possible moves of a chess game, with the
// destination being not one configuration but any one providing a mat
//
// For now, this is not a package but a simple file to copy into your sources
//
// The recommended use is to have Nodes be an id (int), an index in a dynamic
// slice (global variable) of more complex structures states that can be
// created on demand and have a way (a Map?) to find the neigbours of a state,
// auto-creating them on demand. But Node can be any comparable Go type.
// Graph is your "context" type holding all you global info (such as the dynamic
// slice above with Nodes as indexes)
//
// User thus just has to define 4 functions to pass to AStarFindPath:
// - nodesConnected(g any, node Node) []Node, type ConnectedFunc
//   returns a slice of the neighbour nodes, possibly auto-created
// - nodesDistance(g any, node1, node2 Node) float64, type CostFunc
//   returns the length of the path from node1 to neighbour node2
// - nodesHeuristic(g any, node, destination Node) float64, type CostFunc
//   returns how far from the destination we estimate node is
// - nodeEnd(g any, node Node) bool, type EndFunc
//   tells if we have reached a satisfying destination node
// Then the call to the generic AStarFindPath:
// AStarFindPath[Node](nil, start, dest, nodesConnected, nodesDistance, nodesHeuristic, nodeEnd)
// finds the shortest path (a []Node) via the A* algorithm

package main

import "container/heap"

//////////// Types of the Callbacks implemented and provided to AStarFindPath
// A ConnectedFunc returns the neighbour nodes of node n in the graph.
type ConnectedFunc[Graph, Node any] func(g Graph, a Node) []Node
// A CostFunc returns a cost for the transition node a -> node b
type CostFunc[Graph, Node any] func(g Graph, a, b Node) float64
// A EndFunc returns true if node is a/the destination
type EndFunc[Graph, Node any] func(g Graph, a, end Node) bool

//////////// The path-finding function, AStarFindPath

// AStarFindPath finds the least-cost path between start and dest in graph g
// using the cost function d and the cost heuristic function h.
// g can be anything, it is just a context blindly passed to the callback
// functions in argument. This can be either the full graph structure, or just
// nil in trivial programs relying on global variables for the graph and context
// Note that g should be a pointer to an object if the graph is to be
// dynamically created by the calls to your callbacks.
// Returns the shortest path as a list (a slice) of nodes

func AStarFindPath[Graph any, Node comparable](g Graph, start, dest Node, cf ConnectedFunc[Graph, Node], df, hf CostFunc[Graph, Node], ne EndFunc[Graph, Node]) Path[Node] {
	closed := make(map[Node]bool)

	pq := &priorityQueue[Path[Node]]{}
	heap.Init(pq)
	heap.Push(pq, &item[Path[Node]]{value: newPath(start)})

	for pq.Len() > 0 {
		p := heap.Pop(pq).(*item[Path[Node]]).value
		n := p.last()
		if closed[n] {
			continue
		}
		if ne(g, n, dest) {			// Path found
			return p
		}
		closed[n] = true

		for _, nb := range cf(g, n) {
			cp := p.cont(nb)
			heap.Push(pq, &item[Path[Node]]{
				value:    cp,
				priority: -(pathCost[Graph, Node](g, cp, df) + hf(g, nb, dest)),
			})
		}
	}

	// No path found
	return nil
}


// A Path is a sequence of nodes in a graph.
type Path[Node any] []Node

// newPath creates a new path with one start node. More nodes can be
// added with append().
func newPath[Node any](start Node) Path[Node] {
	return []Node{start}
}

// last returns the last node of path p. It is not removed from the path.
func (p Path[Node]) last() Node {
	return p[len(p)-1]
}

// cont creates a new path, which is a continuation of path p with the
// additional node n.
func (p Path[Node]) cont(n Node) Path[Node] {
	cp := make([]Node, len(p), len(p)+1)
	copy(cp, p)
	cp = append(cp, n)
	return cp
}

// Cost calculates the total cost of path p by applying the cost function d
// to all path segments and returning the sum.
func pathCost[Graph any, Node any](g Graph, p []Node, df CostFunc[Graph, Node]) (c float64) {
	for i := 1; i < len(p); i++ {
		c += df(g, p[i-1], p[i])
	}
	return c
}

//////////////////////// Priority Queues ////////////////////////

// An item is something we manage in a priority queue.
type item[T any] struct {
	value    T       // The value of the item; arbitrary.
	priority float64 // The priority of the item in the queue.
}

// A priorityQueue implements heap.Interface and holds items.
type priorityQueue[T any] []*item[T]

func (pq priorityQueue[T]) Len() int { return len(pq) }

func (pq priorityQueue[T]) Less(i, j int) bool {
	return pq[i].priority > pq[j].priority
}

func (pq priorityQueue[T]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *priorityQueue[T]) Push(x any) {
	*pq = append(*pq, x.(*item[T]))
}

func (pq *priorityQueue[T]) Pop() any {
	old := *pq
	n := len(old)
	it := old[n-1]
	*pq = old[0 : n-1]
	return it
}
