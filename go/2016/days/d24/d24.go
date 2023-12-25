// Adventofcode 2016, d24, in go. https://adventofcode.com/2016/day/24
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 14
// TEST: example 20
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"flag"
	"fmt"
	//"regexp"
)

var verbose bool
const MaxUint = ^uint(0)
const MaxInt = int(MaxUint >> 1)

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
	grid, pois := parse(lines)
	VPf("Parsed %d x %d Grid, with %d POIs at pos %v\n", grid.w, grid.h, len(pois), pois)
	// the matrix of the shortest distances between POIs
	dists := POIdists(grid, pois)
	VPgridInt("POIs distances matrix", dists)
	tour, tourlen := shortestOpenTour(dists)
	VPf("Shortest tour %d: %v\n", tourlen, tour)
	return tourlen
}

// Possible open paths going through all POIs, but starting at 0
func shortestOpenTour(grid Grid[int]) (tour []int, tourlen int) {
	l := grid.w - 1				// iterate on all the numbers except 0
	nums := make([]int, l, l)	// 1,2,3,...
	for i := 0; i < l; i++ {
		nums[i] = i+1
	}
	tourlen = MaxInt
	for p := make([]int, l); p[0] < len(p); nextPerm(p) {
		t := getPerm(nums, p)
		tlen := grid.Get(0, t[0]) // first step 0->
		for i := 0; i < l-1; i++ {
			tlen += grid.Get(t[i], t[i+1])
			if tlen >= tourlen {
				continue
			}
		}
		if tlen < tourlen {
			tour = t
			tourlen = tlen
		}
	}
	tour = append([]int{0}, tour...)
	return
}

//////////// Part 2

func part2(lines []string) int {
	grid, pois := parse(lines)
	VPf("Parsed %d x %d Grid, with %d POIs at pos %v\n", grid.w, grid.h, len(pois), pois)
	// the matrix of the shortest distances between POIs
	dists := POIdists(grid, pois)
	VPgridInt("POIs distances matrix", dists)
	tour, tourlen := shortestLoopTour(dists)
	VPf("Shortest loop %d: %v\n", tourlen, tour)
	return tourlen
}

// Possible closed paths going through all POIs, but starting and ending at 0 
func shortestLoopTour(grid Grid[int]) (tour []int, tourlen int) {
	l := grid.w - 1				// iterate on all the numbers except 0
	nums := make([]int, l, l)	// 1,2,3,...
	for i := 0; i < l; i++ {
		nums[i] = i+1
	}
	tourlen = MaxInt
	for p := make([]int, l); p[0] < len(p); nextPerm(p) {
		t := getPerm(nums, p)
		tlen := grid.Get(0, t[0]) // first step 0->
		for i := 0; i < l-1; i++ {
			tlen += grid.Get(t[i], t[i+1])
			if tlen >= tourlen {
				continue
			}
		}
		tlen += grid.Get(t[l-1], 0)
		if tlen < tourlen {
			tour = t
			tourlen = tlen
		}
	}
	tour = append([]int{0}, tour...)
	return
}

//////////// Common Parts code
// in the grid, a wall is true, a free space is false
// Points Of Interest (POI) positions are gathered into a list of their positions

func parse(lines []string) (g Grid[bool], pois []int) {
	g = makeGrid[bool](len(lines[0]), len(lines))
	pois = make([]int, 10, 10)	// pios are single digits, thus inside [0-9]
	maxpoi := 0
	for y, line := range lines {
		for x, c := range line {
			if c == '#' {		 // Wall
				g.Set(x, y, true)
			} else if c != '.' { // Points Of Interest
				poi := int(c - '0')
				pois[poi] = g.Pos(x, y)
				if poi > maxpoi {
					maxpoi = poi
				}
			}
		}
	}
	pois = pois[0:maxpoi+1]		// trim unused positions
	return
}

//////////// Paths

// computes the matrix of the shortest distances between POIs

func POIdists(grid Grid[bool], pois []int) (dists Grid[int]) {
	dists = makeGrid[int](len(pois), len(pois))
	for i1, p1 := range pois {
		for i2, p2 := range pois[0:i1] {
			d := len(shortestPath(grid, p1, p2)) - 1 // path includes start & end
			dists.g[dists.Pos(i1, i2)] = d
			dists.g[dists.Pos(i2, i1)] = d
		}
	}
	return
}

func shortestPath(g Grid[bool], p1, p2 int) (path []int) {
	path = AStarFindPath[Grid[bool], int](g, p1, p2, connectedPOIs[int], manhattanPOIs, manhattanPOIs, reachedPOI)
	return
}

// Callbacks for AStarFindPath for paths between POIs

func connectedPOIs[Node int](g Grid[bool], n Node) (conns []Node) {
	p := int(n)
	for _, d := range g.Dirs(p) {
		if ! g.g[p + d] {
			conns = append(conns, Node(p + d))
		}
	}
	return
}
func manhattanPOIs(g Grid[bool], p1, p2 int) float64 {
	x1, y1 := g.Coords(p1)
	x2, y2 := g.Coords(p2)
	return float64(intAbs(x2-x1) + intAbs(y2-y1))
}
func reachedPOI(g Grid[bool], p, end int) bool {
	return p == end
}

	

//////////// PrettyPrinting & Debugging functions

