package main

import (
	"testing"
	"reflect"
)

/////////// Setups

/////////// Tests to pass

func Test(t *testing.T) {
	TEST("runion")
	//  calls to test function(s) args (adapt "args" and "expected")
	T(t, []Range{{1, 2}, {4, 5}}, []Range{{1, 2}, {4, 5}})
	T(t, []Range{{4, 5}, {1, 2}}, []Range{{1, 2}, {4, 5}})
	T(t, []Range{{1, 2}, {2, 3}}, []Range{{1, 3}})
	T(t, []Range{{2, 3}, {1, 2}}, []Range{{1, 3}})
	T(t, []Range{{1, 2}, {3, 4}}, []Range{{1, 2}, {3, 4}})
	T(t, []Range{{3, 4}, {1, 2}}, []Range{{1, 2}, {3, 4}})
	T(t, []Range{{1, 10}, {5, 15}}, []Range{{1, 15}})
	T(t, []Range{{5, 15}, {1, 10}}, []Range{{1, 15}})
	T(t, []Range{{1, 6}, {2, 4}}, []Range{{1, 6}})
	T(t, []Range{{2, 4}, {1, 6}}, []Range{{1, 6}})
	T(t, []Range{{1, 8}, {1, 4}}, []Range{{1, 8}})
	T(t, []Range{{1, 4}, {1, 8}}, []Range{{1, 8}})
	T(t, []Range{{1, 7}, {3, 7}}, []Range{{1, 7}})
	T(t, []Range{{3, 7}, {1, 7}}, []Range{{1, 7}})

	TEST("rappend")
	TA(t, []Range{{1, 3}}, Range{2, 5}, []Range{{1, 5}})
	TA(t, []Range{{4, 6}}, Range{2, 5}, []Range{{2, 6}})
	TA(t, []Range{{1, 3}, {4, 6}}, Range{2, 5}, []Range{{1, 6}})
}

/////////// Functions to test

// adapt to your case the number and types of args and expected
func T(t *testing.T, args []Range, expected []Range) {
	TNUM++
	// your optional setup here
	t.Run(TNAME + " " + ranges2string(args), func(t *testing.T) {
		got := runion(args)
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("\nexp= '%v'\ngot= '%v'", expected, got)
		}})
}

func TA(t *testing.T, args []Range, a Range, expected []Range) {
	TNUM++
	// your optional setup here
	t.Run(TNAME + " " + ranges2string(args) +  " + " + range2string(a), func(t *testing.T) {
		got := rappend(args, a)
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("\nexp= '%v'\ngot= '%v'", expected, got)
		}})
}

//////////// Automatic Labels: string-argument-to-TEST + # + test-number

var TNAME = "T"
var TNUM = 0

func TEST(s string) {
	TNAME = s
	TNUM = 0
}
