// Adventofcode 2015, day18, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 input 814
// TEST: input 924
package main

import (
	"flag"
	"fmt"
	// "regexp"
)

// Grid is 2-dim array of lights, 1=on, 0=off, with a border of 0
type Grid = [][]int

var verbose bool

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	stepsFlag := flag.Int("s", 100, "number of steps to run")
	flag.Parse()
	verbose = *verboseFlag
	steps := *stepsFlag
	infile := "input.txt"
	if flag.NArg() > 0 {
		infile = flag.Arg(0)
	}
	lines := fileToLines(infile)
	grid := readGrid(lines)

	var result int
	if *partOne {
		VP("Running Part1")
		result = part1(grid, steps)
	} else {
		VP("Running Part2")
		result = part2(grid, steps)
	}
	fmt.Println(result)
}

//////////// Part 1

func part1(grid Grid, steps int) int {
	for i := 0; i < steps; i++ {
		grid = gridStep(grid)
	}
	return gridOns(grid)
}

//////////// Part 2
func part2(grid Grid, steps int) int {
	turnOnCorners(grid)
	for i := 0; i < steps; i++ {
		grid = gridStep(grid)
		turnOnCorners(grid)
	}
	return gridOns(grid)
}

//////////// Common Parts code

// we read a grid and add a "dead" border of width 1 around it
func readGrid(lines []string) Grid {
	size := len(lines)
	VP("Greating grid of size", size)
	grid := make(Grid, size+2)
	grid[0] = make([]int, size+2)
	grid[size+1] = make([]int, size+2)
	for y := 0; y < size; y++ {
		row := make([]int, size+2)
		for x := 0; x < size; x++ {
			if lines[y][x] == '#' {
				row[x+1] = 1
			}
		}
		grid[y+1] = row
	}
	return grid
}

func gridStep(grid Grid) Grid {
	size := len(grid)
	// creta a new grid
	new := make(Grid, size)
	for y := 0; y < size; y++ {
		new[y] = make([]int, size)
	}
	// populate it
	for x := 1; x < size-1; x++ {
		for y := 1; y < size-1; y++ {
			// adj = number of neigbours (adjacent lights)
			adj := grid[x-1][y-1] + grid[x][y-1] + grid[x+1][y-1] + grid[x+1][y] + grid[x+1][y+1] + grid[x][y+1] + grid[x-1][y+1] + grid[x-1][y]
			if grid[x][y] == 1 { // light was on
				if adj == 2 || adj == 3 {
					new[x][y] = 1
				}
			} else { // light was off
				if adj == 3 {
					new[x][y] = 1
				}
			}
		}
	}
	return new
}

// we count the "on" lights in grid
func gridOns(grid Grid) int {
	count := 0
	size := len(grid)
	for x := 1; x < size-1; x++ {
		for y := 1; y < size-1; y++ {
			if grid[x][y] == 1 {
				count++
			}
		}
	}
	return count
}

//////////// Part1 functions

//////////// Part2 functions

// force 4 corners to be on
func turnOnCorners(grid Grid) {
	size := len(grid)
	grid[1][1] = 1
	grid[1][size-2] = 1
	grid[size-2][1] = 1
	grid[size-2][size-2] = 1
}
