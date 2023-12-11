// Adventofcode 2023, d11, in go. https://adventofcode.com/2023/day/11
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 374
// TEST: example 82000210
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"flag"
	"fmt"
	//"regexp"
)

var verbose bool

// For convenience, we use global variables for the universe
// a galaxy is defined by its ID, its index in the univ array
// a galaxy position is a number: x + y * uw
var univ []int					// the universe, list of galaxies positions
// the univ dimensions: width and height of its grid. thus usize = uw*uh
var uw, uh int
// the list of galaxies in each row and col, as their ID + 1 (0 meaning no gal) 
var rows, cols [][]int

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

//////////// Part 1

func part1(lines []string) int {
	return computeDistances(lines, 1)
}

//////////// Part 2

func part2(lines []string) int {
	return computeDistances(lines, 1000000 - 1)
}

//////////// Common Parts code

func computeDistances(lines []string, expansion int) (sum int) {
	parse(lines)
	VPuniv("Initial")
	expandUniv(expansion)
	VPuniv(fmt.Sprintf("Expansion by %d", expansion))
	// unordered pairs: we take all (g1, g2) with g2 < g1
	for g1 := 0; g1 < len(univ) - 1; g1++ {
		for g2 := g1+1; g2 < len(univ); g2++ {
			sum += distance(g1, g2)
		}
	}
	return

}

func parse(lines []string) {
	uw = len(lines[0])
	uh = len(lines)
	for y := 0; y < len(lines); y++ {
		for x := 0; x < len(lines[0]); x++ {
			if lines[y][x] == '#' {
				galaxy := x + y*uw
				univ = append(univ, galaxy)
			}
		}
	}
}

func expandUniv(expansion int) {
	// compute rows and cols
	rows = make([][]int, uh)
	for y, _ := range rows {
		rows[y] = []int{}
	}
	cols = make([][]int, uw)
	for x, _ := range cols {
		cols[x] = []int{}
	}
	for gid, gal := range univ {
		rows[gal/uw] = append(rows[gal/uw], gid)
		cols[gal%uw] = append(cols[gal%uw], gid)
	}
	VPf("Rows:\n %v\n", rows)
	VPf("Cols:\n %v\n", cols)
	VPunivcoords()
	// delay expand for after we computed the new uw
	// lists of amounts to move down and right for all gids
	tmd := make([]int, len(univ), len(univ))
	tmr := make([]int, len(univ), len(univ))
	ouw := uw
	for y, exp := 0, 0; y < len(rows); y++ {
		if len(rows[y]) == 0 {
			exp += expansion				// increase the expansion
			uh += expansion
			continue
		}
		if exp == 0 {
			continue
		}
		for _, gid := range rows[y] { // shift later, when new uw is known
			tmd[gid] = exp
		}
	}
	// expand cols rightwards
	for x, exp := 0, 0; x < len(cols); x++ {
		if len(cols[x]) == 0 {
			exp += expansion				// increase the expansion
			uw += expansion
			continue
		}
		if exp == 0 {
			continue
		}
		for _, gid := range cols[x] { // shift right later
			tmr[gid] = exp
		}
	}
	VPf("To move Down: %v\n", tmd)
	VPf("To move Right: %v\n", tmr)
	
	// perform the moves
	for gid, pos := range univ {
		ox := pos%ouw
		oy := pos/ouw
		x := ox + tmr[gid]
		y := oy + tmd[gid]
		VPf("  [%d] (%d,%d)%d --> (%d,%d)%d\n", gid, ox, oy, pos, x, y, x+y*uw)
		univ[gid] = x + y*uw
	}
	VPunivcoords()
}

// mahattan distance
func distance(g1, g2 int) int {
	return intAbs(univ[g2]%uw - univ[g1]%uw) + intAbs(univ[g2]/uw - univ[g1]/uw)
}

//////////// PrettyPrinting & Debugging functions

func VPuniv(label string) {
	if !verbose {
		return
	}
	fmt.Printf("%s: Universe %d x %d, with %d galaxies:\n", label, uw, uh, len(univ))
	// grid is a ready-to-print byte string, we add newlines at the end of rows
	gw := uw+1
	grid := make([]byte, gw*uh, gw*uh)
	for i := range grid {
		grid[i] = '.'
	}
	for gid, gal := range univ {
		gpos := gal%uw + gal/uw * gw
		grid[gpos] = '0' + byte(gid % 10)
	}
	for y := 0; y < uh; y++ {
		grid[(y+1)*gw - 1] = '\n'
	}
	fmt.Print(string(grid))
}

func VPunivcoords() {
	if !verbose {
		return
	}
	for gid, gal := range univ {
		fmt.Printf(" [%d]=%d(%d,%d)", gid, gal, gal%uw, gal/uw)
	}
	fmt.Println()
}
