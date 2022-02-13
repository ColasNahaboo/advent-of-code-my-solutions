// Adventofcode 2015, day10, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 input 492982
// TEST: input
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	// "regexp"

	"strconv"
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
	_, len := Game(lines[0], 40)
	return len
}

//////////// Part 2
// the Conway contant, see https://en.wikipedia.org/wiki/Look-and-say_sequence
// it allows us to determine the size needed to pre-allocate the new sequence
// for efficiency
const conway = float64(1.303577269034296)

func Part2(lines []string) int {
	_, len := Game(lines[0], 50)
	return len
}

//////////// Common Parts code

func Game(s string, iterations int) ([]int, int) {
	lseq := len(s)
	seq := make([]int, lseq)
	for i := 0; i < lseq; i++ {
		c, err := strconv.Atoi(string(s[i]))
		if err != nil {
			log.Fatal(err)
		}
		seq[i] = c
	}
	for i := 0; i < iterations; i++ {
		seq, lseq = LookAndSaySeq(seq, lseq)
	}
	return seq, lseq
}

func LookAndSaySeq(s []int, l int) ([]int, int) {
	newsize := int(float64(l)*conway*1.2 + 10) // add 10% safety margin + 10
	new := make([]int, newsize)
	n := 0
	o := -1
	newi := 0
	for i := 0; i < l; i++ {
		if s[i] == o {
			n++
		} else {
			if o >= 0 {
				new[newi] = n
				newi++
				new[newi] = o
				newi++
			}
			o = s[i]
			n = 1
		}
	}
	if n > 0 {
		new[newi] = n
		newi++
		new[newi] = o
		newi++
	}
	return new, newi
}

// wrapper for easier testing
func LookAndSay(s string) string {
	seq, len := Game(s, 1)
	new := ""
	for i := 0; i < len; i++ {
		new += strconv.Itoa(seq[i])
	}
	return new
}

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
