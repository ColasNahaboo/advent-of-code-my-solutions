// Adventofcode 2023, d19, in go. https://adventofcode.com/2023/day/19
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 19114
// TEST: example 167409079868000
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

// Part 3 is a alternative part2 solution, a slow, brute force solution where we
// pre-compute all the possible spans.
// Turns out we then must look at a lot of never reached values, wasting time.
// On our input, it means testing 4365450180 cat combos: 252 x 245 x 273 x 259
// taking 2mn30, whereas part2 explores only 514 cat combos in less than 0.05s
// It is just here for reference.

package main

import (
	"flag"
	"fmt"
	"regexp"
	"sort"
)

const ( X = 0; M = 1; A = 2; S = 3) // index in Part arrays
var catnames = []string{"x", "m", "a", "s"}
const ( OK = -1; KO = -2)		// pseudo-worflows used as final destinations
var finalNames = []string{"", "A", "R"}
type Part [4]int
type Rule struct {
	cat int
	sup	bool						// true for '>', false for '<'
	val int							// numeric value to compare to
	dest int						// -1 for A, -2 for rej, or workflow index
	destname string					// name of the dest field
}
type Workflow struct {
	name string
	rules []Rule
	last int					// destination at the end of the workflow
	lastname string				// its name
}

var verbose bool

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	partThree := flag.Bool("3", false, "run exercise part3, (default: part2)")
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

func part1(lines []string) (sum int) {
	workflows, parts := parse(lines)
	in := workflowName2id(workflows, "in")
	for _, part := range parts {
		if accept(workflows, in, part) {
			VPf("  OK part %v\n", part)
			sum += partScore(part)
		}
	}
	return
}

func parse(lines []string) (workflows []Workflow, parts []Part) {
	workflows = parseWorkflows(lines)
	partsStart := len(workflows) + 1
	parts = parseParts(lines[partsStart:], workflows)
	return
}

func parseWorkflows(lines []string) (workflows []Workflow) {
	reWF := regexp.MustCompile("^([[:lower:]]+)[{]([^}]+),([^,]+)[}]")
	reRule := regexp.MustCompile("([xmas])([<>])([0-9]+):([[:alpha:]]+)")
	var lineno int
	var m []string
	var mm [][]string
	// parse workflows
	for lineno = 0; lineno < len(lines); lineno++ {
		if lines[lineno] == "" {
			break
		}
		if m = reWF.FindStringSubmatch(lines[lineno]); m == nil {
			panic(fmt.Sprintf("Workflow Syntax error line %d: \"%s\"\n", lineno, lines[lineno]))
		}
		name := m[1]
		rulestext := m[2]
		last := m[3]
		if mm = reRule.FindAllStringSubmatch(rulestext, -1); mm == nil {
			panic(fmt.Sprintf("Syntax error  in rules list line %d: \"%s\"\n", lineno, lines[lineno]))
		}
		rules := []Rule{}
		for _, m = range mm {
			cat := indexOfString(catnames, m[1])
			rules = append(rules, Rule{cat: cat, sup: m[2] == ">", val: atoi(m[3]), dest: 0, destname: m[4]})
		}
		workflows = append(workflows, Workflow{name: name, rules: rules, lastname: last})
	}
	// now set the destination IDs in workflows and their rules
	for i, wf := range workflows {
		workflows[i].last = workflowName2id(workflows, wf.lastname)
		for j := range wf.rules {
			workflows[i].rules[j].dest = workflowName2id(workflows, workflows[i].rules[j].destname)
		}
	}
	return
}
	
func parseParts(lines []string, workflows []Workflow) (parts []Part) {
	rePart := regexp.MustCompile("[{]x=([0-9]+),m=([0-9]+),a=([0-9]+),s=([0-9]+)[}]")
	var m []string
	for lineno := 0; lineno < len(lines); lineno++ {
		if m = rePart.FindStringSubmatch(lines[lineno]); m == nil {
			panic(fmt.Sprintf("Parts Syntax error line %d: \"%s\"\n", lineno, lines[lineno]))
		}
		parts = append(parts, Part{atoi(m[X+1]), atoi(m[M+1]), atoi(m[A+1]), atoi(m[S+1])})
	}
	return
}

func workflowName2id(workflows []Workflow, name string) int {
	if name == "A" {
		return OK
	} else if name == "R" {
		return KO
	}
	for i := 0; i < len(workflows); i++ {
		if name == workflows[i].name {
			return i
		}
	}
	panic(fmt.Sprintf("No workflow named \"%s\"", name))
}

func partScore(part Part) (sum int) {
	for _, cat := range part {
		sum += cat
	}
	return
}

func accept(workflows []Workflow, wf int, part Part) bool {
	for wf >= 0 {
		wf = applyWorkflow(workflows[wf], part)
	}
	return wf == OK
}

