// Adventofcode 2023, d15, in go. https://adventofcode.com/2023/day/15
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 1320
// TEST: example 145
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
	steps := parse1(lines)
	for _, step := range steps {
		sum += Hash(step)
	}
	return
}

func parse1(lines []string) (steps []string) {
	re := regexp.MustCompile("[^,[:space:]]+")
	for _, line := range lines {
		steps = append(steps, re.FindAllString(line, -1)...)
	}
	return
}

//////////// Part 2

type Lens struct {
	label string
	isfocal bool
	lenfocal int
}

type Box []Lens					// a box is a sparse list of lenses
var NOLENS = Lens{"",false,0}	// placeholder in box holes

func part2(lines []string) (sum int) {
	steps := parse2(lines)
	boxes := make([]Box, 256, 256)
	for i := range boxes {
		boxes[i] = Box{}
	}
	for _, step := range steps {
		boxnum := Hash(step.label)
		if step.isfocal {
			if li := boxes[boxnum].indexLens(step); li != -1 { // replace
				boxes[boxnum][li] = step
			} else {			// append to end
				boxes[boxnum] = append(boxes[boxnum], step)
			}
		} else {				// remove lens if present
			if li := boxes[boxnum].indexLens(step); li != -1 {
				boxes[boxnum][li] = NOLENS
			}
		}
	}
	for boxnum, box := range boxes {
		slotnum := 1
		for _, lens := range box {
			if lens.label == "" { // șkip holes
				continue
			}
			power := (boxnum+1) * slotnum * lens.lenfocal
			sum += power
			slotnum++			// do not increment in holes
		}
	}
	return
}

func parse2(lines []string) (lenses []Lens) {
	re := regexp.MustCompile("([[:alpha:]]+)([-=])([[:digit:]]*)")
	for _, line := range lines {
		for _, m := range re.FindAllStringSubmatch(line, -1) {
			var lens Lens
			if m[2] == "=" { 	// set focal
				lens = Lens{m[1], true, atoi(m[3])}
			} else {
				lens = Lens{m[1], false, 0}
			}
			lenses = append(lenses, lens)
		}
	}
	return
}

func (box Box)indexLens(lens Lens) int {
	for i := 0; i < len(box); i++ {
		if box[i].label == lens.label {
			return i
		}
	}
	return -1
}

//////////// Common Parts code

func Hash(step string) (hash int) {
	for _, c := range step {
		hash += int(c)
		hash *= 17
		hash %= 256
	}
	return hash
}

//////////// PrettyPrinting & Debugging functions
