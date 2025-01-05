// Adventofcode 2018, d08, in go. https://adventofcode.com/2018/day/08
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 138
// TEST: example 66
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"fmt"
	"regexp"
	// "flag"
	// "slices"
)

type Node struct {				// a simple tree
	sons []Node
	meta []int
}	

//////////// Options parsing & exec parts

func main() {
	ExecOptions(2, NoXtraOpts, part1, part2)
}

//////////// Part 1

var metasum int

func part1(lines []string) int {
	nums := parse(lines)
	root, p := ReadNode(nums, 0)
	VP(root)
	if p != len(nums) {
		panicf("Extra %d nums", p - len(nums))
	}
	return metasum
}

//////////// Part 2

func part2(lines []string) (res int) {
	root, _ := ReadNode(parse(lines), 0)
	return root.Value()
}

func (n Node) Value() (value int) {
	if len(n.sons) == 0 {		// no sons? sum of meta
		for _, meta := range n.meta {
			value += meta
		}
	} else {
		for _, meta := range n.meta { // metas are sons number (starting at 1)
			ix := meta - 1
			if ix >=0 && ix < len(n.sons) {
				value += n.sons[ix].Value()
			}
		}
	}
	return
}

//////////// Common Parts code

// parse list on numbers via recusrive descent

func ReadNode(nums []int, p int) (node Node, np int) {
	if p >= len(nums) {
		panic("short read!")
	}
	nsons, nmeta := nums[p], nums[p+1]
	VPf("  reading node @%d, %d sons, %d meta\n", p, nsons, nmeta)
	p += 2
	var son Node
	for s := 0; s < nsons; s++ {
		son, p = ReadNode(nums, p)
		node.sons = append(node.sons, son)
	}
	for m := 0; m < nmeta; m++ {
		node.meta = append(node.meta, nums[p])
		metasum += nums[p]
		p++
	}
	np = p
	return
}

//  simply get the list of numbers

func parse(lines []string) []int {
	renum := regexp.MustCompile("[[:digit:]]+")
	return atoil(renum.FindAllString(lines[0], -1))
}

//////////// PrettyPrinting & Debugging functions. See also the VPx functions.

func DEBUG() {
	if ! debug { return }
	fmt.Println("DEBUG!")
	// ad hoc debug code here
}
