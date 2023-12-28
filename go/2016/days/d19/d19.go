// Adventofcode 2016, d19, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 example 3
// TEST: -1 input 1815603
// TEST: example 2
// TEST: input 1410630

// part1 is done naively, computing the algorithm as given
// the circle of N elves is a slice of size N integers, the number representing
// the amount of elves to go left to find the next active one, skipping the
// ones in between. All start at 1, we stop when one elf has value N.

// part2:
// Running d19 with a number argument (.g: ./d19 12) finds the remaining elf
// the naive way, but it is too slow to find the actual input.
// By looking at the data thus generated by (bash):
//   for i in {1..5000}; do echo  "$i [$(./d19 $i)]"; done >/tmp/L
// I could find an heuristic, the result could be seen as following a pattern:
//   i,r:  4,1  5,2  6,3  7,5  8,7  9,9  10,1
// - for a circle of size i, if the result is r (elf number at the end)
// - for i=4, r is 1
// - then, each time we increment i by 1, we increment r by 1
// - until i becomes twice r, e.g. i=6, r=3
// - then, for incrementing i by 1, r is incremented by 2
// - until r becomes equal to i, e.g. i=9, r=9
// - then, we restart with r=1, e.g: i=10, r=1
// - repeat...
// This heuristic is implemented in the bash script heuristics.sh that served
// me to prototype the heuristic by comparing its generated results to d19 naive.
// it gives the exact same results as the for loop above.
// So I implemented this same heuristic for part 2 as part2heuristicSimple

// By reading the solutions, BOT-Brad on reddit found a variant:
// He noticed that r becames 1 for i being 1,2,4,10,28,82... i.e. where
// i+1 was 3*i-2. So he first found the closest i for r==1 this way,
// and applied my heuristic from there
// I implemented this in part2heuristic to find directly the value


package main

import (
	"flag"
	"fmt"
	"regexp"
)

var size int					//  the number of elves

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
	renum := regexp.MustCompile("^[[:digit:]]+$")
	if renum.MatchString(infile) { // real part2 computation, naive
		size = atoi(infile)
		fmt.Println(part2(size))
		return
	} else {					// default: use heuristic for fastest answer
		lines := fileToLines(infile)
		size = atoi(lines[0])
	}
	
	var result int
	if *partOne {
		VP("Running Part1")
		result = part1(size)
	} else {
		VP("Running Part2")
		result = part2heuristic(size)
	}
	fmt.Println(result)
}

//////////// Part 1

func part1(l int) int {
	r := make([]int, l, l)
	for i := 0; i < l; i++ {
		r[i] = 1
	}
	i := 0
	for r[i] != l {
		i = turn1(i, l, r)
	}
	return i+1
}

//////////// Part 2

// the real computation, too slow to go above 50 000
func part2(l int) int {
	r := make([]int, l, l)
	for i := 0; i < l; i++ {
		r[i] = 1
	}
	i, n := 0, l
	for r[i] != l {
		i = turn2(i, l, n, r)
		n--						// each turn, one elf is removed
	}
	return i+1
}

//////////// Common Parts code

//////////// Part1 functions

// one turn, stealing from left neighbor
// elf #i in circle r of size l take presents from neighbor at i + r[i],
// updates its r[i] to next elf, and returns its position
func turn1(i, l int, r []int) int {
	j := (i + r[i]) % l
	r[i] += r[j]				// next one after i in circle is the one that was after j
	r[j] = 0 					// j is out
	return (i + r[i]) % l
}

//////////// Part2 functions

// one turn, stealing from opposite elf
// elf #i in circle r of size l, with only n active elves take presents from neighbor at i + n/2 actives
// updates its r[i] to next elf, and returns its position
func turn2(i, l, n int, r []int) int {
	j, prev := i, i
	for step := 0; step < (n/2); step++ {
		prev = j
		j = (j + r[j]) % l
	}
	r[prev] += r[j]				// "unlink" j by making previous elf link to its next
	r[j] = 0 					// j is out
	return (i + r[i]) % l
}

// fast heuristic: going through all i (circle lengths) values sequentially
// not used, for reference only
func part2heuristicSimple(l int) int {
	i, r := 4, 1
	for {
		for r*2 < i {
			i++
			r++
			if i == l { return r;}
		}
		for r < i {
			i++
			r += 2
			if i == l { return r;}
		}
		i++
		r=1
		if i == l { return r;}
	}
}

// fastest heuristic
func part2heuristic(l int) int {
	var i int
	for i = 4; l > 3*i-2; i = 3*i-2 { // find the i with r=1 just below l
	}
	VPf("We start with i=%d and r=%d, just below l=%d\n", i, 1, l)
	// then we see if l is in the =+1 part or =+2 part above i-1 (separated at i2)
	i2 := (i-1) * 2
	if l <= i2 {
		VP("l is in the \"increment by\" lowest half above i")
		return l - (i-1)
	} else {	
		VP("l is in the \"increment by\" lowest half above i")
		return (l - i2) * 2 + (i-1)
	}
}