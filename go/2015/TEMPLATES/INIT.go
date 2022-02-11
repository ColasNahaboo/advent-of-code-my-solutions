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
	"log"
	"os"
	"strings"
)

func main() {
	partOne := flag.Bool("1", false, "run part one code, instead of part 2 (default)")
	flag.Parse()
	infile := "input.txt"
	if flag.NArg() > 0 {
		infile = flag.Arg(0)
	}
	lines := FileToLines(infile)

	var result int
	if *partOne {
		fmt.Println("Running Part1")
		result = Part1(lines)
	} else {
		fmt.Println("Running Part2")
		result = Part2(lines)
	}
	fmt.Println(result)
}

//////////// Part 1

func Part1(lines []string) int {
	for _, line := range lines {
		fmt.Println(line)
	}
	return 0
}

//////////// Part 2
func Part2(lines []string) int {
	for _, line := range lines {
		fmt.Println(line)
	}
	return 0
}

//////////// Common Parts code

//////////// Generic code
func StringToLines(s string) (lines []string) {
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	err := scanner.Err()
	if err != nil {
		log.Fatal(err)
	}
	return
}

func FileToLines(filePath string) (lines []string) {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	// optionally, resize scanner's capacity for lines over 64K (65536)
	const maxCapacity = 1000000 // your required line length
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)
	// end optional
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	err = scanner.Err()
	if err != nil {
		log.Fatal(err)
	}

	return
}
