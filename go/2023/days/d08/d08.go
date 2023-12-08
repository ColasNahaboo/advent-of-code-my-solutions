// Adventofcode 2023, d08, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 2
// TEST: -1 example2 6
// TEST: example 2
// TEST: example2 6
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

// Part 2 is solvable in reasonable time only if you state that the input
// follows stringent rules: each ghost path loops through Z-ending nodes with
// a sing;le period, starting at 0
// The flag -3 enables using solving the general case as stated on the problem
// text, but it only works in less than a minute for 4 ghosts at most
// For "normal" Part 2 we check the input follows the specific looping pattern,
// and automatically fallback to the genral case solution (aka "part 3") if we
// cannot use the simplified solution.

package main

import (
	"flag"
	"fmt"
	"regexp"
)

var verbose bool

// globals
var ids map[string]int			// node name -> id
var label []string				// node id -> name
var steps string				// instructions
var left []int					// node id, going left -> next node id
var right []int					// node id, going right -> next node id
var endz []int					// part2: 1 if node ends in 'Z', else 0
var ghosts []int				// part2: current nodes of the parallel paths

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	partThree := flag.Bool("3", false, "run exercise part2, general solution")
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
	} else if *partThree {
		VP("Running Part3")
		result = part3(lines)
	} else {
		VP("Running Part2")
		result = part2(lines)
	}
	fmt.Println(result)
}

//////////// Part 1

func part1(lines []string) (n int) {
	parse(lines)
	start := ids["AAA"]
	goal := ids["ZZZ"]
	l := len(steps)
	for node := start; node != goal; n++ {
		if steps[n % l] == 'L' {
			VPf("  Node %s, going %s --> %s\n", ppn(node), "L", ppn(left[node]))
			node = left[node]
		} else {
			VPf("  Node %s, going %s --> %s\n", ppn(node), "R", ppn(right[node]))
			node = right[node]
		}
	}
	return
}

//////////// Part 2
func part2(lines []string) int {
	parse(lines)
	l := len(steps)
	loops := []int{}				// loop sizes per ghost
	// first, check the input is in the special case of fixed-size loops
	// starting at 0 for ALL the ghsots. Otherwise, switch to part3 code
	maxloops := 10
	VPf("%d ghosts, checking %d loops:\n", len(ghosts), maxloops)
	for i, g := range ghosts {
		VPf("  Ghost %d:", i)
		last := 0
		loop := 0
		for n, nloop := 0, 0; nloop < maxloops; n++ {
			if endz[g] != 0 {
				VPf(" %d", n - last)
				if loop != 0 && (n - last) != loop {
					VPf("\n### Warning: Switching to the general method, as for ghost %d, the loop size varies: %d %d\n", i, loop, n-last)
					return part3(lines)
				}
				last = n
				if loop == 0 {
					loop = n
					loops = append(loops, loop)
				}
				nloop++
			}
			if steps[n % l] == 'L' {
				g = left[g]
			} else {
				g = right[g]
			}
		}
		VP("")
	}

	// In this simple case, the solution is the least common denominator, LCM
	VP("We are in the simple case: all ghosts having fized size 0-starting loops")
	return LCM(loops...)
}

// General Variant, "Part 3"

func part3(lines []string) (n int) {
	parse(lines)
	l := len(steps)
	VPf("%d ghosts\n", len(ghosts))
	for {
		sum := 0
		for _, g := range ghosts {
			sum += endz[g]
		}
		if verbose {
			fmt.Printf("%2d:", n)
			for _, g := range ghosts {
				pad := ""
				if g < 10 {
					pad = "  "
				} else if g < 100 {
					pad = " "
				}
				fmt.Printf(" %s%s", pad, ppn(g))
			}
			fmt.Printf(" (%d)\n", sum)
		}
		if sum == len(ghosts) {
			break				// all ghost paths are on nodes ending with Z
		}
		if steps[n % l] == 'L' {
			for i, g := range ghosts {
				ghosts[i] = left[g]
			}
		} else {
			for i, g := range ghosts {
				ghosts[i] = right[g]
			}
		}
		n++
	}
	return
}


//////////// Common Parts code

func parse(lines []string) {
	steps = lines[0]
	ids = make(map[string]int, 0)
	re := regexp.MustCompile("^([[:upper:]]{3}) = [(]([[:upper:]]{3}), ([[:upper:]]{3})[)]")
	for i := 2; i < len(lines); i++ {
		line := lines[i]
		if len(line) == 0 {
			continue
		}
		m := re.FindStringSubmatch(line)
		if m == nil {
			panic(fmt.Sprintf("Syntax Error line %d: \"%s\"", i, line))
		}
		node := label2id(m[1])
		nodeleft := label2id(m[2])
		noderight := label2id(m[3])
		left[node] = nodeleft
		right[node] = noderight
	}
}

// returns ids[label], auto-allocating a new one if needed
func label2id(name string) (id int) {
	val, ok := ids[name]
	if ok == false {
		id = len(ids)
		ids[name] = id
		// also extend the size of other slices indexed by ids
		label = append(label, name)
		left = append(left, 0)
		right = append(right, 0)
		last := 0
		if name[2] == 'A' {
			ghosts = append(ghosts, id) // part2: add ghost path starting here
		} else if name[2] == 'Z' {
			last = 1
		}
		endz = append(endz, last)
		return
	}
	id = val
	return
}
		
func ppn(id int) string {
	return fmt.Sprintf("%s[%d]", label[id], id)
}

//////////// Part1 functions

//////////// Part2 functions


//////////// classical computation of the LCM via the GCD

// greatest common divisor (GCD) via Euclidean algorithm
func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
              a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func LCM(arg ...int) int {
	result := arg[0] * arg[1] / GCD(arg[0], arg[1])
	for i := 2; i < len(arg); i++ {
		result = LCM(result,
			arg[i])
	}
	return result
}
