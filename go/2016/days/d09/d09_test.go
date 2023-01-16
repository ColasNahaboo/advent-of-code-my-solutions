package main

import (
	"testing"
)

func Test_1(t *testing.T) {
	// setup code

	// table-driven tests: calls to test function(s) with various args
	// adapt args and expected as needed
	t_dc("dc1", t, "ADVENT", "ADVENT")
	t_dc("dc2", t, "A(1x5)BC", "ABBBBBC")
	t_dc("dc3", t, "(3x3)XYZ", "XYZXYZXYZ")
	t_dc("dc4", t, "A(2x2)BCD(2x2)EFG", "ABCBCDEFEFG")
	t_dc("dc5", t, "(6x1)(1x3)A", "(1x3)A")
	t_dc("dc6", t, "X(8x2)(3x3)ABCY", "X(3x3)ABC(3x3)ABCY")
	// cleanup code
}

// test function
func t_dc(label string, t *testing.T, arg string, expected string) {
	// setup here
	t.Run(label, func(t *testing.T) {
		got := decompress(arg)
		if got != expected {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
}

func Test_2(t *testing.T) {
	// setup code

	// table-driven tests: calls to test function(s) with various args
	// adapt args and expected as needed
	t_est("est1", t, "(3x3)XYZ", 9)
	t_est("est2", t, "X(8x2)(3x3)ABCY", 20)
	t_est("est3", t, "(27x12)(20x12)(13x14)(7x10)(1x12)A", 241920)
	t_est("est4", t, "(25x3)(3x3)ABC(2x3)XY(5x2)PQRSTX(18x9)(3x2)TWO(5x7)SEVEN", 445)
	// cleanup code
}

// test function
func t_est(label string, t *testing.T, arg string, expected int) {
	// setup here
	t.Run(label, func(t *testing.T) {
		got := estimate(arg)
		if got != expected {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
}
