// Adventofcode 2015, day03, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 input 2592
// TEST: input 2360
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
	i := ReadInput(infile)

	var result int
	if *partOne {
		fmt.Println("Running Part1")
		result = Part1(i)
	} else {
		fmt.Println("Running Part2")
		result = Part2(i)
	}
	fmt.Println(result)
}

func ReadInput(infile string) []byte {
	i, err := ioutil.ReadFile(infile)
	if err != nil {
		os.Exit(1)
	}
	return i
}

func Part1(i []byte) int {
	x := 0
	y := 0

	houses := map[string]int{"0,0": 1}
	var coords string
	for _, c := range i {
		switch {
		case c == '^':
			y++
		case c == 'v':
			y--
		case c == '>':
			x++
		case c == '<':
			x--
		}
		coords = fmt.Sprintf("%v,%v", x, y)
		houses[coords]++
	}

	return len(houses)
}

func Part2(i []byte) int {
	sx := 0
	sy := 0
	rx := 0
	ry := 0
	santa := true

	houses := map[string]int{"0,0": 2}
	var coords string
	for _, c := range i {
		if santa {
			sx, sy, coords = move(sx, sy, c)
			santa = false
		} else {
			rx, ry, coords = move(rx, ry, c)
			santa = true
		}
		houses[coords]++
	}

	return len(houses)
}

func move(x, y int, c byte) (int, int, string) {
	switch {
	case c == '^':
		y++
	case c == 'v':
		y--
	case c == '>':
		x++
	case c == '<':
		x--
	}
	coords := fmt.Sprintf("%v,%v", x, y)
	return x, y, coords
}
