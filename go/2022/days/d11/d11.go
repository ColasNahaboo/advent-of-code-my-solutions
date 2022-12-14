// Adventofcode 2022, d11, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 example 10605
// TEST: -1 input 102399
// TEST: example 2713310158
// TEST: input 23641658401
package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"
	"sort"
)

type Monkey struct {
	id int
	opmult bool					// true = "*", false => "+"
	oparg int					// -1 = old value
	testdiv int
	testif int
	testelse int
	inspections int
	items []int
}

var verbose bool

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
	monkeys := parseMonkeys(lines)
	for r := 0; r < 20; r++ {
		doRound(monkeys, r, 0)
	}
	if verbose {
		for m := 0; m < len(monkeys); m++ { fmt.Println(monkeys[m]);}
	}
	// sort so we can find the biggest 2 as first ones
	sort.Slice(monkeys, func(i, j int) bool {
		return monkeys[i].inspections > monkeys[j].inspections
	})
	return monkeys[0].inspections * monkeys[1].inspections
}

//////////// Part 2

// we remark that wl (worry level) ends up being used only by its modulo vs .testdiv
// so by capping iy by taking the modulo of the product of all .testdiv fields
// we can keep the size of wl under it (commondiv) while having the same results
func part2(lines []string) int {
	monkeys := parseMonkeys(lines)
	commondiv := 1				// we can mod wl by this to cap it
	for m := 0; m < len(monkeys); m++ { commondiv *= monkeys[m].testdiv;}
	
	for r := 0; r < 10000; r++ {
		doRound(monkeys, r, commondiv)
	}
	if verbose {
		for m := 0; m < len(monkeys); m++ { fmt.Println(monkeys[m]);}
	}
	sort.Slice(monkeys, func(i, j int) bool {
		return monkeys[i].inspections > monkeys[j].inspections
	})
	return monkeys[0].inspections * monkeys[1].inspections
}

//////////// Common Parts code

func parseMonkeys(lines []string) []Monkey {
	monkeys := make([]Monkey, 0)
	for i := 0; i < len(lines);  {
		monkeys = append(monkeys, parseMonkey(lines, i))
		for i += 6; i < len(lines) && lines[i] == ""; i++ {} // skip empty lines
	}
	return monkeys
}

// we brute-force parse with regexps, with minimal flexibility
var reid = regexp.MustCompile("^[[:space:]]*Monkey[[:space:]]+([[:digit:]]+)")
var reitems = regexp.MustCompile("^[[:space:]]*Starting items:[[:space:]]+([, [:digit:]]+)")
var reitemlist = regexp.MustCompile("[[:digit:]]+")
var reop = regexp.MustCompile("^[[:space:]]*Operation: new = old ([*+]) (([[:digit:]]+|old))")
var retestdiv = regexp.MustCompile("^[[:space:]]*Test: divisible by ([[:digit:]]+)")
var retestif = regexp.MustCompile("^[[:space:]]*If true: throw to monkey ([[:digit:]]+)")
var retestelse = regexp.MustCompile("^[[:space:]]*If false: throw to monkey ([[:digit:]]+)")

func parseMonkey(lines []string, i int) (monkey Monkey) {
	m := reid.FindStringSubmatch(lines[i]); i++
	if m == nil {log.Fatalf("Syntax error for id line %d: %s\n", i, lines[i-1]);}
	monkey.id = atoi(m[1])
	
	m = reitems.FindStringSubmatch(lines[i]); i++
	if m == nil {log.Fatalf("Syntax error for items line %d: %s\n", i, lines[i-1]);}
	itemstrings := reitemlist.FindAllString(m[1], -1)
	items := make([]int, len(itemstrings), len(itemstrings))
	for ii, is := range itemstrings { items[ii] = atoi(is);}
	monkey.items = items

	m = reop.FindStringSubmatch(lines[i]); i++
	if m == nil {log.Fatalf("Syntax error for op line %d: %s\n", i, lines[i-1]);}
	monkey.opmult = m[1] == "*"
	if m[2] == "old" {
		monkey.oparg = -1
	} else {
		monkey.oparg = atoi(m[2])
	}

	m = retestdiv.FindStringSubmatch(lines[i]); i++
	if m == nil {log.Fatalf("Syntax error for testdiv line %d: %s\n", i, lines[i-1]);}
	monkey.testdiv = atoi(m[1])

	m = retestif.FindStringSubmatch(lines[i]); i++
	if m == nil {log.Fatalf("Syntax error for testif line %d: %s\n", i, lines[i-1]);}
	monkey.testif = atoi(m[1])

	m = retestelse.FindStringSubmatch(lines[i])
	if m == nil {log.Fatalf("Syntax error for testelse line %d: %s\n", i, lines[i-1]);}
	monkey.testelse = atoi(m[1])

	return
}
	
func doRound(monkeys []Monkey, r int, commondiv int) {
	for m := range monkeys {
		doTurn(monkeys, m, commondiv)
	}
}

func doTurn(monkeys []Monkey, m int, commondiv int) {
	var wl, oparg int					// Worry Level
	items := monkeys[m].items
	monkeys[m].items = make([]int, 0)
	for _, item_wl := range items {
		// Monkey inspects an item with a worry level of item_wl
		if monkeys[m].oparg >= 0 {
			oparg = monkeys[m].oparg
		} else {
			oparg = item_wl
		}
		if monkeys[m].opmult {
			wl = item_wl * oparg
		} else {
			wl = item_wl + oparg
		}
		monkeys[m].inspections++
		// Monkey gets bored with item. Worry level is divided by 3
		if commondiv == 0 {
			wl /= 3
		} else {
			wl %= commondiv
		}
		// Monkey tests and throws
		if wl % monkeys[m].testdiv == 0 {
			throw(monkeys, monkeys[m].testif, wl)
		} else {
			throw(monkeys, monkeys[m].testelse, wl)
		}
	}
}

func throw(monkeys []Monkey, m int, wl int) {
	monkeys[m].items = append(monkeys[m].items, wl)
}

//////////// Part1 functions

//////////// Part2 functions
