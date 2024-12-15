// Adventofcode 2024, d15, in go. https://adventofcode.com/2024/day/15
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 10092
// TEST: -1 example1 2028
// TEST: example 9021
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

// For part2, we use a non-optimal solution, going through all the blocks
// to push recursiveley, so we can go through some many times. E.g: pushing down
//     [aa]        on [aa] will go through b] and [c twice, through d] 4 times...
//   [bb][cc]      A more optimal solution would be to iterate by rows, marking
// [dd][ee][ff]    the blocks to move and recursing only on these
//                 But the added complexity did not seem worth it, as the naive
// solution was quite light and fast anyways.

package main

import (
	"flag"
	"fmt"
	// "slices"
)

var verbose, debug bool

const (
	FREE = 0
	BOX = 1
	WALL = 2
	LPART = 4
	LBOX = 5 					// BOX | LPART
	RPART = 8
	RBOX = 9					// BOX | RPART
)

type WH struct {				// warehouse, our problem world
	b *Board[byte]				// the floor map with walls & boxes
	r Point						// position of the robot in the map
	code []byte					// the robot program
	i int						// next instruction to execute index in code
}

var DirsOrthoGlyphs = [4]string{"^", ">", "v", "<"}

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
	b, r, code := parse(lines)
	wh := &(WH{b, r, code, 0})
	wh.Print("Initial state")
	for wh.Step() {
		if verbose {
			wh.Print(fmt.Sprintf("Move [%d] %s", wh.i-1, DirsOrthoGlyphs[wh.code[wh.i-1]]))
		}
	}
	return GPS(wh.b)
}

// robot executes one program step. Returns false if could not execute

func (wh *WH) Step() bool {
	if wh.i >= len(wh.code) {
		return false
	}
	p := wh.r.Add(DirsOrtho[wh.code[wh.i]])
	if wh.MakeRoom(p, wh.code[wh.i]) {
		wh.r = p				// move there
	}
	wh.i++
	return true
}

// pushes things at p in dir d to make room
// if p is now free, modify wh.b map in place and return true

func (wh *WH) MakeRoom(p Point, d byte) bool {
	c := wh.b.Get(p)
	switch c {
	case FREE: return true
	case WALL: return false
	case BOX, LBOX, RBOX:		// recurse to push all blocks
		q := p.Add(DirsOrtho[d])
		if wh.MakeRoom(q, d) {	// place freed, move box there
			wh.b.Set(q, c)
			wh.b.Set(p, FREE)
			return true
		}
	}
	return false
}

func GPS (b *Board[byte]) (gps int) {
	for x := range b.w {
		for y := range b.h {
			if b.a[x][y] == BOX || b.a[x][y] == LBOX {
				gps += x + 100 * y
			}
		}
	}
	return
}

//////////// Part 2

func part2(lines []string) (res int) {
	b, r, code := parse(lines)
	wh := &(WH{Boardx2(b), Point{r.x * 2, r.y}, code, 0})
	wh.Print("Initial state")
	for wh.Step2() {
		if verbose {
			wh.Print(fmt.Sprintf("Move [%d] %s", wh.i-1, DirsOrthoGlyphs[wh.code[wh.i-1]]))
		}
	}
	return GPS(wh.b)
}

func Boardx2(b *Board[byte]) *Board[byte] {
	b2 := MakeBoard[byte](b.w * 2, b.h)
	for x := range b.w {
		for y := range b.h {
			switch (*b).a[x][y] {
			case WALL: b2.a[x*2][y], b2.a[x*2+1][y] = WALL, WALL
			case BOX: b2.a[x*2][y], b2.a[x*2+1][y] = LBOX, RBOX
			}
		}
	}
	return &b2
}


// robot executes one program step. Returns false if could not execute

func (wh *WH) Step2() bool {
	if wh.i >= len(wh.code) {
		return false
	}
	p := wh.r.Add(DirsOrtho[wh.code[wh.i]])
	if HorizOrtho(int(wh.code[wh.i])) { // horizontal move: same code as part1
		if wh.MakeRoom(p, wh.code[wh.i]) {
			wh.r = p				// move there
		}
	} else {
		if wh.MakeRoom2V(p, wh.code[wh.i], false) {
			wh.MakeRoom2V(p, wh.code[wh.i], true)
			wh.r = p				// move there
		}
	}
	wh.i++
	return true
}

// pushes things at p in dir d to make room, [] boxes moving together
// if p is now free, modify wh.b map in place and return true
// only works for part2 and vertical directions
// Just check without moving boxes if doit is false

func (wh *WH) MakeRoom2V(p Point, d byte, doit bool) bool {
	c := wh.b.Get(p)
	var q1, q2, p2  Point			// the [ and ] positions

	VPf("  == MakeRoom2V: O %v [%v] pushes %d\n", p, DirsOrthoGlyphs[d], c)

	switch c {
	case FREE: return true
	case WALL: return false
	case LBOX:
		q1 = p.Add(DirsOrtho[d])		   // [
		q2 = q1.Add(DirsOrtho[DirsOrthoE]) // ]
		p2 = p.Add(DirsOrtho[DirsOrthoE])
	case RBOX:
		q2 = p.Add(DirsOrtho[d])		   // ]
		q1 = q2.Add(DirsOrtho[DirsOrthoW]) // [
		p2 = p.Add(DirsOrtho[DirsOrthoW])
	default:
		panicf("Bad cell value at %v: %v\n", p, c)
	}
	VPf("  == p=%v, p2=%v, q1=%v, q2=%v\n", p, p2, q1, q2) 
	if wh.MakeRoom2V(q1, d, doit) && wh.MakeRoom2V(q2, d, doit) { // move both
		if doit {
			wh.b.Set(p,  FREE)
			wh.b.Set(p2, FREE)
			wh.b.Set(q1, LBOX)
			wh.b.Set(q2, RBOX)
		}
		return true
	}
	return false
}

//////////// Common Parts code

func parse(lines []string) (b *Board[byte], robot Point, code []byte) {
	blocks := LineBlocks(lines)
	b = parseBoard[byte](blocks[0], func (x, y int, r rune) byte {
		switch r {
		case '#': return WALL
		case 'O': return BOX
		case '@':
			robot = Point{x, y}
			return FREE
		default: return FREE
		}
	})
	var di byte 					// direction index
	for _, line := range blocks[1] {
		for _, dc := range line {
			switch dc {
			case '^': di = DirsOrthoN
			case '>': di = DirsOrthoE
			case 'v': di = DirsOrthoS
			case '<': di = DirsOrthoW
			default: panicf("Unknown instruction: %s", string(dc))
			}
			code = append(code, di)
		}
	}
	return
}

//////////// PrettyPrinting & Debugging functions

func DEBUG() {
	if ! debug { return }
	// ad hoc debug code here
}

func (wh *WH) Print(titles ...string) {
	if ! verbose {
		return
	}
	fmt.Printf("%dx%d", wh.b.w, wh.b.h)
	for _, title := range titles {
		fmt.Printf(" %s", title)
	}
	fmt.Println(":")
	for y := range wh.b.h {
		for x := range wh.b.w {
			switch wh.b.a[x][y] {
			case WALL: fmt.Print("#")
			case BOX:  fmt.Print("O")
			case LBOX: fmt.Print("[")
			case RBOX: fmt.Print("]")
			case FREE:
				if wh.r.x == x && wh.r.y == y {
					fmt.Print("@")
				} else {
					fmt.Print(".")
				}
			default: panicf("Unknown board cell value: %v", wh.b.a[x][y])
			}			
		}
		fmt.Println()
	}
}

