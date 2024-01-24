// Adventofcode 2017, d03, in go. https://adventofcode.com/2017/day/03
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example1 0
// TEST: -1 example2 3
// TEST: -1 example3 2
// TEST: -1 example4 31
// TEST: example
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"flag"
	"fmt"
)

var verbose bool
var inputarg int

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	numFlag := flag.Int("n", 0, "input number")
	flag.Parse()
	verbose = *verboseFlag
	var infile string
	if flag.NArg() > 0 {
		infile = flag.Arg(0)
	} else {
		infile = fileMatch("input,[0-9]*,[0-9]*.test")
	}
	lines := fileToLines(infile)
	inputarg = *numFlag
	var n int
	if inputarg != 0 {
		n = inputarg
	} else {
		n = atoi(lines[0])
	}

	var result int
	if *partOne {
		VP("Running Part1")
		result = part1(n)
	} else {
		VP("Running Part2")
		result = part2(n)
	}
	fmt.Println(result)
}

//////////// Part 1

func part1(n int) int {
	//spiralCornersTest(); return 0			// DEBUG
	
	x, y := spiralCorners(n)
	d := intAbs(x) + intAbs(y)
	VPf("%d is at (%d, %d), distance %d\n", n, x, y, d)
	return d
}

func spiralCorners(n int) (x, y int) {
	// x, y = coords of current corner
	//VPc(1, x, y)
	if n <= 1 {				// i <= n <= k: n on the current segment
			return
	}
	l := 1				  // the loop number
	i := 2				  // i the number in its SE start corner
	x, y = 1, 0
	var k int			  // k is the value of i will have at next Korner
	for {
		//VPc(i, x, y)
		k = i + l * 2 - 1		// i = its SE corner, going N to k (NW)
		if n <= k {				// i <= n <= k: n on the current segment
			y -= n-i			// coords computed from the i-corner towards k
			return
		}

		y = - l
		i = k
		//VPc(i, x, y)
		k = i + l * 2			// next corner: NE, going W to NW
		if n <= k {
			x -= n-i			
			return
		}

		x = -l
		i = k
		//VPc(i, x, y)
		k = i + l * 2			// next corner: NW, going S to DW
		if n <= k {
			y += n-i
			return
		}

		y = l
		i = k
		//VPc(i, x, y)
		k = i + l * 2 + 1		// next corner: SW, going E to next loop SE
		if n <= k {
			x += n-i
			return
		}

		l++						// go to new loop start point: 
		x = l
		i = k
	}
}

//////////// Part 2
func part2(n int) int {
	tabInit(12)					// should suffice for nearly all values
	return turtleWalk(n)
}

type V2 struct {				// 2D vector
	x, y int
}

// walk the turtle in tab, adding until we found return value larger than n
func turtleWalk(n int) int {
	dirs := [4]V2{{1,0},{0,-1},{-1,0},{0,1}} // E N W S
	pos := V2{0, 0}			// center
	dir := 0				// E
	val := 1				// point value
	i := 1					// for debug: step number
	tabSet(pos, val)		// initial state
	for val <= n {
		pos = pos.Move(dirs[dir])		// one step in our dir
		i++
		left := pos.Move(dirs[TurnLeft(dir)]) // what is at our left?
		if tabGet(left) == 0 { // nothing? we moved past the prev corner
			dir = TurnLeft(dir)		  // turn left
		}
		val = sumAdjacent(pos)	// compute sum of vals of adjacent points
		tabSet(pos, val)
		VPf("  [%d] %v = %d\n", i, pos, val)
	}
	return val
}

func (pos V2) Move(step V2) V2 {
	return V2{pos.x + step.x, pos.y + step.y}
}
func TurnLeft(dir int) int {
	return (dir + 1) % 4
}
func sumAdjacent(pos V2) int {
	x, y := pos.x+tabn, pos.y+tabn
	return tab[x-1][y-1] + tab[x][y-1] + tab[x+1][y-1] + tab[x-1][y] + tab[x+1][y] + tab[x-1][y+1] + tab[x][y+1] + tab[x+1][y+1]
}


//////////// Common Parts code

//////////// PrettyPrinting & Debugging functions

var tab [][]int
var tabn int

func tabInit(n int) {
	tabn = n
	tab = make([][]int, tabn*2, tabn*2)
	for i := range tab {
		tab[i] = make([]int, tabn*2, tabn*2)
	}
}
func tabGet(pos V2) int {
	return tab[pos.x+tabn][pos.y+tabn]
}
func tabSet(pos V2, val int) {
	tab[pos.x+tabn][pos.y+tabn] = val
}


func spiralCornersTest() {
	tab = make([][]int, 11)
	for i := range tab {
		tab[i] = make([]int, 11)
	}
	spiralCorners(111)
	for _, line := range tab {
		for _, n := range line {
			fmt.Printf("%4d", n)
		}
		fmt.Println()
	}
}

func VPc(i, x, y int) {
	VPf("  %d (%d, %d)\n", i, x, y)
	tab[y+5][x+5] = i
}
