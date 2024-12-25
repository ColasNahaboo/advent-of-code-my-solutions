// Adventofcode 2024, d25, in go. https://adventofcode.com/2024/day/25
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 3
// TEST: example
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"fmt"
	// "regexp"
	// "flag"
	// "slices"
)

type Key [5]byte				// the key elevations or lock spaces

//////////// Options parsing & exec parts

func main() {
	ExecOptions(2, NoXtraOpts, part1, part2)
}

//////////// Part 1

func part1(lines []string) (res int) {
	locks, keys := parse(lines)
	VPf("%d Locks: %v\n", len(locks), locks)
	VPf("%d Keys:  %v\n", len(keys), keys)
	for _, l := range locks {
		for _, k := range keys {
			if k.Fits(l) {
				res++
			}
		}
	}
	return 
}

func (k Key) Fits(l Key) bool {
	for i, pin := range k {
		if pin > l[i] {
			return false
		}
	}
	return true
}

//////////// Part 2

func part2(lines []string) (res int) {
	fmt.Println("No part2 for Day 25, Happy XMas!")
	return 
}

//////////// Common Parts code

func parse(lines []string) (locks, keys []Key) {
	var y int
	for i := 0; i < len(lines) - 7; i += 8 { // 5 key + 2 border + 1 blank line
		b := ParseBoard[bool](lines[i:i+7], ParseCellBool)
		k := Key{}
		if b.a[0][0] {			// lock
			for x := 0; x < 5; x++ {
				for y = 1; y < 7; y++ {
					if ! b.a[x][y] { // top of the hole
						k[x] = byte(6 - y)
						break
					}
				}
			}
			locks = append(locks, k)
		} else {				// key
			for x := 0; x < 5; x++ {
				for y = 1; y < 7; y++ {
					if b.a[x][y] { // top of the key
						k[x] = byte(6 - y)
						break
					}
				}
			}
			keys = append(keys, k)
		}
	}
	return
}

//////////// PrettyPrinting & Debugging functions. See also the VPx functions.

func DEBUG() {
	if ! debug { return }
	fmt.Println("DEBUG!")
	// ad hoc debug code here
}
