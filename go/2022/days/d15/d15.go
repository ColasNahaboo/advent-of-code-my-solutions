// Adventofcode 2022, d15, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 example 26
// TEST: -1 input 5147333
// TEST: example 56000011
// TEST: input 13734006908372
package main

import (
	"flag"
	"fmt"
	"log"
	"sort"
	// "regexp"
)

type Sensor struct {
	x, y int
	cbx, cby int
	cbd int						// closest manhattan distance S/B
}

var verbose bool
var row int					// the row we examine in part1
var max int					// the max coords for part2
const maxint = 8888888888888888888 // easily identifiable in debug

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	rowFlag := flag.Int("c", 2000000, "the y-coord of the \"row\" line for part1")
	maxFlag := flag.Int("m", 4000000, "the max of coords for part2")
	flag.Parse()
	verbose = *verboseFlag
	infile := "input.txt"
	row = *rowFlag
	max = *maxFlag
	if flag.NArg() > 0 {
		infile = flag.Arg(0)
		if infile == "example.txt" { // default values for the in-text example
			row = 10
			max = 20
		}
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

func part1(lines []string) (excluded int) {
	all := parse(lines)
	sensors := inRange(all, row)
	excluded = - onRow(sensors, row)
	VPf("%d relevant sensors out of %d wrt row %d:\n", len(sensors), len(all), row)
	VP(sensors)
	x1, x2 := xRangeRow(sensors, row)
	// we use a brute force approach of examining each x
	// it is ok for the one row of part1, but we will need to be smarter in part2
	for x := x1; x <= x2; x++ {
		for _, s := range sensors {
			if manDist(x, row, s.x, s.y) <= s.cbd {
				VPf("[%d,%d] inside (%d) range (%d) of sensor [%d,%d]\n", x, row, manDist(x, row, s.x, s.y), s.cbd, s.x, s.y)
				excluded++
				break
			}
		}
	}
	return
}

//////////// Part 2
func part2(lines []string) (tf int) {
	all := parse(lines)
	for row := 0; row <= max; row++ {
		tf = rangeRow(all, row)
		if tf > 0 {
			return
		}
	}
	return 0
}
	
//////////// Common Parts code

func parse(lines []string) (sensors []Sensor) {
	var s Sensor
	for lineno, line := range lines {
		n, err := fmt.Sscanf(line, "Sensor at x=%d, y=%d: closest beacon is at x=%d, y=%d", &s.x, &s.y, &s.cbx, &s.cby)
		if err != nil || n != 4 {
			log.Fatalf("Syntax error line %d: %s\n", lineno + 1, line)
		}
		s.cbd = s.sbDist()
		sensors = append(sensors, s)
		VP(s)
	}
	return
}

// return sensors that are in range of row
func inRange(sensors []Sensor, row int) (sir []Sensor) {
	for _, s := range sensors {
		if absdiff(s.y, row) <= s.cbd {
			sir = append(sir, s)
		}
	}
	return
}

// return the number of sensor beacons on the row
func onRow(sensors []Sensor, row int) int {
	xs := make(map[int]bool)
	for _,s := range sensors {
		if s.cby == row {
			if xs[s.cbx] { continue;}
			VPf("Sensor Beacon [%d,%d] is on row %d\n", s.cbx, row, row)
			xs[s.cbx] = true			
		}
	}
	return len(xs)
}


func manDist(x1, y1, x2, y2 int) int {
	if x1 > x2 {
		if y1 > y2 {
			return x1 - x2 + y1 - y2
		} else {
			return x1 - x2 - y1 + y2
		}
	} else {
		if y1 > y2 {
			return -x1 + x2 + y1 - y2
		} else {
			return -x1 + x2 - y1 + y2
		}
	}
}

func (s Sensor) sbDist() int {
	return manDist(s.x, s.y, s.cbx, s.cby)
}

func absdiff(x, y int) int {
	if y > x {
		return y - x
	} else {
		return x - y
	}
}

//////////// Part1 functions

// return the portion [x1,x2] of row that is in range of any of the sensors
func xRangeRow(sensors []Sensor, row int) (x1, x2 int) {
	x1 = maxint
	for _,s := range sensors {
		dx := s.cbd - absdiff(s.y, row) // range of covered x on row from sensor
		if s.x - dx < x1 { x1 = s.x - dx;}
		if s.x + dx > x2 { x2 = s.x + dx;}
	}
	return
}

//////////// Part2 functions

// builds the sequence "ex" of x segments [x1,x2] of row in range of any of the sensors
// between 0 and max.
// ex is then sorted
// then for each segment in ex, we "push" the available free position to the right
// but if a segment is to the right of x, we then know we have a gap, as they are sorted
// return the Tuning Frequency of the hidden beacon if found, or zero if not

func rangeRow(sensors []Sensor, row int) int {
	ex := [][2]int{}
	for _,s := range sensors {
		dx := s.cbd - absdiff(s.y, row) // range of covered x on row from sensor
		if dx < 0 { continue;}
		x1 := s.x - dx
		x2 := s.x + dx
		if x2 < 0 || x1 > max {	// outside the search region, ignore
			continue
		}
		ex = append(ex, [2]int{x1, x2})
	}

	sort.Slice(ex, func(i, j int) bool {
		return ex[i][0] < ex[j][0] || (ex[i][0] == ex[j][0] && ex[i][1] < ex[j][1])
	})

	x := 0						// smallest x possible on row
	for i := 0; i < len(ex); i++ {
		if ex[i][0] <= x {		// segment potentially covering x
			if ex[i][1] < x { 	// segment left of x, misses it, ignore
				continue
			}
			x = ex[i][1] + 1	// segment covers x -> pushes x to its right
			if x > max {		// x pushed outside of allowed range -> row done, next!
				return 0
			}
		} else {
			// found a gap, and the the possible position of the hidden beacon!
			VPf("Found! [%d,%d], gap in row %d: %v\n", x, row, row, ex)
			return x * 4000000 + row
		}
	}
	// should not happen, as it would mean all the row is acceptable!
	return row
}

