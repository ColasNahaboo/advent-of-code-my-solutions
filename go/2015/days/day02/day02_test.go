package main

import (
	"reflect"
	"testing"
)

func TestReadInput(t *testing.T) {
	t.Run("Test ReadInput #1", func(t *testing.T) {
		got := ReadInput("small.txt")
		expected := []Present{Present{2, 3, 4}}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
	t.Run("Test ReadInput #2", func(t *testing.T) {
		got := ReadInput("exemple.txt")
		expected := []Present{Present{2, 3, 4}, Present{1, 1, 10}}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
}

func TestPart1(t *testing.T) {
	t.Run("Test Part1", func(t *testing.T) {
		got := Part1([]Present{Present{2, 3, 4}})
		expected := 58
		if got != expected {
			t.Errorf("expected '%d' but got '%d'", expected, got)
		}
	})
	t.Run("Test Part1", func(t *testing.T) {
		got := Part1([]Present{Present{1, 1, 10}})
		expected := 43
		if got != expected {
			t.Errorf("expected '%d' but got '%d'", expected, got)
		}
	})
}

func TestPart2(t *testing.T) {
	t.Run("Test Part2", func(t *testing.T) {
		got := Part2([]Present{Present{2, 3, 4}})
		expected := 34
		if got != expected {
			t.Errorf("expected '%d' but got '%d'", expected, got)
		}
	})
	t.Run("Test Part1", func(t *testing.T) {
		got := Part2([]Present{Present{1, 1, 10}})
		expected := 14
		if got != expected {
			t.Errorf("expected '%d' but got '%d'", expected, got)
		}
	})
}
