// Adventofcode 2023, d12, in go. https://adventofcode.com/2023/day/12
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 21
// TEST: example 525152
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"flag"
	"fmt"
	"regexp"
	"strings"
)

var verbose, single bool

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	partTwoAlt := flag.Bool("a", false, "run exercise part2, but alternate version")
	partTwoAlt2 := flag.Bool("b", false, "run exercise part2, but alternate version 2")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	singleFlag := flag.Bool("s", false, "part2: single fold:: do not unfold 5 times")
	argData := flag.String("e", "", "take input from the string argument")
	flag.Parse()
	verbose = *verboseFlag
	single = *singleFlag
	var infile string
	if flag.NArg() > 0 {
		infile = flag.Arg(0)
	} else {
		infile = fileMatch("input,[0-9]*,[0-9]*.test")
	}
	var lines []string
	if *argData == "" {
		lines = fileToLines(infile)
	} else {
		lines = []string{*argData}
	}

	var result int
	if *partOne {
		VP("Running Part1")
		result = part1(lines)
	} else if *partTwoAlt {
		VP("Running Part2a")
		result = part2a(lines, single)
	} else if *partTwoAlt2 {
		VP("Running Part2b")
		result = part2b(lines, single)
	} else {
		// VP("Running Part2")
		result = part2(lines, single)
	}
	fmt.Println(result)
}

//////////// Part 1

// for part1, we use brute force: generate all combinations, and test all of them
func part1(lines []string) (sum int) {
	for _, line := range lines {
		sum += arrangements(line)
	}
	return
}

// A combination of N boolean values is a number of length N in binary format
// reminder conversions of N
//   decimal -> binary: strconv.FormatInt(N, 2)
//   binary -> decimal: strconv.ParseInt(N, 2, 64)

var reline = regexp.MustCompile("^([.#?]+) ([,0-9]+)")
var renum = regexp.MustCompile("[0-9]+")
var respan = regexp.MustCompile("[#]+")

func arrangements(line string) (na int) {
	lineparts := reline.FindStringSubmatch(line)
	conditions := lineparts[1]  // the condition records themselves: . # ?
	spans := []int{}				// list of numbers of consecutive #
	for _, num := range renum.FindAllString(lineparts[2], -1) {
		spans = append(spans, atoi(num))
	}
	arr := []byte(conditions)	// the arrangement candidate to test
	unkpos := []int{}				// list of positions of the unknowns
	for i, b := range line {
		if b == '?' {
			unkpos = append(unkpos, i)
		}
	}
	// now test all possible . or # combinations for the N unknows: 2^N
	ncombs := intPower(2, len(unkpos))
	VPf("Testing %d combinations for %s\n", ncombs, line)
testCombination:
	for comb := 0; comb < ncombs; comb++ {
		fmtstr := fmt.Sprintf("%%0%db", len(unkpos)) // a N long string of 0 or 1
		binstr := fmt.Sprintf(fmtstr, comb) // 0='.', 1='#'
		for i, b := range binstr {
			if b == '1' {
				arr[unkpos[i]] = '#'
			} else {
				arr[unkpos[i]] = '.'
			}
		}
		VPf("    Testing %s\n", string(arr))
		brokenSpans := respan.FindAll(arr, -1)
		if len(brokenSpans) != len(spans) {
			continue
		}
		for i, l := range spans {
			if len(brokenSpans[i]) != l {
				continue testCombination
			}
		}
		VPf("         OK %s satisfies %s\n", string(arr), lineparts[2])
		na++				// found a possible arrangement
	}
	VPf("  ==> %d arrangements\n", na)
	return	
}

//////////// Part 2
// logic copied from https://gist.github.com/sanyi/96ccaf6d3c0a67536b4fe3e99bc53bb3

var maxspan int

func part2(lines []string, single bool) (sum int) {
	times := 5
	if single {
		times = 1
	}
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		sum += arrnumsAll(line, times)
	}
	return
}

