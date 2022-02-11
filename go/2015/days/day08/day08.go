// Adventofcode 2015, day08, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 input 1342
// TEST: input 2074
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
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
	disk := 0
	mem := 0
	for _, line := range lines {
		if len(line) > 1 { //  skip empty lines
			d, m := lengthsDec(line)
			disk += d
			mem += m
		}
	}
	return disk - mem
}

var hexre = regexp.MustCompile("^x[0-9a-fA-F][0-9a-fA-F]")

func lengthsDec(s string) (int, int) {
	m := 0
loop:
	for i := 1; i < len(s); i++ { // start at 1 to skip first "
		switch {
		case s[i] == '\\':
			switch {
			case hexre.MatchString(s[i+1:]):
				i += 3
			case s[i+1] == '"' || s[i+1] == '\\':
				i++
			}
		case s[i] == '"':
			break loop
		}
		m++
	}
	return len(s), m
}

//////////// Part 2
func Part2(lines []string) int {
	diff := 0
	for _, line := range lines {
		if len(line) > 1 { //  skip empty lines
			d := lengthsEnc(line)
			diff += d
		}
	}
	return diff
}

func lengthsEnc(s string) int {
	d := 2 // the " "
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' || s[i] == '"' {
			d++
		}
	}
	// fmt.Printf("  %v ==> %v\n", s, d)
	return d
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
