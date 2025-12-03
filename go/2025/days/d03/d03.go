// Adventofcode 2025, d03, in go. https://adventofcode.com/2025/day/03
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 357
// TEST: example 3121910778619
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"fmt"
)

// a battery bank is a list of batteries (indexes) joltages (values)
type Bank []int					

//////////// Options parsing & exec parts

func main() {
	ExecOptions(2, NoXtraOpts, part1, part2)
}

//////////// Part 1

func part1(lines []string) (res int) {
	banks := parse(lines)
	for _, bank := range banks {
		res += maxBankJolt(bank, 2)
	}
	return 
}

//////////// Part 2

func part2(lines []string) (res int) {
	banks := parse(lines)
	for _, bank := range banks {
		res += maxBankJolt(bank, 12)
	}
	return 
}

//////////// Common Parts code

// find max joltage for batslen batteries to turn on
func maxBankJolt(bank Bank, batslen int) (joltage int) {
	b := 0						// start looking at this position
	bl := batslen - 1			// we must leave room for bl other batteries
	var j int
	for _ = range batslen {
		b, j = maxJolt(bank, b, bl)
		joltage = joltage * 10 + j
		b++
		bl--
	}
	VPf("joltage = %d\n", joltage)
	return
}

// find max jolt value in bank in all positions i and more, with room for l
func maxJolt(bank Bank, b, l int) (maxbat, maxjolt int) {
	for i := b; i < len(bank) - l; i ++ {
		jolt := bank[i]
		if jolt > maxjolt {
			maxbat = i
			maxjolt = jolt
		}
	}
	VPf("  maxJolt %d(%d) ==> %d[%d]\n", b, l, maxbat, maxjolt)
	return
}
		

func parse(lines []string) (banks []Bank) {
	for _, line := range lines {
		banks = append(banks,  bankParse(line))
	}
	return
}

func bankParse(line string) (bank Bank) {
	for _, bat := range line {
		bank = append(bank, int(bat - '0'))
	}
	return
}

//////////// PrettyPrinting & Debugging functions. See also the VPx functions.

func DEBUG() {
	if ! debug { return }
	fmt.Println("DEBUG!")
	// ad hoc debug code here
}
