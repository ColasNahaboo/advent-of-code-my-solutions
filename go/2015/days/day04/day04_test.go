package main

import (
	"testing"
)

func TestPart1(t *testing.T) {
	t.Run("Test Part1.1", func(t *testing.T) {
		got := Part1([]byte("abcdef"))
		expected := 609043
		if got != expected {
			t.Errorf("expected '%d' but got '%d'", expected, got)
		}
	})
	t.Run("Test Part1.2", func(t *testing.T) {
		got := Part1([]byte("pqrstuv"))
		expected := 1048970
		if got != expected {
			t.Errorf("expected '%d' but got '%d'", expected, got)
		}
	})
}
