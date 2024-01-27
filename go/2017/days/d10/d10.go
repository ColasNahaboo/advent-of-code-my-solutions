// Adventofcode 2017, d10, in go. https://adventofcode.com/2017/day/10
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 12
// TEST: example-p2-1 a2582a3a0e66e6e86e3812dcb672a272
// TEST: example-p2-2 33efeb34ea91902bb2f59c9920caa6cd
// TEST: example-p2-3 3efbe78a8d82f29979031a4aa0b16a9d
// TEST: example-p2-4 63960835bcdc130f0b66d7ff4f6a5a8e
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

// in the input files we added the list size as on the second line

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
	lens, list := parse1(lines)
	p, skip := 0, 0
	twistAll(p, skip, lens, list)
	return int(list[0]) * int(list[1])
}

func parse1(lines []string) (lengths []int, list []byte) {
	re := regexp.MustCompile("[[:digit:]]+")
	m := re.FindAllString(lines[0], -1)
	for _, s := range m {
		lengths = append(lengths, atoi(s))
	}
	listsize := atoi(lines[1])
	for i := 0; i < listsize; i++ {
		list = append(list, byte(i))
	}
	return
}

func twistAll(p, skip int, lens []int, list []byte) (int, int) {
	VPf("lens: %v\nlist: %v\n", lens, list)
	VPlist(p, list)
	buf = make([]byte, len(list), len(list))
	for _, l := range lens {
		p = twist(p, skip, l, list)
		VPlist(p, list)
		skip++
	}
	return p, skip
}

var buf []byte

func twist(p, skip, l int, list []byte) (newpos int) {
	// copy the section at p of length len, wrapping index in list
	for i, j := p, 0; i < p+l; i, j = i+1, j+1 {
		buf[j] = list[i % len(list)]
	}
	// and re-paste it in place in reverse
	for i, j := p, l-1; i < p+l; i, j = i+1, j-1 {
		list[i % len(list)] = buf[j]
	}
	return (p + l + skip) % len(list)
}
	
//////////// Part 2

func part2(lines []string) string {
	lens, list := parse2(lines)
	p, skip := 0, 0
	for i := 0; i < 64; i++ {
		p, skip = twistAll(p, skip, lens, list)
	}
	return denseHash(list)
}

func parse2(lines []string) (lengths []int, list []byte) {
	for _, b := range lines[0] {
		lengths = append(lengths, int(b))
	}
	lengths = append(lengths, []int{17, 31, 73, 47, 23}...)
	listsize := atoi(lines[1])
	for i := 0; i < listsize; i++ {
		list = append(list, byte(i))
	}
	return
}

func denseHash(list []byte) (s string) {
	dh := make([]byte, 16, 16)	// hash as array of ints
	for i := range dh {
		x := list[i*16]
		for j := 1; j < 16; j++ {
			x ^= list[i*16+j]
		}
		dh[i] = x
		s += fmt.Sprintf("%02x", x)
	}
	return
}

//////////// Common Parts code

//////////// PrettyPrinting & Debugging functions

func VPlist(p int, list []byte) {
	if ! verbose { return }
	for i, b := range list {
		if i == p {
			fmt.Printf(" [%d]", int(b))
		} else {
			fmt.Printf(" %d", int(b))
		}
	}
	fmt.Println()
}
