// Adventofcode 2022, d04, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 input 560
// TEST: input 839
package main

import (
	"flag"
	"fmt"
	"regexp"
	"log"
)

var verbose bool
var reline *regexp.Regexp

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
	reline = regexp.MustCompile("^([[:digit:]]+)-([[:digit:]]+),([[:digit:]]+)-([[:digit:]]+)")

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
	inclusions := 0
	for _, line := range lines {
		vals := reline.FindStringSubmatch(line)
		if vals  == nil {
			if line != "" {
				log.Fatalf("syntax error on line: \"%s\"", line)
			}
			continue			// ignore empty lines
		}
		if include(atoi(vals[1]), atoi(vals[2]), atoi(vals[3]), atoi(vals[4])) || include(atoi(vals[3]), atoi(vals[4]), atoi(vals[1]), atoi(vals[2])) {
			inclusions += 1
		}
	}
	return inclusions
}

//////////// Part 2

func part2(lines []string) int {
	overlaps := 0
	for _, line := range lines {
		vals := reline.FindStringSubmatch(line)
		if vals  == nil {
			if line != "" {
				log.Fatalf("syntax error on line: \"%s\"", line)
			}
			continue			// ignore empty lines
		}
		if !disjoint(atoi(vals[1]), atoi(vals[2]), atoi(vals[3]), atoi(vals[4])) {
			overlaps += 1
		}
	}
	return overlaps
}

//////////// Common Parts code

//////////// Part1 functions

func include(a1, a2, b1, b2 int) bool {
	if (a1 <= b1) && (a2 >= b2) {
		return true
	} else {
		return false
	}
}

//////////// Part2 functions

func disjoint(a1, a2, b1, b2 int) bool {
	if (a2 < b1) || (b2 < a1) {
		return true
	} else {
		return false
	}
}
