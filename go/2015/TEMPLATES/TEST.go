package main

import (
	"testing"
)

func TestPart1(t *testing.T) {
	t.Run("Test Part1.1", func(t *testing.T) {
		got := Part1([]byte(">"))
		expected := 2
		if got != expected {
			t.Errorf("expected '%d' but got '%d'", expected, got)
		}
	})
}

func TestPart2(t *testing.T) {
	t.Run("Test Part2.1", func(t *testing.T) {
		got := Part2([]byte("^v"))
		expected := 3
		if got != expected {
			t.Errorf("expected '%d' but got '%d'", expected, got)
		}
	})
}
