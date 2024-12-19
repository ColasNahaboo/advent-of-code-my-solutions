// Adventofcode 2024, d19, in go. https://adventofcode.com/2024/day/19
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 6
// TEST: example 16
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

// We simply represent towels by their string
// The only trick is in part2 to cache the number of possible pattern arrangments
// for each design and sub-design, to avoid a combinatory explosion.

package main

import (
	"regexp"
	"strings"
)

// Towels, designs, patterns are simply strings of Colors
var Colors = []string{"w", "u", "b", "r", "g"}

//////////// Options parsing & exec parts

func main() {
	ExecOptions(2, NoXtraOpts, part1, part2)
}

//////////// Part 1

func part1(lines []string) (res int) {
	patterns, designs := parse(lines)
	for _, design := range designs {
		if Possible(patterns, design) {
			res++
		}
	}
	return 
}

func Possible(patterns []string, design string) bool {
	for _, p := range patterns {
		if strings.HasPrefix(design, p) {
			 // match pattern at start and recurse on rest
			if len(p) == len(design) ||  Possible(patterns, design[len(p):]) {
				return true
			}
		}
	}
	return false
}				



//////////// Part 2

var cache = make(map[string]int)

func part2(lines []string) (res int) {
	patterns, designs := parse(lines)
	for _, design := range designs {
		res += Layouts(patterns, design)
	}
	return 
}

func Layouts(patterns []string, design string) (layouts int) {
	if len(design) == 0 {
		return 1
	}
	if cachedLayouts, ok := cache[design]; ok {
		return cachedLayouts
	}
	for _, p := range patterns {
		if strings.HasPrefix(design, p) {
			layouts += Layouts(patterns, design[len(p):])
		}
	}
	cache[design] = layouts
	return
}				


//////////// Common Parts code

func parse(lines []string) (patterns, designs []string)  {
	reword := regexp.MustCompile("[wubrg]+")
	patterns = reword.FindAllString(lines[0], -1)
	for _, line := range lines[2:] {
		designs = append(designs, reword.FindString(line))
	}
	return
}
