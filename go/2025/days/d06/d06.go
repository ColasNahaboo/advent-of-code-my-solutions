// Adventofcode 2025, d06, in go. https://adventofcode.com/2025/day/06
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 4277556
// TEST: example 3263827
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"fmt"
	"regexp"
	// "flag"
	// "slices"
)

type Problem struct {
	add bool					// add or mult operator?
	nums []int					// the values
}

//////////// Options parsing & exec parts

func main() {
	ExecOptions(2, NoXtraOpts, part1, part2)
}

//////////// Part 1

func part1(lines []string) (res int) {
	probs := parse1(lines)
	for _, prob := range probs {
		res += problemResult(prob)
	}
	return 
}

func parse1(lines []string) []Problem {
	probs := parse(lines)
	// fill the values in reading order
	renum := regexp.MustCompile("[[:digit:]]+")
	for v, line := range lines[:len(lines) - 1] { // omit the ops line
		for p, aval := range renum.FindAllString(line, -1) {
			VPf("Prob[%d].nums[%s] = %d\n", p, v, aval)
			probs[p].SetNum(v, atoi(aval))
		}
	}
	return probs
}

//////////// Part 2

func part2(lines []string) (res int) {
	probs := parse2(lines)
	for _, prob := range probs {
		res += problemResult(prob)
	}
	return 
}

func parse2(lines []string) []Problem {
	probs := parse(lines)
	renum := regexp.MustCompile("[[:digit:]]+")
	// we just rotate the input char matrix by a quarter turn to the left
	// 1 2 3           3 6 9 
	// 4 5 6    ==>    2 5 8 
	// 7 8 9           1 4 7
	// so that we can just read the numbers easily: cols => lines
	rlines := make([]string, len(lines[0]))
	for _, row := range lines[:len(lines) - 1] {
		i := len(lines[0]) - 1
		for _, digit := range row {
			rlines[i] = rlines[i] + string(digit)
			i--
		}
	}
	p := len(probs) -1			// problem number, starting with last
	n := 0						// number number
	for _, rline := range rlines {
		aval := renum.FindString(rline)
		if aval == "" {			// blank r-line separating problems
			p--					// switch to next (previous) problem
			n = 0
			continue
		}
		probs[p].SetNum(n, atoi(aval))
		n++						// next num for problem p
	}
	return probs
}

//////////// Common Parts code

// just allocate the problems, values inits are done in parse1 or parse2
// Note: the nums slices are not pre-allocated, as they are not of the same size
func parse(lines []string) []Problem {
	reop := regexp.MustCompile("[+*]")
	arity := len(lines) - 1		// number of values for each problem
	ops := reop.FindAllString(lines[arity], -1) // problems
	probs := make([]Problem, len(ops))
	for p, op := range ops {
		if op == "+" {
			probs[p].add = true
		}
	}
	return probs
}

// Sets value for number n in problem, re-allocating the slice if needed
func (prob *Problem) SetNum(n, val int) {
	if n >= len(prob.nums) {
		prob.nums = append(prob.nums, make([]int, n - len(prob.nums) + 1)...)
	}
	prob.nums[n] = val
}

func problemResult(p Problem) int {
	if p.add {
		return problemResultAdd(p)
	} else {
		return problemResultMult(p)
	}
}

func problemResultAdd(p Problem) (res int) {
	for _, v := range p.nums {
		res += v
	}
	return
}

func problemResultMult(p Problem) (res int) {
	res = 1
	for _, v := range p.nums {
		res *= v
	}
	return
}

//////////// PrettyPrinting & Debugging functions. See also the VPx functions.

func DEBUG() {
	if ! debug { return }
	fmt.Println("DEBUG!")
	// ad hoc debug code here
}
