// Adventofcode 2022, d03, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 input 8105
// TEST: input 2363
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
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
	sum := 0
	for _, line := range lines {
		sum += findDuplicatePriority(line)
	}
	return sum
}

//////////// Part 2
func part2(lines []string) int {
	sum := 0
	for i :=0; i < len(lines) - 2; i += 3 { // -2 skips optional extra line at end
		sum += groupBadgePriority(lines, i)
	}
	return sum
}

//////////// Common Parts code

func priority(b byte) int {
	if b <= 90 {				// A-Z
		return int(b) - 64 + 26
	} else {					// a-z
		return int(b) - 96
	}
}

//////////// Part1 functions

func findDuplicatePriority(rucksack string) int {
	l := len(rucksack) / 2
	if l == 0 {
		return 0				// ignore empty lines
	}
	prios := make([]int, l*2, l*2)
	for i, byte := range []byte(rucksack) {
		prios[i] = priority(byte)
	}
	for i := 0; i < l; i++ {
		for j := l; j < l*2; j++ {
			if prios[i] == prios[j] {
				VP("rucksack", rucksack, "==>", string(rucksack[i]), prios[i], "at positions", i, j)
				return prios[i]
			}
		}
	}
	log.Fatalf("Common element not found in rucksack: \"%s\"", rucksack)
	os.Exit(1)
	return 0					// unreached
}

//////////// Part2 functions

func groupBadgePriority(lines []string, idx int) int {
	var groupPrios [3][]int
	for i := 0; i < 3; i++ {
		line := lines[idx + i]
		groupPrios[i] = make([]int, 53, 53)
		if i < 2 {				// register prios of rucksacks of elves #1 and #2
			for j := 0; j < len(line); j++ {
				groupPrios[i][priority(line[j])] = 1
			}
		} else {				// as soon as we found a common one in #3, return
			for j := 0; j < len(line); j++ {
				p := priority(line[j])
				if groupPrios[0][p] == 1 && groupPrios[1][p] == 1 {
					VP("group #", idx/3, "==> badge:", string(line[j]), p)
					return p
				}
			}
		}			
	}
	log.Fatalf("Badge not found in group starting at rucksack: \"%s\"", lines[idx])
	os.Exit(2)
	return 0					// unreached
}
		
