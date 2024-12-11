// Adventofcode 2024, d11, in go. https://adventofcode.com/2024/day/11
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example2 55312
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

// for part2, we use a cache, since the transformation on a stone depends only
// on its own value, not its position nor neighbours, we can say that if we have
// n stones marked S, and S becomes S1 and S2, we have n*S1 and n*S2
// Its cuts greatly on the number of computations

// part2 uses bignums. It was not actually necessary, but I got some practice.

package main

import (
	"flag"
	"fmt"
	"regexp"
	"math/big"
	// "slices"
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
		part2(lines)
	}
}

//////////// Part 1

func part1(lines []string) int {
	return len(Blink(parse(lines), 25))
}

// Blink N times
func Blink (s []int, n int) []int {
	var ns []int
	for _ = range n {
		for _, s := range s {
			if s == 0 {				// engraved with the number 0 => 1
				ns = append(ns, 1)
			} else if ss := itoa(s); len(ss) % 2 == 0 { // even digits ==> split
				ns = append(ns, atoi(ss[:len(ss)/2]))
				ns = append(ns, atoi(ss[len(ss)/2:]))
			} else {
				ns = append(ns, s * 2024)
			}
		}
		s, ns = ns, s			// we alternate and reuse 
		ns = ns[:0]
	}
	return s
}

//////////// Part 2

func part2(lines []string) {
	stones := map[int]*big.Int{} // map[stone-marked-value] => number-of-them
	one := big.NewInt(1)
	for _, s := range parse(lines) {
		IncStones(&stones, s, one)
	}
	for _ = range 75 {
		stones2 := map[int]*big.Int{} // new one
		for s, n := range stones {
			s1, s2  := BlinkOne(s)
			IncStones(&stones2, s1, n)
			if s2 != -1 {
				IncStones(&stones2, s2, n)
			}
		}
		stones, stones2 = stones2, stones
	}
	res := big.NewInt(0)
	for _, n := range stones {
		res.Add(res, n)
	}
	fmt.Println(res.String())
}

// in map stones, increment stones[s] by n
func IncStones (stones *map[int]*big.Int, s int, n *big.Int) {
	if v, ok := (*stones)[s]; ok {
		nbi := &big.Int{}
		(*stones)[s] = nbi.Add(v, n)
	} else {
		(*stones)[s] = n
	}
}

// blink on one stone, return the 2 new stones (or the new stone and -1)
func BlinkOne(s int) (int, int) {
	if s == 0 {				// engraved with the number 0 => 1
		return 1, -1
	} else if ss := itoa(s); len(ss) % 2 == 0 { // even digits ==> split
		return atoi(ss[:len(ss)/2]), atoi(ss[len(ss)/2:])
	} else {
		return s * 2024, -1
	}
}

//////////// Common Parts code

func parse(lines []string) (res []int) {
	renum := regexp.MustCompile("[[:digit:]]+")
	for _, s := range renum.FindAllString(lines[0], -1) {
		res = append(res, atoi(s))
	}
	return 
}

//////////// PrettyPrinting & Debugging functions

func DEBUG() {
	if ! debug { return }
	// ad hoc debug code here
}