func arrnumsAll(line string, folds int) int {
	lineparts := reline.FindStringSubmatch(line)
	conditions := lineparts[1]  // the condition records themselves: . # ?
	rec1 := []byte(conditions)	// the spring records: array of . # ?
	spans1 := []int{}				// list of numbers of consecutive #
	for _, num := range renum.FindAllString(lineparts[2], -1) {
		span := atoi(num)
		spans1 = append(spans1, span)
		if span > maxspan {
			maxspan = span
		}
	}
	maxspan++
	// fmt.Printf("Exploring the line %s, unfolded %d times\n", line, folds)
	recb := []byte{}
	spans := []int{}
	for i := 0; i < folds; i++ {
		if len(recb) > 0 {
			recb = append(recb, '?')
		}
		recb = append(recb, rec1...)
		spans = append(spans, spans1...)
	}
	rec := string(recb)
	skb := make([]byte, len(spans), len(spans))
	for i, s := range spans {
		skb[i] = '0' + byte(s)
	}
	return arrnums(rec, spans)	// recurse
}

// set an element in the sparse array a, with auto-reallocation
// virtually: arrangements(spanidx, span) ==> a[spanidx * maxspan + span]
// values must be positive or null, holes have value -1
func sparse_set(a []int, spanidx, span, v int) []int {
	i := spanidx * maxspan + span
	if i >= len(a) {
		for j := len(a); j <= i; j++ {
			a = append(a, -1)
		}			
	}
	a[i] = v
	return a
}
// get an element value or 0
func sparse_get(a []int, spanidx, span int) int {
	i := spanidx * maxspan + span
	if i >= len(a) || a[i] == -1 {
		return 0
	} else {
		return a[i]
	}
}
// increment an element value by incr
func sparse_inc(a []int, spanidx, span, incr int) []int {
	v := sparse_get(a, spanidx, span)
	newa := sparse_set(a, spanidx, span, v + incr)
	// return sparse_set(a, spanidx, span, v + incr)
	for _, n := range newa {
		if n != -1 {
			goto OK
		}
	}
	panic("empty sparse!")
OK:
	return newa
}
// pretty print, imitates python print of a defaultdict
func VPsparse(a []int) {
	if !verbose {
		return
	}
	sep := ""
	fmt.Printf("defaultdict(<class 'int'>, {")
	for i, v := range a {
		if v == -1 {
			continue
		}
		fmt.Printf("%s(%d, %d): %d", sep, i/maxspan, i%maxspan, v)
		sep = ", "
	}
	fmt.Printf("})\n")
}
	
	
func arrnums(rec string, groups []int) int {
	//VPf("  ArrNums: %q %v\n", rec, groups)
	group_count := len(groups)
	arangements := []int{1}	// arangements(0,0) = 1

    for _, c := range rec {
        new_arangements := []int{}
        for idx, count := range arangements {
			if count < 0 {
				continue
			}
			group_index := idx / maxspan
			current_group_size := idx % maxspan
			valid := true
            if c == '#' {
                current_group_size += 1
                valid = is_valid(groups, group_index, current_group_size, false)
            } else if c == '?' {
                // option 1: what if it is a dot?
                if current_group_size > 0 {
					// check if is a valid termination (strict)
                    if is_valid(groups, group_index, current_group_size, true) {
						new_arangements = sparse_inc(new_arangements, group_index + 1, 0, count)
					}
				} else {
					new_arangements = sparse_inc(new_arangements, group_index, current_group_size, count)
				}
                // option 2: what if it is a hash?
				if is_valid(groups, group_index, current_group_size + 1, false) {
					new_arangements = sparse_inc(new_arangements, group_index, current_group_size + 1, count)
				}
				continue
            } else if current_group_size > 0 { // c == '.'
				// check if is a valid termination (strict)
				valid = is_valid(groups, group_index, current_group_size, true)
				current_group_size = 0
				group_index += 1
			}
            if valid {
				new_arangements = sparse_inc(new_arangements, group_index, current_group_size, count)
			}
		}
		arangements = new_arangements
	}
    c := 0
	for idx, count := range arangements {
		if count < 0 {
			continue
		}
		group_index := idx / maxspan
		current_group_size := idx % maxspan
        if current_group_size > 0 {
			if group_index >= len(groups) {
				continue
			}
			if current_group_size != groups[group_index] {
				continue
			}
            group_index += 1
		}
        if group_index == group_count {
            c += count
		}
	}
    return c
}

