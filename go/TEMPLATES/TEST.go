package main

import (
	"testing"
	"reflect"
)

// set of unit tests
func Test(t *testing.T) {
	// setup code

	//  calls to test function(s) args (adapt "args" and "expected")
	t_1("1.1", t, args, expected)
	t_1("1.2", t, args, expected)
	t_1("1.3", t, args, expected)

	// cleanup code
}

// test function
func t_1(label string, t *testing.T, args, expected int) {
	// setup here
	t.Run(label, func(t *testing.T) {
		got := theFunctionToTest(args)
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}})}
