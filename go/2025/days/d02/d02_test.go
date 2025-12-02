package main

import (
	"testing"
	//	"reflect"					
)

/////////// Setups

/////////// Tests to pass

func Test(t *testing.T) {
	//  calls to test function(s) args (adapt "args" and "expected")
	t_rs("1", t, 1, false)
	t_rs("123", t, 123, false)
	t_rs("55", t, 55, true)
	t_rs("6464", t, 6464, true)
	t_rs("123123", t, 123123, true)
	t_rs("1111", t, 1111, true)
	t_rs("121212", t, 121212, true)
	t_rs("121213", t, 121213, false)
	t_rs("1188511885", t, 1188511885, true)
}

/////////// Functions to test

func t_rs(label string, t *testing.T, args int, expected bool) {
	// setup here
	t.Run(label, func(t *testing.T) {
		got := invalidRepeatedSequence(args)
		if got != expected {
			t.Errorf("exp '%v'\ngot '%v'", expected, got)
		}})}
