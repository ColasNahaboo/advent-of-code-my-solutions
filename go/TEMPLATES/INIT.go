// Adventofcode YYYY, dNN, in go. https://adventofcode.com/YYYY/day/NN
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example
// TEST: example
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"fmt"
	"regexp"
	// "flag"
	// "slices"
)

//////////// Options

func main() {
	ParseOptions(2, part1, part2)
}
func ProcessXtraOptions() {} //  extra options, see ParseOptions in utils.go

//////////// Part 1

func part1(lines []string) (res int) {
	parse(lines)
	return 
}

//////////// Part 2

func part2(lines []string) (res int) {
	parse(lines)
	return 
}

//////////// Common Parts code

func parse(lines []string) []string {
	renum := regexp.MustCompile("[[:digit:]]+") // example code body, replace.
	for _, line := range lines {
		m := renum.FindAllString(line, -1)
		return m
	}
	return []string{}
}

//////////// PrettyPrinting & Debugging functions. See also the VPx functions.

func DEBUG() {
	if ! debug { return }
	fmt.Println("DEBUG!")
	// ad hoc debug code here
}
