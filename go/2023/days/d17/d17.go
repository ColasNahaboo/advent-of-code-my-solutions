// Adventofcode 2023, d17, in go. https://adventofcode.com/2023/day/17
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 102
// TEST: example 94
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

// This is typical of the adventcode solutions, where we find the shortest
// path in states of "things" on a 2D board
// For this, we combine scalarray.go (2D board) and astar.go (shortest path)

package main

import (
	"flag"
	"fmt"
)

var verbose bool

type State struct {			   // an (unique) state of the crucible
	pos int					   // position of the crucible in the city Scalarray
	dir int					   // its direction: NESW = 0 1 2 3
	steps int				   // how many steps have we just done in dir?
}

// context global vars: these could also be put in a struct whose pointer to
// would have been passed as first parameter (the "graph") to AStarFindPath
var states = []State{}				// the dynamically-maintained set of states
var stateIDs = make(map[State]int, 0) // map States to their index in states: ID
var city Scalarray[int]				  // the map of the heatlosses
var dirs [4]int					// position offsets NESW

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	partThree := flag.Bool("3", false, "run exercise part3, (default: part2)")
	partFour := flag.Bool("4", false, "run exercise part4, (default: part2)")
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
	} else if *partThree {
		VP("Running Part3")
		result = part3(lines)
	} else if *partFour {
		VP("Running Part4")
		result = part4(lines)
	} else {
		VP("Running Part2")
		result = part2(lines)
	}
	fmt.Println(result)
}

//////////// Part 1

func part1(lines []string) (sum int) {
	city = parse(lines)
	dirs = city.Dirs()
	// dir field is not important for start and end
	start := State{0, 0, 0}
	startnode := stateNode(start)
	end := State{city.w*city.h - 1, 0, 0}
	endnode := stateNode(end)
	path := AStarFindPath[any, int](nil, startnode, endnode, connectedNodes, distNodes, distEnd, samePosNodes)
	for i := 1; i < len(path); i++ { //  dont count start heat loss
		sum += city.a[states[path[i]].pos]
	}
	VPmapPath(path[1:])
	return
}

// return Node (ID) of the state, auto-allocating it if needed
func stateNode(s State) (id int) {
	var ok bool
	if id, ok = stateIDs[s]; !ok {
		id = len(states)
		states = append(states, s)
		stateIDs[s] = id
	}
	return
}

// neighbour states. Where can the crucible go?
func connectedNodes(g any, n int) (cns []int) {
	s := states[n]
	// go straight?
	if s.steps < 3 && city.stepDirInside(s.pos, s.dir) {
		cns = append(cns, stateNode(State{s.pos+dirs[s.dir], s.dir, s.steps+1}))
	}
	// turn right?
	newdir := (s.dir + 1) % 4
	if city.stepDirInside(s.pos, newdir) {
		cns = append(cns, stateNode(State{s.pos+dirs[newdir], newdir, 1}))
	}
	// turn left?
	newdir = (s.dir + 3) % 4
	if city.stepDirInside(s.pos, newdir) {
		cns = append(cns, stateNode(State{s.pos+dirs[newdir], newdir, 1}))
	}
	return
}

// reached end, any direction is OK
func samePosNodes(g any, n1, n2 int) bool {
	return states[n1].pos == states[n2].pos
}

// "cost", between 2 points: Mahattan distance + destination heatloss
// but here case all neighbours are adjacent so we only count the heatloss
func distNodes(g any, n1, n2 int) float64 {
	return float64(city.a[states[n2].pos])
}

// distance to end: Mahattan distance
func distEnd(g any, n1, n2 int) float64 {
	x1, y1 := city.Coords(states[n1].pos)
	x2, y2 := city.Coords(states[n2].pos)
	return float64(intAbs(x2-x1) + intAbs(y2-y1))
}

