// Adventofcode 2024, d23, in go. https://adventofcode.com/2024/day/23
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 7
// TEST: example1 "aa_ab_ae"  // the example on wikipedia https://en.wikipedia.org/wiki/Bron%E2%80%93Kerbosch_algorithm#Example
// TEST: example "co_de_ka_ta"
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

// part2 is coded with my own set simple routines
// part3 is part2 coded with the bits-and-bloom bitset package
//       ... which is much faster!

package main

import (
	"fmt"
	"regexp"
	"flag"
	"sort"
	"github.com/bits-and-blooms/bitset"
)

//////////// Options parsing & exec parts
var commaFlag *bool
var commaSep = "_"

func main() {
	commaFlag = flag.Bool("c", false, "outputs numbers separated by comma instead of underscores")
	ExecOptionsString(2, NoXtraOpts, part1, part2, part3)
}

func XtraOpts() { // extra options, see ParseOptions in utils.go
	if *commaFlag {
		commaSep = ","
	}
}

//////////// Part 1

func part1(lines []string) string {
	res := 0
	conns := parse(lines)
	triplets := make(map[[3]int]bool)
	for _, conn := range conns {
		Find3rd(&triplets, conns, conn[0], conn[1])
		Find3rd(&triplets, conns, conn[1], conn[0])
	}
	for t := range triplets {
		if t[0] / 26 == 19 || t[1] / 26 == 19 || t[2] / 26 == 19 { // t-starts
			res++
			VPf("  %s,%s,%s\n", i2n(t[0]), i2n(t[1]), i2n(t[2]))
		}
	}
	return itoa(res)
}

func Find3rd(triplets *map[[3]int]bool, conns [][]int, c1, c2 int) {
	for _, c3 := range ConnectedTo(conns, c1) {
		if c3 == c2 { continue }
		if IsConnected(conns, c3, c2) {
			AddTriplet(triplets, c1, c2, c3)
		}
	}
}

func ConnectedTo(conns [][]int, to int) (cs []int) {
	for _, c := range conns {
		if c[0] == to {
			cs = append(cs, c[1])
		} else if c[1] == to {
			cs = append(cs, c[0])
		}
	}
	return
}

func IsConnected(conns [][]int, c0, c1 int) bool {
	for _, c := range conns {
		if (c[0] == c0 && c[1] == c1) || (c[0] == c1 && c[1] == c0) {
			return true
		}
	}
	return false
}


func AddTriplet(triplets *map[[3]int]bool, elts ...int) {
	sort.Slice(elts, func(i, j int) bool { return elts[i] < elts[j] })
	(*triplets)[[3]int{elts[0], elts[1], elts[2]}] = true
}

//////////// Part 2

func part2(lines []string) string {
	conns := parse(lines)
	edges := make([][]int, 26*26, 26*26) // [comp]-> list of neigbours
	for _, c := range conns {
		edges[c[0]] = append(edges[c[0]], c[1])
		edges[c[1]] = append(edges[c[1]], c[0])
	}
	return NodesPrint(BronKerboschEdges(edges))
}

//////////// Bron–Kerbosch algorithm using my sets
// With pivot version, from:
// https://en.wikipedia.org/wiki/Bron%E2%80%93Kerbosch_algorithm#With_pivoting

func BronKerboschEdges(edges [][]int) []int {
	bkedges = edges
	cliques := BronKerbosch(BKNew(), BKSetMake(edges), BKNew())
	maxlen := 0
	var maxclique []int
	for _, s := range cliques {
		nodes := s.Nodes()
		if len(nodes) > maxlen {
			maxlen = len(nodes)
			maxclique = nodes
		}
	}
	return maxclique
}

