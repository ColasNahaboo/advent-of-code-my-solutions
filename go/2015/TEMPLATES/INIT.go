// Adventofcode YYYY, dayNN, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 input
// TEST: input
package main

import (
	"bufio"
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
	input, err := ioutil.ReadFile(infile)
	if err != nil {
		os.Exit(1)
	}

	var result int
	if *partOne {
		fmt.Println("Running Part1")
		result = Part1(input)
	} else {
		fmt.Println("Running Part2")
		result = Part2(input)
	}
	fmt.Println(result)
}

//func ReadInput(input *bufio.Scanner) []string {
//	var text []string
//    for input.Scan() {
//        text = append(text, scanner.Text())
//    }
//	return text
//}

func Part1(input *bufio.Scanner) int {
	for input.Scan() {
		line := input.Text()
	}
	return 0
}

func Part2(input *bufio.Scanner) int {
	for input.Scan() {
		line := input.Text()
	}
	return 0
}
