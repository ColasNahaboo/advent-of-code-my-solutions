package main

import (
	"testing"
	"reflect"
)

/////////// Setups

/////////// Tests to pass

func Test(t *testing.T) {
	//  calls to test function(s) args (adapt "args" and "expected")
	t_1("1.1", t, args, expected)
	t_1("1.2", t, args, expected)
	t_1("1.3", t, args, expected)
}

/////////// Functions to test

// adapt to your case the types of args and expected
func t_1(label string, t *testing.T, args string, expected int) {
	// setup here
	t.Run(label, func(t *testing.T) {
		got := theFunctionToTest(args)
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("exp '%v'\ngot '%v'", expected, got)
		}})}
