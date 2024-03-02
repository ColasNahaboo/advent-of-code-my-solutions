// Adventofcode 2017, d20, in go. https://adventofcode.com/2017/day/20
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 0
// TEST: example 2
// TEST: example2 1
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"flag"
	"fmt"
	"regexp"
	// "golang.org/x/exp/slices"
)

var verbose bool

// an easier to spot maxint in debug than 9223372036854775807 (^uint(0) >> 1)
const maxint = 8888888888888888888

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

// we just find the smallest acceleration (manhattan-style)

func part1(lines []string) int {
	partnum := 0
	partmin := -1
	distmin := maxint
	// ignore minus signs, as we only consider the absolute value anyways
	re := regexp.MustCompile(", a=<-?([0-9]+),-?([0-9]+),-?([0-9]+)")
	for _, line := range lines {
		m := re.FindStringSubmatch(line)
		if m == nil {
			panic("Syntax error: " + line)
		}
		md := atoi(m[1]) + atoi(m[2]) + atoi(m[3])
		if md < distmin {
			distmin = md
			partmin =  partnum
		}
		partnum++
	}
	return partmin
}

//////////// Part 2

// We use a Map indexed by the 3D position to detect collisions

const QUIET = 100				//  if no collisions in last QUIET times, stop

func part2(lines []string) int {
	parts := parse(lines)
	nparts := len(parts)
	time := 0
	quiet := 0					//  number of moves without any collision
	VPf("  [%d] %d particules\n", time, nparts)
	for quiet < 100 {
		time++
		quiet++
		// move all particules, obtain new state
		parts2 := make(map[Pos3]*Part)
		for _, p := range parts {
			if p.deleted {
				continue
			}
			p.Move()
			if pp, exists := parts2[p.p]; exists { // collision with a PreviousP
				quiet = 0
				if ! pp.deleted {
					nparts--	// forget pp if not already deleted
				}
				nparts--		  // forget p
				pp.deleted = true		   // make previous as deleted
				VPf("  [%d] Collision at %v, remain: %d\n", time, pp.p, nparts)
			} else {
				parts2[p.p] = p
			}
		}
		parts = parts2
	}
	return nparts
}

type Part struct {
	deleted bool
	p, v, a Pos3
}
type Pos3 struct {
	x, y, z int
}

func (p *Part)Move() {
	p.v.x += p.a.x
	p.v.y += p.a.y
	p.v.z += p.a.z
	p.p.x += p.v.x
	p.p.y += p.v.y
	p.p.z += p.v.z
}

var	reline = regexp.MustCompile("p=<([^>]+)>, *v=<([^>]+)>, *a=<([^>]+)")
var	retuple = regexp.MustCompile("([-0-9]+),([-0-9]+),([-0-9]+)")

func parse(lines []string) (parts map[Pos3]*Part) {
	parts = make(map[Pos3]*Part)
	for _, line := range lines {
		tuples := MustFindStringSubmatch(reline, line)
		p := Part{p: parseTuple(tuples[1]), v: parseTuple(tuples[2]), a: parseTuple(tuples[3])}
		parts[p.p] = &p
	}
	return
}

func parseTuple(s string) Pos3 {
	coords := MustFindStringSubmatch(retuple, s)
	return Pos3{atoi(coords[1]), atoi(coords[2]), atoi(coords[3])}
}

func MustFindStringSubmatch(re *regexp.Regexp, s string) (m []string) {
	m = re.FindStringSubmatch(s)
	if m == nil {
		panic("No match on: " + s)
	}
	return
}

//////////// Common Parts code

//////////// PrettyPrinting & Debugging functions
