// Adventofcode 2017, d06, in go. https://adventofcode.com/2017/day/06
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 5
// TEST: example 4
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

func part1(lines []string) int {
	banks := parse(lines)
	history := History{}
	steps := 1
	for {
		banks.Spread(banks.Fullest())
		if history.Has(banks) {
			return steps
		}
		history.Add(banks)
		steps++
	}
	return 0
}


type Banks []int
type History []Banks

func (banks Banks) Fullest() (fi int) {
	maxblocks := -1
	for b, blocks := range banks {
		if blocks > maxblocks {
			fi = b
			maxblocks = blocks
		}
	}
	return
}

func (banks Banks) Spread(i int) {
	blocks := banks[i]
	banks[i] = 0
	lb := len(banks)
	xtra := blocks - (blocks / lb) * lb
	for o := 0; o < lb; o++ {
		n := i + 1 + o			// start with bank after i
		n = n % lb				// wrap around
		banks[n] += blocks / lb	// distribute evenly blocks/lb
		if xtra > 0 {			// distribute remains as much as we can
			banks[n]++
			xtra--
		}
	}
}

func (banks Banks) Equal(b Banks) bool {
	for i, bb := range b {
		if banks[i] != bb {
				return false
		}
	}
	return true
}

func (banks Banks) Clone() Banks {
	clone := make(Banks, len(banks), len(banks))
	copy(clone, banks)
	return clone
}

func (h *History) Add(b Banks) {
	*h = append(*h, b.Clone())
}

func (h *History) Has(b Banks) bool {
BANKS:
	for _, hb := range *h {
		if ! b.Equal(hb) {
			continue BANKS
		}
		return true
	}
	return false
}

//////////// Part 2

func part2(lines []string) int {
	banks := parse(lines)
	history := History{}
	steps := 1
	for {						// first, find start of loop, as in part1
		banks.Spread(banks.Fullest())
		if history.Has(banks) {
			break
		}
		history.Add(banks)
		steps++
	}
	steps = 1
	old := banks.Clone()
	for {						// then, find first repetition
		banks.Spread(banks.Fullest())
		if banks.Equal(old) {
			return steps
		}
		steps++
	}
}

//////////// Common Parts code

func parse(lines []string) (banks Banks) {
	re := regexp.MustCompile("[[:digit:]]+")
	for _, word :=range re.FindAllString(lines[0], -1) {
		banks = append(banks, atoi(word))
	}
	return
}

//////////// PrettyPrinting & Debugging functions
