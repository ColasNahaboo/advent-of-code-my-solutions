//go:build exclude
// This is an old implementation, with states being ints used as bitfields
// each object (G, then M) floor number being encoded in 2 bits,
// the floor number of E being the first two
// time to run was 1mn without the metric optim (calls to metric(ss))
// 31s with it

// Adventofcode 2016, d11, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 example 11
// TEST: -1 input 33
// #TEST: example N/A: example fails if we add the 4 new objects of part2
// TEST: input 57

// TODO: implement A* ourselves, so we can avoid re-visiting equivalent states
// (states with same metric)

package main

import (
	"flag"
	"fmt"
	"regexp"
	"strconv"	
	// we could use gonum A*, but it requires creating all the states beforehand
    //   https://pkg.go.dev/gonum.org/v1/gonum/graph/path#AStar
	//   "gonum.org/v1/gonum/graph"
	//   "gonum.org/v1/gonum/graph/path"
	//   "gonum.org/v1/gonum/graph/simple"
	// So we use astar, that just asks for neighbours on demand
	"github.com/fzipp/astar"
)

var verbose bool

// a state is a series of 2-bit numbers that are the floor the object is on
// object ids are 0 for E, i*2+1, i*2+2 for metal #i G and M
// this is because there are only 4 floors and 5 metals max = 22 bits used.
// reminder conversions of N
//   decimal -> binary: strconv.FormatInt(N, 2)
//   binary -> decimal: strconv.ParseInt(N, 2, 64)

var names []string				// debug:  names of the objects by their IDs
var metalids map[string]int		// IDs of the G for a name
var nmetals int					// its length
var nobjs int					// number of objects: E, MG, MM, XG, XM...
type graph []int
var nodes int					// debug: number of nodes examined for neighbours

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	flag.Parse()
	verbose = *verboseFlag
	infile := "input.txt"
	if flag.NArg() > 0 {
		infile = flag.Arg(0)
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
	start := parse(lines)
	dest := 0
	// goal: everything on 4th floor
	for i := 0; i < nobjs; i++ {
		dest = floorObjSet(dest, i, 3)
	}

	VPf("Looking for path to %d %v\n", dest, state2slice(dest))
	states := graph{start, dest}
	path := astar.FindPath[int](states, start, dest, stateDist, stateDist)
	if verbose {
		for i, s := range path {
			VPf("  [%d] %v\n", i, state2slice(s))
		}
	}
	fmt.Printf("Number of explored nodes: %d\n", nodes)
	return len(path) - 1		// path includes both ends
}

//////////// Part 2
func part2(lines []string) int {
	// we just insert the elerium and dilithium Gs & Ms in floor 0
	lines[0] = lines[0] + " And also a elerium generator, a elerium-compatible microchip, a dilithium generator, a dilithium-compatible microchip."
	return part1(lines)
}

//////////// Common Parts code

func parse(lines []string) (state int) {
	metalids = make(map[string]int, 0)
	re := regexp.MustCompile(" [Aa]n? (([[:alpha:]]+) generator|([[:alpha:]]+)-compatible microchip)")
	// Elevator starts on the first floor (#0)
	names = append(names, "elevator")
	for floor := 0; floor < 3; floor++ { // inputs describe only floors #0 to #2
		for _, m := range re.FindAllStringSubmatch(lines[floor], -1) { // examine all mentions of objects
			VPf("Parsing G = \"%s\", M = \"%s\"\n", m[2], m[3])
			var isMC, id int
			var ok bool
			metal := m[2]
			if metal == "" {
				isMC = 1
				metal = m[3]
			}
			if id, ok = metalids[metal]; !ok { // metal not yet seen, allocate room for it
				id = len(metalids)
				metalids[metal] = id
				names = append(names, metal + " generator")
				names = append(names, metal + " microchip")
			}
			state = floorSet(state, id, isMC, floor)
			VPf("    ==> state = %d (id=%d, isMC=%d, floor=%d\n", state, id, isMC, floor)
		}
	}
	nmetals = len(metalids)
	nobjs = 1 + nmetals * 2		// elevator + one G and one M per metal
	VPf("%d metals, Initial state: %d %v\n", nmetals, state, state2slice(state))
	return
}

func floorGet(state, metal, isMc int) int {
	return (state >> ((metal*2+isMc+1)*2)) & 3
}
func floorSet(state, metal, isMc, floor int) int {
	return (state &^ (3 << ((metal*2+isMc+1)*2))) | (floor << ((metal*2+isMc+1)*2))
}
func floorObjGet(state, i int) int {
	return (state >> (i*2)) & 3
}
func floorObjSet(state, i, floor int) int {
	return (state &^ (3 << (i*2))) | (floor << (i*2))
}

