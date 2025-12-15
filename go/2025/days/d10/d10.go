// Adventofcode 2025, d10, in go. https://adventofcode.com/2025/day/10
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 7
// TEST: example 33
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

// Input properties:
// - less than 200 machines
// - each 12 leds or less
// - led starting status is never all off ([....])
// - joltages are 0..255

package main

import (
	"fmt"
	"regexp"
	"math"
	// "flag"
	"slices"
	// for part2withZ3Solver
	"github.com/aclements/go-z3/z3"
)

// Implementation:
// We call indicator lights "leds" for terseness
// We represent the state of the leds as simple bitfields in an int
// Buttons are also bitfield-uints of what leds they toggle
// Thus pushing a button is just XOR-ing it with the led state
// We use ints and not uints, as is the usage in Go

type Machine struct {
	nleds int
	leds int
	buttons []int				// buttons are bitfields of leds toggled
	buttonsInt [][]int			// human readable form: list of led indexes
	jolts []int
}
const MAXLEDS = 12				// to use fixed size arrays to represent states

//////////// Options parsing & exec parts

func main() {
	ExecOptions(3, NoXtraOpts, part1, part2, part2withZ3Solver, stats)
}

//////////// Part 1

func part1(lines []string) (res int) {
	machines := parse(lines)
	for _, m := range machines {
		VPmachine(m)
		if n := FewestPressesStart(m); n > 0 {
			VPf("  %d\n", n)
			res += n
		}
	}
	return 
}

// recurse breadth-first into the possibilties of button presses
// An action to do:
type Todo struct {
	leds int					// status of lit leds
	button int					// button to press
	presses int					// number of buttons already pressed to get there
}
func FewestPressesStart(m Machine) (minpresses int) {
	// the FIFO queue of actions to do. En-queue the first possble presses
	todo := []Todo{}
	seen := make(map[int]bool)	// already seen?
	seen[0] = true
	for _, nextb := range m.buttons {
		todo = append(todo, Todo{0, nextb, 0})
	}
	for {
		// de-queue next action to try from the todo queue
		leds, b, presses := todo[0].leds, todo[0].button, todo[0].presses
		todo = todo[1:]
		// we push the button
		pushed := leds ^ b
		seen[pushed] = true
		if debug {VPf("  [%s] level %d, [%s] after pushing %s\n", leds2string(m.nleds, leds), presses, leds2string(m.nleds, pushed), button2string(m.nleds, b))}
		if pushed == m.leds {	// found!
			return presses + 1	// This is the minimum number, as we BFS
		}
		if pushed == 0 {		// we loop, abort this branch
			continue
		}
		// else en-queue the exploration of next pushes for the next depth level
		for _, nextb := range m.buttons {
			if ! seen[pushed ^ nextb] {
				todo = append(todo, Todo{pushed, nextb, presses+1})
			}
		}
	}
}

//////////// Part 2
// A graph-exploration implementation with A*
// A node is a the state of the joltages (slice of ints)
// The graph is the machine
// connected nodes are the ones reachable by one button press

func part2(lines []string) (res int) {
	machines := parse(lines)
	for i, m := range machines {
		VPf("Testing #%d: ", i)
		VPmachine(m)
		if n := FewestPressesJolts(m); n > 0 {
			VPf("  ==> %d\n", n)
			res += n - 1
		}
	}
	return 
}

// recurse breadth-first into the possibilties of button presses
// An action to do:

type State [MAXLEDS]byte

type TodoJolts struct {
	counters [MAXLEDS]byte	  // joltages of counters, fastest to handle
	button int				  // button to press
	presses int				  // number of buttons already pressed to get there
}
func FewestPressesJolts(m Machine) (minpresses int) {
	zero := State{}
	target := MakeState(m.jolts)
	
	pressed := AStarFindPath[*Machine, State](&m, zero, target,
		NextStates, AdjacentDist, StatesDist, EndState)

	return len(pressed)
}

func MakeState(jolts []int) (state State) {
	for i, c := range jolts {
		state[i] = byte(c)
	}
	return
}

func State2Ints(state State) (jolts []int) {
	for _, c := range state {
		if c != byte(0) {
			jolts = append(jolts, int(c))
		}
	}
	return jolts
}

func NextStates(m *Machine, s State) (nexts []State) {
BUTTON:
	for _, button := range m.buttonsInt {
		state := s
		for _, led := range button {
			if int(state[led]) >= m.jolts[led] { // overflows are fatal
				continue BUTTON
			}
			state[led] = s[led] + 1
		}
		nexts = append(nexts, state)
	}
	return
}

