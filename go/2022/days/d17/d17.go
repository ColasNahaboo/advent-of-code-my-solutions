// Adventofcode 2022, d17, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 example 3068
// TEST: -1 input 3191
// TEST: example 1514285714288
// TEST: input 1572093023267

// 

package main

import (
	"flag"
	"fmt"
	// "log"
	"regexp"
	"sort"
)

// directions
const right = 0
const down = 1
const left = 2
const up = 3

// coordinates in grid are scalar from bottom left, of width 7, (9 with padding)
// that is bottom left is at 0, point right is 1, just up is at pos+gw
const width = 7					// usable internal width
const gw = width + 2			// allocated width
var height = 0					// height of used space by rocks
var gh = 0						// total allocated grid height
const nrocks = 5
var emptyGridRow []bool
var dirdelta = [4]int{1, -gw, -1, gw}

type Rock struct {
	name string	   // for debugging: - + ⅃ | ■
	id int		   // 0 - 4 in the list above
	height	int	   // total height
	body []int	   // list of positions of parts, from bottom-left edge
}

var grid []bool
var jets []int
var currock int
var curjet int

var rocks = [nrocks]Rock{
	{
		name: "-",
		id: 0,
		height: 1,
		body: []int{0, 1, 2, 3},
	},
	{
		name: "+",
		id: 1,
		height: 3,
		body: []int{1, gw, gw+1, gw+2, gw*2+1},
	},
	{
		name: "⅃",
		id: 2,
		height: 3,
		body: []int{0, 1, 2, gw+2, gw*2+2},
	},
	{
		name: "|",
		id: 3,
		height: 4,
		body: []int{0, gw, gw*2, gw*3},
	},
	{
		name: "■",
		id: 4,
		height: 2,
		body: []int{0, 1, gw, gw+1},
	},
}

var verbose bool

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	preFlag := flag.Bool("p", false, "preprocess file for part2")
	flag.Parse()
	verbose = *verboseFlag
	infile := "input.txt"
	if flag.NArg() > 0 {
		infile = flag.Arg(0)
	}
	lines := fileToLines(infile)
	jets = parseJets(lines[0])
	VPf("Jets loop size: %d\n", len(jets))
	initGrid()
	// side walls
	emptyGridRow = make([]bool, gw, gw)
	emptyGridRow[0] = true
	emptyGridRow[gw - 1] = true
	
	var result int
	if *preFlag {
		part2pre(infile)
	} else if *partOne {
		VP("Running Part1")
		result = part1()
	} else {
		VP("Running Part2")
		result = part2(infile)
	}
	fmt.Println(result)
}

//////////// Part 1

func part1() int {
	time := 0
	for i := 0; i < 2022; i++ {
		time = dropRock(time)
	}
	return height
}

//////////// Part 2
func part2pre(infile string) {
	time := 0
	loopsize := len(jets) * 5	// common looping between rocks and jet cycles
	dhs := []int{}				// delta heights (dh) for loops
	dhf := make(map[int]int, 0)	// number of occurences of each dh
	// we sample on the first loops
	oh := height
	for l := 0; l < 1000; l++ {
		for i := 0; i < loopsize; i++ {
			time = dropRock(time)
		}
		dhs = append(dhs, height - oh)
		dhf[height - oh]++
		oh = height
	}
	// we sort loop numbers per number of occurences
	sorted := make([][2]int, 0, len(dhf))
	for l, n := range dhf {
		sorted = append(sorted, [2]int{l, n})
	}
	sort.Slice(sorted, func(i, j int) bool { return sorted[i][1] < sorted[j][1];})
	// take rarest,but repeating dh
	i := 0
	for sorted[i][1] <= 1 { i++;}
	dh := sorted[i][0]
	vo := 0
	for dhs[vo] != dh { vo++;}
	vd := 1
	for dhs[vo + vd] != dh { vd++;}
	fmt.Printf("Edit the part2 function in d17.go with these values for infile \"%s\":\n    vd = %d\n    vo = %d\n", infile, vd, vo)
}

