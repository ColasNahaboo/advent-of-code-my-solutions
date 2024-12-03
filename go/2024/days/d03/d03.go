// Adventofcode 2024, d03, in go. https://adventofcode.com/2024/day/03
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 161
// TEST: example2 48
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"flag"
	"fmt"
	"regexp"
	// "golang.org/x/exp/slices"
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
	sum, _ := parse1(lines)
	return sum
}

// general code building the muls array, but actually not needed for part2
func parse1(lines []string) (sum int, muls [][2]int) {
	remul := regexp.MustCompile("mul[(]([[:digit:]]{1,3}),([[:digit:]]{1,3})[)]")
	for _, line := range lines {
		mul := [2]int{}
		for _, m := range remul.FindAllStringSubmatch(line, -1) {
			mul[0] = atoi(m[1])
			mul[1] = atoi(m[2])
			muls = append(muls, mul)
			sum += mul[0] * mul[1]
		}
	}
	return
}

//////////// Part 2

func part2(lines []string) (sum int) {
	reinst := regexp.MustCompile("(" +
		"(" +
		"don't[(][)]" +			// #2
		")|(" +
		"do[(][)]" +			// #3
		")|(" +
		"mul[(]([[:digit:]]{1,3}),([[:digit:]]{1,3})[)]" + // #4(#5,#6)
		"))")
	doing := true				// state is kept across lines
	for _, line := range lines {
		for _, m := range reinst.FindAllStringSubmatch(line, -1) {
			VPf("  == %v %#v\n", doing, m)
			if len(m[4]) > 0 {
				if doing {
					m0 := atoi(m[5])
					m1 := atoi(m[6])
					sum += m0 * m1
				}
			} else if len(m[3]) > 0 {
				doing = true
			} else if len(m[2]) > 0 {
				doing = false
			} else {
				panic("Syntax error: \"" + m[0] + "\"")
			}
		}
	}
	return
}

//////////// Common Parts code

//////////// PrettyPrinting & Debugging functions

func DEBUG() {
	if ! debug { return }
	// ad hoc debug code here
}
