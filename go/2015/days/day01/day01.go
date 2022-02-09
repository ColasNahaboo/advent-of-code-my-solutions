// Adventofcode 2015, day01, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: input 1795
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	partOne := flag.Bool("1", false, "run part one code, instead of part 2 (default)")
	flag.Parse()
	infile := "input.txt"
	if flag.NArg() > 0 {
		infile = flag.Arg(0)
	}
	i, err := ioutil.ReadFile(infile)
	if err != nil {
		os.Exit(1)
	}
	var result int
	if *partOne {
		fmt.Println("Running Part1")
		result = Part1(string(i))
	} else {
		fmt.Println("Running Part2")
		result = Part2(string(i))
	}
	fmt.Println(result)
}

func Part1(input string) int {
	floor := 0
	for _, c := range input {
		switch {
		case c == '(':
			floor++
		case c == ')':
			floor--
		}
	}
	return floor
}

func Part2(input string) int {
	floor := 0
	position := 0
	for p, c := range input {
		switch {
		case c == '(':
			floor++
		case c == ')':
			floor--
		}
		if floor == -1 {
			position = p + 1
			break
		}
	}
	return position
}
