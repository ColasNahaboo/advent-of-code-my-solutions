package main

import (
	"testing"
)

// reflect.DeepEqual(got, expected)
func Test_TurnOn(t *testing.T) {
	t.Run("Test TurnOn.1", func(t *testing.T) {
		var grid Grid
		got := grid.TurnOn(0, 0, 0, 2, 2)
		expected := 9
		if got != expected {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
}

func Test_TurnOn2(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		var grid Grid
		got := grid.TurnOn2(1, 0, 0, 0, 0)
		expected := 2
		if got != expected {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
	t.Run("2", func(t *testing.T) {
		var grid Grid
		got := grid.TurnOn2(4, 0, 0, 999, 999)
		expected := 1000004
		if got != expected {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
}

func Test_Toggle2(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		var grid Grid
		got := grid.Toggle2(4, 0, 0, 999, 999)
		expected := 2000004
		if got != expected {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
}
