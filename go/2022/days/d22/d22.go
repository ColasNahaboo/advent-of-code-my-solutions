// Adventofcode 2022, d22, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 example 6032
// TEST: -1 input 126350
// TEST: example 5031
// TEST: input 129339                         // 129338 too low, 120241 too high
package main

import (
	"flag"
	"fmt"
	"regexp"
)

var gw, gh, garea int			// grid dims. Indexes of (x,y) in 1D arrays is x+y*gw
var grid []bool					// including a border mirroring other side
var board []bool			    // grid points defined in input ('.' or '#' but not ' ')
var program string				// original string
var instructions []int			// direction (R=0, L=-1) or distance (> 0) to run
var deltapos [4]int				// for dirs: R=0, D=1, L=2, U=3 positions increments
var mapCubeFile string			// part2: name of the cube mapping file
var edgeCube int				// part2: the length of cube edges
var sides [][2]int				// the coords of the 6 cube sides on the map
type Edge struct {
	fside, fdir int
	tside, tdir int
	invert int
}
var edges []Edge

var verbose bool

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	testingFlag := flag.Bool("t", false, "debug: enter some tests")
	flag.Parse()
	verbose = *verboseFlag
	if *testingFlag { testrun(); return;}
	infile := "input.txt"
	if flag.NArg() > 0 {
		infile = flag.Arg(0)
	}
	mapCubeFile = infile[0:len(infile)-3] + "map"
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
	parse(lines)
	return run(startPos(), 0, instructions, false)
}

//////////// Part 2
func part2(lines []string) int {
	parse(lines)
	parse3Dmap(mapCubeFile)
	return run(startPos(), 0, instructions, true)
}

//////////// Common Parts code

func parse(lines []string) {
	for h, line := range lines {
		if len(line) > gw { gw = len(line); }
		if len(line) == 0 { gh = h; break; }
	}
	gw += 2
	gh += 2
	garea = gw * gh
	program = lines[gh-1]
	grid = make([]bool, garea, garea)
	board = make([]bool, garea, garea)

	for y := 1; y < gh-1; y++ {
		for x := 1; x < len(lines[y-1]) + 1; x++ {
			if lines[y-1][x-1] == byte('#') {grid[x + y*gw] = true;}
			if lines[y-1][x-1] != byte(' ') {board[x + y*gw] = true;}
		}
	}

	// compute position changes for one step in each direction
	deltapos[0] = 1				// R
	deltapos[1] = gw			// D
	deltapos[2] = -1			// L
	deltapos[3] = -gw			// U

	// now parse program
	re := regexp.MustCompile("([RL]|[[:digit:]]+)")
	for _, s := range re.FindAllString(program, -1) {
		if s == "R" {
			instructions = append(instructions, 0)
		} else if s == "L" {
			instructions = append(instructions, -1)
		} else {
			instructions = append(instructions, atoi(s))
		}
	}
	VPf("grid %dx%d(%d), program: %d instructions\n", gh, gw, garea, len(instructions))
}

func startPos() (sp int) {
	for x := 1; x < gw-1; x++ {	// first free pos, not on a wrap-border
		if !grid[x + 1*gw] && board[x + 1*gw] {
			sp = x + 1*gw
			break
		}
	}
	return
}

// fill a slice with int
func fillInt(s *[]int, v int) {
	for i := 0; i < len(*s); i++ {
		(*s)[i] = v
	}
}

// position to coords
func p2xy(p int) (x,y int) {
	x = p % gw
	y = p / gw
	return
}

// coord to position
func xy2p(x, y int) int {
	return x +y*gw
}

