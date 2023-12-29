// Adventofcode 2023, d14, in go. https://adventofcode.com/2023/day/14
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 136
// TEST: example 64
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"flag"
	"fmt"
)

var verbose bool

const (							// what is on the 2D grid?
	NONE = 0
	CUBE = 1
	ROLL = 2
)

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
		//VP("Running Part2")
		result = part2(lines)
	}
	fmt.Println(result)
}

//////////// Part 1

func part1(lines []string) int {
	sa := parse(lines)
	tiltN(sa)
	return load(sa)
}

func parse(lines []string) (sa Scalarray[int]) {
	sa = makeScalarray[int](len(lines[0]), len(lines))
	p := 0
	for _, line := range lines {
		for _, c := range line {
			switch c {
			case 'O': sa.a[p] = ROLL
			case '#': sa.a[p] = CUBE
			}
			p++
		}
	}
	return
}

func load(sa Scalarray[int]) (sum int) {
	for p := 0; p < sa.w * sa.h; p++ {
		if sa.a[p] == ROLL {
			sum += sa.h - p/sa.w
		}
	}
	return
}

//////////// Part 2
func part2(lines []string) int {
	sa := parse(lines)
	past := []Scalarray[int]{sa.Clone()}
	for i := 1; i < 1000000000; i++ {
		cycle(sa)
		for old := range past {
			if sa.Equal(past[old]) {		// we loop
				offset := ((1000000000 - old) % (i-old))
				VPf("Cycle %d, DejaVu on %d, cycle length: %d, offset %d\n", i, old, i-old, offset)
				return load(past[old + offset])
			}
		}
		past = append(past, sa.Clone())
	}
	return load(sa)
}

func cycle(sa Scalarray[int]) {
	tiltN(sa)
	tiltW(sa)
	tiltS(sa)
	tiltE(sa)
}

func saaEqual(a1, a2 []int) bool {
	for i, v := range a1 {
       if v != a2[i] {
           return false
       }
   }
	return true
}

func tiltN(sa Scalarray[int]) {
	for p := 0; p < sa.w * sa.h; p++ {
		if sa.a[p] == ROLL {
			var pos int
			// roll north (up) while on platform and we have free space
			for pos = p; pos - sa.w >= 0 && sa.a[pos - sa.w] == NONE
			pos -= sa.w {
			}
			if pos != p {		// if actually moved, update platform
				sa.a[pos] = ROLL
				sa.a[p] = NONE
			}
		}
	}
	return
}

func tiltS(sa Scalarray[int]) {
	for p :=  sa.w * sa.h - 1; p >= 0; p-- {
		if sa.a[p] == ROLL {
			var pos int
			// roll south (down) while on platform and we have free space
			for pos = p; pos + sa.w < sa.w*sa.h && sa.a[pos + sa.w] == NONE
			pos += sa.w {
			}
			if pos != p {		// if actually moved, update platform
				sa.a[pos] = ROLL
				sa.a[p] = NONE
			}
		}
	}
	return
}

func tiltE(sa Scalarray[int]) {
	for p :=  sa.w * sa.h - 1; p >= 0; p-- {
		if sa.a[p] == ROLL {
			var pos int
			// roll east (right) while on platform and we have free space
			for pos = p; (pos+1) / sa.w == pos / sa.w && sa.a[pos + 1] == NONE
			pos++ {
			}
			if pos != p {		// if actually moved, update platform
				sa.a[pos] = ROLL
				sa.a[p] = NONE
			}
		}
	}
	return
}

func tiltW(sa Scalarray[int]) {
	for p := 0; p < sa.w * sa.h; p++ {
		if sa.a[p] == ROLL {
			var pos int
			// roll west (left) while on platform and we have free space
			// we add sa.w because 1 and -1 both give 0 divided by sa.w
			for pos = p; (pos+sa.w-1) / sa.w == (pos+sa.w) / sa.w && sa.a[pos - 1] == NONE
			pos-- {
			}
			if pos != p {		// if actually moved, update platform
				sa.a[pos] = ROLL
				sa.a[p] = NONE
			}
		}
	}
	return
}

//////////// Common Parts code

//////////// PrettyPrinting & Debugging functions

func VPplatform(label string, sa Scalarray[int]) {
	if ! verbose {
		return
	}
	if label != "" {
		fmt.Println(label)
	}
	for p := 0; p < sa.w * sa.h; p++ {
		if p > 0 && p % sa.w == 0 {
			fmt.Println()
		}
		switch sa.a[p] {
		case ROLL: fmt.Print("O")
		case CUBE: fmt.Print("#")
		default: fmt.Print(".")
		}
	}
	fmt.Println()
}

