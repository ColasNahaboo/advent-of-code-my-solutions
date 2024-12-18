// Adventofcode 2024, d18, in go. https://adventofcode.com/2024/day/18
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 22
// TEST: example "6_1"
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

// we added the grid size and number of falled blocks
// in the first line of input files

// Part1 is a simple AStar shortest path search
// Part2: we know the value tested in Part1 is OK (has a path). We suppose
// the last fall blocks the path. So we find the first fall blocking the path
// by dichotomy between these values

package main

import (
	"fmt"
	"regexp"
	"flag"
	// "slices"
)

//////////// Options

var commaFlag *bool
var commaSep = "_"

func main() {
	commaFlag = flag.Bool("c", false, "outputs numbers separated by comma instead of underscores")
	ParseOptionsString(2, part1, part2)
}

func ProcessXtraOptions() { // extra options, see ParseOptions in utils.go
	if *commaFlag {
		commaSep = ","
	}
}

//////////// Part 1

func part1(lines []string) string {
	b, _ := parse(lines)
	b.VPBoard("initial", VPBoardBool)
	start, end := Point{0, 0}, Point{b.w - 1, b.h - 1}
	path := AStarFindPath[*Board[bool], Point](b, start, end,
		Neigbours, NeigboursDist, Distance, SamePoint)
	return itoa(len(path) - 1)		// do not count the start point
}

//// The Callbacks needed for astar 

// neighbour nodes
func Neigbours(b *Board[bool], p Point) (nbs []Point) {
	for _, d := range DirsOrtho {
		q := p.Add(d)
		if b.Inside(q) && ! b.Get(q) {
			nbs = append(nbs, q)
		}
	}
	return
}

func NeigboursDist(b *Board[bool], p, q Point) float64 {
	return 1.0
}

func Distance(b *Board[bool], p, q Point) float64 {
	return float64(q.x - p.x + q.y - p.y) // manhattan, supposing q (end) > p
}

func SamePoint(b *Board[bool], p, q Point) bool {
	return p.Equal(q)
}

//////////// Part 2

func part2(lines []string) string {
	b, fok, falls := parse2(lines)
	start, end := Point{0, 0}, Point{b.w - 1, b.h - 1}
	var fko = len(falls)		// index of first fall found to prevent exit
	// we seek fist fko by dichotomy
	for {
		d := fko - fok
		if d <= 1 {
			return fmt.Sprintf("%d%s%d\n", falls[fko].x, commaSep, falls[fko].y)
		}
		f := fok + d/2
		VPf("== [%d %d] Testing %d\n", fok, fko, f)
		for _, fall := range falls[fok+1:f+1] {
			b.Set(fall, true)
		}
		if nil != AStarFindPath[*Board[bool], Point](b, start, end,
			Neigbours, NeigboursDist, Distance, SamePoint) {
			fok = f
		} else {
			// undo the previous block falls
			for _, fall := range falls[fok+1:f+1] {
				b.Set(fall, false)
			}
			fko = f
		}
	}
	return "FAILED"
}

//////////// Common Parts code

func parse(lines []string) (*Board[bool], int) {
	renum := regexp.MustCompile("[[:digit:]]+")
	head := renum.FindAllString(lines[0], -1)
	size := atoi(head[0]) + 1
	falls := atoi(head[1])
	b := MakeBoard[bool](size, size)
	for _, line := range lines[1:falls+1] {
		m := renum.FindAllString(line, -1)
		b.a[atoi(m[0])][atoi(m[1])] = true
	}
	return &b, falls
}

// returns also all the falls, and the board filled up to fall # fok
// we suppose that there was still a path at fok (that part1 succeeded)
func parse2(lines []string) (b *Board[bool], fok int, falls []Point) {
	renum := regexp.MustCompile("[[:digit:]]+")
	b, fok = parse(lines)
	for _, line := range lines[1:] {
		m := renum.FindAllString(line, -1)
		falls = append(falls, Point{atoi(m[0]), atoi(m[1])})
	}
	return b, fok, falls
}

//////////// PrettyPrinting & Debugging functions

func DEBUG() {
	if ! debug { return }
	fmt.Println("DEBUG!")
	// ad hoc debug code here
}


func VPBoardBool(c bool) string {
	if c {
		return "#"
	} else {
		return "."
	}
}
