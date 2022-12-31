//go:build exclude
// Adventofcode 2022, d16, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 example 1651
// TEST: -1 ex1 931
// TEST: -1 input 1376
// TEST: example 1707
// TEST: ex1 836
// TEST: input 1933

// This is a smarter implementation by first finding all the routes between all valves
// by a Floyd-Warshall algorythm , and then exploring only the routes leading
// to an openable valve.
// I keep the naive implementation for reference in d16.go-naive that explores
// everything step by step, but part2 took more than 1 hour to run, for 20s here.

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
	to [][]int					// what can I do? arrays of ids [dest,..,+3,+2,+1]
	next []int					// neighbours, one step distance
}

var valves []Valve				// all the valves. Global variable for convenience
var nopenable int				// how many have non-0 flow if opened
var nvalves int					// size of the slice valves
var tempmax int					// the current best pressure found up to now
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
	return explore(flows, 0, 0, 30, 0, openable, nil)
}

//////////// Part 2
func part2(lines []string) int {
	parse(lines)
	flows := make([]int, nvalves, nvalves) // current flows
	openable := 0
	for _, v := range valves {
		openable += v.flow
	}
	VPf("Maximum potential pressure: %d x %d = %d\n", openable, 26, openable * 26)
	return explore2(flows, 0, 0, 0, 26, 0, openable, nil, nil,)
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
		tos := make([][]int, 0)
		for n := range valves {
			if n == v || valves[n].flow == 0 || dist[v][n] >= maxint {
				continue		// skip self. or non reachable, or non openable
			}
			path := pathFWreverse(v, n)
			if len(path) > 0 {
				tos = append(tos, path)
			}
		}
		// heuristic: sort them by biggest and closest flow first
		sort.Slice(tos, func(i, j int) bool {
			return len(tos[i]) - valves[tos[i][0]].flow < len(tos[j]) - valves[tos[j][0]].flow
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
// curflow = sum of the above
// id = id of the valve we are on
// time = remaining seconds until timeout
// pressure = the cumulated pressure upto now
// available = sum of flows of still closed valves
// path = reverse path to follow, just passing through this valve
// returns the total cumulated pressure at timeout
// uses global tempmax

func explore(flows []int, curflow, id, time, pressure, available int, path []int) (maxp int) {
	// first, see if we are at the end
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
	// then, see if we are on a path to a destination, and follow it
	if len(path) > 0 {
		p := explore(flows, curflow, path[len(path)-1], time - 1, pressure + curflow, available, path[:len(path)-1])
		if p > maxp { maxp = p;}
		goto DONE
	}
	// then, if we are openable, we must do it!
	if flows[id] == 0 && valves[id].flow > 0 {
		vflow := valves[id].flow
		flows2 := make([]int, len(flows), len(flows))
		copy(flows2, flows)
		flows2[id] = vflow
		p := explore(flows2, curflow + vflow, id, time - 1, pressure + curflow, available - vflow, nil)
		if p > maxp { maxp = p;}
		goto DONE
	}
	
	// else, choice time! Recurse on all possible paths to remaining openable valves in the to field
	for _, to := range valves[id].to {
		if flows[to[0]] == 0 {
			// set the path to destination, and go to its first step
			p := explore(flows, curflow, to[len(to)-1], time - 1, pressure + curflow, available, to[:len(to)-1])
			if p > maxp {
				maxp = p
				if maxp > tempmax { tempmax = maxp;}
			}
		}
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
// path = reverse path to follow, just passing through this valve
// epath = same for El
// returns the total cumulated pressure at timeout
// uses global tempmax


func explore2(flows []int, curflow, id, eid, time, pressure, available int, path, epath []int) (maxp int) {
	//VPf("explore2: (%d,%d)@%d flows=%v, curf=%d, press=%d, avail=%d, paths: %v %v\n", id, eid, time, flows, curflow, pressure, available, path, epath)
	// first, 3 tests to see if we are at the end of this exploration
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

	{
		// Generate actions: -1 => open this valve, -2 => follow existing path, -3 => nothing
		// i>=0 => enter path to valve at self.to[i]
		// compute the possible actions for me
		myactions := make([]int, 0, nopenable)
		elactions := make([]int, 0, nopenable)
		idopened := -1
		if len(path) > 0 {
			// if we are on a path to a destination, follow it and do nothing else
			myactions = append(myactions, -2)
		} else {
			// then, if we are openable, we must do it! Dont look for paths.
			if flows[id] == 0 && valves[id].flow > 0 {
				myactions = append(myactions, -1)
				idopened = id
			} else {
				// else, explore all possible paths to remaining openable valves in "to" field
				for toi, to := range valves[id].to {
					if flows[to[0]] == 0 {
						myactions = append(myactions, toi)
					}
				}
			}
		}
		// compute the possible actions for El in the same way
		if len(epath) > 0 {
			elactions = append(elactions, -2)
		} else {
			// exclusion principe: do not open the same valve than me at the same time
			if flows[eid] == 0 && valves[eid].flow > 0 && eid != idopened {
				elactions = append(elactions, -1)
			} else {
				for toi, to := range valves[eid].to {
					// and do not go to the same valve as me
					if flows[to[0]] == 0 && to[0] != idopened  && notIn(&myactions, to[0]) {
						elactions = append(elactions, toi)
					}
				}
			}
		}
		
		// at the minimum do nothing
		if myactions == nil { myactions = []int{-3};}
		if elactions == nil { elactions = []int{-3};}
		
		// Explore actions: then play the combination of all registered possible actions
		for _, my := range myactions {
			for _, el := range elactions {
				// reset params
				nid := id
				neid := eid
				pnflows := &flows			// pointer to the flows to use
				npath := path
				nepath := epath
				ncurflow := curflow
				navailable := available
				// play my action to change explore params
				if my == -1 {
					vflow := valves[id].flow
					flows2 := make([]int, len(flows), len(flows))
					copy(flows2, *pnflows)
					flows2[id] = vflow
					pnflows = &flows2
					navailable = navailable - vflow
					ncurflow = ncurflow + vflow
				} else if my == -2 {
					nid = path[len(path)-1]
					npath = path[:len(path)-1]
				} else if my >= 0 {
					tpath := valves[id].to[my]
					nid = tpath[len(tpath) - 1]
					npath = tpath[:len(tpath) - 1]
				}
				// play el action to change explore params
				if el == -1 {
					vflow := valves[eid].flow
					flows2 := make([]int, len(flows), len(flows))
					copy(flows2, *pnflows)
					flows2[eid] = vflow
					pnflows = &flows2
					navailable = navailable - vflow
					ncurflow = ncurflow + vflow
				} else if el == -2 {
					neid = epath[len(epath)-1]
					nepath = epath[:len(epath)-1]
				} else if el >= 0 {
					tpath := valves[eid].to[el]
					neid = tpath[len(tpath) - 1]
					nepath = tpath[:len(tpath) - 1]
				}
				p := explore2(*pnflows, ncurflow, nid, neid, time - 1, pressure + curflow, navailable, npath, nepath)
				if p > maxp { maxp = p; if maxp > tempmax { tempmax = maxp;}}
			}
		}
	}
DONE:
	if maxp > tempmax {
		tempmax = maxp
	}
	return
}
