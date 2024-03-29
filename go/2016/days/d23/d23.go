// Adventofcode 2016, d23, in go. https://adventofcode.com/2016/day/23
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example
// TEST: example
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"flag"
	"fmt"
	"regexp"
)

var verbose bool

// code instructions are 3 ints: opcode and 2 parameters
type Instr struct {
	op int
	x, y Value
	valid bool			   // used to skip invalid expressions created by a TGL
}
// parameter values can be either an integer in NONE..REGV (excluded),
// or an index in registers, as REGV plus the id of the register
// or NONE to indicate this parameter was not present at parsing
const REGV =  9000000000000000000
const NONE = -8888888888888888888				// non-existent parameter at parsing
type Value int
var regnames map[string]Value = map[string]Value{
	"a": REGV, "b": REGV+1, "c": REGV+2, "d": REGV+3,
}
	
// opcodes, op IDs
const (
	CPY = 0
	INC = 1
	DEC = 2
	JNZ = 3
	TGL = 4
)
var opnames = []string{"cpy", "inc", "dec", "jnz", "tgl"}

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	flag.Parse()
	verbose = *verboseFlag
	var infile string
	if flag.NArg() > 0 {
		infile = flag.Arg(0)
	} else {
		infile = fileMatch("input,[0-9]*,[0-9]*.test")
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
	return execCode(parseCode(lines), [4]int{7,0,0,0})
}

//////////// Part 2
func part2(lines []string) int {
	return execCode(parseCode(lines), [4]int{12,0,0,0})
}

//////////// Common Parts code

//////////// Parsing
// assembunny code and monorail computer similar to d12, but coded differently

func parseCode(lines []string) (code []Instr) {
	re := regexp.MustCompile("([[:alpha:]]{3}) +([-[:alnum:]]+) *([-[:alnum:]]*)")
	for lineno, line := range lines {
		m := re.FindStringSubmatch(line)
		if m == nil {
			panic(fmt.Sprintf("Syntax error line %d: %s\n", lineno, line))
		}
		instr := instrMake(m[0], m[1], m[2], m[3])
		if ! isInstrValid(instr) {
			panic(fmt.Sprintf("Invalid operation line %d: %s %v\n", lineno, line, instr))
		}
		code = append(code, instr)
		if verbose {
			VPinstr(lineno, instr, [4]int{0,0,0,0})
		}
	}
	return
}

func instrMake(line, opname, x, y string) Instr {
	op := IndexOf[string](opnames, opname)
	if op == -1 {
		panic(fmt.Sprintf("Unknown op \"%s\", line %s\n", opname, line))
	}
	return Instr{op, valueMake(x), valueMake(y), true}
}

func valueMake(name string) Value {
	val, isreg := regnames[name]
	if isreg {
		return val
	} else if name == "" {
		return NONE
	} else {
		return Value(atoi(name))
	}
}

//////////// Misc utils

// check validity of instruction: type and number of args for each op
func isInstrValid(i Instr) bool {
	argnums := 2
	if i.y == NONE {
		if i.x == NONE {
			argnums = 0
		} else {
			argnums = 1
		}
	}
	switch i.op {
	case CPY:
		if ! isReg(i.y) {
			return false
		}
		if argnums != 2 {
			return false
		}
	case INC, DEC:
		if ! isReg(i.x) {
			return false
		}
		if argnums != 1 {
			return false
		}
	case JNZ:
		if argnums != 2 {
			return false
		}
	case TGL:
		if argnums != 1 {
			return false
		}
	default:
		return false
	}
	return true
}

func isReg(val Value) bool {
	return val >= REGV
}

func isInt(val Value) bool {
	return val < REGV && val > NONE
}

func isNone(val Value) bool {
	return val == NONE
}

func regId(val Value) int {
	return int(val - REGV)
}

func evalValue(val Value, regs [4]int) int {
	if isReg(val) {
		return regs[regId(val)]
	} else {
		return int(val)
	}
}

func printValue(val Value, regs [4]int) string {
	if isReg(val) {
		return string('a' + int(val) - REGV) + "=" + itoa(regs[regId(val)])
	} else {
		return itoa(int(val))
	}
}


//////////// Executing

func execCode(code []Instr, regs [4]int) int {
	codeSize := len(code)
	lineno := 0
	for lineno < codeSize {
		i := code[lineno]
		if verbose {
			VPinstr(lineno, i, regs)
		}
		if ! i.valid {		// skip silently invalid Instrs generated by TGLs
			goto NEXTLINE
		}
		switch i.op {
		case CPY:
			regs[regId(i.y)] = evalValue(i.x, regs)
		case INC:
			regs[regId(i.x)]++
		case DEC:
			regs[regId(i.x)]--
		case JNZ:
			if evalValue(i.x, regs) != 0 {
				lineno += evalValue(i.y, regs)
				continue
			}
		case TGL:
			targetno := lineno + evalValue(i.x, regs)
			if targetno < 0 || targetno >= codeSize {
				goto NEXTLINE
			}
			t := code[targetno]
			if t.y == NONE {	// one argument --> inc
				if t.op == INC {
					code[targetno].op = DEC
				} else {
					code[targetno].op = INC
				}
			} else {			// two arguments --> jnz
				if t.op == JNZ {
					code[targetno].op = CPY
				} else {
					code[targetno].op = JNZ
				}
			}
			// mark invalid results, but leave them as is for future TGLs
			if isInstrValid(code[targetno]) {
				code[targetno].valid = true
			} else {
				code[targetno].valid = false
			}
		}
	NEXTLINE:
		lineno++
	}
	return regs[0]
}

//////////// PrettyPrinting & Debugging functions

func VPinstr(lineno int, i Instr, regs [4]int) {
	s := opnames[i.op] + " " + printValue(i.x, regs)
	if ! isNone(i.y) {
		s += " " + printValue(i.y, regs)
	}
	fmt.Printf("  [%d] %v %s\n", lineno, regs, s)
}

