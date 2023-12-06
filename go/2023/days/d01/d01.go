// Adventofcode 2023, d01, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 example 142
// TEST: example2 281

package main

import (
	"flag"
	"fmt"
	"strings"
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

// for part1 we use a simple char index approach
func part1(lines []string) int {
	var i, high, low, sum, value int
	digits := "0123456789"
	for _, line := range lines {
		if i = strings.IndexAny(line, digits); i < 0 {
			VP("No match on line: " + line)
			continue
		}
		high = atoi(line[i:i+1])
		if i = strings.LastIndexAny(line, digits); i < 0 {
			continue
		}
		low = atoi(line[i:i+1])
		value = high * 10 + low
		VPf("  value = (%d, %d) = %d\n", high, low, value)
		sum += value
	}
	return sum
}

//////////// Part 2
// Part 2 is tricky: "eighthree" should be parsed as 83, but just doing a FindAll
// will only see {"8"}, not {"8", "3"}, as the t will be gobbled by the 8 and
// not available to match "three"
// Thus, for the righmost digit, we look in reverse by hand iteration

var digitAnyRe, digitRe *regexp.Regexp
var  digitNames = []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}

// for part2 we use the more general regexp approach
func part2(lines []string) int {
	digitAnyRe = regexp.MustCompile("[0-9]|one|two|three|four|five|six|seven|eight|nine")
	digitRe = regexp.MustCompile("[0-9]")
	var sum, value int
	var match string
	for _, line := range lines {
		if match = digitAnyRe.FindString(line); len(match) == 0 {
			VP("No match on line: " + line)
			continue
		}
		high := digitValue(match)
		low := 0
		l := len(line)
		for i := len(line) - 1; i >= 0; i-- {
			if match = digitAnyRe.FindString(line[i:l]); len(match) == 0 {
				continue
			}
			low = digitValue(match)
			break
		}
		value = high * 10 + low
		VPf("  value = (%d, %d) = %d\n", high, low, value)
		sum += value
	}
	return sum
}

func digitValue(s string) int {
	if digitRe.MatchString(s) {
		return atoi(s)
	}
	var i int
	if i = IndexOf(digitNames, s); i < 0 {
		panic("Name of digit unknown: " + s)
	}
	return i + 1
}

//////////// Common Parts code

//////////// Part1 functions

//////////// Part2 functions
