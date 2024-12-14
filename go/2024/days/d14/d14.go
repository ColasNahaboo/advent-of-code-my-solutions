// Adventofcode 2024, d14, in go. https://adventofcode.com/2024/day/14
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 12
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

// NOTE: I have added a line w=width, h=height at the start of input files
// to specify the board size. E.g. for example.txt: w=11, h=7

// For part 2, we just look for a consecutive horizontal line of at least
// 24 guards, as the xmas tree is drawn inside a box

package main

import (
	"flag"
	"fmt"
	"regexp"
	// "slices"
)

var verbose, debug bool

type Guard struct {
	p, v Point
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

func part1(lines []string) int {
	seconds := 100
	w, h, guards := parse(lines)
	quadrants := [5]int{}		// quadrants[4] is used as trashcan
	for _, g := range guards {
		x100 := AfterTime(g.p.x, g.v.x, w, seconds)
		y100 := AfterTime(g.p.y, g.v.y, h, seconds)
		qi := QuadrantIndex(x100, y100, w, h)
		VPf("  Guard %v ==> %d, %d = Q%d\n", g, x100, y100, qi)
		quadrants[qi]++
	}
	VP(quadrants)
	return quadrants[0] * quadrants[1] * quadrants[2] * quadrants[3]
}

func AfterTime(x, v, size, t int) int {
	nx := (x + v * t) % size
	if nx < 0 {
		return nx + size
	}
	return nx
}

// 0 1
// 2 3   + 4: trashcan

func QuadrantIndex(x, y, w, h int) int {
	if x < (w/2) {
		if y < (h/2) {
			return 0
		} else if y >= (h/2+1) {
			return 2
		}
	} else if x >= (w/2+1) {
		if y < (h/2) {
			return 1
		} else if y >= (h/2+1) {
			return 3
		}
	}
	return 4
}

//////////// Part 2

func part2(lines []string) (res int) {
	seconds := 10000
	minlinelength := 24
	w, h, guards := parse(lines)
	s := MakeSheet(w, h)		// a blank sheet to write on, rows concatenated
	s0 := MakeSheet(w, h)		// blank version to reset it fast by copy
	for t := range seconds {
		copy(s, s0)
		for _, g := range guards { // draw guards
			x := AfterTime(g.p.x, g.v.x, w, t)
			y := AfterTime(g.p.y, g.v.y, h, t)
			s[x + w * y] = '#'
		}
		if HasConsecutive(s, minlinelength) {
			if verbose {
				PrintSheet(s, w, h, t)
			}
			return t
		}
	}
	return 0
}

func MakeSheet(w, h int) (s []byte) {
	s = make([]byte, w*h, w*h)
	for i := range s {
		s[i] = '.'
	}
	return
}

func HasConsecutive(s []byte, l int) bool {
NEXTGUARD:
	for i := 0; i < len(s) - l; i++ {
		if s[i] == '#' {
			for range l {
				i++
				if s[i] != '#' {
					continue NEXTGUARD
				}
			}
			return true
		}
	}
	return false
}

func PrintSheet(s []byte, w, h, t int) {
	for y := range h {
		fmt.Print("[")
		fmt.Print(string(s[y*w:(y+1)*w]))
		fmt.Println("]")
	}
	fmt.Printf("After second: %d\n", t)
}
	

//////////// Common Parts code

func parse(lines []string) (w, h int, guards []Guard) {
	renum := regexp.MustCompile("-?[[:digit:]]+")
	dims := atoil(renum.FindAllString(lines[0], -1))
	w, h = dims[0], dims[1]
	for _, line := range lines[1:] {
		m := atoil(renum.FindAllString(line, -1))
		g := Guard{Point{m[0], m[1]}, Point{m[2], m[3]}}
		guards = append(guards, g)
	}
	return
}

//////////// PrettyPrinting & Debugging functions

func DEBUG() {
	if ! debug { return }
	// ad hoc debug code here
}

func VPguards(title string, b Board[int], guards []Guard, seconds int) {
	if ! verbose {
		return
	}
	fmt.Println(title)
	b.Fill(0)
	for _, g := range guards {
		x := AfterTime(g.p.x, g.v.x, b.w, seconds)
		y := AfterTime(g.p.y, g.v.y, b.h, seconds)
		b.a[x][y]++
	}
	for y := range b.h {
		for x := range b.w {
			if b.a[x][y] != 0 {
				fmt.Print(b.a[x][y])
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}
