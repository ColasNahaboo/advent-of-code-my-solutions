// Adventofcode 2024, d20, in go. https://adventofcode.com/2024/day/20
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 -s 44 example 
// TEST: -s 50 example 285
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

// For part1, we use a brute force: try to remove ALL walls one by one and
// recompute best route, but having pre-calculated distance from S and from E
// for evaluating the potential entry and exit shortcut points

// For part2, it was too long (5mn).
// So by reading that the input maze had only a single path, I just
// took all pairs of points in this shortest path that were manhattan-distant
// of 20 or less, and see if the shortcut was shorter than the difference of
// pre-computed BFS distance from end of the two points. It takes 0.1s.

package main

import (
	"fmt"
	// "regexp"
	"flag"
)

var start, end Point
var bfs *Board[int]
var mustsave = 100
var mustsaveFlag *int

//////////// Options parsing & exec parts

func main() {
	mustsaveFlag = flag.Int("s", 100, "Must save at least these ps (def 100)")
	ExecOptions(2, XtraOpts, part1, part2)
}

func XtraOpts() {
	mustsave = *mustsaveFlag
}

//////////// Part 1

func part1(lines []string) (res int) {
	b := parse(lines)
	bfsS := BFSInit(b)
	bfsE := BFSInit(b)
	BFSfill(b, bfsS, start)		// forward BFS
	BFSfill(b, bfsE, end)		// reverse BFS
	fairtime := bfsS.Get(end)
	maxsteps := fairtime - mustsave // look only for these paths
	allcheats := AllCheats(b)
	VPf("== Testing %d cheats to go under %d\n", len(allcheats), maxsteps)
	// better paths must go through the cheat: S -path1-> Cheat -path2-> E
	// we look if we can reach cheat via two of its adjacent spaces
	for _, cheat := range allcheats {
		rs, re := []Point{}, []Point{} // reachables from start and end
		for _, d := range DirsOrtho {
			q := cheat.Add(d)
			if qds := bfsS.Get(q); qds >= 0 {
				rs = append(rs, q)
			}
			if qde := bfsE.Get(q); qde >= 0 {
				re = append(re, q)
			}
		}
		for _, ps := range rs {
			for _, pe := range re {
				if pe != ps {
					if bfsS.Get(ps) + bfsE.Get(pe) + 2 <= maxsteps {
						ns := BFSnpaths(bfsS, start, ps)
						ne := BFSnpaths(bfsE, end, pe)
						VPf("== Joining paths %v %v: %d * %d\n", ps, pe, ns, ne)
						res += ns * ne
					}
				}
			}
		}
	}
	return
}

// list all valid cheats
func AllCheats(b *Board[bool]) (cheats []Point) {
	for x := range b.w {
		for y := range b.h {
			p := Point{x, y}
			if b.GetOr(p, false) && FreeAdjCount(b, p) > 1 {
				cheats = append(cheats, p)
			}
		}
	}
	return
}

func FreeAdjCount(b *Board[bool], p Point) (adjs int) {
	for _, d := range DirsOrtho {
		q := p.Add(d)
		if b.GetOr(q, true) {
			continue			// skip if q not on free floor
		}
		adjs++
	}
	return
}

//////////// Part 2

func part2(lines []string) (res int) {
	showtime()
	b := parse(lines)
	bfs := BFSInit(b)
	BFSfill(b, bfs, end)		// reverse BFS
	bfs.Get(start)
	spath := BFSpath(bfs, start, end)
	for p := spath; p != nil; p = p.next {
		for q := p.next; q != nil; q = q.next {
			// how much do we gain by cutting p->q in manhattan?
			fairdist := bfs.Get(p.p) - bfs.Get(q.p)
			mandist := p.p.ManDist(q.p)
			//VPf("== cut %v -> %v == %d gain\n", p.p, q.p, fairdist - mandist)
			if mandist <= 20 && fairdist - mandist >= mustsave {
				res++
			}
		}
	}
	showtime("part2")
	return
}


// list all walls
func AllCheatsStarts(b *Board[bool]) (cheats []Point) {
	for x := range b.w {
		for y := range b.h {
			p := Point{x, y}
			if b.GetOr(p, false) {
				cheats = append(cheats, p)
			}
		}
	}
	return
}


//////////// Common Parts code

func parse(lines []string) *Board[bool] {
	b := ParseBoard[bool](lines, ParseCellBoolSE)
	start = boardStart
	end = boardEnd
	return b
}

//////////// PrettyPrinting & Debugging functions. See also the VPx functions.

func DEBUG() {
	if ! debug { return }
	fmt.Println("DEBUG!")
	// ad hoc debug code here
}
