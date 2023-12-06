// Adventofcode 2023, d04, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: d04-input,RESULT1,RESULT2.test
// TEST: -1 example 13
// TEST: example 30
// And any file named d04-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"flag"
	"fmt"
	"regexp"
)

var verbose bool

var re = regexp.MustCompile("([0-9]+|[|:])")


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

func part1(lines []string) (total int) {
	
	for _, line := range lines {
		win := [100]bool{}		// table of winning numbers
		tokens := re.FindAllString(line, -1)
		i := 0
		points := 0
		// skip the "Card NN:" header
		for ; tokens[i] != ":"; i++ {
		}
		i++
		
		// first, read winning numbers and fill the win table
		for ; i < len(tokens); i++ {
			if tokens[i] == "|" {
				i++
				break
			}
			win[atoi(tokens[i])] = true
		}
		// then look at all the numbers we have and check with the win table
		for ; i < len(tokens); i++ {
			if win[atoi(tokens[i])] {
				if points == 0 {
					points = 1
				} else {
					points *= 2
				}
			}
		}
		total += points
			
	}
	return
}

//////////// Part 2
func part2(lines []string) (ncards int) {
	nc := make([]int, 0)		// for each card ID, how many do we have?
	wc := make([]int, 0)		// how many wins on each card?

	// first, parses all the cards
	for c, line := range lines {
		win := [100]bool{}		// table of winning numbers on card
		tokens := re.FindAllString(line, -1)
		i := 0
		wins := 0
		// skip the "Card NN:" header
		for ; tokens[i] != ":"; i++ {
		}
		i++
		// first, read winning numbers and fill the win table
		for ; i < len(tokens); i++ {
			if tokens[i] == "|" {
				i++
				break
			}
			win[atoi(tokens[i])] = true
		}
		// then look at all the numbers we have and check with the win table
		for ; i < len(tokens); i++ {
			if win[atoi(tokens[i])] {
				wins++
			}
		}
		nc = append(nc, 1)
		wc = append(wc, wins)
		VPf("  Card %d: Initially %d wins\n", c+1, wins)
	}

	// then, scratch them and add earned extra cards
	for c := range wc {
		for next := c+1; next <= c+wc[c]; next++ {
			nc[next] += nc[c]	// add one card to wc following cards, nc times
		}
	}
	
	// count the total of cards
	for c, n := range nc {
		VPf("  Card %d: %d cards with %d wins\n", c+1, n, wc[c])
		ncards += n
	}
	return 
}

//////////// Common Parts code

//////////// Part1 functions

//////////// Part2 functions
