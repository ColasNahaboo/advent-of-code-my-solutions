// Adventofcode 2023, d05, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: d05-input,RESULT1,RESULT2.test
// TEST: -1 example 35
// TEST: example
// And any file named d05-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"flag"
	"fmt"
	"regexp"
	"sort"
)

var verbose bool

// The maps are stored in a array of Maps in order, so we can chain them
// seeds are just an array of numbers
type Map struct {
	id int						// its index in maps array
	name string					// its symbolic name in the input file
	mappings []Mapping			// the mappings sorted by increasing from values
}
type Mapping struct {
	from int					// map from this value (included)
	to int						// to this value (excluded)
	delta int					// mapping adds this offset to value
}

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

func part1(lines []string) (lowloc int) {
	seeds, maps := parse(lines)
	lowloc = MaxInt
	for _, seed := range seeds {
		loc := mapToLocation(maps, seed)
		if loc < lowloc {
			lowloc = loc
		}
	}
	return
}

//////////// Part 2
func part2(lines []string) (lowloc int) {
	seeds, maps := parse(lines)
	lowloc = MaxInt
	for i := 0; i < len(seeds); i += 2 { // for all seed ranges
		VPf("Seed range: %s\n", VPR(seeds[i], seeds[i+1]))
		for s, l := seeds[i], seeds[i+1]; l > 0; {
			loc, ll := mapToLocationRange(maps, s, l)
			if loc < lowloc {
				lowloc = loc
			}
			s += ll			// process next sub-range
			l -= ll
			VPf(" splitted, next range: %s\n", VPR(s, l))
		}
	}
	return
}

//////////// Common Parts code

func parse(lines []string) (seeds []int, maps []Map) {
	renum := regexp.MustCompile("[0-9]+")
	remapname := regexp.MustCompile("(^[-[:lower:]]+) map:")
	remapping := regexp.MustCompile("^([0-9]+) ([0-9]+) ([0-9]+)")
	var mapping []string
	// first line: seeds
	for _, seed := range renum.FindAllString(lines[0], -1) {
		seeds = append(seeds, atoi(seed))
	}
	// then maps
	for mapno, lineno := 0, 2; lineno < len(lines); mapno++ {
		mappings := []Mapping{}
		maptitle := remapname.FindStringSubmatch(lines[lineno])
		for lineno++; lineno < len(lines); lineno++ {
			if mapping = remapping.FindStringSubmatch(lines[lineno]); mapping == nil {
				break			// map end reached
			}
			from := atoi(mapping[2])
			to := from + atoi(mapping[3])
			delta := atoi(mapping[1]) - from
			mappings = append(mappings, Mapping{from, to, delta})
		}
		if len(mappings) > 0 {	// safeguard at the eof
			sortMappings(mappings)
			m := Map{mapno, maptitle[1], mappings}
			maps = append(maps, m)
			VPf("Map[%d] %s: %v\n", m.id, m.name, m.mappings)
		}
		lineno++				// skip empty line between maps
	}
	return
}

func sortMappings(mappings []Mapping) {
	sort.Slice(mappings, func(i, j int) bool {
		return mappings[i].from < mappings[j].from
	})
}

// Pretty print a range (given by startpos, length) into a string
func VPR(s, l int) string {
	return fmt.Sprintf("%d{%d}..%d", s, l, s+l)
}

//////////// Part1 functions

func mapToLocation(maps []Map, seed int) int {
	v := seed
	VPf("Seed %d --> ", seed)
	for i := 0; i < len(maps); i++ {
		v = remap(maps[i], v)
		VPf("%d, ", v)
	}
	VPf("==> loc: %d\n", v)
	return v
}

func remap(m Map, v int) int {
	for _, mapping := range m.mappings {
		if v >= mapping.from && v < mapping.to {
			v += mapping.delta
			return v
		}
	}
	return v
}

//////////// Part2 functions

// returns the location, and the range length "covered" by it,
// that is that are guaranteed to have a larger location
func mapToLocationRange(maps []Map, s, l int) (v, ll int) {
	v = s
	ll = l
	for i := 0; i < len(maps); i++ {
		VPf("  Map#%d: %s", i, VPR(v, ll))
		v, ll = remapRange(maps[i], v, ll)
	}
	return
}

// return mapped value and the length of the range of similarly mapped values
func remapRange(m Map, s, l int) (v, ll int) {
	for _, mapping := range m.mappings {
		if s >= mapping.from && s < mapping.to {
			v = s + mapping.delta
			if (s + l) < mapping.to {
				ll = l
			} else {
				ll = mapping.to - s
			}
			VPf("  remapped %s --> %s\n", VPR(s, l), VPR(v, ll))
			return
		}
	}
	v = s
	ll = l
	VPf(" no remap\n")
	return
}
