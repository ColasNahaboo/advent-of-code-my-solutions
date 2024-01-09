// Adventofcode 2023, d21, in go. https://adventofcode.com/2023/day/21
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example
// TEST: example
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

// For part1, we explore step by step, on step N marking as "N+1" all the plots
// reachable in one step from any plot already marked "N", and then we iterate,
// starting again from all the plots marked N+1, until we reache the number of
// desired steps

package main

import (
	"flag"
	"fmt"
)

var verbose bool

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	stepsFlag := flag.Int("n", 64, "number of mandatory steps to take")
	flag.Parse()
	verbose = *verboseFlag
	var infile string
	if flag.NArg() > 0 {
		infile = flag.Arg(0)
	} else {
		infile = fileMatch("input,[0-9]*,[0-9]*.test")
	}
	lines := fileToLines(infile)

	var result int
	if *partOne {
		VP("Running Part1")
		result = part1(lines, *stepsFlag)
	} else {
		//VP("Running Part2")
		result = part2(lines, *stepsFlag)
	}
	fmt.Println(result)
}

//////////// Part 1

func part1(lines []string, steps int) int {
	sa, start := parse(lines)
	VPgarden("Start", sa, steps)
	starts := []int{start}
	for i := 1; i <= steps; i++ {
		starts = explore(sa, starts, byte(i) + STEP0)
	}
	VPgarden("Result", sa, steps)
	return len(starts)
}

// explore all PLOTs reachable in one step from position in starts list,
// mark them as step on the list (avoid doubles), returns their list
func explore(sa Scalarray[byte], starts []int, steps byte) (stepped []int) {
	for _, pos := range starts {
		for _, dir := range sa.Dirs() {
			p := pos + dir
			if sa.a[p] == WALL { 
				continue		// our exploration aborts
			}
			if sa.a[p] != steps { // not already marked reachable for this step
				sa.a[p] = steps		// p is reachable in N steps
				stepped = append(stepped, p)
			}
		}
		}
		return
}

func parse(lines []string) (sa Scalarray[byte], start int) {
	sa = makeScalarrayB[byte](len(lines[0]), len(lines), 1)
	sa.DrawBorder(WALL)
	for y, line := range lines {
		for x, char := range line {
			switch char {
			case '#': sa.SetB(x, y, WALL)
			case 'S': start = sa.PosB(x, y)
			}
		}
	}
	return
}

//////////// Part 2

func part2(lines []string, steps int) int {
	garden, start := parse2(lines)
	var coeffs [3]int			// quadratic function f coefficients
	explore2(garden, start, 26501365, &coeffs)
	return quadFunc(26501365/len(garden), coeffs)
}

// parse2: a "true" 2D slice, note a Scalarray
func parse2(lines []string) (garden [][]bool, start Point) {
	garden = make([][]bool, len(lines[0]), len(lines[0]))
	for i := range lines {
		garden[i] = make([]bool, len(lines), len(lines))
	}
	for y, line := range lines {
		for x, char := range line {
			switch char {
			case '#': garden[x][y] = true
			case 'S': start = Point{x, y}
			}
		}
	}
	return
}

var pdirs = [4]Point{Point{0, -1}, Point{1, 0}, Point{0, 1}, Point{-1, 0}}

func explore2(garden [][]bool, start Point, maxMoves int, coeffs *[3]int) int {
	var visited = make(map[int][]Point)
	// move := 0
	visited[0] = append(visited[0], start)
	found := 0
	prevLen := 0
	for move := 0; move < maxMoves; move++ {
		for _, currentPos := range visited[move] {
			for _, dir := range pdirs {
				targetPos := Point{currentPos.x + dir.x, currentPos.y + dir.y}
				if isMoveValid(garden, targetPos) {
					if ! targetPos.isIn(visited[move+1]) {
						visited[move+1] = append(visited[move+1], targetPos)
					}
				}
			}
		}
		if (move % len(garden)) == (maxMoves % len(garden)) {
			VPf("Move %d %d, prevLen %d\n", move, len(visited[move]), len(visited[move]) - prevLen)
			prevLen = len(visited[move])
			(*coeffs)[found] = prevLen
			found++
			}
		if found == 3 {
			break
		}
	}
	return len(visited[len(visited)-1])
}

func isMoveValid(garden [][]bool, targetPos Point) bool {
	realX := ((targetPos.x % len(garden))+ len(garden)) % len(garden)
	realY := ((targetPos.y % len(garden))+ len(garden)) % len(garden)

	if realX < 0 || realX > len(garden[0])-1 || realX < 0 || realY < 0 || realY > len(garden)-1 || garden[realY][realX] {
		return false
	}
	return true
}

func quadFunc(x int, a [3]int) int {
	b0 := a[0]
	b1 := a[1]-a[0]
	b2 := a[2]-a[1]
	return b0 + b1*x + (x*(x-1)/2)*(b2-b1)

}

type Point struct {
	x, y int
}

type Move struct {
	pos   Point
	steps int
}

func (p Point) Add(q Point) Point {
	return Point{p.x + q.x, p.y + q.y}
}

func (p Point) isIn(list []Point) bool {
	for _, v := range list {
		if v == p {
			return true
		}
	}
	return false
}

//////////// Common Parts code

const (
	PLOT = 0					// a "virgin" plot
	WALL = 1					// impassable wall (or border)
	STEP0 = 2
	// N					    // plot reachable after N-STEP0 steps
)

//////////// PrettyPrinting & Debugging functions

func VPgarden(label string, sa Scalarray[byte], steps int) {
	if ! verbose { return }
	fmt.Printf("%s: array %dx%d, with border %d:", label, sa.w, sa.h, sa.b)
	for p, b := range sa.a {
		if p % sa.w == 0 {
			fmt.Println()
		}
		if b == WALL {
			fmt.Print("#")
		} else if b == byte(steps) + STEP0 {
			fmt.Print("O")
		} else {
			fmt.Print(".")
		}
	}
	fmt.Println()
}

func VPgarden2(label string, garden [][]bool, starts map[Point]bool, nrocks int) {
	if ! verbose { return }
	sa := makeScalarray[byte](len(garden), len(garden[0]))
	fmt.Printf("%s: array %dx%d with %d rocks and %d reachable plots:", label, len(garden), len(garden[0]), nrocks, len(starts))
	for x := range garden {
		for y := range garden[x] {
			if garden[x][y] {
				sa.Set(x, y, byte('#'))
			} else {
				sa.Set(x, y, byte('.'))
			}
		}
	}
	for p := range starts {
		if p.x >= 0 && p.x < len(garden) && p.y >= 0 && p.y < len(garden[0]) {
			if garden[p.x][p.y] {
				panic("Step on rock!")
			}
			sa.Set(p.x, p.y, byte('O'))
		}
	}
	
	for p, b := range sa.a {
		if p % sa.w == 0 {
			fmt.Println()
		}
		fmt.Printf("%s", string(byte(b)))
	}
	fmt.Println()
}
