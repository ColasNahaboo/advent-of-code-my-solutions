package main

import (
	"testing"
	//"reflect"
)

/////////// Setups

var locs = []Point{Point{1,1}, Point{1,6}, Point{8,3}, Point{3,4}, Point{5,5}, Point{8,9}}

/////////// Tests to pass

func Test(t *testing.T) {
	//  calls to test function(s) args (adapt "args" and "expected")
	t_ClosestXY("Closest-1", t, locs, 2, 2, 0)
}

/////////// Functions to test

func t_ClosestXY(label string, t *testing.T, locs []Point, x, y int, expected int) {
	t.Run(label, func(t *testing.T) {
		got := ClosestXY(locs, x, y)
		if got != expected {
			t.Errorf("exp '%v'\ngot '%v'", expected, got)
		}})}
