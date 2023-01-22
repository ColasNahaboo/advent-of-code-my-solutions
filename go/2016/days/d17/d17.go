// Adventofcode 2016, d17, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 input RDDRLDRURD
// TEST: ex1.txt 370
// TEST: ex2.txt 492
// TEST: ex3.txt 830
// TEST: input 448

package main

import (
	"flag"
	"fmt"
	"crypto/md5"
	// "regexp"
)

// the grid is 4x4 = 16 positions.
const gs = 4					// Grid Size
const ga = gs*gs				// Grid Area
// the 4 directions: 0=U, 1=D, 2=L, 3=R
var dirnames = [4]byte{'U','D','L','R'} // letters for hashing input
var diroffs = [4]int{-gs,gs,-1,1}		  // position offset in grid
var dirhoriz = [4]bool{false, false, true, true} // dir is horizontal?
var dirorder = [4]int{1,3,2,0}		  // explore in this order
var minlen = 888888888888888888	// shortest path length found, starts at maxint
var minpath []byte
var maxlen = 0					// longuest path length found
var maxpath []byte

var verbose bool

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	flag.Parse()
	verbose = *verboseFlag
	infile := "input.txt"
	if flag.NArg() > 0 {
		infile = flag.Arg(0)
	}
	lines := fileToLines(infile)
	passcode := []byte(lines[0])

	if *partOne {
		VP("Running Part1")
		fmt.Println(part1(passcode))
	} else {
		VP("Running Part2")
		fmt.Println(part2(passcode))
	}
}

//////////// Part 1

func part1(pass []byte) string {
	explore1(0, pass)
	if len(minpath) > len(pass) {
		return string(minpath[len(pass):])
	} else {
		return "***FAILURE!***"
	}
}

//////////// Part 2
func part2(pass []byte) int {
	explore2(0, pass)
	return len(maxpath) - len(pass)
}

//////////// Common Parts code

// fmt.Sprintf("%x", harray) but as a []byte for efficiency instead of string
var xdigits = [16]byte{'0','1','2','3','4','5','6','7','8','9','a','b','c','d','e','f'}
func x2b(harray [16]byte) []byte {
	h := make([]byte, 32, 32) 
	for i := 0; i < 16; i++ {
		h[i*2] = xdigits[harray[i]/16]
		h[i*2+1] = xdigits[harray[i]%16]
	}
	return h
}

//////////// Part1 functions

func explore1(pos int, path []byte) {
	VPf("  exploring %d, %v\n", pos, string(path))
	h := x2b(md5.Sum(path))
	for _, d := range dirorder {
		p := pos + diroffs[d]
		if p < 0 || p >= ga || (dirhoriz[d] && (p/gs) != (pos/gs)) { // wall
			continue
		}
		if h[d] <= 'a' {		// cksum char in [0-a], locked door
			continue
		}
		mypath := make([]byte, len(path)+1) // mypath must not point to path
		copy(mypath, path)
		mypath[len(path)] = dirnames[d]
		if p == ga-1 {								// V goal reached!
			if len(mypath) < minlen {
				minlen = len(mypath)
				minpath = mypath
				VPf("==> New Shortest! %d, %v\n", p, string(mypath))
			}
			continue
		}
		explore1(p, mypath)		// open door and recurse
	}
}

//////////// Part2 functions

func explore2(pos int, path []byte) {
	VPf("  exploring %d, %v\n", pos, string(path))
	h := x2b(md5.Sum(path))
	for _, d := range dirorder {
		p := pos + diroffs[d]
		if p < 0 || p >= ga || (dirhoriz[d] && (p/gs) != (pos/gs)) { // wall
			continue
		}
		if h[d] <= 'a' {		// cksum char in [0-a], locked door
			continue
		}
		mypath := make([]byte, len(path)+1) // mypath must not point to path
		copy(mypath, path)
		mypath[len(path)] = dirnames[d]
		if p == ga-1 {								// V goal reached!
			if len(mypath) > maxlen {
				maxlen = len(mypath)
				maxpath = mypath
				VPf("==> New Longuest! %d, %v\n", p, string(mypath))
			}
			continue
		}
		explore2(p, mypath)		// open door and recurse
	}
}
