// Adventofcode 2022, d16, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 example 1651
// TEST: -1 ex1 28
// TEST: -1 input 1376
// TEST: example 1707
// TEST: ex1 24
// TEST: input 1933

// This is a smarter implementation by first finding all the routes between all valves
// by a Floyd-Warshall algorythm , and then exploring only the routes leading
// to an openable valve.
// I keep the naive implementation for reference in d16-naive.go that explores
// everything step by step, but part2 took more than 1 hour to run, for 20s here.
// d16-old.go go directly to the destinations
// I now directly add to the flow the total future flow of opened valves until end

package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"
	"sort"
)

type Valve struct {
	id int
	name string
	flow int
	to []int					// where can I go? array of openable valves ids
	next []int					// neighbours, one step distance
}

var valves []Valve				// all the valves. Global variable for convenience
var nopenable int				// how many have non-0 flow if opened
var nvalves int					// size of the slice valves
var maxpressure int				// the current best pressure found up to now
var limit int					// the maximum theoretical limit for pressure
var dist [][]int				// shortest distance between valves
var next [][]int				// to reconstruct the path of Floyd-Warshall

var verbose bool
const maxint = 888888888888888888 // easily identifiable in debug

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	flag.Parse()
	verbose = *verboseFlag
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

func part1(lines []string) int {
	parse(lines)
	flows := make([]int, nvalves, nvalves) // current flows
	openable := 0
	for _, v := range valves {
		openable += v.flow
	}
	VPf("Maximum potential pressure: %d x %d = %d\n", openable, 30, openable * 30)
	explore(flows, 0, 30, 0, openable)
	return maxpressure
}

//////////// Part 2
func part2(lines []string) int {
	parse(lines)
	flows := make([]int, nvalves, nvalves) // current flows
	openable := 0
	closed := 0					// closed but with potential flow
	for _, v := range valves {
		openable += v.flow
		if v.flow > 0 { closed++;}
	}
	VPf("Maximum potential pressure: %d x %d = %d\n", openable, 26, openable * 26)
	explore2(flows, 0, 0, 26, 26, 0, openable, closed)
	return maxpressure
}

//////////// Common Parts code

var reline = regexp.MustCompile("Valve ([[:upper:]]+) has flow rate=([[:digit:]]+); tunnels? leads? to valves? ([, [:upper:]]+)")
var retos = regexp.MustCompile("[[:upper:]][[:upper:]]")

func parse(lines []string) {
	var name string
	names := []string{}
	for _, line := range lines {
		fmt.Sscanf(line, "Valve %2s ", &name)
		names = append(names, name)
	}
	// we sort by names, so that AA is id 0, rest is in order, easier debug
	sort.Slice(names, func(i, j int) bool { return names[i] < names[j];})
	name2id := make(map[string]int, 0)
	for id, name := range names {
		name2id[name] = id
	}
	VPf("name2id: %v\n", name2id)
	for _, line := range lines {
		vals := reline.FindStringSubmatch(line)
		if len(vals) == 0 {
			log.Fatalf("Syntax error on line: %s\n", line)
		}
		tos := retos.FindAllString(vals[3], -1)
		tois := []int{}
		for _, name := range tos {
			tois = append(tois, name2id[name])
		}
		valves = append(valves, Valve{id:name2id[vals[1]], name:vals[1], flow:atoi(vals[2]), next:tois})
		if vals[2] != "0" { nopenable++;}
	}
	nvalves = len(valves)
	sort.Slice(valves, func(i, j int) bool { return valves[i].id < valves[j].id;})
	computePaths()
	VP(valves)
}

func flowsPressure(flows []int) (sum int){
	for _, flow := range flows {
		sum += flow
	}
	return
}

// compute the to fields of all valves, modifies valves in place
func computePaths() {
	// initialize 2D slices
	l := len(valves)
	dist = make([][]int, l, l)
	next = make([][]int, l, l)
	for i := range dist {
		dist[i] = make([]int, l, l)
		next[i] = make([]int, l, l)
		for j :=0; j < l; j++ {
			dist[i][j] = maxint
		}
	}
	// Floyd-Warshall https://en.wikipedia.org/wiki/Floyd%E2%80%93Warshall_algorithm
	// compute shortests distances (dist) and way to recreate paths (next)
	for v := range valves {
		dist[v][v] = 0
		next[v][v] = v
		for _, n := range valves[v].next {
			dist[v][n] = 1
			next[v][n] = n
		}
	}
	for k := range valves {
		for i := range valves {
			for j := range valves {
				if dist[i][j] > dist[i][k] + dist[k][j] {
					dist[i][j] = dist[i][k] + dist[k][j]
					next[i][j] = next[i][k]
				}
			}
		}
	}
	// now, each to field is the reverse path to all openable valves
	for v := range valves {
		tos := make([]int, 0)
		for n := range valves {
			if n == v || valves[n].flow == 0 || dist[v][n] >= maxint {
				continue		// skip self. or non reachable, or non openable
			}
			tos = append(tos, n)
		}
		// heuristic: sort them by biggest and closest flow first
		sort.Slice(tos, func(i, j int) bool {
			return dist[v][tos[i]] - valves[tos[i]].flow < dist[v][tos[j]] - valves[tos[j]].flow
		})
		valves[v].to = tos
	}
}

