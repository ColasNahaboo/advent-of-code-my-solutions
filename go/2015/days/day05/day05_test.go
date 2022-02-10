package main

import (
	"bufio"
	"log"
	"os"
	"reflect"
	"testing"
)

func Test_hasDouble(t *testing.T) {
	t.Run("Test hasDouble.1", func(t *testing.T) {
		got := hasDouble("abbc")
		expected := true
		if got != expected {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})

	t.Run("Test hasDouble.2", func(t *testing.T) {
		got := hasDouble("abdbc")
		expected := false
		if got != expected {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
}

func Test_lettersOf(t *testing.T) {
	t.Run("Test lettersOf.1", func(t *testing.T) {
		got := lettersOf("abbc")
		expected := "abc"
		if got != expected {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
}

func Test_pairsOf(t *testing.T) {
	t.Run("Test pairsOf.1", func(t *testing.T) {
		got := pairsOf("abbc")
		expected := []string{"ab", "bb", "bc"}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
	t.Run("Test pairsOf.2", func(t *testing.T) {
		got := pairsOf("abbcabbab")
		expected := []string{"ab", "bb", "bc", "ca", "ba"}
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
}

func Test_twoPairs(t *testing.T) {
	t.Run("Test twoPairs.1", func(t *testing.T) {
		got := twoPairs("xyxy")
		expected := true
		if got != expected {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
	t.Run("Test twoPairs.2", func(t *testing.T) {
		got := twoPairs("aabcdefgaa")
		expected := true
		if got != expected {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
	t.Run("Test twoPairs.3", func(t *testing.T) {
		got := twoPairs("aaa")
		expected := false
		if got != expected {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
}

func Test_spacedPair(t *testing.T) {
	t.Run("Test spacedPair.1", func(t *testing.T) {
		got := spacedPair("xyzxy")
		expected := false
		if got != expected {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
	t.Run("Test spacedPair.2", func(t *testing.T) {
		got := spacedPair("abcdefeghi")
		expected := true
		if got != expected {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
}

func Test_Part2(t *testing.T) {
	t.Run("Test Part2.1", func(t *testing.T) {
		file, err := os.Open("part2test1.txt")
		if err != nil {
			log.Fatalf("failed opening file: %s", err)
		}
		input := bufio.NewScanner(file)

		got := Part2(input)
		expected := 1
		if got != expected {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
	t.Run("Test Part2.2", func(t *testing.T) {
		file, err := os.Open("part2test2.txt")
		if err != nil {
			log.Fatalf("failed opening file: %s", err)
		}
		input := bufio.NewScanner(file)

		got := Part2(input)
		expected := 1
		if got != expected {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
	t.Run("Test Part2.3", func(t *testing.T) {
		file, err := os.Open("part2test3.txt")
		if err != nil {
			log.Fatalf("failed opening file: %s", err)
		}
		input := bufio.NewScanner(file)

		got := Part2(input)
		expected := 0
		if got != expected {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
	t.Run("Test Part2.4", func(t *testing.T) {
		file, err := os.Open("part2test4.txt")
		if err != nil {
			log.Fatalf("failed opening file: %s", err)
		}
		input := bufio.NewScanner(file)

		got := Part2(input)
		expected := 0
		if got != expected {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}
	})
}
