// Adventofcode 2016, d14, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 example 22728
// TEST: -1 input 16106
// TEST: example 22551
// TEST: input 22423
package main

import (
	"flag"
	"fmt"
	"crypto/md5"
	// "regexp"
)

var salt string
var cksums []string				// cache of md5 sums for each index
var stretching bool				// are we computing the hash 2017 times?
var xdigits = [16]byte{'0','1','2','3','4','5','6','7','8','9','a','b','c','d','e','f'}

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
	salt = lines[0]
	cksums = make([]string, 0, 2000)

	var result int
	if *partOne {
		VP("Running Part1")
		result = part1()
	} else {
		VP("Running Part2")
		stretching = true
		result = part1()
	}
	fmt.Println(result)
}

//////////// Part 1

func part1() int {
	index, keynum := 0, 0
	var hash string
	for keynum < 64 {
		hash = idxMD5(index)
		if isKey(hash, index) {
			keynum++
			VPf("  OK, %d keys found! at %d\n", keynum, index)
		}
		index++
	}
	return index - 1
}

//////////// Part 2
func part2(lines []string) int {
	for _, line := range lines {
		fmt.Println(line)
	}
	return 0
}

//////////// Common Parts code

// we are sure to be called in sequential order: cksum is not sparse
func idxMD5(idx int) string {
	if idx < len(cksums) {
		return cksums[idx]
	}
	var s string
	if stretching {				// Part2: compute MD5 2017 times
		// we do the []byte => hex string conversion by hand to try to be the fastest possible
		// avoid copying the [16]byte array as much as possible, minimize conversions
		h := make([]byte, 32, 32) 
		harray := md5.Sum([]byte(salt + itoa(idx)))
		n := 1
		for {
			for i := 0; i < 16; i++ {
				h[i*2] = xdigits[harray[i]/16]
				h[i*2+1] = xdigits[harray[i]%16]
			}
			if n >= 2017 {break;}
			harray = md5.Sum(h)
			n++
		}
		s =	string(h)
	} else {					// Part1: compute MD5 once
		s = fmt.Sprintf("%x", (md5.Sum([]byte(salt + itoa(idx)))))
	}
	cksums = append(cksums, s)	// cache it
	VPf("cksums[%d] = %s\n", idx, s)
	//if idx != len(cksums) - 1 { panic(fmt.Sprintf("i = %d, len(cksums) = %d\n", idx, len(cksums)));}
	return s
}

func isKey(h string, idx int) bool {
	var c byte
	for i := 0 ; i <= len(h) - 3; i++ { // do not check when no room anymore for 3 chars
		if c = h[i]; h[i+1] == c && h[i+2] == c {
			goto FOUND3
		}
	}
	return false
FOUND3:
	VPf("triplet found at %d... checking [%d, %d]\n", idx, idx+1, idx+1000)
	for idx2 := idx+1; idx2 <= idx+1000; idx2++ {
		h2 := idxMD5(idx2)
		for i := 0; i <= len(h2) - 5; i++ {
			if h2[i] == c && h2[i+1] == c && h2[i+2] == c && h2[i+3] == c && h2[i+4] == c {
				return true
			}
		}
	}
	return false
}

//////////// Part1 functions

//////////// Part2 functions
