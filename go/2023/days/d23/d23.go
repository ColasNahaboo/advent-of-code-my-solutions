// Adventofcode 2023, d23, in go. https://adventofcode.com/2023/day/23
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 94
// TEST: example 154
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

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

func part1(lines []string) int {
	parse(lines)
	hike := makeScalarray[bool](area.w, area.h)
	hike.Set(1, 0, true) 		// dont go back there
	hike.Set(1, 1, true)		// mandatory 1st step: one down
	maxhikelen = 0
	explore(hike, hike.Pos(1, 1), 1)
	return maxhikelen
}

func parse(lines []string) {
	area = makeScalarray[byte](len(lines[0]), len(lines))
	for y, line := range lines {
		for x, char := range line {
			switch char {
			case '#': area.Set(x, y, WALL)
			case '^': area.Set(x, y, NSLOPE)
			case '>': area.Set(x, y, ESLOPE)
			case 'v': area.Set(x, y, SSLOPE)
			case '<': area.Set(x, y, WSLOPE)
			}
		}
	}
	end = len(area.a) - 2		// bottom right less one wall
	start = 1
	areadirs = area.Dirs()
}

// explore from pos
func explore(hike Scalarray[bool], pos, steps int) {
	for _, dir := range areadirs {
		p := pos + dir
		if p == end {			// exit reached!
			steps++
			if steps > maxhikelen {
				maxhikelen = steps
			}
			return
		} else if area.a[p] == PATH && ! hike.a[p] { // we can go there
			hike2 := hike.Clone()
			hike2.a[p] = true
			explore(hike2, p, steps + 1)
		} else if area.a[p] >= NSLOPE && ! hike.a[p] { // slide on slope
			sdir := areadirs[area.a[p] - NSLOPE]
			if area.a[p+sdir] != WALL && ! hike.a[p+sdir] { // slope bottom free
				hike2 := hike.Clone()
				hike2.a[p] = true
				hike2.a[p+sdir] = true
				explore(hike2, p+sdir, steps + 2)
			}
		}
	}
}

//////////// Part 2

func part2(lines []string) int {
	parse2(lines)
	area2.Set(1, 0, true) 		// dont go back there
	area2.Set(1, 1, true)		// mandatory 1st step: one down
	maxhikelen = 0
	explore2(area2.Pos(1, 1), 1)
	return maxhikelen
}

// part2: ignore slopes
func parse2(lines []string) {
	area2 = makeScalarray[bool](len(lines[0]), len(lines))
	for y, line := range lines {
		for x, char := range line {
			if char == '#' {
				area2.Set(x, y, true)
			}
		}
	}
	end = len(area2.a) - 2		// bottom right less one wall
	start = 1
	areadirs = area2.Dirs()
}

// simplified, as no slopes
// for efficiency, we do not "fork" a new hike state, reuse the same one
// and work directly inside the area, marking steps as walls
func explore2(pos, steps int) {
	for _, dir := range areadirs {
		p := pos + dir
		if ! area2.a[p] { // we can go there
			if p == end {			// exit reached!
				steps++
				if steps > maxhikelen {
					maxhikelen = steps
				}
				return
			}
			area2.a[p] = true
			explore2(p, steps + 1)
			area2.a[p] = false
		}
	}
}

//////////// Common Parts code

const (
	PATH = 0
	WALL = 1
	NSLOPE = 2
	ESLOPE = 3
	SSLOPE = 4
	WSLOPE = 5
)

var area Scalarray[byte]		// the static landscape, hike is only the steps
var area2 Scalarray[bool]		// same for part2
var areadirs [4]int				// Up Right Down Left (N E S W)
var start, end int				// entry and exit points of the area
var maxhikelen int				// the max length of the hike paths found
	

//////////// PrettyPrinting & Debugging functions
