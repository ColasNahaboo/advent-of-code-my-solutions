// Adventofcode 2025, d04, in go. https://adventofcode.com/2025/day/04
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 13
// TEST: example 43
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"fmt"
)

//////////// Options parsing & exec parts

func main() {
	ExecOptions(2, NoXtraOpts, part1, part2)
}

//////////// Part 1

func part1(lines []string) int {
	b := parse(lines)
	return reachables(b)
}

//////////// Part 2

func part2(lines []string) (res int) {
	b := parse(lines)
	for {						// iterate, removing reachables at each pass
		removed := 0
		b2 := b.Map(func (b *Board[bool], p Point) bool {
			if b.Get(p) {
				if isReachable(b, p) {
					removed++
				} else {
					return true
				}
			}
			return false
		})
		if removed == 0 {
			break
		}
		res += removed
		b = b2
	}
	return
}

//////////// Common Parts code

func reachables(b *Board[bool]) (res int) {
	b.Apply(func (b *Board[bool], p Point) {
		if b.Get(p) && isReachable(b, p) {
			res++
		}})
	return
}
	

func isReachable(b *Board[bool], p Point) bool {
	adjacents := 0
	for _, dir := range DirsAll {
		if b.GetOr(p.Add(dir), false) {
			adjacents++
		}
	}
	if adjacents < 4 {
		return true
	}
	return false
}
			
func parse(lines []string) *Board[bool] {
	b := MakeBoard[bool](len(lines[0]), len(lines))
	for y, line := range lines {
		for x, roll := range line {
			if roll == '@' {
				b.Set(Point{x, y}, true)
			}
		}
	}
	return &b
}

//////////// PrettyPrinting & Debugging functions. See also the VPx functions.

func DEBUG() {
	if ! debug { return }
	fmt.Println("DEBUG!")
	// ad hoc debug code here
}

func VPcell(c bool) string {
	if c {
		return "@"
	}
	return "."
}
