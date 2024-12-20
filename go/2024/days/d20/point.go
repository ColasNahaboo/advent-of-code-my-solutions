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

// Manhattan distance
func (p Point) ManDist(q Point) int {
	return intAbs(q.x - p.x) + intAbs(q.y - p.y)
}

func (p Point) Mult(n int) Point {
	return Point{p.x * n, p.y * n}
}

func (p Point) Div(n int) Point {
	return Point{p.x / n, p.y / n}
}

func (p Point) Equal(q Point) bool {
	return p.x == q.x && p.y == q.y
}

func (p Point) Before(q Point) bool { // rows then col order
	return p.y < q.y || (p.y == q.y && p.x < q.x)
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

func (b *Board[T]) GetOr(p Point, defval T) T { // return value, or provided default
	if b.Inside(p) {
		return b.a[p.x][p.y]
	} else {
		return defval
	}
}

func (b *Board[T]) Set(p Point, v T) {
	b.a[p.x][p.y] = v
}

func (b *Board[T]) SetOr(p Point, v T) bool {
	if b.Inside(p) {
		b.a[p.x][p.y] = v
		return true
	}
	return false
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

// Scalar coords: sequential count x then y (rows by row) from 0

func (b *Board[T]) XYtoS(x, y int) int {
	return x + b.w * y
}

func (b *Board[T]) PtoS(p Point) int {
	return p.x + b.w * p.y
}
	
func (b *Board[T]) StoXY(i int) (int, int) {
	return i % b.w, i / b.w
}
	
func (b *Board[T]) StoP(i int) Point {
	return Point{i % b.w, i / b.w}
}

func (b *Board[T]) GetS(i int) T {
	return b.a[i % b.w][i / b.w]
}

func (b *Board[T]) GetSOr(i int, defval T) T { // return value, or provided default
	if b.InsideS(i) {
		return b.a[i % b.w][i / b.w]
	} else {
		return defval
	}
}

func (b *Board[T]) SetS(i int, v T) {
	b.a[i % b.w][i / b.w] = v
}

func (b *Board[T]) SetSOr(i int, v T) bool {
	if b.InsideS(i) {
		b.a[i % b.w][i / b.w] = v
		return true
	}
	return false
}

// Point is inside Board
func (b *Board[T]) InsideS(i int) bool {
	return i >= 0 && i < b.h * b.w
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

/////////// BFS

// fills a cache (bfs) of all manhattan-style distances from p on a Board
// thus shortest path length from p to any q is bfs.a[q.x][q.y]
// Warning: bfs depends on p, do not reuse with another p
// bfs must be created beforehand and filled with -1,use BFSInit
// return max distance found

func BFSfill(b *Board[bool], bfs *Board[int], p Point) (dist int) {
	bfs.Set(p, dist)
	next, todo := []Point{p}, []Point{}
	for len(next) > 0 {
		dist++
		todo, next = next, todo // alternate to avoid allocations
		next = next[:0]
		for _, p := range todo {
			for _, d := range DirsOrtho {
				q := p.Add(d)
				if bfs.GetOr(q, -2) == -1 && ! b.Get(q) {
					// valid, not visited, and free floor
					bfs.Set(q, dist)	   // cache distance
					next = append(next, q) // and recurse on it
				}
			}
		}
	}
	return
}

func BFSInit(b *Board[bool]) (bfs *Board[int]) {
	bfsI := MakeBoard[int](b.w, b.h)
	bfs = &bfsI
	bfsSeed := bfsI.ClearInitValue(-1) //  fill with -1 as default
	bfsI.Clear(bfsSeed)
	return
}

// return the number of possible distinct paths from p to end
func BFSnpaths(bfs *Board[int], p, end Point) (res int) {
	if p == end {
		return 1
	}
	dist := bfs.Get(p) + 1
	for _, d := range DirsOrtho {
		q := p.Add(d)
		if bfs.GetOr(q, -1) == dist {
			res += BFSnpaths(bfs, q, end)
		}
	}
	return
}

type PointChain struct {
	p Point
	next *PointChain
}

// return the shortest path p to end. bfs is distance to the end
func BFSpath(bfs *Board[int], p, end Point) *PointChain {
	if p == end {
		return &PointChain{p: p}
	}
	dist := bfs.Get(p) - 1
	for _, d := range DirsOrtho {
		q := p.Add(d)
		if bfs.GetOr(q, -1) == dist {
			if qpath := BFSpath(bfs, q, end); qpath != nil {
				return &PointChain{p, qpath}
			}
		}
	}
	return nil
}

//////////// Parsing

// Creates a board from an ASCII map with func f to create board cells values

func ParseBoard[T comparable](lines []string, parsecell func (x, y int, r rune) T) (bp *Board[T]) {
	b := MakeBoard[T](len(lines[0]), len(lines))
	for y, line := range lines {
		for x, r := range line {
			b.a[x][y] = parsecell(x, y, r)
		}
	}
	return &b
}

// And some common convenient parsecell functions:
// E.g. use by ParseBoard[int](lines, ParseCellInt)
func ParseCellInt(x, y int, r rune) int { // digits
     return int(r - '0')
}
func ParseCellBool(x, y int, r rune) bool { // ".#"
     return r == '#'
}
var boardStart, boardEnd Point	// globals set with marked "S" and "E" in bool map
func ParseCellBoolSE(x, y int, r rune) bool { // ".#SE"
		switch r {
		case '#': return true
		case '.': return false
		case 'S': boardStart = Point{x, y}
		case 'E': boardEnd = Point{x, y}
		default: panicf("Bad cell glyph: %v", r)
		}
		return false
}

////////////////////// DEBUG
func (b *Board[T]) VP(title string, printcell func(c T) string) {
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
