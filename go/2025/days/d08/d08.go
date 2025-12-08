// Adventofcode 2025, d08, in go. https://adventofcode.com/2025/day/08
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 40
// TEST: example 25272
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

// Input properties:
// all coordinates are between 0 and 100 000 (10e5), exclusive
// input file is 1000 boxes

package main

import (
	"fmt"
	"regexp"
	// "flag"
	"slices"
	"cmp"
)

// Implementation:
// To save space, we represent a point xyz by a single number:
// 10e10 x + 10e5 y + z
// Since each coord is < 10e5, and 10e15 is less than MaxInt (more than 10e18)
// we have room to fit into a 64-bit integer
// We also sort boxes, so that pairs are ordered p[0] < p[1] for compares

type Box int
type Pair [2]Box
type Circuit []Box
type Circuits []Circuit

const PN = 100000

//////////// Options parsing & exec parts

func main() {
	ExecOptions(2, NoXtraOpts, part1, part2)
}

// Heuristics on problem statement:
// How many shortest connections to connect based on the number of boxes in input
// 10 for example.txt which has 20 boxes, and 1000 for input with 1000 boxes.
func NumberToConnect(n int) int {
	if n <= 100 {
		return n / 2
	} else if n >= 10 {
		return n
	} else {
		return n*n / 2
	}
}

//////////// Part 1

func part1(lines []string) (res int) {
	boxes := parse(lines)
	slices.Sort(boxes)			// This ensures pairs are ordered
	pairs := MakePairs(boxes)
	slices.SortFunc(pairs, CmpPairs) // Pairs with shortest distance first
	circuits := Circuits{}
	circuits.ConnectPairsN(pairs, NumberToConnect(len(boxes)))
	slices.SortFunc(circuits, RCmpCircuits)
	return len(circuits[0]) * len(circuits[1]) * len(circuits[2])
}

func (circuits *Circuits) ConnectPairsN(pairs []Pair, max int) {
	for i := 0; i < len(pairs) && i < max; i++ {
		circuits.ConnectPair(pairs[i])
	}
	return
}

//////////// Part 2

func part2(lines []string) (res int) {
	boxes := parse(lines)
	slices.Sort(boxes)			// This ensures pairs are ordered
	pairs := MakePairs(boxes)
	slices.SortFunc(pairs, CmpPairs) // Pairs with shortest distance first
	circuits := Circuits{}
	lastpair := circuits.ConnectPairsFull(pairs, boxes)
	return lastpair[0].x() * lastpair[1].x()
}

func (circuits *Circuits) ConnectPairsFull(pairs []Pair, boxes []Box) Pair {
	connected := make(map[Box]bool, len(boxes)) // is the box already connected?
	for _, p := range pairs {
		circuits.ConnectPair(p)
		for i := range 2 {
			if ! connected[p[i]] {
				connected[p[i]] = true
			}
		}
		// Once all boxes are connected, we are done.
		// And we are sure there is only one circuit.
		if len(connected) == len(boxes) {
			return p
		}
	}
	panic("Could not connect all boxes")
}

//////////// Common Parts code

func MakePairs(boxes []Box) (pairs []Pair) {
	for i, p := range boxes {
		for j := i+1; j < len(boxes); j++ {
			q := boxes[j]
			pairs = append(pairs, Pair{p, q})
		}
	}
	return
}

func (circuits *Circuits) ConnectPair(p Pair) {
	VPf("Adding pair %s of distance %d\n", p.String(), p.Dist())
	added := -1			// index of circuit we connected the pair to
	for cix, c := range *circuits {
		if c.HasBox(p[0]) {
			if ! c.HasBox(p[1]) {
				VPf("  adding p1 since p0 already in %s\n", c.String())
				(*circuits)[cix] = append(c, p[1])
			}
			added = circuits.Merge(added, cix)
		} else {
			if c.HasBox(p[1]) {
				VPf("  adding p0 since p1 already in %s\n", c.String())
				(*circuits)[cix] = append(c, p[0])
				added = circuits.Merge(added, cix)
			}
		}
	}
	if added == -1 {				// create a new circuit with the orphan pair
		orphan := Circuit{p[0], p[1]}
		VPf("  made new circuit: %s\n", orphan.String())
		*circuits = append(*circuits, orphan)
	}
}

// merge from into to (by indexes), deleter from, create to if none (index -1)
func (circuits *Circuits) Merge(tox, fromx int) int {
	if tox == -1 {
		return fromx
	}
	for _, box := range (*circuits)[fromx] {
		if ! (*circuits)[tox].HasBox(box) {
			(*circuits)[tox] = append((*circuits)[tox], box)
		}
	}
	(*circuits)[fromx] = Circuit{}
	return tox
}

func (circuit Circuit) HasBox(box Box) bool {
	for _, b := range circuit {
		if b == box {
			return true
		}
	}
	//VPf("    box %v not in %v\n", box.String(), circuit.String(6))
	return false
}

func MakeBox(x, y, z int) Box {
	return Box(PN * PN * x + PN * y + z)
}
	
func (b *Box) xyz() (x, y, z int) {
	return b.x(), b.y(), b.z()
}
func (b *Box) x() int {
	return int(*b) / (PN * PN)
}
func (b *Box) y() int {
	return (int(*b) / PN) % PN
}
func (b *Box) z() int {
	return int(*b) % PN
}

func CmpPairs(p, q Pair) int {
	return cmp.Compare(p.Dist(), q.Dist())
}

func RCmpCircuits(c1, c2 Circuit) int {
	return cmp.Compare(len(c2), len(c1))
}

func (p *Pair) Dist() int {
	x0, y0, z0 := p[0].xyz()
	x1, y1, z1 := p[1].xyz()
	x, y, z := x1 - x0, y1 - y0, z1 - z0
	return x*x + y*y + z*z
}

func parse(lines []string) (points []Box) {
	reline := regexp.MustCompile("^[[:space:]]*([[:digit:]]+)[,[:space:]]+([[:digit:]]+)[,[:space:]]+([[:digit:]]+)[[:space:]]*$")
	for _, line := range lines {
		m := reline.FindStringSubmatch(line)
		points = append(points, MakeBox(atoi(m[1]),atoi(m[2]),atoi(m[3])))
	}
	return
}

//////////// PrettyPrinting & Debugging functions. See also the VPx functions.

func (b *Box) String() string {
	x,y,z := b.xyz()
	return fmt.Sprintf("%d,%d,%d", x, y, z)
}

func (p *Pair) String() string {
	return fmt.Sprintf("<%s-%s>", p[0].String(), p[1].String())
}

func (c *Circuit) String(maxs ...int) (s string) {
	maxlen := MaxInt
	if len(maxs) > 0 {
		maxlen = maxs[0]
	}
	first := true
	s = "["
	for i := 0; i < len(*c) && i < maxlen; i++ {
		if ! first {
			s += "|"
		}
		s += (*c)[i].String()
		first = false
	}
	if maxlen < len(*c) {
		s += "|..."
	}
	s += "]"
	return
}

func (circuits *Circuits) String() (s string) {
	first := true
	s = "{"
	for _, c := range (*circuits) {
		if ! first {
			s += " "
		}
		s += c.String()
		first = false
	}
	s += "}"
	return
}

func (circuits *Circuits) VP() {
	if !verbose {
		return
	}
	for _, c := range (*circuits) {
		fmt.Printf("[%d] %s\n", len(c), c.String())
	}
}

func DEBUG() {
	if ! debug { return }
	fmt.Println("DEBUG!")
	// ad hoc debug code here
}
