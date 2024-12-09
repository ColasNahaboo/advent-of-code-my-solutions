// Adventofcode 2024, d09, in go. https://adventofcode.com/2024/day/09
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 1928
// TEST: example 2858
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"flag"
	"fmt"
	//"regexp"
	"slices"
)

var verbose, debug bool

type Span struct { pos, size int }

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	debugFlag := flag.Bool("V", false, "debug: even more verbose")
	flag.Parse()
	verbose = *verboseFlag
	debug = *debugFlag
	var infile string
	if flag.NArg() > 0 {
		infile = flag.Arg(0)
	} else {
		infile = fileMatch("input,[[:alnum:]]*,[[:alnum:]]*.test")
	}
	lines := fileToLines(infile)
	if *partOne {
		VP("Running Part1")
		fmt.Println(part1(lines))
	} else {
		VP("Running Part2")
		fmt.Println(part2(lines))
	}
}

//////////// Part 1

func part1(lines []string) int {
	disk, _, _ := parse(lines)
	VP(*disk)
	disk = Compact(disk)
	VP(*disk)
	return Checksum(disk)
}

func Checksum(disk *[]int) (sum int) {
	for i, id := range *disk {
		if id != -1 {
			sum += i * id
		}
	}
	return
}

func Compact(diskp *[]int) (*[]int) {
	disk := *diskp				// we modify disk it in place
	for i, end := 0, len(disk)-1; i < end; i++ {
		id := disk[i]
		if id == -1 {			// we find a hole, fill it with last elt
			VPf("== moving %d@%d to %d@%d\n", disk[end], end, disk[i], i)
			disk[i] = disk[end]
			end--				// trim trailing holes
			for disk[end] == -1 {
				end--
			}
			disk = disk[:end+1]
		}
	}
	return &disk
}

//////////// Part 2

func part2(lines []string) int {
	VP(lines[0])
	disk, blocks, holes := parse(lines)
	VPDisk(*disk)
	VP(blocks)
	VP(holes)
	disk = Compact2(disk, blocks, holes)
	VPDisk(*disk)
	return Checksum(disk)
}

func Compact2(diskp *[]int, blocks, holes []Span) (*[]int) {
	disk := *diskp				// we modify disk it in place
	// Attempt to move each file exactly once in order of decreasing file ID
	for _, b := range blocks {
		// leftmost span of free space blocks that could fit the file
		for hi, h := range holes {
			if h.pos >= b.pos {	// no more holes left of block
				break
			}
			if h.size >= b.size {
				id := disk[b.pos]
				VPf("== Copy %d from %d(%d) to %d(%d)\n", id, b.pos, b.size, h.pos, h.size)
				for i := range b.size { // copy block at start of hole in map
					disk[h.pos + i] = id
				}
				for i := range b.size { // remove block from its old position
					disk[b.pos + i] = -1
				}
				// moves hole to the space after block (can be null)
				holes[hi].pos = h.pos + b.size
				holes[hi].size = h.size - b.size
				VPDisk(disk)
				break
			}
		}
	}
	return &disk
}

//////////// Common Parts code

// returns the disk block map, blocks list in reverse order, holes list
func parse(lines []string) (*[]int, []Span, []Span) {
	disk := []int{}
	blocks := []Span{}
	holes := []Span{}
	id := 0
	free := false
	pos := 0
	val := 0
	for _, s := range lines[0] {
		size := int(byte(s) - '0')
		if free {
			val = -1
			VPf("  Adding hole {%d %d}\n", pos, size)
			holes = append(holes, Span{pos, size})
		} else {
			val = id
			id++
			VPf("  Adding block {%d %d}\n", pos, size)
			blocks = append(blocks, Span{pos, size})
		}
		for _ = range size {
			disk = append(disk, val)
		}
		pos += size
		free = ! free
	}
	slices.Reverse(blocks)		// reverse blocks list
	return &disk, blocks, holes
}

//////////// PrettyPrinting & Debugging functions

func DEBUG() {
	if ! debug { return }
	// ad hoc debug code here
}

func VPDisk(disk []int) {
	if verbose {
		fmt.Printf("== 01234567 10 234567 20 234567 30 234567 40 234567 50\n== ")
		for _, id := range disk[:min(50, len(disk))] {
			if id == -1 {
				fmt.Print(".")
			} else if id >= 10 {
				fmt.Print("#")
			} else {
				fmt.Print(id)
			}
		}
		fmt.Printf("\n")
	}
}
