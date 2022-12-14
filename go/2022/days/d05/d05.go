// Adventofcode 2022, d05, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 input TPGVQPFDH
// TEST: input DMRDFRHHH
package main

import (
	"flag"
	"fmt"
	"regexp"
	"log"
)

var verbose bool

type Stack []string
type Stacks []Stack
type Move struct {
	amount int
	from int
	to int
}
type Moves []Move

var remove = regexp.MustCompile("^move ([[:digit:]]+) from ([[:digit:]]+) to ([[:digit:]]+)") 

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
	stacks, moves := parseDrawing(lines) 
	
	var result string
	if *partOne {
		VP("Running Part1")
		result = part1(stacks, moves)
	} else {
		VP("Running Part2")
		result = part2(stacks, moves)
	}
	fmt.Println(result)
}

//////////// Part 1

func part1(stacks Stacks, moves Moves) string {
	VP("Stacks:", stacks)
	VP("Moves: ", moves)
	runMovesCrateMover9000(&stacks, moves)
	return stacksTops(stacks)
}

//////////// Part 2
func part2(stacks Stacks, moves Moves) string {
	VP("Stacks:", stacks)
	VP("Moves: ", moves)
	runMovesCrateMover9001(&stacks, moves)
	return stacksTops(stacks)
}

//////////// Common Parts code

// Parse the input file. We expect a correct one, minimal checks are performed.
// The "cut", an empty line, separates the stacks of crates layout from the moves
func parseDrawing(lines []string) (stacks Stacks, moves Moves) {
	var cut, i, c int
	// find empty line
	for i = 0; i < len(lines); i++ {
		if lines[i] == "" {
			cut = i
			goto CUT
		}
	}
	log.Fatalln("Parse error: cannot find cut line")
CUT:
	// parse stacks of crates above the cut
	nstacks := (len(lines[cut-1]) + 1 ) / 4
	stacks = make(Stacks, nstacks, nstacks)
	// We know now the maximum size of a stack, so pre-allocate them.
	for c = 0; c < nstacks; c++ {
		stacks[c] = make(Stack, 0, nstacks * (cut - 1))
	}
	for i = cut - 2; i >= 0; i-- {
		if len(lines[i]) != nstacks * 4 - 1 {
			log.Fatalln("Parse error: bad line length", len(lines[i]), lines[i])
		}
		parseStacksRow(&stacks, lines[i])
	}

	// parse moves after the cut
	moves = make(Moves, 0, len(lines) - cut - 1)
	for i = cut + 1; i < len(lines); i++ {
		parseMove(&moves, lines[i])
	}
	return
}

func parseStacksRow(stacks *Stacks, line string) {
	for c := 0; c < len(*stacks); c++ {
		crate := string(line[c * 4 + 1])
		if crate != " " {
			(*stacks)[c] = append((*stacks)[c], crate)
		}
	}
}

func parseMove(moves *Moves, line string) {
	m := remove.FindStringSubmatch(line)
	if m == nil {
		log.Fatalln("Move syntax error: " + line)
	}
	// be careful: stacks names in input start at 1, but we start at 0
	*moves = append(*moves, Move{atoi(m[1]), atoi(m[2]) - 1, atoi(m[3]) - 1})
}

func stacksTops(stacks Stacks) (tops string) {
	for c := 0; c < len(stacks); c++ {
		tops = tops + stacks[c][len(stacks[c]) - 1]
	}
	return
}
	

//////////// Part1 functions

func runMovesCrateMover9000(stacks *Stacks, moves Moves) {
	for _, move := range moves {
		for i := 0; i < move.amount; i++ {
			// appends last elt of from to to
			(*stacks)[move.to] = append((*stacks)[move.to], (*stacks)[move.from][len((*stacks)[move.from]) - 1])
			// remove last element from from
			(*stacks)[move.from] = 	(*stacks)[move.from][:len((*stacks)[move.from]) - 1]
		}
	}
	return
}

//////////// Part2 functions

func runMovesCrateMover9001(stacks *Stacks, moves Moves) {
	for _, move := range moves {
		// appends "amount" last elt of from to to, bottom to top
		for i := len((*stacks)[move.from]) - move.amount; i < len((*stacks)[move.from]); i++ {
			(*stacks)[move.to] = append((*stacks)[move.to], (*stacks)[move.from][i])
		}
		// remove "amount" last elements from from
		(*stacks)[move.from] = 	(*stacks)[move.from][:len((*stacks)[move.from]) - move.amount]
	}
	return
}
