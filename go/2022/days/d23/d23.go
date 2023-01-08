// Adventofcode 2022, d23, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 example 110
// TEST: -1 input 4034
// TEST: example 20
// TEST: input 960
package main

import (
	"flag"
	"fmt"
	// "regexp"
)

var verbose bool
// coordinates x,y start at top left. A border of size rounds is added around it for growth
var elves []int			    // current positions of elves (scalar): position = x + y * gw
var gw, gh int				// grid width & height
var pgrid []int			   // 2d map of props. Value is id of elf proposing + 1
var	egrid []bool		   // 2D map of elves
var delta [4]int			// delta pos for each dir: 0=-gw, 1=gw, 2=-1, E=1
var rim [8]int				// the 8 positions around a space
var scans [4][3]int			// for each dir the 3 spaces to scan: for 0: N, NE, NW
const maxint = 8888888888888888888 // easily identifiable in debug
var maxrounds int				   // max expected rounds for part2

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	roundsFlag := flag.Int("r", 1000, "max expected rounds for part2")
	flag.Parse()
	verbose = *verboseFlag
	maxrounds = *roundsFlag
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
	parse(lines, 10)
	VPgrid(-1)
	for i := 0; i < 10; i++ {
		elves, _ = round(i)
		VPgrid(i)
	}
	return score()
}

//////////// Part 2
func part2(lines []string) int {
	parse(lines, maxrounds)
	var moved int
	VPgrid(-1)
	for i := 0; i < maxrounds; i++ {
		elves, moved = round(i)
		VPgrid(i)
		if moved == 0 {
			return i+1
		}
	}
	panic(fmt.Sprintf("Not enough rounds allocated: %d. Please use more with -r\n", maxrounds))
}

//////////// Common Parts code

func parse(lines []string, rounds int) {
	gw = len(lines[0]) + 2 * rounds + 2 // anticipate growing once per round + margin
	gh = len(lines) + 2 * rounds + 2
	delta=[4]int{-gw, gw, -1, 1}
	scans=[4][3]int{
		[3]int{-gw, -gw+1, -gw-1},
		[3]int{gw, gw+1, gw-1},
		[3]int{-1, -gw-1, gw-1},
		[3]int{1, -gw+1, gw+1},
	}
	rim = [8]int{-gw, -gw+1, 1, gw+1, gw, gw-1, -1, -gw-1}
	for y, line := range lines {
		for x, b := range line {
			if b == '#' {
				p := x+rounds+1 + (y+rounds+1)*gw
				elves = append(elves, p)
			} else if b != '.' {
				panic(fmt.Sprintf("Syntax error line %d, char: %s\n", y+1, string(b)))
			}
		}
	}
	pgrid = make([]int, gw*gh, gw*gh)			// 2D map of props. Value is id of elf proposing + 1
	egrid = make([]bool, gw*gh, gw*gh)			// 2D map of elves
	for _, p := range elves {
		egrid[p] = true
	}
	VPf("%d elves on a grid %d x %d, room for %d rounds\n", len(elves), gw, gh, rounds)
}

// firstdir = first dir to consider for proposal: 0=N, 1=S, 2=W, 3=E
func round(firstdir int) (props []int, moved int) {
	props = make([]int, len(elves), len(elves)) // proposed position for elves
	pgridset := make([]int, 0)
NEXTELF:
	for e, p := range elves {
		for i := 0; i < 8; i++ {
			if egrid[p + rim[i]] {
				goto MUSTMOVE
			}
		}
		props[e] = p			// no adjacent elves, don't move
		continue NEXTELF
	MUSTMOVE:
		for d := 0; d < 4; d++ { // probe 4 directions
			dir := (firstdir + d) % 4
			if !egrid[p + scans[dir][0]] && !egrid[p + scans[dir][1]] && !egrid[p + scans[dir][2]] {
				if other := pgrid[p + delta[dir]] - 1; other != -1 { // collision of proposals
					props[other] = elves[other]						 // abort proposal of other elf
					props[e] = p									 // abort ours
				} else {
					pgrid[p + delta[dir]] = e + 1 // mark our proposal on grid, avoid 0
					pgridset = append(pgridset, p + delta[dir])
					props[e] = p + delta[dir]
					moved++
				}
				continue NEXTELF
			}
		}
		props[e] = p			// no free dir to move to, stay here
	}
	// clean grids for next time
	for _, p := range pgridset { // zero pgrid
		pgrid[p] = 0
	}
	for e := range props {		// clear old positions in grid of moved elves
		if props[e] != elves[e] {
			egrid[elves[e]] = false
			egrid[props[e]] = true
		}
	}
	return
}

func score() int {
	xmin, ymin, xmax, ymax := minrectangle()
	VPf("Score: %d elves on a rectangle [%d,%d]x[%d,%d] = %d\n", len(elves), xmin, xmax, ymin, ymax, (xmax - xmin + 1) * (ymax - ymin + 1))
	return (xmax - xmin + 1) * (ymax - ymin + 1) - len(elves)
}

func minrectangle() (xmin, ymin, xmax, ymax int) {
	xmin, ymin = maxint, maxint
	xmax, ymax = 0, 0
	for _, p := range elves {
		x := p % gw
		y := p / gw
		if x < xmin { xmin = x;}
		if x > xmax { xmax = x;}
		if y < ymin { ymin = y;}
		if y > ymax { ymax = y;}
	}
	return
}

// debug
func VPgrid(i int) {
	if !verbose {return;}
	xmin, ymin, xmax, ymax := minrectangle()
	if i < 0 {
		fmt.Printf("== Initial State ==\n")
	} else {
		fmt.Printf("== End of Round %d ==\n", i+1)
	}
	egrid := make([]bool, gw*gh, gw*gh)			// 2D map of elves
	for _, p := range elves {
		egrid[p] = true
	}
	for y := ymin; y <= ymax; y++ {
		for x := xmin; x <= xmax; x++ {
			if egrid[x + y*gw] {
				fmt.Printf("#")
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Printf("\n")
	}
}
	
//////////// Part1 functions

//////////// Part2 functions
