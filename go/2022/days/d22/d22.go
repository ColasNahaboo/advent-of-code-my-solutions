// Adventofcode 2022, d22, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 example 6032
// TEST: -1 input 126350
// TEST: example
// TEST: input
package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"
)

var gw, gh, garea int			// grid dims. Indexes of (x,y) in 1D arrays is x+y*gw
var grid []bool					// including a border mirroring other side
var board []bool			    // grid points defined in input ('.' or '#' but not ' ')
var program string				// original string
var instructions []int			// either direction (<= 0) or distance (> 0) to run
var deltapos [4]int				// for dirs: R=0, D=1, L=2, U=3 positions increments
var mapCubeFile string			// part2: name of the cube mapping file
var edgeCube int				// part2: the length of cube edges
var sides [][2]int				// the coords of the 6 cube sides on the map

var verbose bool

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	flag.Parse()
	verbose = *verboseFlag
	infile := "input.txt"
	mapCubeFile = infile[0:len(infile)-3] + "map"
	if flag.NArg() > 0 {
		infile = flag.Arg(0)
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
	parse(lines)
	return run(startPos(), 0, false)
}

//////////// Part 2
func part2(lines []string) int {
	parse(lines)
	parse3Dmap(mapCubeFile)
	return run(startPos(), 0, true)
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

// run program on grid, starting at p facing d
func run(p, dir int, cube bool) int {
	VPf("Start run:  at pos %4d, x=%d, y=%d, facing: %d\n", p, (p % gw), (p / gw), dir)
	var dir2 int
	for _, inst := range instructions {
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
			}
		} else if inst == 0 {
			dir = (dir + 1) % 4
		} else {
			dir = (dir + 3) % 4	// dir - 1, but positive
		}
	}
	VPf("End of run: at pos %4d, x=%d, y=%d, facing: %d\n", p, (p % gw), (p / gw), dir)
	return 1000 * (p / gw) + 4 * (p % gw) + dir
}

// naive square root for integers
func squareRoot(i int) (r int) {
	for r = 1; r*r < i; r++ {}
	if r*r != i { log.Fatalf("Error: not a square: %d\n", i)}
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

// example.map  input.map
//   1               12
// 234               3
//   56             45
//                  6

func parse3Dmap(mapfilename string) {
	var boardsize int				// number of actual squares defined
	for x := 0; x < gw; x++ {
		for y := 0; y < gh; y++ {
			if board[x + y*gw] { boardsize++}
		}
	}
	if boardsize % 6 != 0 { log.Fatalf("Error: boardsize not multiple of 6: %d\n", boardsize)}
	edgeCube = squareRoot(boardsize / 6)
	if mapfilename == "example.map" {
		// TODO: implement the actual parsing of the files.. or even auto-detection!
		sides = [][2]int{{-1,-1}, {2,0}, {0,1}, {1,1}, {2,1}, {2,2}, {3,2}}
		
	}
		
}

// On the edge at pos, return pos after step once in dir
func wrapCube(pos, dir int) (p, d int) {
	x := (pos % gw - 1) / edgeCube
	y := (pos / gw - 1) / edgeCube
	for side, xy := range sides {
		if x == xy[0] && y == xy[1] {
			break
		}
		fmt.Println(side)
	}
	// TODO: find face to enter, from which side, and which orientation
	// (relative pos on the edge, starting from topleft)
	for !board[p] { p+= deltapos[dir];}
	return
}
	
