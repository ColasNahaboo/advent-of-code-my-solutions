// Adventofcode 2022, d06, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 input 1892
// TEST: input 2313
package main

import (
	"flag"
	"fmt"
	"log"
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
	return findDistinctN(lines[0], 4)
}

//////////// Part 2
func part2(lines []string) int {
	return findDistinctN(lines[0], 14)
}

//////////// Common Parts code

// we keep a map of the last N chars, "last", with value the count of occurences
// on each step we increase the count of the entering one, decrease the leaving
// we delete the entries reaching 0
// the number of unique chars is the the length of the map

func findDistinctN(line string, n int) int {
	if len(line) < n {
		log.Fatalln("buffer length must be at least", n)
	}
	last := map[byte]int{}
	for i := 0; i < n; i++ {
		last[line[i]] ++
	}
	if len(last) == n {
		return n
	}
	for i := n; i < len(line); i++ {
		last[line[i]] ++
		last[line[i - n]] --
		if last[line[i - n]] == 0 {
			delete(last, line[i - n])
		}
		if len(last) == n {
			return i+1
		}
	}
	return 0
}

//////////// Part1 functions

//////////// Part2 functions
