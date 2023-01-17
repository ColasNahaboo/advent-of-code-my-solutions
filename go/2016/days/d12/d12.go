// Adventofcode 2016, d12, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 example 42
// TEST: -1 input 318117
// TEST: example 42
// TEST: input 9227771

// Note at first I tested a naive approach: interpreting the code,
// it ran in 42 seconds. I left the code here as interpretCode() for reference.
// But by compiling the code lines into a simple machine language (parseCode),
// and then executing it (execCode)
// the total ran in less than 0.1s

package main

import (
	"flag"
	"fmt"
	// "regexp"
)

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
	//return interpretCode(lines, [4]int{0,0,0,0})
	return execCode(parseCode(lines), [4]int{0,0,0,0})
}

//////////// Part 2
func part2(lines []string) int {
	//return interpretCode(lines, [4]int{0,0,1,0})
	return execCode(parseCode(lines), [4]int{0,0,1,0})
}

//////////// Common Parts code

func interpretCode(lines []string, regs [4]int) int {
	lineno := 0
	var n, n2 int
	var regname, regname2 string
	for lineno < len(lines) {
		line := lines[lineno]
		if nf, _ := fmt.Sscanf(line, "cpy %d %1s", &n, &regname); nf == 2 {
			regs[regid(regname)] = n
		} else if nf, _ := fmt.Sscanf(line, "cpy %1s %1s", &regname2, &regname); nf == 2 {
			regs[regid(regname)] = regs[regid(regname2)]
		} else if nf, _ := fmt.Sscanf(line, "inc %1s", &regname); nf == 1 {
			regs[regid(regname)]++
		} else if nf, _ := fmt.Sscanf(line, "dec %1s", &regname); nf == 1 {
			regs[regid(regname)]--
		} else if nf, _ := fmt.Sscanf(line, "jnz %1d %d", &n2, &n); nf == 2 {
			if n2 != 0 {
				lineno += n
				continue
			}
		} else if nf, _ := fmt.Sscanf(line, "jnz %1s %d", &regname, &n); nf == 2 {
			if regs[regid(regname)] != 0 {
				lineno += n
				continue
			}
		} else {
			panic(fmt.Sprintf("Syntax error line %d: %s\n", lineno, line))
		}
		lineno++
	}
	return regs[0]
}

func regid(s string) int {
	return int(s[0] - 'a')
}

// code instructions are 3 ints: op reg value
// op 0=cpy int, 1=cpy reg, 2=inc val, 3=jnz int, 4=jnz reg
type Instr struct {
	op, reg, int int
}

func parseCode(lines []string) (code []Instr) {
	lineno := 0
	var n, n2 int
	var regname, regname2 string
	for lineno < len(lines) {
		line := lines[lineno]
		if nf, _ := fmt.Sscanf(line, "cpy %d %1s", &n, &regname); nf == 2 {
			code = append(code, Instr{0, regid(regname), n})
		} else if nf, _ := fmt.Sscanf(line, "cpy %1s %1s", &regname2, &regname); nf == 2 {
			code = append(code, Instr{1, regid(regname), regid(regname2)})
		} else if nf, _ := fmt.Sscanf(line, "inc %1s", &regname); nf == 1 {
			code = append(code, Instr{2, regid(regname), 1})
		} else if nf, _ := fmt.Sscanf(line, "dec %1s", &regname); nf == 1 {
			code = append(code, Instr{2, regid(regname), -1})
		} else if nf, _ := fmt.Sscanf(line, "jnz %d %d", &n2, &n); nf == 2 {
			code = append(code, Instr{3, n2, n})
		} else if nf, _ := fmt.Sscanf(line, "jnz %1s %d", &regname, &n); nf == 2 {
			code = append(code, Instr{4, regid(regname), n})
		} else {
			panic(fmt.Sprintf("Syntax error line %d: %s\n", lineno, line))
		}
		lineno++
	}
	return
}

func execCode(code []Instr, regs [4]int) int {
	size := len(code)
	lineno := 0
	for lineno < size {
		switch code[lineno].op {
		case 0: regs[code[lineno].reg] = code[lineno].int
		case 1: regs[code[lineno].reg] = regs[code[lineno].int]
		case 2: regs[code[lineno].reg] += code[lineno].int
		case 3: if code[lineno].reg != 0 {
			lineno += code[lineno].int
			continue
		}
		case 4:if regs[code[lineno].reg] != 0 {
			lineno += code[lineno].int
			continue
		}
		}
		lineno++
	}
	return regs[0]
}
	
//////////// Part1 functions

//////////// Part2 functions
