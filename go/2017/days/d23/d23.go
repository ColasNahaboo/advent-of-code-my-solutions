// Adventofcode 2017, d23, in go. https://adventofcode.com/2017/day/23
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example
// TEST: example
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

// Similar to the virtual computer "tablets" of d19
// but Part 2 is quite different, requiring to understanding what the algorithm
// is doing

package main

import (
	"flag"
	"fmt"
	"regexp"
	"golang.org/x/exp/slices"
	"github.com/fxtlabs/primes"
)

var verbose bool

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	partThree := flag.Bool("3", false, "run exercise part3, naive slow mode of part2 (default: part2)")
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
	} else if *partThree {
		VP("Running Part3")
		fmt.Println(part3(lines))
	} else {
		VP("Running Part2")
		fmt.Println(part2(lines))
	}
}

//////////// Part 1

func part1(lines []string) int {
	t := parse(lines)
	t.Run()
	return t.mulcount
}

//////////// Part 2

func part2(lines []string) int {
	t := parse(lines)			// we just find the initial value of b
	bseed := t.prog[0].y
	b := bseed * 100 + 100000
	c := b + 17000
	VPf("b seed = %d, b = %d, c = %d\n", bseed, b, c)

	// This implements the code in README.md using a prime-finding function
	// from https://pkg.go.dev/github.com/fxtlabs/primes
	h := 0
	for ; b <= c; b += 17 {
		if ! primes.IsPrime(b) {
			h++
		}
	}
	return h
}


func part3(lines []string) int {
	t := parse(lines)			// we just find the initial value of b
	bseed := t.prog[0].y
	b := bseed * 100 + 100000
	c := b + 17000
	VPf("b seed = %d, b = %d, c = %d\n", bseed, b, c)

	// This implements the naive (slow) code in README.md
	h := 0
	for ; b <= c; b += 17 {
		f := 1
		for d := 2; d <= b; d++ {
			for e := 2; e <= b; e++ {
				if d * e == b {
					f = 0
				}
			}
		}
		if f == 0 {
			h++
		}
	}
	return h
}

//////////// Common Parts code
// Mostly copied from d19

// we mix register names or values in the same int type: Under the REG thresold
// it is a number, above it is the ID of the register.
// Values above are register names, 'a' (ascii 97) being REG, 'b' REG+1, etc
const REG = 9000000000000000000

type Tablet struct {
	regs []int
	prog []Instr
	p int						// current position in prog
	opnames []string			// instruction operators (op) names
	opargs []int				// number of args for each op (debug)
	opexec []OpExec				// exec function for each op
    // d23 specific
	mulcount int				// how many times is the mul instruction invoked?
}

type Instr struct {
	op int						// ID: index in .opnames, .opexec
	x, y int					// params
}

type OpExec func(t *Tablet) bool // changes t.regs and t.p


func parse(lines []string) (t *Tablet) {
	t = &Tablet{}
	t.Init()
	re := regexp.MustCompile("^([[:lower:]]{3}) ([-[:alnum:]]+)( ([-[:alnum:]]+))?")
	for _, line := range lines {
		m := re.FindStringSubmatch(line)
		t.prog = append(t.prog, Instr{parseOp(m[1], t), parseParam(m[2], t), parseParam(m[4], t)})
	}
	return
}

func parseOp(s string, t *Tablet) int {
	op := slices.Index[[]string, string](t.opnames, s)
	if op == -1 {
		panic("Unknown OP: " + s)
	}
	return op
}

func parseParam(s string, t *Tablet) int {
	if len(s) == 1 && s[0] >= 'a' && s[0] <= 'z' {
		reg := int(s[0] - 'a') 
		if reg >= len(t.regs) {
			t.regs = append(t.regs, make([]int, reg - len(t.regs) + 1)...)
		}
		return REG + reg
	} else if len(s) > 0 {
		return atoi(s)
	} else {
		return 0
	}
}

func valueOf(x int, t *Tablet) int {
	if x >= REG {
		return t.regs[x - REG]
	} else {
		return x
	}
}

func regOf(x int) int {
	if x >= REG {
		return x - REG
	} else {
		panic("Not a reg: " + itoa(x))
	}
}

func nameOf(x int) string {
	if x >= REG {
		return string('a' + (x - REG))
	} else {
		panic("Not a reg: " + itoa(x))
	}
}
		
func (t *Tablet) Init() {
	t.opnames = []string{"set", "sub", "mul", "jnz"}
	t.opargs = []int{2, 2, 2, 2}
	t.opexec = []OpExec{setExec, subExec, mulExec, jnzExec}
	t.regs = []int{}
}

// execute one step of a tablet program
func (t *Tablet) Exec() bool {
	return t.opexec[t.prog[t.p].op](t) && t.p < len(t.prog)
}


// run all the program on a tablet
func (t *Tablet) Run() {
	step := 0
	for t.Exec() {
		step++
		VPf("  [%d] @%d %s %v\n", step, t.p, t.opnames[t.prog[t.p].op], t.regs)
	}
}

// The Op implementations

func setExec(t *Tablet) bool {
	i := t.prog[t.p]
	t.regs[regOf(i.x)] = valueOf(i.y, t)
	t.p++
	return true
}

func subExec(t *Tablet) bool {
	i := t.prog[t.p]
	t.regs[regOf(i.x)] = valueOf(i.x, t) - valueOf(i.y, t)
	t.p++
	return true
}

func mulExec(t *Tablet) bool {
	i := t.prog[t.p]
	t.regs[regOf(i.x)] = valueOf(i.x, t) * valueOf(i.y, t)
	t.p++
	t.mulcount++
	return true
}

func jnzExec(t *Tablet) bool {
	i := t.prog[t.p]
	if valueOf(i.x, t) == 0 {
		t.p++
	} else {
		t.p += valueOf(i.y, t)
	}
	return true
}

//////////// PrettyPrinting & Debugging functions