func parse(lines []string) Scalarray[int] {
	sa := makeScalarray[int](len(lines[0]), len(lines))
	for y, line := range lines {
		for x := range line {
			sa.Set(x, y, atoi(line[x:x+1]))
		}
	}
	return sa
}

//////////// Part 2

// Part2 is exactly like Part1, but with different neighbours for each state
// and a different end condition

func part2(lines []string) (sum int) {
	city = parse(lines)
	dirs = city.Dirs()
	// dir field is not important for start and end
	start := State{0, 0, 0}
	startnode := stateNode(start)
	end := State{city.w*city.h - 1, 0, 0}
	endnode := stateNode(end)
	path := AStarFindPath[any, int](nil, startnode, endnode, connectedNodes2, distNodes, distEnd, samePosNodes2)
	for i := 1; i < len(path); i++ { //  dont count start heat loss
		sum += city.a[states[path[i]].pos]
	}
	VPmapPath(path[1:])
	return
}


// neighbour states. Where can the crucible go?
func connectedNodes2(g any, n int) (cns []int) {
	s := states[n]
	// go straight?
	if s.steps < 10 && city.stepDirInside(s.pos, s.dir) {
		cns = append(cns, stateNode(State{s.pos+dirs[s.dir], s.dir, s.steps+1}))
	}
	// cannot turn if we do not have performed 4 steps
	if s.steps >= 1 && s.steps < 4 {
		return
	}
	// turn right?
	newdir := (s.dir + 1) % 4
	if city.stepDirInside(s.pos, newdir) {
		cns = append(cns, stateNode(State{s.pos+dirs[newdir], newdir, 1}))
	}
	// turn left?
	newdir = (s.dir + 3) % 4
	if city.stepDirInside(s.pos, newdir) {
		cns = append(cns, stateNode(State{s.pos+dirs[newdir], newdir, 1}))
	}
	return
}

// reached end, any direction is OK, but we must have done at least 4 steps
func samePosNodes2(g any, n1, n2 int) bool {
	return states[n1].pos == states[n2].pos && states[n1].steps >= 4
}

//////////// Part 3

// Part3 is a Part2 implementation variant, implemented with a 3D array id3d to
// store states directly, without the indirection of a mapTable stateIDs.
// A state ID is thus its position (number) in the array of Scalarray3D
// We do not have an actual State type, we map the fields virtually on the 3D
// coords: pos is x, dir is y, steps is z
// Note that we do not need to actually instanciate the scalar array field of
// id3d, we only use its coordinates/position conversion methods, as we do not
// have to store additional data to Nodes.
// Name of functions rewritten to use id3d insteads of stateIDs end with "3"

var id3d Scalarray3D[int]

func part3(lines []string) (sum int) {
	city = parse(lines)
	dirs = city.Dirs()
	maxsteps := city.w
	if city.h > city.w {
		maxsteps = city.h
	}
	id3d = Scalarray3D[int]{w:len(city.a), h:4, d:maxsteps} // no actuall array
	// dir field is not important for start and end
	startnode := id3d.Pos(0, 0, 0)
	endnode := id3d.Pos(city.w*city.h - 1, 0, 0)
	path := AStarFindPath[any, int](nil, startnode, endnode, connectedNodes3, distNodes3, distEnd3, samePosNodes3)
	for i := 1; i < len(path); i++ { //  dont count start heat loss
		sum += city.a[id3d.X(path[i])]
	}
	return
}

// neighbour states. Where can the crucible go?
func connectedNodes3(g any, n int) (cns []int) {
	pos, dir, steps := id3d.Coords(n) // decode fields of virtual state
	// go straight?
	if steps < 10 && city.stepDirInside(pos, dir) {
		cns = append(cns, id3d.Pos(pos+dirs[dir], dir, steps+1))
	}
	// cannot turn if we do not have performed 4 steps
	if steps >= 1 && steps < 4 {
		return
	}
	// turn right?
	newdir := (dir + 1) % 4
	if city.stepDirInside(pos, newdir) {
		cns = append(cns, id3d.Pos(pos+dirs[newdir], newdir, 1))
	}
	// turn left?
	newdir = (dir + 3) % 4
	if city.stepDirInside(pos, newdir) {
		cns = append(cns, id3d.Pos(pos+dirs[newdir], newdir, 1))
	}
	return
}

