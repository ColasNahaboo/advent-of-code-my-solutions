// Adventofcode 2016, d16, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 example 01100
// TEST: -1 input 10101001010100001
// TEST: example 10111110011110111
// TEST: input 10100001110101001
package main

import (
	"flag"
	"fmt"
	// "regexp"
)

// a data is []int, a slice of bits

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
	size := atoi(lines[0])
	init := lines[1]

	var result string
	if *partOne {
		VP("Running Part1")
		result = part1(size, init)
	} else {
		VP("Running Part2")
		result = part1(35651584, init)
	}
	fmt.Println(result)
}

//////////// Part 1

func part1(size int, init string) string {
	data := []byte(init)
	for i := 0; i < len(init); i++ {
		if init[i] == '1' {
			data[i] = '1'
		} else {
			data[i] = '0'
		}
	}
	data = fill(size, data)
	return checksum(data[:size]) // checksum only the data fitting on the disc
}

//////////// Part 2

//////////// Common Parts code

func fill(size int, a []byte) []byte {
	for len(a) < size {
		b := make([]byte, len(a))
		for i := range a {		// b = reverse a and bit-invert
			if a[len(a) - 1 - i] == '0' {
				b[i] = '1'
			} else {
				b[i] = '0'
			}
		}
		a = append(a, '0')
		a = append(a, b...)
	}
	return a
}

func checksum(data []byte) string {
	for {
		chk := data[:len(data)/2]
		VPf("  len = %d\n", len(chk))
		for i := 0; i < len(data); i+=2 {
			if data[i] == data[i+1] {
				chk[i/2] = '1'
			} else {
				chk[i/2] = '0'
			}
		}
		if len(chk) % 2 == 1 {
			return string(chk)
		}
		data = chk				// even ==> loop
	}
}
	

//////////// Part1 functions

//////////// Part2 functions
