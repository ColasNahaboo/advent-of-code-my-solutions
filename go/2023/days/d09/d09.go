// Adventofcode 2023, d09, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example
// TEST: example
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"flag"
	"fmt"
	"regexp"
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

//////////// Part 1

func part1(lines []string) (nextvalues int) {
	for _, line := range lines {
		nextvalue, _ := find(line)
		VPf("  Next value for %s ==> %d\n", line, nextvalue)
		nextvalues += nextvalue
	}
	return
}

//////////// Part 2

func part2(lines []string) (prevvalues int) {
	for _, line := range lines {
		_, prevvalue := find(line)
		VPf("  Prev value for %s ==> %d\n", line, prevvalue)
		prevvalues += prevvalue
	}
	return
}

//////////// Common Parts code

var renum = regexp.MustCompile("[-[:digit:]]+")

// find next value for a sequence. We compute its successive rows
// return next & prev values
func find(line string) (int, int) {
	numbersvalues := renum.FindAllString(line, -1)
	sequence := []int{}
	for _, nv := range numbersvalues {
		sequence = append(sequence, atoi(nv))
	}
	// compute rows until we get an all-zero one
	rows := [][]int{sequence}	// the sequence is the first row
	// r = index of the row, n = index of a number in a row
	for r := 0; r < len(sequence); r++ { // number of rows is limited by seq. len
		nextrow := make([]int, len(rows[r])-1, len(rows[r])-1)
		nonzero := 0					// jow many non-zero values in row?
		for n := 0; n < len(nextrow); n++ {
			nextrow[n] = rows[r][n+1] - rows[r][n]
			if nextrow[n] != 0 {
				nonzero++
			}
		}
		rows = append(rows, nextrow)
		if nonzero == 0 {			// all numbers in row are zero
			break
		}
	}
	// check we actually did find a non-empty row of all zeros
	if len(rows[len(rows)-1]) == 0 {
		panic("Sequence " + line + "Do not result in a all-zero line")
	}
	
	// Part1: then compute then next values, climbing back the rows
	nexts := make([]int, len(rows), len(rows))
	nexts[len(rows)-1] = 0
	for r := len(rows) - 2; r >= 0; r-- {
		nexts[r] = rows[r][len(rows[r])-1] + nexts[r+1]
	}
	// Part2: then compute then prev values, climbing back the rows
	prevs := make([]int, len(rows), len(rows))
	prevs[len(rows)-1] = 0
	for r := len(rows) - 2; r >= 0; r-- {
		prevs[r] = rows[r][0] - prevs[r+1]
	}
	VProws(rows, nexts)
	return nexts[0], prevs[0]
}
	
func VProws(rows [][]int, nexts []int) {
	if verbose {
		for r, row := range rows {
			VPf("    %v --> %d\n", row, nexts[r])
		}
	}
}


//////////// Part1 functions

//////////// Part2 functions
