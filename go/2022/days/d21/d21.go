// Adventofcode 2022, d21, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 example 152
// TEST: -1 input 21120928600114
// TEST: example 301
// TEST: input 3453748220116

// For part2, we just find the value by linear interpolation, as the operators
// are linear.
// We addd a monkey "delta" on to of root that performs a soustraction between
// the two sub-monkeys of root. So we just aim to have this delta monkey yell zero
// We do a linear interpolation on the value of "humn" to get it close to 0

package main

import (
	"flag"
	"fmt"
	"log"
	// "regexp"
)

type Monkey struct {
	name string
	yelled bool					// is the value already defined and cached?
	value int
	m1, m2 *Monkey					// listens to these monkeys
	op string						// operator on them: + - * /
	name1, name2 string				// names of these monkeys
}

var monkeys map[string]*Monkey	// all the monkeys, indexed by name
var verbose bool
var humanval int

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	humanFlag := flag.Int("h", -1, "set value of \"humn\" monkey, return delta")
	flag.Parse()
	verbose = *verboseFlag
	humanval = *humanFlag
	infile := "input.txt"
	if flag.NArg() > 0 {
		infile = flag.Arg(0)
	}
	lines := fileToLines(infile)
	parse(lines)
	
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
	return listenCached(monkeys["root"])
}

//////////// Part 2
func part2(lines []string) int {
	root := monkeys["root"]
	delta := Monkey{name: "delta", m1: root.m1, m2: root.m2, name1: root.name1, name2: root.name2, op: "-"}
	//monkeys["delta"] = delta
	if humanval != -1 {			// juste compute a whatif? value
		human := monkeys["humn"]
		human.yelled = true
		human.value = humanval
		return listenCached(&delta)
	}
	// we find the value of human making delta == 0 by dichotomy
	return humanZeroing(&delta)
}

//////////// Common Parts code

func parse(lines []string) {
	var name, name1, op, name2 string
	var num int
	monkeys = make(map[string]*Monkey, 0)
	// 1st pass, gather names
	for _, line := range lines {
		if n, _ := fmt.Sscanf(line, "%s %d", &name, &num); n == 2 {
			m := Monkey{name: name[:len(name)-1], value: num, yelled: true}
			monkeys[name[:len(name)-1]] = &m
		} else if n, _ = fmt.Sscanf(line, "%s %s %s %s", &name, &name1, &op, &name2); n == 4 {
			m := Monkey{name: name[:len(name)-1], op: op, name1:name1,  name2:name2}
			monkeys[name[:len(name)-1]] = &m
		} else {
			log.Fatalf("Parse error: %s\n", line)
		}
	}
	// 2nd pass, fill the sub-monkey fields
	for _, m := range monkeys {
		if !m.yelled {
			m.m1 = monkeys[m.name1]
			m.m2 = monkeys[m.name2]
		}
	}
}

//////////// Part1 functions

// resolve and cache result. Faster but cannot be reused.
func listenCached(m *Monkey) int {
	if m.yelled {
		return m.value
	}
	n1 := listenCached(m.m1)
	n2 := listenCached(m.m2)
	m.yelled = true
	switch m.op {
	case "+": m.value = n1 + n2
	case "-": m.value = n1 - n2
	case "*": m.value = n1 * n2
	case "/": m.value = n1 / n2
	}
	return m.value
}

//////////// Part2 functions

// no cache, suitable for many attempts.
func listen(m *Monkey) int {
	if m.yelled {
		return m.value
	}
	n1 := listen(m.m1)
	n2 := listen(m.m2)
	switch m.op {
	case "+": return n1 + n2
	case "-": return n1 - n2
	case "*": return n1 * n2
	case "/": return n1 / n2
	}
	return 0					// not reached
}

func listenFor(d, h *Monkey, value int) int {
	h.value = value
	return listen(d)
}	

func humanZeroing(delta *Monkey) int {
	human := monkeys["humn"]
	human.yelled = true
	h1 := 0
	h2 := 1000
	for {
		d1 := listenFor(delta, human, h1)
		d2 := listenFor(delta, human, h2)
		d2d1 := d2 - d1
		VPf("%d, %d --> %d, %d (diff: %d)\n", h1, h2, d1, d2, d2d1)
		if d1 == 0 { return humanSmallest(delta, human, h1);} // found
		if d2 == 0 { return humanSmallest(delta, human, h2);} // found
		if d2d1 == 0 { log.Fatalf("d2 == d1\n");}
		// linear regression on h to minimize d
		if absInt(d1) < absInt(d2) {
			// h1 was closest to the goal: interpolate from it
			h := h1 - (d1 * (h2-h1)) / d2d1
			VPf("==> h = %d (h1=%d + %d)\n", h, h1, - (d1 * (h2-h1)) / d2d1)
			h2 = h
		} else {
			// h2 was closest to the goal: interpolate from it
			h := h2 - (d2 * (h1-h2)) / (-d2d1)
			VPf("==> h = %d (h2=%d + %d)\n", h, h2, - (d2 * (h1-h2)) / (-d2d1))
			h1 = h
		}
	}
	
}

// As many humn values can give the same result, find the smallest one.
func humanSmallest(delta, human *Monkey, hval int) (h int) {
	for h = hval; listenFor(delta, human, h-1) == 0; h-- {}		
	return
}

func absInt(i int) int {
	if i >= 0 {
		return i
	}
	return -i
}
