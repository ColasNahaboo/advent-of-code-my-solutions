// Adventofcode 2025, d11, in go. https://adventofcode.com/2025/day/11
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 5
// TEST: example 0
// TEST: example2 2
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

// Input properties:
// less than 640 devices max 25 outputs each

package main

import (
	"fmt"
	"regexp"
	// "flag"
	// "slices"
)

// Implementation:
// devices (nodes in the directed graph) are represented by their numeric ID
// IDs are indexes in the array of possible node names
// we force "you" (start) at ID 0 and "out" (end) at ID 1


// we cache the results of the PathsTo functions
var cacheto = map[[2]int]int{}
var cachevia1 = map[[3]int]int{}
var cachevia2 = map[[4]int]int{}

//////////// Options parsing & exec parts

func main() {
	ExecOptions(2, NoXtraOpts, part1, part2)
}

//////////// Part 1

func part1(lines []string) (res int) {
	names, outs, _ := parse(lines)
	return PathsTo(0, 1, names, outs, []int{0})
}

// the path argument is optional, only needed for tracing and debugging

func PathsTo(from, to int, names []string, outs [][]int, path []int) (npaths int) {
	if res, ok := cacheto[[2]int{from, to}]; ok {
		return res
	}
	for _, next := range outs[from] {
		if next == to {
			VPf("  PathsTo %v\n", append(path, to))
			npaths++
			continue
		}
		npaths += PathsTo(next, to, names, outs, append(path, next))
	}
	cacheto[[2]int{from, to}] = npaths
	return
}

//////////// Part 2

func part2(lines []string) (res int) {
	names, outs, ids := parse(lines)
	svr := ids["svr"]
	out := 1
	dac := ids["dac"]
	fft := ids["fft"]
	VPf("svr=%d, dac=%d, fft=%d\n", svr, dac, fft)

	dacfft := PathsToVia2(svr, out, dac, fft, names, outs, []int{svr})
	fftdac := PathsToVia2(svr, out, fft, dac, names, outs, []int{svr})

	return dacfft + fftdac
}

func PathsToVia2(from, to, via1, via2 int, names []string, outs [][]int, path []int) (npaths int) {
	if res, ok := cachevia2[[4]int{from, to, via1, via2}]; ok {
		return res
	}
	VPf("  PathsToVia2 %v\n", path)
	for _, next := range outs[from] {
		if next == via1 {
			npaths += PathsToVia(via1, to, via2, names, outs, append(path,via1))
		}
		if next == via2 || next == to {
			continue
		}
		npaths += PathsToVia2(next, to, via1, via2, names, outs, append(path, next))
	}
	cachevia2[[4]int{from, to, via1, via2}] = npaths
	return
}

func PathsToVia(from, to, via int, names []string, outs [][]int, path []int) (npaths int) {
	if res, ok := cachevia1[[3]int{from, to, via}]; ok {
		return res
	}
	VPf("  PathsToVia %v\n", path)
	for _, next := range outs[from] {
		if next == via {
			npaths += PathsTo(via, to, names, outs, append(path, via))
			continue
		}
		if next == to {
			continue
		}
		npaths += PathsToVia(next, to, via, names, outs, append(path, next))
	}
	cachevia1[[3]int{from, to, via}] = npaths
	return
}

//////////// Common Parts code

func parse(lines []string) (names []string, outs [][]int, ids map[string]int) {
	redev := regexp.MustCompile("[[:lower:]]{3}")
	ids = make(map[string]int)
	ids["you"] = 0
	ids["out"] = 1
	names = []string{"you", "out"}
	for _, line := range lines { // 1st pass: get int IDs of all device names
		conns := redev.FindAllString(line, -1)
		for _, conn := range conns {
			if _, ok := ids[conn]; ! ok {
				id := len(names)
				ids[conn] = id
				names = append(names, conn)
			}
		}
	}
	outs = make([][]int, len(names))
	for _, line := range lines { // 2nd pass: get all out connections between IDs
		conns := redev.FindAllString(line, -1)
		dev := ids[conns[0]]
		outsdev := []int{}
		for _, conn := range conns[1:] {
			outdev := ids[conn]
			outsdev = append(outsdev, outdev)
		}
		outs[dev] = outsdev
	}
	return
}

//////////// PrettyPrinting & Debugging functions. See also the VPx functions.

func DEBUG() {
	if ! debug { return }
	fmt.Println("DEBUG!")
	// ad hoc debug code here
}
