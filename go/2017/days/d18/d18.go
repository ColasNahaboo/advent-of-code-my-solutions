// Adventofcode 2017, d18, in go. https://adventofcode.com/2017/day/18
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 4
// TEST: example
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"flag"
	"fmt"
	"regexp"
	"os"
	"golang.org/x/exp/slices"
	// "sync/atomic"
)

// we mix register names or values in the same int type: Under the REG thresold
// it is a number, above it is the ID of the register.
// Values above are register names, 'a' (ascii 97) being REG, 'b' REG+1, etc
const REG = 9000000000000000000

var verbose bool

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	partThree := flag.Bool("3", false, "run exercise part3, (default: part2)")
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
	return 0
}

//////////// Part 2
func part2(lines []string) int {
	// create the 2 connected tablets
	t1 := parse2(0, lines)
	t2 := parse2(1, lines)
	c12 := LinkedList[int]{}		// t1 -> t2
	t1.out = &c12
	t2.in = &c12
	c21 := LinkedList[int]{}		// t2 -> t1
	t1.in = &c21
	t2.out = &c21
	// run alternatively till the 2 have stopped
	for {
		s1 := t1.Exec()
		s2 := t2.Exec()
		if ! s1 && ! s2 {
			return t2.nsends
		}
    }
}

func parse2(id int, lines []string) (t *Tablet) {
	t = parse(lines)
	t.id = id
	p := regOf(parseParam("p", t))
	t.regs[p] = t.id
	t.opexec[parseOp("snd", t)] = snd2Exec
	t.opexec[parseOp("rcv", t)] = rcv2Exec
	return
}

//////////// Common Parts code

type Tablet struct {
	regs []int
	prog []Instr
	p int						// current position in prog
	sound int 					// freq of last sound played
	opnames []string			// instruction operators (op) names
	opargs []int				// number of args for each op (debug)
	opexec []OpExec				// exec function for each op
	// part2  fields
	id int						// ID of the tablet (program)
	out *LinkedList[int]			// channel to other tablet
	in *LinkedList[int]			// channel from other tablet
	nsends int					// how many times did we send something on out?
	waiting bool				// is waiting? (used only in verbose mode)
	// part3 fields
	cout chan int				// channel to other tablet
	cin chan int				// channel from other tablet
	cobs chan int				// wait status sent to observer
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
	t.opnames = []string{"snd", "set", "add", "mul", "mod", "rcv", "jgz"}
	t.opargs = []int{1, 2, 2, 2, 2, 1, 2}
	t.opexec = []OpExec{sndExec, setExec, addExec, mulExec, modExec, rcvExec, jgzExec}
	t.regs = []int{}
}

// execute one step of a tablet program, return false if program stopped
func (t *Tablet) Exec() bool {
	if t.waiting && t.in.IsEmpty() { // in verbose mode, dont duplicate traces
		return false
	}
	VPf("%sTablet [%d]:%d %s\n", tindent[t.id], t.id, t.p, t.PCurrentInstr())
	return t.opexec[t.prog[t.p].op](t) && t.p < len(t.prog)
}

// run all the program on a tablet
func (t *Tablet) Run() {
	for t.Exec() {
	}
}

// The Op implementations

func sndExec(t *Tablet) bool {
	t.sound = valueOf(t.prog[t.p].x, t)
	t.p++
	return true
}

func setExec(t *Tablet) bool {
	i := t.prog[t.p]
	t.regs[regOf(i.x)] = valueOf(i.y, t)
	t.p++
	return true
}

func addExec(t *Tablet) bool {
	i := t.prog[t.p]
	t.regs[regOf(i.x)] = valueOf(i.x, t) + valueOf(i.y, t)
	t.p++
	return true
}

func mulExec(t *Tablet) bool {
	i := t.prog[t.p]
	t.regs[regOf(i.x)] = valueOf(i.x, t) * valueOf(i.y, t)
	t.p++
	return true
}

func modExec(t *Tablet) bool {
	i := t.prog[t.p]
	t.regs[regOf(i.x)] = valueOf(i.x, t) % valueOf(i.y, t)
	t.p++
	return true
}

func rcvExec(t *Tablet) bool {
	if t.sound != 0 {
		fmt.Println(t.sound)
		os.Exit(0)
	}
	t.p++
	return true
}

func jgzExec(t *Tablet) bool {
	i := t.prog[t.p]
	if valueOf(i.x, t) <= 0 {
		t.p++
	} else {
		t.p += valueOf(i.y, t)
	}
	return true
}

// part2 codes

func snd2Exec(t *Tablet) bool {
	value := valueOf(t.prog[t.p].x, t)
	t.out.Put(value)
	t.nsends++
	VPf("%sTablet [%d] SEND(#%d) %d\n", tindent[t.id], t.id, t.nsends, value)
	t.p++
	return true
}

func rcv2Exec(t *Tablet) bool {
	if t.in.IsEmpty() {
		VPf("%sTablet [%d] WAIT\n", tindent[t.id], t.id)
		t.waiting = true
		return false			// this tablet waits for input
	}
	t.waiting = false
	value := t.in.Pop()
	i := t.prog[t.p]
	t.regs[regOf(i.x)] = value
	VPf("%sTablet [%d] READ %d -> %s\n", tindent[t.id], t.id, value, string(byte(regOf(i.x)) + 'a'))
	t.p++
	return true
}

//////////// PrettyPrinting & Debugging functions

var tindent = [2]string{"  ", "                    "}

func (t *Tablet) PCurrentInstr() string {
	arg2 := ""
	if t.opargs[t.prog[t.p].op] == 2 {
		arg2 = t.PParam(t.prog[t.p].y)
	}
	return fmt.Sprintf("%s %v %v", t.opnames[t.prog[t.p].op], t.PParam(t.prog[t.p].x), arg2)
}

func (t *Tablet) PParam(x int) string {
	if x < REG {
		return itoa(x)
	}
	return fmt.Sprintf("%s(%d)", nameOf(x), valueOf(x, t))
}

func VPqueue(ll *LinkedList[int]) (s string) {
	s = "["
	if ll.head == nil {
		return "[]"
	}
	for c := ll.head; c != nil; c = c.next {
		s += " " + itoa(c.val)
	}
	return s + " ]"
}
