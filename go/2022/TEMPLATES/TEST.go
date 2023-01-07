package main

import (
	"reflect"
	"testing"
)

// reflect.DeepEqual(got, expected)
func Test_1(t *testing.T) {
	// setup code

	// table-driven tests: calls to test function(s) with various args
	// adapt args and expected as needed
	t_1("1.1", t, args, expected)
	t_1("1.2", t, args, expected)
	t_1("1.3", t, args, expected)

	// cleanup code
}

// test function
func t_wc(label string, t *testing.T, args, expected int) {
	// setup here
	t.Run(label, func(t *testing.T) {
		got := theFunctionToTest(args)
		if reflect.DeepEqual(got, expected) {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
}
