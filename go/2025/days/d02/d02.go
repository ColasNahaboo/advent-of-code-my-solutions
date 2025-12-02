// Adventofcode 2025, d02, in go. https://adventofcode.com/2025/day/02
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 1227775554
// TEST: example 4174379265
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"fmt"
	"regexp"
	// "flag"
	// "slices"
)

// we use the naming IDR (ID Range), easier as range is a reserved keyword.
// All IDs are >0
type IDR [2]int			

//////////// Options parsing & exec parts

func main() {
	ExecOptions(2, NoXtraOpts, part1, part2)
}

//////////// Part 1

func part1(lines []string) (res int) {
	idrs := parse(lines)
	for _, idr := range idrs {
		for i := idr[0]; i <= idr[1]; i++ {
			res += invalidIdValueTwice(i)
		}
	}
	return 
}

//////////// Part 2

func part2(lines []string) (res int) {
	idrs := parse(lines)
	for _, idr := range idrs {
		for i := idr[0]; i <= idr[1]; i++ {
			res += invalidIdValueRepeated(i)
		}
	}
	return 
}

//////////// Common Parts code

// return 0 if ID valid, its integer value otherwise

func invalidIdValueTwice(id int) (res int) {
	if invalidTwiceSequence(id) {
		return id
	}
	return 0
}

func invalidTwiceSequence(id int) bool {
	ids := itoa(id)
	if len(ids) % 2 != 0 {
		return false
	}
	hl := len(ids) / 2
	if ids[hl:] == ids[0:hl] {
		return true
	}
	return false
}

func invalidIdValueRepeated(id int) (res int) {
	if invalidRepeatedSequence(id) {
		return id
	}
	return 0
}

func invalidRepeatedSequence(id int) bool {
	ids := itoa(id)
SUBSEQ:
	for ssl := range len(ids) / 2 { // for all possible subsequences lengths ssl
		ssl++
		if len(ids) % ssl != 0 {	// first, can we fit perfectly N ssl into s?
			continue SUBSEQ
		}
		for i := range len(ids) / ssl - 1{ // i-th subseq all duplicates of #0 ?
			i++							  // start with ss #1
			if ids[i*ssl:(i+1)*ssl] != ids[0:ssl] {
				continue SUBSEQ
			}
		}
		return true
	}
	return false
}

func parse(lines []string) (idrs []IDR) {
	reidr := regexp.MustCompile("([[:digit:]]+)-([[:digit:]]+)")
	for _, line := range lines {
		m := reidr.FindAllStringSubmatch(line, -1)
		for _, idrm := range m {
			idr := IDR{atoi(idrm[1]), atoi(idrm[2])}
			idrs = append(idrs, idr)
		}
	}
	return
}

//////////// PrettyPrinting & Debugging functions. See also the VPx functions.

func DEBUG() {
	if ! debug { return }
	fmt.Println("DEBUG!")
	// ad hoc debug code here
}
