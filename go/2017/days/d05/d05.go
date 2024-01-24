// Adventofcode 2017, d05, in go. https://adventofcode.com/2017/day/05
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 5
// TEST: example 10
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

func part1(lines []string) (i int) {
	for p, list := 0, parse(lines); p >= 0; i++ {
		VPf("  [%d] @%d %v\n", i, p, list)
		p = jump1(p, list)
	}
	return
}

func jump1(p int, list []int) (q int) {
	q = p + list[p]
	if q < 0 || q >= len(list) {
		return -1
	}
	list[p]++
	return
}

//////////// Part 2

func part2(lines []string) (i int) {
	for p, list := 0, parse(lines); p >= 0; i++ {
		VPf("  [%d] @%d %v\n", i, p, list)
		p = jump2(p, list)
	}
	return
}

func jump2(p int, list []int) (q int) {
	q = p + list[p]
	if q < 0 || q >= len(list) {
		return -1
	}
	if list[p] >= 3 {
		list[p]--
	} else {
		list[p]++
	}
	return
}

//////////// Common Parts code

func parse(lines []string) (list []int) {
	list = make([]int, len(lines), len(lines))
	for i, line := range lines {
		list[i] = atoi(line)
	}
	return
}

//////////// PrettyPrinting & Debugging functions
