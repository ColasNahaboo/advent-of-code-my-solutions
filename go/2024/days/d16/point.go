// Some simple operations on 2D integer coords for board games & puzzles
// I use it for solving the AdventOfCode
// Cells can be of any type T, unlike common Go matrix packages using floats

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

const (							// directions as indexes in the above slices
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
var DirsOrthoGlyph = []string{"N", "E", "S", "W"}
var DirsDiagGlyph = []string{"NE", "SE", "SW", "NW"}
var DirsAllGlyph = []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"}

// returns a new point moved one cell from p in Direction index d
func (p Point) StepOrtho(d int) Point {
	return Point{p.x + DirsOrtho[d].x, p.y + DirsOrtho[d].y}
}
func (p Point) StepDiag(d int) Point {
	return Point{p.x + DirsDiag[d].x, p.y + DirsDiag[d].y}
}
func (p Point) StepAll(d int) Point {
	return Point{p.x + DirsAll[d].x, p.y + DirsAll[d].y}
}

func RotateDir(d, r, granularity int) int {
	if d + r < 0 {
		return (d + r) % granularity + granularity
	} else {
		return (d + r) % granularity
	}
}
func RotateDirOrtho(d, r int) int {
	return RotateDir(d, r, 4)
}
func RotateDirDiag(d, r int) int {
	return RotateDir(d, r, 4)
}
func RotateDirAll(d, r int) int {
	return RotateDir(d, r, 8)
}

func HorizOrtho(d int) bool {
	return d & 1 != 0
}
func VertOrtho(d int) bool {
	return d & 1 == 0
}

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

func MakeBoard[T comparable](w, h int) (b Board[T]) {
	b.w = w
	b.h = h
	b.a = make([][]T, w)
	for x := range w {
		b.a[x] = make([]T, h)
	}
	return
}

func (b *Board[T]) Get(p Point) T {
	return b.a[p.x][p.y]
}

func (b *Board[T]) Set(p Point, v T) {
	b.a[p.x][p.y] = v
}

// Point is inside Board
func (b *Board[T]) Inside(p Point) bool {
	return p.x >= 0 && p.x < b.w && p.y >= 0 && p.y < b.h
}

// apply f to all cells, and returns ("maps") a new Board of the results
func (b *Board[T])Map (f func (b *Board[T], x, y int) T) (bp *Board[T]) {
	bb := MakeBoard[T](b.w, b.h)
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

// Fast clearing of boards to default values for T.
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

// Fill(v) is More convenient: auto-use first column as seed

func (b *Board[T]) Fill (v T) {
	// fills first column, then use it as seed
	for y := range b.h {
		b.a[0][y] = v
	}
	for x := 1; x < b.w; x++  {
		copy(b.a[x], b.a[0])
	}
}


// Creates a board from an ASCII map with func f to create board cells values
// Example of a function returning an int from a digit map
// func(x, y int, r rune) int {
//     return int(r - '0')
// }

func parseBoard[T comparable](lines []string, f func (x, y int, r rune) T) (bp *Board[T]) {
	b := MakeBoard[T](len(lines[0]), len(lines))
	for y, line := range lines {
		for x, r := range line {
			b.a[x][y] = f(x, y, r)
		}
	}
	return &b
}

////////////////////// DEBUG
func (b *Board[T]) VPBoard(title string, printcell func(c T) string) {
	if !verbose {
		return
	}
	fmt.Printf("Board %dx%d %s\n", b.w, b.h, title)
	for y := range b.h {
		for x  := range b.w {
			fmt.Print(printcell(b.a[x][y]))
		}
		fmt.Print("\n")
	}
}
