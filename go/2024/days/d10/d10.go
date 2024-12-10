// Adventofcode 2024, d10, in go. https://adventofcode.com/2024/day/10
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example1 1
// TEST: -1 example2 36
// TEST: example1 16
// TEST: example2 81
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

// we explore trails starting from 0s, when we reach a top (z value == 9) we
// increment the integer value for this top in a Map
// For part 1, we just count the number of tops in the map (its size)
// Fpr part 2, we add all the values in the map, this counts all the trails

package main

import (
	"flag"
	"fmt"
	// "regexp"
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
		fmt.Println(part2(lines))
	}
}

//////////// Part 1

func part1(lines []string) (score int) {
	b, trailheads := parse(lines)
	for _, th := range trailheads {
		tops := map[Point]int{}
		VPf("== Looking for trails starting at %v\n", th)
		TrailsFrom(b, th, tops)
		score += len(tops)
	}
	return
}

func TrailsFrom(b *Board[int], s Point, tops map[Point]int) {
	for _, step := range StepsFrom(b, s) {
		if b.a[step.x][step.y] == 9 {
			VPf("   Trail found ending at %v\n", step)
			tops[step]++
		} else {
			TrailsFrom(b, step, tops)
		}
	}
	return
}

// list of points ortho-adjacent to s with elevation (z) +1 
func StepsFrom(b *Board[int], s Point) (steps []Point) {
	for _, d := range DirsOrtho {
		step := s.Add(d)
		if b.Inside(step) && b.a[step.x][step.y] == b.a[s.x][s.y] + 1 {
			steps = append(steps, step)
		}
	}
	VPf("     Steps possibles from %v: %v\n", s, steps)
	return
}

//////////// Part 2

func part2(lines []string) (score int) {
	b, trailheads := parse(lines)
	for _, th := range trailheads {
		tops := map[Point]int{}
		VPf("== Looking for trails starting at %v\n", th)
		TrailsFrom(b, th, tops)
		for _, rating := range tops {
			score += rating
		}
	}
	return
}

//////////// Common Parts code

func parse(lines []string) (bp *Board[int], trailheads []Point) {
	b := parseBoard[int](lines, func(x, y int, r rune) int {
		if r == '0' {
			trailheads = append(trailheads, Point{x, y})
		}
		if r == '.' {			// for some examples, '.' is not a th but its z=0
			return 0
		} else {
			return int(r - '0')
		}
	})
	b.VPBoard(fmt.Sprint(trailheads))
	return b, trailheads
}

//////////// PrettyPrinting & Debugging functions

func DEBUG() {
	if ! debug { return }
	// ad hoc debug code here
}
