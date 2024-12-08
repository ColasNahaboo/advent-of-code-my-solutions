// Adventofcode 2024, d08, in go. https://adventofcode.com/2024/day/08
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 14
// TEST: example 34
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input



package main

import (
	"flag"
	"fmt"
	// "golang.org/x/exp/slices"
)

var verbose, debug bool

type Grid struct {
	w, h int				// dimensions of the grid
	freqs []Antenna			// list of antennas for each frequency
}

type Antenna struct {
	freq byte					// its [0-9a-zA-Z] name, the frequency
	places []Point			   // list of antennas positions for this frequency
}

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	debugFlag := flag.Bool("V", false, "debug: even more verbose")
	flag.Parse()
	verbose = *verboseFlag
	debug = *debugFlag
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
	grid := parse(lines)
	locations := make(map[Point]bool) // locations of all antinodes
	for _, ant := range grid.freqs {
		for i, posi := range ant.places {
			for j := i+1; j < len(ant.places); j++ {
				FindAntinodes(grid, locations, posi, ant.places[j])
			}
		}
	}
	return len(locations)
}

// adds valid antinodes to locations
func FindAntinodes(grid *Grid, locations map[Point]bool, p1, p2 Point) {
	for _, antinode := range []Point{
		p1.Add(p2.Delta(p1)), p2.Add(p1.Delta(p2))} {
		if antinode.IsInGrid(grid.w, grid.h) {
			locations[antinode] = true
		}
	}
}

//////////// Part 2

func part2(lines []string) int {
	grid := parse(lines)
	locations := make(map[Point]bool) // locations of all antinodes
	for _, ant := range grid.freqs {
		for i, posi := range ant.places {
			for j := i+1; j < len(ant.places); j++ {
				FindAllAntinodes(grid, locations, posi, ant.places[j])
				FindAllAntinodes(grid, locations, ant.places[j], posi)
			}
		}
	}
	return len(locations)
}

// adds valid antinodes on harmonics on direction p1 -> p2, starting at p2
func FindAllAntinodes(grid *Grid, locations map[Point]bool, p1, p2 Point) {
	step := p1.Delta(p2)
	for anti := p2; anti.IsInGrid(grid.w, grid.h); anti = anti.Add(step) {
		locations[anti] = true
	}
}

//////////// Common Parts code

func parse(lines []string) (* Grid) {
	g := Grid{}
	g.w = len(lines[0])
	g.h = len(lines)
	ids := make(map[byte]int) // indexes of frequency names in g.freqs
	for y, line := range lines {
		for x, r := range line {
			if r == '.' {
				continue
			}
			b := byte(r)
			var id int
			var ok bool
			if id, ok = ids[b]; ! ok {
				id = len(g.freqs) // declare new frequency
				ids[b] = id
				g.freqs = append(g.freqs, Antenna{b, []Point{}})
			}
			g.freqs[id].places = append(g.freqs[id].places, Point{x, y})
		}
	}
	return &g
}

//////////// PrettyPrinting & Debugging functions

func DEBUG() {
	if ! debug { return }
	// ad hoc debug code here
}
