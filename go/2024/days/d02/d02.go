// Adventofcode 2024, d02, in go. https://adventofcode.com/2024/day/02
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 2
// TEST: example 4
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

// for testing the record with one less level, instead of creating a new record
// we use a template ReportLess1 on the record masking one level
// and provide an iterator on it

package main

import (
	"flag"
	"fmt"
	"regexp"
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

func part1(lines []string) (safe int) {
	reports := parse(lines)
NEXT_REPORT:
	for n, r := range reports {
		dir := dirOf(r[0], r[1])
		for i := 0; i < len(r) - 1; i++ {
			if ! sameDir(dir, r[i], r[i+1]) {
				continue NEXT_REPORT
			}
			if ! differIn(1, 3, r[i], r[i+1]) {
				continue NEXT_REPORT
			}
		}
		VPf("  == [%d] %v is safe\n", n+1, r)
		safe++
	}
	return
}


//////////// Part 2
func part2(lines []string) (safe int) {
	reports := parse(lines)
NEXT_REPORT:
	for n, r := range reports {
		// check record, and all the records less 1 elt
		for x := -1; x < len(r); x++ { // x == -1 for the whole record
			r1 := ReportLess1{r, x, 0}
			if r1.isSafe() {
				VPf("  == [%d] %v is safe\n", n+1, r1)
				safe++
				continue NEXT_REPORT
			}
		}
	}
	return
}

// ReportLess1 is a report with one element removed at index x
// no element is removed for x == -1

type ReportLess1 struct {
	r []int
	x int
	p int
}

// iterator on r: return elements in order one per one, skips x, -1 at end
func (r *ReportLess1) next() (e int) {
	if r.p == r.x {
		r.p++
	}
	if r.p < len(r.r) {
		e = r.r[r.p]
	} else {
		e = -1
	}
	r.p++
	return
}		

func (r *ReportLess1) isSafe() bool {
	e0 := r.next()
	e1 := r.next()
	dir := dirOf(e0, e1)
	VPf("  == Check: %v\n", r)
	for {
		if ! sameDir(dir, e0, e1) {
			return false
		}
		if ! differIn(1, 3, e0, e1) {
			return false
		}
		e0 = e1
		e1 = r.next()
		if e1 == -1 {
			VPf("  ==      %v is safe\n", r)
			return true
		}
	}
}

//////////// Common Parts code

func parse(lines []string) (reports [][]int) {
	renum := regexp.MustCompile("[[:digit:]]+")
	for _, line := range lines {
		reportString := renum.FindAllString(line, -1)
		report := []int{}
		for _, levelString := range reportString {
			report = append(report, atoi(levelString))
		}
		reports = append(reports, report)
	}
	return
}

func dirOf(i, j int) int {
	if j > i {
		return 1
	} else if j < i {
		return -1
	}
	return 0
}

func sameDir(d, i, j int) bool {
	return d == dirOf(i, j)
}

func differIn(f, t, i, j int) bool {
	d := intAbs(i - j)
	return (d >= f) && (d <= t)
}

//////////// PrettyPrinting & Debugging functions

func DEBUG() {
	if ! debug { return }
	// ad hoc debug code here
}
