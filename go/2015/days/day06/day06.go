// Adventofcode 2015, day06, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 input 543903
// TEST: input
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Grid [1000][1000]int

func main() {
	partOne := flag.Bool("1", false, "run part one code, instead of part 2 (default)")
	flag.Parse()
	infile := "input.txt"
	if flag.NArg() > 0 {
		infile = flag.Arg(0)
	}
	lines := FileToLines(infile)

	var result int
	if *partOne {
		fmt.Println("Running Part1")
		result = Part1(lines)
	} else {
		fmt.Println("Running Part2")
		result = Part2(lines)
	}
	fmt.Println(result)
}

//////////// Part 1

func Part1(lines []string) int {
	var grid Grid
	on := 0
	reInst := regexp.MustCompile("(turn on|turn off|toggle)[[:space:]]+([0-9]+),([0-9]+)[[:space:]]+through[[:space:]]+([0-9]+),([0-9]+)")
	for _, line := range lines {
		s := reInst.FindStringSubmatch(line)
		if s != nil {
			x1, err1 := strconv.Atoi(s[2])
			y1, err2 := strconv.Atoi(s[3])
			x2, err3 := strconv.Atoi(s[4])
			y2, err4 := strconv.Atoi(s[5])
			if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
				os.Exit(2)
			}
			switch {
			case s[1] == "turn on":
				on = grid.TurnOn(on, x1, y1, x2, y2)
			case s[1] == "turn off":
				on = grid.TurnOff(on, x1, y1, x2, y2)
			case s[1] == "toggle":
				on = grid.Toggle(on, x1, y1, x2, y2)
			default:
				fmt.Printf("== Unknown instruction code: %v\n", line)
			}
		} else {
			fmt.Printf("== Cannot parse line: %v\n", line)
		}
	}
	return on
}

func (grid *Grid) TurnOn(on int, x1, y1, x2, y2 int) int {
	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			if grid[x][y] == 0 {
				grid[x][y] = 1
				on++
			}
		}
	}
	return on
}

func (grid *Grid) TurnOff(on int, x1, y1, x2, y2 int) int {
	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			if grid[x][y] != 0 {
				grid[x][y] = 0
				on--
			}
		}
	}
	return on
}

func (grid *Grid) Toggle(on int, x1, y1, x2, y2 int) int {
	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			if grid[x][y] != 0 {
				grid[x][y] = 0
				on--
			} else {
				grid[x][y] = 1
				on++
			}
		}
	}
	return on
}

//////////// Part 2
func Part2(lines []string) int {
	var grid Grid
	bright := 0
	reInst := regexp.MustCompile("(turn on|turn off|toggle)[[:space:]]+([0-9]+),([0-9]+)[[:space:]]+through[[:space:]]+([0-9]+),([0-9]+)")
	for _, line := range lines {
		if s := reInst.FindStringSubmatch(line); s != nil {
			x1, err1 := strconv.Atoi(s[2])
			y1, err2 := strconv.Atoi(s[3])
			x2, err3 := strconv.Atoi(s[4])
			y2, err4 := strconv.Atoi(s[5])
			if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
				os.Exit(2)
			}
			switch {
			case s[1] == "turn on":
				bright = grid.TurnOn2(bright, x1, y1, x2, y2)
			case s[1] == "turn off":
				bright = grid.TurnOff2(bright, x1, y1, x2, y2)
			case s[1] == "toggle":
				bright = grid.Toggle2(bright, x1, y1, x2, y2)
			default:
				fmt.Printf("== Unknown instruction code: %v\n", line)
			}
		} else {
			fmt.Printf("== Cannot parse line: %v\n", line)
		}
	}
	return bright
}

func (grid *Grid) TurnOn2(bright int, x1, y1, x2, y2 int) int {
	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			grid[x][y]++
			bright++
		}
	}
	return bright
}

func (grid *Grid) TurnOff2(bright int, x1, y1, x2, y2 int) int {
	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			if grid[x][y] > 0 {
				grid[x][y]--
				bright--
			}
		}
	}
	return bright
}

func (grid *Grid) Toggle2(bright int, x1, y1, x2, y2 int) int {
	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			grid[x][y] += 2
			bright += 2
		}
	}
	return bright
}

//////////// Common Parts code

//////////// Generic code
func StringToLines(s string) (lines []string) {
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	err := scanner.Err()
	if err != nil {
		os.Exit(1)
	}
	return
}

func FileToLines(filePath string) (lines []string) {
	f, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	// optionally, resize scanner's capacity for lines over 64K (65536)
	const maxCapacity = 1000000 // your required line length
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)
	// end optional
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	err = scanner.Err()
	if err != nil {
		os.Exit(1)
	}

	return
}
