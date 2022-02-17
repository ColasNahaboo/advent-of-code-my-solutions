package main

import (
	"testing"
)

// reflect.DeepEqual(got, expected)
func Test_part1(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		got := part1(150)
		expected := 8
		if got != expected {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
	t.Run("2", func(t *testing.T) {
		got := part1(29000000)
		expected := 665280
		if got != expected {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
	t.Run("3", func(t *testing.T) {
		got := part1(36000000)
		expected := 831600
		if got != expected {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
}

func Test_part2(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		got := part2(150)
		expected := 8
		if got != expected {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
	t.Run("2", func(t *testing.T) {
		got := part2(29000000)
		expected := 705600
		if got != expected {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
	t.Run("3", func(t *testing.T) {
		got := part2(36000000)
		expected := 884520
		if got != expected {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
}
