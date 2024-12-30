// Adventofcode 2024, d24, in go. https://adventofcode.com/2024/day/24
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example1 4
// TEST: -1 example 2024
// TEST: -1 example2 9
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

// WARNING: Part2 is not done yet!
// I solved it by hand for now by printing the graph
// I plan to design a proper solution later, after the holidays with the family
// The wires+gates implement a full adder
// see https://www.geeksforgeeks.org/full-adder-in-digital-logic/
// The solution is to recognise this pattern for each Z digit, and spot
// the misplaced wires.

package main
	
import (
	"fmt"
	"regexp"
	"flag"
	"slices"
	"os"
	"github.com/olekukonko/tablewriter"
)

type Gate struct {
	id Gid
	op OPid
	opname string
	in1, in2 Wid
	out Wid
}
type Wire struct {
	id Wid
	name string
	ready bool
	value bool
	ingate Gid
}
type Gid int
type Wid int
type OPid int
var gates = []Gate{}
var wires = []Wire{}
var	xnums = []Wid{}				// for part2
var	ynums = []Wid{}				// for part2
var	znums = []Wid{}
var wiresID = make(map[string]Wid) // wire name ==> index in wires
var opsID = []string{"AND", "OR", "XOR"}
var opsFunc = []func(i, j bool)bool{opAND, opOR, opXOR}
const (
	AND = 0
	OR = 1
	XOR = 2
)

//////////// Options parsing & exec parts

var usage = `
part3 prints in binary format x, y, x+y, and z, to the errors manually
part4 prints a graphical representation of the wires and gates
`
var commaFlag *bool
var commaSep = "_"
var swapsFlag string
func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:%s", os.Args[0], usage)
		flag.PrintDefaults()
	}
	commaFlag = flag.Bool("c", false, "outputs numbers separated by comma instead of underscores")
	flag.StringVar(&swapsFlag, "s", "", "swaps the pairs of wires, comma-separated, e.g: -s z00-z05,z02-z01")
	flag.BoolVar(&NSPsilent, "q", false, "quiet operation for part3 (testing adder)")
	ExecOptions(2, XtraOpts, part1, part2, part3, part4)
}

func XtraOpts() { // extra options, see ParseOptions in utils.go
	if *commaFlag {
		commaSep = ","
	}
}

//////////// Part 1

func part1(lines []string) (res int) {
	parse(lines)
	return DigitsValue(znums)
}

// reverse-chain: force upstream computations "on demand"
func On(wid Wid, level int) bool {
	if wires[wid].ready {
		return wires[wid].value
	}
	if wires[wid].ingate < 0 {
		if NSPsilent {
			os.Exit(1)
		} else {
			panicf("Cannot find a gate to feed wire %s", wires[wid].name)
		}
	}
	if level < 0 {
		if NSPsilent {
			os.Exit(2)
		} else {
			panicf("On: computing wires stack overflow")
		}
	}
	gate := gates[wires[wid].ingate]
	VPf("== On %s from %s %s %s\n", wires[wid].name, wires[gate.in1].name, gate.opname, wires[gate.in2].name)
	wires[wid].value = opsFunc[gate.op](On(gate.in1, level-1), On(gate.in2, level-1))
	wires[wid].ready = true
	return wires[wid].value
}

func DigitsValue(nums []Wid) (v int) {
	bitposition := 1
	for id := range nums {
		On(nums[id], 1000)
		if wires[nums[id]].value {
			v += bitposition
		}
		bitposition *= 2
	}
	return
}

func opAND(i, j bool) bool { return i && j }
func opOR(i, j bool) bool { return i || j }
func opXOR(i, j bool) bool { return i != j}

//////////// Part 2

func part2(lines []string) (res int) {
	parse(lines)
	return 
}

//////////// Common Parts code

func DeclareWire(name string, values ...bool) (wid Wid) {
	var ok bool
	if wid, ok = wiresID[name]; ! ok {
		wid = Wid(len(wires))
		wire := Wire{id: Wid(wid), name: name, ingate: -1}
		wires = append(wires, wire)
	}
	if len(values) > 0 {		// some juice!
		wires[wid].ready = true
		wires[wid].value = values[0]
	}
	switch name[0] {
	case 'x': xnums = DeclareDigit(xnums, name, wid)
	case 'y': ynums = DeclareDigit(ynums, name, wid)
	case 'z': znums = DeclareDigit(znums, name, wid)
	}
	wiresID[name] = wid
	return		
}

