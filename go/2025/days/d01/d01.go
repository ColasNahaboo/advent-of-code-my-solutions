// Adventofcode 2025, d01, in go. https://adventofcode.com/2025/day/01
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 3
// TEST: example 6
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"fmt"
	"regexp"
	// "flag"
	// "slices"
)

//////////// Options parsing & exec parts

func main() {
	ExecOptions(2, NoXtraOpts, part1, part2)
}

//////////// Part 1

func part1(lines []string) (res int) {
	rots := parse(lines)
	n := 50
	for _, rot := range rots {
		n = turn(n, rot)
		VPf("%d ==> %d\n", rot, n)
		if n == 0 {
			res++
		}
	}
	return 
}

//////////// Part 2

func part2(lines []string) (zeros int) {
	rots := parse(lines)
	n := 50
	for _, rot := range rots {
		if rot > 0 {
			zeros += (n + rot) / 100
		} else {
			next := ((100 - n) % 100) - rot // this flips the dial horizontally
			zeros += next / 100
		}
		n = turn(n, rot)
		VPf("Rot = %d, dial at n = %d, zeros = %d\n", rot, n, zeros)
	}
	return 
}

//////////// Common Parts code

// turn the dial
func turn(n, rot int) (n2 int) {
	n2 = (n + rot) % 100
	if n2 < 0 {
		n2 += 100
	}
	return
}

func parse(lines []string) (rots []int) {
	reline := regexp.MustCompile("([LR])([[:digit:]]+)")
	for _, line := range lines {
		line := reline.FindStringSubmatch(line)
		n := atoi(line[2])
		if n == 0 {
			panic("Rotation by zero")
		}
		if line[1] == "L" {
			n *= -1
		}
		rots = append(rots, n)
	}
	return
}

//////////// PrettyPrinting & Debugging functions. See also the VPx functions.

func DEBUG() {
	if ! debug { return }
	fmt.Println("DEBUG!")
	// ad hoc debug code here
}
