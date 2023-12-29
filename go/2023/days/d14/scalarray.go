// A "Scalar Array" package, which stores a 2D array as a one-dimensional slice
// We call position the index in this slice (a number), and it is
// easier to manage than the (x, y) virtual 2D array coordinates.
// A point a (x, y) in the scalarray being at position p = x + y * width
//
// #..#..#   Here @ is at 2D coords (4, 1),
// #...@#.   but at scalar position 11
// #.#....   in the slice: #..#..##...@#.#.#....
// Coordinates start at 0 (top left) and increase downwards and to the right

// This system is quite useful for managing problems of the "adventofcode",
// and more generally fixed size 2-dimensional arrays.
// The type is not opaque on purpose for simplicity and efficieny, feel free to
// directly use the fields of the types.

// You can specify to add a border of width b around the actual data
// in this case (x,y) is at position x+b + (y+b) * (width+2*b)
// It is often easier to add "walls" around a 2D board rather to always check
// if the coordinates stay in the board space
// No work with the original coordinates before theborder was added use the
// methods ending in B such as GetB instead of Get

// Scalarray is the default type: 2D array with an optional border
// Scalarray3D is the same concept, but for fixed size 3 dimensions arrays

// Scallaray4D, ScallarayN... may be implemented later for N-dimensional arrays

package main

import "fmt"

//////////// Scalarray, 2D with optional borders

type Scalarray[T comparable] struct {
	w, h, b int					// width, height and borderwidth of the array
	a []T						// the array elements in a slice
}

// by default, stay simple, do not add a border
func makeScalarray[T comparable](w, h int) (sa Scalarray[T]) {
	sa.w = w
	sa.h = h
	sa.a = make([]T, w*h, w*h)
	return
}

// the following functions work on the whole array, including the border
func (sa *Scalarray[T]) Set(x, y int, v T) {
	sa.a[x + y * sa.w] = v
}

func (sa *Scalarray[T]) Get(x, y int) T {
	return sa.a[x + y * sa.w]
}

func (sa *Scalarray[T]) Pos(x, y int) int {
	return x + y * sa.w
}

func (sa *Scalarray[T]) X(pos int) int {
	return pos % sa.w
}

func (sa *Scalarray[T]) Y(pos int) int {
	return pos / sa.w
}

func (sa *Scalarray[T]) Coords(pos int) (x, y int) {
	x = pos % sa.w
	y = pos / sa.w
	return
}

func (sa *Scalarray[T]) Vector(pos int) (v [2]int) {
	v[0] = pos % sa.w
	v[1] = pos / sa.w
	return
}

func  (sa *Scalarray[T]) isValid(pos int) bool {
	return pos >= 0 && pos < sa.w * sa.h
}

// returns array of position offsets for going Up Right Down Left
func (sa *Scalarray[T]) Dirs(pos int) [4]int {
	return [4]int{-sa.w, 1, sa.w, -1}
}

// Clone also the underlying array
func (sa *Scalarray[T]) Clone() (sc Scalarray[T]) {
	sc.w = sa.w
	sc.h = sa.h
	sc.b = sa.b
	sc.a = make([]T, len(sa.a), len(sa.a))
	copy(sc.a, sa.a)
	return
}

func (s1 *Scalarray[T]) Equal(s2 Scalarray[T]) bool {
	if s1.w != s2.w || s1.h != s2.h {
		return false	   // we dont care about border size, only its contents
	}
	for i, v := range s1.a {
		if s2.a[i] != v {
			return false
		}
	}
	return true
}

//////////// Convenience 2D functions to manage a border, names end in B

// make room for the border, but do not draw it (yet)
func makeScalarrayB[T comparable](w, h, b int) (sa Scalarray[T]) {
	sa.w = w + 2*b
	sa.h = h + 2*b
	sa.b = b
	sa.a = make([]T, sa.w*sa.h, sa.w*sa.h)
	return
}

// draw the border in an existing Scalarray by setting its elements to v
func (sa *Scalarray[T]) DrawBorder(v T) {
	for b := 0; b < sa.b; b++ {	// draw border rows, starting from outside
		for x := 0; x < sa.w; x++ {
			sa.a[x + b*sa.w] = v			// top borders
			sa.a[x + (sa.h - b - 1) * sa.w] = v // bottom borders
		}
		for y := sa.b; y < sa.h - sa.b; y++ {
			sa.a[b + y*sa.w] = v			// left borders
			sa.a[sa.w - b - 1 + y*sa.w] = v // right borders
		}
	}
}

// These methods convert from/to relative coords (x, y) relative to the inside
// of the border to Scalarray positions including the border
func (sa *Scalarray[T]) SetB(x, y int, v T) {
	sa.a[x + sa.b + (y+sa.b) * sa.w] = v
}

