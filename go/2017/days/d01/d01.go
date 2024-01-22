// Adventofcode 2017, d01, in go. https://adventofcode.com/2017/day/01
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example1 3
// TEST: -1 example2 4
// TEST: -1 example3 0
// TEST: -1 example4 9
// TEST: example21 6
// TEST: example22 0
// TEST: example23 4
// TEST: example24 12
// TEST: example25 4
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"flag"
	"fmt"
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

func part1(lines []string) (sum int) {
	line := lines[0]
	for i, c := range line {
		if byte(c) == line[(i+1) % len(line)] {
			sum += int(c - '0')
		}
	}
	return
}

//////////// Part 2
func part2(lines []string) (sum int) {
	line := lines[0]
	offset := len(line) / 2
	for i, c := range line {
		if byte(c) == line[(i + offset) % len(line)] {
			sum += int(c - '0')
		}
	}
	return
}

//////////// Common Parts code

//////////// PrettyPrinting & Debugging functions
