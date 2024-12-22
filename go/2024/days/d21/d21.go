// Adventofcode 2024, d21, in go. https://adventofcode.com/2024/day/21
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 126384
// TEST: -1 example1 1972
// TEST: -1 example2 24256
// TEST: example
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

// This was hard and insane to debug.
// the trick is to cache the expansion of all the one-step moves, but for
// ALL recursion depths

package main

import (
	"fmt"
	"flag"
	//"slices"
)

type Keypad struct {
	b *Board[int]
	glyphs string
	keys []Point
	hole Point
}

var knA = 10						// keys 0->keyA, -1 is forbidden cell
var kdA = 4
// the maps of forbidden positions (walls) on a Numeric KeyPad and Directional
var nkpMap = Board[int]{3, 4,
	[][]int{
		[]int{7, 4, 1, -1},
		[]int{8, 5, 2, 0},
		[]int{9, 6, 3, 10},
	}}
var nhole = Point{0, 3}
var nkpGlyphs = "0123456789A"
var nkpKeys = MakeKeys(nkpMap, knA)

var dkpMap = Board[int]{3, 2,	// we use the DirIndexes [0:4] NESW
	[][]int{
		[]int{-1, 3},
		[]int{0, 2},
		[]int{4, 1},
	}}
var dkpGlyphs = "^>v<A"
var dkpKeys = MakeKeys(dkpMap, kdA)
var dhole = Point{0, 0}
var nkp = Keypad{&nkpMap, nkpGlyphs, nkpKeys, nhole} // n-kpad on door
var dkp = Keypad{&dkpMap, dkpGlyphs, dkpKeys, dhole} // d-kpad on robots

//////////// Options parsing & exec parts

var inputString *string
func main() {
	inputString = flag.String("s", "", "take string argument as input")
	ExecOptions(2, NoXtraOpts, part1, part2)
}

//////////// Part 1

func part1(lines []string) (res int) {
	codes := parse(lines)
	for _, code := range codes {
		VPf("== code: %v\n", code)
		dseq := NSeq(code)		// cache only the dseq
		VPseq("  == dseq1", dseq)
		slen := SeqLen(dseq, 2)
		VPf("  Code %v: seq len = %d * %d = %d\n", code, slen, CodeNum(code), slen * CodeNum(code))
		res += slen * CodeNum(code)
	}
	return
}

// this function expands up to final level a pair of d-keys, with caching
// only returns the length to avoid building the full huge sequences

func SeqLen(seq []int, d int) (slen int) {
	ck := cachekey(seq, d)
	if cachedlen, ok := cache[ck]; ok {
		return cachedlen
	}
	if d == 0 { return len(seq) } // reached final depth
	prevk := kdA
	for _, k := range seq {
		slen += PairLen(prevk, k, d)
		prevk = k
	}
	cache[ck] = slen
	return
}
	
func PairLen(o, k, d int) (slen int) {
	if k == o { return 1 }
	nseq := DSeqStep(dkpKeys, dhole, o, k)
	//if d == 1 { VPdpair(o, k, nseq) }
	return SeqLen(nseq, d-1)	// recurse
}

func NSeq(s []int) (ns []int) {
	old := knA
	var ndseq []int
	for _, c := range s {
		ndseq = DSeqStep(nkpKeys, nhole, old, c)
		ns = append(ns, ndseq...)
		old = c
	}
	return
}

