// Adventofcode 2017, d16, in go. https://adventofcode.com/2017/day/16
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example baedc
// TEST: example abcde
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

// in the input files, we added a second line with the number of programs
// E.g example.txt is the 2 lines:
// s1,x3/4,pe/b
// 5

// For part2, we just first determine the number of dances after which we get
// back to the start position, and only count the dances in a loop.

package main

import (
	"flag"
	"fmt"
)

var verbose bool

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	flag.Parse()
	verbose = *verboseFlag
	var infile string
	if flag.NArg() > 0 {
		infile = flag.Arg(0)
	} else {
		infile = fileMatch("input,[[:alnum:]]*,[[:alnum:]]*.test")
	}
	lines := fileToLines(infile)
	if *partOne {
		VP("Running Part1")
		fmt.Println(part1(lines))
	} else {
		VP("Running Part2")
		fmt.Println(part2(lines))
	}
}

//////////// Part 1

func part1(lines []string) string {
	progs, moves := parse(lines)
	VPf("Progs: %s\nMoves: %v\n", progNames(progs), moves)
	for _, move := range moves {
		progs = dance(progs, move)
		VPf("  Move %v ==> %s\n", move, progNames(progs))
	}
	return progNames(progs)
}

//////////// Part 2

const DANCES = 1000000000

func part2(lines []string) string {
	progs, moves := parse(lines)
	var loop int
	for i := 0; i < 10000; i++ {
		for _, move := range moves {
			progs = dance(progs, move)
		}
		if isStart(progs) {
			loop = i+1
			break
		}
	}
	steps := DANCES % loop
	fmt.Printf("Dances loop size: %d, %d dances is the same as %d\n", loop, DANCES, steps)
	for i := range progs {		// reset progs
		progs[i] = i
	}
	for i := 0; i < steps; i++ {
		for _, move := range moves {
			progs = dance(progs, move)
		}
	}
	return progNames(progs)
}

func isStart(progs []int) bool {
	for i, p := range progs {
		if p != i {
			return false
		}
	}
	return true
}

//////////// Common Parts code

type Move struct {
	op int						// SPIN, EXCHANGE, PARTNER
	x, y int					// params
}
const (
	SPIN = 0
	EXCHANGE = 1
	PARTNER = 2
)
	
	
func parse(lines []string) (progs []int, moves []Move) {
	nprogs := atoi(lines[1])
	progs = make([]int, nprogs, nprogs)
	for i := range progs {
		progs[i] = i
	}
	i, s, endprog := 0, lines[0], byte('a' + nprogs)
	var n1, n2 int
	for i < len(s) {
		VPf("  [%d] (%s)\n", i, s[i:])
		switch s[i] {
		case ',':
			i++
		case 's':
			n1, i = parseInt(s, i+1)
			moves = append(moves, Move{op: SPIN, x: n1})
		case 'x':
			n1, i = parseInt(s, i+1)
			if s[i] != '/' {panic("x: Expected / got " + s[i:i+1])}
			n2, i = parseInt(s, i+1)
			moves = append(moves, Move{op: EXCHANGE, x: n1, y: n2})
		case 'p':
			n1, i = parseProg(s, i+1, endprog)
			if s[i] != '/' {panic("p: Expected / got " + s[i:i+1])}
			n2, i = parseProg(s, i+1, endprog)
			moves = append(moves, Move{op: PARTNER, x: n1, y: n2})
		}
		VPprogs(progs)
	}
	return		
}

func parseInt(s string, i int) (int, int) {
	val := 0
	for s[i] >= '0' && s[i] <= '9' {
		val = val * 10 + int(s[i] - '0')
		i++
	}
	return val, i
}

func parseProg(s string, i int, endprog byte) (int, int) {
	if s[i] < 'a' && s[i] >= endprog {
		panic("Not a progname: " + s[i:i+1])
	}
	return int(s[i] - 'a'), i+1
}

func progNames(progs []int) (s string) {
	for _, p := range progs {
		s = s + string('a' + p)
	}
	return
}

func dance(start []int, move Move) (progs []int) {
	l := len(start)
	switch move.op {
	case SPIN:
		progs = make([]int, l, l)
		copy(progs, start[l - move.x:])
		copy(progs[move.x:], start[0:l - move.x])
	case EXCHANGE:
		progs = start
		progs[move.x], progs[move.y] = progs[move.y], progs[move.x]
	case PARTNER:
		progs = start
		px, py := posOf(progs, move.x), posOf(progs, move.y)
		progs[px], progs[py] = progs[py], progs[px]
	}
	return
}

func posOf(list []int, number int) (int) {
	for i, v := range list {
		if number == v {
           return i
		}
	}
	panic(fmt.Sprintf("%d not found in %v", number, list))
}

//////////// PrettyPrinting & Debugging functions

func VPprogs(progs []int) {
	if ! verbose { return }
	seen := make([]bool, len(progs), len(progs))
	for _, p := range progs {
		if seen[p] {
			panic(fmt.Sprintf("%d duplicated in %v", p, progs))
		}
		seen[p] = true
	}
}
