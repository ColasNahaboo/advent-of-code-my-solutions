// Some simple operations on 2D integer coords for board games & puzzles

package main

import "fmt"

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

func RotateDirOrtho(stepright int) Point {
	return DirsOrtho[(stepright + 1) % 4]
}

func RotateDirDiag(stepright int) Point {
	return DirsDiag[(stepright + 1) % 4]
}

func RotateDirAll(stepright int) Point {
	return DirsAll[(stepright + 1) % 8]
}

////////////////////// The board, grid of cells at points

type Board[T comparable] struct {
	w, h int					// width, height
	a [][]T						// the board cells as slices of slices [X][Y]
}

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

// Creates a board from an ASCII map with func f to create board cells values
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
