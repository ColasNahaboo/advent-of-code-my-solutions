// Adventofcode 2022, d16, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 example 1651
// TEST: -1 input 1376
// TEST: example 1707
// TEST: ex1 836
// TEST: input 1933
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
	to []int
}

var verbose bool
var tempmax int					// the current best pressure found up to now
var limit int					// the maximum theoritcal limit for pressure

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
	valves := parse(lines)
	flows := make([]int, len(valves), len(valves)) // current flows
	openable := 0
	for _, v := range valves {
		openable += v.flow
	}
	VPf("Maximum potential pressure: %d x %d = %d\n", openable, 30, openable * 30)
	return explore(&valves, flows, 0, 0, 30, 0, openable, 0)
}

//////////// Part 2
func part2(lines []string) int {
	valves := parse(lines)
	flows := make([]int, len(valves), len(valves)) // current flows
	openable := 0
	for _, v := range valves {
		openable += v.flow
	}
	VPf("Maximum potential pressure: %d x %d = %d\n", openable, 26, openable * 26)
	return explore2(&valves, flows, 0, 0, 0, 26, 0, openable, 0, 0)
}

//////////// Common Parts code

var reline = regexp.MustCompile("Valve ([[:upper:]]+) has flow rate=([[:digit:]]+); tunnels? leads? to valves? ([, [:upper:]]+)")
var retos = regexp.MustCompile("[[:upper:]][[:upper:]]")

func parse(lines []string) (valves []Valve) {
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
		valves = append(valves, Valve{id:name2id[vals[1]], name:vals[1], flow:atoi(vals[2]), to:tois})
	}
	sort.Slice(valves, func(i, j int) bool { return valves[i].id < valves[j].id;})
	// sort to tunnels, to try first the higest flow valves
	for _, v := range valves {
		tov := []Valve{}
		for _, id := range v.to {
			tov = append(tov, valves[id])
		}
		sort.Slice(tov, func(i, j int) bool { return tov[i].flow > tov[j].flow; })
		for _, tovv := range tov {
			v.to = append(v.to, tovv.flow)
		}
	}
	VP(valves)
	return
}

func flowsPressure(flows []int) (sum int){
	for _, flow := range flows {
		sum += flow
	}
	return
}

//////////// Part1 functions

// flows = list of currently opened flows
// curflow = sum of the above
// id = id of the valve we are on
// time = remaining seconds until timeout
// pressure = the cumulated pressure upto now
// available = sum of flows of still closed valves
// from = if we come from a node (id) without opening a valve, dont go back there
// returns the total cumulated pressure at tiemout
// uses global tempmax

func explore(valves *[]Valve, flows []int, curflow, id, time, pressure, available, from int) (maxp int) {
	if time <= 0 {
		// time expired
		maxp = pressure
		goto DONE
	}
	if pressure + (curflow + available) * time <= tempmax {
		// there is no way we can beat tempmax now, abort this exploration
		goto DONE
	}
	if available <= 0 {
		// no remaining openable valves, stay put and compute the final pressure
		maxp = pressure + curflow * time
		goto DONE
	}
	if (*valves)[id].flow > 0 && flows[id] == 0 { // can we open current one?
		// open current valve. set from to -1 because we then can go back
		vflow := (*valves)[id].flow
		flows2 := make([]int, len(flows), len(flows))
		copy(flows2, flows)
		flows2[id] = vflow
		p := explore(valves, flows2, curflow + vflow, id, time - 1, pressure + curflow, available - vflow, -1)
		if p > maxp { maxp = p;}
	}
	// try moving to all the "to" tunnels,
	for _, to := range (*valves)[id].to {
		if to == from { continue;}
		p := explore(valves, flows, curflow, to, time - 1, pressure + curflow, available, id)
		if p > maxp { maxp = p;}
	}
DONE:
	if maxp > tempmax {
		tempmax = maxp
	}
	return
}

//////////// Part2 functions


// flows = list of currently opened flows
// curflow = sum of the above
// id = id of the valve we are on
// eid = id the elephant (El) is on
// time = remaining seconds until timeout
// pressure = the cumulated pressure upto now
// available = sum of flows of still closed valves
// from = if we come from a node (id) without opening a valve, dont go back there
// efrom = same as from but for location of elephant
// returns the total cumulated pressure at tiemout
// uses global tempmax


func explore2(valves *[]Valve, flows []int, curflow, id, eid, time, pressure, available, from, efrom int) (maxp int) {
	if time <= 0 {
		// time expired
		maxp = pressure
		goto DONE
	}
	if pressure + (curflow + available) * time <= tempmax {
		// there is no way we can beat tempmax now, abort this exploration
		goto DONE
	}
	if available <= 0 {
		// no remaining openable valves, stay put and compute the final pressure
		maxp = pressure + curflow * time
		goto DONE
	}
	if (*valves)[id].flow > 0 && flows[id] == 0 { // can we open current one?
		// open current valve. set from to -1 because we then can go back
		vflow := (*valves)[id].flow
		flows2 := make([]int, len(flows), len(flows))
		copy(flows2, flows)
		flows2[id] = vflow
		if (*valves)[eid].flow > 0 && flows2[eid] == 0 { // can El open current one?
			evflow := (*valves)[eid].flow
			eflows2 := make([]int, len(flows), len(flows))
			copy(eflows2, flows2)
			eflows2[eid] = evflow
			// both me and El opened their valves
			p := explore2(valves, eflows2, curflow + vflow + evflow, id, eid, time - 1, pressure + curflow, available - vflow - evflow, -1, -1)
			if p > maxp { maxp = p;}
		}
		// I opened my valve, and try moving El to all the "to" tunnels,
		for _, eto := range (*valves)[eid].to {
			if eto == efrom { continue;}
			p := explore2(valves, flows2, curflow + vflow, id, eto, time - 1, pressure + curflow, available - vflow, -1, eid)
			if p > maxp { maxp = p;}
		}
	}
	// try moving to all the "to" tunnels,
	for _, to := range (*valves)[id].to {
		if to == from { continue;}
		if (*valves)[eid].flow > 0 && flows[eid] == 0 { // can El open current one?
			evflow := (*valves)[eid].flow
			eflows2 := make([]int, len(flows), len(flows))
			copy(eflows2, flows)
			eflows2[eid] = evflow
			// I move and El opens
			p := explore2(valves, eflows2, curflow + evflow, to, eid, time - 1, pressure + curflow, available - evflow, -1, eid)
			if p > maxp { maxp = p;}
		}
		// try moving El to all the "to" tunnels,
		for _, eto := range (*valves)[eid].to {
			if eto == efrom { continue;}
			// I and El both move
			p := explore2(valves, flows, curflow, to, eto, time - 1, pressure + curflow, available, id, eid)
			if p > maxp { maxp = p;}
		}
	}
DONE:
	if maxp > tempmax {
		tempmax = maxp
	}
	return
}
