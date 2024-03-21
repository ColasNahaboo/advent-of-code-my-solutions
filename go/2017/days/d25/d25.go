// Adventofcode 2017, d25, in go. https://adventofcode.com/2017/day/25
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 3
// TEST: example
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"flag"
	"fmt"
	"regexp"
	"strings"
	// "golang.org/x/exp/slices"
)

var verbose, debug bool

type Tablet struct {
	tape Tape
	cursor int
	state int
	states []State
}
type State [2]Action
type Action struct {
	write bool
	move int
	state int
}
type Tape struct {
	pos, neg []bool				// two half-tapes, one for >= 0 and one for < 0
}
	

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	debugFlag := flag.Bool("V", false, "debug: even more verbose")
	flag.Parse()
	verbose = *verboseFlag
	debug = *debugFlag
	var infile string
	if flag.NArg() > 0 {
		infile = flag.Arg(0)
	} else {
		infile = fileMatch("input,[[:alnum:]]*,[[:alnum:]]*.test")
	}
	lines := fileToLines(infile)
	if *partOne {
		VP("Running Part1")
		fmt.Println(part1(lines))
	} else {
		VP("Running Part2")
		fmt.Println(part2(lines))
	}
}

//////////// Part 1

func part1(lines []string) int {
	t, diagnostic := parse(lines)
	VPf("diagnostic at %d\n", diagnostic)
	t.VP(0)
	for step := 1; step <= diagnostic; step++ {
		t.Run()
		if verbose {
			t.VP(step)
		}
	}
	return t.Checksum()
}

//////////// Part 2
func part2(lines []string) int {
	for _, line := range lines {
		fmt.Println(line)
	}
	return 0
}

//////////// Common Parts code

func (t *Tablet) Run() {
	action := 0
	if t.Read() {
		action = 1
	}
	t.Write(t.states[t.state][action].write)
	t.cursor += t.states[t.state][action].move
	t.state = t.states[t.state][action].state
}
	
func (t *Tablet) SlotGet(i int) bool {
	if i >= 0 {
		if i >= len(t.tape.pos) {
			return false
		}
		return t.tape.pos[i]
	}
	if -i-1 >= len(t.tape.neg) {
		return false
	}
	return t.tape.neg[-i-1]
}

func (t *Tablet) SlotSet(i int, val bool) {
	if i >= 0 {
		if i >= len(t.tape.pos) {
			t.tape.pos = extendTape(t.tape.pos, i+1) // realloc
		}
		t.tape.pos[i] = val
		return
	}
	if -i-1 >= len(t.tape.neg) {
		t.tape.neg = extendTape(t.tape.neg, -i) // realloc
	}
	t.tape.neg[-i-1] = val
}

func (t *Tablet) Write(val bool) {
	t.SlotSet(t.cursor, val)
}

func (t *Tablet) Read() bool {
	return t.SlotGet(t.cursor)
}

func (t *Tablet) Checksum() (cksum int) {
	for _, b := range t.tape.pos {
		if b {
			cksum++
		}
	}
	for _, b := range t.tape.neg {
		if b {
			cksum++
		}
	}
	return
}

func extendTape(htape []bool, n int) []bool {
	for i := len(htape); i < n; i++ {
		htape = append(htape, false)
	}
	return htape
}

func parse(lines []string) (t Tablet, diagnostic int) {
	diagnostic = numIn(lines[1])
	t.tape.pos = []bool{}
	t.tape.neg = []bool{}
	for i := 3; i < len(lines); i += 10 {
		a0 := parseAction(lines, i+2)
		a1 := parseAction(lines, i+6)
		state := State{a0, a1}
		t.states = append(t.states, state)
	}
	return
}

func parseAction(lines []string, i int) (a Action) {
	if numIn(lines[i]) == 1 {
		a.write = true
	}
	if strings.Index(lines[i+1], "to the right") != -1 {
		a.move = 1
	} else {
		a.move = -1
	}
	stateIdx := strings.Index(lines[i+2], " state ") + len(" state ")
	a.state = int(lines[i+2][stateIdx] - 'A')
	return
}

var numInRe = regexp.MustCompile("[[:digit:]]+")
func numIn(s string) int {
	return atoi(numInRe.FindString(s))
}

//////////// PrettyPrinting & Debugging functions

func DEBUG() {
	if ! debug { return }
	// ad hoc debug code here
}

func (t *Tablet) VP(step int) {
	if ! verbose {
		return
	}
	fmt.Printf("...")
	for i := -6; i <= 6; i++ {
		t.VPslot(i)
	}
	fmt.Printf("... (after %d steps; about to run state %s)\n", step, string('A' + byte(t.state)))
}

func (t *Tablet) VPslot(i int) {
	if i == t.cursor {
		fmt.Printf("[%d]", bool2int(t.SlotGet(i)))
	} else {
		fmt.Printf(" %d ", bool2int(t.SlotGet(i)))
	}
}
		
func bool2int(b bool) int {
	if b {
		return 1
	} else {
		return 0
	}
}
