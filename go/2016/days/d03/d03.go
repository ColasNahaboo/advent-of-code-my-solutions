// Adventofcode 2016, d03, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 input 917
// TEST: input 1649
package main

import (
	"flag"
	"fmt"
	"regexp"
	"log"
)

var reline *regexp.Regexp
var verbose bool

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	flag.Parse()
	verbose = *verboseFlag
	infile := "input.txt"
	if flag.NArg() > 0 {
		infile = flag.Arg(0)
	}
	lines := fileToLines(infile)
	reline = regexp.MustCompile("([[:digit:]]+)[[:space:]]+([[:digit:]]+)[[:space:]]+([[:digit:]]+)")

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

func part1(lines []string) (possibles int) {
	for _, line := range lines {
		vals := reline.FindStringSubmatch(line)
		if vals  == nil {
			if line != "" {
				log.Fatalf("syntax error on line: \"%s\"", line)
			}
			continue			// ignore empty lines
		}
		if possible(atoi(vals[1]), atoi(vals[2]), atoi(vals[3])) {
			possibles ++
			VPf("possible:   %d %d %d\n", atoi(vals[1]), atoi(vals[2]), atoi(vals[3]))
		} else {
			VPf("IMPOSSIBLE: %d %d %d\n", atoi(vals[1]), atoi(vals[2]), atoi(vals[3]))
		}
	}
	return possibles
}

//////////// Part 2
func part2(lines []string) (possibles int) {
	var vals [3][3]int			// a block of 3 triangles
	for i := 0; i < len(lines); i += 3 {
		// gather values
		for row := 0; row < 3; row++ {
			valsrow := reline.FindStringSubmatch(lines[i+row])
			for col := 0; col < 3; col++ {
				vals[row][col] = atoi(valsrow[col+1])
			}
		}
		// examine them by vertical slices
		for col := 0; col < 3; col++ {
			if possible(vals[0][col], vals[1][col], vals[2][col]) {
				possibles++
			}
		}
	}
	return possibles
}

//////////// Common Parts code

//////////// Part1 functions

func possible(x, y, z int) bool {
	if x + y > z && x + z > y && y + z > x {
		return true
	} else {
		return false
	}
}

//////////// Part2 functions
