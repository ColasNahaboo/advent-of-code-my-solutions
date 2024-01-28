// Adventofcode 2017, d13, in go. https://adventofcode.com/2017/day/13
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 24
// TEST: example 10
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

// At time t, a scanner at depth d is at position: range * 2 - 2

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
		infile = fileMatch("input,[[:alnum:]]*,[[:alnum:]]*.test")
	}
	lines := fileToLines(infile)
	if *partOne {
		VP("Running Part1")
		fmt.Println(part1(lines))
	} else {
		VP("Running Part2")
		fmt.Println(part2(lines))
	}
}

//////////// Part 1

func part1(lines []string) (severity int) {
	ranges := parse(lines)
	for t := 0; t < len(ranges); t++ {
		if ranges[t] == 0 {		//  skip layers without scanners
			continue
		}
		if t % (ranges[t] * 2 - 2) == 0 { // scanner is at top
			VPf("  Hit scanner at %d of range %d\n", t, ranges[t])
			severity += t * ranges[t]
		}
	}
	return
}

//////////// Part 2
func part2(lines []string)  int {
	ranges := parse(lines)
TRY_DELAY:
	for delay := 0; true; delay++ {
		for t := 0; t < len(ranges); t++ {
			if ranges[t] == 0 {		//  skip layers without scanners
				continue
			}
			if (t + delay) % (ranges[t] * 2 - 2) == 0 { // scanner is at top
				VPf("  Delay %d: hit scanner at %d of range %d\n", delay, t, ranges[t])
				continue TRY_DELAY
			}
		}
		return delay
	}
	return -1
}

//////////// Common Parts code

func parse(lines []string) (ranges []int) {
	re := regexp.MustCompile("[[:digit:]]+")
	m := re.FindAllString(lines[len(lines)-1], -1)
	lranges := atoi(m[0]) + 1
	ranges = make([]int, lranges, lranges)
	for _, line := range lines {
		m = re.FindAllString(line, -1)
		ranges[atoi(m[0])] = atoi(m[1])
	}
	return
}

	
//////////// PrettyPrinting & Debugging functions