func is_valid(groups []int, group_index, current_group_size int, strict bool) bool {
	if group_index >= len(groups) {
		return false
	}
	if strict {
		return current_group_size == groups[group_index]
	} else {
		return current_group_size <= groups[group_index]
	}
}

//////////// Part 2 ALTERNATE 2 (part2a with cache). But still too slow for input
// for part2, brute force wont scale. We recurse and prune as soon as possible

func part2b(lines []string, single bool) (sum int) {
	times := 5
	if single {
		times = 1
	}
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		sum += countAll(line, times)
	}
	return
}

const SUCCESS = true
const FAIL = false
const NONE = 0
const OK = 1

func countAll(line string, folds int) int {
	lineparts := reline.FindStringSubmatch(line)
	conditions := lineparts[1]  // the condition records themselves: . # ?
	rec1 := []byte(conditions)	// the spring records: array of . # ?
	spans1 := []int{}				// list of numbers of consecutive #
	for _, num := range renum.FindAllString(lineparts[2], -1) {
		spans1 = append(spans1, atoi(num))
	}
	fmt.Printf("Exploring the line %s, unfolded %d times\n", line, folds)
	recb := []byte{}
	spans := []int{}
	for i := 0; i < folds; i++ {
		if len(recb) > 0 {
			recb = append(recb, '?')
		}
		recb = append(recb, rec1...)
		spans = append(spans, spans1...)
	}
	rec := string(recb)
	skb := make([]byte, len(spans), len(spans))
	for i, s := range spans {
		skb[i] = '0' + byte(s)
	}
	return count(rec, spans)	// recurse
}

// we try to fit the first span only, and recurse on the rest
func count(rec string, spans []int) int {
	VPf("  Count: %q %v\n", rec, spans)
	if len(rec) == 0 {
		if len(spans) == 0 {
			return OK
		} else {
			return NONE
		}
	} else if len(spans) == 0 {
		// it was the last span: rest of rec must not contain a span-starting #
		if strings.IndexByte(rec, '#') == -1 {
			return OK
		} else {
			return NONE
		}
	}
	if SumOfInts(spans) + len(spans) - 1 > len(rec) {
		// no need to cache obvious fails
		return NONE				// cannot fit spans and separators in rec
	}
	// from now on, uncached case
	if rec[0] == '#' {
		// we can cache only the sub-matching from a leading # to a single span
		span := spans[0]			// the span we try to fit
		// then all the chars of the next span must fit from here (no '.')
		if span > len(rec) {
			return NONE			// cannot physically fit
		}
		// any "hard" '.' prevents placement of span
		for i := 0; i < span ; i++ {
			if rec[i] == '.' {
				return NONE
			}
		}
		// if span fits it must not be adjacent to a following #
		if span < len(rec) && rec[span] == '#' {
			return NONE
		}
		// at this point we are sure that span fits in rec by itself
		if len(rec) > span {
			VPf("    %d span fits, recurse on %q %v\n", span,rec[span+1:],spans[1:])
			return count(rec[span+1:], spans[1:]) // recurse after span + sep
		} else {
			VPf("    %d span fits at end\n", span)
			return count(rec[span:], spans[1:]) // recurse after span + sep
		}
	} else if rec[0]  == '.' {			// skip leadings '.' to end, # or ?
		i := 0
		for ; i < len(rec) && rec[i] == '.'; i++ {
		}
		VPf("    skipping %d dots, recurse on %q %v\n", i, rec[i:], spans)
		return count(rec[i:], spans)
	} else {  // '?', exanime both options. we can skip the leading '.' in one.
		// duplicate rec for '#' branch, as we must not modify the underlying rec
		VPf("    ? as #, recurse on %q %v\n", "#" + rec[1:], spans)
		c1 := count("#" + rec[1:], spans)
		VPf("    ? as ., recurse on %q %v\n", rec[1:], spans)
		c2 := count(rec[1:], spans)
		return c1 + c2
	}
}

// Sum of ints in a slice of ints
func SumOfInts(s []int) (sum int) {
	for _, i := range s {
		sum += i
	}
	return
}


//////////// Part 2 ALTERNATE (complex, faster than par1t method but still slow)
// for part2, brute force wont scale.