func DeclareDigit(nums []Wid, name string, wid Wid) []Wid {
	id := atoi(name[1:])
	for len(nums) <= id {
		nums = append(nums, Wid(0))
	}
	nums[id] = wid
	return nums
}

func DeclareGate(outname, opname, inname1, inname2 string) (gid Gid) {
	in1 := DeclareWire(inname1)
	in2 := DeclareWire(inname2)
	if inname1 > inname2 {		// keep args in alphabetical order
		in1, in2 = in2, in1
	}
	out := DeclareWire(outname)
	op := OPid(slices.Index(opsID, opname))
	gid = Gid(len(gates))
	gate := Gate{id: gid, op: op, opname: opname, in1: in1, in2: in2, out: out}
	gates = append(gates, gate)
	wires[out].ingate = gid
	return
}

func ParseSwaps(swaps string) map[string]string {
	reswap := regexp.MustCompile("([a-z0-9]{3})-([a-z0-9]{3})")
	sm := make(map[string]string)
	if ms := reswap.FindAllStringSubmatch(swapsFlag, -1); ms != nil {
		for _, m := range ms {
			sm[m[1]] = m[2]
			sm[m[2]] = m[1]
		}
	}
	return sm
}

func parse(lines []string) {
	reinit := regexp.MustCompile("^([a-z0-9]{3}): +([[:digit:]]+)")
	regate := regexp.MustCompile("^([a-z0-9]{3}) (AND|OR|XOR) ([a-z0-9]{3}) -> ([a-z0-9]{3})")
	swaps := ParseSwaps(swapsFlag)
	for lineno, line := range lines {
		if m := reinit.FindStringSubmatch(line); m != nil {
			DeclareWire(m[1], m[2] == "1")
		} else if m = regate.FindStringSubmatch(line); m != nil {
			if repl := swaps[m[4]]; len(repl) > 0 {
				NSPf("Swapping %s ==> %s\n", m[4], repl)
				DeclareGate(repl, m[2], m[1], m[3])
			} else {
				DeclareGate(m[4], m[2], m[1], m[3])
			}
		} else if line != "" {
			panicf("Syntax error line %d: \"%s\"", lineno+1, line)
		}
	}
	VPf("== Parsed %d wires, %d gates, %d z-digits\n", len(wires), len(gates), len(znums))
}

//////////// Part3
// Part3 is for debugging: it just checks that we get X + Y = Z on the input

func part3(lines []string) (res int) {
	parse(lines)
	errs := ValidAdderErrs(NSPsilent)
	if ! NSPsilent{
		fmt.Println(errs)
	}
	os.Exit(errs)
	return
}

//  returns number of bits in error in z compared to x+y
func ValidAdderErrs(silent bool) int {
	NSPsilent = silent
	x := DigitsValue(xnums)
	y := DigitsValue(ynums)
	z := DigitsValue(znums)
	NSPf("x:  %64b\ny:  %64b\nz:  %64b\nx+y:%64b\n", x, y, z, x+y)
	if x + y == z {
		NSPf("OK! %d + %d == %d\n", x, y, z)
	} else {
		got := fmt.Sprintf("%64b", z)
		exp := fmt.Sprintf("%64b", x+y)
		dif := ""
		for i, c := range exp {
			if byte(c) != got[i] {
				dif += "#"
			} else {
				dif += " "
			}
		}
		NSPf("err:%s\n", dif)
		NSPf("*BAD* %d+%d is %d, must be %d\n", x, y, z, x+y)
	}
	return DifferentBits(z, x+y)
}

func DifferentBits(i, j int) (db int) {
	for _ = range 64 {
		if (i & 1) != (j & 1) {
			db++
		}
		i >>= 1
		j >>= 1
	}
	return
}

var NSPsilent bool					// set to use ValidAdder internally
func NSPf(f string, v ...interface{}) {
	if ! NSPsilent {
		fmt.Printf(f, v...)
	}
}

//////////// Part4
// Part4 is for debugging: prints the graph of the gates and wires

type TWCell struct {
	x, y, h int
}
func part4(lines []string) (res int) {
	parse(lines)
	table := tablewriter.NewWriter(os.Stdout)
	for _, zid := range znums {
		v := []string{wires[zid].name}
		table.Append(v)
	}
	table.Render() 
	os.Exit(0)					// avoid printing the res
	return
}

//////////// PrettyPrinting & Debugging functions. See also the VPx functions.

func DEBUG() {
	if ! debug { return }
	fmt.Println("DEBUG!")
	// ad hoc debug code here
}
