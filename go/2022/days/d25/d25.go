// Adventofcode 2022, d25, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 example 2=-1=0
// TEST: -1 input 2=1-=02-21===-21=200
// TEST: example
// TEST: input
package main

import (
	"flag"
	"fmt"
	// "regexp"
)

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

	var result int
	if *partOne {
		VP("Running Part1")
		fmt.Println(part1(lines))
	} else {
		VP("Running Part2")
		result = part2(lines)
		fmt.Println(result)
	}
}

//////////// Part 1

func part1(lines []string) string {
	n := 0
	for _, line := range lines {
		n += snafu2dec(line)
	}
	return dec2snafu(n)
}

//////////// Part 2
func part2(lines []string) int {
	for _, line := range lines {
		fmt.Println(line)
	}
	return 0
}

//////////// Common Parts code

func snafu2dec(s string) (n int) {
	for i, place := len(s) - 1, 1; i >= 0; i-- {
		switch s[i] {
		case '2': n += place * 2
		case '1': n += place
		case '-': n += place * (-1)
		case '=': n += place * (-2)
		}
		place *= 5
	}
	return
}

func dec2snafu(n int) (s string) {
	place := 1
	for {
		d := (n / place) % 5
		switch d {
		case 0: s = "0" + s
		case 1: s = "1" + s
		case 2: s = "2" + s
		case 3: s = "=" + s; n += place * 5 // carry
		case 4: s = "-" + s; n += place * 5 // carry
		}
		place *= 5
		if n < place {
			return
		}
	}
}

//////////// Part1 functions

//////////// Part2 functions
