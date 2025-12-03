package main

import (
	"testing"
	//"reflect"
)

/////////// Setups

/////////// Tests to pass

func Test(t *testing.T) {
	//  calls to test function(s) args (adapt "args" and "expected")
	TEST("maxBankJolt, len 2")
	T(t, "987654321111111", 2, 98)
	T(t, "811111111111119", 2, 89)
	T(t, "234234234234278", 2, 78)
	T(t, "818181911112111", 2, 92)
	TEST("maxBankJolt, len 12")
	T(t, "987654321111111", 12, 987654321111)
	T(t, "811111111111119", 12, 811111111119)
	T(t, "234234234234278", 12, 434234234278)
	T(t, "818181911112111", 12, 888911112111)
}

/////////// Test function(s)

// adapt to your case the types of args and expected
func T(t *testing.T, banklabel string, banklen int, expected int) {
	// setup here
	TNUM++
	t.Run(TNAME + "#" + itoa(TNUM), func(t *testing.T) {
		got := maxBankJolt(bankParse(banklabel), banklen)
		if got != expected {
			t.Errorf("\nexp= '%v'\ngot= '%v'", expected, got)
		}})}

//////////// Automatic Labels: string-argument-to-TEST + # + test-number

var TNAME = "T"
var TNUM = 0

func TEST(s string) {
	TNAME = s
	TNUM = 0
}
