// Adventofcode 2018, d14, in go. https://adventofcode.com/2018/day/14
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 5158916779
// TEST: -1 example1 124515891
// TEST: -1 example2 9251071085
// TEST: -1 example3 5941429882
// TEST: example
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"fmt"
	// "regexp"
	// "flag"
	"slices"
)

//////////// Options parsing & exec parts

func main() {
	ExecOptions(2, NoXtraOpts, part1, part2)
}

//////////// Part 1

func part1(lines []string) (res int) {
	todo := parse(lines)
	scores := []int{3, 7}
	elf1 := 0
	elf2 := 1
	for len(scores) < todo + 10 {
		scores = Combine(scores, elf1, elf2) //  create recipe(s) and append them
		elf1 = (elf1 + scores[elf1] + 1) % len(scores)
		elf2 = (elf2 + scores[elf2] + 1) % len(scores)
		VP(scores)
	}
	return Score(scores[todo:todo+10])
}

func Combine(scores []int, elf1, elf2 int) []int {
	sum := scores[elf1] + scores[elf2]
	if sum >= 10 {
		scores = append(scores, sum / 10)
	}
	return append(scores, sum % 10)
}

func Score(scores []int) (res int) {
	for i := 0; i < len(scores); i++ {
		res = res * 10 + scores[i]
	}
	return
}

//////////// Part 2

func part2(lines []string) (res int) {
	recipes := []int{}
	for i := range lines[0] {
		recipes = append(recipes, atoi(lines[0][i:i+1]))
	}
	VPscore(recipes)
	scores := []int{3, 7}
	elf1 := 0
	elf2 := 1
	olen := 0
	for {
		scores = Combine(scores, elf1, elf2) //  create recipe(s) and append them
		// see if recipes has been added: so at end or just before
		if len(scores) > len(recipes) {
			if slices.Equal(recipes, scores[len(scores)-len(recipes):]) {
				return len(scores) - len(recipes)
			}
			if len(scores) - olen > 1 {
				if slices.Equal(recipes, scores[len(scores)-len(recipes)-1:len(scores)-1]) {
					return len(scores) - len(recipes) - 1
				}
			}
		}
		elf1 = (elf1 + scores[elf1] + 1) % len(scores)
		elf2 = (elf2 + scores[elf2] + 1) % len(scores)
		olen = len(scores)	
		VPscore(scores)
	}
}

//////////// Common Parts code

func parse(lines []string) int {
	return atoi(lines[0])
}

//////////// PrettyPrinting & Debugging functions. See also the VPx functions.

func VPscore(s []int) {
	if ! verbose { return }
	fmt.Print("[]int{")
	for i, n := range s {
		if i > 0 {
			fmt.Print(",")
		}
		fmt.Printf("%d", n)
	}
	fmt.Println("}")
}
