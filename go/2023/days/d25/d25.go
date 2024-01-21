// Adventofcode 2023, d25, in go. https://adventofcode.com/2023/day/25
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 54
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"flag"
	"fmt"
	"regexp"
	"math/rand"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/emirpasic/gods/queues/arrayqueue"
)

var verbose bool

type Link struct {
	i, j int					// ids of comps, in order
}

var compname = []string{}			// components id -> name
var compid = map[string]int{}		// name -> id
var conns = [][]int{}				// id -> list of connected ids
var links = []Link{}				// list of links (from, to), from < to

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	flag.Parse()
	verbose = *verboseFlag
	var infile string
	if flag.NArg() > 0 {
		infile = flag.Arg(0)
	} else {
		infile = fileMatch("input,[0-9]*,[0-9]*.test")
	}
	lines := fileToLines(infile)

	var result int
	if *partOne {
		VP("Running Part1")
		result = part1(lines)
	} else {
		VP("Running Part2")
		result = part2(lines)
	}
	fmt.Println(result)
}

//////////// Part 1

func part1(lines []string) int {
	graph := parse(lines)
	return KargerMinCutSolve(graph)
}

//////////// Part 2


func part2(lines []string) int {
	for _, line := range lines {
		fmt.Println(line)
	}
	return 0
}

//////////// Common Parts code

func parse(lines []string) [][]int {
	re := regexp.MustCompile("[[:alpha:]]+")
	for _, line := range lines {
		names := re.FindAllString(line, -1)
		ids := make([]int, len(names))
		for i, name := range names {
			id, ok := compid[name]
			if ! ok {
				id = len(compname)
				compname = append(compname, name)
				compid[name] = id
				conns = append(conns, []int{})
			}
			ids[i] = id
		}
		conns[ids[0]] = append([]int{}, ids[1:]...)
	}
	// check for double links
	for i, cos := range conns {
		for _, j := range cos {
			if i != j && indexOfInt(conns[j], i) != -1 {
				panic(fmt.Sprintf("Dual link found: %s / %s", compname[i], compname[j]))
			}
		}
	}
	return conns
}

var id int

func KargerMinCutSolve(graph [][]int) int {
	adj, edges := make(map[int]mapset.Set[int]), make([]Edge, 0)
	// build our ad hoc  edges data
	for src, dests := range graph {
		for _, d := range dests {
			edge := addLink(adj, src, d)
			edges = append(edges, edge)
		}
	}
	mincut, vertices := mapset.NewSet[Edge](), len(graph)
	for mincut.Cardinality() != 3 {
		mincut = KargerMinCut(edges, vertices)
	}
	for rem := range mincut.Iter() {
		adj[rem.i].Remove(rem.j)
		adj[rem.j].Remove(rem.i)
	}
	queue, visited := arrayqueue.New(), mapset.NewSet[int]()
	queue.Enqueue(0)
	visited.Add(0)
	for !queue.Empty() {
		vertex, _ := queue.Dequeue()
		for n := range adj[vertex.(int)].Iter() {
			if !visited.Contains(n) {
				visited.Add(n)
				queue.Enqueue(n)
			}
		}
	}
	return visited.Cardinality() * (vertices - visited.Cardinality())
}

func addLink(adj map[int]mapset.Set[int], src, dest int) Edge {
	if _, hs := adj[src]; !hs {
		adj[src] = mapset.NewSet[int]()
	}
	if _, hd := adj[dest]; !hd {
		adj[dest] = mapset.NewSet[int]()
	}
	adj[src].Add(dest)
	adj[dest].Add(src)
	return Edge{src, dest}
}

//////////// Karger algo

type Edge struct {
	i int
	j int
}

func KargerMinCut(edges []Edge, vertices int) mapset.Set[Edge] {
	dsu, v, e := Initialize(vertices), vertices, len(edges)

	for v > 2 {
		// Getgting a random integer in the range [0, e-1].
		i := rand.Intn(e)
		set1, set2 := dsu.find(edges[i].i), dsu.find(edges[i].j)
		if set1 != set2 {
			dsu.union(edges[i].i, edges[i].j)
			v--
		}
	}
	cutset := mapset.NewSet[Edge]()
	for _, edge := range edges {
		set1, set2 := dsu.find(edge.i), dsu.find(edge.j)
		if set1 != set2 {
			cutset.Add(edge)
		}
	}
	return cutset
}

//////////// union-find via disjoint sets trees
// from https://github.com/vipul0092/advent-of-code-2023/blob/main/utils/disjointset.go

type Arr []int

type DisjointSet struct {
	parent Arr
	size   Arr
}

func Initialize(sz int) *DisjointSet {
	parent, size := make(Arr, sz), make(Arr, sz)
	dsu := DisjointSet{parent, size}
	for i := 0; i < sz; i++ {
		parent[i], size[i] = i, 1
	}
	return &dsu
}

func (dsu *DisjointSet) find(item int) int {
	if item == dsu.parent[item] {
		return item
	}
	res := dsu.find(dsu.parent[item])
	dsu.parent[item] = res
	return res
}

func (dsu *DisjointSet) union(v1, v2 int) {
	p1, p2 := dsu.find(v1), dsu.find(v2)
	if p1 == p2 {
		return
	}
	size1, size2 := dsu.size[v1], dsu.size[v2]
	// Make the item having higher tree size the parent of the one having lower size
	if size1 <= size2 {
		dsu.parent[p1] = p2
		dsu.size[p2] += dsu.size[p1]
	} else {
		dsu.parent[p2] = p1
		dsu.size[p1] += dsu.size[p2]
	}
}

//////////// PrettyPrinting & Debugging functions

func VPgroup(label string, group []int) {
	if ! verbose { return }
	VP(label + ":")
	if group == nil {
		group = make([]int, len(compname))
	}
	for i, c := range conns {
		VPf("  (%d)%s {%d} %v\n", i, compname[i], group[i], c)
	}
}

func comboString(combo []int) (s string) {
	for _, lid := range combo {
		link := links[lid]
		s = s + fmt.Sprintf(" (%d)%s/(%d)%s", link.i, compname[link.i], link.j, compname[link.j])
	}
	return
}
