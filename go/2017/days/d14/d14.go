// Adventofcode 2017, d14, in go. https://adventofcode.com/2017/day/14
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 8108
// TEST: example 1242
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

// for convenience, -x and -b just print the Knot hash in hexa or binary

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var verbose bool

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	hashHexFlag := flag.Bool("x", false, "just print the Knot Hash of argument in hex")
	hashBitFlag := flag.Bool("b", false, "just print the Knot Hash of argument in binary")
	flag.Parse()
	verbose = *verboseFlag
	var infile string
	if *hashHexFlag {
		fmt.Println(KnotHashHex(flag.Arg(0)))
		os.Exit(0)
	} else if *hashBitFlag {
		fmt.Println(KnotHashBit(flag.Arg(0)))
		os.Exit(0)
	}
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

func part1(lines []string) (used int) {
	key := lines[0]
	for row := 0; row < 128; row++ {
		rowkey := key + "-" + itoa(row)
		kh := KnotHash(rowkey)
		if row < 8 {			// debug output similar to the one in text
			VProw(kh)
		}
		used += numOfBits(kh)
	}
	return
}

func numOfBits(kh [16]byte) (n int) {
	for _, b := range kh {
		bs := fmt.Sprintf("%b", b)
		for _, c := range bs {
			if c == '1' {
				n++
			}
		}
	}
	return
}

//////////// Part 2
func part2(lines []string) (regions int) {
	key := lines[0]
	hashGrid := make([]string, 128, 128)
	regionsGrid := make([][]int, 128, 128)
	for i := range regionsGrid {
		regionsGrid[i] = make([]int, 128, 128)
	}
	for row := 0; row < 128; row++ {
		rowkey := key + "-" + itoa(row)
		hashGrid[row] = KnotHashBit(rowkey)
	}
	return findRegions(hashGrid, regionsGrid)
}

func findRegions(hg []string, rg [][]int) (rn int) {
	for y := range rg {
		for x := range rg {
			if rg[y][x] != 0 {	// already mapped, skip
				continue
			}
			if hg[y][x] == '1' {
				rn++
				mapRegion(rn, x, y, hg, rg)
			}
		}
	}
	return
}

var dirs = [4][2]int{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}

func mapRegion(rn, x, y int, hg []string, rg [][]int) {
	rg[y][x] = rn
	for _, d := range dirs {
		nx, ny := x + d[0], y + d[1]
		if nx < 0 || nx >= 128 || ny < 0 || ny >= 128 { //  out of grid
			continue
		}
		if rg[ny][nx] != 0 {	// already mapped, skip
				continue
		}
		if hg[ny][nx] == '1' {	// mark in region, and recurse
			mapRegion(rn, nx, ny, hg, rg)
		}
	}
}


//////////// Common Parts code

//////////// KnotHash code, copied/adapted from ../d10/d10.go

const kbufhsize = 256
var khbuf []byte = make([]byte, kbufhsize, kbufhsize)

func KnotHash(s string) (hash [16]byte) {
	lengths := []byte(s)
	lengths = append(lengths, []byte{17, 31, 73, 47, 23}...)
	list := make([]byte, kbufhsize, kbufhsize)
	for i := 0; i < kbufhsize; i++ {
		list[i] = byte(i)
	}
	p, skip := 0, 0
	for i := 0; i < 64; i++ {
		p, skip = khtwistAll(p, skip, lengths, list)
	}
	for i := range hash {
		x := list[i*16]
		for j := 1; j < 16; j++ {
			x ^= list[i*16+j]
		}
		hash[i] = x
	}
	return
}

func khtwistAll(p, skip int, lens, list []byte) (int, int) {
	for _, l := range lens {
		p = khtwist(p, skip, int(l), list)
		skip++
	}
	return p, skip
}

func khtwist(p, skip, l int, list []byte) (newpos int) {
	// copy the section at p of length len, wrapping index in list
	for i, j := p, 0; i < p+l; i, j = i+1, j+1 {
		khbuf[j] = list[i % len(list)]
	}
	// and re-paste it in place in reverse
	for i, j := p, l-1; i < p+l; i, j = i+1, j-1 {
		list[i % len(list)] = khbuf[j]
	}
	return (p + l + skip) % len(list)
}

func KnotHashHex(s string) (h string) {
	kh := KnotHash(s)
	for _, x := range kh {
		h += fmt.Sprintf("%02x", x)
	}
	return
}

func KnotHashBit(s string) (h string) {
	kh := KnotHash(s)
	for _, x := range kh {
		h += fmt.Sprintf("%08b", x)
	}
	return
}

//////////// PrettyPrinting & Debugging functions

func VProw(kh [16]byte) {
	if ! verbose { return }
	bs := fmt.Sprintf("%08b", kh[0])
	bs = strings.ReplaceAll(bs, "0", ".")
	bs = strings.ReplaceAll(bs, "1", "#")
	fmt.Println(bs)
}




