package main

import (
	"testing"
)

//------------------------------------------------------------------

func Test_a2h(t *testing.T) {
	t_a2h("32T3K", t, "32T3K", 20302100313)
	t_a2h("T55J5", t, "T55J5", 41005051105)
	t_a2h("KK677", t, "KK677", 31313060707)
	t_a2h("KTJJT", t, "KTJJT", 31310111110)
	t_a2h("QQQJA", t, "QQQJA", 41212121114)
}
func t_a2h(label string, t *testing.T, arg string, expected int) {
	// setup here
	t.Run(label, func(t *testing.T) {
		got := a2h(arg)
		if got != expected {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
}

//------------------------------------------------------------------

func Test_a2hi(t *testing.T) {
	t_a2hi("32T3K", t, "32T3K", [5]int{3,2,10,3,13})
	t_a2hi("T55J5", t, "T55J5", [5]int{10,5,5, 11,5})
	t_a2hi("KK677", t, "KK677", [5]int{13,13,6,7,7})
	t_a2hi("KTJJT", t, "KTJJT", [5]int{13,10,11,11,10})
	t_a2hi("QQQJA", t, "QQQJA", [5]int{12,12,12,11,14})
}
func t_a2hi(label string, t *testing.T, arg string, expected [5]int) {
	// setup here
	t.Run(label, func(t *testing.T) {
		got := a2hi(arg)
		if got != expected {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
}

//------------------------------------------------------------------

func Test_h2a(t *testing.T) {
	t_h2a("32T3K", t, "32T3K", 20302100313)
	t_h2a("T55J5", t, "T55J5", 41005051105)
	t_h2a("KK677", t, "KK677", 31313060707)
	t_h2a("KTJJT", t, "KTJJT", 31310111110)
	t_h2a("QQQJA", t, "QQQJA", 41212121114)
}
func t_h2a(label string, t *testing.T, expected string, arg int) {
	// setup here
	t.Run(label, func(t *testing.T) {
		got := h2a(arg)
		if got != expected {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
}

//------------------------------------------------------------------

// set of unit tests
func Test_cardsOf(t *testing.T) {
	// setup code

	// table-driven tests: calls to test function(s) with various args
	// adapt args and expected as needed
	t_cardsOf("32T3K", t, "32T3K", []int{3, 2, 10, 13})
	t_cardsOf("T55J5", t, "T55J5", []int{10, 5})
	t_cardsOf("KK677", t, "KK677", []int{13, 6, 7})
	t_cardsOf("KTJJT", t, "KTJJT", []int{13, 10})
	t_cardsOf("QQQJA", t, "QQQJA", []int{12, 14})

	// cleanup code
}

// test function
func t_cardsOf(label string, t *testing.T, arg string, expected []int) {
	// setup here
	t.Run(label, func(t *testing.T) {
		hand := a2h(arg)
		got := cardsOf(hand)
		if ! sliceIntEquals(got, expected) {
			t.Errorf("for %v, expected '%v' but got '%v'", hand, expected, got)
		}
	})
}
