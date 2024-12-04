// Adventofcode 2024, d04, in go. https://adventofcode.com/2024/day/04
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 18
// TEST: example 9
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

// We use pattern matcher functions passed as parameters.
// We do not use methods of Scalarray[int], as in Go, mixing generics and methods
// is opening a whole Pandora can of worms

package main

import (
	"flag"
	"fmt"
	// "golang.org/x/exp/slices"
)

var verbose, debug bool

var dirs8 [8]int

type PatternMatcher func(*Scalarray[int], int) int

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
	sa := parse(lines)
	return CountPattern(sa, isLinePattern)
}

func isLinePattern(sa *Scalarray[int], p int) (matches int) {
	for _, d := range dirs8 {
		if sa.a[p] == 1 &&
			sa.a[p+d] == 2 &&
			sa.a[p+2*d] == 3 &&
			sa.a[p+3*d] == 4 {
			matches++
		}
	}
	return
}

//////////// Part 2
func part2(lines []string) int {
	sa := parse(lines)
	return CountPattern(sa, isCrossPattern)
}

func isCrossPattern(sa *Scalarray[int], p int) int {
	if sa.a[p] != 3 {		// central A
		return 0
	}
	wb := sa.w
	if isCrossBranch(sa, p - wb - 1, p + wb + 1) &&
		isCrossBranch(sa, p - wb + 1, p + wb - 1) {
		return 1
	}
	return 0
}

// diagonal branch M A S or S A M?
func isCrossBranch(sa *Scalarray[int], p1, p2 int) bool {
	if (sa.a[p1] == 2 && sa.a[p2] == 4) ||
		(sa.a[p1] == 4 && sa.a[p2] == 2) {
		return true
	}
	return false
}

//////////// Common Parts code

func parse(lines []string) (*Scalarray[int]) {
	sa := makeScalarrayB[int](len(lines[0]), len(lines), 2)
	VPf("  == ARRAY %d x %d\n", sa.w - 2*sa.b, sa.h - 2*sa.b)
	var v int
	for y, line := range lines {
		for x, b := range line {
			switch b {
			case 'X' : v = 1
			case 'M' : v = 2
			case 'A' : v = 3
			case 'S' : v = 4
			default: v = 0
			}
			sa.SetB(x, y, v)
		}
	}
	wb := sa.w												 // width + border
	dirs8 = [8]int{1, wb+1, wb, wb-1, -1, -wb-1, -wb, -wb+1} // straight + diags
	return &sa
}

func CountPattern(sa *Scalarray[int], pm PatternMatcher) (total int) {
	for y := 0; y < sa.h - 2*sa.b; y++ {
		for x := 0; x < sa.w - 2*sa.b; x++ {
			p := sa.PosB(x, y)
			total += pm(sa, p)
		}
	}
	return
}

//////////// PrettyPrinting & Debugging functions

func DEBUG() {
	if ! debug { return }
	// ad hoc debug code here
}
