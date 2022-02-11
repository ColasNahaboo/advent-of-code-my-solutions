// Adventofcode 2015, day07, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 input 46065
// TEST: input 14134
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Board struct { // maps indexed by the wire label
	values map[string]uint16 // the electrical value of wires
	wiring map[string]Rule   // the rules to compute them
}
type Rule struct { // a closure, to solve values by reverse-chaining
	gate func(*Board, string, string) uint16 // the GateXxx() instruction
	v1   string                              // first param, name or int
	v2   string                              // optional param, name or int
}
type Parsers map[*regexp.Regexp]func(*Board, []string) // parse input instructions

func main() {
	partOne := flag.Bool("1", false, "run part one code, instead of part 2 (default)")
	flag.Parse()
	infile := "input.txt"
	if flag.NArg() > 0 {
		infile = flag.Arg(0)
	}
	lines := FileToLines(infile)

	var result uint16
	if *partOne {
		fmt.Println("Running Part1")
		result = Part1(lines)
	} else {
		fmt.Println("Running Part2")
		result = Part2(lines)
	}
	fmt.Println(result)
}

//////////// Part 1

func Part1(lines []string) uint16 {
	board := InitBoard(lines)
	return WireValue(board, "a")
}

//////////// Part 2
func Part2(lines []string) uint16 {
	board := InitBoard(lines)
	value := WireValue(board, "a") // the value of a from Part1
	board = InitBoard(lines)       // reset board
	board.values["b"] = value      // override b with ex value of a
	return WireValue(board, "a")
}

//////////// Common Parts code

func InitBoard(lines []string) *Board {
	board := new(Board)
	board.values = make(map[string]uint16)
	board.wiring = make(map[string]Rule)

	parsers := InitParsers()
	for _, line := range lines {
		ok := false
		for re, parser := range parsers {
			if s := re.FindStringSubmatch(line); s != nil {
				parser(board, s)
				ok = true
				break
			}
		}
		if !ok {
			fmt.Println("UNMATCHED:", line)
		}
	}
	return board
}

// create the Parsers, the interpreting machine as a map:
// instruction regexp => instruction parsing function
func InitParsers() Parsers {
	parsers := make(Parsers)
	AddParser(&parsers, "x", ParseSet)
	AddParser(&parsers, "x AND x", ParseAnd)
	AddParser(&parsers, "x OR x", ParseOr)
	AddParser(&parsers, "x LSHIFT x", ParseLShift)
	AddParser(&parsers, "x RSHIFT x", ParseRShift)
	AddParser(&parsers, "NOT x", ParseNot)
	return parsers
}

// declare a parser for a list of tokens matching an instruction (less the '-> x' part)
// for tokens, 'n' means a number, 'v' a variable, rest is literal
func AddParser(parsers *Parsers, instruction string, parser func(*Board, []string)) {
	pattern := `^`
	tokens := strings.Split(instruction, " ")
	for _, token := range tokens {
		switch {
		case token == `x`:
			pattern += `\s*([[:digit:][:lower:]]+)`
		default:
			pattern += `\s*` + token
		}
	}
	pattern += `\s*->\s*([[:lower:]]+)\s*$`
	(*parsers)[regexp.MustCompile(pattern)] = parser
}

// easier to use than raw Atoi
func isNum(s string) (uint16, bool) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return uint16(n), true
	} else {
		return 0, false
	}
}

// get the value of a wire, chaining back the rules recursively
func WireValue(board *Board, wire string) uint16 {
	i, err := strconv.Atoi(wire)
	if err == nil { // value is already determined
		return uint16(i)
	} else { // fire the rule to compute it
		n, okv := board.values[wire]
		if !okv {
			if wiring, okw := board.wiring[wire]; okw {
				n = wiring.gate(board, wiring.v1, wiring.v2)
				board.values[wire] = n
			} else {
				log.Fatalf("No rule for wire %v\n", wire)
			}
		}
		return n
	}
}

// the parsing and gate-ing (rule) functions for each instruction

// n -> v
func ParseSet(board *Board, s []string) {
	board.wiring[s[2]] = Rule{GateSet, s[1], ``}
}
func GateSet(board *Board, v1, v2 string) uint16 {
	n1 := WireValue(board, v1)
	return n1
}

// v AND v -> v
func ParseAnd(board *Board, s []string) {
	board.wiring[s[3]] = Rule{GateAnd, s[1], s[2]}
}
func GateAnd(board *Board, v1, v2 string) uint16 {
	n1 := WireValue(board, v1)
	n2 := WireValue(board, v2)
	return n1 & n2
}

// v OR v -> v
func ParseOr(board *Board, s []string) {
	board.wiring[s[3]] = Rule{GateOr, s[1], s[2]}
}
func GateOr(board *Board, v1, v2 string) uint16 {
	n1 := WireValue(board, v1)
	n2 := WireValue(board, v2)
	return n1 | n2
}

// v LSHIFT n -> v
func ParseLShift(board *Board, s []string) {
	board.wiring[s[3]] = Rule{GateLShift, s[1], s[2]}
}
func GateLShift(board *Board, v1, v2 string) uint16 {
	n1 := WireValue(board, v1)
	n2 := WireValue(board, v2)
	return n1 << n2
}

// v RSHIFT n -> v
func ParseRShift(board *Board, s []string) {
	board.wiring[s[3]] = Rule{GateRshift, s[1], s[2]}
}
func GateRshift(board *Board, v1, v2 string) uint16 {
	n1 := WireValue(board, v1)
	n2 := WireValue(board, v2)
	return n1 >> n2
}

// NOT v -> v
func ParseNot(board *Board, s []string) {
	board.wiring[s[2]] = Rule{GateNot, s[1], ``}
}
func GateNot(board *Board, v1, v2 string) uint16 {
	n1 := WireValue(board, v1)
	return ^n1
}

//////////// Generic code
func StringToLines(s string) (lines []string) {
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	err := scanner.Err()
	if err != nil {
		log.Fatal(err)
	}
	return
}

func FileToLines(filePath string) (lines []string) {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	// optionally, resize scanner's capacity for lines over 64K (65536)
	const maxCapacity = 1000000 // your required line length
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)
	// end optional
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	err = scanner.Err()
	if err != nil {
		log.Fatal(err)
	}

	return
}
