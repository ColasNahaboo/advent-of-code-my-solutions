// Adventofcode 2023, d10, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example1s 4   // 1st example, simplified
// TEST: -1 example1f 4   // 1st example, full
// TEST: -1 example2s 8   // 2nd example, simplified
// TEST: -1 example2f 8   // 2nd example, full
// TEST: example-p2e1 4   // 1st example of part2 description text 
// TEST: example-p2e2 4   // 2nd example of part2 description text 
// TEST: example-p2e3 8   // 3rd example of part2 description text 
// TEST: example-p2e4 10   // 4th example of part2 description text 
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"flag"
	"fmt"
	//"regexp"
)

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

// We use global variables for simplicity rather than having to pass the
// state of the "world" around all the time.
//
// No real parsing, we just concatenate all the input strings (the rows) into
// one single byte array (a slice, actually), "grid" of size (gw * gh), with
// gw being the length of input lines, and gh their number.
// A position (x, y) is thus just an integer index x+gw*y in this array
// We also replace the character 'S' with the actual pipe underneath
var grid []byte				// the grid, gw x gh, the raw input, with S replaced
var gw, gh, area int		// grid size
var start int				// start position, pos of 'S'
var dirs [4]int				// the 4 directions as position offsets
const (
	UP = 1
	RIGHT = 2
	DOWN = 4
	LEFT = 8
)

//////////// Part 1

// the max distance is just half the longest path length, rounded up

func part1(lines []string) int {
	path := findLoop(lines)
	return (len(path)+1)/2
}

//////////// Part 2

// A point is inside the loop if we must cross the loop an odd number of times
// to reach a border. For simplicity we chose to just draw a line upwards
// instead of to the closest border
// The trick is then to not count | as a crossing, and choose only one horizontal
// side (i.e left: "7 and J" or right: "F and L") to not count twice the same
// crossing. Here we chose right, "F and L".

func part2(lines []string) int {
	path := findLoop(lines)
	// mark all the places occupied by the loop
	board := make([]bool, area, area)
	inside := make([]bool, area, area) // DEBUG
	for _, pos := range path {
		board[pos] = true
	}
	// find places that can go to the top border crossing the loop an odd times
	insiders := 0
	for p, _ := range board {
		if board[p] {			// on the loop itself means not inside
			continue
		}
		crosses := 0
		for i := p; i >= 0; i -= gw {
			// we are moving up, so moving along '|' does not count as a crossing
			if board[i] {
				crosses += isCrossing(i)
			}
		}
		if crosses % 2 == 1 {
			insiders++
			inside[p] = true
		}
	}
	VPf("Start pipe is a: \"%s\"\n", string(grid[start]))
	VPboard(board, inside)
	return insiders
}

// we count the pipes with only an horizontal part connecting to the right
// otherwise 7 above a L would count as a double crossing
func isCrossing(pos int) int {
	switch grid[pos] {
	case '-', 'L', 'F': return 1
	default: return 0
	}
}

//////////// Common Parts code

// No real parsing, we just concatenate all the input strings (the rows) into
// one single byte array (a slice, actually), the grid of size (gw x gh).
// A position (x, y) is thus just an integer index x+gw*y in this array

func findLoop(lines []string) (lpath []int) {
	gw = len(lines[0])
	gh = len(lines)
	area = gw*gh
	grid = make([]byte, area, area)
	for y, line := range lines {
		for x, c := range line {
			grid[y*gw+x] = byte(c)
		}
	}
	dirs = [4]int{-gw, +1, gw, -1}
	start = startPos()
	VPf("Grid %d x %d, start at %d\n", gw, gh, start)
	ns := neighbours(start)
	for _, n := range ns {
		path := followLoop(start, n)
		if len(path) > len(lpath) {
			lpath = path
		}
	}
	grid[start] = pipeConnecting(start, lpath[0], lpath[len(lpath) - 2])
	return
}

func startPos() int {
	for p, c := range grid {
		if c == 'S' {
			return p
		}
	}
	panic("S not found")
}

// follows a (potential) loop completely
func followLoop(pos, next int) (path []int) {
	for {
		path = append(path, next)
		if next == start {		// we went back to start
			return
		}
		from := pos
		pos = next
		next = nextStep(pos, from)
		if next == -1 {			// not a loop!
			path = []int{}
			return
		}
	}
	return
}

// follows the loop for one step from "pos", not going back via "from"
// we suppose we have one and only one possibility
// returns -1 if no neigbours can be found
func nextStep(pos, from int) (next int) {
	ns := neighbours(pos)
	for _, n := range ns {
		if n == from {
			continue
		}
		next = n
		return
	}
	next = -1
	return
}

// List the possible neighbours. S has all 4 dirs, other pipes only 2.
func neighbours(pos int) []int {
	switch grid[pos] {
	case '|': return validPair(pos - gw, pos + gw)
	case '-': return validPair(pos - 1, pos + 1)
	case 'L': return validPair(pos - gw, pos + 1)
	case 'J': return validPair(pos - gw, pos - 1)
	case '7': return validPair(pos - 1, pos + gw)
	case 'F': return validPair(pos + 1, pos + gw)
	case 'S':					// S is special, be tolerant
		ns := []int{}
		for _, d := range dirs {
			p := pos + d
			if d >= 0 && d < area {
				ns = append(ns, p)
			}
		}
		return ns
	default: return []int{}
	}
}

func validPair(i, j int) (pair []int) {
	if i >= 0 && i < area && j >=0 && j < area {
		pair = []int{i, j}
	}
	return
}

// return the symbol of the pipe at p that connects neighbours to and from
func pipeConnecting(p int, ns ...int) byte {
	sides := 0
	for _, n := range ns {
		switch n - p {
		case -gw: sides |= UP
		case 1: sides |= RIGHT
		case gw: sides |= DOWN
		case -1: sides |= LEFT
		default: panic(fmt.Sprintf("%d is not a neighbour of %d", n, p))
		}
	}
	switch sides {
	case 5: return '|'
	case 10: return '-'
	case 3: return 'L'
	case 9: return 'J'
	case 12: return '7'
	case 6: return 'F'
	default: panic(fmt.Sprintf("Cannot connect %d to %v", p, ns))
	}
}

//////////// Debug code

func VPboard(board, inside []bool) {
	if ! verbose {
		return
	}
	for pos, inside := range inside {
		if inside {
			fmt.Print("#")
		} else if board[pos] {
			fmt.Print(string(grid[pos]))
		} else {
			fmt.Print(".")
		}
		if pos % gw == gw - 1 {
			fmt.Println()
		}
	}
}
