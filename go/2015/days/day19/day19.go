// Adventofcode 2015, day19, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 input 509
// TEST: input 195

// the part2 computations diverge too rapidly to solve by brute force.
// So we apply the following heuristics:
// - we search in reverse, starting with the goal and applying rules in reverse
// - for each step, we only keep the N shortest molecules
// In our input, N=5 failed to find a solution, but N>5 worked
// N>32 took more than 1s.

package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

var verbose bool
var maxlen int

type (
	rule struct {
		from   string
		refrom regexp.Regexp
		to     string
		reto   regexp.Regexp
	}
)

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	maxlenFlag := flag.Int("m", 20, "heuristic for part2: recurse only on the N shortest strings")
	flag.Parse()
	verbose = *verboseFlag
	maxlen = *maxlenFlag
	infile := "input.txt"
	if flag.NArg() > 0 {
		infile = flag.Arg(0)
	}
	lines := fileToLines(infile)
	rules, molecule := parseInput(lines)

	var result int
	if *partOne {
		VP("Running Part1")
		result, _ = part1(rules, molecule)
	} else {
		VP("Running Part2")
		result = part2(rules, molecule)
	}
	fmt.Println(result)
}

//////////// Part 1

func part1(rules []rule, molecule string) (count int, results []string) {
	total := 0
	for _, rule := range rules {
		matches := rule.refrom.FindAllStringIndex(molecule, -1)
		if matches != nil {
			for _, match := range matches {
				new := molecule[:match[0]] + rule.to + molecule[match[1]:]
				if !stringMember(results, new) {
					results = append(results, new)
					count++
				}
				total++
			}
		}
	}
	VP("Total of possible replacements:", total)
	return
}

//////////// Part 2
func part2(rules []rule, goal string) int {
	return stepsToMolecule([]string{goal}, rules, 0)
}

//////////// Common Parts code

func parseInput(lines []string) (rules []rule, molecule string) {
	for i, line := range lines {
		if line == "" {
			molecule = lines[i+1]
			break
		}
		s := strings.Split(line, " => ")
		refrom := regexp.MustCompile(s[0])
		reto := regexp.MustCompile(s[1])
		rules = append(rules, rule{s[0], *refrom, s[1], *reto})
	}
	return
}

func stringMember(ss []string, s string) bool {
	for _, sse := range ss {
		if sse == s {
			return true
		}
	}
	return false
}

//////////// Part1 functions

//////////// Part2 functions

// we chain backwards: start from goal an apply the rules in reverse
// because forward-chaining like in part1 does not scale

func stepsToMolecule(molecules []string, rules []rule, steps int) int {
	VPf("  stepsToMolecule[%v]: %v molecules\n", steps, len(molecules))
	results := make([]string, 0)
	for _, molecule := range molecules {
		for _, rule := range rules {
			matches := rule.reto.FindAllStringIndex(molecule, -1)
			if matches != nil {
				for _, match := range matches {
					new := molecule[:match[0]] + rule.from + molecule[match[1]:]
					if new == "e" {
						return steps + 1
					} else {
						if !stringMember(results, new) {
							results = append(results, new)
						}
					}
				}
			}
		}
	}
	if len(results) == 0 {
		fmt.Println("***ERROR: No rule match at step! Cannot generate molecule.", steps)
		os.Exit(1)
	}
	// now heuristics: keep only the shortest results. sort and keep lower maxlen
	sort.Slice(results,
		func(i, j int) bool {
			return len(results[i]) < len(results[j])
		})
	if len(results) > maxlen {
		results = results[:maxlen]
	}
	return stepsToMolecule(results, rules, steps+1)
}
