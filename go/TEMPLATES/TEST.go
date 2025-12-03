package main

import (
	"testing"
	// "reflect"
)

/////////// Setups

/////////// Tests to pass

func Test(t *testing.T) {
	TEST("My Test")
	//  calls to test function(s) args (adapt "args" and "expected")
	T(t, args, expected)
	T(t, args, expected)
	// ...
}

/////////// Functions to test

// adapt to your case the number and types of args and expected
func T(t *testing.T, args string, expected int) {
	TNUM++
	// your optional setup here
	t.Run(TNAME + "#" + itoa(TNUM), func(t *testing.T) {
		got := theFunctionToTest(args)
		// or: if !reflect.DeepEqual(got, expected) {
		if got != expect {
			t.Errorf("\nexp= '%v'\ngot= '%v'", expected, got)
		}}
	)
}

//////////// Automatic Labels: string-argument-to-TEST + # + test-number

var TNAME = "T"
var TNUM = 0

func TEST(s string) {
	TNAME = s
	TNUM = 0
}