// run program on grid, starting at p facing d
func run(p, dir int, code []int, cube bool) int {
	VPf("Start run:  at pos %4d, x=%d, y=%d, facing: %d (cube: %d)\n", p, (p % gw), (p / gw), dir, edgeCube)
	var dir2 int
	for _, inst := range code {
		if inst > 0 {			// move straight ahead
			for i := 0; i < inst; i++ {
				p2 := p + deltapos[dir]
				dir2 = dir
				if !board[p2] {	// oops, outside grid, wrap
					if cube {
						p2, dir2 = wrapCube(p, dir)
					} else {
						p2 = wrapFlat(p, dir)
					}
				}
				if grid[p2] {	// wall, abort, stay in place
					break
				}
				p = p2			// OK, move or teleport
				dir = dir2
				//VPf("    (%d, %d) %d\n", p % gw - 1, p / gw - 1, dir)
			}
		} else if inst == 0 {	// R
			dir = (dir + 1) % 4
		} else {				// L
			dir = (dir + 3) % 4	// dir - 1, but positive
		}
		if verbose {tracepos(p, dir, inst);}
	}
	VPf("End of run: at pos %4d, x=%d, y=%d, facing: %d\n", p, (p % gw), (p / gw), dir)
	return 1000 * (p / gw) + 4 * (p % gw) + dir
}

// naive square root for integers
func squareRoot(i int) (r int) {
	for r = 1; r*r < i; r++ {}
	if r*r != i { panic(fmt.Sprintf("Error: not a square: %d\n", i))}
	return
}		

//////////// Part1 functions

// On the edge at pos, return pos after step once in dir
func wrapFlat(pos, dir int) (p int) {
	switch dir {
	case 0: p = (pos / gw) * gw
	case 1: p = pos % gw
	case 2: p = (pos / gw) * gw + (gw - 1)
	case 3: p = pos % gw + (gh - 1) * gw
	}
	for !board[p] { p+= deltapos[dir];}
	return
}
	
//////////// Part2 functions

// sidemaps
// example.map  input.map
//   1               12
// 234               3
//   56             45
//                  6

// for now the mapfilename is not really read, is used as an index in precomputed data
func parse3Dmap(mapfilename string) {
	var boardsize int				// number of actual squares defined
	for x := 0; x < gw; x++ {
		for y := 0; y < gh; y++ {
			if board[x + y*gw] { boardsize++}
		}
	}
	if boardsize % 6 != 0 {
		panic(fmt.Sprintf("Error: boardsize not multiple of 6: %d\n", boardsize))
	}
	edgeCube = squareRoot(boardsize / 6)
	if mapfilename == "example.map" {
		// TODO: implement building sides and edges by the actual parsing of the files
		// sides[i], for i in [1..6] gives coords of sides in the sidemap
		sides = [][2]int{{-1,-1}, {2,0}, {0,1}, {1,1}, {2,1}, {2,2}, {3,2}}
		// edges for each side and out dir, give target sizde and in dir + must invert edge pos?
		edges = []Edge{
			Edge{1, 0, 6, 2, 1},
			Edge{1, 2, 3, 1, 0},
			Edge{1, 3, 2, 1, 1},
			Edge{2, 1, 5, 3, 1},
			Edge{2, 2, 6, 3, 1},
			Edge{2, 3, 1, 1, 1},
			Edge{3, 1, 5, 0, 1},
			Edge{3, 3, 1, 0, 0},
			Edge{4, 0, 6, 1, 1},
			Edge{5, 1, 2, 3, 1},
			Edge{5, 2, 3, 3, 1},
			Edge{6, 0, 1, 2, 1},
			Edge{6, 1, 2, 0, 1},
			Edge{6, 3, 4, 2, 1},
		}
	} else if mapfilename == "input.map" {
		sides = [][2]int{{-1,-1}, {1,0}, {2,0}, {1,1}, {0,2}, {1,2}, {0,3}}
		edges = []Edge{
			Edge{1, 2, 4, 0, 1},
			Edge{1, 3, 6, 0, 0},
			Edge{2, 0, 5, 2, 1},
			Edge{2, 1, 3, 2, 0},
			Edge{2, 3, 6, 3, 0},
			Edge{3, 0, 2, 3, 0},
			Edge{3, 2, 4, 1, 0},
			Edge{4, 2, 1, 0, 1},
			Edge{4, 3, 3, 0, 0},
			Edge{5, 0, 2, 2, 1},
			Edge{5, 1, 6, 2, 0},
			Edge{6, 0, 5, 3, 0},
			Edge{6, 1, 2, 1, 0},
			Edge{6, 2, 1, 1, 0},
		}
	}
		
}

