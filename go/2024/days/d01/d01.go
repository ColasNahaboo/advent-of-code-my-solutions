// Adventofcode 2024, d01, in go. https://adventofcode.com/2024/day/01
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 11
// TEST: example 31
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"flag"
	"fmt"
	"regexp"
	"sort"
	// "golang.org/x/exp/slices"
)

var verbose, debug bool

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	debugFlag := flag.Bool("V", false, "debug: even more verbose")
	flag.Parse()
	verbose = *verboseFlag
	debug = *debugFlag
	var infile string
	if flag.NArg() > 0 {
		infile = flag.Arg(0)
	} else {
		infile = fileMatch("input,[[:alnum:]]*,[[:alnum:]]*.test")
	}
	lines := fileToLines(infile)
	if *partOne {
		VP("Running Part1")
		fmt.Println(part1(lines))
	} else {
		VP("Running Part2")
		fmt.Println(part2(lines))
	}
}

//////////// Part 1

func part1(lines []string) (dist int) {
	l, r := parse(lines)
	// sort both lists
	sort.Slice(l, func(i, j int) bool { return l[i] < l[j] })
	sort.Slice(r, func(i, j int) bool { return r[i] < r[j] })

	for i := range l {
		VPf("  == [%d ] = %d, %d\n", i, l[i], r[i])
		dist += intAbs(l[i] - r[i])
	}
	return
}

//////////// Part 2
func part2(lines []string) (simils int) {
	l, r := parse(lines)
	for i := range l {
		simil := similOf(l[i], r)
		VPf("  == [%d ] = %d : %d, %d\n", i, simil, l[i], r[i])
		simils += l[i] * simil
	}
	return
}

// how many times n appears in list l?
func similOf(n int, l []int) (simil int) {
	for i := range l {
		if n == l[i] {
			simil++
		}
	}
	return
}	

//////////// Common Parts code

func parse(lines []string) (l, r []int) {
	renums := regexp.MustCompile("^([[:digit:]]+)[[:space:]]+([[:digit:]]+)")
	for lineno, line := range lines {
		nums := renums.FindStringSubmatch(line)
		if nums == nil {
			panic(fmt.Sprintf("Syntax error line %d: %s\n", lineno+1, line))
		}
		l = append(l, atoi(nums[1]))
		r = append(r, atoi(nums[2]))
	}
	return
}

//////////// PrettyPrinting & Debugging functions

func DEBUG() {
	if ! debug { return }
	// ad hoc debug code here
}
