// Adventofcode 2024, d06, in go. https://adventofcode.com/2024/day/06
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 41
// TEST: example 6
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"flag"
	"fmt"
)

// if guard lands at pos gp and at direction gd (pos offset), then:
//   grid.path[gd] = true
//   grid.dirs[index-of(gd)].a[gd] = true
//   where index-of returns 0,1,2,3 from a pos offset dir
// this one-bool-array-per-dir is an alternative to a single array with
// bitfields for storing dir tracks

type Grid struct {				// the problem data world
	labo, path Scalarray[bool]	// lab obstruction map and positions visited
	gp, gd int					// guard position and direction
	dirs []Scalarray[bool]		// for part2: guard 4 directions on path NESW
}

var verbose, debug bool

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	debugFlag := flag.Bool("V", false, "debug: even more verbose")
	flag.Parse()
	verbose = *verboseFlag
	debug = *debugFlag
	var infile string
	if flag.NArg() > 0 {
		infile = flag.Arg(0)
	} else {
		infile = fileMatch("input,[[:alnum:]]*,[[:alnum:]]*.test")
	}
	lines := fileToLines(infile)
	if *partOne {
		VP("Running Part1")
		fmt.Println(part1(lines))
	} else {
		VP("Running Part2")
		fmt.Println(part2(lines))
	}
}

//////////// Part 1

func part1(lines []string) int {
	grid := parse(lines)
	for grid.Step() {}
	return VisitsCount(grid)
}

func (grid *Grid) Step() bool {
	gnp := grid.gp + grid.gd	// next position for the guard
	if ! grid.labo.StepDirInside(grid.gp, grid.gd) {	// left the lab
		return false
	}
	if grid.labo.a[gnp] {		// bumps into an obstacle, turn right
		gnd := grid.labo.RotateDir(grid.gd, 90)
		grid.gd = gnd
		gnp = grid.gp			//  but stays in place
	} else {
		grid.gp = gnp			// moves ahead
		grid.path.a[gnp] = true // marks new pos into path
	}
	return true					// ready for next step
}

func VisitsCount(grid *Grid) (visits int) {
	for _, v := range grid.path.a {
		if v {
			visits++
		}
	}
	return
}

//////////// Part 2

func part2(lines []string) (loops int) {
	grid := parse(lines)
	gp := grid.gp				// remember start position for reset
	gd := grid.gd
	for p := range grid.labo.a {
		grid.gp, grid.gd = gp, gd // reset grid guard pos & path
		GridPathInit(grid)
		if ObstacleCreatesLoop(grid, p) {
			loops++
		}
	}
	return
}

// direction offset => 0, 1, 2, 3 for N, E, S, W, index of grid.dirs[]
func DirIndex(grid *Grid, d int) int {
	return grid.labo.DirToDegrees(d) / 90
}

// does adding an obstacle at pos p creates a loop?
func ObstacleCreatesLoop(grid *Grid, p int) (loop bool) {
	if grid.gp == p || grid.labo.a[p] {
		return					// if p is occupied by guard or obstacle, skip
	}
	grid.labo.a[p] = true		// place obstacle and test a run
	var ok bool
	for {
		ok, loop = grid.StepCheckLoop()
		if ! ok {		// guard lefts the lab
			break
		}
		if loop {
			break
		}
	}
	grid.labo.a[p] = false		// resets lab obstacles map
	return
}

// returns: ok-to-continue?, loop-detected?
func  (grid *Grid) StepCheckLoop() (bool, bool) {
	gnp := grid.gp + grid.gd	// next position for the guard
	if ! grid.labo.StepDirInside(grid.gp, grid.gd) {	// left the lab
		return false, false
	}
	if grid.labo.a[gnp] {		// bumps into an obstacle, turn right
		gnd := grid.labo.RotateDir(grid.gd, 90)
		grid.gd = gnd
		gnp = grid.gp			//  but stays in place
		return true, false
	} else {
		if grid.path.a[gnp] && grid.dirs[DirIndex(grid, grid.gd)].a[gnp] {
			return false, true	// guard already went through in same dir
		}
		grid.gp = gnp			// moves ahead
		grid.path.a[gnp] = true // marks new pos into path
		grid.dirs[DirIndex(grid, grid.gd)].a[gnp] = true
		return true, false
	}
}


func GridPathInit(grid *Grid) {
	grid.path = grid.labo.New()
	grid.dirs = []Scalarray[bool]{ // for part2, init this optional field
		grid.labo.New(), grid.labo.New(), grid.labo.New(), grid.labo.New(),
	}
	grid.path.a[grid.gp] = true
	grid.dirs[DirIndex(grid, grid.gd)].a[grid.gp] = true
}

//////////// Common Parts code

func parse(lines []string) (*Grid) {
	labo := makeScalarray[bool](len(lines[0]), len(lines))
	path := makeScalarray[bool](len(lines[0]), len(lines))
	gd := - labo.w				// guard direction is up
	var gp int
	for y, line := range lines {
		for x, b := range line {
			switch b {
			case '#' : labo.Set(x, y, true) // block
			case '^' :						// guard pos
				gp = path.Pos(x, y)
				path.Set(x, y, true)
			}
		}
	}
	grid := Grid{labo, path, gp, gd, []Scalarray[bool]{}}
	return &grid
}

//////////// PrettyPrinting & Debugging functions

func DEBUG() {
	if ! debug { return }
	// ad hoc debug code here
}
