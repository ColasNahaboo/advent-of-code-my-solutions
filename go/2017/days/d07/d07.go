// Adventofcode 2017, d07, in go. https://adventofcode.com/2017/day/07
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example tknk
// TEST: example 60
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

// Code is overly complex, I was afraid it may be slow at scale, so I commited
// the sin of premature optimization. I should have stayed simple and computed
// more on the fly and not cache in structure fields
// We build a tree, and traverse it in multiple passes:
// - to allocate IDs
// - to build the graph hierarchy
// - to find the root
// - to compute the total weights of all direct sons for each disk
// - to determine if one son is wrong, and what should be the good value
// - to find the node whose weight should be changed

package main

import (
	"flag"
	"fmt"
	"regexp"
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
		infile = fileMatch("input,[[:alnum:]]*,[[:alnum:]]*.test")
	}
	lines := fileToLines(infile)

	if *partOne {
		VP("Running Part1")
		sresult := part1(lines)
		fmt.Println(sresult)
	} else {
		VP("Running Part2")
		result := part2(lines)
		fmt.Println(result)
	}
}

//////////// Part 1

func part1(lines []string) string {
	discs := parse(lines)
	return discsTree(discs).name
}

//////////// Part 2

func part2(lines []string) int {
	discs := parse(lines)
	root := discsTree(discs)
	discWeights(root.id, discs)
	return balance(root.id, discs)
}

//////////// Common Parts code

type Disc struct {
	id int						// ID, index in the array of discs
	name string					// name read
	weight int					// weight of the node itself
	sons []string				// list of sons names
	sids []int					// list of sons IDs
	parent int					// ID of parent
	sonsweights int				// total weigth of sons
	wrong int					// ID of which son is unbalanced? (or -1 if OK)
	good int					// the weight all sons should all have
}

func parse(lines []string) (discs []Disc) {
	re := regexp.MustCompile("^([[:lower:]]+) [(]([[:digit:]]+)[)](.*)")
	rename := regexp.MustCompile("[[:lower:]]+")
	for _, line := range lines {
		m := re.FindStringSubmatch(line)
		if m == nil {panic("Syntax error")}
		name := m[1]
		weight := m[2]
		sons := rename.FindAllString(m[3], -1)
		if sons == nil {
			sons = []string{}
		}
		disc := Disc{id: len(discs), name: name, weight: atoi(weight), sons: sons, parent: -1, wrong: -1}
		discs = append(discs, disc)
	}
	// compute sons ids and parent
	for i, disc := range discs {
		discs[i].sids = make([]int, len(disc.sons), len(disc.sons))
		for j, son := range disc.sons {
			discs[i].sids[j] = discID(son, discs)
			discs[discs[i].sids[j]].parent = disc.id
		}
	}
	return
}

func discID(name string, discs []Disc) int {
	for i, d := range discs {
		if d.name == name {
			return i
		}
	}
	panic("DiskID not found for: " + name)
}
		
func discsTree(discs []Disc) Disc {
	for _, d := range discs {
		if d.parent == -1 {
			return d
		}
	}
	panic("Root disc not found")
}

// compute disc weights
func discWeights(d int, discs []Disc) int {
	discs[d].sonsweights = 0
	for _, sid := range discs[d].sids {
		discs[d].sonsweights += discWeights(sid, discs)
	}
	discs[d].wrong = wrongWeight(discs[d].sids, discs)
	if discs[d].wrong != -1 {
		sid := 0
		if discs[d].wrong == discs[d].sids[0] {
			sid = 1
		}
		discs[d].good =  discs[discs[d].sids[sid]].sonsweights + discs[discs[d].sids[sid]].weight
	}
	return discs[d].weight + discs[d].sonsweights
}

// find the correct weight for balancing, d being wrong
func balance(d int, discs []Disc) int {
	if discs[d].wrong == -1 {
		// our sons are OK, thus we (d) are the node that must compensate
		// our total weight .weight + .sonsweigths should be parent's .good
		return discs[discs[d].parent].good - discs[d].sonsweights
	}
	// else, go find the culprit in our "wrong" sons lineage
	return balance(discs[d].wrong, discs)
}

func weight(d int, discs []Disc) int {
	return discs[d].weight + discs[d].sonsweights
}

// return the id of the wrong one, or -1 if balanced
func wrongWeight(ids []int, discs []Disc) int {
	switch len(ids) {
	case 0, 1, 2 : return -1
	}
	good := weight(ids[0], discs)
	if  weight(ids[0], discs) !=  weight(ids[1], discs) {
		good = weight(ids[2], discs)
	}
	for _, id := range ids {
		if weight(id, discs) != good {
			return id
		}
	}
	return -1
}

	
//////////// PrettyPrinting & Debugging functions
