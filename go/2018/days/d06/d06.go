// Adventofcode 2018, d06, in go. https://adventofcode.com/2018/day/06
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 17
// TEST: example 16
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"fmt"
	"regexp"
	"sort"
	// "flag"
	// "slices"
)

//////////// Options parsing & exec parts

func main() {
	ExecOptions(2, NoXtraOpts, part1, part2)
}

var infinity = maxint/10		// everything above it is considered infinite

//////////// Part 1

func part1(lines []string) (res int) {
	locs := parse(lines)
	ox, oy, ow, oh := BoundingBox(locs) // the area enclosing all the locations
	ex, ey := ox + ow, oy + oh			// coords of far corner, excluded
	areas := make([]int, len(locs), len(locs))
	// first list infinities by checking the points on the border
	for x := ox; x < ex; x++ {
		if closest := ClosestXY(locs, x, oy); closest >= 0 { // top row
			areas[closest] = infinity
		}
		if closest := ClosestXY(locs, x, ey-1); closest >= 0 { // bottom row
			areas[closest] = infinity
		}
	}
	for y := oy+1; y < ey-1; y++ { // sides, less corners
		if closest := ClosestXY(locs, ox, y); closest >= 0 { // left side
			areas[closest] = infinity
		}
		if closest := ClosestXY(locs, ex-1, y); closest >= 0 { // right side
			areas[closest] = infinity
		}
	}
	// then compute the areas inside the border point per point
	for x := ox+1; x < ex-1; x++ {
		for y := oy+1; y < ey-1; y++ {
			if closest := ClosestXY(locs, x, y); closest >= 0 {
				VPf("  [%d %d] closest to %s\n", x, y, LocName(closest))
				areas[closest]++
			}
		}
	}
	for i, area := range areas {
		VPf("  area of #%d = %d\n", i, area)
		if area < infinity && area > res {
			res = area
		}
	}
	return 
}

type DistToLoc struct {d, l int} // d = manhattan distance to location #l

func ClosestXY(locs []Point, x, y int) int {
	dists := make([]DistToLoc, len(locs), len(locs))
	p := Point{x, y}
	for i, loc := range locs {
		dists[i].d = p.ManDist(loc)
		dists[i].l = i
	}
	sort.Slice(dists, func(i, j int) bool { return dists[i].d < dists[j].d })
	if dists[0].d == dists[1].d {	// more than one are closest
		VPf("  [%d %d] not closest to single loc\n", x, y)
		return -1
	}
	return dists[0].l
}
		

//////////// Part 2

func part2(lines []string) (res int) {
	locs := parse(lines)
	ox, oy, ow, oh := BoundingBox(locs) // the area enclosing all the locations
	ex, ey := ox + ow, oy + oh			// coords of far corner, excluded
	maxdist := 10000
	if ow < 10 { maxdist = 32 }	//  reduced size for example
	VPf("  %d, %d, %d, %d, maxdist = %d\n", ox, oy, ex, ey, maxdist)

	for x := ox; x < ex; x++ {
	POINT:
		for y := oy; y < ey; y++ {
			p := Point{x, y}
			dist := 0
			for _, loc := range locs {
				dist += p.ManDist(loc)
				if dist >= maxdist {
					continue POINT
				}
			}
			VPf("  [%d %d] is only at %d, ok\n", x, y, dist)
			res++
		}
	}
	return
}

//////////// Common Parts code

func parse(lines []string) (locations []Point) {
	renum := regexp.MustCompile("[[:digit:]]+") // example code body, replace.
	for _, line := range lines {
		coords := atoil(renum.FindAllString(line, -1))
		locations = append(locations, Point{coords[0], coords[1]})
	}
	return
}

//////////// PrettyPrinting & Debugging functions. See also the VPx functions.

func LocName(i int) string {
	return string(rune('A' + i))
}

func DEBUG() {
	if ! debug { return }
	fmt.Println("DEBUG!")
	// ad hoc debug code here
}
