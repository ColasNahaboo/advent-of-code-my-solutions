// Adventofcode 2017, d08, in go. https://adventofcode.com/2017/day/08
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 1
// TEST: example 10
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"flag"
	"fmt"
	"regexp"
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

	var result int
	if *partOne {
		VP("Running Part1")
		result = part1(lines)
	} else {
		VP("Running Part2")
		result = part2(lines)
	}
	fmt.Println(result)
}

//////////// Part 1

func part1(lines []string) int {
	c := parse(lines)
	c.Run()
	return c.LargestRegValue()
}

//////////// Part 2
func part2(lines []string) int {
	c := parse(lines)
	c.Run()
	return c.maxregval
}

//////////// Common Parts code

type CPU struct {
	regs []int					// register values
	regnames []string			// their names
	ptr int						// index in prog of instruction to execute next
	prog []Instr				// the program
	opnames []string			// the list of operators
	maxregval int				// for part2
}

type Instr struct {
	reg int						// register ID, index in CPU.regs
	add int						// value to add
	ifreg int					// register ID to test condition on
	ifop int					// ID of condition operator == != < <= > >=
	ifval int					// condition value
}

func parse(lines []string) (c CPU) {
	c.regs = []int{}
	c.regnames = []string{}
	c.prog = []Instr{}
	c.opnames = []string{"==", "!=", "<", "<=", ">", ">="}
	re := regexp.MustCompile("([[:alpha:]]+) (inc|dec) ([-[:digit:]]+) if ([[:alpha:]]+) ([=!<>]+) ([-[:digit:]]+)")
	for _, line := range lines {
		m := re.FindStringSubmatch(line)
		if m == nil {
			panic("Syntax Error: " + line)
		}
		regname := m[1]
		reg := c.RegId(regname)
		add := atoi(m[3])
		if m[2] == "dec" {
			add = -add
		}
		ifregname := m[4]
		ifreg :=  c.RegId(ifregname)
		ifop := IndexOf[string](c.opnames, m[5])
		if ifop == -1 {
			panic("Unknown op: \"" + m[5] + "\" at: " + line)
		}
		ifval := atoi(m[6])
		instr := Instr{reg, add, ifreg, ifop, ifval}
		c.prog = append(c.prog, instr)
	}
	return
}

func (c *CPU) RegId(name string) (r int) {
	r = IndexOf[string](c.regnames, name)
	if r == -1 {				// allocate new reg
		r = len(c.regnames)
		c.regnames = append(c.regnames, name)
		c.regs = append(c.regs, 0)
	}
	return
}

func (c *CPU) Run() {
	for p := range c.prog {
		c.Exec(p)
	}
}

func (c *CPU) Exec(p int) {
	if ! c.OK(p) {
		return
	}
	c.regs[c.prog[p].reg] += c.prog[p].add
	if c.regs[c.prog[p].reg] > c.maxregval {
		c.maxregval = c.regs[c.prog[p].reg]
	}
}

func (c *CPU) OK(p int) bool {
	i := c.prog[p]
	switch i.ifop {				// "==", "!=", "<", "<=", ">", ">="
	case 0: return c.regs[i.ifreg] == i.ifval
	case 1: return c.regs[i.ifreg] != i.ifval
	case 2: return c.regs[i.ifreg] < i.ifval
	case 3: return c.regs[i.ifreg] <= i.ifval
	case 4: return c.regs[i.ifreg] > i.ifval
	case 5: return c.regs[i.ifreg] >= i.ifval
	default: panic("Bad operator")
	}
}


func (c *CPU) LargestRegValue() (v int) {
	for _, rv := range c.regs {
		if rv > v {
			v = rv
		}
	}
	return
}
	


//////////// PrettyPrinting & Debugging functions
