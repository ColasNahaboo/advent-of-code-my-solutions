// Adventofcode 2024, d22, in go. https://adventofcode.com/2024/day/22
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 37327623
// TEST: example2 23  // sequence [-2 1 -1 3]
// TEST: example3 9   // sequence [2 -2 -4 7]
// TEST: example4 27  // sequence [3 1 4 1]
// TEST: example5 27  // sequence [-1 0 -1 8]

// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"fmt"
	// "regexp"
	// "flag"
	// "slices"
)

//////////// Options parsing & exec parts

func main() {
	ExecOptions(2, NoXtraOpts, part1, part2, dotests)
}

//////////// Part 1

func part1(lines []string) (res int) {
	secrets := parse(lines)
	for _, secret := range secrets {
		res += Generate(secret, 2000)
	}
	return res
}

func Generate(s int, n int) int {
	for _ = range n {
		s = Evolve(s)
	}
	return s
}

//////////// Part 2

func part2(lines []string) (res int) {
	// inits
	initSecrets := parse(lines)
	nbuyers := len(lines)
	prices := make([][]int, nbuyers, nbuyers)
	changes := make([][]int, nbuyers, nbuyers)
	sequences := make(map[int]int) // total of prices across all buyers
	
	var ok bool
	var price int

	for b, s := range initSecrets {
		// create the buyers sequences of prices and changes
		prices[b] = make([]int, 2001, 2001)
		changes[b] = make([]int, 2001, 2001)
		prevprice := 0
		for i := range 2001 {
			price = s % 10		// price is last digit
			prices[b][i] = price
			changes[b][i] = price - prevprice
			prevprice = price
			s = Evolve(s)
		}
		
		// record all the sequences and resulting prices for each buyer

		bseqs := make(map[int]bool) // seq already seen for this buyer?
		for i := 1; i < 2001 - 3; i++ {
			// to use sequences of changes as Map keys, we convert into digits
			// from 1 to 19 )by adding 10) in base 32 (each digit takes 5 bits)
			// so that we index by integers
			seq := 0
			for j := i; j < i+4; j++ {
				seq = 32 * seq + changes[b][j]+10
			}
			
			if ! bseqs[seq] { // only store first occurence
				bseqs[seq] = true
				price := prices[b][i+3]
				//VPf("== SEQ %v @%d => %d\n", seq, i, price)
				if _, ok = sequences[seq]; ok {
					sequences[seq] += price
				} else {
					sequences[seq] = price
				}
			}
		}
	}

	// for each sequence, calculate the total price across all buyers
	maxprice := 0
	for _, price := range sequences {
		if price > maxprice {
			maxprice = price
		}
	}
	res = maxprice
	return res
}

//////////// Common Parts code

// 32 = 2^5, 64 = 2^6, 2048 = 2^11, 16777216 = 2^24
func Evolve(i int) int {
	i = ((i * 64) ^ i) % 16777216
	i = ((i / 32) ^ i) % 16777216
	i = ((i * 2048) ^ i) % 16777216
	return i
}	

func parse(lines []string) (secrets []int) {
	for _, line := range lines {
		secrets = append(secrets, int(atoi(line)))
	}
	return
}

//////////// PrettyPrinting & Debugging functions. See also the VPx functions.

func dotests(lines []string) int {
	return SubSliceIntIndex([]int{0,1,2,5,6,3,4,5,6,7,8,5,6,7,9}, []int{5, 6, 7}, 0)
}

// test a specific change sequence
// e.g. via: res = TestSeqExample([]int{-2, 1, -1, 3}, changes, prices)

func TestSeqExample(cs []int, bchanges, bprices [][]int) (res int) {
	for b, changes := range bchanges {
		//fmt.Printf("#%d Changes: %v\n", b+1, changes)
		//fmt.Printf("#%d  Prices:  %v\n", b+1, bprices[b])
		i := SubSliceIntIndex(changes, cs, 0)
		if i> 0 { fmt.Printf("  Subslice @%d: %v, prices %d\n", i, changes[i:i+5], bprices[b][i:i+5]) }
		if i >0 {
			price := bprices[b][i+len(cs)-1] // price is at i+3
			fmt.Printf("  Buyer #%d: price @%d = %d\n", b+1, i, price)
			res += price
		} else {
			fmt.Printf("  Buyer #%d: price does not occur\n", b+1)
		}
	}
	return
}

func DEBUG() {
	if ! debug { return }
	fmt.Println("DEBUG!")
	// ad hoc debug code here
}
