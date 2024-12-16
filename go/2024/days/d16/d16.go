// Adventofcode 2024, d16, in go. https://adventofcode.com/2024/day/16
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 7036
// TEST: -1 example2 11048
// TEST: example
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

type Step struct {
	p Point						// position and direction index
	d int
}

type Path []Step

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
var minpaths []Path				// the list of all best paths
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
	path := Path{}
	done = MakeBoard[[4]int](maze.b.w, maze.b.h)
	done.Fill([4]int{maxint, maxint, maxint, maxint})
	// Now, forward-explore
	Explore(maze, &tracks, path, rp, rd, 0)
	VPpath(maze, minpaths)
	return minscore
}

// path is actually only used for debug in part1

func Explore(m *Maze, tracks *Board[bool], path Path, p Point, d, score int) {
	if p == m.end {				// ok, path found
		VPf("== Path found: %d %d\n", len(path), score)
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
			npath := make(Path, len(path)+1, len(path)+1)
			copy(npath, path)
			npath[len(path)] = Step{q, nd}
			VPf("== Exploring %v->%v %s->%s\n", p, q, PathGlyphs(path), DirsOrthoGlyph[nd])
			Explore(m, tracks, npath, q, nd, nscore)
			tracks.Set(q, false)
		}
	}
}

//////////// Part 2

func part2(lines []string) (res int) {
	maze, _, _ := parse(lines)
	part1(lines)
	b := MakeBoard[bool](done.w, done.h)
	for _, path := range minpaths {
		for _, p := range path {
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
	VPbestpaths(maze, b)
	return
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

func VPpath(m *Maze, paths []Path) {
	b := MakeBoard[string](m.b.w, m.b.h)
	for x := range b.w {
		for y := range b.h {
			if m.b.a[x][y] {
				b.a[x][y] = "#"
			} else {
				b.a[x][y] = "."
			}
		}
	}
	path := paths[0]
	for _, s := range path {
		b.Set(s.p, dirglyphs[s.d])
	}
	b.Set(m.start, "S")
	b.Set(m.end, "E")
	title := fmt.Sprintf("%d min paths %d steps %d turns", len(paths), len(path), Turns(path))
	b.VPBoard(title, func (c string) string { return c })
}

func VPbestpaths(m *Maze, bp Board[bool]) {
	b := MakeBoard[string](m.b.w, m.b.h)
	for x := range b.w {
		for y := range b.h {
			if bp.a[x][y] {
				b.a[x][y] = "O"
			} else if m.b.a[x][y] {
				b.a[x][y] = "#"
			} else {
				b.a[x][y] = "."
			}
		}
	}
	title := "Best places"
	b.VPBoard(title, func (c string) string { return c })
}

func Turns(p Path) (turns int) {
	od := 1
	for _, s := range p {
		if s.d != od {
			turns++
		}
		od = s.d
	}
	return
}

func PathGlyphs(p Path) string {
	b := []byte{}
	for _, s := range p {
		b = append(b, byte(DirsOrthoGlyph[s.d][0]))
	}
	return string(b)
}