// reached end, any direction is OK, but we must have done at least 4 steps
func samePosNodes3(g any, n1, n2 int) bool {
	return id3d.X(n1) == id3d.X(n2) && id3d.Z(n1) >= 4
}

// "cost", between 2 points: Mahattan distance + destination heatloss
// but here case all neighbours are adjacent so we only count the heatloss
func distNodes3(g any, n1, n2 int) float64 {
	return float64(city.a[id3d.X(n2)])
}

// distance to end: Mahattan distance
func distEnd3(g any, n1, n2 int) float64 {
	x1, y1 := city.Coords(id3d.X(n1))
	x2, y2 := city.Coords(id3d.X(n2))
	return float64(intAbs(x2-x1) + intAbs(y2-y1))
}

//////////// Part 4

// Part4 is a Part2 implementation variant, also without the indirection of a
// mapTable stateIDs, by using the states themselves directly as Nodes.

func part4(lines []string) (sum int) {
	city = parse(lines)
	dirs = city.Dirs()
	// dir field is not important for start and end
	start := State{0, 0, 0}
	end := State{city.w*city.h - 1, 0, 0}
	path := AStarFindPath[any, State](nil, start, end, connectedNodes4, distNodes4, distEnd4, samePosNodes4)
	for i := 1; i < len(path); i++ { //  dont count start heat loss
		sum += city.a[path[i].pos]
	}
	return
}

// neighbour states. Where can the crucible go?
func connectedNodes4(g any, s State) (cns []State) {
	// go straight?
	if s.steps < 10 && city.stepDirInside(s.pos, s.dir) {
		cns = append(cns, State{s.pos+dirs[s.dir], s.dir, s.steps+1})
	}
	// cannot turn if we do not have performed 4 steps
	if s.steps >= 1 && s.steps < 4 {
		return
	}
	// turn right?
	newdir := (s.dir + 1) % 4
	if city.stepDirInside(s.pos, newdir) {
		cns = append(cns, State{s.pos+dirs[newdir], newdir, 1})
	}
	// turn left?
	newdir = (s.dir + 3) % 4
	if city.stepDirInside(s.pos, newdir) {
		cns = append(cns, State{s.pos+dirs[newdir], newdir, 1})
	}
	return
}

// reached end, any direction is OK, but we must have done at least 4 steps
func samePosNodes4(g any, s1, s2 State) bool {
	return s1.pos == s2.pos && s1.steps >= 4
}

// "cost", between 2 points: Mahattan distance + destination heatloss
// but here case all neighbours are adjacent so we only count the heatloss
func distNodes4(g any, s1, s2 State) float64 {
	return float64(city.a[s2.pos])
}

// distance to end: Mahattan distance
func distEnd4(g any, s1, s2 State) float64 {
	x1, y1 := city.Coords(s1.pos)
	x2, y2 := city.Coords(s2.pos)
	return float64(intAbs(x2-x1) + intAbs(y2-y1))
}

//////////// Common Parts code

//////////// PrettyPrinting & Debugging functions

var dirSymbols = [4]string{"^", ">", "v", "<"}

func VPmapPath(path []int) {
	if ! verbose {
		return
	}
	sa := makeScalarray[string](city.w, city.h)
	for p, i := range city.a {
		sa.a[p] = itoa(i)
	}
	for _, n := range path {
		sa.a[states[n].pos] = dirSymbols[states[n].dir]
	}
	for p, s := range sa.a {
		if p % sa.w == 0 {
			fmt.Println()
		}
		fmt.Print(s)
	}
	fmt.Println()
}
