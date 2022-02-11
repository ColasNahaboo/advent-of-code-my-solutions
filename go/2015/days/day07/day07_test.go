package main

import (
	"testing"
)

// reflect.DeepEqual(got, expected)
func Test_Part1(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		got := int(Part1([]string{"3 -> a"}))
		expected := 3
		if got != expected {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
	t.Run("1", func(t *testing.T) {
		got := int(Part1([]string{"3 -> x", "6 -> y", "x AND y -> a"}))
		expected := 2
		if got != expected {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
}
