// Adventofcode 2024, d12, in go. https://adventofcode.com/2024/day/12
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example1 140
// TEST: -1 example2 772
// TEST: -1 example3 1930
// TEST: example1 80
// TEST: example2 436
// TEST: example3 1206
// TEST: example4 236
// TEST: example5 368
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

// We "color" the garden by filling regions with adjacent plots with same plants
// Then we can look at each plot and raise a fence for each adjacent plot with
// a different plant, and count them.

// For part2, we look at all the fences, and see if they are aligned with any
// other fence of the same region. In which case, we just do not count the
// fence by decreasing the perimeter field (now a sides field) by one.
// But, to take into account the exemple5.txt, we do not consider two
// fences as aligned if there is also a 3rd fence joining at their intersection

package main

import (
	"flag"
	"fmt"
	//"regexp"
	// "slices"
)

var verbose, debug bool

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

type Plot struct {
	plant byte
	region *Region
}

type Region struct {			// distinct regions
	plant byte					// 1 to 27
	area, peri int
	fences []Fence				// for part2
}

type Fence struct {
	dir int						// 0 = inactive, 1 = horizontal, 2 = vertical
	a, b Point					// start and end points. b always >= a
}

var noregion = Region{}
var nilregion = &noregion		// the "nil" of *Region

func part1(lines []string) (res int) {
	g := parse(lines)
	regions := []*Region{} 
	for x, col := range g.a {
		for y, plot := range col {
			if plot.region == nilregion { // unprocessed, start new region
				region := Region{plant: plot.plant}
				regions = append(regions, &region)
				plot.Flow(g, x, y, &region)
			}
		}
	}
	ComputePeris(g)
	VPregions("Regions",regions)
	return ComputePrice(g, regions)
}

func (p *Plot) Flow(g *Board[*Plot], x, y int, r *Region) {
	r.area++
	p.region = r
	for d := range 4 {		 // flow into adjacent plots
		nx, ny := x + DirsOrtho[d].x, y + DirsOrtho[d].y
		if g.InsideXY(nx, ny) && // flow into:
			g.a[nx][ny].plant == p.plant &&	  // same plant
			g.a[nx][ny].region == nilregion { // not yet belonging to region
			g.a[nx][ny].Flow(g, nx, ny, r)
		}
	}
}

func ComputePeris(g *Board[*Plot]) {
	for x, col := range g.a {
		for y, plot := range col {
			for d := range 4 {		 // look into adjacent plots
				nx, ny := x + DirsOrtho[d].x, y + DirsOrtho[d].y
				if ! g.InsideXY(nx, ny) ||
					g.a[nx][ny].region != plot.region { // raise a fence
					g.a[x][y].region.peri++
				}
			}
		}
	}
}

func ComputePrice(g *Board[*Plot], regions []*Region) (total int) {
	for _, r :=range regions {
		total += r.area * r.peri
	}
	return
}
	


//////////// Part 2

// Part2 is exactly Part1, except the peri field now count the number of sides,
// as we remove continuation fences from its count via DetectSides()

func part2(lines []string) (res int) {
	g := parse(lines)
	regions := []*Region{} 
	for x, col := range g.a {
		for y, plot := range col {
			if plot.region == nilregion { // unprocessed, start new region
				region := Region{plant: plot.plant}
				regions = append(regions, &region)
				plot.Flow(g, x, y, &region)
			}
		}
	}
	RaiseFences(g)
	DetectSides(regions)
	VPregions("Regions",regions)
	return ComputePrice(g, regions)
}

func RaiseFences(g *Board[*Plot]) {
	for x, col := range g.a {
		for y, plot := range col {
			for d := range 4 {		 // look into adjacent plots
				nx, ny := x + DirsOrtho[d].x, y + DirsOrtho[d].y
				if ! g.InsideXY(nx, ny) ||
					g.a[nx][ny].region != plot.region { // raise a fence
					g.a[x][y].region.fences = append(g.a[x][y].region.fences, MakeFence(x, y, d))
					g.a[x][y].region.peri++
				}
			}
		}
	}
}

func DetectSides(regions []*Region) {
	for _, r := range regions {
		for i, f := range r.fences {
			for j := i + 1; j < len(r.fences); j++ {
				if f.dir == r.fences[j].dir {
					if (f.a == r.fences[j].b && NoCorner(f.dir, f.a, r.fences)) ||
						(f.b == r.fences[j].a && NoCorner(f.dir, f.b, r.fences)) {
						r.peri--
						break
					}
				}
			}
		}
	}
}

func NoCorner(dir int, p Point, fences []Fence) bool {
	for _, f := range fences {
		if f.dir != dir && (p == f.a || p == f.b) {
			return false
		}
	}
	return true
}

func MakeFence(x, y, d int) Fence {
	var nx, ny, dir int
	switch d {
	case DirsOrthoN: dir = 1
		nx, ny = x + 1, y
	case DirsOrthoE: dir = 2
		x++
		nx, ny = x, y + 1
	case DirsOrthoS: dir = 1
		y++
		nx, ny = x + 1, y
	case DirsOrthoW: dir = 2
		nx, ny = x, y + 1
	default: panic("Bad DirsOrthoX:" + itoa(d))
	}
	return Fence{dir, Point{x, y}, Point{nx, ny}}
}

//////////// Common Parts code

func parse(lines []string) *Board[*Plot] {
	b := parseBoard[*Plot](lines, func(x, y int, r rune) *Plot {
		plant := byte(r - 'A') // A=0, Z=25
		return &Plot{plant, nilregion}
	})
	return b
}

//////////// PrettyPrinting & Debugging functions

func DEBUG() {
	if ! debug { return }
	// ad hoc debug code here
}

func VPregions(title string, regions []*Region) {
	if verbose {
		fmt.Printf("%s:", title)
		for _, r := range regions {
			fmt.Printf(" {%s %d %d}=%d", string('A' + rune(r.plant)), r.area, r.peri, r.area * r.peri)
		}
		fmt.Println()
	}
}
