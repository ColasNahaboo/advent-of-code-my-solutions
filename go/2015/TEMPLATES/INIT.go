// Adventofcode YYYY, dayNN, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 input
// TEST: input
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
)

func main() {
	partOne := flag.Bool("1", false, "run part one code, instead of part 2 (default)")
	flag.Parse()
	infile := "input.txt"
	if flag.NArg() > 0 {
		infile = flag.Arg(0)
	}
	presents := ReadInput(infile)

	var result int
	if *partOne {
		fmt.Println("Running Part1")
		result = Part1(presents)
	} else {
		fmt.Println("Running Part2")
		result = Part2(presents)
	}
	fmt.Println(result)
}

func ReadInput(infile string) []Present {
	i, err := ioutil.ReadFile(infile)
	if err != nil {
		os.Exit(1)
	}
	presents := new([]Present)
	var present Present
	re := regexp.MustCompile("([[:digit:]]+)x([[:digit:]]+)x([[:digit:]]+)")
	presentsStrings := re.FindAllStringSubmatch(string(i), -1)
	for _, presentString := range presentsStrings {
		l, err := strconv.Atoi(presentString[1])
		w, err := strconv.Atoi(presentString[2])
		h, err := strconv.Atoi(presentString[3])
		if err != nil {
			os.Exit(2)
		}
		present = Present{l, w, h}
		*presents = append(*presents, present)
	}
	return *presents
}
