package main

import (
	"testing"
)

// unit tests based on the example
func TestInstructions(t *testing.T) {
	expected := "abcde"
	t.Run("Test instrSwapPos", func(t *testing.T) {
		got := string(instrSwapPos([]byte(expected), 4, 0))
		expected = "ebcda"
		if got != expected {
			t.Errorf("expected '%s' but got '%s'", expected, got)
		}
	})

	t.Run("Test instrSwapLetter", func(t *testing.T) {
		got := string(instrSwapLetter([]byte(expected), 'd', 'b'))
		expected = "edcba"
		if got != expected {
			t.Errorf("expected '%s' but got '%s'", expected, got)
		}
	})

	t.Run("Test instrReversePos", func(t *testing.T) {
		got := string(instrReversePos([]byte(expected), 0, 4))
		expected = "abcde"
		if got != expected {
			t.Errorf("expected '%s' but got '%s'", expected, got)
		}
	})

	t.Run("Test instrRotate", func(t *testing.T) {
		got := string(instrRotate([]byte(expected), -1))
		expected = "bcdea"
		if got != expected {
			t.Errorf("expected '%s' but got '%s'", expected, got)
		}
	})

	t.Run("Test instrMovePos", func(t *testing.T) {
		got := string(instrMovePos([]byte(expected), 1, 4))
		expected = "bdeac"
		if got != expected {
			t.Errorf("expected '%s' but got '%s'", expected, got)
		}
	})

	t.Run("Test instrMovePosBack", func(t *testing.T) {
		got := string(instrMovePos([]byte(expected), 3, 0))
		expected = "abdec"
		if got != expected {
			t.Errorf("expected '%s' but got '%s'", expected, got)
		}
	})
	
	t.Run("Test instrRotateLetter", func(t *testing.T) {
		got := string(instrRotateLetter([]byte(expected), 'b'))
		expected = "ecabd"
		if got != expected {
			t.Errorf("expected '%s' but got '%s'", expected, got)
		}
	})
	
	t.Run("Test instrRotateLetterBack", func(t *testing.T) {
		got := string(instrRotateLetter([]byte(expected), 'd'))
		expected = "decab"
		if got != expected {
			t.Errorf("expected '%s' but got '%s'", expected, got)
		}
	})
}