// On the edge at pos, return pos after step once in dir
func wrapCube(pos, dir int) (p, d int) {
	// TODO: find face to enter, from which side, and which orientation
	// (relative pos on the edge, starting from topleft)
	var side, s, invert, epos, x, y int
	side = pos2side(pos)
	s, d, invert = wrapEdge(side, dir)
	// exiting edge: compute epos, distance of pos from the edge end of lowest coord 
	if dir % 2 == 0 {			// horiz dir: edge is vertical, take lowest y
		y0 := (sides[side][1] * edgeCube + 1) * gw // top y of side
		epos = (pos - y0) / gw
	} else {					// vert dir: edge is horiz, take lowest x
		x0 := sides[side][0] * edgeCube + 1 // left x of side
		epos = pos % gw - x0
	}
	// now compute p as pos on to edge at distance epos from lowest or biggest (invert=1) corner
	if d % 2 == 0 {				// entering side s via horiz dir
		if invert == 0 {
			y = epos + sides[s][1] * edgeCube + 1
		} else {
			y = - epos + (sides[s][1] + 1) * edgeCube
		}
		if d == 0 {				// entering from left
			x = sides[s][0] * edgeCube + 1
		} else {				// entering from right
			x = (sides[s][0] + 1) * edgeCube
		}
	} else {					// entering s via vert dir
		if invert == 0 {
			x = epos + sides[s][0] * edgeCube + 1
		} else {
			x = - epos + (sides[s][0] + 1) * edgeCube
		}
		if d == 1 {				// entering from top
			y =  sides[s][1] * edgeCube + 1
		} else {				// entering from bottom
			y = (sides[s][1] + 1) * edgeCube
		}
	}
	p = x + y*gw
	return
}

// find id [0..6] of cube side we are in
func pos2side(pos int) int {
	x := (pos % gw - 1) / edgeCube
	y := (pos / gw - 1) / edgeCube
	for side, xy := range sides {
		if x == xy[0] && y == xy[1] {
			return side
		}
	}
	panic(fmt.Sprintf("Cannot find side {%d, %d} of (%d, %d), pos=%d, gw=%d, cube=%d\n", x, y, pos % gw - 1, pos / gw - 1, pos, gw, edgeCube))
}

// which side are we entering, and how, after wrapping out of a border?
func wrapEdge(inside, indir int) (int, int, int) {
	for _, edge := range edges {
		if inside == edge.fside && indir == edge.fdir {
			return edge.tside, edge.tdir, edge.invert
		}
	}
	panic(fmt.Sprintf("Cannot find edge for side %d and dir %d\n", inside, indir))
}

////// Some temporary testing code copied from d22_test.go, for using under vld
// beware not using the same function names: t_wc ==> twc
func testrun() {
	lines := fileToLines("example.txt")
	parse(lines)
	parse3Dmap("example.map")
	t := 0						// dummy
	twc("1ul", t,  9, 1, 3,  4,  5, 1)
}

func twc(label string, t int, x1, y1, d1, x2, y2, d2 int) {
	p1 := x1 + y1*gw
	p2 := x2 + y2*gw
	p, d := wrapCube(p1, d1)
	if p != p2 || d != d2 {
		panic(fmt.Sprintf("expected %d:%d (%d,%d) but got %d:%d (%d,%d)", p2, d2, x2, y2, p, d, p % gw, p / gw))
	}
}

func tracepos(p, d, inst int) {
	if edgeCube == 0 {return;}	// not in 3D mode
	var i string
	if inst > 0 {
		i = itoa(inst)
	} else if inst == 0 {
		i = " R"
	} else {
		i = " L"
	}
	side := pos2side(p)
	// verbose: instr: face-coord (face-relative-coord) (absolute-coords) dir
	fmt.Printf("%s: [%d] (%d, %d) (%d, %d) %d\n", i, side, p % gw - 1 - sides[side][0]*edgeCube, p / gw - 1 - sides[side][1]*edgeCube, p % gw - 1, p / gw - 1, d)
}