func StatesDist(m *Machine, s1, s2 State) (dist float64) {
	for _, button := range m.buttonsInt {
		for _, led := range button {
			dist += math.Pow(float64(intAbs(int(s1[led]) - int(s2[led]))), 2)
		}
	}
	return math.Sqrt(dist)
}

func EndState(m *Machine, s, e State) bool {
	return s == e
}

func SumSliceInt(s []int) (sum int) {
	for _, e := range s {
		sum += e
	}
	return
}

//////////// Optimization Research approaches

// a + b + d = 43  <=>  SPEquation{{0,1,3},43}
type SPEquation struct {			// Set Partitiong Equations, for Integer Prog
	vars []int						// List of indexes of var with 1 coeffs
	rhs  int						// integer non-variable
}

//////////// Part 3
// We use the Z3 theorem solver

func part2withZ3Solver(lines []string) (res int) {
	machines := parse(lines)
	for i, m := range machines {
		VPf("Testing #%d: ", i)
		VPmachine(m)
		equations := SetPartitioningEquations(m)
		VPSetPartitioningEquations(equations)

		// Z3
		config := z3.NewContextConfig()
		ctx := z3.NewContext(config)
		// Create the solver. Alas, no z3.NewOptimize in this Go API!
		slv := z3.NewSolver(ctx)
		// Declare useful constants
		zero := ctx.FromInt(0, ctx.IntSort()).(z3.Int)
		// Declare the variables: the number of pushes on each button
		vars := []z3.Int{}
		varnames := []string{}
		for vx := range len(m.buttons) {
			vname := string('a' + vx)
			varnames = append(varnames, vname)
			v := ctx.IntConst(vname)
			vars = append(vars, v)
			slv.Assert(v.GE(zero))
		}
		// Declare Equations
		for _, equation := range equations {
			rhs := ctx.FromInt(int64(equation.rhs), ctx.IntSort()).(z3.Int)
			lhs := []z3.Int{}	// the variables summed in the left hand side
			for _, vx := range equation.vars {
				lhs = append(lhs, vars[vx])
			}
			slv.Assert(lhs[0].Add(lhs[1:]...).Eq(rhs))
		}
		sum := ctx.IntConst("sum")
		slv.Assert(sum.Eq(vars[0].Add(vars[1:]...)))
		// Optimize by hand: Solve, and look for solutions with smaller sums
		presses := 0
		for {
			sat, err := slv.Check()
			if err != nil {
				VPf("Machine #%d Invalid\n", i)
				return
			}
			if ! sat {
				if presses == 0 {
					VPf("Machine #%d Unsolvable\n", i)
					return
				}
				break
			}
			solution := slv.Model()
			presses = atoi(solution.Eval(sum, true).String())
			// tighten the constraints to re-try with smaller sum
			slv.Assert(sum.LT(solution.Eval(sum, true).(z3.Int)))
		}
		// Get result: the sum of variable values in solution
		VPf(" total=%d\n", presses)
		res += presses
	}
	return 
}

func VPSetPartitioningEquations(eqs []SPEquation) {
	if ! verbose {return}
	for _, eq := range eqs {
		for i, varx := range eq.vars {
			if i > 0 {
				fmt.Print(" + ")
			}
			fmt.Print(string('a' + varx))
		}
		fmt.Printf(" = %d\n", eq.rhs)
	}
}

// create the Partitioning Equations to solve a machine
// They are Systems of Linear Equations with coeffs only 0 or 1
// bounds is the range of values that each non-negative var can take, inclusives

func SetPartitioningEquations(m Machine) (equations []SPEquation) {
	for led, jolt := range m.jolts {
		coeffs := make([]int, len(m.buttonsInt))
		for bi, b := range m.buttonsInt {
			if slices.Contains(b, led) {
				coeffs[bi]++
			}
		}
		eq := SPEquation{rhs: jolt}
		for varx, coeff := range coeffs {
			if coeff > 0 {
				eq.vars = append(eq.vars, varx)
			}
		}
		equations = append(equations, eq)
	}
	return
}

// return the possible values for each variables: inside the inclusive bounds
func VarBounds(m Machine) (bounds [][]int) {
	for vx, jolt := range m.jolts {
		for _ = range m.leds {
			// 8...8 is nearly "MaxInt" but human-friendly for debugging
			bounds = append(bounds, []int{0, 8888888888888888888})
		}
		for _, b := range m.buttonsInt {
			if slices.Contains(b, vx) {
				if jolt < bounds[vx][1] {
					bounds[vx][1] = jolt
				}
			}
		}
	}
	return
}

