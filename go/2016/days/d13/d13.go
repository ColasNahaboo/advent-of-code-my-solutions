// Adventofcode 2016, d13, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 example 11
// TEST: -1 input 82
// #TEST: example 82
// TEST: input 138
package main

import (
	"flag"
	"fmt"
	"time"
	"image"
	//"math"
	"strconv"
	"github.com/fzipp/astar" // easier than gonum/path to build the nodes on demand
)

// We solve the shortest path problem by the A* (A-star) algorithm.
// A node is defined by its coordinates [x, y]. For simplicity we use the
// Points of the image package just like in the examples in fzipp/astar doc:
// https://pkg.go.dev/github.com/fzipp/astar#example-FindPath-Maze

// Part2: This is relatively simple, except that the part2 problem text is
// ambiguous. I interpreted it as: try all the locations on the floor that have
// a shortes path to them of less than 50 steps. For my input, the shortest
// path to `(31,39)` was `82`, and by tring all the 1352 locations within 50
// steps of `(1,1)` for the Manhattan distance (in a triangle) that had an
// actual shorter path of 50 or less from `(1,1)` I found 615 positions. But
// this was not the expected answer! I thus kept the code for reference and
// made it callable via a new `-3` command line option, and implemented what
// seemed the interpretation of the puzzle author by looking at the solutions:
// it was, when looking for the shortest path to `(31,39)`, all the locations
// examined during this specific search of the shorted path to this single
// location... and none other! This yield the expected result of `138` in my
// case. But I still think that the problem was wrongly formulated!


var verbose bool
type FloorPlan int				// just used for typing, int is the DFN
const maxsteps = 50				// for part2

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	partThree := flag.Bool("3", false, "part3: a different interpretation of the part2 problem")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	flag.Parse()
	verbose = *verboseFlag
	infile := "input.txt"
	lines := fileToLines(infile)

 	start := image.Pt(1, 1)
	dest := image.Pt(31, 39)
	if flag.NArg() > 0 && flag.Arg(0) != "input.txt" {
		dest = image.Pt(7,4)
	}
	dfn := atoi(lines[0])		// the Designer Favorite Number (input.txt)
	
	var result int
	if *partOne {
		VP("Running Part1")
		result = part1(start, dest, dfn)
	} else if *partThree {
		VP("Running Part3")
		result = part3(start, dfn)
	} else {
		VP("Running Part2")
		result = part2(start, dest, dfn)
	}
	fmt.Println(result)
}

//////////// Part 1

func part1(start, dest image.Point, dfn int) int {
	floor := FloorPlan(dfn)
	startTime := time.Now().UnixMicro()
	path := astar.FindPath[image.Point](floor, start, dest, distance, distance)
	endTime := time.Now().UnixMicro()
	fmt.Printf("Ran in %.0fus\n", float64(endTime - startTime))
	return len(path) - 1
}

//////////// Part 2

const mes = 53					// the size of the edge of the maze square: 50 steps from pos 1
var dirs = [4]int{-mes, 1, mes, -1} // the delta in pos for the 4 dirs N E S W

func part2(startp, destp image.Point, dfn int) int {
	// scalar array: int, pos = x + y * mes, values is the number of steps to reach pos
	maze := make([]int, mes*mes, mes*mes)
	start := startp.X + startp.Y * mes
	dest := destp.X + destp.Y * mes
	return step1(maze, []int{start}, 1, dest, 0, maxsteps+1, dfn)
}

// BFS search
func step1(maze, todo []int, dist, dest, count, cap, dfn int) int {
	if dist <= cap {			//  do not count visited
		count += len(todo)
	}
	nexts := []int{}
	for _, p := range todo {
		maze[p] = dist
		for _, n := range possibleMoves(maze, p, dfn) {
			if n == dest {
				return count
			}
			if indexOfInt(nexts, n) == -1 {
				nexts = append(nexts, n)
			}
		}
	}
	return step1(maze, nexts, dist+1, dest, count, cap, dfn)	
}

// to move in dir, we check that we did not cross the maze edges,
// that we haven't yet step1-ed into this square, and there is no wall
func possibleMoves(maze []int, pos, dfn int) (l []int) {
	var p int
	p = pos - mes				// Up
	if p >= 0 && p < mes*mes && maze[p] == 0 && freeAt(p%mes, p/mes, dfn) {
		l = append(l, p)
	}
	p = pos + mes 				// Down
	if p >= 0 && p < mes*mes && maze[p] == 0 && freeAt(p%mes, p/mes, dfn) {
		l = append(l, p)
	}
	p = pos + 1					// Right, check we are on same row (no edge wrap)
	if p/mes == pos/mes && maze[p] == 0 && freeAt(p%mes, p/mes, dfn) {
		l = append(l, p)
	}
	p = pos - 1					// Left, check we are on same row (no edge wrap)
	if p/mes == pos/mes && maze[p] == 0 && freeAt(p%mes, p/mes, dfn) {
		l = append(l, p)
	}
	return
}

//////////// Common Parts code

//// callbacks on graph Nodes for use by fzipp/astar: distance and neighbours

// distance is our cost function, the Manhattan distance.
func distance(p, q image.Point) float64 {
	d := q.Sub(p)
	return float64(intAbs(d.X) + intAbs(d.Y))
}

// Neighbours implements the astar.Graph interface
func (f FloorPlan) Neighbours(p image.Point) []image.Point {
	offsets := []image.Point{
		image.Pt(0, -1), // North
		image.Pt(1, 0),  // East
		image.Pt(0, 1),  // South
		image.Pt(-1, 0), // West
	}
	res := make([]image.Point, 0, 4)
	for _, off := range offsets {
		q := p.Add(off)
		if f.isFreeAt(q) {
			res = append(res, q)
		}
	}
	return res
}

func (f FloorPlan) isFreeAt(p image.Point) bool {
	sum := p.X * p.X + 3 * p.X + 2 * p.X * p.Y + p.Y + p.Y * p.Y + int(f)
	bitstring := strconv.FormatInt(int64(sum),2)
	bitsum := 0
	for _, b := range bitstring {
		if b == '1' {
			bitsum++
		}
	}
	return p.X >= 0 && p.Y >= 0 && (bitsum % 2 == 0)
}
		
//////////// Part1 functions

//////////// Part2 functions

// lighter version of FloorPlan.freeAt for standalone use

func freeAt(x, y, dfn int) bool {
	sum := x * x + 3 * x + 2 * x * y + y + y * y + dfn
	bitstring := strconv.FormatInt(int64(sum),2)
	bitsum := 0
	for _, b := range bitstring {
		if b == '1' {
			bitsum++
		}
	}
	return x >= 0 && y >= 0 && (bitsum % 2 == 0)
}	

//////////// Part 3
func part3(start image.Point, dfn int) int {
	floor := FloorPlan(dfn)
	locations := 0
	startTime := time.Now().UnixMicro()
	x0 := start.X
	y0 := start.Y
	for x := 0; x <= x0 + maxsteps; x++ {
		for y := 0; y <= y0 + maxsteps - x; y++ {
			if freeAt(x, y, dfn) {
				path := astar.FindPath[image.Point](floor, start, image.Pt(x, y), distance, distance)
				if (len(path) - 1) <= maxsteps {
					locations++
				} else {
					VPf("failure for (%d, %d)\n", x, y)
				}
			}
		}
	}
	endTime := time.Now().UnixMicro()
	fmt.Printf("Ran in %.0fus\n", float64(endTime - startTime))
	return locations
}
