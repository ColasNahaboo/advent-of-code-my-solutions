// Some simple operations on 2D integer coords for board games & puzzles

package main

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

