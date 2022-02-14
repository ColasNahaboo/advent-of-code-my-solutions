package main

import (
	"reflect"
	"testing"
)

// reflect.DeepEqual(got, expected)
func Test_ParseReindeers(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		got := ParseReindeers([]string{("Comet can fly 14 km/s for 10 seconds, but then must rest for 127 seconds.")})
		expected := []Reindeer{{"Comet", 14, 10, 127, 0, 0, 1, 10}}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
}

func Test_ReindeerDistance(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		got := ReindeerDistance(Reindeer{"Comet", 14, 10, 127, 0, 0, 1, 10}, 10)
		expected := 140
		if got != expected {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
	t.Run("2", func(t *testing.T) {
		got := ReindeerDistance(Reindeer{"Comet", 14, 10, 127, 0, 0, 1, 10}, 20)
		expected := 140
		if got != expected {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
	t.Run("3", func(t *testing.T) {
		got := ReindeerDistance(Reindeer{"Comet", 14, 10, 127, 0, 0, 1, 10}, 137)
		expected := 140
		if got != expected {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
	t.Run("Comet-1000", func(t *testing.T) {
		got := ReindeerDistance(Reindeer{"Comet", 14, 10, 127, 0, 0, 1, 10}, 1000)
		expected := 1120
		if got != expected {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
	t.Run("Dancer-1000", func(t *testing.T) {
		got := ReindeerDistance(Reindeer{"Dancer", 16, 11, 162, 0, 0, 1, 10}, 1000)
		expected := 1056
		if got != expected {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
}

func Test_WinnerScoreIncr(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		rs := []Reindeer{(Reindeer{"Comet", 14, 10, 127, 0, 0, 1, 10})}
		WinnerScoreIncr(rs, 0)
		got := rs[0]
		expected := Reindeer{"Comet", 14, 10, 127, 1, 14, 1, 10}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
}
