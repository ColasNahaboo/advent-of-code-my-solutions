// Adventofcode 2017, d04, in go. https://adventofcode.com/2017/day/04
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 2
// TEST: example 2
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"flag"
	"fmt"
	"regexp"
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
		infile = fileMatch("input,[0-9]*,[0-9]*.test")
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

func part1(lines []string) (sum int) {
	re := regexp.MustCompile("[[:lower:]]+")
LINE:
	for _, line := range lines {
		words := re.FindAllString(line, -1)
		seen := make(map[string]bool, 0)
		for _, word :=range words {
			if seen[word] {
				continue LINE
			}
			seen[word] = true
		}
		sum++
	}
	return
}

//////////// Part 2
func part2(lines []string) (sum int) {
	re := regexp.MustCompile("[[:lower:]]+")
LINE:
	for _, line := range lines {
		words := re.FindAllString(line, -1)
		seen := make(map[string]bool, 0)
		for _, word :=range words {
			letters := wordLetters(word)
			if seen[letters] {
				continue LINE
			}
			seen[letters] = true
		}
		sum++
	}
	return
}

// String of letters present in word followed by the number of their occurences
// E.g: ioii ==> i3o1, msukkku ==> k3m1s1u2
func wordLetters(word string) (occurs string) {
	letters := [26]int{}
	for _, c := range word {
		letters[c - 'a']++
	}
	for i, n := range letters {
		if n > 0 {
			occurs = occurs + string('a' + i) + itoa(n)
		}
	}
	return
}


//////////// Common Parts code

//////////// PrettyPrinting & Debugging functions
