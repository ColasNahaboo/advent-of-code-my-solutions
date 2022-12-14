// Adventofcode 2022, d08, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 ex 21
// TEST: ex 8
// TEST: -1 input 1711
// TEST: input 301392
package main

import (
	"flag"
	"fmt"
	//"log"
	// "regexp"
)

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
	// grid is the array of tree heights. index of [x,y] is x + y*width
	width := len(lines[0])
	height := len(lines)
	grid := make([]int, width * height, width * height)
	cur := 0
	for _, line := range lines {
		for _, b := range line {
			grid[cur] = int(b - '0')
			cur++
		}
	}
	
	var result int
	if *partOne {
		VP("Running Part1")
		result = part1(grid, width, height)
	} else {
		VP("Running Part2")
		result = part2(grid, width, height)
	}
	fmt.Println(result)
}

//////////// Part 1

func part1(grid []int, width, height int) (nlit int) {
	// "raytracing", we mark true all trees "painted" by a light ray from the edges
	// in a "light grid": "lit"
	area := width * height
	lit := make([]bool, area, area)

	// lit each row from left then right
	for pos := 0; pos < area; pos += width {
		z := -1					 // shadows height from left
		for p := pos ; p < pos + width; p++ {
			if grid[p]  > z {
				lit[p] = true
				z = grid[p]
			}
		}
		z = -1					 // shadows height from right
		for p := pos + width -1; p >= pos; p-- {
			if grid[p] > z {
				lit[p] = true
				z = grid[p]
			}
		}
	}
		
	// lit each col from top then bottom
	for pos := 0; pos < width; pos ++ {
		z := -1					 // shadows height from top
		for p := pos ; p < area; p += width {
			if grid[p] > z {
				lit[p] = true
				z = grid[p]
			}
		}
		z = -1					 // shadows height from bottom
		for p := pos + area - width; p >= 0; p -= width {
			if grid[p] > z {
				lit[p] = true
				z = grid[p]
			}
		}
	}

	// count trees lit
	for p := 0 ; p < area; p++ {
		if lit[p] {
			nlit++
		}
	}
	if verbose {
		for y := 0; y < height; y++ {
			for x := 0 ; x < width; x++ {
				if lit[x + y*width] {
					fmt.Printf("*")
				} else {
					fmt.Printf(".")
				}
			}
			fmt.Printf("\n")
		}
	}
	return
}

//////////// Part 2
func part2(grid []int, width, height int) (maxss int) {
	// now we "raytrace", but in reverse, from the tree
	// we do not consider border trees as their viewing distances wil be 0
	var d, p int
	area := width * height
	
	for y := 1; y < height-1; y++ {
		x0 := y * width			// row limits
		x1 := (y+1) * width
		for x := 1 ; x < width-1; x++ {
			ss, tl, tr, td, tu := 1, 0, 0, 0, 0
			pos := x + y*width
			z := grid[pos]
			// look right
			d = 1 				// pos increment
			for p = pos + d; p < x1; p += d {
				tr++
				if grid[p] >= z { break; }
			}
			ss *= tr
			// look left
			d = -1 				// pos increment
			for p = pos + d; p >= x0; p += d {
				tl++
				if grid[p] >= z { break; }
			}
			ss *= tl
			// look down
			d = width			// pos increment
			for p = pos + d; p < area; p += d {
				td++
				if grid[p] >= z { break; }
			}
			ss *= td
			// look up
			d = -width			// pos increment
			for p = pos + d; p >= 0; p += d {
				tu++
				if grid[p] >= z { break; }
			}
			ss *= tu

			VPf("[%d,%d](%d) ss = %d: l=%d r=%d d=%d u=%d\n", x, y, z, ss, tl, tr, td, tu)
			if ss > maxss {
				maxss = ss
			}
		}
	}
	return
}

//////////// Common Parts code

//////////// Part1 functions

//////////// Part2 functions
