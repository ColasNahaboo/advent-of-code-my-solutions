package main

import (
	"strconv"
	"testing"
)

func test_LookAndSay(t *testing.T, in, expected string) {
	got := LookAndSay(in)
	if got != expected {
		t.Errorf("for %v: expected '%v' but got '%v'", in, expected, got)
	}
}

func Test_LookAndSay(t *testing.T) {
	t.Run("1", func(t *testing.T) { test_LookAndSay(t, "1", "11") })
	t.Run("2", func(t *testing.T) { test_LookAndSay(t, "11", "21") })
	t.Run("1", func(t *testing.T) { test_LookAndSay(t, "21", "1211") })
	t.Run("1", func(t *testing.T) { test_LookAndSay(t, "1211", "111221") })
	t.Run("1", func(t *testing.T) { test_LookAndSay(t, "111221", "312211") })
	t.Run("1", func(t *testing.T) { test_LookAndSay(t, "312211", "13112221") })
}

func test_Game(t *testing.T, in string, n int, expected string) {
	gotseq, len := Game(in, n)
	got := ""
	for i := 0; i < len; i++ {
		got += strconv.Itoa(gotseq[i])
	}

	if got != expected {
		t.Errorf("for %v(%c): expected '%v' but got '%v'", in, n, expected, got)
	}
}

func Test_Game(t *testing.T) {
	t.Run("1", func(t *testing.T) { test_Game(t, "1", 5, "312211") })
}
