// Adventofcode 2015, day25, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 input 2650453
// TEST: input
package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"
)

var verbose bool

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	flag.Parse()
	verbose = *verboseFlag
	infile := "input.txt"
	var x, y int
	switch flag.NArg() {
	case 1:
		infile = flag.Arg(0)
		fallthrough
	case 0:
		lines := fileToLines(infile)
		re := regexp.MustCompile(`row\s+([[:digit:]]+),\s*column\s+([[:digit:]]+)`)
		s := re.FindStringSubmatch(lines[0])
		if s != nil {
			x = atoi(s[2]) // reversed: row, column = y, x
			y = atoi(s[1])
		} else {
			log.Fatalf("Syntax error: %v\n", lines[0])
		}
	case 2:
		x = atoi(flag.Arg(0))
		y = atoi(flag.Arg(1))
	default:
		log.Fatalf("Bad number of Arguments: %v\n", flag.NArg())
	}

	var result int
	if *partOne {
		VP("Running Part1")
		result = part1(x, y)
	}
	fmt.Println(result)
}

//////////// Part 1

func part1(x, y int) int {
	n := numAt(x, y)
	VPf("Computing code at (%v, %v), the %vth position\n", x, y, n)
	return codeN(numAt(x, y) - 1)
}

//////////// Common Parts code

// from a position, walk back to 1,1 and count steps: The number of the position

func numAt(x, y int) int {
	v := 1
	for {
		if x == 1 {
			if y == 1 {
				return v
			} else {
				x = y - 1
				y = 1
			}
		} else {
			x--
			y++
		}
		v++
	}
}

// generates codes from 1,1 to a position number, return last.

func codeN(n int) int {
	code := 20151125
	for i := 1; i <= n; i++ {
		code = (code * 252533) % 33554393
	}
	return code
}

//////////// Part1 functions

//////////// Part2 functions
