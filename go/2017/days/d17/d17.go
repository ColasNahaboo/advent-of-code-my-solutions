// Adventofcode 2017, d17, in go. https://adventofcode.com/2017/day/17
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 638
// TEST: example
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

// part3 is a brute force implementation of part2, but too slow. Takes 90s
// I kept it as reference

package main

import (
	"flag"
	"fmt"
)

var verbose bool

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	partThree := flag.Bool("3", false, "run exercise part3, (default: part2)")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	flag.Parse()
	verbose = *verboseFlag
	var infile string
	if flag.NArg() > 0 {
		infile = flag.Arg(0)
	} else {
		infile = fileMatch("input,[[:alnum:]]*,[[:alnum:]]*.test")
	}
	lines := fileToLines(infile)
	if *partThree {
		VP("Running Part3")
		fmt.Println(part3(lines))
	} else if *partOne {
		VP("Running Part1")
		fmt.Println(part1(lines))
	} else {
		VP("Running Part2")
		fmt.Println(part2(lines))
	}
}

//////////// Part 1

func part1(lines []string) int {
	step := parse(lines)
	buf := LinkedList[int]{}
	buf.Push(0)
	pos := buf.head
	for i := 1; i <= 2017; i++ {
		pos = insert(step, i, pos, &buf)
		VPbuf(i, pos, &buf)
	}
	if pos == buf.tail {
		return buf.head.val
	}
	return pos.next.val
}

func insert(step, val int, pos *LinkedCell[int], buf *LinkedList[int]) *LinkedCell[int] {
	actualSteps := step % val
	for i := 0; i < actualSteps; i++ {
		if pos == buf.tail {
			pos = buf.head
		} else {
			pos = pos.next
		}
	}
	cell := LinkedCell[int]{val: val, next: pos.next}
	pos.next = &cell
	if pos == buf.tail {
		buf.tail = &cell
	}
	if buf.head.val !=0 {panic("Non-zero at first pos for " + itoa(val))}
	return &cell
}

//////////// Part 2

func part2(lines []string) int {
	step := parse(lines)
	pos := 0
	for i := 1; i <= 50000000; i++ {
		pos = insert2(step, i, pos)
	}
	return v1
}

var v1 int						// value at pos 1, v0 is always 0

func insert2(step, val, pos int) int {
	np := (pos + step) % val
	if np == 0 {
		v1 = val
	}
	return np+1
}

//////////// Part 3

// for reference, this is the brute force approach for part2 similar to part1

func part3(lines []string) int {
	step := parse(lines)
	buf := LinkedList[int]{}
	buf.Push(0)
	pos := buf.head
	for i := 1; i <= 50000000; i++ {
		pos = insert(step, i, pos, &buf)
		VPf("  %d\n", buf.head.next.val)
	}
	return buf.head.next.val
}

//////////// Common Parts code

func parse(lines []string) int {
	return atoi(lines[0])
}

//////////// PrettyPrinting & Debugging functions

func VPbuf(val int, pos *LinkedCell[int], buf *LinkedList[int]) {
	if ! verbose { return }
	c := buf.head
	for {
		if c == pos {
			fmt.Printf(" (%d)", c.val)
		} else {
			fmt.Printf("  %d ", c.val)
		}
		if c == buf.tail {
			break
		}
		c = c.next
	}
	fmt.Println()
}
