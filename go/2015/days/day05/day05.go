// Adventofcode 2015, day05, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 input 258
// TEST: input 53
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	partOne := flag.Bool("1", false, "run part one code, instead of part 2 (default)")
	flag.Parse()
	infile := "input.txt"
	if flag.NArg() > 0 {
		infile = flag.Arg(0)
	}
	file, err := os.Open(infile)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	input := bufio.NewScanner(file)

	var result int
	if *partOne {
		fmt.Println("Running Part1")
		result = Part1(input)
	} else {
		fmt.Println("Running Part2")
		result = Part2(input)
	}
	fmt.Println(result)
}

func Part1(input *bufio.Scanner) int {
	reThreeVowels := regexp.MustCompile(`[aeiou].*[aeiou].*[aeiou]`)
	reExcluded := regexp.MustCompile(`(ab|cd|pq|xy)`)
	matches := 0
	for input.Scan() {
		line := input.Text()
		if reThreeVowels.MatchString(line) &&
			hasDouble(line) &&
			!reExcluded.MatchString(line) {
			matches++
		}
	}
	return matches
}

func Part2(input *bufio.Scanner) int {
	matches := 0
	for input.Scan() {
		line := input.Text()
		if twoPairs(line) && spacedPair(line) {
			matches++
		}
	}
	return matches
}

// string matching functions
// we cannot regexps like `(.)\1` as GO regexp are RE2, without backreferences

// contains at least one letter that appears twice in a row, like xx
func hasDouble(s string) bool {
	char := []byte(s)
	len := len(s) - 1
	for i := 0; i < len; i++ {
		if char[i+1] == char[i] {
			return true
		}
	}
	return false
}

// returns the list (string) of unique letters in a string
func lettersOf(s string) string {
	letters := ""
	for _, r := range s {
		if !strings.ContainsRune(letters, r) {
			letters += string(r)
		}
	}
	return letters
}

// returns the list (slice) of unique pairs of letters in a string
func pairsOf(s string) []string {
	pairs := make([]string, 0)
	for i := 0; i < len(s)-1; i++ {
		pair := s[i : i+2]
		in := false
		for j := 0; j < len(pairs); j++ {
			if pair == pairs[j] {
				in = true
				break
			}
		}
		if !in {
			pairs = append(pairs, pair)
		}
	}
	return pairs
}

// contains a pair of any two letters that appears at least twice in the string
// without overlapping, like xyxy or aabcdefgaa
func twoPairs(s string) bool {
	pairs := pairsOf(s)
	for _, pair := range pairs {
		re := regexp.MustCompile(pair + ".*" + pair)
		if re.MatchString(s) {
			return true
		}
	}
	return false
}

// contains at least one letter which repeats with exactly one letter between them
func spacedPair(s string) bool {
	letters := lettersOf(s)
	for i := 0; i < len(letters); i++ {
		letter := letters[i : i+1]
		re := regexp.MustCompile(letter + "." + letter)
		if re.MatchString(s) {
			return true
		}
	}
	return false
}
