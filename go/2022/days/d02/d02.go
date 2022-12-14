// Adventofcode 2022, d02, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 input 13682
// TEST: input 12881
package main

import (
	"flag"
	"fmt"
	"regexp"
	"log"
	"os"
)

var verbose bool
var reline *regexp.Regexp

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
	reline = regexp.MustCompile("^([ABC])[[:space:]]+([XYZ])")

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
	score := 0
	for _, line := range lines {
		round := reline.FindStringSubmatch(line)
		if round  == nil {
			if line != "" {
				log.Fatalf("syntax error on line: \"%s\"", line)
			}
		} else {
			score += line_score(round[1], round[2])
		}
	}
	return score
}

//////////// Part 2

func part2(lines []string) int {
	score := 0
	for _, line := range lines {
		round := reline.FindStringSubmatch(line)
		if round  == nil {
			if line != "" {
				log.Fatalf("syntax error on line: \"%s\"", line)
			}
		} else {
			score += line_score(round[1], myplay(round[1], round[2]))
		}
	}
	return score
}

//////////// Common Parts code

func line_score(elf, me string) int {
	score := 0
	// we consider only wins & draw, as losing adds 0
	if me == "X" {		// you played rock
		score += 1
		if elf == "C" {
			score += 6 			// win
		} else if elf == "A" {
			score += 3
		}
	} else if me == "Y" {	// you played paper
		score += 2
		if elf == "A" {
			score += 6 			// win
		} else if elf == "B" {
			score += 3
		}
	} else if me == "Z" {	// you played scissors
		score += 3
		if elf == "B" {
			score += 6 			// win
		} else if elf == "C" {
			score += 3
		}
	}
	return score
}

//////////// Part1 functions

//////////// Part2 functions

// what should I play?
func myplay(elf, goal string) (play string) {
	switch goal {
	case "X":					// must lose
		switch elf {
		case "A" : return "Z"
		case "B" : return "X"
		case "C" : return "Y"
		}
	case "Y":					// must draw
		switch elf {
		case "A" : return "X"
		case "B" : return "Y"
		case "C" : return "Z"
		}
	case "Z":					// must win
		switch elf {
		case "A" : return "Y"
		case "B" : return "Z"
		case "C" : return "X"
		}
	}
	log.Fatalf("Invalid order: \"%s\"", goal)
	os.Exit(2)
	return ""					// not reached
}
		
	
