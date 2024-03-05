package main

import (
        "testing"
)

//------------------------------------------------------------------

func Test_PatternTransform(t *testing.T) {
	t_PT("size2 r90", t, "12/34", []int{1, 4, 2, 0, 3}, "31/42")
	t_PT("size2 fh", t, "12/34", []int{1, 0, 2, 4, 3}, "21/43")
	t_PT("size2 fv", t, "12/34", []int{3, 4, 2, 0, 1}, "34/12")

	t_PT("size3 r90", t, "012/456/89a", []int{1, 2, 6, 3, 0, 5, 10, 7, 4, 8, 9}, "401/852/9a6")
	t_PT("size3 fh", t, "012/456/89a", []int{2, 1, 0, 3, 6, 5, 4, 7, 10, 9, 8}, "210/654/a98")
	t_PT("size3 fv", t, "012/456/89a", []int{8, 9, 10, 3, 4, 5, 6, 7, 0, 1, 2}, "89a/456/012")
}
		
func t_PT(label string, t *testing.T, s string, transform []int, expected string) {
	// setup here
	t.Run(label, func(t *testing.T) {
		got := PatternTransform(s, transform...)
		if got != expected {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
}

//------------------------------------------------------------------

func Test_Extract(t *testing.T) {
	i1 := MakeImage("#..#/..../..../#..#")     // #..#
	t_Extract("i1q1", t, 2, 0, 0, i1, "#./..") // ....
	t_Extract("i1q1", t, 2, 2, 0, i1, ".#/..") // ....
	t_Extract("i1q1", t, 2, 0, 2, i1, "../#.") // #..#
	t_Extract("i1q1", t, 2, 2, 2, i1, "../.#")
	t_Extract("i1q1", t, 2, 1, 1, i1, "../..")
}

func t_Extract(label string, t *testing.T, size, x, y int, i Image, e string) {
	// setup here
	t.Run(label, func(t *testing.T) {
		lit, got := i.Extract(size, x, y)
		if string(got) != e {
			t.Errorf("expected '%v' but got '%v'", e, got)
		}
		if lit != PatternLit(e) {
			t.Errorf("expected lit '%d' but got '%d'", PatternLit(e), lit)
		}
	})
}
