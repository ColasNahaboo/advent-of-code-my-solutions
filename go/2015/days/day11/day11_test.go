package main

import (
	"reflect"
	"testing"
)

func Test_SeqToString(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		got := SeqToString([]int{5, 14, 14, 1, 0, 17})
		expected := "foobar"
		if got != expected {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
}

func Test_StringToSeq(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		got := StringToSeq("foobar")
		expected := []int{5, 14, 14, 1, 0, 17}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
}

func Test_SeqIncPos(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		got := SeqIncPos([]int{5, 14, 14, 1, 0, 17}, 5)
		expected := []int{5, 14, 14, 1, 0, 18}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
	t.Run("2", func(t *testing.T) {
		got := SeqIncPos([]int{5, 25, 14}, 1)
		expected := []int{6, 0, 14}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
	t.Run("3", func(t *testing.T) {
		got := SeqIncPos([]int{7, 13, 10}, 2)
		expected := []int{7, 13, 12}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
}

func Test_PassIsValid(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		got := PassIsValid([]int{5, 14, 14, 1, 0, 17})
		expected := false
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
	t.Run("2", func(t *testing.T) {
		got := PassIsValid([]int{5, 1, 0, 17})
		expected := false
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
	t.Run("3", func(t *testing.T) {
		got := PassIsValid(StringToSeq("aabcc"))
		expected := true
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
	t.Run("ex1", func(t *testing.T) {
		got := PassIsValid(StringToSeq("abceeaaa"))
		expected := true
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
	t.Run("ex2", func(t *testing.T) {
		got := PassIsValid(StringToSeq("ghikkaaa"))
		expected := true
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
}
