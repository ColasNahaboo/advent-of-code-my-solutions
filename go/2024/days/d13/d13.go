// Adventofcode 2024, d13, in go. https://adventofcode.com/2024/day/13
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 480
// TEST: example 875318608908
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

// We express the buttons values as 2 equations
// Input:                Matrix: X    Y      names
// Button A: X+94, Y+34     A   94   34      x0 y0
// Button B: X+22, Y+67     B   22   67      x1 y1
// Prize: X=8400, Y=5400      8400 5400      x2 y2
//
// we have thus to solve 2 equations: a*x0 + b*x1 = x2 and a*y0 + b*y1 = y2
//
// We represent by the type CM that is an array of the 2 equations for X and Y
// e.g. for the example above:  {{94 22 8400} {34 67 5400}}
//      cm[0] = {94 22 8400}, cm[1] = {34 67 5400}
// The determinant (d) is thus  = x0*y1 - x1*y0     (never 0 for our inputs)
// And the solution: a = (x2*y1 - x1*y2) / d
//                   b = (x0*y2 - x2*y0) / d
// We use the big.Rat Fractions package to check that a and b are integers

package main

import (
	"flag"
	"fmt"
	"regexp"
	"math/big"					// used for integer fractions
	// "slices"
)

var verbose, debug bool

type CM [2][3]int // Claw Machine: Equations for X and Y: a*x[0] + b*x[1] = x[2]

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	debugFlag := flag.Bool("V", false, "debug: even more verbose")
	flag.Parse()
	verbose = *verboseFlag
	debug = *debugFlag
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

func part1(lines []string) (tokens int64) {
	cms := parse(lines)
	for _, cm := range cms {
		tokens += cm.Tokens()
	}
	return 
}

func (cm CM) Tokens() int64 {
	x, y := cm[0], cm[1]		   // for readability below
	d := x[0] * y[1] - x[1] * y[0] // determinant
	an := x[2] * y[1] - x[1] * y[2]
	bn := x[0] * y[2] - x[2] * y[0]
	a := Fraction(an, d)
	b := Fraction(bn, d)
	if a.IsInt() && b.IsInt() {
		return a.Num().Int64() * 3 + b.Num().Int64()
	}
	return 0
}

//////////// Part 2

func part2(lines []string) (tokens int64) {
	cms := parse(lines)
	for _, cm := range cms {
		cm.Fix(10000000000000)	//  we fix the values of the prizes to reach
		tokens += cm.Tokens()
	}
	return 
}

func (cm *CM) Fix(n int) {
	cm[0][2] += n
	cm[1][2] += n
}

//////////// Common Parts code

func parse(lines []string) (cms []CM) {
	renum := regexp.MustCompile("[[:digit:]]+")
	for i := 0; i < len(lines); i += 2 { // process one CM per blocks of 3 lines
		cm := CM{}
		al := atoil(renum.FindAllString(lines[i], -1)) // Button A
		cm[0][0] = al[0]
		cm[1][0] = al[1]
		i++
		al = atoil(renum.FindAllString(lines[i], -1)) // Button B
		cm[0][1] = al[0]
		cm[1][1] = al[1]
		i++
		al = atoil(renum.FindAllString(lines[i], -1)) // Prizes
		cm[0][2] = al[0]
		cm[1][2] = al[1]
		cms = append(cms, cm)
	}
	return
}

func Fraction(x, y int) (*big.Rat) {
	r := big.Rat{}
	return (&r).Quo(big.NewRat(int64(x), 1), big.NewRat(int64(y), 1))
}

//////////// PrettyPrinting & Debugging functions

func DEBUG() {
	if ! debug { return }
	// ad hoc debug code here
}
