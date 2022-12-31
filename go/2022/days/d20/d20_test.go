package main

import (
	"reflect"
	"testing"
)

// Initial arrangement:
// 1, 2, -3, 3, -2, 0, 4
// 
// 1 moves between 2 and -3:
// 2, 1, -3, 3, -2, 0, 4
// 
// 2 moves between -3 and 3:
// 1, -3, 2, 3, -2, 0, 4
// 
// -3 moves between -2 and 0:
// 1, 2, 3, -2, -3, 0, 4
// 
// 3 moves between 0 and 4:
// 1, 2, -2, -3, 0, 3, 4
// 
// -2 moves between 4 and 1:
// 1, 2, -3, 0, 3, 4, -2
// 
// 0 does not move:
// 1, 2, -3, 0, 3, 4, -2
// 
// 4 moves between -3 and 0:
// 1, 2, -3, 4, 0, 3, -2

// reflect.DeepEqual(got, expected)
func Test_move(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		data := []int{1,2,-3,3,-2,0,4}
		got := move(data, 1)
		expected := []int{2,1,-3,3,-2,0,4}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}})
	t.Run("2", func(t *testing.T) {
		data := []int{2,1,-3,3,-2,0,4}
		got := move(data, 2)
		expected := []int{1, -3, 2, 3, -2, 0, 4}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}})
	t.Run("-3", func(t *testing.T) {
		data := []int{1, -3, 2, 3, -2, 0, 4}
		got := move(data, -3)
		expected := []int{1, 2, 3, -2, -3, 0, 4}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}})
	t.Run("3", func(t *testing.T) {
		data := []int{1, 2, 3, -2, -3, 0, 4}
		got := move(data, 3)
		expected := []int{1, 2, -2, -3, 0, 3, 4}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}})
	t.Run("-2", func(t *testing.T) {
		data := []int{1, 2, -2, -3, 0, 3, 4}
		got := move(data, -2)
		expected := []int{1, 2, -3, 0, 3, 4, -2}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}})
	t.Run("0", func(t *testing.T) {
		data := []int{1, 2, -3, 0, 3, 4, -2}
		got := move(data, 0)
		expected := []int{1, 2, -3, 0, 3, 4, -2}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}})
	t.Run("4", func(t *testing.T) {
		data := []int{1, 2, -3, 0, 3, 4, -2}
		got := move(data, 4)
		expected := []int{1, 2, -3, 4, 0, 3, -2}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}})
}