func BronKerbosch(r, p, x BKSet) (cliques []BKSet) {
	// if P and X are both empty then report R as a maximal clique
	if p.IsEmpty() && x.IsEmpty() {
		return []BKSet{r}
	}
	// else choose pivot vertex (with lots of edges) q from P ⋃ X:
	u := p.Union(x).Pivot()
	// for each vertex v in P \ N(u) do
	for _, v := range p.Restrict(Neighbors(u)).Nodes() {
		if p.HasNode(v) {
			nbs := Neighbors(v)
			p = p.RemNode(v)
			cliques = append(cliques,
				BronKerbosch(r.AddNode(v), p.Inter(nbs), x.Inter(nbs))...)
			x = x.AddNode(v)
		}
	}
	return
}

//////////// Set routines
// Sets are bitfields in chunks of uint64
// They indexes the node IDs, which are n2i indexes of the "bkedges" array
// or lists of neighbors for each node

type BKSet []uint64
var bksetmaxsize = 26*26
var bksetnchuncks = (bksetmaxsize+63)/64
var bkedges [][]int

// BKSet operations

func BKNew() BKSet {
	return make([]uint64, bksetnchuncks, bksetnchuncks)
}

func (s BKSet) IsEmpty() bool {
	for _, i := range s {
		if i != 0 { return false}
	}
	return true
}

func (s BKSet) Clone() BKSet {
	t := BKNew()
	for i := range s {
		t[i] = s[i]
	}
	return t
}

func (s BKSet) Union(t BKSet) BKSet {
	u := BKNew()
	for i := range s {
		u[i] = s[i] | t[i]
	}
	return u
}

func (s BKSet) Restrict(t BKSet) BKSet {
	u := BKNew()
	for i := range s {
		u[i] = s[i] ^ t[i]
	}
	return u
}

func (s BKSet) Inter(t BKSet) BKSet {
	u := BKNew()
	for i := range s {
		u[i] = s[i] & t[i]
	}
	return u
}

func BKMake(ns ...int) BKSet {
	return BKNew().AddNode(ns...)
}	

func (s BKSet) HasNode(n int) bool {
	i := n / 64					  // the 64-bit chunk
	mask := uint64(1) << (n % 64) // the bit in it
	return s[i] & mask != 0
}

func (s BKSet) AddNode(ns ...int) BKSet {
	t := s.Clone()
	for _, n := range ns {			  // each node to add
		i := n / 64					  // the 64-bit chunk
		mask := uint64(1) << (n % 64) // the bit in it
		t[i] |= mask
	}
	return t
}

func (s BKSet) RemNode(ns ...int) BKSet {
	t := s.Clone()
	for _, n := range ns {			  // each node to remove
		i := n / 64					  // the 64-bit chunk
		mask := uint64(1) << (n % 64) // the bit in it
		t[i] |= mask
	}
	return t
}

// specific to our comps graph

func BKSetMake(edges [][]int) BKSet {
	s := BKNew()
	nl := []int{}
	for i, e := range edges {
		if len(e) > 0 {
			nl = append(nl, int(i))
		}
	}
	return s.AddNode(nl...)
}

func (s BKSet) Nodes() (nodes []int) { // list of nodes in set
	var node int
	for _, n := range s {
		for _ = range 64 {
			if n & 1 == 1 {
				nodes = append(nodes, node)
			}
			n >>= 1
			node++
		}
	}
	return
}

func (s BKSet) Pivot() (pnode int) { // node with most edges
	var node int
	plen := 0
	for _, n := range s {
		for _ = range 64 {
			if n & 1 == 1 {
				if len(bkedges[node]) > plen {
					plen = len(bkedges[node])
					pnode = node
				}
			}
			n >>= 1
			node++
		}
	}
	return
}

func Neighbors(n int) BKSet {
	s := BKNew()
	for _, i := range bkedges[n] {
		//VPf("      adding %d to s %v =>", i, s)
		s = s.AddNode(i)
		//VPf(" %v\n", s)
	}
	return s
}

func BKSetPrint(s BKSet) (p string) {
	for i, n := range s.Nodes() {
		if i > 0 {
			p += commaSep
		}				
		p += i2n(int(n))
	}
	return
}
	
