// Some simple operations on 2D integer coords for board games & puzzles

package main

import "fmt"
import "slices"

type Point struct {
	x, y int
}

// returns a new point moved by q from p
func (p Point) Add(q Point) Point {
	return Point{p.x + q.x, p.y + q.y}
}

// vector q - p: from p to q
func (p Point) Delta(q Point) Point {
	return Point{q.x - p.x, q.y - p.y}
}

func (p Point) Mult(n int) Point {
	return Point{p.x * n, p.y * n}
}

func (p Point) Div(n int) Point {
	return Point{p.x / n, p.y / n}
}

// inside grid starting at (0,0) and of (width, height)
func (p Point) IsInGrid(w, h int) bool {
	return p.x >= 0 && p.x < w && p.y >= 0 && p.y < h
}

// inside (x0,y0) inclusive and (x1, y1) exclusive, with x0,y0 < x1,y1
func (p Point) IsInRect(x0, y0, x1, y1 int) bool {
	return p.x >= x0 && p.x < x1 && p.y >= y0 && p.y < y1
}

// inside x0,y0,x1,y1 in any order, inclusive on lowers, exclusive on highers
func (p Point) IsInside(xa, ya, xb, yb int) bool {
 	x0 := min(xa, xb)
	x1 := max(xa, xb)
 	y0 := min(ya, yb)
	y1 := max(ya, yb)
	return p.x >= x0 && p.x < x1 && p.y >= y0 && p.y < y1
}

func (p Point) IsInList(list []Point) bool {
	for _, v := range list {
		if v == p {
			return true
		}
	}
	return false
}

func (p Point) Index(list []Point) int {
	for i, v := range list {
		if v == p {
			return i
		}
	}
	return -1
}

////////////////////// Directions

var DirsOrtho = []Point{Point{0,-1},Point{1,0},Point{0,1},Point{-1,0}}
var DirsDiag = []Point{Point{1,-1},Point{1,1},Point{-1,1},Point{-1,-1}}
var DirsAll = []Point{Point{0,-1},Point{1,-1},Point{1,0},Point{1,1},
	Point{0,1},Point{-1,1},Point{-1,0},Point{-1,-1}}

const (
	DirsOrthoN = 0
	DirsOrthoE = 1
	DirsOrthoS = 2
	DirsOrthoW = 3
	DirsDiagNE = 0
	DirsDiagSE = 1
	DirsDiagSW = 2
	DirsDiagNW = 3
	DirsAllN = 0
	DirsAllNE = 1
	DirsAllE = 2
	DirsAllSE = 3
	DirsAllS = 4
	DirsAllSW = 5
	DirsAllW = 6
	DirsAllNW = 7
)

func (p Point) RotateDirOrtho(n int) Point {
	i := slices.Index(DirsOrtho, p)
	if i == -1 {
		panic(fmt.Sprintf("Not a DirsOrtho: %v", p))
	}
	return DirsOrtho[(i + n) % 4]
}

func (p Point) RotateDirDiag(n int) Point {
	i := slices.Index(DirsDiag, p)
	if i == -1 {
		panic(fmt.Sprintf("Not a DirsDiag: %v", p))
	}
	return DirsDiag[(i + n) % 4]
}

func (p Point) RotateDirAll(n int) Point {
	i := slices.Index(DirsAll, p)
	if i == -1 {
		panic(fmt.Sprintf("Not a DirsAll: %v", p))
	}
	return DirsAll[(i + n) % 8]
}


////////////////////// The board, grid of cells at points

type Board[T comparable] struct {
	w, h int					// width, height
	a [][]T						// the board cells as slices of slices [X][Y]
}

// To go through all board cells {x y}
// for x := range b.w {
//     for y := range b.h {

func makeBoard[T comparable](w, h int) (b Board[T]) {
	b.w = w
	b.h = h
	b.a = make([][]T, h)
	for y := range h {
		b.a[y] = make([]T, w)
	}
	return
}

// Point is inside Board
func (b *Board[T]) Inside(p Point) bool {
	return p.x >= 0 && p.x < b.w && p.y >= 0 && p.y < b.h
}

// apply f to all cells, and returns ("maps") a new Board of the results
func (b *Board[T])Map (f func (b *Board[T], x, y int) T) (bp *Board[T]) {
	bb := makeBoard[T](b.w, b.h)
	for x, col := range b.a {
		for y := range col {
			bb.a[x][y] = f(b, x, y)
		}
	}
	return &bb
}

// apply f to all cells, that can modify them in place if needed
func (b *Board[T])Apply (f func (b *Board[T], x, y int)) {
	for x, col := range b.a {
		for y := range col {
			f(b, x, y)
		}
	}
}

	
/////////// Clears

// Clear the board to default values for T.
// Needs a clearcol seed previously created by ClearInit or ClearInitValue
func (b *Board[T]) Clear(clearcol *[]T) {
	for x := range b.w {
		copy(b.a[x], *clearcol)
	}
}

func (b *Board[T]) ClearInit() *[]T {
	cc := make([]T, b.h)
	return &cc
}

func (b *Board[T]) ClearInitValue(v T) *[]T {
	cc := make([]T, b.h)
	for i := range cc {
		cc[i] = v
	}
	return &cc
}

// Creates a board from an ASCII map with func f to create board cells values
// Example of a function returning an int from a digit map
// func(x, y int, r rune) int {
//     return int(r - '0')
// }

func parseBoard[T comparable](lines []string, f func (x, y int, r rune) T) (bp *Board[T]) {
	b := makeBoard[T](len(lines[0]), len(lines))
	for y, line := range lines {
		for x, r := range line {
			b.a[x][y] = f(x, y, r)
		}
	}
	return &b
}

////////////////////// DEBUG
func (b *Board[T]) VPBoard(titles ...string) {
	if !verbose {
		return
	}
	title := ""
	if len(titles) > 0 {
		title = titles[0]
	}
	fmt.Printf("Board %dx%d %s\n", b.w, b.h, title)
	for y := range b.h {
		for x  := range b.w {
			fmt.Print(b.a[x][y])
		}
		fmt.Print("\n")
	}
}
