// Adventofcode 2017, d24, in go. https://adventofcode.com/2017/day/24
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 31
// TEST: example 19
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"flag"
	"fmt"
	"regexp"
	// "golang.org/x/exp/slices"
)

var verbose, debug bool

type Comp [2]int
type Comps []Comp

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
	comps := parse(lines)
	return MaxComps(0, 0, comps)
}

// we recurse until we canmnot fit any component. return max strength
func MaxComps(start, strength int, availableComps Comps) int {
	VPf("  [%d]%d %v\n", strength, start, availableComps)
	comps := availableComps.Fits(start)
	maxstrength := strength
	for _, c := range comps {
		startN := c.Other(start)
		compsN, _ := deleteElt[Comp](availableComps, c)
		strengthN := MaxComps(startN, strength + c[0] + c[1], compsN)
		if strengthN > maxstrength {
			maxstrength = strengthN
		}
	}
	return maxstrength
}

func (c Comp) Fits(n int) bool {
	return c[0] == n || c[1] == n
}

func (comps Comps) Fits(n int) (fits Comps) {
	for _, c := range comps {
		if c.Fits(n) {
			fits = append(fits, c)
		}
	}
	return
}

func (c Comp) Other(n int) int {
	if c[0] == n {
		return c[1]
	} else {
		return c[0]
	}
}

func (c Comp) Strength() int {
	return c[0] + c[1]
}

func (b Comps) Strength() (strength int) {
	for _, c := range b {
		strength +=  c.Strength()
	}
	return
}

//////////// Part 2

func part2(lines []string) int {
	comps := parse(lines)
	maxlength, maxstrength := MaxLenComps(0, 0, 0, comps)
	VPf("Longest bridge length: %d\n", maxlength)
	return maxstrength
}

// we recurse until we cannot fit any component.
// return length of longests bridges and their max strength
func MaxLenComps(start, length, strength int, availableComps Comps) (maxlen, maxstr int) {
	comps := availableComps.Fits(start)
	maxlen = length
	maxstr = strength
	for _, c := range comps {
		startN := c.Other(start)
		strN := c.Strength()
		compsN, _ := deleteElt[Comp](availableComps, c)
		ml, ms := MaxLenComps(startN, length + 1, strength + strN, compsN)
		if ml > maxlen {
			maxlen = ml			// new longest sub-bridge
			maxstr = ms			// reset maxstr for this new length
		} else if ml == maxlen && ms > maxstr {
			maxstr = ms			// same max length, look for best strength
		}
	}
	VPf("[%d]%d: %d, %d\n", start, strength, maxlen, maxstr)
	return 
}

//////////// Common Parts code

func parse(lines []string) (comps []Comp) {
	renum := regexp.MustCompile("[[:digit:]]+")
	for _, line := range lines {
		nums := renum.FindAllString(line, -1)
		comps = append(comps, Comp{atoi(nums[0]), atoi(nums[1])})
	}
	return
}

//////////// PrettyPrinting & Debugging functions

func DEBUG() {
	if ! debug { return }
	// ad hoc debug code here
}
