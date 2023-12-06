// Adventofcode 2023, d06, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 288
// TEST: example 71503
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

// For a race of duration D, keeping the button pressed for time T
// makes the boat goes to length L:
// speed = T; L = speed * (D - T); L = T * (D - T); L = -1 T2 + D T
// Wins are T * (D - T) > R, values above solutions of -1 T2 + D T - R = 0
// Solutions are (-D +/- sqrt(D2 - 4R) ) / -2 = (D +/- sqrt(D2 - 4R)) / 2
// E.g for example race 1, D=7, R=9: (7 +- sqrt(49 - 36)) / 2
// sqrt(13) = 3.6... so we retain 3 (rounded low) to remove and 4 (rounded high)
// to add to get integer solutions inside the allowed range: 2 to 5.

package main

import (
	"flag"
	"fmt"
	"regexp"
	"math"
)

var verbose bool

var	durations []int					// the (fixed) duration of each race
var	records []int					// its current distance record

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

func part1(lines []string) (score int) {
	parse1(lines)
	score = 1
	for i, _ := range durations {
		wins := winsNum(i)
		score *= wins
	}
	return
}

//////////// Part 2
func part2(lines []string) int {
	parse2(lines)
	return winsNum(0)
}

//////////// Common Parts code

var margin = 0.00000000001		// a tiny offset to avoid delta being an int
func winsNum(race int) int {
	// sqrt(D2 - 4R))
	delta := math.Sqrt(float64(durations[race] * durations[race] - 4 * records[race]))
	VPf("    Delta = %f\n", delta)
	from := int((float64(durations[race]) - delta + margin) / 2) + 1
	to := int((float64(durations[race]) + delta - margin) / 2)
	VPf("  Race[%d](%d, %d): %d = %d..%d\n", race, durations[race], records[race], to - from + 1, from, to)
	return to - from + 1
}

//////////// Part1 functions

func parse1(lines []string) {
	renum := regexp.MustCompile("[0-9]+")
	for _, n := range renum.FindAllString(lines[0], -1) {
		durations = append(durations, atoi(n))
	}
	for _, n := range renum.FindAllString(lines[1], -1) {
		records = append(records, atoi(n))
	}
}

//////////// Part2 functions

func parse2(lines []string) {
	renum := regexp.MustCompile("[0-9]+")
	var line string
	for _, n := range renum.FindAllString(lines[0], -1) {
		line = line + n
	}
	durations = append(durations, atoi(line))
	line = ""
	for _, n := range renum.FindAllString(lines[1], -1) {
		line = line + n
	}
	records = append(records, atoi(line))
}
