// Adventofcode 2017, d02, in go. https://adventofcode.com/2017/day/02
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 18
// TEST: example2 9
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"flag"
	"fmt"
	"regexp"
)

var verbose bool

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	flag.Parse()
	verbose = *verboseFlag
	var infile string
	if flag.NArg() > 0 {
		infile = flag.Arg(0)
	} else {
		infile = fileMatch("input,[0-9]*,[0-9]*.test")
	}
	lines := fileToLines(infile)

	var result int
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

func part1(lines []string) (sum int) {
	re := regexp.MustCompile("[0-9]+")
	for _, line := range lines {
		m := re.FindAllString(line, -1)
		min := atoi(m[0])
		max := atoi(m[0])
		for _, s := range m[1:] {
			n := atoi(s)
			if n > max {
				max = n
			} else if n < min {
				min = n
			}
		}
		sum += max - min
	}
	return
}

//////////// Part 2
func part2(lines []string) (sum int) {
	re := regexp.MustCompile("[0-9]+")
	for lineno, line := range lines {
		m := re.FindAllString(line, -1)
		n := make([]int, len(m), len(m))
		for i, s := range m {
			n[i] = atoi(s)
		}
		VPf(" [%d]: %v\n", lineno, n)
		for i := 0; i < len(n); i++ {
			for j := 0; j < len(n); j++ {
				if i != j && n[j] > n[i] && n[j] % n[i] == 0 {
					VPf("  %d multiple of %d\n", n[j], n[i])
					sum += n[j]/n[i]
				}
			}
		}
	}
	return
}

//////////// Common Parts code

//////////// PrettyPrinting & Debugging functions