func (sa *Scalarray[T]) GetB(x, y int) T {
	return sa.a[x + sa.b + (y+sa.b) * sa.w]
}

func (sa *Scalarray[T]) PosB(x, y int) int {
	return x + sa.b + (y+sa.b) * sa.w
}

func (sa *Scalarray[T]) XB(pos int) int {
	return pos % sa.w - sa.b
}

func (sa *Scalarray[T]) YB(pos int) int {
	return pos / sa.w - sa.b
}

func (sa *Scalarray[T]) CoordsB(pos int) (x, y int) {
	x = pos % sa.w - sa.b
	y = pos / sa.w - sa.b
	return
}

func (sa *Scalarray[T]) VectorB(pos int) (v [2]int) {
	v[0] = pos % sa.w - sa.b
	v[1] = pos / sa.w - sa.b
	return
}

//////////// Specialized 2D PrettyPrinting & Debugging functions

// Array of boolean values, true = '#', false = '.'
func VPScallarrayBool(label string, sa Scalarray[bool]) {
	if ! verbose {
		return
	}
	fmt.Printf("%s: array %dx%d, with border %d:", label, sa.w, sa.h, sa.b)
	for p, b := range sa.a {
		if p % sa.w == 0 {
			fmt.Println()
		}
		if b {
			fmt.Print("#")
		} else {
			fmt.Print(".")
		}
	}
	fmt.Println()
}

// Array of integers. We set the cell display width to the largest number
func VPScallaryInt(label string, sa Scalarray[int]) {
	if ! verbose {
		return
	}
	fmt.Printf("%s: array %dx%d, with border %d:",label, sa.w, sa.h, sa.b)
	max := 0
	for _, i := range sa.a {
		if i > max {
			max = i
		} else if -1 > max {
			max = -i
		}
	}
	maxs := itoa(max)
	cellformat := "%" + itoa(len(maxs) + 1) + "d"
	for p, i := range sa.a {
		if p % sa.w == 0 {
			fmt.Println()
		}
		fmt.Printf(cellformat, i)
	}
	fmt.Println()
}

//////////// Scallarray3D
// Minimal support for now: no borders.

type Scalarray3D[T comparable] struct {
	w, h, d int					// width, height and depth of the array
	a []T						// the array elements in a slice
}

func makeScalarray3D[T comparable](w, h, d int) (sa Scalarray3D[T]) {
	sa.w = w
	sa.h = h
	sa.d = d
	sa.a = make([]T, w*h*d, w*h*d)
	return
}

func (sa *Scalarray3D[T]) Set(x, y, z int, v T) {
	sa.a[x + y * sa.w + z * sa.w*sa.h] = v
}

func (sa *Scalarray3D[T]) Get(x, y, z int) T {
	return sa.a[x + y * sa.w + z * sa.w*sa.h]
}

func (sa *Scalarray3D[T]) Pos(x, y, z int) int {
	return x + y * sa.w + z * sa.w*sa.h
}

func (sa *Scalarray3D[T]) X(pos int) int {
	return pos % sa.w
}

func (sa *Scalarray3D[T]) Y(pos int) int {
	return (pos % (sa.w*sa.h)) / sa.w
}

func (sa *Scalarray3D[T]) Z(pos int) int {
	return pos / (sa.w*sa.h)
}

func (sa *Scalarray3D[T]) Coords(pos int) (x, y, z int) {
	x = sa.X(pos)
	y = sa.Y(pos)
	z = sa.Z(pos)
	return
}

func (sa *Scalarray3D[T]) Vector(pos int) [3]int {
	return [3]int{sa.X(pos), sa.Y(pos), sa.Z(pos)}
}

func  (sa *Scalarray3D[T]) isValid(pos int) bool {
	return pos >= 0 && pos < sa.w * sa.h * sa.d
}

// returns array of position offsets for going Up Right Down Left Front Back
func (sa *Scalarray3D[T]) Dirs(pos int) [6]int {
	return [6]int{-sa.w, 1, sa.w, -1, sa.w*sa.h, -sa.w*sa.h}
}

// Clone also the underlying array
func (sa *Scalarray3D[T]) Clone() (sc Scalarray3D[T]) {
	sc.w = sa.w
	sc.h = sa.h
	sc.d = sa.d
	sc.a = make([]T, len(sa.a), len(sa.a))
	copy(sc.a, sa.a)
	return
}

func (s1 *Scalarray3D[T]) Equal(s2 Scalarray3D[T]) bool {
	if s1.w != s2.w || s1.h != s2.h || s1.d != s2.d {
		return false
	}
	for i, v := range s1.a {
		if s2.a[i] != v {
			return false
		}
	}
	return true
}