//// Metric functions on States for use by astar

// we use the sum of floor distances per object, without counting E
func stateDist(s1, s2 int) (float64) {
	var od, d int
	for i := 0; i < nmetals; i++ {
		for isMc := 0; isMc < 2; isMc++ {
			od = floorGet(s1, i, isMc) - floorGet(s2, i, isMc)
			if od < 0 { od = -od;}
			d += od
		}
	}
	return float64(d)
}

// here we implement all the constraints on object moves, to list only possible neighbour states
// - look for moves to next or previous floor of E only
// - E goes to tofloor
// - of 1 or 2 objects (not counting E)
// - on fromfloor and tofloor, there must be no G or all M are connected to their G

func (g graph) Neighbours(s int) (nexts []int) {
	nodes++
	f := floorObjGet(s, 0)
	if f > 0 {
		nexts = append(nexts, floorNeighbours(floorObjSet(s, 0, f - 1), f, f-1)...)
	}
	if f < 3 {
		nexts = append(nexts, floorNeighbours(floorObjSet(s, 0, f + 1), f, f+1)...)
	}
	return
}

// [high optim] ALL PAIRS ARE INTERCHANGEABLE - The following two states are
// EQUIVALENT: x 0 1 2 2 and x 2 2 0 1

// returns list of all neighbours on to floor (tf) of state s, from floor ff
func floorNeighbours(s, ff, tf int) (nf []int) {
	sso := state2slice(s) // create slice-version of tofloor
	var nfmetrics []int
	for _, pair := range objPairsOnFloor(s, ff) {
		ss := make([]int, nobjs, nobjs)
		copy(ss, sso)
		for _, o := range pair { // move the pair to tf in ss
			ss[o] = tf
		}
		if isStateSliceValid(ss, ff) && isStateSliceValid(ss, tf) {
			m := metric(ss)
			if indexOfInt(nfmetrics, m) != -1 { continue;} // skip equivalents
			nfmetrics = append(nfmetrics, m)
			nf = append(nf, slice2state(ss))
		}
	}
	return
}

// same metric number for equivalent states
func metric(ss []int) (m int) {
	pairs := make([]int, 16, 16)
	for i := 1; i < nobjs; i += 2 {
		pair := ss[i] + ss[i+1] * 4
		pairs[pair]++
	}
	for i, n := range pairs {
		m |= n << (i*3)
	}
	m |= ss[0] << (16*3) 		// E
	return
}

func isStateSliceValid(ss []int, floor int) bool {
	for i := 1; i < nobjs; i += 2 { // No Gs on floor? ==> valid
		if ss[i] == floor {
			goto HAS_G
		}
	}
	return true
HAS_G:
	for i := 2; i < nobjs; i += 2 { // look at all Ms on floor
		if ss[i] == floor && ss[i-1] != floor {
			return false // a M without its G? invalid!
		}
	}
	return true
}

func objsOnFloor(s, f int) (objs []int) {
	for i := 1; i < nobjs; i++ {
		if floorObjGet(s, i) == f {
			objs = append(objs, i)
		}
	}
	return
}

// return all the sets of 1 or 2 non-E objects on floor f
func objPairsOnFloor(s, f int) (sets [][]int) {
	for i := 1; i < nobjs; i++ {
		if floorObjGet(s, i) == f {
			sets = append(sets, []int{i})
			for j := i + 1; j < nobjs; j++ {
				if floorObjGet(s, j) == f {
					sets = append(sets, []int{i, j})
				}
			}
		}
	}
	return
}

//// slices (ss) are unpacked states (s):
// slices of floors as integers indexed by objects
// they are easier to handle for some operations

// state int ==> slice of floor numbers
func state2slice(s int) []int {
	ss := make([]int, nobjs, nobjs)
	for i := 0; i < nobjs; i++ {
		ss[i] = (s >> (i * 2)) & 3
	}
	return ss
}

// slice of floor numbers ==> state int
func slice2state(ss []int) (s int) {
	for i := 0; i < nobjs; i++ {
		s |= ss[i] << (i * 2)
	}
	return
}

//////////// Part1 functions

//////////// Part2 functions

//////////// Debugging tools

// decimal <-> string of bits
func b2d(s string) int {
	i, err := strconv.ParseInt(s,2,64)
	if err != nil { panic(fmt.Sprintf("Not a string of bits: \"%s\"\n", s));}
	return int(i)
}
func d2b(i int) string {
	return strconv.FormatInt(int64(i),2)
}
