// Adventofcode 2017, d19, in go. https://adventofcode.com/2017/day/19
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example ABCDEF
// TEST: example 38
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"flag"
	"fmt"
	"strings"
	// "golang.org/x/exp/slices"
)

type Pos struct {
	x, y int
}

type Puck struct {
	pos Pos
	dir Pos
	what byte					// what is under the puck?
	grid []string				// the grid on which the path is drawn
	path []byte					// the letters encountered in the path
	steps int					// for part2, number of steps in path
}

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

func part1(lines []string) string {
	puck := PuckStart(lines)
	for puck.Next() {
	}
	return string(puck.path)
}

//////////// Part 2
func part2(lines []string) int {
	puck := PuckStart(lines)
	for puck.Next() {
	}
	return puck.steps
}

//////////// Common Parts code

func Rotated(b byte) byte {
	if b == '|' {
		return '-'
	} else {
		return '|'
	}
}

////// Puck

// start on the | on first row, towards south
func PuckStart(lines []string) (p Puck) {
	p.pos.x = strings.IndexByte(lines[0], '|')
	p.dir.y = 1
	p.grid = lines
	p.what = p.What()
	p.steps = 1
	return
}

// what symbol is under the puck?
func (p *Puck) What() byte {
	return p.pos.What(p.grid)
}

// move puck in dir until we are out of a tunnel
func (p *Puck) Move(dir Pos) {
	crossing := Rotated(p.what)
	for {
		p.pos.x += p.dir.x
		p.pos.y += p.dir.y
		p.steps++
		if what := p.What(); what != crossing {
			VPf("Moving to %s at %v, dir %v\n", string(p.What()), p.pos, p.dir)
			if what >= 'A' && what <= 'Z' {
				p.path = append(p.path, what)
			}
			return
		}
	}
}

// turn puck to face direction
func (p *Puck) Face(dir Pos) {
	p.dir = dir
}

// move the puck once
func (p *Puck) Next() bool {
	np := p.pos.Add(p.dir)
	nw := Rotated(p.what)
	switch np.What(p.grid) {
	case p.what:				// continue on straight path
		p.Move(p.dir)
	case ' ':					// end of path reached
		return false
	case nw:					// crossing, tunnel under it
		p.Move(p.dir)
	case '+':					// turn
		if rd := p.dir.Right(); p.Legit(np, rd, nw)  { // after crosses, find letter, '+' or same path
			p.Move(p.dir)		  // move onto '+' and turn right
			p.what = nw
			p.Face(rd)
		} else if ld := p.dir.Left(); p.Legit(np, ld, nw)  {
			p.Move(p.dir)		  // move onto '+' and turn right
			p.what = nw
			p.Face(ld)
		} else {
			panic(fmt.Sprintf("At %d, %d, + without exit", np.x, np.y))
		}
	default:					// Letter
		p.Move(p.dir)
	}
	return true
}

// skip crossings to check for a valid continuation path at the end of the tunnel
func (p *Puck) Legit(pos, dir Pos, crossing byte) bool {
	for {
		pos = pos.Add(dir)
		if ! pos.InGrid(p.grid) {
			return false
		}
		what := pos.What(p.grid)
		if what == ' ' {
			return false
		}
		if what != crossing {
			return true
		}
	}
}

////// Pos

func (p *Pos) What(grid []string) byte {
	return grid[p.y][p.x]
}

func (p *Pos) InGrid(grid []string) bool {
	if p.x < 0 || p.y < 0 || p.y >= len(grid) || p.x >= len(grid[p.y]) {
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
