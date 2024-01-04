// Adventofcode 2023, d20, in go. https://adventofcode.com/2023/day/20
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example1 32000000
// TEST: -1 example2 11687500
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

// Part2 is not solvable in time in the general case.
// However, for the input we see that the target rx is the only successor to a
// predecessor flip-flop module, and this module itself is the destination of N
// flip-flop pre-predecessors modules.
// So we just look at the cycles for each pre-predecessors modules where they are
// being sent a low pulse, and consider that rx will receive a low pulse at
// the Least Common Multiple (LCM) of these cycles.

package main

import (
	"flag"
	"fmt"
	"regexp"
	"os"
)

var verbose bool

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	flag.Parse()
	verbose = *verboseFlag
	var infile string
	if flag.NArg() > 0 {
		infile = flag.Arg(0)
	} else {
		infile = fileMatch("input,[0-9]*,[0-9]*.test")
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
	broadcasterID = parse(lines)
	VP(modules)
	targetMod = -1				// dont bother
	for i :=0; i < 1000; i++ {
		VPf("############ Cycle [%d]\n", i)
		Cycle()
	}
	VPf("%d low pulses, %d high pulses\n", nPulsesLow, nPulsesHigh)
	return nPulsesLow * nPulsesHigh
}

func parse(lines []string) (broadcaster ModID) {
	re := regexp.MustCompile("^([&%]?)([[:lower:]]+)[[:space:]]*->[[:space:]]*([[:lower:]].*)$")
	reName := regexp.MustCompile("[[:lower:]]+")
	// NONE module
	modules = append(modules, Mod{id: 0, t: NONE, name: "NONE", prefix: ""})
	// 1st pass, create and allocate all modules IDs in order (nodes)
	for lineno, line := range lines {
		rule := re.FindStringSubmatch(line)
		if rule == nil {
			panic(fmt.Sprintf("Syntax error line %d: %s\n", lineno+1, line))
		}
		prefix := rule[1]
		name := rule[2]
		mod := Mod{name: name, prefix: prefix}
		switch prefix {
		case "%": mod.t = FLIP
		case "&": mod.t = CONJ
		case "": mod.t = BROAD
		}
		mod.id = ModID(len(modules))
		modules = append(modules, mod)
		moduleIDs[name] = mod.id
		if mod .t == BROAD {
			broadcaster = mod.id
		}
	}
	// 2nd pass, initialize all the modules destinations (outbound edges)
	for _, line := range lines {
		rule := re.FindStringSubmatch(line)
		name := rule[2]
		id := moduleIDs[name]
		destslist := rule[3]
		for _, destname := range reName.FindAllString(destslist, -1) {
			did, ok := moduleIDs[destname]
			if ok == false {		// NONE module
				VPf("NONE module detected: %s as dest of %s\n", destname, name)
				modules[id].out = append(modules[id].out, 0)
				continue
			}
			modules[id].out = append(modules[id].out, did)
		}
	}
	// 3rd pass, initialise all the modules origins (inbound edges)
	for oid, omod := range modules {
		for _, did := range omod.out {
			modules[did].in = append(modules[did].in, ModID(oid))
		}
	}
	// 4th pass, initialise all the FLIP modules memories
	for id, mod := range modules {
		if mod.t == CONJ {
			modules[id].mem = make([]bool, len(mod.in), len(mod.in))
		}
	}
	return
}

//////////// Part 2
func part2(lines []string) int {
	broadcasterID = parse(lines)
	VP(modules)

	// we know that rx is the dest of a "Pre Mode" with some ins.
	// we detect the cycles in these pre-mode-ins
	// and get the LCM
	isPart2 = true
	targetMod = moduleIDs["rx"]
	targetModPre := modules[targetMod].in[0]
	if modules[targetModPre].t != CONJ {
		panic("Predecessor of rx is not a Conjunction!")
	}
	monitoredMods = modules[targetModPre].in
	monitoredLasts = make([]int, len(monitoredMods))
	monitoredCycles = make([]int, len(monitoredMods))
	monitoredDone = make([]int, len(monitoredMods))
	monitoredNames := make([]string, len(monitoredMods))
	for i := range monitoredNames {
		monitoredNames[i] = modules[monitoredMods[i]].name
		if modules[monitoredMods[i]].t != CONJ {
			panic("A Pre-Predecessor of rx is not a Conjunction!")
		}
	}
	fmt.Printf("Looking for cycles for inputs %v to %s, precursor of rx:\n", monitoredNames, modules[targetModPre].name)

	cycleno = 1
	for {
		VPf("############ Button press [%d]\n", cycleno)
		Cycle()
		if targetCount > 0 {
			return cycleno
		}
		cycleno++
	}
}

func isInModidList(l []ModID, m ModID) int {
	for i, id := range l {
		if id == m {
			return i
		}
	}
	return -1
}

