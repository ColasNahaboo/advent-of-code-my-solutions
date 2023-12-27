// Adventofcode 2023, d13, in go. https://adventofcode.com/2023/day/13
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 405
// TEST: example 400
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

// input are grids of odd sides between 7 and 17

// I implemented part1 in a smart wat, pre-computing rows and columns
// But when I saw part2, I just did it in a straightforward way

package main

import (
	"flag"
	"fmt"
	//"regexp"
)

var verbose bool

type Pattern struct {
	w, h int 					// dims
	rows []string				// the grid, row by row
	cols []string				// the grid, col by col flipped on the NW-SE axis
}
var patterns []Pattern			// the list of all patterns, the input

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	flag.Parse()
	verbose = *verboseFlag
	var infile string
	if flag.NArg() > 0 {
		infile = flag.Arg(0)
	} else {
		infile = fileMatch("input,[0-9]*,[0-9]*.test")
	}
	lines := fileToLines(infile)

	var result int
	if *partOne {
		VP("Running Part1")
		result = part1(lines)
	} else {
		VP("Running Part2")
		result = part2(lines)
	}
	fmt.Println(result)
}

//////////// Part 1

func part1(lines []string) (note int) {
	parse(lines)
	for i, p := range patterns {
		VPf("  [%d] %d x %d: ", i, p.w, p.h)
		vnote := vNote(p)
		VPf("/ ")
		hnote := hNote(p)
		note += vnote + 100 * hnote
		VPf("= %d + 100 * %d = %d\n", vnote, hnote, vnote + 100 * hnote)
	}
	return
}

func parse(lines []string) {
	patterns = []Pattern{}
	start := 0
	for l, line := range lines {
		if line == "" {
			patterns = append(patterns, parsePattern(lines[start:l]))
			start = l+1
		}
	}
}

func parsePattern(lines []string) (p Pattern) {
	p.w = len(lines[0])
	p.h = len(lines)
	p.rows = make([]string, p.h, p.h)
	for i, line := range lines {
		p.rows[i] = line
	}
	p.cols = make([]string, p.w, p.w)
	for x := 0; x < p.w; x++ {
		col := make([]byte, p.h, p.h)
		for y := 0; y < p.h; y++ {
			col[y] = lines[y][x]
		}
		p.cols[x] = string(col)
	}
	return
}

// note for symmetries around vertical lines
func vNote(p Pattern) (note int) {
	for i := 0; i < p.w - 1; i++ {
		for d := 0; d < min(i+1, p.w - i - 1); d++ {
			if p.cols[i - d] != p.cols[i + d + 1] {
				goto NO
			}
		}
		note += i + 1
	NO:
	}
	return
}

// note for symmetries around horizontal lines
func hNote(p Pattern) (note int) {
	for i := 0; i < p.h - 1; i++ {
		for d := 0; d < min(i+1, p.h - i - 1); d++ {
			if p.rows[i - d] != p.rows[i + d + 1] {
				goto NO
			}
		}
		note += i + 1
	NO:
	}
	return
}

func Note(p Pattern) (note int) {
	if note = vNote(p); note == 0 {
		note = 100 * hNote(p)
	}
	return
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

//////////// Part 2

func part2(lines []string) (notes int) {
	patterns := parse2(lines)
PATTERN:
	for saidx, sa := range patterns {
		VPf("  [%d] %d x %d\n", saidx, sa.w, sa.h)
		var note0, note int
		note0 = saNote(sa, 0)
		if note0 == 0 {
			panicSA0("No symmetry found", sa)
		}
		// now, try changing a smudge one by one
		for p, b := range sa.a {
			if b {	// invert potential smudge
				sa.a[p] = false
			} else {
				sa.a[p] = true
			}
			note = saNote(sa, note0)
			sa.a[p] =b
			if note != 0 { // OK, found
				VPf("  Pattern [%d], smudge at %v, new note: %d (was %d)\n", saidx, sa.Vector(p), note, note0)
				notes += note
				continue PATTERN
			}
		}
		panic("No alternate reflection line found for pattern " + itoa(saidx))
	}
	return
}

func parse2(lines []string) (patterns []Scalarray0[bool]) {
	start := 0
	for l, line := range lines {
		if line == "" {
			patterns = append(patterns, parseScalarray(lines[start:l]))
			start = l+1
		}
	}
	return
}
	
func parseScalarray(lines []string) (p Scalarray0[bool]) {
	p = makeScalarray0[bool](len(lines[0]), len(lines))
	pos := 0
	for _, line := range lines {
		for _, c := range line {
			if c == '#' {
				p.a[pos] = true
			}
			pos++
		}
	}
	return
}

// find the dividing line, other than the one specified by its note in avoid
func saNote(sa Scalarray0[bool], avoid int) int {
	// find vertical line
	for i := 0; i < sa.w - 1; i++ {
		for d := 0; d < min(i+1, sa.w - i - 1); d++ {
			if sa.colsEqual(i - d, i + d + 1) == 0 {
				goto NOVERT
			}
		}
		// found. Note would be i+1, but continue if we found avoid
		if i+1 != avoid {
			return i + 1
		}
	NOVERT:
	}
	// find horizontal line
	for i := 0; i < sa.h - 1; i++ {
		for d := 0; d < min(i+1, sa.h - i - 1); d++ {
			if sa.rowsEqual(i - d, i + d + 1) == 0 {
				goto NOHORIZ
			}
		}
		// found. Note would be (i+1)*100, but continue if we found avoid
		if (i+1)*100 != avoid {
			return (i+1)*100
		}
	NOHORIZ:
	}
	return 0					// 0 is not a valid note
}

// comparisons of rows & cols. I cannot return a bool because of a type error
func (sa *Scalarray0[bool]) rowsEqual(y1, y2 int) int {
	for x, p1, p2 := 0, sa.Pos(0, y1), sa.Pos(0, y2); x < sa.w; x, p1, p2 = x+1, p1+1, p2+1 {
		if sa.a[p2] != sa.a[p1] {
			return 0
		}
	}
	return 1
}

func (sa *Scalarray0[bool]) colsEqual(x1, x2 int) int {
	for p1, p2 := sa.Pos(x1, 0), sa.Pos(x2, 0); p1 < sa.w*sa.h; p1, p2 = p1+sa.w, p2+sa.w {
		if sa.a[p2] != sa.a[p1] {
			return 0
		}
	}
	return 1
}


func panicSA0(label string, sa Scalarray0[bool]) {
	s := fmt.Sprintf("%s for pattern;\n", label)
	VPScallarray0Bool(s, sa)
	panic("PANIC!")
}


//////////// Common Parts code

//////////// PrettyPrinting & Debugging functions

// print the pattern
func VPpattern(p Pattern) {
	for _, r := range p.rows {
		fmt.Println(r)
	}
}

// print with a dividing vertical line at i
func VPpatternV(p Pattern, i int) {
	for _, r := range p.rows {
		for x, c := range r {
			fmt.Print(string(c))
			if x == i {
				fmt.Print("|")
			}
		}
		fmt.Println()
	}
}

// print with a dividing horizontal line at i
func VPpatternH(p Pattern, i int) {
	for y, r := range p.rows {
		fmt.Println(r)
		if y == i {
			for j := 0; j < p.w; j++ {
				fmt.Print("-")
			}
			fmt.Println()
		}
	}
}

func panicPattern(label string, p Pattern) {
	fmt.Printf("%s for pattern;\n", label)
	VPpattern(p)
	panic("PANIC!")
}
