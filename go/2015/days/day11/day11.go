// Adventofcode 2015, day11, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 input hepxxyzz
// TEST: input heqaabcc
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	// "regexp"
	// "strconv"
)

func main() {
	partOne := flag.Bool("1", false, "run part one code, instead of part 2 (default)")
	flag.Parse()
	infile := "input.txt"
	if flag.NArg() > 0 {
		infile = flag.Arg(0)
	}
	lines := FileToLines(infile)
	var result string
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

func Part1(lines []string) string {
	pass := StringToSeq(lines[0])
	for {
		pass = NextPass(pass)
		if PassIsValid(pass) {
			return SeqToString(pass)
		}
	}
}

//////////// Part 2
func Part2(lines []string) string {
	pass := StringToSeq(lines[0])
	for {
		pass = NextPass(pass)
		if PassIsValid(pass) {
			for {
				pass = NextPass(pass)
				if PassIsValid(pass) {
					return SeqToString(pass)
				}
			}
		}
	}
}

//////////// Common Parts code
const seqbase = int('a') // letters are from base to high
const seqhigh = int('z')
const idx_i = int('i') - seqbase
const idx_o = int('o') - seqbase
const idx_l = int('l') - seqbase

func NextPass(pass []int) []int {
	return SeqIncPos(pass, len(pass)-1)
}

// we increment "digits" but skip the forbidden letters "i o l"
func SeqIncPos(pass []int, pos int) []int {
	if pass[pos] >= seqhigh-seqbase { // carry 1 to higher rank pos
		if pos > 0 {
			pass[pos] = 0
			pass = SeqIncPos(pass, pos-1)
		} else {
			fmt.Println("Fatal error: sequence overflow", pass)
			os.Exit(1)
		}
	} else {
		// this works because i o l are not z and not consecutive
		c := pass[pos] + 1
		if c == idx_i || c == idx_o || c == idx_l {
			c++
		}
		pass[pos] = c
	}
	return pass
}

func PassIsValid(pass []int) bool {
	//fmt.Println("Testing validity of:", pass, SeqToString(pass)) // DDD
	l := len(pass)
	// must include one increasing straight of at least three letters
	ok := false
	for i := 0; i+2 < l; i++ {
		if pass[i+1] == pass[i]+1 && pass[i+2] == pass[i]+2 {
			ok = true
			break
		}
	}
	if !ok {
		return false
	}
	// i o l not present because NextPass/SeqIncPos skip them

	// must contain at least two different, non-overlapping pairs
	ok = false
	for i := 0; i+3 < l; i++ {
		if pass[i+1] == pass[i] && HasPair(pass, i+2, pass[i]) {
			ok = true
			break
		}
	}
	if !ok {
		return false
	}
	return true
}

func HasPair(pass []int, pos, excluded int) bool {
	for i := pos; i+1 < len(pass); i++ {
		if pass[i+1] == pass[i] && pass[i] != excluded {
			return true
		}
	}
	return false
}

// convert a string of letters to a sequence of numbers, diffs from base
func StringToSeq(s string) []int {
	lseq := len(s)
	seq := make([]int, lseq)
	for i := 0; i < lseq; i++ {
		seq[i] = int(s[i]) - seqbase
	}
	return seq
}

// convert a  sequence of numbers (diff from base), to a string
// we use a Builder to avoid multiple string re-allocations
func SeqToString(seq []int) string {
	lseq := len(seq)
	var sb strings.Builder
	for i := 0; i < lseq; i++ {
		sb.WriteRune(rune(seq[i] + seqbase))
	}
	return sb.String()
}

//////////// Generic code

// useful in tests to feed Part1 & Part2 with a simple string (with newlines)
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

// read the input file into a string array for feeding Parts
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

// simplified functions to not bother with error handling. Just abort.

func Atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal(err)
	}
	return i
}

// for completeness
func Itoa(i int) string {
	return strconv.Itoa(i)
}
