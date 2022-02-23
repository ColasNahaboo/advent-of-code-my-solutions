// Adventofcode 2015, day23, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 exemple 2
// TEST: -1 input 184
// TEST: input 231
package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	// "regexp"
)

type (
	Code []Instruction

	Instruction struct {
		src    string // useful for debugging
		op     Operator
		r      Register
		offset int
	}

	Operator int

	Register int
)

const (
	hlf = iota
	tpl
	inc
	jmp
	jie
	jio
	InstructionsCount
)

var code Code
var register [2]uint

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
	code = readCode(lines)

	var result int
	if *partOne {
		VP("Running Part1")
		result = part1()
	} else {
		VP("Running Part2")
		register[0] = 1
		result = part2()
	}
	fmt.Println(result)
}

//////////// Part 1

func part1() int {
	for p := 0; p >= 0 && p < len(code); p++ {
		p += execInst(p)
	}
	return int(register[1])
}

//////////// Part 2
func part2() int {
	for p := 0; p >= 0 && p < len(code); p++ {
		p += execInst(p)
	}
	return int(register[1])
}

//////////// Common Parts code

func readCode(lines []string) (code Code) {
	for _, line := range lines {
		if len(line) != 0 {
			var i Instruction
			i.src = line
			s := strings.Fields(line)
			switch s[0] {
			case "hlf":
				i.op = hlf
				i.r = regidx(s[1])
			case "tpl":
				i.op = tpl
				i.r = regidx(s[1])
			case "inc":
				i.op = inc
				i.r = regidx(s[1])
			case "jmp":
				i.op = jmp
				i.offset = atoi(s[1])
			case "jie":
				i.op = jie
				i.r = regidx(s[1])
				i.offset = atoi(s[2])
			case "jio":
				i.op = jio
				i.r = regidx(s[1])
				i.offset = atoi(s[2])
			default:
				log.Fatalf("Syntax error: %v", line)
			}
			code = append(code, i)
		}
	}
	return
}

func regidx(s string) Register {
	if s[0] == 'a' {
		return 0
	} else {
		return 1
	}
}

func execInst(p int) (offset int) {
	switch code[p].op {
	case hlf:
		register[code[p].r] /= 2
	case tpl:
		register[code[p].r] *= 3
	case inc:
		register[code[p].r]++
	case jmp:
		offset = code[p].offset - 1 // -1 negates the "p++" of the main loop
	case jie:
		if register[code[p].r]%2 == 0 {
			offset = code[p].offset - 1
		}
	case jio:
		if register[code[p].r] == 1 {
			offset = code[p].offset - 1
		}
	}
	VPf("%-10v(%v) --> %v ==> %v\n", code[p].src, p, register, p+offset+1)
	return
}

//////////// Part1 functions

//////////// Part2 functions
