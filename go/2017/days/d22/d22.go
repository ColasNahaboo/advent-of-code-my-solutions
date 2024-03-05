// Adventofcode 2017, d22, in go. https://adventofcode.com/2017/day/22
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 5587
// TEST: example 2511944
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"flag"
	"fmt"
	// "golang.org/x/exp/slices"
)

var verbose, debug bool

const (							// part 1 uses only CLEAN and INFECTED
	CLEAN = 0
	WEAKENED = 1
	INFECTED = 2
	FLAGGED = 3
)

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
	puck := parse(lines)
	for burst := 0; burst < 10000; burst++ {
		puck.Burst1()
	}
	return puck.infections
}

// a move, a burst of activity of the virus
func (p *Puck) Burst1() {
	if p.grid[p.pos] == INFECTED { // infected ==> turn right and clean
		p.dir = p.dir.Right()
		p.grid[p.pos] = CLEAN
	} else {					// clean ==> turn left and infect
		p.dir = p.dir.Left()
		p.grid[p.pos] = INFECTED
		p.infections++
	}
	p.pos = p.pos.Add(p.dir)	// advance forward 1 step
}	

//////////// Part 2

func part2(lines []string) int {
	puck := parse(lines)
	for burst := 0; burst < 10000000; burst++ {
		puck.Burst2()
	}
	return puck.infections
}

// a move, a burst of activity of the virus
func (p *Puck) Burst2() {
	switch p.grid[p.pos] {
	case CLEAN:
		p.dir = p.dir.Left()
		p.grid[p.pos] = WEAKENED
	case WEAKENED:
		// keep direction
		p.grid[p.pos] = INFECTED
		p.infections++
	case INFECTED:
		p.dir = p.dir.Right()
		p.grid[p.pos] = FLAGGED
	case FLAGGED:
		p.dir = p.dir.UTurn()
		p.grid[p.pos] = CLEAN
	}
	p.pos = p.pos.Add(p.dir)	// advance forward 1 step
}	

//////////// Common Parts code

func parse(lines []string) Puck {
	size := len(lines[0])
	start := size / 2
	grid := IGrid{}
	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			if lines[y][x] == '#' {
				grid[Pos{x, y}] = INFECTED
			}
		}
	}
	return Puck{pos: Pos{start, start}, dir: Pos{0, -1}, grid: grid}
}

////// Puck (copied from d19)

type Puck struct {
	pos Pos
	dir Pos
	grid IGrid					// the infinite grid
	infections int				// number of bursts having caused an infection
}

////// Infinite 2D grid

type IGrid map[Pos] int

////// Pos

type Pos struct {
	x, y int
}

func (p *Pos) Inside(x, y int) bool {
	if p.x < 0 || p.y < 0 || p.y >= y || p.x >= x {
		return false
	}
	return true
}

func (p *Pos) Add(dir Pos) Pos {
	return Pos{p.x + dir.x, p.y + dir.y}
}

// smart way to turn a direction [x, y] right or left
func (p *Pos) Right() Pos {
	return Pos{-p.y, p.x}
}
func (p *Pos) Left() Pos {
	return Pos{p.y, -p.x}
}
func (p *Pos) UTurn() Pos {
	return Pos{-p.x, -p.y}
}

//////////// PrettyPrinting & Debugging functions

func DEBUG() {
	if ! debug { return }
	// ad hoc debug code here
}
