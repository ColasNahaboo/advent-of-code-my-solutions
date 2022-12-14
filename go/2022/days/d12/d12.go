// Adventofcode 2022, d12, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 example 31
// TEST: -1 input 391
// TEST: example 29
// TEST: input 386
package main

import (
	"flag"
	"fmt"
	"log"
	"github.com/gammazero/deque" // fast FIFO & LIFO
	// "regexp"
)

// the grid is a contiguous slice of the rows
// we pad by a border of elevation 1000 to avoid testing for edges

type Grid struct {
	width, height, area int
	grid []int
}
// some elevations
const start = 0					// S has elevation 'a'
const goal = 'z' - 'a'			// E has elevation 'z'
const border = 1000
var startp, goalp int			// their positions on the grid
var dirs [4]int					// deltas to the adjacent positions

var verbose bool

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
	grid := parseGrid(lines)
	VP(grid)

	var result int
	if *partOne {
		VP("Running Part1")
		result = part1(grid)
	} else {
		VP("Running Part2")
		result = part2(grid)
	}
	fmt.Println(result)
}

//////////// Part 1

func part1(grid Grid) int {
	VPf("Starting at %d, looking for %d\n", startp, goalp)
	path := bfsSearch(grid, startp, goalp)
	VPf("Shortest Path: %v\n", path)
	return len(path) - 1		// dont count starting point in # of steps
}

//////////// Part 2
func part2(g Grid) (best int) {
	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			p := ((y+1)*(g.width+2)) + x + 1
			if g.grid[p] == 0 { // 'a' place
				pathlen := len(bfsSearch(g, p, goalp)) - 1
				if pathlen >0 && pathlen < best || best == 0 {
					best = pathlen
				}
			}
		}
	}
	return
}

//////////// Common Parts code

func parseGrid(lines []string) (g Grid) {
	g.height = len(lines)
	if g.height < 1 {log.Fatalln("No grid!");}
	g.width = len(lines[0])
	g.area = (g.width + 2) * (g.height + 2) // room for the edges
	g.grid = make([]int, g.area, g.area)
	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			g.grid[((y+1)*(g.width+2)) + x + 1] = int(lines[y][x]) - 'a'
			if lines[y][x] == 'S' {
				startp = ((y+1)*(g.width+2)) + x + 1
				g.grid[startp] = start
			} else if lines[y][x] == 'E' {
				goalp = ((y+1)*(g.width+2)) + x + 1
				g.grid[goalp] = goal
			}
		}
		//  left & right borders
		g.grid[((y+1)*(g.width+2))] = border
		g.grid[((y+1)*(g.width+2)) + g.width + 1] = border
	}
	//  top and bottom borders
	for x := 0; x < (g.width + 2); x++ {
		g.grid[x] = border
		g.grid[((g.height+1)*(g.width+2)) + x] = border
	}
	// delta of position for one step in R D U L dirs
	dirs = [4]int{1, g.width + 2, - g.width - 2, -1}
	return
}

// from position p, applies a BFS search to the grid
// prevpath is the sequence of pos we went through before p
func bfsSearch(g Grid, startp, goalp int) (path []int) {
	frontier := deque.New[int](1024,1024)
	frontier.PushBack(startp)
	previous := make([]int, g.area, g.area)
	visited := make([]bool, g.area, g.area)
	visited[startp] = true

	// look at all nodes in the FIFO frontier
	for frontier.Len() != 0 {
		p := frontier.PopFront()
		if p == goalp {
			VPf("Goal %d reached from %d\n", p, previous[p])
			goto GOAL_REACHED
		}
		for _, dir := range dirs {
			new := p + dir
			if visited[new] {	// already processed
				continue
			}
			if g.grid[new] > g.grid[p] + 1 { // check new is a valid step
				continue
			}
			// ok, push new onto the frontier FIFO
			previous[new] = p
			frontier.PushBack(new)
			visited[new] = true
			VPf("Considering %d from %d\n", new, p)
		}
	}
	return
GOAL_REACHED:
	// retrace our steps to get how we got to the goal
	path = []int{goalp}
	for p := goalp; p != startp; p = previous[p] {
		path = prependInt(path, previous[p])
	}
	return
}

func prependInt(x []int, y int) []int {
    x = append(x, 0)
    copy(x[1:], x)
    x[0] = y
    return x
}

	
//////////// Part1 functions

//////////// Part2 functions
