// Adventofcode 2016, d02, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 input 74921
// TEST: input A6B35
package main

import (
	"flag"
	"fmt"
	// "regexp"
)

var verbose bool

// the pad of part1: each key has a map to neighbors indexed by direction
// 1 2 3
// 4 5 6
// 7 8 9
var pad1 = []map[string]int{
	{},
	{"U":1, "D":4, "L":1, "R":2},
	{"U":2, "D":5, "L":1, "R":3},
	{"U":3, "D":6, "L":2, "R":3},
	{"U":1, "D":7, "L":4, "R":5},
	{"U":2, "D":8, "L":4, "R":6},
	{"U":3, "D":9, "L":5, "R":6},
	{"U":4, "D":7, "L":7, "R":8},
	{"U":5, "D":8, "L":7, "R":9},
	{"U":6, "D":9, "L":8, "R":9},
}

// the pad of part2:
//     1
//   2 3 4
// 5 6 7 8 9
//   A B C         10 11 12
//     D              13 
var pad2 = []map[string]int{
	{},
	{"D":3},
	{"D":6, "R":3},
	{"U":1, "D":7, "L":2, "R":4},
	{"D":8, "L":3},
	{"R":6},
	{"U":2, "D":10, "L":5, "R":7},
	{"U":3, "D":11, "L":6, "R":8},
	{"U":4, "D":12, "L":7, "R":9},
	{"L":8},
	{"U":6, "R":11},
	{"U":7, "D":13, "L":10, "R":12},
	{"U":8, "L":11},
	{"U":11},
}

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

	if *partOne {
		VP("Running Part1")
		fmt.Println(part1(lines))
	} else {
		VP("Running Part2")
		fmt.Println(part2(lines))
		
	}
}

//////////// Part 1

func part1(lines []string) (code int) {
	key := 5
	for _, line := range lines {
		for _, next := range []rune(line) {
			nextkey, ok := pad1[key][string(next)]
			if ok {
				key = nextkey
			}
			VPf("%s: key %d, code: %d\n", string(next), key, code)			
		}
		code = 10 * code + key
	}
	return code
}

//////////// Part 2
func part2(lines []string) (code string) {
	key := 5
	for _, line := range lines {
		for _, next := range []rune(line) {
			nextkey, ok := pad2[key][string(next)]
			if ok {
				key = nextkey
			}
		}
		code = code + xtoa(key)
		VPf("%s: code: %s, key:%d (%s)\n", line, code, key, xtoa(key))
	}
	return code
}

//////////// Common Parts code

//////////// Part1 functions

//////////// Part2 functions

func xtoa(i int) string {
	if i < 10 {
		return itoa(i)
	} else {
		return string([]byte{byte(int([]byte("A")[0]) + i - 10)})
	}
}