var nilseq = []int{kdA}
// d-keys typed on the user d-pad to move from keys i to j on the n-pad or d-pad
func DSeqStep(keys []Point, hole Point, i, j int) (dkeys []int) {
	if i == j {
		return nilseq
	}
	dx, keydirx := StepTo(keys[i].x, keys[j].x, 1, 3) // dkeys: > <
	dy, keydiry := StepTo(keys[i].y, keys[j].y, 2, 0) // dkeys: v ^

	// TODO: needed exceptions to the path finding rules below. but WHY?
	if i == 2 && j == kdA && hole == dhole { // for d-pad
		return []int{0, 1, kdA}	// vA ==> v>A
	}
	if hole == nhole {			// for n-pad
		if i == 0 && j == 3 {
			return []int{0, 1, kdA} // 03 ==> v>A
		}
		if i == 1 && j == 9 {
			return []int{0, 0, 1, 1, kdA} // 19 ==> vv>>A
		}			
	}
	
	if (keydiry > keydirx &&	// start by the y dir. if hole is not in the way
		! (keys[i].x == 0 && keys[j].y == hole.y)) ||
		(keys[i].y == hole.y && keys[j].x == 0) {
		for y := keys[i].y; y != keys[j].y; y += dy {
			dkeys = append(dkeys, keydiry)
		}
		for x := keys[i].x; x != keys[j].x; x += dx {
			dkeys = append(dkeys, keydirx)
		}
	} else {
		for x := keys[i].x; x != keys[j].x; x += dx {
			dkeys = append(dkeys, keydirx)
		}
		for y := keys[i].y; y != keys[j].y; y += dy {
			dkeys = append(dkeys, keydiry)
		}
	}
	VPpair(dkeys, hole, i, j)
	return append(dkeys, kdA)
}

func StepTo(i, j, kpos, kneg int) (int, int) {
	if j > i { return 1, kpos
	} else if j < i { return -1, kneg
	} else { return 0, -1
	}
}	

// the numeric part of the code: 029A => 29
func CodeNum(code []int) (res int) {
	for _, n := range code {
		if n == knA || (res == 0 && n == 0) {
			continue
		}
		res = res * 10 + n
	}
	return
}

//////////// Part 2

func part2(lines []string) (res int) {
	codes := parse(lines)
	for _, code := range codes {
		VPf("== code: %v\n", code)
		dseq := NSeq(code)		// cache only the dseq
		VPseq("  == dseq1", dseq)
		slen := SeqLen(dseq, 25)	// 25 robots
		VPf("  Code %v: seq len = %d * %d = %d\n", code, slen, CodeNum(code), slen * CodeNum(code))
		res += slen * CodeNum(code)
	}
	return
}

//////////// Common Parts code

type SeqKey struct {
	seq string					// string made of the dkpKeys of the []int seq
	depth    int
}
var cache = make(map[SeqKey]int)

func cachekey(seq []int, d int) SeqKey {
	bs := make([]byte, len(seq), len(seq))
	for i, c := range seq {
		bs[i] = byte(dkpGlyphs[c])
	}
	return SeqKey{string(bs), d}
}

func MakeKeys(b Board[int], kA int) (keys []Point) {
	keys = make([]Point, kA+1, kA+1)
	for x, col := range b.a {
		for y, k := range col {
			if k >= 0 {
				keys[k] = Point{x, y}
			}
		}
	}
	return
}

func parse(lines []string) (codes [][]int) {
	if len(*inputString) > 0 {
		lines = []string{*inputString}
	}
	for _, line := range lines {
		code := []int{}
		for _, r := range line {
			if r == 'A' {
				code = append(code, knA)
			} else {
				code = append(code, int(r - '0'))
			}
		}
		codes = append(codes, code)
	}
	return
}

//////////// PrettyPrinting & Debugging functions. See also the VPx functions.

func VPseq(title string, s []int) {
	if ! verbose { return }
	fmt.Printf("%s: ", title)
	for _, c := range s {
		fmt.Print(dkpGlyphs[c:c+1])
	}
	fmt.Println()
}	

func VPpair(keys []int, hole Point, i, j int) {
	if ! verbose { return }
	kp := dkp
	if hole == nhole { kp = nkp }
	fmt.Printf("    {'%s', '%s'}: \"", kp.glyphs[i:i+1], kp.glyphs[j:j+1])
	for _, k := range keys {
		fmt.Print(dkpGlyphs[k:k+1])
	}
	fmt.Println("A\",")
}	

func VPdpair(i, j int, seq []int) {
	if ! verbose { return }
	fmt.Printf("    {'%s', '%s'}: \"", dkp.glyphs[i:i+1], dkp.glyphs[j:j+1])
	for _, k := range seq {
		fmt.Print(dkpGlyphs[k:k+1])
	}
	fmt.Println("A\",")
}	


func DEBUG() {
	if ! debug { return }
	fmt.Println("DEBUG!")
	// ad hoc debug code here
}