// a low pulse has been send to a monitored mod?
func monitorSentLow(id ModID) {
	i := isInModidList(monitoredMods, id)
	if i == -1 {
		return
	}
	fmt.Printf("Cycle [%d]: low pulse to %s\n", cycleno, modules[id].name)
	if monitoredLasts[i] > 0 {	// we already say a low pulse for i
		if monitoredCycles[i] > 0 {	// we already saw a cycle
			if cycleno - monitoredLasts[i] == monitoredCycles[i] { // we cycle!
				monitoredDone[i] = 1
			} else {			// not yet stabilized cycle length
				monitoredCycles[i] = cycleno - monitoredLasts[i]
				monitoredLasts[i] = cycleno
			}
		} else { 				// found first cycle, record it
			monitoredCycles[i] = cycleno - monitoredLasts[i]
		}
	} else { 					// first time we see a pulse, record
		monitoredLasts[i] = cycleno
	}
	for _, idone := range monitoredDone {
		if idone == 0 {
			return
		}
	}
	fmt.Printf("Cycles found: %v\n", monitoredCycles)
	rxCycles := leastCommonMultiple(monitoredCycles[0], monitoredCycles[1], monitoredCycles[2:]...)
	fmt.Println(rxCycles)
	os.Exit(0)
}


//////////// Common Parts code

// the Module type
type Mod struct {
	id ModID
	t int						// BROAD, FLIP, CONJ
	name string					// without the % or & prefix
	prefix string				// for convenience: "", "%", "&"
	out []ModID					// destination modules
	in []ModID					// reverse: modules having self as destination
	state bool					// for FLIP
	mem []bool					// last pulse received from each in
}
type ModID int
type Pulse struct {
	pulse bool					// true = high, false = low
	orig ModID					// the module emitting the pulse
	dest ModID					// its (single) destination module
}
const (							// the Mod subtypes
	NONE = 0					// untyped modules such as "output"
	BROAD = 1
	FLIP = 2
	CONJ = 3
)

// global vars
var modules []Mod
var moduleIDs = make(map[string]ModID, 0)
var broadcasterID ModID
var pulses LinkedList[Pulse]			// the FIFO of sent pulses
var nPulsesLow int						// global counts
var nPulsesHigh int
// for part2
var isPart2 bool
var targetMod ModID				// count low pulses on this Mod, if != -1
var targetCount int
var monitoredMods []ModID		// print each time they emit a low pulse
var cycleno int
var monitoredLasts, monitoredCycles, monitoredDone []int

func (p Pulse) Send() {
	pulses.Put(p)
	if p.pulse {
		nPulsesHigh++
	} else {
		nPulsesLow++
		if isPart2 {
			monitorSentLow(p.dest)
		}
	}
	VPpulse("  Sent: ", p)
}

func pulseEmitNextOne() {
	p := Pulse(pulses.Pop())
	VPpulse("        Act:  ", p)
	p.Act()
}

func (p Pulse) Act() {
	switch modules[int(p.dest)].t {
	case BROAD: modules[int(p.dest)].ActBroad(p)
	case FLIP: modules[int(p.dest)].ActFlip(p)
	case CONJ: modules[int(p.dest)].ActConj(p)
		// NONE: do nothing
	}
}

func (mod *Mod) ActBroad(p Pulse) {
	for _, did := range mod.out {
		Pulse{pulse: p.pulse, orig: mod.id, dest: did}.Send()
	}
}

func (mod *Mod) ActFlip(p Pulse) {
	if p.pulse {				// high pulse, it is ignored
		return
	}
	mod.state = ! mod.state		// flip state and re-sends it
	for _, did := range mod.out {
		Pulse{pulse: mod.state, orig: mod.id, dest: did}.Send()
	}
}

func (mod *Mod) ActConj(p Pulse) {
	for i, oid := range mod.in { // update memory
		if oid == p.orig {
			mod.mem[i] = p.pulse
			break
		}
	}
	outpulse := false
	for _, status := range mod.mem { // are all inputs high?
		if ! status {
			outpulse = true
			break
		}
	}
	for _, did := range mod.out {
		Pulse{pulse: outpulse, orig: mod.id, dest: did}.Send()
	}
}

func Cycle() {
	Pulse{dest: broadcasterID}.Send() // push button
	// route pulses, wait for all pulses to be delivered
	for pulses.isEmpty() == false {
		pulseEmitNextOne()
	}
}

//////////// PrettyPrinting & Debugging functions

func VPpulse(label string, p Pulse) {
	if ! verbose {
		return
	}
	s := "low"
	if p.pulse {
		s  = "high"
	}
	fmt.Printf("%s%s%s -%s-> %s\n", label, modules[p.orig].prefix, modules[p.orig].name, s, modules[p.dest].name)
}
