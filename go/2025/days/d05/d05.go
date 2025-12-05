// Adventofcode 2025, d05, in go. https://adventofcode.com/2025/day/05
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 3
// TEST: example 14
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"fmt"
	"regexp"
	// "flag"
	// "slices"
)

type Range [2]int				// inclusive, exclusive: [i, j[

//////////// Options parsing & exec parts

func main() {
	ExecOptions(2, NoXtraOpts, part1, part2)
}

//////////// Part 1

func part1(lines []string) (res int) {
	ranges, ids := parse(lines)
	for _, id := range(ids) {
		if isFresh(ranges, id) {
			res++
		}
	}
	return 
}

//////////// Part 2

func part2(lines []string) (res int) {
	ranges, _ := parse(lines)
	fresh := runion(ranges)		// we build a list of distinct ranges, ordered
	VP(fresh)
	return runionLen(fresh)		// and return the sum of its ranges lengths
}

func runionLen(ranges []Range) (l int) {
	for _, r := range ranges {
		l += r[1] - r[0]
	}
	return
}

func runion(ranges []Range) []Range {
	ru := []Range{ranges[0]}
	for _, r := range ranges[1:] {
		ru = rappend(ru, r)
	}
	return ru
}
	
func rappend(ru []Range, ra Range) (ru2 []Range) {
	// we scan all ru to find where to insert ra: just before a bigger one
	for i, r := range ru {
		if rless(ra, r) { // bigger distinct r found: insert ra before r
			ru2 = append(ru2, ra)
			ru2 = append(ru2, ru[i:]...) // copy rest
			return
		} else if rless(r, ra) { // not yet found a bigger, continue looking
			ru2 = append(ru2, r)
		} else {				//  merge ra with r and absorbable followers
			// we absorb what ra overlaps
			j := i
			maxid := ra[1]
			for j < len(ru) && ! rless(ra, ru[j]) {
				maxid = max(maxid, ru[j][1])
				j++
			}
			r2 := Range{min(ra[0], r[0]), maxid}
			ru2 = append(ru2, r2)
			if j < len(ru) {
				ru2 = append(ru2, ru[j:]...)
			}
			return
		}
	}
	ru2 = append(ru2, ra)	// append at end if we could not fit it in
	VPf("%v + %v = %v\n", ru, ra, ru2)
	return
}

// r1 is before r2, with a non-null gap
func rless(r1, r2 Range) bool {
	if r1[1] < r2[0] {
		return true
	}
	return false
}

//////////// Common Parts code

func isFresh(ranges []Range, id int) bool {
	for _, r := range ranges {
		if id >= r[0] && id < r[1] {
			return true
		}
	}
	return false
}

func parse(lines []string) (ranges []Range, ids []int) {
	rerange := regexp.MustCompile("^([[:digit:]]+)-([[:digit:]]+)$")
	renum := regexp.MustCompile("^[[:digit:]]+$")
	for _, line := range lines {
		if m := rerange.FindStringSubmatch(line); m != nil {
			i1 := atoi(m[1])
			i2 := atoi(m[2]) + 1
			if i2 <= i1 {
				panic("invalid range")
			}
			ranges = append(ranges, Range{i1, i2})
		} else if m := renum.FindString(line); m != "" {
			ids = append(ids, atoi(m))
		} else if line != "" {
			panic("input error")
		}
	}
	if len(ranges) == 0 {
		panic("no ranges in input")
	}
	return
}

//////////// PrettyPrinting & Debugging functions. See also the VPx functions.

func DEBUG() {
	if ! debug { return }
	fmt.Println("DEBUG!")
	// ad hoc debug code here
}

func VPrange(r Range) {
	if ! verbose {
		return
	}
	fmt.Print(range2string(r))
}

func range2string(r Range) string {
	return fmt.Sprintf("[%d %d[", r[0], r[1])
}
	
func ranges2string(rs []Range) (s string) {
	for _, r := range rs {
		s = s + " " + fmt.Sprintf("[%d %d[", r[0], r[1])
	}
	return
}
	
