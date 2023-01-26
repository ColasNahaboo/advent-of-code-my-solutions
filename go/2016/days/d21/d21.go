// Adventofcode 2016, d21, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 example decab
// TEST: -1 input cbeghdaf
// TEST: example abcde
// TEST: input bacdefgh

// modified input:
// the first line of the input are the letters to scramble for part1
// the second line is the passoerd to de-scramble for part2
// the rest is the instructions

package main

import (
	"flag"
	"fmt"
	// "regexp"
)

var verbose bool

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	flag.Parse()
	verbose = *verboseFlag
	infile := "input.txt"
	if flag.NArg() > 0 {
		infile = flag.Arg(0)
	}
	lines := fileToLines(infile)

	if *partOne {
		VP("Running Part1")
		fmt.Println(part1(lines))
		return
	} else {
		VP("Running Part2")
		fmt.Println(part2(lines))
		return
	}
}

//////////// Part 1

func part1(lines []string) string {
	return interpret(lines[0], lines[2:])
}

//////////// Part 2
func part2(lines []string) string {
	return interpretR(lines[1], lines[2:])
}

//////////// Common Parts code

//////////// Part1 functions

// swap position X with position Y
//     means that the letters at indexes X and Y (counting from 0) should be swapped.
// swap letter X with letter Y
//     means that the letters X and Y should be swapped (regardless of where they appear in the string).
// rotate left/right X steps
//     means that the whole string should be rotated; for example, one right rotation would turn abcd into dabc.
// rotate based on position of letter X
//     means that the whole string should be rotated to the right based on the index of letter X (counting from 0) as determined before this instruction does any rotations. Once the index is determined, rotate the string to the right one time, plus a number of times equal to that index, plus one additional time if the index was at least 4.
// reverse positions X through Y
//     means that the span of letters at indexes X through Y (including the letters at X and Y) should be reversed in order.
// move position X to position Y
//     means that the letter which is at index X should be removed from the string, then inserted such that it ends up at index Y.


func interpret(seed string, lines []string) string{
	var n1, n2 int
	var s1, s2 string
	s := []byte(seed)
	for lineno := 0; lineno < len(lines); lineno++ {
		line := lines[lineno]
		if nf, _ := fmt.Sscanf(line, "swap position %d with position %d", &n1, &n2); nf == 2 {
			s = instrSwapPos(s, n1, n2)
		} else if nf, _ := fmt.Sscanf(line, "swap letter %1s with letter %1s", &s1, &s2); nf == 2 {
			s = instrSwapLetter(s, s1[0], s2[0])
		} else if nf, _ := fmt.Sscanf(line, "rotate left %d steps", &n1); nf == 1 {
			s = instrRotate(s, -n1)
		} else if nf, _ := fmt.Sscanf(line, "rotate right %d steps", &n1); nf == 1 {
			s = instrRotate(s, n1)
		} else if nf, _ := fmt.Sscanf(line, "rotate based on position of letter %1s", &s1); nf == 1 {
			s = instrRotateLetter(s, s1[0])
		} else if nf, _ := fmt.Sscanf(line, "reverse positions %d through %d", &n1, &n2); nf == 2 {
			s = instrReversePos(s, n1, n2)
		} else if nf, _ := fmt.Sscanf(line, "move position %d to position %d", &n1, &n2); nf == 2 {
			s = instrMovePos(s, n1, n2)
		} else {
			panic(fmt.Sprintf("Syntax error line %d: %s\n", lineno, line))
		}
	}
	return string(s)
}

// the implementation of instructions

func instrSwapPos(s []byte, n1, n2 int) []byte {
	s[n1], s[n2] = s[n2], s[n1]
	return s
}
	
func instrSwapLetter(s []byte, b1, b2 byte) []byte {
	n1 := IndexOf[byte](s, b1)
	n2 := IndexOf[byte](s, b2)
	return instrSwapPos(s, n1, n2)
}

func instrRotate(s []byte, n int) []byte {
	r := make([]byte, len(s), len(s))
	for i := 0; i < len(s); i++ {
		r[i] = s[(i-n + 2*len(s)) % len(s)]
	}
	return r
}

func instrRotateLetter(s []byte, b byte) []byte {
	n := IndexOf[byte](s, b) + 1
	if n >= 5 { n++;}
	return instrRotate(s, n)
}

