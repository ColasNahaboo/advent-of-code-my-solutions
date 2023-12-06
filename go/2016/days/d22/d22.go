// Adventofcode 2016, d22, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,,.test
// TEST: -1 example 5
// TEST: example
package main

import (
	"flag"
	"fmt"
	"regexp"
	//"time"
	//"github.com/fzipp/astar"
)

type Node struct {				// the initial state
	x, y, size, used, avail int
}
var grid []Node					// linear of the 2d grid: pos = x + y*gw
var gw, gh, garea int			// its dims

// data of nodes for part2 A* graph search of shortest path
// for each pos in grid, delta of avail to it
// e.g: avail at p is grid[p].avail + state[p]
// an extra int is the position of data G, thus at state[garea]
// BUT we use 2-byte numbers instead of ints, to be able to use strings
// as indexes in a map of states, using the optimisation in
// https://pthevenet.com/posts/programming/go/bytesliceindexedmaps/
// v := m[string(byteSlice)]
// Thus avail at p is grid[p].avail + state[p*2]*256 + state[p*2+1]
type State []byte				// of length (garea + 1) * 2
var states map[string]State

var verbose bool

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
	parse(lines[2:])
	
	var result int
	if *partOne {
		VP("Running Part1")
		result = part1()
	} else {
		VP("Running Part2")
		result = part2()
	}
	fmt.Println(result)
}

//////////// Part 1

func part1() (viable int) {
	for a := 0; a < garea; a++ {
		for b := 0; b < garea; b++ {
			if grid[a].used != 0 &&
				a != b &&
				grid[a].used < grid[b].avail {
				viable++
			}
		}
	}
	return			
}

//////////// Part 2

func part2() int {
	return 0
}

//////////// Common Parts code

func parse(lines []string) {
	re := regexp.MustCompile("/dev/grid/node-x([[:digit:]]+)-y([[:digit:]]+)[[:space:]]+([[:digit:]]+)T[[:space:]]+([[:digit:]]+)T[[:space:]]+([[:digit:]]+)T[[:space:]]+")
	nodes := []Node{}				// list, not ordered in a 2D grid
	for lineno := 0; lineno < len(lines); lineno++ {
		line := lines[lineno]
		if m := re.FindStringSubmatch(line); m != nil {
			node := Node{
				x: atoi(m[1]),
				y: atoi(m[2]),
				size: atoi(m[3]),
				used: atoi(m[4]),
				avail: atoi(m[5]),
			}
			if node.x > gw -1 { gw = node.x +1; }
			if node.y > gh -1 { gh = node.y +1; }
			nodes = append(nodes, node)
		} else {
			panic(fmt.Sprintf("Syntax error line %d: %s\n", lineno, line))
		}
	}
	// build the 2D grid
	garea = gw * gh
	grid = make([]Node, garea, garea)
	for _, node := range nodes {
		grid[node.x + node.y * gw] = node
	}
	// check we do not have holes
	for i := 0; i < garea; i++ {
		if grid[i].size == 0 {
			panic(fmt.Sprintf("grid missing a node at (%d, %d)\n", i%gw, i/gw))
		}
	}
	VPf("Grid %d x %d\n", gw, gh)
}


//////////// Part1 functions

//////////// Part2 functions
