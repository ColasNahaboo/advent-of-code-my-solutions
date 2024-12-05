// Adventofcode 2024, d05, in go. https://adventofcode.com/2024/day/05
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 143
// TEST: example 123
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

// by looking at the input, pages numbers are between 11 and 99 inclusive
// as they are two-digits numbers without using the digit zero
// so we use a slice [0:100] to store indexes of pages in update, plus one
// as we use the 0 value to mean "not present"
// E.g: if update is 4,2,6 orders is 0,0,2,0,1,0,3,0,0...
// because update[0] == 4, orders[4] == 1, etc...

// for part 2, we just switch the pages that make rules fail two per two
// until the update is in order

package main

import (
	"flag"
	"fmt"
	"regexp"
	// "golang.org/x/exp/slices"
)

var verbose, debug bool
var nilrule = [2]int{}

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

func part1(lines []string) (sum int) {
	rules, updates, orders := parse(lines)
	for i, u := range updates {
		if ok, _ := UpdateIsValid(rules, u, orders[i]); ok {
			sum += UpdateMiddlePage(u)
		}
	}
	return
}

//////////// Part 2
func part2(lines []string) (sum int) {
	var r [2]int
	var ok bool
	rules, updates, orders := parse(lines)
	for i, u := range updates {
		if ok, r = UpdateIsValid(rules, u, orders[i]); ok {
			continue			// skip correct ones
		}
		o := orders[i]
		for {
			SwitchPages(r, u, o) // switches elts in place in u and o
			if ok, r = UpdateIsValid(rules, u, orders[i]); ok {
				sum += UpdateMiddlePage(u)
				break
			}
		}
	}
	return
}

//////////// Common Parts code

// returns true, or false and which rule failed
func UpdateIsValid(rules [][2]int, u []int, o []int) (bool, [2]int) {
	for _, r := range rules {
		if RuleFailed(r, u, o) {
			return false, r
		}
	}
	return true, nilrule
}

func RuleFailed(r [2]int, u []int, o []int) bool {
	if o[r[0]] > 0 && o[r[1]] > 0 && o[r[0]] > o[r[1]] {
		return true
	} else {
		return false
	}
}

func UpdateMiddlePage(u []int) int {
	if (len(u) % 2) == 0 {
    	panic(fmt.Sprintf("Update has even number of pages: %v\n", u))
    }
	return u[len(u) / 2]
}

func SwitchPages(r [2]int, u []int, o []int) {
	i0, i1 := o[r[0]] - 1, o[r[1]] - 1	// indexes of rule pages in u
	u[i0], u[i1] = u[i1], u[i0]			// switch pages positions
	o[r[0]], o[r[1]] = o[r[1]], o[r[0]]			// update their indexes
}	

func parse(lines []string) (rules [][2]int, updates [][]int, orders [][]int) {
	rerule := regexp.MustCompile("^([[:digit:]]+)[|]([[:digit:]]+)")
	renum := regexp.MustCompile("[[:digit:]]+")
	var lineno int
	var line string
	for lineno, line = range lines {
		if m := rerule.FindStringSubmatch(line); m != nil {
			rules = append(rules, [2]int{atoi(m[1]), atoi(m[2])})
		} else {
			break
		}
	}
	for lineno++; lineno < len(lines); lineno++ {
		m := renum.FindAllString(lines[lineno], -1)
		update := []int{}
		order := make([]int, 100) // we known pages are in 0:100 range
		for _, pagename := range m {
			page := atoi(pagename)
			update = append(update, page)
			// order is index of page in update plus one, to use 0 as not there
			order[page] = len(update)
		}
		updates = append(updates, update)
		orders = append(orders, order)
	}
	return
}

//////////// PrettyPrinting & Debugging functions

func DEBUG() {
	if ! debug { return }
	// ad hoc debug code here
}
