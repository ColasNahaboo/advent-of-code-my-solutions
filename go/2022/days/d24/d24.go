// Adventofcode 2022, d24, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 example 18
// TEST: -1 input 308
// TEST: example 
// TEST: input 908

// This problem is quite interesting, because we can see that the blizzard
// paths are deterministic: they do not depend on the positions of the other
// blizzards nor us. So for each time, the blizzards places are the same! This
// means than we can cache all state of all the blizzards in an array "grids"
// of 2D maps indexed by time only. The state of the system then becomes only
// two integers: the time and our position! This makes "forking" the state to
// explore a new branch quite light.

// Also, the blizzards state loops with a frequency of the Least Common Multiple
// of the inside dims of the grid, so we can pre-compute this fixed number
// of states and index via time modulo the LCM

// Then, we do not need to explore again the
// same pairs (time, pos), so we cache the already explored positions into a 2D
// map for each time in `dones`. Add the small heuristic of trying the
// exploration with the most likely directions first (right, then left), and
// the search ends up being super fast. We do a DFS (Depth First Search) in
// order to come up with a solution as fast as possible, so we can abort the
// subsequent explorations for better ones as early as possible as soon as they
// take more time.

package main

import (
	"flag"
	"fmt"
	// "regexp"
)

// A grid ([]int) is a scalar array of the rows appended to each other
// position of coordinates (x,y) is p = x + y*gw
// grid value: 0=free, 1=wall, bitwise union of blizzards: 2=> 4=v 8=< 16=^
const FREE = 0
const WALL = 1
const R = 2
const D = 4
const L = 8
const U = 16
var glyph = [31]rune{'.', '#', '>', '?', 'v', '?', '2', '?', '<',
	'?', '2', '?', '2', '?', '3', '?', '^',
	'?', '2', '?', '2', '?', '3', '?', '2',
	'?', '3', '?', '3', '?', '4'}
type Dir struct {
	mask int					// bitwise mask: R D L U
	step int					// position delta to move one step in dir
	vertical bool				// is D or U ?
	label rune					// U D L R O
}
var dirs [5]Dir 				// dirs[4] is special: stay in place, "O"
// exploration order of dirs U D R L O to DFS-explore
// we try to go first right, then down, and stay in place last
var exploredirs = []int{3, 1, 0, 2, 4}
var exploredirsrev = []int{2, 0, 1, 3, 4} // for reverse exploring

// here are the globals that stay the same for all the exploration branches
var gw, gh, area int			//  global dims of the grid
var lcm int						// least common multiple of the inside dims
var start, end int				// position of entry and exit
// the blizzard positions are deterministic, so are the same in all branches
// for a given time! So we can cache them in this array indexed by time
var grids [][]int				// grid states at time
// and thus we do not need to re-explore from the same (time,pos) pairs
var dones [][]bool				// positions already done at time
const maxint = 8888888888888888888 // easily identifiable in debug
var mintime = maxint

var verbose bool

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	capFlag := flag.Int("c", 0, "cap the result: do not explore further. 0 = no cap")
	flag.Parse()
	verbose = *verboseFlag
	if *capFlag !=0 { mintime = *capFlag + 1;}
	infile := "input.txt"
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
	donesInit()
	explore(0, start, end, exploredirs)
	return mintime
}

//////////// Part 2
func part2(lines []string) int {
	parse(lines)
	// 1rst trip, start -> end
	time := 0
	donesInit()
	explore(time, start, end, exploredirs)
	fmt.Printf("  run:   %d\n", mintime)
	// 2nd trip, get back to start
	time = mintime
	mintime = maxint
	donesInit()
	explore(time, end, start, exploredirsrev)
	fmt.Printf("  back:  %d\n", mintime)
	// 3rd trip, start -> end
	time = mintime
	mintime = maxint
	donesInit()
	explore(time, start, end, exploredirs)
	fmt.Printf("  rerun: %d\n", mintime)

	return mintime
}

//////////// Common Parts code

