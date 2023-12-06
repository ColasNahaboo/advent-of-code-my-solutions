// Adventofcode 2023, d03, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: d03-input,RESULT1,RESULT2.test
// TEST: -1 example 4361
// TEST: example 467835
// And any file named d03-DESCRIPTION,RESULT1,RESULT2.test containing an input

// grid []int is a grid with an added border of 1 width, with non-zero value
// marking symbol positions. If the input grid is 10x10, our grid is 12x12


package main

import (
	"flag"
	"fmt"
	// "regexp"
)

type Num struct {
	value int		// the value of the number
	pos int			// its starting pos in grid
	len int			// its length
}
type Gear struct {
	pos int			// position of the (potential) gear in grid
	numidx int		// index in nums array of previous num, or -1
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
		infile = fileMatch("input,[0-9]*,[0-9]*.test")
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

func part1(lines []string) (sum int) {
	grid, nums, _, gw, gh := readGrid(lines)
	VPf("  %d numbers in grid %dx%d\n", len(nums), gw-2, gh-2)
	for _, num := range nums {
		if ! alone(num, grid, gw, gh) {
			VPf("  number OK: %d\n", num.value)
			sum += num.value
		} else {
			VPf("  number ko: %d\n", num.value)
		}
	}
	return
}

func alone(num Num, grid []int, gw, gh int) bool {
	o := num.pos - gw - 1 // upper left corner of border around number
	for i := 0; i < num.len + 2; i++ {
		if grid[o + i] != 0 { // top row
			return false
		}
		if grid[o + 2*gw + i] != 0 { // bottom row
			return false
		}
	}
	if grid[o + gw] != 0 {	// left side
		return false
	}
	if grid[o + gw + num.len + 1] != 0 { // right side
		return false
	}
	return true
}

//////////// Part 2
func part2(lines []string) (sum int) {
	grid, nums, gears, gw, gh := readGrid(lines)
	VPf("  %d numbers & %d proto-gears in grid %dx%d\n", len(nums), len(gears), gw-2, gh-2)
	for i, gear := range gears {
		adjnums, adjvalue := adjacentToGear(gear, grid, nums, gw, gh)
		if adjnums == 2 {
			VPf("  gear#%d OK, value %d\n", i+1, adjvalue)
			sum += adjvalue
		}
			
	}
	return
}

// count nums adjacent to gear
func adjacentToGear(g Gear, grid []int, nums []Num, gw, gh int) (n, value int) {
	adjacents := make([]int, 0)
	for i := g.numidx; i >= 0; i-- { // look backwards
		if nums[i].pos + nums[i].len + 1 < g.pos - gw { //  too far
			break
		}
		if isAdjacent(g.pos, nums[i], gw) {
			adjacents = append(adjacents, nums[i].value)
		}
	}
	for i := g.numidx + 1; i < len(nums); i++ { // look forwards
		if nums[i].pos - 1 > g.pos + gw { //  too far
			break
		}
		if isAdjacent(g.pos, nums[i], gw) {
			adjacents = append(adjacents, nums[i].value)
		}
	}
	n = len(adjacents)
	if n == 2 {
		value = adjacents[0] * adjacents[1]
	}
	return	
}

// look if proto-gear at pos gpos is adjacent to a number
func isAdjacent(gpos int, num Num, gw int) bool {
	o := num.pos - gw - 1 // upper left corner of border around number
	for i := 0; i < num.len + 2; i++ {
		if gpos == o + i { // top row
			return true
		}
		if gpos == o + 2*gw + i { // bottom row
			return true
		}
	}
	if gpos == o + gw {	// left side
		return true
	}
	if gpos == o + gw + num.len + 1 { // right side
		return true
	}
	return false
}

//////////// Common Parts code

func readGrid(lines []string) (grid []int, nums []Num, gears []Gear, gw, gh int) {
	gw = len(lines[0]) + 2
	gh = len(lines) + 2
	grid = make([]int, gw * gh)
	for y, line := range lines {
		for x := 0; x < len(line); x++ {
			c := line[x]
			if c >= '0' && c <= '9' {
				// grid[(y+1)*gw + x+1] = int(c - '0' + 1)
				value := 0
				pos := (y+1)*gw + x + 1
				size := 0
				for  x < len(line) && line[x] >= '0' && line[x] <= '9' {
					value = value * 10 + int(line[x] - '0')
					x++
					size++
				}
				x-- // avoid eating the char after the number
				nums = append(nums, Num{value, pos, size})
			} else if c != '.' {
				grid[(y+1)*gw + x+1] = int(c)
				if c == '*' { // potential gear
					gears = append(gears, Gear{(y+1)*gw + x+1, len(nums) - 1})
				}
			}
		}
	}
	return
}

//////////// Part1 functions

//////////// Part2 functions
