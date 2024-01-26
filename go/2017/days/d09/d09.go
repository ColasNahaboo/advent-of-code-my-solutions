// Adventofcode 2017, d09, in go. https://adventofcode.com/2017/day/09
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example1 1
// TEST: -1 example2 6
// TEST: -1 example3 5
// TEST: -1 example4 16
// TEST: -1 example5 1
// TEST: -1 example6 9
// TEST: -1 example7 9
// TEST: -1 example8 3
// TEST: example
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"flag"
	"fmt"
)

var verbose bool

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	flag.Parse()
	verbose = *verboseFlag
	var infile string
	if flag.NArg() > 0 {
		infile = flag.Arg(0)
	} else {
		infile = fileMatch("input,[[:alnum:]]*,[[:alnum:]]*.test")
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
	return parse(lines[0])
}

//////////// Part 2
var garbage int

func part2(lines []string) int {
	parse(lines[0])
	return garbage
}

//////////// Common Parts code

func parse(line string) int {
	l := []byte(line)
	score, _ := parseGroup(0, 0, l)
	return score
}

// parses a group starting at pos with a "{" in l
// returns its score plus the score of all contained groups
func parseGroup(pos, level int, l []byte) (score, p int) {
	level++
	score += level
	p = pos + 1
	for {
		switch l[p] {
		case '}': return
		case '{':
			var subscore int
			subscore, p = parseGroup(p, level, l)
			score += subscore
		case '!': p++
		case '<': p = skipGarbage(p, l)
		}
		p++
	}
}
			
func skipGarbage(pos int, l []byte) (p int) {
	p = pos + 1
	for {
		switch l[p] {
		case '!': p++
		case '>': return
		default: garbage++		// for part2
		}
		p++
	}
}
	

//////////// PrettyPrinting & Debugging functions