func instrReversePos(s []byte, n1, n2 int) []byte {
	r := make([]byte, len(s), len(s))
	for i := 0; i < n1; i++ {
		r[i] = s[i]
	}
	for i := n2 + 1; i < len(s); i++ {
		r[i] = s[i]
	}
	for i := n1; i <= n2; i++ {
		r[i] = s[n2 + n1 - i]
	}
	return r
}

func instrMovePos(s []byte, n1, n2 int) []byte {
	if n2 < n1 {
		return instrMovePosBack(s, n1, n2)
	}
	r := make([]byte, len(s), len(s))
	for i := 0; i < n1; i++ {
		r[i] = s[i]
	}
	for i := n2 + 1; i < len(s); i++ {
		r[i] = s[i]
	}
	for i := n1; i < n2; i++ {
		r[i] = s[i+1]
	}
	r[n2] = s[n1]
	return r
}

func instrMovePosBack(s []byte, n1, n2 int) []byte {
	r := make([]byte, len(s), len(s))
	for i := 0; i < n2; i++ {
		r[i] = s[i]
	}
	for i := n1 + 1; i < len(s); i++ {
		r[i] = s[i]
	}
	for i := n2 + 1; i <= n1; i++ {
		r[i] = s[i-1]
	}
	r[n2] = s[n1]
	return r
}

//////////// Part2 functions

// we reverse the functions of part1, and apply them in reverse

func interpretR(seed string, lines []string) string{
	var n1, n2 int
	var s1, s2 string
	s := []byte(seed)
	for lineno := len(lines) -1 ; lineno >= 0; lineno-- {
		line := lines[lineno]
		if nf, _ := fmt.Sscanf(line, "swap position %d with position %d", &n1, &n2); nf == 2 {
			s = instrSwapPosR(s, n1, n2)
		} else if nf, _ := fmt.Sscanf(line, "swap letter %1s with letter %1s", &s1, &s2); nf == 2 {
			s = instrSwapLetterR(s, s1[0], s2[0])
		} else if nf, _ := fmt.Sscanf(line, "rotate left %d steps", &n1); nf == 1 {
			s = instrRotateR(s, -n1)
		} else if nf, _ := fmt.Sscanf(line, "rotate right %d steps", &n1); nf == 1 {
			s = instrRotateR(s, n1)
		} else if nf, _ := fmt.Sscanf(line, "rotate based on position of letter %1s", &s1); nf == 1 {
			s = instrRotateLetterR(s, s1[0])
		} else if nf, _ := fmt.Sscanf(line, "reverse positions %d through %d", &n1, &n2); nf == 2 {
			s = instrReversePosR(s, n1, n2)
		} else if nf, _ := fmt.Sscanf(line, "move position %d to position %d", &n1, &n2); nf == 2 {
			s = instrMovePosR(s, n1, n2)
		} else {
			panic(fmt.Sprintf("Syntax error line %d: %s\n", lineno, line))
		}
	}
	return string(s)
}

// the implementation of instructions

func instrSwapPosR(s []byte, n1, n2 int) []byte {
	return instrSwapPos(s, n2, n1)
}
	
func instrSwapLetterR(s []byte, b1, b2 byte) []byte {
	return instrSwapLetter(s, b2, b1)
}

func instrRotateR(s []byte, n int) []byte {
	return instrRotate(s, -n)
}

// We just brute force it: try all possible inputs that could give s
// among the len(s) rotations of s
func instrRotateLetterR(s []byte, b byte) []byte {
	for i := 1; i <= len(s); i++ {
		r := instrRotate(s, -i)
		rr := instrRotateLetter(r, b)
		VPf("instrRotateLetterR, to %s(%s) testing %s\n", string(s), string(b), string(r))
		if sliceEquals[byte](rr, s) {
			VPf("instrRotateLetterR ==> %s\n", string(r))
			return r
		}
	}
	panic(fmt.Sprintf("Cannot instrRotateLetterR \"%s\" based on '%s'\n", string(s), string(b)))
}

func instrReversePosR(s []byte, n1, n2 int) []byte {
	return instrReversePos(s, n1, n2)
}

func instrMovePosR(s []byte, n1, n2 int) []byte {
	return instrMovePos(s, n2, n1)
}
