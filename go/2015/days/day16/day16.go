// Adventofcode 2015, day16, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 input 213
// TEST: input 323
package main

import (
	"flag"
	"fmt"
	"regexp"
)

var verbose bool

func main() {
	partOne := flag.Bool("1", false, "run part one code, instead of part 2 (default)")
	verboseFlag := flag.Bool("v", false, "verbose: print routes")
	flag.Parse()
	verbose = *verboseFlag
	infile := "input.txt"
	if flag.NArg() > 0 {
		infile = flag.Arg(0)
	}
	lines := fileToLines(infile)

	mfcsam := map[string]int{
		"children":    3,
		"cats":        7,
		"samoyeds":    2,
		"pomeranians": 3,
		"akitas":      0,
		"vizslas":     0,
		"goldfish":    5,
		"trees":       3,
		"cars":        2,
		"perfumes":    1,
	}

	ranges := map[string]rune{
		"children":    '=',
		"cats":        '>',
		"samoyeds":    '=',
		"pomeranians": '<',
		"akitas":      '=',
		"vizslas":     '=',
		"goldfish":    '<',
		"trees":       '>',
		"cars":        '=',
		"perfumes":    '=',
	}

	var result int
	if *partOne {
		fmt.Println("Running Part1")
		result = Part1(mfcsam, lines)
	} else {
		fmt.Println("Running Part2")
		result = Part2(mfcsam, ranges, lines)
	}
	fmt.Println(result)
}

//////////// Part 1

func Part1(mfcsam map[string]int, lines []string) int {
	return matchSue(mfcsam, lines)
}

//////////// Part 2
func Part2(mfcsam map[string]int, ranges map[string]rune, lines []string) int {
	return matchSueRanges(mfcsam, ranges, lines)
}

//////////// Common Parts code

//////////// Part1 functions

func matchSue(mfcsam map[string]int, lines []string) int {
	re := regexp.MustCompile(`[[:space:]]([[:alpha:]]+):[[:space:]]*([[:digit:]]+)`)
	for i, line := range lines {
		VPf("Parsing line %v\n", i)
		keyvals := re.FindAllStringSubmatch(line, -1)
		if keyvals != nil {
			ok := true
			for i, keyval := range keyvals {
				if i != 0 { // skip keyval[0], the whole string
					if atoi(keyval[2]) != mfcsam[keyval[1]] {
						VPf("Bad %v for %v\n", keyval, line)
						ok = false
						break
					}
				}
			}
			if ok {
				fmt.Println(line)
				return i + 1
			}
		}
	}
	return -1
}

//////////// Part2 functions

func matchSueRanges(mfcsam map[string]int, ranges map[string]rune, lines []string) int {
	re := regexp.MustCompile(`[[:space:]]([[:alpha:]]+):[[:space:]]*([[:digit:]]+)`)
	for i, line := range lines {
		VPf("Parsing line %v\n", i)
		keyvals := re.FindAllStringSubmatch(line, -1)
		if keyvals != nil {
			ok := true
			for _, keyval := range keyvals {
				got := atoi(keyval[2])
				expected := mfcsam[keyval[1]]
				switch op := ranges[keyval[1]]; op {
				case '<':
					if got >= expected {
						VPf("Bad %v for %v\n", keyval, line)
						ok = false
						break
					}
				case '>':
					if got <= expected {
						VPf("Bad %v for %v\n", keyval, line)
						ok = false
						break
					}
				default:
					if got != expected {
						VPf("Bad %v for %v\n", keyval, line)
						ok = false
						break
					}
				}
			}
			if ok {
				fmt.Println(line)
				return i + 1
			}
		}
	}
	return -1
}
