// Adventofcode 2018, d11, in go. https://adventofcode.com/2018/day/11
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example "33_45"
// TEST: -1 example2 "21_61"
// TEST: example
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

// we output x_y instead of x,y because of our test system than use commas

package main

import (
	"fmt"
	// "regexp"
	// "flag"
	// "slices"
)

//////////// Options parsing & exec parts

func main() {
	XOptsUsage(3, "part2, but coded naively (slower: more than 1mn)")
	ExecOptionsString(2, NoXtraOpts, part1, part2, part3)
}

//////////// Part 1

// no need to actually build the grid for this simple case

func part1(lines []string) string {
	serial := parse(lines)
	var maxpl, mx, my int
	for x := 1; x < 299; x++ {
		for y := 1; y < 299; y++ {
			pl := 0
			for sx := 0; sx < 3; sx++ { //  the 3x3 square
				for sy :=0; sy < 3; sy++ {
					pl += PowerLevel(serial, x+sx, y+sy)
				}
			}
			if pl > maxpl {
				maxpl, mx, my = pl, x, y
			}
		}
	}
	return itoa(mx) + "_" + itoa(my)
}

//////////// Part 2

// We use a power level cache per coord
// and for each coord, we compute each increasing square size incrementally
// and as an heuristics, once the square size is big enough (10), since all power
// levels of a coord are in [-5, 4], we abort if adding the next col+row could
// not possibly beat the max size record, considering it will be the same for
// all subsequent sizes

func part2(lines []string) (res string) {
	serial := parse(lines)
	var maxpl, mx, my, ms int
	b := MakeBoard[int](301, 301) // the grid of power levels used as cache
	for x := 1; x < 300; x++ {
		for y := 1; y < 300; y++ {
			b.a[x][y] = PowerLevel(serial, x, y)
		}
	}
	// now, examine all the possible squares
	for x := 1; x < 300; x++ {	// all positions x,y of top corner
		for y := 1; y < 300; y++ {
			availableSize := min(301 - x, 301 - y)
			ppl := 0			// cached previous size power level
			for s := 1; s < availableSize; s++ { // all sizes fitting free space
				pl := ppl
				// compute only the added right & bottom sides for s*s square
				for sx := 0; sx < s; sx++ { //  the bottom row: sy = y+s-1
					pl += b.a[x+sx][y+s-1]
				}
				for sy := 0; sy < s-1; sy++ { // the right side, less SE corner
					pl += b.a[x+s-1][y+sy]
				}
				if pl > maxpl {
					maxpl, mx, my, ms = pl, x, y, s
				}
				ppl = pl
				if s > 9 && maxpl - ppl <= (s*2+1) * 4 {
					// there is no hope to reach maxpl on next size, abort
					break
				}
			}
		}
	}
	return itoa(mx) + "_" + itoa(my) + "_" + itoa(ms)
}

//////////// Part 3

// brute force naive implmentation, with just a power level cache per coord

func part3(lines []string) (res string) {
	serial := parse(lines)
	var maxpl, mx, my, ms int
	b := MakeBoard[int](301, 301) // we build the grid of power levels as a cache
	for x := 1; x < 300; x++ {
		for y := 1; y < 300; y++ {
			b.a[x][y] = PowerLevel(serial, x, y)
		}
	}
	// now, examine all the possible squares
	for x := 1; x < 300; x++ {	// all positions x,y of top corner
		for y := 1; y < 300; y++ {
			availableSize := min(301 - x, 301 - y)
			for s := 1; s < availableSize; s++ { // all sizes fitting free space
				pl := 0
				for sx := 0; sx < s; sx++ { //  the s x s square
					for sy :=0; sy < s; sy++ {
						pl += b.a[x+sx][y+sy]
					}
				}
				if pl > maxpl {
					maxpl, mx, my, ms = pl, x, y, s
				}
			}
		}
	}
	return itoa(mx) + "_" + itoa(my) + "_" + itoa(ms)
}

//////////// Common Parts code

func PowerLevel(serial, x, y int) (pl int) {
	rackid := x + 10
	pl = rackid * y
	pl += serial
	pl *= rackid
	pl = (pl / 100) % 10
	pl -= 5
	return
}

func parse(lines []string) int {
	return atoi(lines[0])
}

//////////// PrettyPrinting & Debugging functions. See also the VPx functions.

func DEBUG() {
	if ! debug { return }
	fmt.Println("DEBUG!")
	// ad hoc debug code here
}
