// Adventofcode 2018, d03, in go. https://adventofcode.com/2018/day/03
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 4
// TEST: example 3
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"fmt"
	"regexp"
	// "flag"
	// "slices"
)

type Claim struct {
	id int						// IDs start with 1. So index in claims[] is id-1
	x, y int
	w, h int
	overlap bool				// for part 2: are we overlapping any other?
}

//////////// Options parsing & exec parts

func main() {
	ExecOptions(2, NoXtraOpts, part1, part2)
}

//////////// Part 1

func part1(lines []string) (res int) {
	w, h, claims := parse(lines)
	b := MakeBoard[int](w, h)	// number of claims on each point
	for _, c := range claims {
		for x := c.x; x < c.x + c.w; x++ {
			for y := c.y; y < c.y + c.h; y++ {
				if b.a[x][y] == 1 {
					res++
				}
				b.a[x][y]++
			}
		}
	}
	return 
}

//////////// Part 2

func part2(lines []string) (res int) {
	w, h, claims := parse(lines)
	b := MakeBoard[int](w, h)	// ID of the first claim on the point
	for ic, c := range claims {
		for x := c.x; x < c.x + c.w; x++ {
			for y := c.y; y < c.y + c.h; y++ {
				if b.a[x][y] > 0 { // claimed
					VPf("  [%d] and [%d] overlap at %d,%d\n",b.a[x][y],c.id,x,y)
					claims[ic].overlap = true // mark both as overlapping
					claims[b.a[x][y]-1].overlap = true
				} else {
					b.a[x][y] = c.id
				}
			}
		}
	}
	VP(claims)
	for _, c := range claims {
		if ! c.overlap {
			return c.id
		}
	}
	return
}

//////////// Common Parts code

func parse(lines []string) (maxw, maxh int, claims []Claim) {
	renum := regexp.MustCompile("[[:digit:]]+") // example code body, replace.
	for _, line := range lines {
		m := atoil(renum.FindAllString(line, -1))
		claims = append(claims, Claim{m[0], m[1], m[2], m[3], m[4], false})
		maxw = max(maxw, m[1] + m[3])
		maxh = max(maxh, m[2] + m[4])
	}
	return
}

//////////// PrettyPrinting & Debugging functions. See also the VPx functions.

func DEBUG() {
	if ! debug { return }
	fmt.Println("DEBUG!")
	// ad hoc debug code here
}
