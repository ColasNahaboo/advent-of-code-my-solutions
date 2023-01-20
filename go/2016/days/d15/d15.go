// Adventofcode 2016, d15, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 example 5
// TEST: -1 input 122318
// TEST: example 85
// TEST: input 3208583
package main

import (
	"flag"
	"fmt"
	// "regexp"
)

type Disc struct {
	id, size, offset int		  // number of positions and pos when capsule hits if launched at t=0
}
var discs []Disc // list of discs defined by their [initial pos, number of pos]

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
	parse(lines)

	var result int
	if *partOne {
		VP("Running Part1")
		result = part1()
	} else {
		VP("Running Part2")
		result = part2()
	}
	fmt.Println(result)
}

//////////// Part 1

func part1() int {
	t := 0
	for {
		for _, d := range discs {
			if (t + d.offset) % d.size != 0 {
				goto BOUNCE
			}
		}
		return t
	BOUNCE:
		t++
	}
}

//////////// Part 2
func part2() int {
	// add extra disc
	discnum := len(discs)+1
	size := 11
	discs = append(discs, Disc{id: discnum, size: size, offset: (discnum + 0) % size})
	return part1()
}

//////////// Common Parts code

func parse(lines []string) {
	var lineno, discnum, size, start int
	for _, line := range lines {
		if n, _ := fmt.Sscanf(line, "Disc #%d has %d positions; at time=0, it is at position %d.", &discnum, &size, &start); n == 3 {
			if discnum != lineno+1 {
				panic(fmt.Sprintf("Bad discnum %d at line %d: %s\n", discnum, lineno, line))
			}
			discs = append(discs, Disc{id: discnum, size: size, offset: (discnum + start) % size})
		} else {
			panic(fmt.Sprintf("Syntax error line %d: %s\n", lineno, line))
		}
		lineno++
	}
}
//////////// Part1 functions

//////////// Part2 functions
