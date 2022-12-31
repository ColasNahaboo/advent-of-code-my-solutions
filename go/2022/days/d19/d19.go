// Adventofcode 2022, d19, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 example 33
// TEST: -1 input 1404
// TEST: input 5880
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
const glyphs = "ocbg_"			// letter of robot id + nothing
const maxint = 888888888888888888 // easily identifiable in debug

var blueprints []Blueprint
// globals for heuristics to cut short dead ends
var totaltime int
var maxgeodes int				// temporary current max reached

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
		maxgeodes = 0
		VPf("Simulating Blueprint#%d for %ds: %v\n", id+1, totaltime, blueprints[id])
		run(id, totaltime, -1, initState, "")
		ql := (id+1) * maxgeodes
		now := time.Now().UnixMilli()
		duration := now - lastTime
		lastTime = now
		VPf("BP#%2d => %2d geodes, %3d quality level, duration: %6.3f\n", id+1, maxgeodes, ql, float64(duration)/1000)
		qls += ql
	}
	fmt.Printf("Total duration: %.3f\n", float64(time.Now().UnixMilli() - startTime)/1000)
	return
}

//////////// Part 2
func part2(lines []string) int {
	parse(lines)
	initState := State{robots: [4]int{1,0,0,0}}
	lastTime := time.Now().UnixMilli()
	startTime := lastTime
	totaltime = 32
	res := 1
	for id := range blueprints {
		maxgeodes = 0
		VPf("Simulating Blueprint#%d for %ds: %v\n", id+1, totaltime, blueprints[id])
		run(id, totaltime, -1, initState, "")
		now := time.Now().UnixMilli()
		duration := now - lastTime
		lastTime = now
		VPf("BP#%2d => %2d geodes, duration: %6.3f\n", id+1, maxgeodes, float64(duration)/1000)
		fmt.Printf("BP#%d: %d\n", id, maxgeodes)
		res *= maxgeodes
		if id >= 2 {
			return res
		}
	}
	fmt.Printf("Total duration: %.3f\n", float64(time.Now().UnixMilli() - startTime)/1000)
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
		if id != lineno + 1 && verbose {
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

// explore the possibilities.
// trace is here for debuuging only

func run(id, time, fact int, s State, trace string) {
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
	time--	  // ok, consumed our time
	VPf("%s@%d %v\n", trace, time, s)
	// out of time, cannot start a new cycle, stop successfully here.
	if time <= 0 {
		if s.mats[geode] > maxgeodes {
			maxgeodes = s.mats[geode]
			VPf("%d geodes for %s@%d (timeout)\n", s.mats[geode], trace, time)
		}
		return
	}
	// if we could not possibly beat the current max, abort this run
	rg := s.robots[3]
	mg := s.mats[3]
	for t := 0; t < time; t++ {	// suppose we from now build a geode robot per turn
		mg += rg
		rg++
	}
	if mg <= maxgeodes {
		VPf("%s@%d %v ==> no time to build geodes\n", trace, time, s)
		return
	}
	
	// Now explore each possible factory action:
	// we do not limit to just the next step(second), we pursue the goal

	// first, no robots: how many geodes if we just create nothing until end
	newgeodes := s.mats[geode] + time * s.robots[geode]
	if newgeodes > maxgeodes {
		maxgeodes = newgeodes
		VPf("%s@%d ->%d_ ### %d geodes ###\n", trace, time, time, maxgeodes)
	}
	
	// then DFS-explore the 4 branches where we decide to do a robot
	// in reverse as we try to build the "better" robots first
NEXT_ROBOT:
	for r := 3; r >= 0; r-- {
		// first check we are not already at max capacity for this type
		if s.robots[r] >= blueprints[id].maxrobots[r] {
			continue
		}
		// how much steps we will have enough of all mats to build this robot?
		wait := 0
		for mat := 0; mat < 4; mat++ {
			if blueprints[id].cost[r][mat] > s.mats[mat] { // currently missing
				if s.robots[mat] == 0 { // no bot -> never, abort!
					continue NEXT_ROBOT
				}
				// + s.robots[mat]-1 ==> round up to get time needed
				mwait := (blueprints[id].cost[r][mat] - s.mats[mat] + s.robots[mat]-1) / s.robots[mat]
				if mwait > wait {
					if mwait >= time {
						continue NEXT_ROBOT // not enough time left
					}
					wait = mwait
				}
			}
		}
		if wait > 0 {
			s2 := State{[4]int{s.mats[0] + wait*s.robots[0], s.mats[1] + wait*s.robots[1], s.mats[2] + wait*s.robots[2], s.mats[3] + wait*s.robots[3]}, s.robots}
			VPf("%s@%d ->%d%s (@+%d)\n", trace, time, wait, glyphs[r:r+1], time-wait)
			run(id, time - wait, r, s2, trace + itoa(wait) + glyphs[r:r+1])

		} else {
			VPf("%s@%d -> %s\n", trace, time, glyphs[r:r+1])
			run(id, time, r, s, trace + glyphs[r:r+1])
		}
	}
	return
}

//////////// Part1 functions

//////////// Part2 functions
