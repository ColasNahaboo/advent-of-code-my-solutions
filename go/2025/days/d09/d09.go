// Adventofcode 2025, d09, in go. https://adventofcode.com/2025/day/09
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 50
// TEST: example 24
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

// Input properties:
// numbers (red tiles coordinates) between 1900 and 99000 exclusive
// a line is x,y, with x or y the same between lines (and it wraps last to first)
// less than 500 points. less than 125000 rectangles (pairs of points)
// a map of the area would have more than 6 billion cells
// however, the x take only less than 500 distinct values, and the y 250
// so we can normalize into a 500x250 grid of 125000 cells max

package main

import (
	"fmt"
	"regexp"
	// "flag"
	"slices"
	"cmp"
)

// Implementation:
type Norm []int					// normed coords are indexes to in(put) coords

//////////// Options parsing & exec parts

func main() {
	ExecOptions(2, NoXtraOpts, part1, part2)
}

//////////// Part 1

func part1(lines []string) (maxarea int) {
	redtiles := parse(lines)
	slices.SortFunc(redtiles, CmpPoints)
	for i, p := range redtiles {
		for j := i+1; j < len(redtiles); j++ {
			q := redtiles[j]
			area := Area(p, q)
			if area > maxarea {
				maxarea = area
			}
		}
	}
	return 
}

//////////// Part 2

func part2(lines []string) (res int) {
	normx, normy, reds, inside := parseNormed(lines)
	return inside.BiggestRectangleInside(reds, normx, normy)
}

func parseNormed(lines []string) (normx, normy Norm, reds []Point, inside *Board[bool]) {
	redsin := parse(lines)
	// create the Norm tables with the list of all inputted x and y
	for _, p := range redsin {
		normx = AppendUniqNorm(normx, p.x)
		normy = AppendUniqNorm(normy, p.y)
	}
	slices.Sort(normx)
	slices.Sort(normy)
	inside = MakeBoard[bool](len(normx), len(normy))
	// since we wrap, make the previous point start with the end point
	p0 := Point{slices.Index(normx, redsin[len(redsin)-1].x),
		slices.Index(normy, redsin[len(redsin)-1].y)}
	for _, pin := range redsin {
		// normalize the coordinates
		p := Point{slices.Index(normx, pin.x), slices.Index(normy, pin.y)}
		reds = append(reds, p)
		// draw the perimeter into the inside board
		inside.DrawLine(p0, p, true)
		p0 = p
	}
	// fill the inside, from one point we find inside
	pi := inside.Seed(true)
	VPf("Inside Seed: %v\n", pi)
	inside.FillFrom(pi, true)
	// setup done. Now find the biggest rectangle
	return
}

func AppendUniqNorm(norm Norm, n int) Norm {
	if slices.Contains(norm, n) {
		return norm
	} else {
		return append(norm, n)
	}
}

// using a generic (t T) instead of directly "true" avoids the error:
// "cannot use true (untyped bool constant) as bool value in argument to b.Set"
func (b *Board[T]) DrawLine(p0, p Point, t T) {
	if p0.x == p.x {
		for i := min(p.y, p0.y); i <= max(p.y, p0.y); i++ {
			b.Set(Point{p.x, i}, t)
		}
	} else if p0.y == p.y {
		for i := min(p.x, p0.x); i <= max(p.x, p0.x); i++ {
			b.Set(Point{i, p.y}, t)
		}
	} else {
		panic("Not an orthogonal line!")
	}
}

// fill with t (true) by propagation to neighbours, also to avoid the error:
// "cannot use true (untyped bool constant) as bool t in argument to b.GetOr"
func (b *Board[T]) FillFrom(p Point, t T) {
	b.Set(p, t)
	for _, dir := range DirsOrtho {
		q := p.Add(dir)
		if b.GetOr(q, t) != t {
			b.FillFrom(q, t)
		}
	}
}

// heuristic: we start at [0,1], scan horizontally, to the first edge, then
// get the first empty cell
// A general solution would check that this empty cell is actually inside
// by looking at all the lines crosssed to reach the edges in all directions
func (b *Board[T]) Seed(t T) Point {
	y := 1						// we work only on line #1
	x := 0	
	for b.Get(Point{x, y}) != t { // skip initial empties if any
		x++
	}
	for b.Get(Point{x, y}) == t { // skip edges
		x++
		if x >= b.w {
			panic("Could not find an inside point in row 1")
		}
	}
	return Point{x, y}
}

func (b *Board[T]) BiggestRectangleInside(reds []Point, normx, normy Norm) (maxarea int) {
	for i, p := range reds {
		for j := i+1; j < len(reds); j++ {
			q := reds[j]
			area := AreaXY(normx[p.x], normy[p.y], normx[q.x], normy[q.y])
			if area > maxarea && b.RectIsInside(p, q) {
				maxarea = area
			}
		}
	}
	return
}

func (b *Board[T]) RectIsInside(p, q Point) bool {
	var null T					// get the default value of generic type T
	for x := min(p.x, q.x); x <= max(p.x, q.x); x++ {
		for y := min(p.y, q.y); y <= max(p.y, q.y); y++ {
			if b.Get(Point{x, y}) == null {
				return false
			}
		}
	}
	return true
}

//////////// Common Parts code

func Area(p, q Point) int {
	return AreaXY(p.x, p.y, q.x, q.y)
}

func AreaXY(x0, y0, x1, y1 int) int {
	return (max(x0, x1) - min(x0, x1) + 1) * (max(y0, y1) - min(y0, y1)+ 1)
}

func parse(lines []string) (redtiles []Point) {
	reline := regexp.MustCompile("^([[:digit:]]+),([[:digit:]]+)")
	for _, line := range lines {
		if m := reline.FindStringSubmatch(line); m != nil {
			redtiles = append(redtiles, Point{atoi(m[1]), atoi(m[2])})
		}
	}
	return
}

const HalfMax = 100000000
func CmpPoints(p, q Point) int {
	return cmp.Compare(p.x * HalfMax + p.y, q.x * HalfMax + q.y)
}

//////////// PrettyPrinting & Debugging functions. See also the VPx functions.

func DEBUG() {
	if ! debug { return }
	fmt.Println("DEBUG!")
	// ad hoc debug code here
}
