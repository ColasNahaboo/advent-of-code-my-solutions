// Adventofcode 2016, d05, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 example 18f47a30
// TEST: -1 input d4cd2ee1
// TEST: input f2c730e5
package main

import (
	"flag"
	"fmt"
	//"log"
	"crypto/md5"
	//"regexp"
)

var verbose bool

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	flag.Parse()
	verbose = *verboseFlag
	infile := "input.txt"
	if flag.NArg() > 0 {
		infile = flag.Arg(0)
	}
	lines := fileToLines(infile)

	var result string
	if *partOne {
		VP("Running Part1")
		result = part1(lines)
	} else {
		VP("Running Part2")
		result = part2(lines)
	}
	fmt.Println(result)
}

//////////// Part 1

func part1(lines []string) (pass string) {
	doorId := lines[0]
	salt := 0
	n := 8
	var cksum string
	for i := 0; i < n; i++ {
		for ;; {
			cksum = fmt.Sprintf("%x", md5.Sum([]byte(doorId + itoa(salt))))
			salt++
			if cksum[0:5] == "00000" {
				VP(cksum)
				pass += cksum[5:6]
				break
			}
		}
	}
	return
}

//////////// Part 2
func part2(lines []string) string {
	doorId := lines[0]
	salt := 0
	nkeys := 0
	pass := []byte("        ")
	var position int
	var cksum string
	for ;; {
		cksum = fmt.Sprintf("%x", md5.Sum([]byte(doorId + itoa(salt))))
		salt++
		if cksum[0:5] == "00000" {
			VP(cksum)
			switch posx := cksum[5]; {
			case posx <= '9': position = int(posx - '0')
			default: position = 10 + int(posx - 'a')
			}
			if position < 8 && pass[position] == byte(32) {
				pass[position] = cksum[6]
				nkeys ++
				if nkeys >= 8 {
					return string(pass)
				}
			}
		}
	}
}

//////////// Common Parts code

//////////// Part1 functions

//////////// Part2 functions
