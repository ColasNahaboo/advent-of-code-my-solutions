// Adventofcode 2016, d18, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 example 38
// TEST: -1 input 1987
// TEST: example 1935478
// TEST: input 19984714

// Rules trap if the 3 above are either: ^^. / .^^ / ^.. / ..^
// that is, A and C are different.

// I changed the input files to have two rows: the number of expected rows, and the first row

// rows are slices of bools (true = trap), with an element added at strat and end (walls)

package main

import (
	"flag"
	"fmt"
	// "regexp"
)

var verbose bool

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	flag.Parse()
	verbose = *verboseFlag
	infile := "input.txt"
	if flag.NArg() > 0 {
		infile = flag.Arg(0)
	}
	lines := fileToLines(infile)

	var result int
	if *partOne {
		VP("Running Part1")
		result = part1(atoi(lines[0]), lines[1])
	} else {
		VP("Running Part2")
		result = part2(lines[1])
	}
	fmt.Println(result)
}

//////////// Part 1

func part1(rows int, row0s string) (safe int) {
	row := parserow(row0s)
	if verbose {VProw(row);}	// "if" here for a tad more speed in non-verbose case 
	safe = safetiles(row)
	for i := 1; i < rows; i++ {
		row = nextrow(row)
		safe += safetiles(row)
		if verbose {VProw(row);}
	}
	return safe
}

//////////// Part 2
func part2(row0s string) (safe int) {
	return part1(400000, row0s)
}

//////////// Common Parts code

// parse the textual representation of a row into a padded []bool
func parserow(s string) []bool {
	row := make([]bool, len(s)+2, len(s)+2) // pad with false for walls on both sides
	for i := 1; i <= len(s); i++ {
		if s[i-1] == '^' {
			row[i] = true
		}
	}
	return row
}

// count the number of safetiles in a eow, excluding walls
func safetiles(r []bool) (safe int) {
	for i := 1; i < len(r)-1; i++ {
		if ! r[i] {
			safe++
		}
	}
	return
}

// compute next row
func nextrow(r []bool) []bool {
	n := make([]bool, len(r), len(r))
	for i := 1; i < len(r)-1; i++ {
		if r[i-1] != r[i+1] {	// 2 is trap if A and C are different
			n[i] = true
		}
	}
	return n
}

// debug: returns the textual representation of a row
func printrow(r []bool) string {
	s := []byte{}
	for i := 1; i < len(r)-1; i++ {
		if r[i] {	// 2 is trap if A and C are different
			s = append(s, '^')
		} else {
			s = append(s, '.')
		}
	}
	return string(s)
}

// debug: prints a row and the number of safe tiles in it
func VProw(r []bool) {
	if verbose {
		fmt.Printf("  |%s|  [%d]\n", printrow(r), safetiles(r))
	}
}

//////////// Part1 functions

//////////// Part2 functions
