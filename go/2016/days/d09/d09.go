// Adventofcode 2016, d09, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 example 18
// TEST: -1 input 110346
// TEST: example 20
// TEST: input 10774309173
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

func part1(lines []string) int {
	d := decompress(lines[0])
	VP(d)
	return len(d)
}

//////////// Part 2
func part2(lines []string) int {
	l := estimate(lines[0])
	return l
}

//////////// Common Parts code


//////////// Part1 functions

func decompress(s string) (d string) {
	reMarker := regexp.MustCompile("[(]([[:digit:]]+)x([[:digit:]]+)[)]")
	end := 0
	for {
		marker := reMarker.FindStringSubmatchIndex(s)
		if marker == nil {
			break
		}
		d += s[0:marker[0]]		// string part before first marker
		datalen := atoi(s[marker[2]:marker[3]])
		repeats := atoi(s[marker[4]:marker[5]])
		end = marker[1]+datalen
		data := s[marker[1]:end]
		for i := 0; i < repeats; i++ {
			d += data
		}
		s = s[end:]
	}
	d += s
	return
}

//////////// Part2 functions

// the input is composed of sequences {marker}{data} interspread into a string
// decompression replaces (LxR){data} by R times the expansion of data
// we recurse in the "data" parts of markers, only returning the length of expansion
func estimate(s string) (l int) {
	reMarker := regexp.MustCompile("[(]([[:digit:]]+)x([[:digit:]]+)[)]")
	end := 0
	for {
		marker := reMarker.FindStringSubmatchIndex(s)
		if marker == nil {
			break
		}
		l += marker[0]		// count string part before first marker
		datalen := atoi(s[marker[2]:marker[3]])
		repeats := atoi(s[marker[4]:marker[5]])
		end = marker[1]+datalen
		data := s[marker[1]:end]
		// we replace the marker+data substring of length "end-marker[0]" by R*estimate(data)
		l += repeats * estimate(data)
		s = s[end:]
	}
	l += len(s)
	return
}
