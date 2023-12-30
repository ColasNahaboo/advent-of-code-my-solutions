// Adventofcode 2023, d16, in go. https://adventofcode.com/2023/day/16
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 46
// TEST: example 51
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

// we use to Scalarrays: one for the mirrors, one for the light beams
// for each lit position,,we also store the direction of the beam so we can
// detect if we have entered a loop and stop propaging the beam

package main

import (
	"flag"
	"fmt"
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

// con is the raw map, with the . / \ | - chars "as read" in it
// lits is a map of lights, as bytes used as bitfields

func part1(lines []string) (sum int) {
	con := parse(lines)			// contraption
	lits := makeScalarray[byte](con.w, con.h) // map of lit places
	lit(&con, &lits, 0, 1)		// send a lightbeam to the right at 0
	VPlit(&lits)
	for _, l := range lits.a {
		if l & LIT != 0 {
			sum++
		}
	}
	return
}

func lit(con *Scalarray[byte], lits *Scalarray[byte], pos, dir int) {
	var out int
	// if a beam has already gone through with the same dir, we are looping, stop
	dirmask := dirMask(dir, con.w)
	if lits.a[pos] & LIT != 0 {
		if lits.a[pos] & dirmask != 0 {
			VPf("  Loop at %v\n", lits.Vector(pos))
			return
		}
	}
	lits.a[pos] |= LIT
	lits.a[pos] |= dirmask

	VPlit(lits)
	
	switch con.a[pos] {
	case '/': 
		out = slashRedir(dir, con.w)
		if con.stepOnceInside(pos, out) {
			lit(con, lits, pos + out, out)
		}
	case '\\':
		out = - slashRedir(dir, con.w)
		if con.stepOnceInside(pos, out) {
			lit(con, lits, pos + out, out)
		}
	case '|':
		for _, out = range barRedir(dir, con.w) {
			if con.stepOnceInside(pos, out) {
				lit(con, lits, pos + out, out)
			}
		}
	case '-':
		for _, out = range dashRedir(dir, con.w) {
			if con.stepOnceInside(pos, out) {
				lit(con, lits, pos + out, out)
			}
		}
	case '.':
		out = dir
		if con.stepOnceInside(pos, out) {
			lit(con, lits, pos + dir, dir)
		}
	}
}

const (
	LIT = 1
	UP = 2
	RIGHT = 4
	DOWN = 8
	LEFT = 16
)

func dirMask(dir, w int) byte {
	switch dir {
	case -w: return UP
	case 1: return RIGHT
	case w: return DOWN
	case -1: return LEFT
	}
	return 0
}

// \ is the opposite of /
func slashRedir(d, w int) int {
	switch d {
	case -w: return 1
	case 1: return -w
	case w: return -1
	case -1: return w
	}
	panic("Bad dir: " + itoa(d))
}

func barRedir(d, w int) []int {
	if d == 1 || d == -1 {
		return []int{-w, w}
	}
	return []int{d}
}

func dashRedir(d, w int) []int {
	if d == w || d == -w {
		return []int{-1, 1}
	}
	return []int{d}
}

func parse(lines []string) Scalarray[byte] {
	con := makeScalarray[byte](len(lines[0]), len(lines))
	for y, line := range lines {
		for x, c := range line {
			con.Set(x, y, byte(c))
		}
	}
	return con
}

//////////// Part 2
func part2(lines []string) (maxenergy int) {
	con := parse(lines)			// contraption
	var x, y, energy int
	var lits Scalarray[byte]
	for x = 0; x < con.w; x++ {
		lits = makeScalarray[byte](con.w, con.h) // map of lit places
		lit(&con, &lits, con.Pos(x, 0), con.w) // lightbeam down from top row
		energy = litCount(lits.a)
		if energy > maxenergy {
			maxenergy = energy
		}

		lits = makeScalarray[byte](con.w, con.h)
		lit(&con, &lits, con.Pos(x, con.h-1), -con.w) // up from bottom row
		energy = litCount(lits.a)
		if energy > maxenergy {
			maxenergy = energy
		}
	}
	for y = 0; y < con.h; y++ {
		lits = makeScalarray[byte](con.w, con.h) // map of lit places
		lit(&con, &lits, con.Pos(0, y), 1) // lightbeam right from left col
		energy = litCount(lits.a)
		if energy > maxenergy {
			maxenergy = energy
		}

		lits = makeScalarray[byte](con.w, con.h)
		lit(&con, &lits, con.Pos(con.w-1, y), -1) // left from right row
		energy = litCount(lits.a)
		if energy > maxenergy {
			maxenergy = energy
		}
	}

	return
}

func litCount (a []byte) (sum int) {
	for _, b := range a {
		if b & LIT != 0 {
			sum++
		}
	}
	return
}
	
//////////// Common Parts code

//////////// PrettyPrinting & Debugging functions

func VPlit(lits *Scalarray[byte]) {
	if ! verbose {
		return
	}
	VPf("Contraption %d x %d:\n", lits.w, lits.h)
	for y := 0; y < lits.h; y++ {
		for x := 0; x < lits.w; x++ {
			if lits.Get(x, y) & LIT != 0 {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}