func NodesPrint(ns []int) (p string) {
	for i, n := range ns {
		if i > 0 {
			p += commaSep
		}				
		p += i2n(int(n))
	}
	return
}


//////////// Part 3

func part3(lines []string) string {
	conns := parse(lines)
	edges := make([][]int, 26*26, 26*26) // [comp]-> list of neigbours
	for _, c := range conns {
		edges[c[0]] = append(edges[c[0]], c[1])
		edges[c[1]] = append(edges[c[1]], c[0])
	}
	return BSNodesPrint(BSBronKerboschEdges(edges))
}
	
func BSNodesPrint(ns []uint) (p string) {
	for i, n := range ns {
		if i > 0 {
			p += commaSep
		}				
		p += i2n(int(n))
	}
	return
}

//////////// Bron–Kerbosch algorithm using the bitset package
// With pivot version, from:
// https://en.wikipedia.org/wiki/Bron%E2%80%93Kerbosch_algorithm#With_pivoting

func BSBronKerboschEdges(edges [][]int) []uint {
	bkedges = edges
	var r, p, x bitset.BitSet
	for i, e := range edges {
		if len(e) > 0 {
			p.Set(uint(i))
		}
	}
	cliques := BSBronKerbosch(&r, &p, &x)
	maxlen := 0
	var maxclique []uint
	for _, s := range cliques {
		nodes := BSNodes(&s)
		if len(nodes) > maxlen {
			maxlen = len(nodes)
			maxclique = nodes
		}
	}
	return maxclique
}

func BSBronKerbosch(r, p, x *bitset.BitSet) (cliques []bitset.BitSet) {
	// if P and X are both empty then report R as a maximal clique
	if p.Count() == 0 && x.Count() == 0 {
		return []bitset.BitSet{*r}
	}
	// else choose pivot vertex (with lots of edges) q from P ⋃ X:
	u := BSPivot(p.Union(x))
	// for each vertex v in P \ N(u) do
	for _, v := range BSNodes(p.SymmetricDifference(BSNeighbors(u))) {
		if p.Test(v) {
			nbs := BSNeighbors(v)
			p = p.Clear(v)
			cliques = append(cliques,
				BSBronKerbosch(r.Clone().Set(v), p.Intersection(nbs), x.Intersection(nbs))...)
			x = x.Set(v)
		}
	}
	return
}

func BSNodes(s *bitset.BitSet) (nodes []uint) { // list of nodes in set
	nodes = make([]uint, s.Count(), s.Count())
	return s.AsSlice(nodes)
}

func BSPivot(s *bitset.BitSet) (pnode uint) { // node with most edges
	plen := 0
	for _, node := range BSNodes(s) {
		if len(bkedges[node]) > plen {
			plen = len(bkedges[node])
			pnode = node
		}
	}
	return
}

func BSNeighbors(n uint) *bitset.BitSet {
	var s bitset.BitSet
	for _, i := range bkedges[n] {
		s.Set(uint(i))
	}
	return &s
}

//////////// Common Parts code

func n2i(s string) int {
	return int(s[1] - 'a') + int(s[0] - 'a') * 26
}

func i2n(i int) string {
	return string([]byte{byte(i/26)+'a', byte(i%26)+'a'})
}

func parse(lines []string) (conns [][]int) {
	recomp := regexp.MustCompile("[[:lower:]][[:lower:]]")
	for _, line := range lines {
		conn := recomp.FindAllString(line, -1)
		// auto-create comps
		c0 := n2i(conn[0])
		c1 := n2i(conn[1])
		conns = append(conns, []int{c0, c1})
	}
	return
}

//////////// PrettyPrinting & Debugging functions. See also the VPx functions.

func DEBUG() {
	if ! debug { return }
	fmt.Println("DEBUG!")
	// ad hoc debug code here
}
