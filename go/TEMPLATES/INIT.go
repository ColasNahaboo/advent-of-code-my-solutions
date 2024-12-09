// Adventofcode YYYY, dNN, in go. https://adventofcode.com/YYYY/day/NN
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example
// TEST: example
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"flag"
	"fmt"
	"regexp"
	// "slices"
)

var verbose, debug bool

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
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
	} else {
		VP("Running Part2")
		fmt.Println(part2(lines))
	}
}

//////////// Part 1

func part1(lines []string) int {
	for _, line := range lines {
		fmt.Println(line)
	}
	return 0
}

//////////// Part 2

func part2(lines []string) int {
	for _, line := range lines {
		fmt.Println(line)
	}
	return 0
}

//////////// Common Parts code

func parse(lines []string) {
	renum := regexp.MustCompile("[[:digit:]]+")

}

//////////// PrettyPrinting & Debugging functions

func DEBUG() {
	if ! debug { return }
	// ad hoc debug code here
}
