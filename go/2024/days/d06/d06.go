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
	"slices"
)

// if guard lands at pos gp and at direction gd (pos offset), then:
//   grid.path[gd] = true
//   grid.dirs[index-of(gd)].a[gd] = true
//   where index-of returns 0,1,2,3 from a pos offset dir
// this one-bool-array-per-dir is an alternative to a single array with
// bitfields for storing dir tracks. called with -2

// Part3 (-3) is an alternative with a simpler implementation of path as ints,
// bitwise-OR of dirs, which consumes a bit more memory... but is twice faster!

// Part4 (-4) is an alternative with a simpler implementation of path as bytes,
// bitwise-OR of dirs, which is both smaller and faster.

// Part5 (-5) is an alternative implementation with point instead of scalarray
// Simpler, ans as fast as -4
// So we make it the default


type Grid struct {				// the problem data world
	labo, path Scalarray[bool]	// lab obstruction map and positions visited
	gp, gd int					// guard position and direction
	dirs []Scalarray[bool]		// for part2: guard 4 directions on path NESW
}

var verbose, debug bool

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	partTwo := flag.Bool("2", false, "run exercise part2, alt. part2")
	partThree := flag.Bool("3", false, "run exercise part3, alt. part2")
	partFour := flag.Bool("4", false, "run exercise part4, alt. for part2")
	partFive := flag.Bool("5", false, "run exercise part5, default for part2")
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
	} else if *partTwo {
		VP("Running Part2")
		fmt.Println(part2(lines))
	} else if *partThree {
		VP("Running Part3")
		fmt.Println(part3(lines))
	} else if *partFour {
		VP("Running Part4")
		fmt.Println(part4(lines))
	} else if *partFive {
		VP("Running Part5")
		fmt.Println(part5(lines))
	} else {
		VP("Running Part5")
		fmt.Println(part5(lines))
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
	if grid.labo.a[gnp] {		// bumps into an obstacle, turn right in place
		gnd := grid.labo.RotateDir(grid.gd, 90)
		grid.gd = gnd
		grid.dirs[DirIndex(grid, gnd)].a[gnp] = true // marks new dir taken
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

//////////// Part3
// with a simpler implementation of path as ints, bitwise-OR of dirs

type Grid3 struct {				// the problem data world
	labo Scalarray[bool]	    // lab obstruction map
	path Scalarray[int]         // positions and dirs visited, OR of DN DE DS DW
	gp, gd int					// guard position and direction
}

const (
	DN = 1
	DE = 2
	DS = 4
	DW = 8
)

func part3(lines []string) (loops int) {
	grid2 := parse(lines)
	gp := grid2.gp				// remember start position for reset
	gd := grid2.gd
	grid := &Grid3{grid2.labo, makeScalarray[int](grid2.labo.w, grid2.labo.h), gp, gd}
	for p := range grid.labo.a {
		grid.gp, grid.gd = gp, gd // reset grid guard pos & path
		Grid3PathInit(grid)
		if Grid3ObstacleCreatesLoop(grid, p) {
			loops++
		}
	}
	return
}

func Grid3DirToDX(grid *Grid3, d int) int {
	switch d {
	case - grid.labo.w: return DN
	case 1: return DE
	case grid.labo.w: return DS
	case -1: return DW
	}
	return 0
}

func Grid3PathInit(grid *Grid3) {
	grid.path = grid.path.New()
	grid.path.a[grid.gp] = Grid3DirToDX(grid, grid.gd)
}

// does adding an obstacle at pos p creates a loop?
func Grid3ObstacleCreatesLoop(grid *Grid3, p int) (loop bool) {
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
func  (grid *Grid3) StepCheckLoop() (bool, bool) {
	gnp := grid.gp + grid.gd	// next position for the guard
	if ! grid.labo.StepDirInside(grid.gp, grid.gd) {	// left the lab
		return false, false
	}
	if grid.labo.a[gnp] {		// bumps into an obstacle, turn right in place
		gnd := grid.labo.RotateDir(grid.gd, 90)
		grid.gd = gnd
		grid.path.a[grid.gp] |= Grid3DirToDX(grid, gnd)
		return true, false
	} else {
		dx := Grid3DirToDX(grid, grid.gd)
		if grid.path.a[gnp] & dx != 0 {
			return false, true	// guard already went through in same dir
		}
		grid.gp = gnp			// moves ahead
		grid.path.a[gnp] = dx // marks new pos into path
		return true, false
	}
}

//////////// Part4
// with a simpler implementation of path as bytes, bitwise-OR of dirs

type Grid4 struct {				// the problem data world
	labo Scalarray[bool]	    // lab obstruction map
	path Scalarray[byte]        // positions and dirs visited, OR of DN DE DS DW
	gp, gd int					// guard position and direction
}

func part4(lines []string) (loops int) {
	grid2 := parse(lines)
	gp := grid2.gp				// remember start position for reset
	gd := grid2.gd
	grid := &Grid4{grid2.labo, makeScalarray[byte](grid2.labo.w, grid2.labo.h), gp, gd}
	for p := range grid.labo.a {
		grid.gp, grid.gd = gp, gd // reset grid guard pos & path
		Grid4PathInit(grid)
		if Grid4ObstacleCreatesLoop(grid, p) {
			loops++
		}
	}
	return
}

func Grid4DirToDX(grid *Grid4, d int) byte {
	switch d {
	case - grid.labo.w: return DN
	case 1: return DE
	case grid.labo.w: return DS
	case -1: return DW
	}
	return 0
}

func Grid4PathInit(grid *Grid4) {
	grid.path = grid.path.New()
	grid.path.a[grid.gp] = Grid4DirToDX(grid, grid.gd)
}

// does adding an obstacle at pos p creates a loop?
func Grid4ObstacleCreatesLoop(grid *Grid4, p int) (loop bool) {
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
func  (grid *Grid4) StepCheckLoop() (bool, bool) {
	gnp := grid.gp + grid.gd	// next position for the guard
	if ! grid.labo.StepDirInside(grid.gp, grid.gd) {	// left the lab
		return false, false
	}
	if grid.labo.a[gnp] {		// bumps into an obstacle, turn right in place
		gnd := grid.labo.RotateDir(grid.gd, 90)
		grid.gd = gnd
		grid.path.a[grid.gp] |= Grid4DirToDX(grid, gnd)
		return true, false
	} else {
		dx := Grid4DirToDX(grid, grid.gd)
		if grid.path.a[gnp] & dx != 0 {
			return false, true	// guard already went through in same dir
		}
		grid.gp = gnp			// moves ahead
		grid.path.a[gnp] = dx // marks new pos into path
		return true, false
	}
}

//////////// Part5  with point.go Point and Board

type Grid5 struct {			   // the problem data world
	labo Board[bool]		   // lab obstruction map
	path Board[byte]		   // positions and dirs visited, OR of DN DE DS DW
	gp Point				   // guard position and direction
	gd Point
}

func part5(lines []string) (loops int) {
	g := Grid5{}
	var gp, gd Point
	g.labo = *parseBoard[bool](lines, func(x, y int, r rune) bool {
		if r == '#' {
			return true
		} else if r == '^' {
			gp = Point{x, y}
			gd = DirsOrtho[DirsOrthoN]
		}
		return false
	})
	g.path = makeBoard[byte](g.labo.w, g.labo.h)
	clearcol := g.path.ClearInit() // used as seed for later Clear()
	
	for x := range g.labo.w {
		for y := range g.labo.h {
			g.gp, g.gd = gp, gd // reset grid guard pos & path
			g.path.Clear(clearcol) // fast fill of path[][] with zeros
			g.path.a[gp.x][gp.y] = DN
			if Grid5ObstacleCreatesLoop(&g, Point{x, y}) {
				loops++
			}
		}
	}
	return
}

// does adding an obstacle at pos p creates a loop?
func Grid5ObstacleCreatesLoop(grid *Grid5, p Point) (loop bool) {
	if grid.gp == p || grid.labo.a[p.x][p.y] {
		return					// if p is occupied by guard or obstacle, skip
	}
	grid.labo.a[p.x][p.y] = true		// place obstacle and test a run
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
	grid.labo.a[p.x][p.y] = false		// resets lab obstacles map
	return
}

// returns: ok-to-continue?, loop-detected?
func  (grid *Grid5) StepCheckLoop() (bool, bool) {
	gnp := grid.gp.Add(grid.gd)	// next position for the guard
	if ! grid.labo.Inside(gnp) {	// left the lab
		return false, false
	}
	if grid.labo.a[gnp.x][gnp.y] { // bumps into an obstacle, turn right in place
		gnd := grid.gd.RotateDirOrtho(1)
		grid.gd = gnd
		grid.path.a[grid.gp.x][grid.gp.y] |= Grid5DirToDX(gnd)
		return true, false
	} else {
		dx := Grid5DirToDX(grid.gd)
		if grid.path.a[gnp.x][gnp.y] & dx != 0 {
			return false, true	// guard already went through in same dir
		}
		grid.gp = gnp			// moves ahead
		grid.path.a[gnp.x][gnp.y] |= dx // marks new pos into path
		return true, false
	}
}

func Grid5DirToDX(d Point) byte {
	i := slices.Index(DirsOrtho, d)
	if i == -1 {
		panic(fmt.Sprintf("Not a DirsOrtho: %v", d))
	}
	return [4]byte{DN, DE, DS, DW}[i]
}

//////////// PrettyPrinting & Debugging functions

func DEBUG() {
	if ! debug { return }
	// ad hoc debug code here
}
