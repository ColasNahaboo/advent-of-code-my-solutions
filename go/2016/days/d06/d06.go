// Adventofcode 2016, d06, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 input afwlyyyq
// TEST: input bhkzekao
package main

import (
	"flag"
	"fmt"
	"log"
	"sort"
	// "regexp"
)

var verbose bool
type LetterCount struct {
	letter byte
	count int
}

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

	var result string
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

func part1(lines []string) (message string) {
	return letterFreqs(lines, false)
}

//////////// Part 2
func part2(lines []string) (message string) {
	return letterFreqs(lines, true)
}
//////////// Common Parts code

// least == true if we get the least frequent letter, otherwise most frequent
func letterFreqs(lines []string, least bool) (message string) {
	ncols := len(lines[0])
	byteCounts := make([]map[byte]int, ncols, ncols)
	for c := 0; c < ncols; c++ {
		byteCounts[c] = make(map[byte]int, 0)
	}
	for _, line := range lines {
		if len(line) != ncols {
			log.Fatalln("Syntax error: " + itoa(len(line)))
		}
		for c := 0; c < ncols; c++ {
			byteCounts[c][line[c]] ++
		}
	}
	for c := 0; c < ncols; c++ {
		// we create a slice of LetterCounts to sort it
		lcs := make([]LetterCount, len(byteCounts[c]))
		i := 0
		for char, count := range byteCounts[c] {
			lcs[i] = LetterCount{char, count}
			i++
		}
		sort.Slice(lcs, func(i, j int) bool {
			if lcs[i].count > lcs[j].count {
				return true
			} else if (lcs[i].count == lcs[j].count) && (lcs[i].letter < lcs[j].letter) {
				return true
			} else {
				return false
			}
		})
		if least {
			message += string(lcs[len(lcs) - 1].letter)
		} else {
			message += string(lcs[0].letter)
		}
	}
	return
}
	

//////////// Part1 functions

//////////// Part2 functions
