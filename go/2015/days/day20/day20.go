// Adventofcode 2015, day20, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 input 831600
// TEST: input 884520
package main

import (
	"flag"
	"fmt"
	"log"
	// "regexp"
)

var verbose bool
var goal int

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	goalFlag := flag.Int("g", 0, "If non-zero, use it as the input number instead of the one in the input.txt file")
	flag.Parse()
	verbose = *verboseFlag
	infile := "input.txt"
	if flag.NArg() > 0 {
		infile = flag.Arg(0)
	}
	lines := fileToLines(infile)
	if *goalFlag != 0 {
		goal = *goalFlag
	} else {
		goal = atoi(lines[0])
	}
	if goal < 100 {
		log.Fatal("Goal must be >= 100")
	}

	var result int
	if *partOne {
		VP("Running Part1")
		result = part1(goal)
	} else {
		VP("Running Part2")
		result = part2(goal)
	}
	fmt.Println(result)
}

//////////// Part 1

func part1(goal int) int {
	presentsAtHouses := sendElves(goal)
	for i := 1; i < len(presentsAtHouses); i++ {
		if presentsAtHouses[i] >= goal {
			return i
		}
	}
	return -1
}

//////////// Part 2
func part2(goal int) int {
	presentsAtHouses := sendElvesLimited(goal, 50)
	for i := 1; i < len(presentsAtHouses); i++ {
		if presentsAtHouses[i] >= goal {
			return i
		}
	}
	return -1
}

//////////// Common Parts code

//////////// Part1 functions

// let loose the elves, make them fill houses with presents
func sendElves(goal int) []int {
	payload := 10
	maxelf := goal / payload //  worst case, one divisor for a prime
	presents := make([]int, maxelf+1)
	for elf := 1; elf <= maxelf; elf++ {
		elfload := elf * payload
		for house := elf; house <= maxelf; house += elf {
			presents[house] += elfload
		}
	}
	return presents
}

//////////// Part2 functions

// let loose the elves, make them fill houses with presents
// but they run only limit times
func sendElvesLimited(goal, limit int) []int {
	payload := 11
	maxelf := goal / payload //  worst case, one divisor for a prime
	presents := make([]int, maxelf+1)
	for elf := 1; elf <= maxelf; elf++ {
		deliveries := 0
		elfload := elf * payload
		for house := elf; house <= maxelf; house += elf {
			presents[house] += elfload
			deliveries++
			if deliveries >= limit { // elf exhausted
				break
			}
		}
	}
	return presents
}
