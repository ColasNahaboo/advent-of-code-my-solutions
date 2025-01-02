// Adventofcode 2018, d01, in go. https://adventofcode.com/2018/day/01
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 3
// TEST: -1 example1 3
// TEST: -1 example2 0
// TEST: -1 example3 -6
// TEST: example 2
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"fmt"
	"regexp"
	// "flag"
	// "slices"
)

//////////// Options parsing & exec parts

func main() {
	ExecOptions(2, NoXtraOpts, part1, part2)
}

//////////// Part 1

func part1(lines []string) (res int) {
	changes := parse(lines)
	for _, c := range changes {
		res += c
	}	
	return 
}

//////////// Part 2

func part2(lines []string) (res int) {
	changes := parse(lines)
	seen := make(map[int]bool)
	i := 0
	for {
		res += changes[i % len(changes)]
		if seen[res] {
			return
		}
		seen[res] = true
		i++
	}
	return 
}

//////////// Common Parts code

func parse(lines []string) (changes []int) {
	renum := regexp.MustCompile("[-+[:digit:]]+")
	for _, line := range lines {
		for _, num := range renum.FindAllString(line, -1) {
			changes = append(changes, atoi(num))
		}
	}
	return
}

//////////// PrettyPrinting & Debugging functions. See also the VPx functions.

func DEBUG() {
	if ! debug { return }
	fmt.Println("DEBUG!")
	// ad hoc debug code here
}
