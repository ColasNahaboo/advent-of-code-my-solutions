// Adventofcode YYYY, dayNN, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 input 654
// TEST: input 57

// we call containers "cans", as it is easier to write :-)

package main

import (
	"flag"
	"fmt"
	// "regexp"
)

var verbose bool

func main() {
	partOne := flag.Bool("1", false, "run part one code, instead of part 2 (default)")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	totalFlag := flag.Int("t", 150, "Total liters to store")
	flag.Parse()
	verbose = *verboseFlag
	total := *totalFlag
	infile := "input.txt"
	if flag.NArg() > 0 {
		infile = flag.Arg(0)
	}
	lines := fileToLines(infile)
	cans := readCans(lines)

	var result int
	if *partOne {
		VP("Running Part1")
		result = part1(cans, total)
	} else {
		VP("Running Part2")
		result = part2(cans, total)
	}
	fmt.Println(result)
}

//////////// Part 1

func part1(cans []int, total int) int {
	// A combination is a binary number of size l=len(cans) bits,
	// max being: 2^l
	// It represents the bitset of the used cans
	// explained in https://www.topcoder.com/blog/generating-combinations/
	fullset := intPower64(2, len(cans))
	count := 0
	var i uint64
	for i = 0; i < fullset; i++ {
		VP("Testing", i)
		allcans := 0
		for canidx := 0; canidx < len(cans); canidx++ {
			if i&(uint64(1)<<canidx) != 0 {
				allcans += cans[canidx]
			}
		}
		if allcans == total {
			count++
		}
	}
	return count
}

//////////// Part 2
func part2(cans []int, total int) int {
	// same code as part1, but now we count the possible combinations
	// separately, per their number of cans, in the counts array of counts
	counts := make([]int, len(cans))
	fullset := intPower64(2, len(cans))
	var i uint64
	for i = 0; i < fullset; i++ {
		VP("Testing", i)
		allcans := 0
		setbits := 0
		for canidx := 0; canidx < len(cans); canidx++ {
			if i&(uint64(1)<<canidx) != 0 {
				allcans += cans[canidx]
				setbits++
			}
		}
		if allcans == total {
			counts[setbits]++
		}
	}
	for i = 0; i < fullset; i++ {
		if counts[i] > 0 {
			fmt.Println("Min number of cans:", i)
			return counts[i]
		}
	}
	return -1
}

//////////// Common Parts code

// read cans into an array with values their capacity
func readCans(lines []string) []int {
	cans := make([]int, 0)
	for _, line := range lines {
		cans = append(cans, atoi(line))
	}
	return cans
}

//////////// Part1 functions

//////////// Part2 functions