func part2(infile string) int {
	var vd, vo int
	// These constants for an input  file are found by running "d17 -p" on it
	if infile == "example.txt" {		// example
		vd = 7
		vo = 1
	} else {					// input.txt
		vd = 344
		vo = 81
	}
	time := 0
	loopsize := len(jets) * 5	// sequence repeats
	nrocks := 1000000000000		// total rocks to run

	// skip non-repeating early vaues
	for i := 0; i < vo * loopsize; i++ {
		time = dropRock(time)
	}
	ho := height
	// determine height gain of a single loop
	for i := vo * loopsize; i < (vo + vd) * loopsize; i++ {
		time = dropRock(time)
	}
	hd := height - ho
	// we skip nloops...
	nloops := (nrocks - vo * loopsize) / (vd * loopsize)
	// compute the trailing part after these nloops
	for i := (vo + nloops * vd) * loopsize; i < nrocks; i++ {
		time = dropRock(time)
	}
	// add the height the skipped loops would have added
	height += hd * (nloops - 1)
	return height
}

//////////// Common Parts code

// return a slice of dirs (0 = right, 2 = left)
func parseJets(s string) (jets []int) {
	re := regexp.MustCompile("[<>]")
	jetnames := re.FindAllString(s, -1)
	for _, jetname := range jetnames {
		if jetname == ">" {
			jets = append(jets, 0)
		} else if jetname == "<" {
			jets = append(jets, 2)
		}
	}
	return
}

// (re-)initialise the global vars that are changed during drops

func initGrid() {
	height =0
	gh = 0
	grid = make([]bool, gw)
	for i := 0; i < gw; i++ { grid[i] = true;}
	height = 0
	gh = 1
}

// drops one rock until it rests. returns new time. Uses globals.

func dropRock(time int) (t int) {
	t = time
	rock := rocks[currock % len(rocks)]
	currock++
	makeRoom(1 + height + 3 + rock.height)
    p := 3 + gw * (height + 4) // rock position, of the bottom-left corner
	VProck(rock, p, "start")
	for {
		ph := stepRock(rock, p, jets[t % len(jets)])
		VProck(rock, ph, "after jet")
		t++
		p = stepRock(rock, ph, down)
		VProck(rock, p, "after fall")
		if p == ph {				// came to rest
			for i := 0; i < len(rock.body); i++ {
				grid[p + rock.body[i]] = true
			}
			rockheight := (p / gw) - 1 + rock.height
			if rockheight > height { height = rockheight;}
			break
		}
	}
	VPgrid()
	return
}

// ensure grid has allocated for at least size height

func makeRoom(size int) {
	for gh < size  {
		grid = append(grid, emptyGridRow...)
		gh++
	}
}

// move the rock one step in the direction, return new bottomleft pos

func stepRock(rock Rock, p, dir int) int {
	delta := dirdelta[dir]
	p2 := p + delta
	for i := 0; i < len(rock.body); i++ {
		if grid[p2 + rock.body[i]] { return p;} // hit a rock
	}
	return p2
}

// pretty print grid

func VPgrid() {
	if verbose {
		fmt.Printf("=== height = %d ===\n", height)
		for y := height; y >= 0; y-- {
			for x := 0; x < gw; x++ {
				if grid[x + y*gw] {
					fmt.Printf("#")
				} else {
					fmt.Printf(".")
				}
			}
			fmt.Printf("\n")
		}
	}
}

// pretty print grid + rock

func VProck(rock Rock, p int, label string) {
	if verbose {
		fmt.Printf("=== height = %d === [%s]\n", height, label)
		g := make([]byte, len(grid))
		for i:= 0; i < len(grid); i++ {
			if grid[i] {
				g[i] = byte('#')
			} else {
				g[i] = byte('.')
			}
		}
		for i := 0; i < len(rock.body); i++ {
			g[p + rock.body[i]] = byte('@')
		}
		for y := gh - 1; y >= 0; y-- {
			fmt.Printf("%s\n", string(g[y*gw:(y+1)*gw]))
		}
	}
}

//////////// Part1 functions

//////////// Part2 functions