func applyWorkflow(wf Workflow, part Part) int {
	for _, rule := range wf.rules {
		if rule.sup {
			if part[rule.cat] > rule.val {
				return rule.dest
			}
		} else {
			if part[rule.cat] < rule.val {
				return rule.dest
			}
		}
	}
	return wf.last
}

//////////// Part 2
// We start with workflow "in", and explore all the possible branches, splitting
// the categories of possible Parts into spans that act the same for conditions.

func part2(lines []string) (sum int) {
	workflows := parseWorkflows(lines)
	// the possible parts combo that can pass this wf, spans of XMAS cat values
	spans := [4][2]int{{1, 4001},{1, 4001}, {1, 4001}, {1,4001}}
	// now we prune this possible space by running through the workflows
	sum = exploreOkSpans(workflows, workflowName2id(workflows, "in"), spans)
	VPf("Explored %d combos\n", combos)
	return
}

var combos int

func exploreOkSpans(workflows []Workflow, wf int, spans [4][2]int) int {
	// reached conclusion. Compute for how many parts combos it applies
	if wf == OK {
		parts := 1
		for cat := 0; cat < 4; cat++ {
			parts *= spans[cat][1] - spans[cat][0]
		}
		combos++
		return parts
	} else if wf == KO {
		return 0
	}
	sum := 0					// sum of the results of the sub-explorations
	// go through rules
	for _, r := range workflows[wf].rules {
		if r.sup {
			if spans[r.cat][0] > r.val { // whole span pass, go to other wf
				return sum + exploreOkSpans(workflows, r.dest, spans)
			} else if spans[r.cat][1] <= r.val { // whole fails, go on next rule
				continue
			}
			// cut inside the span? split in 2 spans and explore each
			subspans := spans		   // create a high half span
			subspans[r.cat][0] = r.val+1 // high half, over the cut, pass
			sum += exploreOkSpans(workflows, r.dest, subspans)
			spans[r.cat][1] = r.val+1 // truncate to low half, fails, next rule
			continue
		} else {
			if spans[r.cat][1] <= r.val { // whole span pass, go to other wf
				return sum + exploreOkSpans(workflows, r.dest, spans)
			} else if spans[r.cat][0] >= r.val { // whole fails, go on next rule
				continue
			}
			// cut inside the span? split in 2 spans and explore each
			subspans := spans		   // create a a low half span
			subspans[r.cat][1] = r.val // low half, under the cut, pass
			sum += exploreOkSpans(workflows, r.dest, subspans)
			spans[r.cat][0] = r.val //  truncate to high half, fails, next rule
			continue
		}
	}
	// end of workflow
	return sum + exploreOkSpans(workflows, workflows[wf].last, spans)
}

//////////// Part 3

// a span of values that behave the same in comparisons
type Span struct {
	start, len int
}

func part3(lines []string) (sum int) {
	workflows := parseWorkflows(lines)
	sum = findOkSpansCombos(workflows, workflowName2id(workflows, "in"))
	return
}

func findOkSpansCombos(workflows []Workflow, in int) (sum int) {
	var part Part
	// find spans in all category XMAS axis
	catSpans := [4][]Span{}
	for cat :=0; cat < 4; cat++ {
		catSpans[cat] = findSpans(workflows, cat)
	}
	VPf("Testing %d cat combos: %d x %d x %d x %d\n", len(catSpans[0])*len(catSpans[1])*len(catSpans[2])*len(catSpans[3]), len(catSpans[0]), len(catSpans[1]), len(catSpans[2]), len(catSpans[3]))
	// now test all combinations of spans
	for _, xspan := range catSpans[X] {
		part[X] = xspan.start
		for _, mspan := range catSpans[M] {
			part[M] = mspan.start
			for _, aspan := range catSpans[A] {
				part[A] = aspan.start
				for _, sspan := range catSpans[S] {
					part[S] = sspan.start
					if accept(workflows, in, part) {
						sum += xspan.len * mspan.len * aspan.len * sspan.len
					}
				}
			}
		}
	}
	return
}

func findSpans(workflows []Workflow, cat int) (spans []Span) {
	seps := []int{1}
	for _, wf := range workflows {
		for _, rule := range wf.rules {
			if cat != rule.cat {
				continue
			}
			if rule.sup { 		// x > val means span starts at val+1
				seps = append(seps, rule.val + 1)
			} else {			// x < val means next span starts at val
				seps = append(seps, rule.val)
			}
		}
	}
	sort.Slice(seps, func(i, j int) bool { return seps[i] < seps[j] })
	seps = append(seps, 4001) 	// upper bound
	for i := 0; i < len(seps)-1; i++ {
		spans = append(spans, Span{seps[i], seps[i+1] - seps[i]})
	}
	return
}

//////////// Common Parts code

func indexOfString(list []string, s string) (int) {
   for i, v := range list {
       if v == s {
           return i
       }
   }
   return -1    //not found.
}

//////////// PrettyPrinting & Debugging functions
