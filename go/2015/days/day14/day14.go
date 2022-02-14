// Adventofcode 2015, day14, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 input 2660
// TEST: input 1256
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	// "regexp"
)

type Reindeer struct {
	// reindeer static properties
	name  string
	speed int // km/s
	run   int // duration of a run state
	rest  int // duratioon of a rest state
	// reindeer dynamic status during a run, at time t
	score int // number of points gained
	pos   int // position in km
	state int // 0 = rest, 1 = run
	end   int // at what time next state starts
}

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
	reindeers := ParseReindeers(lines)
	return WinnerDistance(reindeers, 2503)
}

//////////// Part 2
func Part2(lines []string) int {
	reindeers := ParseReindeers(lines)
	return WinnerPoints(reindeers, 2503)
}

//////////// Common Parts code

func ParseReindeers(lines []string) []Reindeer {
	reindeers := make([]Reindeer, 0)
	for _, line := range lines {
		// Vixen can fly 19 km/s for 7 seconds, but then must rest for 124 seconds.
		// 0     1   2   3  4    5   6 7        8   9    10   11   12  13  14
		tokens := strings.Split(line, " ")
		if len(tokens) != 15 {
			if len(tokens) == 0 { // skip empty lines
				continue
			} else { // crude error check
				log.Fatalf("Parse error (%v tokens) on line: %v", len(tokens), line)
			}
		}
		r := new(Reindeer)
		r.name = tokens[0]
		r.speed = Atoi(tokens[3])
		r.run = Atoi(tokens[6])
		r.rest = Atoi(tokens[13])
		r.state = 1 // reindeers start by running
		r.end = r.run
		reindeers = append(reindeers, *r)
	}
	return reindeers
}

////// Part1 funcs

func WinnerDistance(reindeers []Reindeer, time int) int {
	max := 0
	for _, r := range reindeers {
		d := ReindeerDistance(r, time)
		if d > max {
			max = d
		}
	}
	return max
}

func ReindeerDistance(r Reindeer, time int) int {
	cycles := time / (r.run + r.rest)
	d := cycles * r.run * r.speed // number of full cycles: run+rest
	remain := time - (cycles * (r.run + r.rest))
	if remain >= r.run { // stopped in rest
		d += r.run * r.speed
	} else { // stopped mid-run
		d += remain * r.speed
	}
	return d
}

///// Part2 funcs

func WinnerPoints(reindeers []Reindeer, time int) int {
	score := 0
	for t := 0; t < time; t++ {
		WinnerScoreIncr(reindeers, t)
	}
	for _, r := range reindeers {
		if r.score > score {
			score = r.score
		}
	}
	return score
}

// run the simulation for one second starting a time t,
// updating score and position of the Reindeer objects
func WinnerScoreIncr(reindeers []Reindeer, t int) {
	max := 0
	// advance reindeers
	for n := range reindeers {
		r := &reindeers[n] // to update them, and not copies
		if t >= r.end {    // change state
			if r.state == 0 {
				r.state = 1
				r.end = t + r.run
			} else {
				r.state = 0
				r.end = t + r.rest
			}
		}
		r.pos = r.pos + r.state*r.speed
		if r.pos > max {
			max = r.pos
		}
	}
	// award a point to all currently in the lead
	for n := range reindeers {
		r := &reindeers[n]
		if r.pos == max {
			r.score = r.score + 1
		}
	}
}

//////////// Generic code

// useful in tests to feed Part1 & Part2 with a simple string (with newlines)
func StringToLines(s string) (lines []string) {
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	err := scanner.Err()
	if err != nil {
		log.Fatal(err)
	}
	return
}

// read the input file into a string array for feeding Parts
func FileToLines(filePath string) (lines []string) {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}

	return
}

// simplified functions to not bother with error handling. Just abort.

func Atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal(err)
	}
	return i
}

// for completeness
func Itoa(i int) string {
	return strconv.Itoa(i)
}
