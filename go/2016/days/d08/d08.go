// Adventofcode 2016, d08, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 input 110
// TEST: input
package main

import (
	"flag"
	"fmt"
	// "regexp"
)

type Op struct { code, x, y int;}
const NONE = 0
const RECT = 1
const RROW = 2
const RCOL = 3

// 2D grid of fixed dimensions
const gx = 50
const gy = 6					
var grid [gx*gy]bool
// x and y slices (a row, a column) buffers
var xbuf [gx]bool
var ybuf [gy]bool

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

func part1(lines []string) int {
	ops := parse(lines)
	for _, op := range ops {
		opExec(op)
	}
	return gridLit()
}

//////////// Part 2
func part2(lines []string) int {
	n := part1(lines)
	gridPrint()
	return n
}

//////////// Common Parts code

func parse(lines []string) (prog []Op) {
	var x, y int
	for ln, line := range lines {
		if n, _ := fmt.Sscanf(line, "rect %dx%d", &x, &y); n == 2 {
			prog = append(prog, Op{RECT, x, y})
			continue
		}
		if n, _ := fmt.Sscanf(line, "rotate row y=%d by %d", &y, &x); n == 2 {
			prog = append(prog, Op{RROW, x, y})
			continue
		}
		if n, _ := fmt.Sscanf(line, "rotate column x=%d by %d", &x, &y); n == 2 {
			prog = append(prog, Op{RCOL, x, y})
			continue
		}
		panic(fmt.Sprintf("Syntax error line %d: \"%s\"\n", ln+1, line))
	}
	return
}

func opExec(op Op) {
	switch op.code {
	case RECT:
		for x := 0; x < op.x; x++ {
			for y := 0; y < op.y; y++ {
				grid[x +y*gx] = true
			}
		}
	case RROW:
		for i := 0; i < gx; i++ {
			xbuf[(i + op.x) % gx] = grid[i + op.y*gx]
		}
		for i := 0; i < gx; i++ {
			grid[i + op.y*gx] = xbuf[i]
		}
	case RCOL:
		for i := 0; i < gy; i++ {
			ybuf[(i + op.y) % gy] = grid[op.x + i*gx]
		}
		for i := 0; i < gy; i++ {
			grid[op.x +i*gx] = ybuf[i]
		}
	}
}

//////////// Part1 functions
		
func gridLit() (n int) {
	for _, p := range grid {
		if p { n++;}
	}
	return
}

//////////// Part2 functions

func gridPrint() {
	for y := 0; y < gy; y++ {
		for x := 0; x < gx; x++ {
			if grid[x +y*gx] {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}
	
