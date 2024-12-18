// Adventofcode 2024, d17, in go. https://adventofcode.com/2024/day/17
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example "4_6_3_5_6_3_5_2_1_0"
// TEST: -1 example2 "5_7_3_0"
// TEST: example2 117440
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

// Note: output separates the digits by _ instead of , for our testing tools

// For Part2, the output digits are detemined by groups of 3 bits of a
// So we test all values of these 3 bits (0 to 7) separately
// knowing that adding a higher value to a prepends output digits , keeping
// the last ones

package main

import (
	"flag"
	"fmt"
	"regexp"
	"strconv"
	// "slices"
)

type VM struct {
	a, b, c uint64				// registers
	code []int					// the program instructions opcodes
	p int						// the program current position in code
	out []int					// the current (cumulated) output
	l []Instr					// the language used: list of Instr by opcodes
	opnames []string			// for pretty-printing
}
type Instr func(*VM)			// the opcode actual code

var opnames = []string{"adv", "bxl", "bst", "jnz", "bxc", "out", "bdv", "cdv"}
var language = []Instr{OPadv, OPbxl, OPbst, OPjnz, OPbxc, OPout, OPbdv, OPcdv}
var outputsep = "_"
var commaFlag *bool
var seedFlag *uint64

//////////// Options

func main() {
	commaFlag = flag.Bool("c", false, "outputs numbers separated by comma instead of underscores")
	seedFlag = flag.Uint64("s", uint64(0), "seed: sets initial value of register A for part1")
	ParseOptionsString(2, part1, part2)
}

func ProcessXtraOptions() {
	if *commaFlag {
		outputsep = ","
	}
}

//////////// Part 1

func part1(lines []string) string {
	vm := parse(lines)
	if *seedFlag != 0 {
		vm.a = *seedFlag
	}
	RunVM(&vm)
	return Output(vm.out)
}

func RunVM(vm *VM) {
	for vm.p < len(vm.code) {
		vm.l[vm.code[vm.p]](vm)
		vm.p += 2
	}
}

func ResetVM(vm *VM, a, b, c uint64) {
	vm.a, vm.b, vm.c = a, b, c
	vm.p = 0
	vm.out = vm.out[0:0]
}

func ReRun(vm *VM, a, b, c uint64) {
	ResetVM(vm, a, b, c )
	RunVM(vm)
}

//////////// Part 2

type State struct {				// groups of 3-bit octal numbers making a
	segs []uint64
}

func part2(lines []string) string {
	vm := parse(lines)
	a, b, c := vm.a, vm.b, vm.c // initial state, used as seed for ReRun()
	queue := []State{}
	for i := 0; i < 8; i++ {
		queue = append(queue, State{[]uint64{uint64(i)}})
	}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		a = uint64(0)
		for i := len(cur.segs) - 1; i >= 0; i-- {
			s := cur.segs[i] << (3 * i)
			a = a | s
		}
		ReRun(&vm, a,  b, c)	// run for this value of a, ignore initial a
		if sliceIntEqualsLast(vm.code, vm.out) {
			// another last out digit found, save the 3-bits chunk used
			if len(vm.out) == len(vm.code) {
				return strconv.FormatUint(a, 10)
			}
			for i := 0; i < 8; i++ {
				nseg := make([]uint64, len(cur.segs))
				copy(nseg, cur.segs)
				nseg = append([]uint64{uint64(i)}, nseg...)
				queue = append(queue, State{nseg})
			}
		}
	}
	return ""
}

// does l2 matches the last elements of l1?
func sliceIntEqualsLast(l1, l2 []int) bool {
	from := len(l1) - len(l2)
	if from < 0 {
		return false
	}
	for i, v := range l2 {
		if v != l1[i + from] {
           return false
		}
	}
	return true
}

//////////// Common Parts code

func Combo(vm *VM) int {
	switch vm.code[vm.p+1] {
	case 0,1,2,3: return vm.code[vm.p+1]
	case 4: return int(vm.a)
	case 5: return int(vm.b)
	case 6: return int(vm.c)
	}
	panic(fmt.Sprintf("Invalid operand %d at opcode %d \"%s\"", vm.code[vm.p+1], vm.p, vm.opnames[vm.code[vm.p]]))
}

func Literal(vm *VM) int {
	return vm.code[vm.p+1]
}

func OPadv(vm *VM) {
	vm.a = uint64(int(vm.a) / intPower(2, Combo(vm)))
}

func OPbxl(vm *VM) {
	vm.b = uint64(int(vm.b % 8)  ^ Literal(vm))
}

func OPbst(vm *VM) {
	vm.b = uint64(Combo(vm) % 8)
}

func OPjnz(vm *VM) {
	if vm.a == 0 {
		return
	}
	vm.p = Literal(vm)
	vm.p -= 2					// cancels the automatic +2 on pointer
}

func OPbxc(vm *VM) {
	vm.b = vm.b ^ vm.c
}

func OPout(vm *VM) {
	vm.out = append(vm.out, Combo(vm) % 8)
}

func OPbdv(vm *VM) {
	vm.b = uint64(int(vm.a) / intPower(2, Combo(vm)))
}

func OPcdv(vm *VM) {
	vm.c = uint64(int(vm.a) / intPower(2, Combo(vm)))
}

func Output(numlist []int) (s string) {
	for i, n := range numlist {
		if i != 0 {
			s += outputsep
		}
		s += itoa(n)
	}
	return
}

func parse(lines []string) (vm VM) {
	renum := regexp.MustCompile("[[:digit:]]+")
	// we brute-force the parsing, just look for numbers at fixed places
	vm.a = uint64(atoi(renum.FindAllString(lines[0], -1)[0]))
	vm.b = uint64(atoi(renum.FindAllString(lines[1], -1)[0]))
	vm.c = uint64(atoi(renum.FindAllString(lines[2], -1)[0]))
	vm.code = atoil(renum.FindAllString(lines[4], -1))
	vm.l = language
	vm.opnames = opnames
	return
}

//////////// PrettyPrinting & Debugging functions

func DEBUG() {
	if ! debug { return }
	fmt.Println("DEBUG!")
	// ad hoc debug code here
}
