// A grid package, where a 2D grid is stored as a one-dimensional slice
// We call position the index in this slice (a number)
// Easier to manage than the (x, y) virtual 2D grid coordinates
// A point a (x, y) in the grid being at position p = x + y * grid_width

// TODO: border option

package main

import "fmt"


//////////// Grid

type Grid[T any] struct {
	w, h int					//  width and height of the grid
	g []T
}

func makeGrid[T any](w, h int) (grid Grid[T]) {
	grid.w = w
	grid.h = h
	grid.g = make([]T, w*h, w*h)
	return
}

func (grid *Grid[T]) Set(x, y int, v T) {
	grid.g[x + y * grid.w] = v
}

func (grid *Grid[T]) Get(x, y int) T {
	return grid.g[x + y * grid.w]
}

func (grid *Grid[T]) Pos(x, y int) int {
	return x + y * grid.w
}

func (grid *Grid[T]) Coords(pos int) (x, y int) {
	x = pos % grid.w
	y = pos / grid.w
	return
}

// returns array of position offsets for going Up Right Down Left
func (grid *Grid[T]) Dirs(pos int) [4]int {
	return [4]int{-grid.w, 1, grid.w, -1}
}

//////////// PrettyPrinting & Debugging functions

func VPgridBool(label string, g Grid[bool]) {
	if ! verbose {
		return
	}
	fmt.Printf("%s: grid %dx%d:",label, g.w, g.h)
	for p, b := range g.g {
		if p % g.w == 0 {
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

func VPgridInt(label string, g Grid[int]) {
	if ! verbose {
		return
	}
	fmt.Printf("%s: grid %dx%d:",label, g.w, g.h)
	max := 0
	for _, i := range g.g {
		if i > max {
			max = i
		} else if -1 > max {
			max = -i
		}
	}
	maxs := itoa(max)
	cellformat := "%" + itoa(len(maxs) + 1) + "d"
	for p, i := range g.g {
		if p % g.w == 0 {
			fmt.Println()
		}
		fmt.Printf(cellformat, i)
	}
	fmt.Println()
}