// build a path from the Floyd-Warshall "next" table
// includes start but not end
func pathFW(i, j int) (path []int) {
	path = []int{i}
	for i != j {
		i = next[i][j]
		path = append(path, i)
	}
	return
}

// build a reverse path from the Floyd-Warshall "next" table
// includes end but not start
func pathFWreverse(i, j int) (path []int) {
	path = []int{}
	for i != j {
		i = next[i][j]
		path = append([]int{i}, path...) // prepend i
	}
	return
}

// not member of slice
func notIn(s *[]int, n int) bool {
	for i := range *s {
		if n == (*s)[i] { return false;}
	}
	return true
}

//////////// Part1 functions

// flows = list of currently opened flows
// id = id of the valve we are on
// time = remaining seconds until timeout
// pressure = the cumulated pressure reached a the end with this valves state
// available = sum of flows of still closed valves
// uses and update the global maxpressure, max pressure

func explore(flows []int, id, time, pressure, available int) {
	if pressure > maxpressure {
		maxpressure = pressure
	}
	// first, see if we are at the end
	if time <= 0 {				// time expired
		return
	}
	if pressure + available * time <= maxpressure { // impossible to reach max
		return
	}
	if available <= 0 {
		// no remaining openable valves, stay put and compute the final pressure
		return
	}
	// recurse! Go straight to each openable valves, open it, add total flow
	for _, to := range valves[id].to {
		if flows[to] != 0 { continue;}		// already opened
		if time < dist[id][to] { continue;}	  // not reachable in time
		// pre-compute state and pressure for this branch
		toflow := valves[to].flow
		flows2 := make([]int, len(flows), len(flows))
		copy(flows2, flows)
		flows2[to] = toflow
		// time of opening to is travel time dist[id][to] + 1s to open
		explore(flows2, to, time - dist[id][to] - 1, pressure + toflow * (time - dist[id][to] - 1), available - toflow)
	}
}

//////////// Part2 functions

// flows = list of currently opened flows
// id = id of the valve we are on
// eid = id the elephant (El) is on
// time = remaining seconds until timeout for me
// etime = remaining seconds until timeout for El
// pressure = the cumulated pressure reached a the end with this valves state
// available = sum of flows of still closed valves
// uses and update the global maxpressure, max pressure

func explore2(flows []int, id, eid, time, etime, pressure, available, closed int) {
	// first, check if we at at the end of the exploration
	if pressure > maxpressure {
		maxpressure = pressure
	}
	// no remaining openable valves, stay put and compute the final pressure
	if available <= 0 {
		return
	}
	if time <= 0 {
		if etime <= 0 {			// no time left for any of us, terminate
			return
		} else {				// explore only El: fall back to part1 code
			explore(flows, eid, etime, pressure, available)
			return
		}
	} else {
		if etime <= 0 {			// explore only me: fall back to part1 code
			explore(flows, id, time, pressure, available)
			return
		}						// else go on
	}
	if pressure + available * (time + etime) <= maxpressure { // impossible to reach max
		return
	}
	if closed == 1 {
		// special case: only one valve is openable. The closest (spacetime) one go open it.
		// the general code below do not work in this case, nor part1 exploring.
		// find the only one still openable: to1
		var to1 int				
		for _, to := range valves[id].to {
			if flows[to] == 0 { to1 = to;}
		}
		toflow := valves[to1].flow
		// choose who will be the more efficient to open it: me or el?
		opened_time := time - dist[id][to1] - 1
		if etime - dist[eid][to1] - 1 > opened_time {
			opened_time = etime - dist[eid][to1] - 1
		}
		pressure += toflow * opened_time
		if pressure > maxpressure {
			maxpressure = pressure
		}
		return
	}

	// recurse both! Go straight to each openable valves, open it, add total flow
	for _, to := range valves[id].to {
		if flows[to] != 0 { continue;}		// already opened
		if time < dist[id][to] { continue;}	// not reachable in time
		for _, eto := range valves[eid].to {
			if flows[eto] != 0 { continue;}		// already opened
			if etime < dist[eid][eto] { continue;} // not reachable in time
			if eto == to { continue;}
			// move both
			toflow := valves[to].flow
			etoflow := valves[eto].flow
			flows2 := make([]int, len(flows), len(flows))
			copy(flows2, flows)
			flows2[to] = toflow
			flows2[eto] = etoflow
			explore2(flows2, to, eto,
				time - dist[id][to] - 1, etime - dist[eid][eto] - 1,
				pressure + toflow * (time - dist[id][to] - 1) + etoflow * (etime - dist[eid][eto] - 1),
				available - toflow - etoflow, closed - 2)
		}
	}
}
