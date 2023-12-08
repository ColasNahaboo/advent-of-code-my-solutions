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
	"github.com/fzipp/astar"
)

type Node struct {				// the initial state
	x, y, size, used, avail int
	t int						// taquin type of the node
}
var grid []Node					// linear of the 2d grid: pos = x + y*gw
var gw, gh, garea int			// its dims

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

// Part1 is totally simplistic
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

//////////// Common Parts code

func parse(lines []string) {
	re := regexp.MustCompile("/dev/grid/node-x([[:digit:]]+)-y([[:digit:]]+)[[:space:]]+([[:digit:]]+)T[[:space:]]+([[:digit:]]+)T[[:space:]]+([[:digit:]]+)T[[:space:]]+")
	nodes := []Node{}				// list, not ordered in a 2D grid
	hole_p := -1
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
			if node.used == 0 {
				if hole_p != -1 {
					panic("grid has 2 holes")
				}
				hole_p = lineno
				node.t = HOLE
			} else if node.used > 100 {
				node.t = WALL
			} else {
				node.t = TILE
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
	// our goal is at x=gw-1, y=0
	grid[gw-1].t = GOAL
	// checks: missing nodes
	for i := 0; i < garea; i++ {
		if grid[i].size == 0 {
			panic(fmt.Sprintf("grid missing a node at (%d, %d)\n", i%gw, i/gw))
		}
	}
	if hole_p == -1 {
		panic("grid has no hole!")
	}
	VPf("Grid %d x %d\n", gw, gh)
}

//////////// Part 2

// For part 2, we solve this as a taquin, with A* of github.com/fzipp/astar
// There is one hole, and movable tiles or immovable ealls
// We must bring the goal tile to (0, 0) from its (gw-1, 0) position

type State struct {
	hole int					// position of the hole
	data int					// position of the data
	grid []int8					// the grid, of size gw*gh
}
var gw, gh int					// grid size with a border of walls added

const (							// taquin types. movable tile: >0
	HOLE = 0
	TILE = 1
	GOAL = 2
	WALL = -1
)



func part2() int {
	return 0
}


