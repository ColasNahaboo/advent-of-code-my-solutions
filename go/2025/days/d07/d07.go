// Adventofcode 2025, d07, in go. https://adventofcode.com/2025/day/07
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 21
// TEST: example 40
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

// Input properties:
// S is in the top row
// every even row has no splitters
// the last row has no splitters

package main

import (
	"fmt"
	"strings"
	// "flag"
	// "slices"
)

// Implementation:
// we only consider even rows
// the Manifold is a 2D bool Board provided by point.go, true => splitter
// the beams states are other 2D bool Boards
// E.g:
//            manifold beams      traced beams
// ...S.  ==> _^_^_    ___|_  ==> ___|_
// .....               _____      __|_|
// .^.^.
// .....

//////////// Options parsing & exec parts

func main() {
	ExecOptions(2, NoXtraOpts, part1, part2)
}

//////////// Part 1

func part1(lines []string) (res int) {
	manifold, beams, _ := parse(lines)
	VP(manifold)
	VP(beams)
	return traceBeams(manifold, beams)
}

func traceBeams(m, b *Board[bool]) (splits int) {
	for y := range m.h {
		for x := range m.w {
			if b.a[x][y] {
				if m.a[x][y] { // split
					b.a[x-1][y+1] = true
					b.a[x+1][y+1] = true
					splits++
				} else {		// go through unsplitted
					b.a[x][y+1] = true
				}
			}
		}
	}
	return
}

//////////// Part 2

func part2(lines []string) (res int) {
	manifold, beams, entry := parse(lines)
	VP(manifold)
	VP(beams)
	return traceWorlds(manifold, beams, entry)
}

// part 2 just "charges" each beam with the number of worlds it appears in

func traceWorlds(m, b *Board[bool], entry int) int {
	timelines := MakeBoard[int](b.w, b.h)
	timelines.a[entry][0] = 1
	for y := range m.h {
		for x := range m.w {
			if timelines.a[x][y] > 0 {
				if m.a[x][y] { // split worlds, add to both timelines
					timelines.a[x-1][y+1] += timelines.a[x][y]
					timelines.a[x+1][y+1] += timelines.a[x][y]
				} else {		// go through unsplitted
					timelines.a[x][y+1] += timelines.a[x][y]
				}
			}
		}
	}
	return sumOfInts(timelines.GetRow(timelines.h - 1))
}

func sumOfInts(is []int) (res int) {
	for _, i := range is {
		res += i
	}
	return
}

//////////// Common Parts code

func parse(lines []string) (*Board[bool], *Board[bool], int) {
	// we only store lines with splitters: even ones, starting at 2
	manifold := MakeBoard[bool](len(lines[0]), len(lines)/2 - 1)
	for y := range manifold.h {
		line := lines[y*2 + 2]
		for x, r := range line {
			if r == '^' {
				manifold.a[x][y] = true
			}
		}
	}
	beams := MakeBoard[bool](manifold.w, manifold.h + 1)
	entry := strings.IndexRune(lines[0], 'S')
	beams.a[entry][0] = true
	return manifold, beams, entry
}

//////////// PrettyPrinting & Debugging functions. See also the VPx functions.

func VPworlds(label string, worlds []*Board[bool]) {
	fmt.Printf("%s:\n", label)
	for _, world := range worlds {
		fmt.Printf("  %v\n", world)
	}
}

func DEBUG() {
	if ! debug { return }
	fmt.Println("DEBUG!")
	// ad hoc debug code here
}
