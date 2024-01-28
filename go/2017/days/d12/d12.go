// Adventofcode 2017, d12, in go. https://adventofcode.com/2017/day/12
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 6
// TEST: example 2
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"flag"
	"fmt"
	"regexp"
)

var verbose bool

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
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
	if *partOne {
		VP("Running Part1")
		fmt.Println(part1(lines))
	} else {
		VP("Running Part2")
		fmt.Println(part2(lines))
	}
}

//////////// Part 1

func part1(lines []string) int {
	pipes := parse(lines)
	groupsmap := make([]bool, len(pipes), len(pipes))
	n, _ := inGroup(0, pipes, groupsmap)
	return n
}

func inGroup(p int, pipes [][]int, groupsmap []bool) (int, []bool) {
	n := 0
	todo := append([]int{}, pipes[p]...) // FIFO queue (stack) of nodes to check
	for len(todo) > 0 {
		q := todo[len(todo)-1]	// Pop
		todo = todo[0:len(todo)-1]
		if groupsmap[q] {	// already processed, ignore
			continue
		}
		n++		 // count node q in group, mark it, and look at connected pipes
		groupsmap[q] = true
		todo = append(todo, pipes[q]...) // Push
	}
	return n, groupsmap
}
		

//////////// Part 2
func part2(lines []string) int {
	pipes := parse(lines)
	groupsmap := make([]bool, len(pipes), len(pipes))
	grouped := 0
	ngroups := 0
	ngroup := 0
	for grouped < len(pipes) {
		ngroup, groupsmap = inGroup(nextUnGrouped(groupsmap), pipes, groupsmap)
		grouped += ngroup
		ngroups++
	}
	return ngroups
}

func nextUnGrouped(groupsmap []bool) int {
	for i, grouped := range groupsmap{ 
		if ! grouped {
			return i
		}
	}
	return -1
}

//////////// Common Parts code

func parse(lines []string) (pipes [][]int) {
	re := regexp.MustCompile("[[:digit:]]+")
	for _, line := range lines {
		m := re.FindAllString(line, -1)
		pipe := []int{}
		for _, p := range m[1:] {
			pipe = append(pipe,atoi(p))
		}
		pipes = append(pipes, pipe)
	}
	return
}

	

//////////// PrettyPrinting & Debugging functions
