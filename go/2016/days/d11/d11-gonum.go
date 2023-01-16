// Adventofcode 2016, d11, in go.
// variant of part2, implemented with the gonum/path package instead of fzipp/astar

package main

import (
	"fmt"
    // https://www.gonum.org/ - https://pkg.go.dev/gonum.org/v1/gonum/graph/path#AStar
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
	"time"
)

var tos map[int64][]int64
var statesids map[int64]State		// the cache of dynamically instancied states
var stateslist States

func part3(lines []string) int {
	lines[0] = lines[0] + " And also a elerium generator, a elerium-compatible microchip, a dilithium generator, a dilithium-compatible microchip."
	startid := parse(lines)
	graph := graphBuild(startid)
	start := states[startid]
	VPf("Looking for path to %v\n", goal)
	startTime := time.Now().UnixMilli()
	path, expanded := path.AStar(start, goal, graph, heuristic)
	endTime := time.Now().UnixMilli()
	steps, weight := path.To(s2iID(goal.id))
	fmt.Printf("Number of explored nodes: %d, expanded: %d, path weight: %v, in %.0fms\n", nodes, expanded, weight, float64(endTime - startTime))
	return len(steps) - 1
}

// gonum demands that we first build the graph totally
func graphBuild(sid string) *simple.DirectedGraph {
	tos = make(map[int64][]int64, 0)
	statesids = make(map[int64]State, 0)
	fillMaps(sid)
	g := simple.NewDirectedGraph()
	for id, tolist := range tos {
		for _, to := range tolist {
			g.SetEdge(simple.Edge{statesids[id], statesids[to]})
		}
	}
	return g
}

// fills our cache
func fillMaps(sid1 string) {
	id1 := s2iID(sid1)
	statesids[id1] = states[sid1]
	for _, sid2 := range stateslist.Neighbours(sid1) {
		id2 := s2iID(sid2)
		tos[id1] = append(tos[id1], id2)
		if _, ok := tos[id2]; !ok {
			fillMaps(sid2)
		}
	}
}

//////// Interfaces for gonum/path/AStar
// we define here methods required to use gonum/path/AStar

// see https://pkg.go.dev/gonum.org/v1/gonum/graph#Node
func (s State)ID() int64 {
	return s2iID(s.id)
}

// see https://pkg.go.dev/gonum.org/v1/gonum/graph/path#Heuristic
func heuristic(id1, id2 graph.Node) float64 {
	return stateDist(i2sID(id1.ID()), i2sID(id2.ID()))
}


//////// Utilities

// convert between the int64 IDs for gonum/graph and string IDs used with fzipp/astar
func s2iID(sid string) (id int64) {
	id = int64(sid[0] - '0')				// E floor
	for i :=0; i < nmetals; i++ {
		id = id*16 + int64(sid[i+1] - '0')
	}
	return
}

func i2sID(id int64) (sid string) {
	for i := nmetals; i >=0; i-- {
		sid += string('0' + (id >> (4*i)) & 15)
	}
	return
}