func parse(lines []string) () {
	gw = len(lines[0])
	gh = len(lines)
	area = gw * gh
	lcm = leastCommonMultiple(gw-2, gh-2)
	dirs = [5]Dir{
		Dir{U, -gw, true,  'U'},
		Dir{D, gw,  true,  'D'},
		Dir{L, -1,  false, 'L'},
		Dir{R, 1,   false, 'R'},
		Dir{0, 0,   false, 'O'}, // no move
	}
	grid := make([]int, area, area)
	p := 0
	for y, line := range lines {
		for x, b := range line {
			switch b {
			case '#': grid[p] = 1
			case '.':
			case '>': grid[p] = 2
			case 'v': grid[p] = 4
			case '<': grid[p] = 8
			case '^': grid[p] = 16
			default:
				panic(fmt.Sprintf("Syntax error line %d, char %d: \"%s\"\n", y+1, x+1, string(b)))
			}
			p++
		}
	}
	for p := 0; p < gw; p++ {
		if grid[p] == 0 {
			start = p
			break
		}
	}
	for p := area - gw; p < area; p++ {
		if grid[p] == 0 {
			end = p
			break
		}
	}
	gridsCreate(grid)				// generate grids for a loop of size lcm
	fmt.Printf("Grid %dx%d = %d, start=%d, end=%d, loop=%d\n", gw, gh, area, start, end, lcm)
}

func gridsCreate(grid []int) {
	grids = make([][]int, lcm, lcm)
	grids[0] = grid				// initial state
	for i := 1; i < lcm; i++ {
		grids[i] = gridNext(grids[i-1])
	}
	// debug: check that are really looping after lcm time
	grid2 := gridNext(grids[len(grids)-1])
	if !sliceEquals(grid, grid2) {
		panic(fmt.Sprintf("Grids do not loop after %d time!\n", lcm))
	}
}

// compute state after one tick
func gridNext(prev []int) []int {
	grid := make([]int, area, area)
	for p := 0; p < area; p++ {
		if prev[p] == WALL {
			grid[p] = WALL
			continue
		}
		for d := 0; d < 4; d++ {
			if (prev[p] & dirs[d].mask) != 0 {
				np := p + dirs[d].step
				if prev[np] == WALL { // wrap if hit wall
					if dirs[d].vertical {
						if np / gw == 0 { // wrap top to bottom
							np = (np % gw) + (gh-2)* gw
						} else { // wrap bottom to top
							np = (np % gw) + gw
						}
					} else {
						if np % gw == 0 { // wrap left to right
							np = gw - 2 + (np / gw) * gw
						} else  { // wrap right to left
							np = 1 + (np / gw) * gw
						}
					}
				}
				grid[np] |= dirs[d].mask
			}
		}
	}
	return grid
}

// (re-)init dones
func donesInit() {
	if len(dones) == 0 {		// init
		dones = make([][]bool, lcm, lcm)
		for i := 0; i < lcm; i++ {
			dones[i] = make([]bool, area, area)
		}
	} else {					// reset
		for i := 0; i < lcm; i++ {
			for j := 0; j < area; j++ {
				dones[i][j] = false
			}
		}
	}
}

// greatest common divisor (GCD) via Euclidean algorithm
func greatestCommonDivisor(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func leastCommonMultiple(a, b int, integers ...int) int {
	result := a * b / greatestCommonDivisor(a, b)
	for i := 0; i < len(integers); i++ {
		result = leastCommonMultiple(result, integers[i])
	}
	return result
}
	

// explore from time & pos, to goal, explore order
func explore(time, pos, goal int, expldirs []int) {
	VPf("explore(%d, %d) at [%d, %d] towards %d\n", time, pos, pos%gw, pos/gw, goal)
	dones[time % lcm][pos] = true			// avoid re-exploring from same state: (time, pos)
	time++
	if time >= mintime {		// too long, abort
		return
	}
	grid := grids[time % lcm]
	for _, d := range expldirs {
		p := pos + dirs[d].step
		if p == goal {
			if time < mintime {
				mintime = time
				VPf("explore1: found new best, time = %d / %d\n", mintime, len(grids))
			}
			continue			// a solution found, dont explore further
		}
		if p > 0 && p < area && !dones[time % lcm][p] && grid[p] == FREE {
			explore(time, p, goal, expldirs)
		}
	}
}

//////////// Part1 functions

//////////// Part2 functions
