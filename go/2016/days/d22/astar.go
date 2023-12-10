// Copyright 2013 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package astar implements the A* search algorithm for finding least-cost paths.

// Colas: This is a modified version of https://github.com/fzipp/astar
// with the end condition of FindPath being not that n is equal to dest, but
// that the distance from n to dest is <= 0
// This allows to search for multiple possible destinations
//
// The recommended use is to have Nodes be an id (int), an index in a dynamic
// slice (global variable) of more complex structures states that can be
// created on demand and have way (Map table?) to find the neigbours of a state,
// auto-creating them on demand.
//
// User thus just has to define 3 functions to pass to FindPath:
// - nodesConnected(node Node) []Node, type ConnectedFunc
// - nodesDistance(node1, node2 Node) float64, type CostFunc
//   a negative or null result meants the destination is reached
// - nodesHeuristic(node, destination Node) float64, type CostFunc
// Then the call to the generic FindPath:
// FindPath[Node](start, dest, nodesConnected, nodesDistance, nodesHeuristic)
// finds the shortest path (a []Node) via the A* algorithm

package main

import "container/heap"

// A ConnectedFunc returns the neighbour nodes of node n in the graph.
type ConnectedFunc[Node any] func(a Node) []Node

// A CostFunc is a function that returns a cost for the transition
// from node a to node b.
type CostFunc[Node any] func(a, b Node) float64

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
func (p Path[Node]) Cost(d CostFunc[Node]) (c float64) {
	for i := 1; i < len(p); i++ {
		c += d(p[i-1], p[i])
	}
	return c
}

// FindPath finds the least-cost path between start and dest in graph g
// using the cost function d and the cost heuristic function h.
func FindPath[Node comparable](start, dest Node, c ConnectedFunc[Node], d, h CostFunc[Node]) Path[Node] {
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
		if h(n, dest) <= 0 {
			// Path found
			return p
		}
		closed[n] = true

		for _, nb := range c(n) {
			cp := p.cont(nb)
			heap.Push(pq, &item[Path[Node]]{
				value:    cp,
				priority: -(cp.Cost(d) + h(nb, dest)),
			})
		}
	}

	// No path found
	return nil
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
