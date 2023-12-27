// A "Scalar Array" package, which stores a 2D array as a one-dimensional slice
// We call position the index in this slice (a number)
// Easier to manage than the (x, y) virtual 2D array coordinates
// A point a (x, y) in the scalarray being at position p = x + y * width
//
// #..#..#   Here @ is at 2D coords (4, 1),
// #...@#.   but at scalar position 11
// #.#....   in the slice: #..#..##...@#.#.#....

// You can specify to add a border of width b around the actual data
// in this case (x,y) is at position x+b + (y+b) * (width+2*b)

// Scalarray is the default type: 2D array with an optional border
// Scalarray0 is a simplified version without borders
// Scalarray2 may be used internally for things common to Scalarray & Scalarray0
// Scallarray3, Scallaray4, ... may be implemented later nor N-dimensional arrays

package main

import "fmt"

//////////// Scalarray, 2D with optional borders

type Scalarray[T comparable] struct {
	w, h, b int					// width and height and border of the array
	a []T						// the array elements in a slice
}

func makeScalarray[T comparable](w, h, b int) (sa Scalarray[T]) {
	sa.w = w + 2*b
	sa.h = h + 2*b
	sa.b = b
	sa.a = make([]T, sa.w*sa.h, sa.w*sa.h)
	return
}

func (sa *Scalarray[T]) Set(x, y int, v T) {
	sa.a[x + sa.b + (y+sa.b) * sa.w] = v
}

func (sa *Scalarray[T]) Get(x, y int) T {
	return sa.a[x + sa.b + (y+sa.b) * sa.w]
}

func (sa *Scalarray[T]) Pos(x, y int) int {
	return x + sa.b + (y+sa.b) * sa.w
}

func (sa *Scalarray[T]) Coords(pos int) (x, y int) {
	x = pos % sa.w - sa.b
	y = pos / sa.w - sa.b
	return
}

func (sa *Scalarray[T]) Vector(pos int) (v [2]int) {
	v[0] = pos % sa.w - sa.b
	v[1] = pos / sa.w - sa.b
	return
}

// returns array of position offsets for going Up Right Down Left
func (sa *Scalarray[T]) Dirs(pos int) [4]int {
	return [4]int{-sa.w, 1, sa.w, -1}
}

// draw the border by setting it to v
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


//////////// Scalarray0, 2D with no borders, a bit simpler and faster

type Scalarray0[T comparable] struct {
	w, h int					// width and height of the array
	a []T						// the array elements in a slice
}

func makeScalarray0[T comparable](w, h int) (sa Scalarray0[T]) {
	sa.w = w
	sa.h = h
	sa.a = make([]T, w*h, w*h)
	return
}

func (sa *Scalarray0[T]) Set(x, y int, v T) {
	sa.a[x + y * sa.w] = v
}

func (sa *Scalarray0[T]) Get(x, y int) T {
	return sa.a[x + y * sa.w]
}

func (sa *Scalarray0[T]) Pos(x, y int) int {
	return x + y * sa.w
}

func (sa *Scalarray0[T]) Coords(pos int) (x, y int) {
	x = pos % sa.w
	y = pos / sa.w
	return
}

func (sa *Scalarray0[T]) Vector(pos int) (v [2]int) {
	v[0] = pos % sa.w
	v[1] = pos / sa.w
	return
}

// returns array of position offsets for going Up Right Down Left
func (sa *Scalarray0[T]) Dirs(pos int) [4]int {
	return [4]int{-sa.w, 1, sa.w, -1}
}


//////////// PrettyPrinting & Debugging functions

// Array of boolean values, true = '#', false = '.'
func VPScallarrayBool(label string, sa Scalarray[bool]) {
	if ! verbose {
		return
	}
	fmt.Printf("%s: array %dx%d, with border %d:", label, sa.w, sa.h, sa.b)
	VPScallarray2Bool(label, sa.w, sa.h, sa.a)
}
// Same for Scalarray0
func VPScallarray0Bool(label string, sa Scalarray0[bool]) {
	if ! verbose {
		return
	}
	fmt.Printf("%s: array %dx%d:", label, sa.w, sa.h)
	VPScallarray2Bool(label, sa.w, sa.h, sa.a)
}
func VPScallarray2Bool(label string, saw, sah int, saa []bool) {
	for p, b := range saa {
		if p % saw == 0 {
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

// Array of integers. We set the cell width to the largest number
func VPScallaryInt(label string, sa Scalarray[int]) {
	if ! verbose {
		return
	}
	fmt.Printf("%s: array %dx%d, with border %d:",label, sa.w, sa.h, sa.b)
	VPScallarray2Int(label, sa.w, sa.h, sa.a)
}
// Same for Scallaray0
func VPScallary0Int(label string, sa Scalarray[int]) {
	if ! verbose {
		return
	}
	fmt.Printf("%s: array %dx%d:",label, sa.w, sa.h)
	VPScallarray2Int(label, sa.w, sa.h, sa.a)
}
func VPScallarray2Int(label string, saw, sah int, saa []int) {
	max := 0
	for _, i := range saa {
		if i > max {
			max = i
		} else if -1 > max {
			max = -i
		}
	}
	maxs := itoa(max)
	cellformat := "%" + itoa(len(maxs) + 1) + "d"
	for p, i := range saa {
		if p % saw == 0 {
			fmt.Println()
		}
		fmt.Printf(cellformat, i)
	}
	fmt.Println()
}
