// Adventofcode 2017, d11, in go. https://adventofcode.com/2017/day/11
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 3
// TEST: -1 example1 0
// TEST: -1 example2 2
// TEST: -1 example3 3
// TEST: example
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
		infile = fileMatch("input,[[:alnum:]]*,[[:alnum:]]*.test")
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

func part1(lines []string) int {
	steps := parse(lines[0])
	dest := move3D([3]int{0,0,0}, steps)
	return dist3D([3]int{0,0,0}, dest)
}

//////////// Part 2
func part2(lines []string) int {
	steps := parse(lines[0])
	distance := move3Dmax([3]int{0,0,0}, steps)
	return distance
}

func move3Dmax(p [3]int, steps []int) (dmax int) {
	for _, step := range steps {
		p = add3D(p, dirs[step])
		d := dist3D([3]int{0,0,0}, p)
		if d > dmax {
			dmax = d
		}
	}
	return
}

//////////// Common Parts code

// 3D representation of hex coordinates, see README.md
var dirnames = []string{"n", "ne", "se", "s", "sw", "nw"}
var dirs = [][3]int{{0, -1, 1}, {1, -1, 0}, {1, 0, -1}, {0, 1, -1}, {-1, 1, 0}, {-1, 0, 1}}

func parse(line string) (steps []int) {
	re := regexp.MustCompile("[nesw]+")
	m := re.FindAllString(line, -1)
	for _, dn := range m {
		d := IndexOf[string](dirnames, dn)
		if d == -1 {
			panic("Bad direction: " + dn)
		}
		steps = append(steps, d)
	}
	return
}

func move3D(p [3]int, steps []int) (dest [3]int) {
	dest = p
	for _, step := range steps {
		dest = add3D(dest, dirs[step])
	}
	return
}

func add3D(p, q [3]int) (r [3]int) {
	for i := 0; i < 3; i++ {
		r[i] = p[i] + q[i]
	}
	return
}

// manhattan distance in 3D is twice the hexa one
func dist3D(p, q [3]int) (d int) {
	for i := 0; i < 3; i++ {
		d += intAbs(p[i] - q[i])
	}
	return d/2
}
	
//////////// PrettyPrinting & Debugging functions
