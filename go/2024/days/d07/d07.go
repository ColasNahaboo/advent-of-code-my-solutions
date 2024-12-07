// Adventofcode 2024, d07, in go. https://adventofcode.com/2024/day/07
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 3749
// TEST: example 11387
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

// From the input, we can add the following assertions:
// - all numbers are positive
// - the results of operators are always bigger than any of the operands
// So, we can stop computing as soon as we become bigger than testval

package main

import (
	"flag"
	"fmt"
	"regexp"
	// "golang.org/x/exp/slices"
)

var verbose, debug bool

type Equation struct {
	testval int					// result of the equation to get
	nums []int					// the numbers to combine
}

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

func part1(lines []string) (sum int) {
	equations := parse(lines)
	for _, e := range equations {
		if EqIsTrue(&e, e.nums[0], 1) {
			sum += e.testval
		}
	}
	return
}

// recurse to compute value remaining from nums index i and test it
func EqIsTrue(e *Equation, val, i int) bool {
	if i >= len(e.nums) || val >= e.testval { // all nums done, or early overflow
		return val == e.testval
	}
	return EqIsTrue(e, val + e.nums[i], i+1) || // operator +
		EqIsTrue(e, val * e.nums[i], i+1)		// operator *
}

//////////// Part 2

func part2(lines []string) (sum int) {
	equations := parse(lines)
	for _, e := range equations {
		if EqIsTrue2(&e, e.nums[0], 1) {
			sum += e.testval
		}
	}
	return
}

// recurse to compute value remaining from nums index i and test it
func EqIsTrue2(e *Equation, val, i int) bool {
	if val > e.testval { 		// overflow => fail, no need to continue
		return false
	}
	if i >= len(e.nums) {		// computation complete, do we pass the test?
		return val == e.testval
	}
	return EqIsTrue2(e, val + e.nums[i], i+1) ||	// operator +
		EqIsTrue2(e, val * e.nums[i], i+1) ||		// operator *
		EqIsTrue2(e, OpConcat(val, e.nums[i]), i+1)	// operator ||
}

func OpConcat(i, j int) int {
	return atoi(itoa(i) + itoa(j))
}

//////////// Common Parts code

func parse(lines []string) (equations []Equation) {
	renum := regexp.MustCompile("[[:digit:]]+")
	for _, line := range lines {
		m := renum.FindAllString(line, -1)
		nums := []int{}
		for i := 1; i < len(m); i++ {
			nums = append(nums, atoi(m[i]))
		}
		equations = append(equations, Equation{atoi(m[0]), nums})
	}
	return
}

//////////// PrettyPrinting & Debugging functions

func DEBUG() {
	if ! debug { return }
	// ad hoc debug code here
}
