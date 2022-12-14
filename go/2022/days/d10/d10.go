// Adventofcode 2022, d10, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 input
// TEST: input
package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"
)

var verbose bool
var reinstr = regexp.MustCompile("^(addx|noop)[[:space:]]*([-[:digit:]]*)")

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

func part1(lines []string) (result int) {
	// timeline of X values indexed by cycle number
	// X value at END of cycle number N is x[N-1]
	// we fill all of x first, then compute the signals
	x := xTimeline(lines, 220)
	VP(x)
	for _, i := range []int{20, 60, 100, 140, 180, 220} {
		result += i * x[i - 2]		// indexes start at 0 + start value is value of previous
	}
	return 
}

//////////// Part 2
func part2(lines []string) int {
	x := xTimeline(lines, 240)	// precompute register X values
	t := 0						// time: the cycle number. X at end of cycle is x[t]
	pt := 0						// previous t. X at cycle start is x[pt]
	crt := make([]rune, 240, 240)
	// draw on CRT
	for row := 0; row < 6; row++ {
		for col := 0; col < 40; col++ {
			if col >= x[pt] - 1 && col <= x[pt] + 1 { // pixel in sprite
				crt[t] = '#'
			} else {
				crt[t] = '.'
			}
			// DEBUG
			//fmt.Printf("Cycle %d, pixel %d, sprite [%d-%d]\n", t+1, col, x[t]-1, x[t]+1)
			//fmt.Println(string(crt[row*40:(row+1)*40]))
			pt = t
			t++
		}
	}
	VPcrt(crt)
	return 0
}

//////////// Common Parts code

// we ensure all values are filled at least until size.
func xTimeline(lines []string, size int) []int {
	x := make([]int, 0, size)			
	xv := 1						// cache of last x value
	for lineno, line := range lines {
		instr := reinstr.FindStringSubmatch(line)
		if instr == nil {
			log.Fatalf("Syntax error line %d: %s\n", lineno+1, line)
		}
		if instr[1] == "noop" {
			x = append(x, xv)	// x stays the same for a cycle
		} else {
			n := atoi(instr[2])
			x = append(x, xv)	// a noop for a cycle
			xv += n
			x = append(x, xv)	// then adds addx value n in 2nd cycle
		}
	}
	// fill to size, avoiding errors on x[n]
	for i := 0; i < (size - len(x)); i++ {
		x = append(x, xv)
	}
	return x
}

func VPcrt(crt []rune) {
	if verbose {
		for row := 0; row < 6; row++ {
			fmt.Println(string(crt[row*40:(row+1)*40]))
		}
	}
}
//////////// Part1 functions

//////////// Part2 functions
