package main

import (
	"reflect"
	"testing"
)

// reflect.DeepEqual(got, expected)
func Test_Part1(t *testing.T) {
	t.Run("Test Part1.1", func(t *testing.T) {
		got := Part1([]byte(">"))
		expected := 2
		if got != expected {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
}

func Test_Part2(t *testing.T) {
	t.Run("Test Part2.1", func(t *testing.T) {
		got := Part1([]byte(">"))
		expected := 2
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
}
