// Adventofcode 2022, d13, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 example 13
// TEST: -1 input 5588
// TEST: example 140
// TEST: input 23958
package main

import (
	"flag"
	"fmt"
	"sort"
	//"log"
	//"regexp"
)

type Packet struct {
	i int						// union: either i or p are defined,
	p []Packet					// i defined only if field p is non-nil
}

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
	pair_index := 0
	for n := 0; n < len(lines); n += 3 {
		pair_index++
		p1,_ := parsePacket(lines[n], 0)
		p2,_ := parsePacket(lines[n+1], 0)
		VPf("p1 = %v\np2 = %v\n\n", p1, p2)
		if p1.Less(p2) == 1 {
			VPf("OK, pair_index = %d\n", pair_index)
			sum += pair_index
		}
	}
	return
}

//////////// Part 2
func part2(lines []string) int {
	packets := make([]Packet, 2)
	// the divider packets
	pd2,_ := parsePacket("[[2]]", 0)
	packets[0] = pd2
	pd6,_ := parsePacket("[[6]]", 0)
	packets[1] = pd6
	for _, line := range lines {
		if line == "" { continue;} // ignore blank lines
		p,_ := parsePacket(line, 0)
		packets = append(packets, p)
	}
	// sort the list of packets via the Less method
	sort.Slice(packets, func(i, j int) bool {
		return packets[i].Less(packets[j]) == 1
	})
	key := 1
	for i, p := range packets {
		if p.Less(pd2) == 0 {
			VPf("divider 2 is at line $d\n", i+1)
			key *= i+1
		} else if p.Less(pd6) == 0 {
			VPf("divider 2 is at line $d\n", i+1)
			key *= i+1
		}
	}
	return key
}

//////////// Common Parts code

// dump recursive parsing which doesnt tolerate extra spacing
func parsePacket(s string, i int) (Packet, int) {
	// VPf("parse @%d: %v\n", i, s)
	var p, subp Packet
	if s[i] == '[' {
		p.p = make([]Packet, 0)
		i++
		for s[i] != ']' { 
			subp, i = parsePacket(s, i)
			p.p = append(p.p, subp)
			for s[i] == ',' {
				i++
			}
		}
		return p, i+1
	} else {
		n := 0
		for s[i] >= '0' && s[i] <= '9' {
			n *= 10
			n += int(s[i] - '0')
			i++
		}
		for s[i] == ',' { i++;}
		p.i = n
		return p, i
	}
}

// 1 -> p1 < p2, 0 -> p1 == p2, -1 -> p1 > p2
func (p1 Packet) Less(p2 Packet) int {
	if p1.p == nil && p2.p == nil {
		if p1.i < p2.i {
			return 1
		} else if p1.i > p2.i {
			return -1
		} else {
			return 0
		}
	} else if p1.p != nil && p2.p != nil {
		for i := 0; i < len(p1.p); i++ {
			if i >= len(p2.p) {
				return -1
			}
			res := p1.p[i].Less(p2.p[i])
			if res > 0 {
				return 1
			} else if res < 0 {
				return -1
			}
		}
		if len(p1.p) < len(p2.p) {
			return 1
		}
		return 0
	} else if p1.p == nil && p2.p != nil {
		return iToPacket(p1.i).Less(p2)
	} else {
		return p1.Less(iToPacket(p2.i))
	}
}

func iToPacket(i int) (p Packet) {
	p.p = make([]Packet, 1)
	p.p[0] = Packet{i:i}
	return
}

//////////// Part1 functions

//////////// Part2 functions
