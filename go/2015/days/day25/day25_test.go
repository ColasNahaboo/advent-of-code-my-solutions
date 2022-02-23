package main

import (
	"testing"
)

// reflect.DeepEqual(got, expected)
func Test_numAt(t *testing.T) {
	t.Run("1,1", func(t *testing.T) {
		got := numAt(1, 1)
		expected := 1
		if got != expected {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
	t.Run("6,1", func(t *testing.T) {
		got := numAt(6, 1)
		expected := 21
		if got != expected {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
	t.Run("1,6", func(t *testing.T) {
		got := numAt(1, 6)
		expected := 16
		if got != expected {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
	t.Run("3,3", func(t *testing.T) {
		got := numAt(3, 3)
		expected := 13
		if got != expected {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
}
