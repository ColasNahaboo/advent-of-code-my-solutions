package main

import (
	"testing"
)

func TestPart1(t *testing.T) {
	t.Run("Test Part1 with valid list", func(t *testing.T) {
		got := Part1(")())())")
		expected := -3
		if got != expected {
			t.Errorf("expected '%d' but got '%d'", expected, got)
		}
	})
}

func TestPart2(t *testing.T) {
	t.Run("Test Part2 with valid list", func(t *testing.T) {
		got := Part2(")())())")
		expected := 1
		if got != expected {
			t.Errorf("expected '%d' but got '%d'", expected, got)
		}
	})
}
