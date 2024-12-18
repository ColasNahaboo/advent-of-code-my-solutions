// Adventofcode 2024, d16, in go. https://adventofcode.com/2024/day/16
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 7036
// TEST: -1 example2 11048
// TEST: example 45
// TEST: example2 64
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"flag"
	"fmt"
	// "slices"
)

var verbose, debug bool

type Maze struct {
	b Board[bool]				// the walls
	start, end Point			// start and end points in the maze
}

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

var turns = [3]int{0, -1, 1}	// we try straight, then left, then right
var minscore = maxint			// and their best score
var dirglyphs = [4]string{"^", ">", "v", "<"}
// the done map is a way for explorations to mark each dir at each p:
// "I am already here at score S, dont continue if your score is already higher"
var done Board[[4]int]

// part1 will also be used ffor part2, so we keep extra info such as minpaths
// as part1 is only interested in minscore

func part1(lines []string) int {
	maze, rp, rd := parse(lines)
	VPmaze(maze)
	// keep track of the current path so we do not loop
	tracks := MakeBoard[bool](maze.b.w, maze.b.h)
	tracks.Set(rp, true)
	done = MakeBoard[[4]int](maze.b.w, maze.b.h)
	done.Fill([4]int{maxint, maxint, maxint, maxint})
	// Now, forward-explore
	Explore(maze, &tracks, rp, rd, 0)
	return minscore
}

// Explore but only to find the best score

func Explore(m *Maze, tracks *Board[bool], p Point, d, score int) {
	if p == m.end {				// ok, path found
		VPf("== Path found: %d\n", score)
		if score < minscore {
			minscore = score
		}
		return
	}
	for _, turn := range turns {	// straight, then turn left and right
		nd := RotateDirOrtho(d, turn)
		var nscore int
		if q := p.StepOrtho(nd); ! (tracks.Get(q) || m.b.Get(q)) { // free?
			if turn == 0 {
				nscore = score + 1
			} else {
				nscore = score + 1001
			}
			if nscore >= minscore { // too big already, abort
				return
			}
			if nscore >= done.a[q.x][q.y][nd] { // another exploration is better
				return
			}
			done.a[q.x][q.y][nd] = nscore
				
			tracks.Set(q, true)
			Explore(m, tracks, q, nd, nscore)
			tracks.Set(q, false)
		}
	}
}

//////////// Part 2

// a path is a reverse linked list of steps
type Path *Step

type Step struct {
	p Point						// position and direction index
	d int
	prev *Step					// reverse chain
}

var minpaths []Path				// the list of all best paths
var nilpath Path

func part2(lines []string) (res int) {
	// similar to part1, but we keep all the best paths
	maze, rp, rd := parse(lines)
	startStep := Step{p: maze.start, d: DirsOrthoE}
	nilpath = &startStep
	VPmaze(maze)
	// keep track of the current path so we do not loop
	tracks := MakeBoard[bool](maze.b.w, maze.b.h)
	tracks.Set(rp, true)
	path := nilpath
	done = MakeBoard[[4]int](maze.b.w, maze.b.h)
	done.Fill([4]int{maxint, maxint, maxint, maxint})
	// Now, forward-explore
	ExploreAll(maze, &tracks, path, rp, rd, 0)

	b := MakeBoard[bool](done.w, done.h)
	for _, path := range minpaths {
		for p := path; p != nilpath; p = p.prev {
			b.a[p.p.x][p.p.y] = true
		}
	}
	b.Set(maze.start, true)
	for _, col := range b.a {
		for _, bestplace := range col {
			if bestplace {
				res++
			}
		}
	}
	return
}


// Explore to record ALL the paths with the best score

func ExploreAll(m *Maze, tracks *Board[bool], path Path, p Point, d, score int) {
	if p == m.end {				// ok, path found
		VPf("== Path found: %d\n", score)
		if score < minscore {
			minscore = score
			minpaths = []Path{path}
		} else if score == minscore {
			minpaths = append(minpaths, path)
		}
		return
	}
	for _, turn := range turns {	// straight, then turn left and right
		nd := RotateDirOrtho(d, turn)
		var nscore int
		if q := p.StepOrtho(nd); ! (tracks.Get(q) || m.b.Get(q)) { // free?
			if turn == 0 {
				nscore = score + 1
			} else {
				nscore = score + 1001
			}
			if nscore > minscore { // too big already, abort
				return
			}
			if nscore > done.a[q.x][q.y][nd] { // another exploration is better
				return
			}
			done.a[q.x][q.y][nd] = nscore
				
			tracks.Set(q, true)
			npath := &Step{q, nd, path}
			VPf("== Exploring %v->%v ->%s\n", p, q, DirsOrthoGlyph[nd])
			ExploreAll(m, tracks, npath, q, nd, nscore)
			tracks.Set(q, false)
		}
	}
}


//////////// Common Parts code

func parse(lines []string) (*Maze, Point, int) {
	var start, end Point
	b := *parseBoard[bool](lines, func (x, y int, r rune) bool {
		switch r {
		case '#': return true
		case '.': return false
		case 'S': start = Point{x, y}
		case 'E': end = Point{x, y}
		default: panicf("Bad cell glyph: %v", r)
		}
		return false
	})
	maze := Maze{b, start, end}
	return &maze, start, DirsOrthoE
}

//////////// PrettyPrinting & Debugging functions

func DEBUG() {
	if ! debug { return }
	// ad hoc debug code here
}

func VPmaze(maze *Maze) {
	title := fmt.Sprintf("S=%v, E=%v", maze.start, maze.end)
	maze.b.VPBoard(title, func (c bool) string {
		if c {
			return "#"
		} else {
			return "."
		}
	})
}
