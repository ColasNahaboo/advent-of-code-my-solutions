// Adventofcode 2018, d05, in go. https://adventofcode.com/2018/day/05
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 10
// TEST: example 4
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"fmt"
	"unicode"					// for the upper/lower handling
	"container/list"			// we use a double-linked list
)

type Unit struct {
	r rune						// canonical (lower case) version of the rune
	up bool						// is the unit uppercase?
}

//////////// Options parsing & exec parts

func main() {
	ExecOptions(2, NoXtraOpts, part1, part2)
}

//////////// Part 1

func part1(lines []string) (res int) {
	p := parse(lines)
	VP(*p)
	React(p)
	return p.Len()
}

//////////// Part 2

func part2(lines []string) (res int) {
	runes := RunesOf(lines[0])
	minlen := maxint
	for r := range runes {
		p := parseExcept(lines, r)
		React(p)
		l := p.Len()
		VPf("  Removing rune %s ==> len %d\n", string(r), l)
		if l < minlen {
			minlen = l
		}
	}
	return minlen
}

func parseExcept(lines []string, excl rune) (p *list.List) {
	p = list.New()
	for _, r := range lines[0] {
		lr := unicode.ToLower(r)
		if lr == excl {
			continue
		}
		u := Unit{r: lr}
		if unicode.IsUpper(r) {
			u.up = true
		}
		p.PushBack(u)
	}
	return
}

func RunesOf(line string) (ro map[rune]bool) {
	ro = make(map[rune]bool)
	for _, r := range line {
		ro[unicode.ToLower(r)] = true
	}
	return
}	

//////////// Common Parts code

func React(p *list.List) {
	for e := p.Front(); e.Next() != nil; {
		u := e.Value.(Unit)		// forces Value to be considered as of type Unit
		v := e.Next().Value.(Unit)
		if u.r == v.r && u.up != v.up { // React!
			VPf("  React on %s\n", string(u.r))
			prev := e.Prev()
			p.Remove(e.Next())
			p.Remove(e)
			e = prev
			if e == nil {		// e was at start, use the new start of p
				e = p.Front()
			}
		} else {
			e = e.Next()
		}
	}
}

func parse(lines []string) (p *list.List) {
	p = list.New()
	for _, r := range lines[0] {
		u := Unit{r: unicode.ToLower(r)}
		if unicode.IsUpper(r) {
			u.up = true
		}
		p.PushBack(u)
	}
	return
}

//////////// PrettyPrinting & Debugging functions. See also the VPx functions.

func DEBUG() {
	if ! debug { return }
	fmt.Println("DEBUG!")
	// ad hoc debug code here
}
