// Adventofcode 2016, d07, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 input 110
// TEST: input 242
package main

import (
	"flag"
	"fmt"
	//"log"
	"regexp"
)

var verbose bool
var re = regexp.MustCompile("([a-z]+)([[]([a-z]+)[]]|$)")

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
	ntls := 0
	for _, line := range lines {
		tlsin := false
		tlsout := false
		for _, m := range re.FindAllStringSubmatch(line, -1) {
			if hasAbba(m[1]) {
				tlsin = true
			}
			if hasAbba(m[3]) {
				tlsout = true
				break
			}
		}
		if tlsin && ! tlsout {
			ntls++
		}
	}
	return ntls
}

//////////// Part 2
func part2(lines []string) int {
	nssl := 0
	for _, line := range lines {
		supernet := make(map[string]bool, 0)
		hypernet := make(map[string]bool, 0)
		for _, m := range re.FindAllStringSubmatch(line, -1) {
			addAba(supernet, m[1])
			addAba(hypernet, m[3])
		}
		for aba := range supernet {
			bab := aba[1:2] + aba[0:1] + aba[1:2]
			if hypernet[bab] {
				nssl ++
				break
			}
		}
	}
	return nssl
}

//////////// Common Parts code

func hasAbba(s string) bool {
	if len(s) < 4 {
		return false
	}
	for i := 0; i < len(s) - 3; i++ {
		if s[i] == s[i+3] && s[i+1] == s[i+2] && s[i] != s[i+1] {
			return true
		}
	}
	return false
}

func addAba(abas map[string]bool, s string) {
	if len(s) < 3 {
		return
	}
	for i := 0; i < len(s) - 2; i++ {
		if s[i] == s[i+2] && s[i] != s[i+1] {
			abas[s[i:i+3]] = true
		}
	}
	return
}

//////////// Part1 functions

//////////// Part2 functions