func part2a(lines []string, single bool) (sum int) {
	times := 5
	if single {
		times = 1
	}
	for _, line := range lines {
		sum += exploreAll(line, times)
	}
	return
}

// keep the context (line data) in globals for simplicity
// these stay constant during the processing of each line
var arr []byte					// the .#? arrangement of the processed line
var spans []int					// the #-spans lengths
var maxpos []int				// the max pos to fit the rest of spans 

func exploreAll(line string, folds int) int {
	lineparts := reline.FindStringSubmatch(line)
	conditions := lineparts[1]  // the condition records themselves: . # ?
	arr1 := []byte(conditions)	// the arrangement: array of . # ?
	spans1 := []int{}				// list of numbers of consecutive #
	for _, num := range renum.FindAllString(lineparts[2], -1) {
		spans1 = append(spans1, atoi(num))
	}
	fmt.Printf("Exploring the line %s, unfolded %d times\n", line, folds)
	arr = []byte{}
	spans = []int{}
	for i := 0; i < folds; i++ {
		if len(arr) > 0 {
			arr = append(arr, '?')
		}
		arr = append(arr, arr1...)
		spans = append(spans, spans1...)
	}
	maxpos = make([]int, len(spans), len(spans))
	// if S spans remain to fill, p must be before this position maxpos[S]
	// TODO: count only # and ? positions?
	mp := len(arr) + 1
	for i := len(spans)-1; i >= 0; i-- {
		mp += spans[i]	+ 1// room for N # + one .
		maxpos[i] = mp
	}
	return explore(0, 0, 0, arr[0])
}

// explore all combinations with first "pos" spring states already determined
// pos is the position of the first remaining '?' in arr
// inspan is the number of # in the current span (0 if not in a # span already)
// span is the index of the next span to match
// c is the supposed value of arr[p] at the start of the exploration,
// so we can override it in exploration of options

func explore(p, inspan, span int, c byte) int {
	VPf("  exploring from char %q @%d\n", string(c), p)
	// go to next ?, keeping track of #
	for {
		VPf("    loop: %q @%d, inspan=%d, span=%d\n", string(c), p, inspan, span)
		if c == '.' {
			if inspan > 0 { // close current span
				if inspan != spans[span] { // its length is less than expected
					VPf("    FAILS, span#%d has len %d, LESS than %d\n", span, inspan, spans[span])
					return 0
				}
				inspan = 0		// prepare for next span, but not instancied yet
				span++			// so it can can be >=len(spans) at this point
			}
			p++
		} else if c == '#' {
			if inspan > 0 { // extend current span length
				inspan++
				if inspan > spans[span] { // its length is more than expected
					VPf("    FAILS, span#%d has len %d, MORE than %d\n", span, inspan, spans[span])
					return 0
				}
			} else {		// start the new span
				if span >= len(spans) { // but we would exceed spans count
					VPf("    FAILS, not enough spans: %d, need %d\n", span, len(spans))
					return 0
				}
				inspan = 1
			}
			p++
		} else {				// ?
			// we stay at p, but explore the 2 possible values of c: . #
			return explore(p, inspan, span, '#') + explore(p, inspan, span, '.')
		}
		if p >= len(arr) {		// AT END, examine if current span matches
			if inspan > 0 {		// end current span
				if inspan != spans[span] { // its length is not expected
					VPf("    FAILS at end, span#%d has len %d instead of %d\n", span, inspan, spans[span])
					return 0
				}
				span++
			}
			if span != len(spans) { // we found less spans than expected
				VPf("    FAILS at end, found only %d spans instead of %d\n", span, len(spans))
				return 0
			}
			VPf("SUCCESS!\n")
			return 1			// this branch satisfies!
		}
		// check we have room for the remaining spans
		if span < len(spans) && p >= maxpos[span] + inspan {
			VPf("    FAILS, %s@%d not enough space left %d to fit %d spans, need %d\n", string(arr[p]), p, len(arr) - p, len(spans) - span, len(arr) - maxpos[span])
			VPf("    maxpos = %v\n", maxpos)
			return 0
		}
		c = arr[p]
	}
}

//////////// Common Parts code

//////////// PrettyPrinting & Debugging functions
