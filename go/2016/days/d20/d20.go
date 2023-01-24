// Adventofcode 2016, d20, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 example 3
// TEST: -1 input 14975795
// TEST: example 4294967288
// TEST: input 101

// Part1:
// the lowest (or biggest) non-blocked value is on contiguous to a blocked interval
// I thus make two pass: first collect the possible values, then test them

package main

import (
	"flag"
	"fmt"
	"sort"
	// "regexp"
)

var verbose bool
const maxval = 4294967295

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
	blocks, values := parse(lines)
	sort.Ints(values)			// sort in place in increasing order
	VP(blocks)
	VP(values)
	for _, val := range values {
		for _, block := range blocks {
			if val >= block[0] && val <= block[1] {
				goto BLOCKED
			}
		}
		return val
	BLOCKED:
	}
	return -1					// FAIL
}

//////////// Part 2
func part2(lines []string) (allowed int) {
	blocks, _ := parse(lines)
	intervals := [][2]int{{0, maxval}}	// todo list: a stack of IP intervals to test
	for {
		if len(intervals) == 0 { break;} // no intervals left to test in the todo list
		ips := intervals[len(intervals)-1] // pop next one
		intervals = intervals[:len(intervals)-1]
		for _, block := range blocks {
			if ips[1] < block[0] || ips[0] > block[1] { // disjoints
				continue
			}
			if ips[0] >= block[0] {
				if ips[1] <= block[1] { // ips fully included in block
					goto BLOCKED
				} else {		// truncate low
					ips[0] = block[1] + 1
					if isEmpty(ips) { goto BLOCKED;}
				}
			} else {
				if ips[1] <= block[1] { // truncate high
					ips[1] = block[0] - 1
					if isEmpty(ips) { goto BLOCKED;} // ips empty
				} else {								  // block cuts ips in two intervals
					ips2 := [2]int{block[1]+1, ips[1]}
					ips[1] = block[0] - 1
					intervals = append(intervals, ips2) // push high half onto todo list
					continue							// continue testing the lower half
				}
			}
		}
		allowed += ips[1] - ips[0] + 1 // all IPs in ips are allowed
		BLOCKED:				// interval ips fully blocked
	}
	return
}

func isEmpty(ips [2]int) bool {
	return ips[0] > ips[1]
}

//////////// Common Parts code

func parse(lines []string) (blocks [][2]int, values []int) {
	var low, high, lineno int
	for _, line := range lines {
		lineno++
		if n, _ := fmt.Sscanf(line, "%d-%d", &low, &high); n == 2{
			blocks = append(blocks, [2]int{low, high})
			if low > 0 {
				values = append(values, low - 1)
			}
			if high < maxval {
				values = append(values, high + 1)
			}
		} else {
			panic(fmt.Sprintf("Syntax error line %d: %s\n", lineno, line))
		}
	}
	return
}

//////////// Part1 functions

//////////// Part2 functions
