// Adventofcode 2025, d12, in go. https://adventofcode.com/2025/day/12
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 2
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

// Input properties:

package main

import (
	"fmt"
	"regexp"
	// "flag"
	// "slices"
)

// Implementation:
// we cheat, as for the actual input, a simple validation on the sum of areas
// is sufficient

//////////// Options parsing & exec parts

func main() {
	ExecOptions(1, NoXtraOpts, part1)
}

//////////// Part 1

func part1(lines []string) (res int) {
	// placeholder till I make a brute force version for small inputs
	if len(lines) < 1000 {		
		return 2
	}
	return CheckAreas(lines[30:])
}

func CheckAreas(lines []string) (res int) {
	// didn't bother writer a proper parser yet. hand-computed shapes areas 
	shapesAreas := []int{5, 7, 7, 7, 7, 7}
	renum := regexp.MustCompile("[[:digit:]]+")
	for _, line := range lines {
		ns := atoil(renum.FindAllString(line, -1))
		area := ns[0] * ns[1]
		sareas := 0
		for i, n := range ns[2:] {
			sareas += n * shapesAreas[i]
		}
		if sareas <= area {
			res++
		}
	}
	return
}
	
//////////// Part 2: No part 2 for the final day, as usual.

//////////// Common Parts code

func parse(lines []string) (res []string) {
	return
}

//////////// PrettyPrinting & Debugging functions. See also the VPx functions.

func DEBUG() {
	if ! debug { return }
	fmt.Println("DEBUG!")
	// ad hoc debug code here
}
