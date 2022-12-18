// Adventofcode 2022, d18, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 example 64
// TEST: -1 input
// TEST: example
// TEST: input
package main

import (
	"flag"
	"fmt"
	"log"
	// "regexp"
)

// we modelize the 3D grid space, as a linear array of positions,
// a position p is: x + y*size + z*area
var	size, area, volume int		// dims of world: cube enclosing the lava droplet
var cube []bool					// is there a rock bit at position?
var bits []int					// the positions of rock bits in the droplet
var adjacent [6]int				// the relative position offsets of the 6 adjacent bits
var vapor []bool				// part2: is there vapor at position?

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

func part1(lines []string) (exposed int) {
	parse(lines)
	for _, bit := range bits {
		exposed += 6
		for _, adj := range adjacent {
			p := bit + adj
			if p >= 0 && p < volume && cube[p] {
				exposed--
			}
		}
	}
	return
}

//////////// Part 2
func part2(lines []string) (outside int) {
	parse(lines)
	vapor = make([]bool, volume, volume) // is filled with vapor?
	for z := 0; z < size; z += size-1 {	 // top & bottom sides
		for x := 0; x < size; x++ {
			for y := 0; y < size; y++ {
				injectVapor(x + y*size + z*area)
			}
		}
	}
	for z := 1; z < size-1; z++ { // middle slices: perimeter x,y
		for x := 0; x < size; x += size-1 { // x edges
			for y := 0; y < size; y++ {
				injectVapor(x + y*size + z*area)
			}
		}
		for x := 1; x < size-1; x++ { // y edges
			for y := 0; y < size; y += size {
				injectVapor(x + y*size + z*area)
			}
		}
	}

	for _, bit := range bits {
		for _, adj := range adjacent {
			p := bit + adj
			if p < 0 || p >= volume || vapor[p] {
				outside++
			}
		}
	}
	return
}

//////////// Common Parts code

func parse(lines []string) {
	var x, y, z, max int
	bits = make([]int, 0)
	for lineno, line := range lines {
		_, err := fmt.Sscanf(line, "%d,%d,%d", &x, &y, &z)
		if err != nil {
			log.Fatalf("Syntax error line %d: %s\n", lineno, line)
		}
		if x > max { max = x;}
		if y > max { max = y;}
		if z > max { max = z;}
	}
	size = max + 1
	area = size*size
	volume = area*size
	cube = make([]bool, volume)
	adjacent[0] = 1
	adjacent[1] = -1
	adjacent[2] = size
	adjacent[3] = -size
	adjacent[4] = area
	adjacent[5] = -area

	for _, line := range lines {
		fmt.Sscanf(line, "%d,%d,%d", &x, &y, &z)
		cube[x + y*size + z*area] = true
		bits = append(bits, x + y*size + z*area)
	}
}

//////////// Part1 functions

//////////// Part2 functions

func injectVapor(bit int) {
	if !cube[bit] && !vapor[bit] {
		vapor[bit] = true		// fill it
		for _, adj := range adjacent { // expand to neighbors
			p := bit + adj
			if p >= 0 && p < volume && !cube[p] && !vapor[p] {
				injectVapor(p)
			}
		}
	}
}
