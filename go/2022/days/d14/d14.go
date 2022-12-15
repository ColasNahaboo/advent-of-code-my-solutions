// Adventofcode 2022, d14, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 example 24
// TEST: -1 input 618
// TEST: example 93
// TEST: input 26358
package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"
	"time"
)

type Point struct {
	x int
	y int
}

type Cave struct {
	width, height, area int		// dimensions of the grid
	sourcex int					// source position is [sourcex, 0]
	grid []int					// point [x,y] is at index x+y*width in grid
}

const air = 0
const rock = 1
const sand = 2
const source = 3
var gridglyphs = []string{".", "#", "o", "+"}

var verbose bool
var anim bool
var lastlines int
var repath = regexp.MustCompile("([[:digit:]]+),([[:digit:]]+)")

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	animFlag := flag.Bool("a", false, "animation: simple ascii-art animation")
	lastFlag := flag.Int("l", 0, "in anim and verbose mode, only display last L lines")
	flag.Parse()
	verbose = *verboseFlag
	anim = *animFlag
	lastlines = *lastFlag
	if anim { verbose = true;}
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
	cave := parse(lines)
	return flow(cave)
}

//////////// Part 2
func part2(lines []string) int {
	cave := parse(lines)
	// add floor on last line of grid
	for i := (cave.height - 1) * cave.width; i < cave.area; i++ {
		cave.grid[i] = rock
	}
	return fill(cave)
}

//////////// Common Parts code

// inflexible parsing, returns the Cave object trimmed to dimensions
func parse(lines []string) Cave {
	xmin := 2000000; ymin := xmin
	xmax := 0; ymax := xmax
	// collect the path points, adding a pseudo-point [-1,-1] as line/path termination
	points := make([]Point, 0)
	for _, line := range lines {
		path := repath.FindAllStringSubmatch(line, -1)
		for _, p := range path {
			x := atoi(p[1])
			y := atoi(p[2])
			points = append(points, Point{x:x, y:y})
			if x > xmax {
				xmax = x
			} else if x < xmin {
				xmin = x
			}
			if y > ymax {
				ymax = y
			} else if y < ymin {
				ymin = y
			}
		}
		points = append(points, Point{x:-1, y:-1}) // [-1,-1] marks end of line
	}
	VPf("%d <= x <= %d, %d <= y <= %d\n", xmin, xmax, ymin, ymax)
	// we are sure that everything will happen between xmin-1 and xmax+1
	// reposition the coordinate system so that we can work inside a grid
	// encompassing [xmin-1, xmax+2[ and [0, ymax+1[
	// point [x,y] is at index x+y*width in grid
	// Update: for Part2, we need to ensure that there is ymax+3 room on
	// each side on sourcex, so we pad on the side with air if needed.
	// The extra room is no issue for Part1 as long as we do not add the floor
	offset := xmin - 1
	height := ymax + 3
	if 500 - offset < height {
		offset = 500 - height
	}
	if xmax - 500 < height {
		xmax = 500 + height
	}
	
	for i, p := range points {
		if p.x >= 0 {
			points[i] = Point{p.x - offset, p.y}
		} else {
			points[i] = Point{-1, -1}
		}
	}
	width := xmax - offset + 1
	area := width * height
	grid := make([]int, area, area)
	sourcex := 500 - offset
	grid[sourcex] = source
	VPf("grid %d x %d, size = %d, offset = %d\n", width, height, area, offset)

	// now, draw the rocks along path segments. In the grid, 0=air, 1=rock, 2=sand
	// we draw the start point of the segment, then the segment excluding the
	// start point, but including the last point
	pp := Point{-1,-1}				// Previous Point: starting point of path
	for _, p := range points {
		if p.x != -1 {			// end of path, do nothing special
			if pp.x == -1 {		// start of path, mark start point
				grid[p.x + p.y * width] = rock
			} else if p.y == pp.y { // horiz segment ]pp, p]
				if p.x >= pp.x {	// left to right
					for x := pp.x+1; x <= p.x; x++ {
						grid[x + p.y * width] = rock
					}
				} else {
					for x := pp.x-1; x >= p.x; x-- {
						grid[x + p.y * width] = rock
					}
				}
			} else { 			// vertical segment ]pp, p]
				if p.y >= pp.y {	// bottom to top
					for y := pp.y+1; y <= p.y; y++ {
						grid[p.x + y * width] = rock
					}
				} else {
					for y := pp.y-1; y >= p.y; y-- {
						grid[p.x + y * width] = rock
					}
				}
			}
		}
		pp = p
	}
	return Cave{width, height, area, sourcex, grid}
}

func VPcave(c Cave) {
	if verbose {
		start := 0
		if lastlines > 0 { start = c.height - lastlines;}
		for y := start; y < c.height; y++ {
			for x := 0; x < c.width; x++ {
				fmt.Printf("%s", gridglyphs[c.grid[x+y*c.width]])
			}
			fmt.Printf("\n")
		}
		fmt.Printf("\n")
		if anim {
			time.Sleep(200 * time.Millisecond)
		}
	}
}

// simulate flow of one sand packet. Return true if not overflowed
func flowOnce(c Cave) bool {
	p := c.sourcex
	for true {
		pp := p + c.width		// get down one step
		if pp >= c.area {
			return false		// overflow
		}
		if c.grid[pp] > air {	// blocked
			if c.grid[pp - 1] == air { // look diagonal left
				p = pp - 1 // flow left
			} else if c.grid[pp + 1] == air { // look diagonal right
				p = pp + 1 // flow right
			} else {
				c.grid[p] = sand
				return true		// rest here
			}
		} else {
			p = pp				// free fall down one step in air
		}
	}
	log.Fatalln("Could not resovle flow")
	return false				// notreached
}	

//////////// Part1 functions


// return number of sand units before overflow
func flow(c Cave) (n int) {
	VPcave(c)
	for flowOnce(c) {
		VPcave(c)
		n++
	}
	return
}

//////////// Part2 functions

// return number of sand units before fullup
func fill(c Cave) (n int) {
	VPcave(c)
	for {
		flowOnce(c)
		VPcave(c)
		n++
		if c.grid[c.sourcex] == sand {
			return
		}
	}
	return
}
