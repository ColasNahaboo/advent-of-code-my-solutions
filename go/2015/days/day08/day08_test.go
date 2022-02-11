package main

import (
	"fmt"
	"testing"
)

// reflect.DeepEqual(got, expected)
func Test_lengthsDec(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		arg := `""`
		expected := "2,0"
		d, m := lengthsDec(arg)
		got := fmt.Sprintf("%v,%v", d, m)
		if got != expected {
			t.Errorf("Arg '%v': expected '%v' but got '%v'", arg, expected, got)
		}
	})

	t.Run("2", func(t *testing.T) {
		arg := `"abc"`
		expected := "5,3"
		d, m := lengthsDec(arg)
		got := fmt.Sprintf("%v,%v", d, m)
		if got != expected {
			t.Errorf("Arg '%v': expected '%v' but got '%v'", arg, expected, got)
		}
	})

	t.Run("3", func(t *testing.T) {
		arg := `"aaa\"aaa"`
		expected := "10,7"
		d, m := lengthsDec(arg)
		got := fmt.Sprintf("%v,%v", d, m)
		if got != expected {
			t.Errorf("Arg '%v': expected '%v' but got '%v'", arg, expected, got)
		}
	})

	t.Run("4", func(t *testing.T) {
		arg := `"\x27"`
		expected := "6,1"
		d, m := lengthsDec(arg)
		got := fmt.Sprintf("%v,%v", d, m)
		if got != expected {
			t.Errorf("Arg '%v': expected '%v' but got '%v'", arg, expected, got)
		}
	})

	t.Run("5", func(t *testing.T) {
		arg := `"a\bc"`
		expected := "6,4"
		d, m := lengthsDec(arg)
		got := fmt.Sprintf("%v,%v", d, m)
		if got != expected {
			t.Errorf("Arg '%v': expected '%v' but got '%v'", arg, expected, got)
		}
	})

	t.Run("6", func(t *testing.T) {
		arg := `"\"pat\"\x63kpfc\"\x2ckhfvxk\"uwqzlx"`
		expected := "37,25"
		d, m := lengthsDec(arg)
		got := fmt.Sprintf("%v,%v", d, m)
		if got != expected {
			t.Errorf("Arg '%v': expected '%v' but got '%v'", arg, expected, got)
		}
	})

	t.Run("7", func(t *testing.T) {
		arg := `"a\xhi"`
		expected := "7,5"
		d, m := lengthsDec(arg)
		got := fmt.Sprintf("%v,%v", d, m)
		if got != expected {
			t.Errorf("Arg '%v': expected '%v' but got '%v'", arg, expected, got)
		}
	})

	t.Run("8", func(t *testing.T) {
		arg := `"a\x0"`
		expected := "6,4"
		d, m := lengthsDec(arg)
		got := fmt.Sprintf("%v,%v", d, m)
		if got != expected {
			t.Errorf("Arg '%v': expected '%v' but got '%v'", arg, expected, got)
		}
	})

	t.Run("9", func(t *testing.T) {
		arg := `"oh\\x\\h"`
		expected := "10,6"
		d, m := lengthsDec(arg)
		got := fmt.Sprintf("%v,%v", d, m)
		if got != expected {
			t.Errorf("Arg '%v': expected '%v' but got '%v'", arg, expected, got)
		}
	})

}