//////////// Part 4: just perform some stats on the input

func stats(lines []string) (res int) {
	machines := parse(lines)
	var maxFreeVars, maxfvm int
	for i, m := range machines {
		equations := SetPartitioningEquations(m)
		freevars := len(m.buttonsInt) - m.nleds
		VPSetPartitioningEquations(equations)
		if freevars > maxFreeVars {
			maxFreeVars, maxfvm = freevars, i
		}
	}
	fmt.Printf("max free variables: %d, for machine #%d\n", maxFreeVars, maxfvm)
	return maxFreeVars
}

//////////// Common Parts code

var renum = regexp.MustCompile("[[:digit:]]+")
var rebutton = regexp.MustCompile("[(][,[:digit:]]+[)]")

func parse(lines []string) (machines []Machine) {
	reline := regexp.MustCompile("^[[]([.#]+)[]]" +
		"([(),[:digit:][:space:]]+)" +
		"[{]([,[:digit:][:space:]]+)[}]")
	for _, line := range lines {
		var m []string
		if m = reline.FindStringSubmatch(line); m == nil {
			panic("Syntax error in input")
		}
		nleds := len(m[1])
		leds := parseLeds(m[1])
		buttons := parseButtons(m[2])
		buttonsInt := parseButtonsInt(m[2])
		jolts := parseJolts(m[3])
		if len(jolts) != nleds {
			panic(itoa(len(jolts)) + " joltages, for " + itoa(nleds) + " leds!")
		}
		machines = append(machines, Machine{nleds, leds, buttons, buttonsInt, jolts})
	}
	return
}


func parseLeds(s string) (leds int) {
	bit := 1
	for _, r := range s {
		if r == '#' {
			leds |= bit
		}
		bit <<= 1
	}
	return
}
	
func parseButtons(s string) (buttons []int) {
	var m []string
	if m = rebutton.FindAllString(s, -1); m == nil {
		panic("Syntax error in buttons")
	}
	for _, b := range m {
		var button int
		for _, wire := range renum.FindAllString(b, -1) {
			button |= 1 << atoi(wire) // sets bit "wire"
		}
		buttons = append(buttons, button)
	}
	return
}
	
func parseButtonsInt(s string) (buttons [][]int) {
	var m []string
	if m = rebutton.FindAllString(s, -1); m == nil {
		panic("Syntax error in buttons")
	}
	for _, b := range m {
		button := []int{}
		for _, wire := range renum.FindAllString(b, -1) {
			button = append(button, atoi(wire))
		}
		buttons = append(buttons, button)
	}
	return
}
	
func parseJolts(s string) (jolts []int) {
	var m []string
	if m = renum.FindAllString(s, -1); m == nil {
		panic("Syntax error in jolts")
	}
	for _, n := range m {
		jolts = append(jolts, atoi(n))
	}
	return
}
	

//////////// PrettyPrinting & Debugging functions. See also the VPx functions.

func VPmachine(m Machine) {
	if ! verbose {return}
	fmt.Printf("[%s]", leds2string(m.nleds, m.leds))
	for _, b := range m.buttons {
		fmt.Printf(" %s", button2string(m.nleds, b))
	}
	fmt.Printf(" %s\n", jolts2string(m.jolts))
}

func leds2string(n, leds int) (s string) {
	for i := range n {
		if leds & (1 << i) != 0 {
			s += "#"
		} else {
			s += "."
		}
	}
	return
}

func button2string(n, b int) (s string) {
	s = "("
	for i := range n {
		if b & (1 << i) != 0 {
			if len(s) > 1 {s += ","}
			s += itoa(i)
		}
	}
	s += ")"
	return
}

func jolts2string(jolts []int) (s string) {
	s = "{"
	for _, j := range jolts {
		if len(s) > 1 {s += ","}
		s += itoa(j)
	}
	s += "}"
	return
}

func counters2string(n int, counters [MAXLEDS]byte) (s string) {
	s = "<"
	for i := range n {
		if len(s) > 1 {s += " "}
		s += fmt.Sprintf("%4d", counters[i])
	}
	s += ">"
	return
}


func DEBUG() {
	if ! debug { return }
	fmt.Println("DEBUG!")
	// ad hoc debug code here
}
