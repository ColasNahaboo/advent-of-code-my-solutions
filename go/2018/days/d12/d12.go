// Adventofcode 2018, d12, in go. https://adventofcode.com/2018/day/12
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 325
// TEST: example 999999999374
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

// for the part2, by displaying the generations (with -1 -v -g G -w W) with
// W wide enough (e.g. 500) and G=150
// wee see that after some generation (147 in my input), the plants form a
// "runner", a pattern moving up one pot per generation
// We thus just compute the generation 200 (a safety margin), find the runner,
// and compute the score supposing that it has moved billions up


package main

import (
	"fmt"
	"regexp"
	"flag"
	"slices"
)

var gens = 20					// number of generations
var dwidth = 80

//////////// Options parsing & exec parts

func main() {
	flag.IntVar(&gens, "g", 20, "Number of generations (default 20)")
	flag.IntVar(&dwidth, "w", 80, "for -v, width of the display (def. 80)")
	ExecOptions(2, NoXtraOpts, part1, part2)
}

//////////// Part 1

func part1(lines []string) (res int) {
	off := 2*gens + 2			// manage pots [-off, init + off[ + 2 at end
	init, rules := parse(lines)
	slen :=  len(init) + 2 * off
	s := make([]bool, slen, slen)
	ns := slices.Clone(s)
	n0 := slices.Clone(s)
	copy(s[off:], init)			// initial state starts at off
	VPgen(0, s)
	for i := 0; i < gens; i++ {
		copy(ns, n0)			// fast-clears ns
		Generate(&ns, s, off, gens-i, rules) // create next s generation into ns
		copy(s, ns)
		VPgen(i, s)
	}
	return Score(s, off)
}

func Generate(ns *[]bool,s []bool, off, gens int, rules []bool) {
	for i := gens * 2; i < len(s) - gens * 2; i++ {
		ix := IX5(s, i)
		if  rules[ix] {
			(*ns)[i] = true
		}
	}
}

func IX5(s []bool, p int) (ix int) {
	for i := p+2; i >= p-2; i-- {
		ix *= 2
		if s[i] {
			ix++
		}
	}
	return
}

func Score(s []bool, off int) (res int) {
	for i, p := range s{
		if p {
			res += i - off		// add number of the pot: i-off
		}
	}
	return
}

//////////// Part 2

func part2(lines []string) (res int) {
	gens = 200					// sufficient to reach the stable state
	off := 2*gens + 2			// manage pots [-off, init + off[ + 2 at end
	init, rules := parse(lines)
	slen :=  len(init) + 2 * off
	s := make([]bool, slen, slen)
	ns := slices.Clone(s)
	n0 := slices.Clone(s)
	copy(s[off:], init)			// initial state starts at off
	VPgen(0, s)
	for i := 0; i < gens; i++ {
		copy(ns, n0)			// fast-clears ns
		Generate(&ns, s, off, gens-i, rules) // create next s generation into ns
		copy(s, ns)
		VPgen(i, s)
	}
	x, y := FindRunner(s)
	return ScoreRunner(s[x:y], x + 50000000000 - gens, off)
}

// score of a runner when it begins at position p

func ScoreRunner(s []bool, from, off int) (res int) {
	for i, p := range s {
		if p {
			res += from + i - off		// add number of the pot: i-off
		}
	}
	return
}


// return the enclosing interval of all plants as [x y[
func FindRunner(s []bool) (x, y int) {
	for i := 0; i < len(s); i++ {
		if s[i] {
			x = i-1
			break
		}
	}
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] {
			y = i+1
			break
		}
	}
	return
}


//////////// Common Parts code

func parse(lines []string) (state, rules []bool) {
	reinit := regexp.MustCompile(": *([.#]+)")
	m := reinit.FindStringSubmatch(lines[0])
	state = make([]bool, len(m[1]))
	for i, r := range m[1] {
		if r == '#' {
			state[i] = true
		}
	}
	// keep only rules creating a plant for simplicity
	rules = make([]bool, 32, 32) // all the 5-bit numbers possible
	rerule := regexp.MustCompile("^([.#])([.#])([.#])([.#])([.#]) => #")
	for _, line := range lines {
		if m = rerule.FindStringSubmatch(line); m != nil {
			ix := 0
			for i := 5; i > 0; i-- {
				ix *= 2
				ix += DotHash2Bit(m[i])
			}
			rules[ix] = true
		}
	}
	return
}

func DotHash2Bit(s string) int {
	if s == "#" {
		return 1
	}
	return 0
}

//////////// PrettyPrinting & Debugging functions. See also the VPx functions.

func DEBUG() {
	if ! debug { return }
	fmt.Println("DEBUG!")
	// ad hoc debug code here
}

func VPgen(i int, s []bool) {
	if ! verbose { return }
	margin := (len(s) - (dwidth - 5)) / 2
	fmt.Printf("%2d: ", i)
	for i := margin; i < len(s) - margin; i++ {
		if s[i] {
			fmt.Print("#")
		} else {
			fmt.Print(".")
		}
	}
	fmt.Println()
}
