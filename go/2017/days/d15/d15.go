// Adventofcode 2017, d15, in go. https://adventofcode.com/2017/day/15
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 588
// TEST: example 309
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

func part1(lines []string) (n int) {
	gav, gbv := parse(lines)
	for i := 0; i < 40000000; i++ {
		gav = gen(gav, factorA)
		gbv = gen(gbv, factorB)
		if gav % sixteenbits ==  gbv % sixteenbits {
			VPf("  Match at %d\n", i)
			n++
		}
	}
	return
}

//////////// Part 2

func part2(lines []string) (n int) {
	gav, gbv := parse(lines)
	for i := 0; i < 5000000; i++ {
		gav = gen2(gav, factorA, mulA)
		gbv = gen2(gbv, factorB, mulB)
		VPf("  [%d] %12d %12d\n", i, gav, gbv)
		if gav % sixteenbits ==  gbv % sixteenbits {
			VPf("  Match at %d\n", i)
			n++
		}
	}
	return
}

//////////// Common Parts code

const divider = 2147483647
const factorA = 16807
const factorB = 48271
const mulA = 4
const mulB = 8
const sixteenbits = 65536

func parse(lines []string) (gas, gbs int) {
	re := regexp.MustCompile("[[:digit:]]+")
	gas = atoi(re.FindString(lines[0]))
	gbs = atoi(re.FindString(lines[1]))
	return
}

func gen(value, factor int) int {
	return (value * factor ) % divider
}

func gen2(value, factor, mul int) int {
	for {
		value = (value * factor ) % divider
		if value % mul == 0 {
			return value
		}
	}
}

//////////// PrettyPrinting & Debugging functions
