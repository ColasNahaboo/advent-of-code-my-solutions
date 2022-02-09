// Adventofcode 2015, day02, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 input 1588178
// TEST: input 3783758
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
)

type Present struct {
	l int
	w int
	h int
}

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
	presents := make([]Present, 0)
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
		presents = append(presents, present)
	}
	return presents
}

func min(is ...int) int {
	min := is[0]
	for _, i := range is[1:] {
		if i < min {
			min = i
		}
	}
	return min
}

func Part1(presents []Present) int {
	surface := 0
	for _, p := range presents {
		area := 2*p.l*p.w + 2*p.w*p.h + 2*p.h*p.l
		flap := min(p.l*p.w, p.w*p.h, p.h*p.l)
		surface += area + flap
	}
	return surface
}

func Part2(presents []Present) int {
	ribbons := 0
	for _, p := range presents {
		volume := p.l * p.w * p.h
		minperim := min(p.l+p.w, p.w+p.h, p.h+p.l)
		ribbon := 2*minperim + volume
		ribbons += ribbon
	}
	return ribbons
}
