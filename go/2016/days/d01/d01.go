// Adventofcode 2016, d01, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 input 234
// TEST: input 113
package main

import (
	"flag"
	"fmt"
	"regexp"
)

type Person struct {
	x, y int					// position
	dx, dy int					// orientation: increments if we step once
}

// we use a map of coords to keep trace of visited places
// this is simpler than a 2D slice
type Coord struct {
	x, y int
}

var restep *regexp.Regexp

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
	doc := fileToString(infile)
	restep = regexp.MustCompile("([RL])([[:digit:]]+)")

	var result int
	if *partOne {
		VP("Running Part1")
		result = part1(doc)
	} else {
		VP("Running Part2")
		result = part2(doc)
	}
	fmt.Println(result)
}

//////////// Part 1

func part1(doc string) int {
	me := new(Person)
	me.dy = 1					// face north
	for _, step := range restep.FindAllStringSubmatch(doc, -1) {
		VP("Step:", step)
		turn(me, step[1])
		run(me, atoi(step[2]))
	}
	return cityDistance(me.x, me.y)
}

//////////// Part 2
func part2(doc string) int {
	me := new(Person)
	me.dy = 1					// face north
	places := make(map[Coord]bool)
	places[Coord{0,0}] = true
	for _, step := range restep.FindAllStringSubmatch(doc, -1) {
		turn(me, step[1])
		twice, x2, y2 := visited(places, me, atoi(step[2]))
		VPf("Step %s: [%d,%d]\n", step[0], me.x, me.y)
		if twice {
			return cityDistance(x2, y2)
		}
	}
	return 0
}

//////////// Common Parts code

func cityDistance(x, y int) int {
	return abs( x) + abs(y)
}

func turn(me *Person, dir string) {
	// turn
	if dir == "R" {
		if me.dy == 1 {
			me.dx = 1; me.dy = 0
		} else if me.dx == 1 {
			me.dx = 0; me.dy = -1
		} else if me.dy == -1 {
			me.dx = -1; me.dy = 0
		} else {
			me.dx = 0; me.dy = 1
		}
	} else {
		if me.dy == 1 {
			me.dx = -1; me.dy = 0
		} else if me.dx == 1 {
			me.dx = 0; me.dy = 1
		} else if me.dy == -1 {
			me.dx = 1; me.dy = 0
		} else {
			me.dx = 0; me.dy = -1
		}
	}
}

func abs(i int) int {
	if i >= 0 {
		return i
	} else {
		return -i
	}
}

//////////// Part1 functions

func run(me *Person, dist int) {
	me.x += dist * me.dx
	me.y += dist * me.dy
}

//////////// Part2 functions

func visited(places map[Coord]bool, me *Person, dist int) (bool, int, int) {
	for i := 0; i < dist; i++ {
		me.x += me.dx
		me.y += me.dy
		xy := Coord{me.x, me.y}
		if places[xy] {
			return true, me.x, me.y
		} else {
			places[xy] = true
		}
	}
	return false, me.x, me.y
}
