//go:build exclude_this_is_just_for_reference
// This is a first naive implementation, exploring step by step. Too slow.

// Adventofcode 2022, d19, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 example 33
// TEST: -1 input
// TEST: example
// TEST: input
package main

import (
	"flag"
	"fmt"
	"log"
	"time"
	// "regexp"
)

type Blueprint struct {
	id int
	cost [4][4]int				// cost[i][j] = cost in material j to make robot i
	maxrobots [4]int			// optim: max robots for the factory capacity
}

// whate defines the state of the system?
// we do not include time and the robot being made, as they are more transient
type State struct {
	mats [4]int					// the materials in stock
	robots [4]int				// the fleet of active robots
}

const ore = 0					// abbrev in code: o
const clay = 1					// abbrev in code: c
const obsidian = 2				// abbrev in code: b
const geode = 3					// abbrev in code: g
const glyphs = ".OCBG"			// letter of robot id + 1 for debug
const maxint = 888888888888888888 // easily identifiable in debug

var blueprints []Blueprint
// globals for heuristics to cut short dead ends
var totaltime int
var totalgeodes int				// temporary current max reached

var verbose bool

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	timeFlag := flag.Int("t", 24, "time: seconds to run")
	flag.Parse()
	verbose = *verboseFlag
	totaltime = *timeFlag
	infile := "input.txt"
	if flag.NArg() > 0 {
		infile = flag.Arg(0)
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

func part1(lines []string) (qls int) {
	parse(lines)
	initState := State{robots: [4]int{1,0,0,0}}
	lastTime := time.Now().UnixMilli()
	startTime := lastTime
	for id := range blueprints {
		totalgeodes = 0
		VPf("Simulating Blueprint#%d for %ds: %v\n", id+1, totaltime, blueprints[id])
		geodes := run(id, totaltime, -1, initState)
		ql := (id+1) * geodes
		now := time.Now().UnixMilli()
		duration := now - lastTime
		lastTime = now
		fmt.Printf("BP#%2d => %2d geodes, %3d quality level, duration: %6.3f\n", id+1, geodes, ql, float64(duration)/1000)
		qls += ql
	}
	fmt.Printf("Total duration: %.3f\n", float64(time.Now().UnixMilli() - startTime)/1000)
	return
}

//////////// Part 2
func part2(lines []string) int {
	for _, line := range lines {
		fmt.Println(line)
	}
	return 0
}

//////////// Common Parts code

func parse(lines []string) {
	for lineno, line := range lines {
		var id, oroc, croc, broc, brcc, groc, grbc int // groc = Geode Robot Ore Cost
		n, err := fmt.Sscanf(line, "Blueprint %d: Each ore robot costs %d ore. Each clay robot costs %d ore. Each obsidian robot costs %d ore and %d clay. Each geode robot costs %d ore and %d obsidian", &id, &oroc, &croc, &broc, &brcc, &groc, &grbc)
		if err != nil || n != 7 {
			log.Fatalf("Syntax error line %d: %s\n", lineno, line)
		}
		if id != lineno + 1 {
			log.Printf("Warning line %d: expecting Blueprint #%d, got %v\n", lineno+1, id, line)
		}
		blueprints = append(blueprints, Blueprint{id:id, cost:[4][4]int{
			[4]int{oroc, 0, 0, 0},
			[4]int{croc, 0, 0, 0},
			[4]int{broc, brcc, 0, 0},
			[4]int{groc, 0, grbc, 0},
		}})
	}
	// some optimisations: max robots per type, more exceeds the factory capabilities
	for id := 0; id < len(blueprints); id++ {
		for mat := 0; mat < 4; mat++ {
			for r := 0; r < 4; r++ {
				if blueprints[id].cost[r][mat] > blueprints[id].maxrobots[mat] {
					blueprints[id].maxrobots[mat] = blueprints[id].cost[r][mat]
				}
			}
			if blueprints[id].maxrobots[mat] == 0 {
				blueprints[id].maxrobots[mat] = maxint
			}
		}
	}
}

func run(id, time, fact int, s State) (maxgeodes int) {
	// run mining robots
	for mat := 0; mat < 4; mat++ {
		s.mats[mat] += s.robots[mat]
	}
	// run the factory to finish building its started robot, if any
	if fact >= 0 {
		for mat := 0; mat < 4; mat++ { // consume building costs
			s.mats[mat] -= blueprints[id].cost[fact][mat]
		}
		s.robots[fact]++
	}
	// out of time, cannot start a new cycle, stop successfully here.
	if time <= 1 {
		VPf("%d geodes for %v @%d\n", s.mats[geode], s, time)
		if s.mats[geode] > maxgeodes { maxgeodes = s.mats[geode];}
		return
	}
	// if we could not possibly beat the current max, abort this run
	rg := s.robots[3]
	mg := s.mats[3]
	if rg == 0 && ! (time > blueprints[id].cost[3][2] - s.mats[2]) {
		return					// no time remaining to build at least a geode robot
	}
	for t := 0; t < time; t++ {	// suppose we from now build a geode robot per turn
		mg += rg
		rg++
	}
	if mg <= totalgeodes {
		return
	}
	
	// Now explore each possible factory action:
	// try to be smart and avoid exploring dead ends
	// we try to build the "better" robots first
	todo := []int{}
	for _, r := range [5]int{3, -1, 2, 1, 0} {
		ok := true
		// no need to have more bots than mats than needed to build one robot,
		// as factory can only build one at a time
		if r >=0 && s.robots[r] >= blueprints[id].maxrobots[r] {
			continue
		}
		// check we have enough mats to build robot
		if r >= 0 {
			for mat := 0; mat < 4; mat++ {
				if blueprints[id].cost[r][mat] > s.mats[mat] { 
					ok = false
					break
				}
			}
		}
		if ok {
			todo = append(todo, r)
		}
	}

	// recurse to explore further the possibilities
	for _, r := range todo {
		newgeodes := run(id, time - 1, r, s)
		if newgeodes > maxgeodes { maxgeodes = newgeodes;}
	}
	return
}

//////////// Part1 functions

//////////// Part2 functions
