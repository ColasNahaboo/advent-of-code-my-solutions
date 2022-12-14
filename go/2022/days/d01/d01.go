// Adventofcode 2022, d01, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 input 74198
// TEST: input 209914

package main

import (
	"flag"
	"fmt"
	// "regexp"
)

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

	// we add an empty line at end, to define the last elf
	lines = append(lines, "")
	
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

func part1(lines []string) int {
	cals := 0
	maxcals := 0
	for _, line := range lines {
		if line == "" {			// end of elf, max one?
			if cals > maxcals {
				maxcals = cals
			}
			cals = 0			// start new elf
		} else {				// accumulate calories
			cals = cals + atoi(line)
		}
	}
	return maxcals
}

//////////// Part 2
func part2(lines []string) int {
	cals := 0
	maxcals := []int{0,0,0}
	for _, line := range lines {
		if line == "" {			// end of elf, max one?
			VP("Elf with cals:", cals)
			for i := 0; i < 3; i++ {
				if cals > maxcals[i] {
					// make room, push the rest down
					for j := 2; j > i; j-- {
						maxcals[j] = maxcals[j-1]
					}
					maxcals[i] = cals
					VP("Calories of top three elves:", maxcals[0], maxcals[1], maxcals[2])
					break
				}
			}
			cals = 0			// start new elf
		} else {				// accumulate calories
			cals = cals + atoi(line)
		}
	}
	VP("Final Calories of top three elves:", maxcals[0], maxcals[1], maxcals[2])
	return maxcals[0] + maxcals[1] + maxcals[2]
}

//////////// Common Parts code

//////////// Part1 functions

//////////// Part2 functions
