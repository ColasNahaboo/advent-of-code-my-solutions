// Adventofcode 2022, d09, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 ex 13
// TEST: -1 input 6384
// TEST: ex 1
// TEST: input 2734
package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"
)

type Pos struct {
	x int
	y int
}

var verbose bool
var remotion = regexp.MustCompile("^([RLUD])[[:space:]]+([[:digit:]]+)")
const tails = 9

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

func part1(lines []string) int {
	h := Pos{0,0}				  // head pos 
	t := Pos{0, 0}				  // tail pos
	path := make(map[Pos]bool, 0) // positions visited by t
	path[t] = true
	for lineno, line := range lines {
		motion := remotion.FindStringSubmatch(line)
		if motion == nil {
			log.Fatalf("Syntax error line %d: %s\n", lineno, line)
		}
		steps := atoi(motion[2])
		for i := 0; i < steps; i++ {
			h = h.add(motion[1])
			t = t.follow(h)
			path[t] = true
		}
	}
	return len(path)
}

//////////// Part 2

func part2(lines []string) int {
	h := Pos{0,0}				  // head pos
	t := make([]Pos, tails, tails)	  // the 9(tails) tail knots pos
	for i :=0; i < tails; i++ { t[i] = h;}
	path := make(map[Pos]bool, 0) // positions visited by the last tail
	path[h] = true
	for lineno, line := range lines {
		motion := remotion.FindStringSubmatch(line)
		if motion == nil {
			log.Fatalf("Syntax error line %d: %s\n", lineno, line)
		}
		steps := atoi(motion[2])
		for i := 0; i < steps; i++ {
			h = h.add(motion[1])
			prev := h
			for j :=0; j < tails; j++ {
				t[j] = t[j].follow(prev)
				prev =  t[j]
			}
			path[t[tails - 1]] = true
		}
	}
	return len(path)
}

//////////// Common Parts code

// a single step in one of the 4 directions
func (p Pos) add(dir string) Pos {
	switch dir {
	case "R": return Pos{p.x + 1, p.y}
	case "L": return Pos{p.x - 1, p.y}
	case "U": return Pos{p.x, p.y + 1}
	case "D": return Pos{p.x, p.y - 1}
	default: log.Fatalf("add: bad direction: %s\n", dir)
	}
	return Pos{}				// never reached
}

// make t follow h
func (t Pos) follow(h Pos) Pos {
	if t.touches(h) {
		return t
	} else if t.y == h.y {		// same row, migrate closer on x
		return Pos{(t.x + h.x ) / 2, t.y}
	} else if t.x == h.x {		// same column
		return Pos{t.x, (t.y + h.y ) / 2}
	} else {					// diagonal move
		var dx, dy int			// we know they are != 0
		if h.x > t.x { dx = 1;} else { dx = -1;}
		if h.y > t.y { dy = 1;} else { dy = -1;}
		return Pos{t.x + dx, t.y + dy}
	}
}

// is t on or touches h?
func (t Pos) touches(h Pos) bool {
	return t.x >= (h.x - 1) && t.x <= (h.x + 1) && t.y >= (h.y - 1) && t.y <= (h.y + 1)
}

//////////// Part1 functions

//////////// Part2 functions
