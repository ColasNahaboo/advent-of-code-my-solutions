// Adventofcode 2018, d07, in go. https://adventofcode.com/2018/day/07
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example "CABDFE"
// TEST: example 15
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

// I tried just sorting the steps via smart sorts, but to no avail.
// So I implemented it the naive way.

// for part2, I use a timeleine of events, each event being a step completed
// or a worker being freed

package main

import (
	"fmt"
	"regexp"
	"slices"
	"os"
	// "sort"
	// "flag"
)

type Step struct {
	id byte
	ready, doing, done bool			// ready to do, doing, done.
	reqs []byte					// required step ids
}
// steps = the slice of all Step, in alhpa order
// So step of id X is at index X - 'A' in steps

//////////// Options parsing & exec parts

func main() {
	ExecOptionsString(2, NoXtraOpts, part1, part2)
}

//////////// Part 1

func part1(lines []string) (res string) {
	steps := parse(lines)
	order := []int{}
	for {
		nexts := Nexts(steps)
		if len(nexts) <= 0 {
			break
		}
		order = append(order, nexts[0])
		StepDone(&steps, nexts[0])
	}
	return PrintIXs(order)
}

func StepDone(steps *[]Step, ix int) {
	(*steps)[ix].done = true
	(*steps)[ix].doing = false
STEPS:							// check all steps that this become ready
	for i, s := range (*steps) {
		if ! s.ready && ! s.doing && ! s.done {
			for _, rid := range s.reqs {
				if ! (*steps)[id2ix(rid)].done { // s has still pending reqs
					continue STEPS
				}
			}
			(*steps)[i].ready = true		// all reqs done, we are ready to go!
		}
	}
}

func Readys(steps []Step) (readys []int) {
	for ix, s := range steps {
		if s.ready && ! s.doing && ! s.done {
			readys = append(readys, ix)
		}
	}
	return
}

func Nexts(steps []Step) (nexts []int) {
	nexts = Readys(steps)
	slices.Sort(nexts)
	return
}

//////////// Part 2

type Done struct {				// at a tick, changes hapenning
	active bool					// is there something hapenning?
	steps []int					// the ids of steps completed
	workers int					// the number of workers becoming available
}
var timeline = []Done{}			// sequence of seconds with their events
	
func part2(lines []string) (res string) {
	steps := parse(lines)
	workers := 5
	mindelay := 61
	if len(steps) < 10 {		// smaller scope for examples
		workers = 2
		mindelay = 1
	}
	t := 0
	for {
		if StepsAllDone(&steps) {
			goto DONE
		}
		//VPf("== @t = %d, TODO = %v\n", t, steps)
		// we suppose the done actions have been performed for time t
		todo := Nexts(steps)
		for _, s := range todo {
			if workers > 0 {
				workers--
				delay := mindelay + s
				nt := t + delay
				if nt >= len(timeline) {
					timeline = SliceSetX[Done](timeline, nt, Done{})
				}
				timeline[nt].active = true
				steps[s].ready = false
				steps[s].doing = true
				timeline[nt].steps = append(timeline[nt].steps, s)
				timeline[nt].workers++
				VPf("  @%d, setting Done@%d: step %c, and %d worker\n", t, nt, ix2id(s), 1)
			} else {
				break
			}
		}
		if nt := NextEvent(timeline, t); nt < 0 {
			goto DONE
		} else {
			t = nt
		}
		// go to next tick, and perform the done actions there
		for _, s := range timeline[t].steps {
			VPf("  @%d, Done with step %c\n", t, ix2id(s))
			StepDone(&steps, s)
		}
		workers += timeline[t].workers
	}
DONE:
	fmt.Println(t)
	os.Exit(0)					// part2 returns an int, not a string
	return 
}

func NextEvent(timeline []Done, tick int) int {
	for t := tick + 1; t < len(timeline); t++ {
		if timeline[t].active {
			return t
		}
	}
	return -1
}

func StepsAllDone(steps *[]Step) bool {
	for _, s := range (*steps) {
		if ! s.done {
			return false
		}
	}
	return true
}

//////////// Common Parts code

// steps[ix] <==> Step{id: id}
func id2ix(id byte) int {
	return int(id - 'A')
}
func ix2id(ix int) byte {
	return 'A' + byte(ix)
}

// input lines: Step X must be finished before step Y can begin.

func parse(lines []string) (steps []Step) {
	resteps := regexp.MustCompile(" [[:upper:]] ")
	steps = []Step{}
	reqs := [][2]byte{}
	stepsmap := make(map[byte]bool)
	for _, line := range lines {
		stepsnames := resteps.FindAllString(line, -1)
		reqs = append(reqs, [2]byte{stepsnames[0][1], stepsnames[1][1]})
		stepsmap[stepsnames[0][1]] = true
		stepsmap[stepsnames[1][1]] = true
	}
	steps = make([]Step, len(stepsmap), len(stepsmap))
	for id := range stepsmap {
		step := Step{id: id}
		for _, req := range reqs {
			if req[1] == id {
				step.reqs = append(step.reqs, req[0])
			}
		}
		if len(step.reqs) == 0 {
			step.ready = true
		}
		steps[id2ix(id)] = step
	}
	return
}

func PrintIXs(ixs []int) string {
	s := []byte{}
	for _, i := range ixs {
		s = append(s, ix2id(i))
	}
	return string(s)
}
	

//////////// PrettyPrinting & Debugging functions. See also the VPx functions.

func DEBUG() {
	if ! debug { return }
	fmt.Println("DEBUG!")
	// ad hoc debug code here
}
