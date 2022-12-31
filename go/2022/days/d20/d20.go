// Adventofcode 2022, d20, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 example 3
// TEST: -1 input 9866
// TEST: example 1623178306
// TEST: input 12374299815791
package main

import (
	"flag"
	"fmt"
	// "log"
	// "regexp"
)

// a double linked list to implement a ring buffer
type Link struct {
	id int						// its original position (unique)
	value int					// the value of the number, not unique
	prev, next *Link
}

var sequence [](*Link)
var links *Link
var size int
var verbose bool
const key = 811589153

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	flag.Parse()
	verbose = *verboseFlag
	infile := "input.txt"
	if flag.NArg() > 0 {
		infile = flag.Arg(0)
	}
	lines := fileToLines(infile)
	parse(lines)
	VPf("Len %2d:", size); VPseq()
	
	var result int
	if *partOne {
		VP("Running Part1")
		result = part1()
	} else {
		VP("Running Part2")
		result = part2()
	}
	fmt.Println(result)
}

//////////// Part 1

func part1() int {
	for _, linkp := range sequence {
		VPf("[%d] %2d:", linkp.id, linkp.value)
		move(linkp, 1)
		VPseq()
	}
	return score(1)
}

//////////// Part 2
func part2() int {
	for i := 0; i < 10; i++ {
		VPf("[%d]:", i)
		for _, linkp := range sequence {
			move(linkp, key)
		}
		VPseq()
	}
	return score(key)
}

//////////// Common Parts code

func parse(lines []string) {
	start := Link{}
	links = &start
	l := links
	for id, line := range lines {
		//  fill link
		n := atoi(line)
		l.id = id
		l.value = n
		sequence = append(sequence, l)
		size++
		//  append new empty link
		nl := Link{prev: l}
		l.next = &nl
		nl.prev = l
		l = &nl
	}
	// close the double-linked list. l is empty, cut it.
	l.prev.next = links
	links.prev = l.prev
}

// move a link by its value
func move(l *Link, k int) {
	o := l
	// tricky: the size of the list is size-1, since we remove (temporarily) o
	steps := (l.value * k) % (size - 1)
	if steps > 0 {
		for i := 0; i < steps; i++ {
			l = l.next
		}
		o.prev.next = o.next	// remove from old pos
		o.next.prev = o.prev
		o.next = l.next			// insert right of new pos
		l.next.prev = o	
		o.prev = l
		l.next = o
	} else if steps < 0 {
		for i := 0; i < -steps; i++ {
			l = l.prev
		}
		o.prev.next = o.next	// remove from old pos
		o.next.prev = o.prev
		o.prev = l.prev			// insert left of new pos
		l.prev.next = o
		o.next = l
		l.prev = o
	}
}

func findLinkOf(n int) *Link {
	l := links
	for l.value != n {
		l = l.next
	}
	return l
}

func VPseq() {
	if verbose {
		l := links
		for {
			fmt.Printf(" %2d", l.value)
			l = l.next
			if l == links {
			fmt.Println()
				return
			}
		}
	}
}

func score(k int) int {
	l := findLinkOf(0)
	var sum int
	for t := 0; t < 3; t++ {
		for i := 0; i < 1000; i++ {
			l = l.next
		}
		fmt.Printf("+ %d ", l.value * k)
		sum += l.value * k
	}
	fmt.Println()
	return sum
}

//////////// Part1 functions

//////////// Part2 functions
