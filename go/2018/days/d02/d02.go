// Adventofcode 2018, d02, in go. https://adventofcode.com/2018/day/02
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 12
// TEST: example abcde
// TEST: example2 fgij
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"fmt"
	"os"
	// "regexp"
	// "flag"
	// "slices"
)

//////////// Options parsing & exec parts

func main() {
	ExecOptions(2, NoXtraOpts, part1, part2)
}

//////////// Part 1

func part1(lines []string) (res int) {
	ids := parse(lines)
	ids2, ids3 := 0, 0
	for _, id := range ids {
		letters := make(map[rune]int)
		for _, r := range id {
			letters[r]++
		}
		is2, is3 := false, false
		for _, occurs := range letters {
			switch occurs {
			case 2: is2 = true
			case 3: is3 = true
			}
		}
		if is2 { ids2++ }
		if is3 { ids3++ }
	}
	return ids2 * ids3
}

//////////// Part 2

func part2(lines []string) (res int) {
	for i := range len(lines[0]) {
		seen := make(map[string]bool)
		for _, id := range lines {
			id1 := id[:i] + id[i+1:]
			if seen[id1] {
				fmt.Println(id1)
				os.Exit(0)
			}
			seen[id1] = true
		}
	}
	return 
}

//////////// Common Parts code

func parse(lines []string) (ids [][]rune) {
	for _, line := range lines {
		ids = append(ids, []rune(line))
	}
	return
}

//////////// PrettyPrinting & Debugging functions. See also the VPx functions.

func DEBUG() {
	if ! debug { return }
	fmt.Println("DEBUG!")
	// ad hoc debug code here
}
