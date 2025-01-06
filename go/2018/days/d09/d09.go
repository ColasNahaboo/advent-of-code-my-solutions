// Adventofcode 2018, d09, in go. https://adventofcode.com/2018/day/09
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 32
// TEST: -1 example1 8317
// TEST: -1 example2 146373
// TEST: -1 example3 2764
// TEST: -1 example4 54718
// TEST: -1 example5 37305
// TEST: example
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"fmt"
	"regexp"
	"container/list"			// we use a double-linked list
	// "flag"
	"slices"
)

//////////// Options parsing & exec parts

func main() {
	ExecOptions(2, NoXtraOpts, part1, part2)
}

//////////// Part 1

func part1(lines []string) (res int) {
	nplayers, nmarbles := parse(lines)
	return Play(nplayers, nmarbles)
}

//////////// Part 2

func part2(lines []string) (res int) {
	nplayers, nmarbles := parse(lines)
	return Play(nplayers, nmarbles * 100)
}

//////////// Common Parts code

func Play(nplayers, nmarbles int) int {
	circle := list.New()
	score := make([]int, nplayers) // scores
	player := 0
	cur := circle.PushFront(0)	// first marble
	for marble := 1; marble <= nmarbles; marble++ {
		if marble % 23 != 0 {	// move + 1, insert after it
			cur = circle.InsertAfter(marble, CNext(circle, cur))
		} else {				// move 7 back, get marble, cur is next
			score[player] += marble
			for i := 0; i < 6; i++ {
				cur = CPrev(circle, cur)
			}
			newcur := cur
			score[player] += circle.Remove(CPrev(circle, cur)).(int)
			cur = newcur
		}
		player = (player + 1) % nplayers
		VPlist(circle, cur)
	}
	VPf("Players scores: %v\n", score)
	return slices.Max(score)
}

func parse(lines []string) (players, marbles int) {
	renum := regexp.MustCompile("[[:digit:]]+") 
	ns := atoil(renum.FindAllString(lines[0], 2))
	return ns[0], ns[1]
}

// circular list functions on the underlying double-linked list
func CNext(l *list.List, e *list.Element) *list.Element {
	if e.Next() != nil {
		return e.Next()
	}
	return l.Front()
}

func CPrev(l *list.List, e *list.Element) *list.Element {
	if e.Prev() != nil {
		return e.Prev()
	}
	return l.Back()
}

//////////// PrettyPrinting & Debugging functions. See also the VPx functions.

func DEBUG() {
	if ! debug { return }
	fmt.Println("DEBUG!")
	// ad hoc debug code here
}

func VPlist(l *list.List, c *list.Element) {
	if ! verbose { return }
	for e := l.Front(); e != nil; e = e.Next() {
		if e == c {
			fmt.Printf("(%2d)", e.Value.(int))
		} else {
			fmt.Printf(" %2d ", e.Value.(int))
		}
	}
	fmt.Println()
}
	
