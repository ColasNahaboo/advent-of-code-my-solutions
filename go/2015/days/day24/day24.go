// Adventofcode 2015, day24, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 input 10439961859
// TEST: input 72050269
package main

import (
	"flag"
	"fmt"
	"log"
	"math/big" // just for fun, for QE results
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
	nums := make([]int, len(lines))
	sum := 0
	prev := 0
	for i := 0; i < len(lines); i++ {
		nums[i] = atoi(lines[i])
		if nums[i] < prev {
			log.Fatalf("input is not ordered: %v after %v\n", nums[i], prev)
		}
		sum += nums[i]
		prev = nums[i]
	}
	if len(nums) >= 64 {
		log.Fatalf("Code works only up to 63 numbers, not: %v\n", len(nums))
	}
	if sum%3 != 0 {
		log.Fatalf("Sum of packages is not a multiple of 3: %v\n", sum)
	}
	VPf("Total size: %v, Compartment size: %v\n", sum, sum/3)

	var result big.Int
	if *partOne {
		VPf("Running Part1: %v nums, sets of %v weights: %v .. %v\n", len(nums), sum/3)
		result = part1(nums, sum/3)
	} else {
		VPf("Running Part2: %v nums, sets of %v weights: %v .. %v\n", len(nums), sum/4)
		result = part1(nums, sum/4)
	}
	fmt.Println(result.String())
}

//////////// Part 1

func part1(nums []int, maxweight int) big.Int {
	var minset, maxset int // min and max size of subsets
	weight := 0
	for i := 0; i < len(nums); i++ {
		weight += nums[i]
		if weight >= maxweight {
			maxset = i
			break
		}
	}
	weight = 0
	for i := len(nums) - 1; i >= 0; i-- {
		weight += nums[i]
		if weight >= maxweight {
			minset = len(nums) - i
			break
		}
	}

	for setl := minset; setl <= maxset; setl++ {
		set := []int{}
		for _, setn := range subsetsList(len(nums), setl) {
			w := subsetWeight(nums, setn)
			if w == maxweight {
				set = append(set, setn)
			}
		}
		VPf("%v sets of length %v possible\n", len(set), setl)
		if len(set) > 0 {
			if len(set) == 1 {
				return subsetQE(nums, set[0])
			}
			var minqe, minqe0 big.Int
			minqe0 = minqe
			for _, setn := range set {
				qe := subsetQE(nums, setn)
				if qe.Cmp(&minqe) == -1 || minqe.Cmp(&minqe0) == 0 {
					minqe = qe
				}
			}
			fmt.Println("Found:", minqe.String())
			return minqe
		}
	}
	return *big.NewInt(int64(0))
}

//////////// Part 2
func part2(lines []string) int {
	for _, line := range lines {
		fmt.Println(line)
	}
	return 0
}

//////////// Common Parts code

// sums the weights of a subset bitset
func subsetWeight(nums []int, setn int) (sum int) {
	for i := 0; i < 64; i++ {
		if (1<<i)&setn != 0 {
			sum += nums[i]
		}
	}
	return
}

// product of the weights of a subset bitset: the QE Quantum Entanglement
// we use big numbers as they can grow extra large
func subsetQE(nums []int, setn int) (qe big.Int) {
	qe = *big.NewInt(int64(1))
	for i := 0; i < 64; i++ {
		if big.NewInt(int64(setn)).Bit(i) != 0 {
			qe.Mul(&qe, big.NewInt(int64(nums[i])))
		}
	}
	return
}

// Enumerate all the subsets of size k from a set of size n
// Fast algortihm via bitsets, but only for n < 64
// see http://math0.wvstateu.edu/~baker/cs405/code/Subset.html

func subsetsList(n, k int) []int {
	x := intPower(2, k) - 1
	max := intPower(2, n) - intPower(2, n-k)
	s := make([]int, 1, max-x)
	s[0] = x
	for x < max {
		u := x & (-x)
		v := x + u
		x = v + (((v ^ x) / u) >> 2)
		s = append(s, x)
	}
	return s
}

//////////// Part1 functions

//////////// Part2 functions
