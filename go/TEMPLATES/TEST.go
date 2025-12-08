package main
###*** THIS IS A NON-INSTANCIED TEMPLATE, IT DOES NOT COMPILE ***###

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

// Wapper test functions calling the existing functions to test
func T(t *testing.T, args string, expected int) {
	TNUM++
	// optional setup
	t.Run(TNAME + "#" + itoa(TNUM), func(t *testing.T) {
		// run the test function
		got := theFunctionToTest(args)
		// compare result with expected
		if got != expect {	  // or use: if !reflect.DeepEqual(got, expected) {
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
