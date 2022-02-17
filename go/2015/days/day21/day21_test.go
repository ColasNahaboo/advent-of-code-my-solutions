package main

import (
	"reflect"
	"testing"
)

// reflect.DeepEqual(got, expected)
func Test_day21(t *testing.T) {

	t.Run("read", func(t *testing.T) {
		got := readMob(fileToLines("input.txt"))
		expected := mob{109, 8, 2}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})

	t.Run("winFight", func(t *testing.T) {
		got := winFight(mob{8, 5, 5}, mob{12, 7, 2})
		expected := true
		if got != expected {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})

	t.Run("winFight", func(t *testing.T) {
		got := winFight(mob{100, 3, 8}, mob{109, 8, 2})
		expected := false
		if got != expected {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})

}
