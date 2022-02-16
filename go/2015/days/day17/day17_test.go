package main

import (
	"testing"
)

// reflect.DeepEqual(got, expected)
func Test_part1(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		cans := readCans(fileToLines("exemple.txt"))
		got := part1(cans, 25)
		expected := 4
		if got != expected {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
}

func Test_part2(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		cans := readCans(fileToLines("exemple.txt"))
		got := part2(cans, 25)
		expected := 3
		if got != expected {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
}
