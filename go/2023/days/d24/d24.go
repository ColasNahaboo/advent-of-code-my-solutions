// Adventofcode 2023, d24, in go. https://adventofcode.com/2023/day/24
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 2
// TEST: example
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"flag"
	"fmt"
	"strings"
	"math/big"
)

var verbose bool

type Stonef struct {
	px, py, pz, vx, vy, vz float64
}
var stonesf []Stonef
var testmin, testmax float64

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

func part1(lines []string) (sum int) {
	parse(lines)
	for i, s := range stonesf[0:len(stonesf)-1] {
		for _, s2 := range stonesf[i+1:] {
			x, y := XYintersection(s, s2)
			if inFuture(x, s.px, s.vx) && inFuture(x, s2.px, s2.vx) && inTest(x, y) {
				VPf("  OK for %v + %v\n", s, s2)
				sum++
			}
		}
	}
	return
}

// we add to the input a first line of dims of test area
func parse(lines []string) {
	var px, py, pz, vx, vy, vz float64
	_, err := fmt.Sscanf(lines[0], "TESTAREA-XY: %f %f", &testmin, &testmax)
	if err != nil { panic("Syntax error testarea")}
	for lineno, line := range lines[1:] {
		s := strings.ReplaceAll(line, "  ", " ") // poor man parsing
		n, err := fmt.Sscanf(s, "%f, %f, %f @ %f, %f, %f", &px, &py, &pz, &vx, &vy, &vz)
		if err != nil {
			panic(fmt.Sprintf("Syntax error line %d: \"%s\" (%d parsed)\n", lineno+2, s, n))
		}
		stonesf = append(stonesf, Stonef{px, py, pz, vx, vy, vz})
	}
}

func XYintersection(s1, s2 Stonef) (x, y float64) {
	x = (s1.vx * s2.vx * s2.py - s1.vx * s2.vy * s2.px - s1.py * s1.vx * s2.vx + s1.vy * s1.px * s2.vx) / (s1.vy * s2.vx - s1.vx * s2.vy)
	y = x * s1.vy / s1.vx + s1.py - s1.vy * s1.px / s1.vx
	return
}

func inTest(x, y float64) bool {
	return x >= testmin && x <= testmax && y >= testmin && y <= testmax
}

// intersection in the future if dx and vx have same sign
func inFuture(x, px, vx float64) bool {
	return (x - px) * vx > 0
}

//////////// Part 2

type Stone struct {
	p [3]int
	v [3]int
}
var stones []Stone
var axisname = [3]string{"X", "Y", "Z"}

func part2(lines []string) (sum int) {
	parse2(lines)
	var samepv []int
	sameaxis := -1
	for axis := range axisname {
		samepv = findSamePosVelStones(axis)
		if len(samepv) > 0 {
			sameaxis = axis
			break
		}
	}
	if sameaxis == -1 {
		panic("Could not find 2 hailstones with same pos and vel in one axis")
	}
	VPf("Found HA and HB with same pos and vel on axis %s: %v %v\n", axisname[sameaxis], stones[samepv[0]], stones[samepv[1]])
	rock := solve(sameaxis, stones[samepv[0]], stones[samepv[1]])
	return rock.p[0] + rock.p[1] + rock.p[2]
}

func parse2(lines []string) {
	var px, py, pz, vx, vy, vz int
	_, err := fmt.Sscanf(lines[0], "TESTAREA-XY: %f %f", &testmin, &testmax)
	if err != nil { panic("Syntax error testarea")}
	for lineno, line := range lines[1:] {
		s := strings.ReplaceAll(line, "  ", " ") // poor man parsing
		n, err := fmt.Sscanf(s, "%d, %d, %d @ %d, %d, %d", &px, &py, &pz, &vx, &vy, &vz)
		if err != nil {
			panic(fmt.Sprintf("Syntax error line %d: \"%s\" (%d parsed)\n", lineno+2, s, n))
		}
		stones = append(stones, Stone{[3]int{px, py, pz}, [3]int{vx, vy, vz}})
	}
}

type PV struct { p, v int }

// for each axis, find hailstones with same pos and velocity, return a list
func findSamePosVelStones(axis int) []int {
	pvs := make(map[PV][]int, 0) // table same PV => list of stone IDs
	for i, s := range stones {
		pvs[PV{s.p[axis], s.v[axis]}] = append(pvs[PV{s.p[axis], s.v[axis]}], i)
	}
	// keep only the ones with at least 2 stones
	for pv, samepv := range pvs {
		if len(samepv) < 2 {
			delete(pvs, pv)	
		} else {
			return samepv
		}
	}
	return nil
}

func solve(sa int, ha, hb Stone) (rock Stone) {
	// 2 hailstones HA and HB have the same p&v on axis sa:
	psa := ha.p[sa]
	vsa := ha.v[sa]
	// so the rock R has also the same p&v
	rock.p[sa] = psa
	rock.v[sa] = vsa
	// find two stones H0 and H1 with different p&v@sa
	h0id := stoneNot(0, sa, psa, vsa)
	h0 := stones[h0id]
	h1id := stoneNot(h0id+1, sa, psa, vsa)
	h1 := stones[h1id]
	// determine on sa axis time t0 when R intersects H0 (that we note R+H0)
	// psa + vsa * t0 = h0p + h0v * t0
	// t0 = (h0p - psa) / (vsa - h0v)
	t0 := bigIntDivide(h0.p[sa] - psa, vsa - h0.v[sa])
	// same for t1, time of R+H1
	t1 := bigIntDivide(h1.p[sa] - psa, vsa - h1.v[sa])
	var t0divt1 big.Int
	t0divt1.Quo(t0, t1)
	// now for each axis, we can compute the positions of R+H0 and R+H1
	// and deduce the R positions by factorizing its velocities
	// We use big Ints as a position numbers in the input are already ~ 2^50
	// multiplying two of them overflows int64
	for a := range axisname {
		// We compute the rock pos on axis a. E.g. for the X axis:
		// r.px + t0 * r.vx = h0px + t0 * h0vx
		// r.px + t1 * r.vx = h1px + t1 * h1vx
		// ... some math magic to remove r.vx ...
		// r.px = (t0*t1*h1vx - t0*t1*h0vx + t0*h1px - t1*h0px) / (t0 - t1)
		//              n1    -      n2    +    n3   -    n4    /     n5
		//                    n6           +         n7         /
		//                                 n8                   /     n5
		// r.px =                                               n
		var n, n1, n2, n3, n4, n5, n6, n7, n8, n0 big.Int
		n0.Mul(t1, big.NewInt(int64(h1.v[a])))
		n1.Mul(t0, &n0)
		n0.Mul(t1, big.NewInt(int64(h0.v[a])))
		n2.Mul(t0, &n0)
		n3.Mul(t0, big.NewInt(int64(h1.p[a])))
		n4.Mul(t1, big.NewInt(int64(h0.p[a])))
		n5.Sub(t0, t1)
		n6.Sub(&n1, &n2)
		n7.Sub(&n3, &n4)
		n8.Add(&n6, &n7)
		n.Quo(&n8, &n5)
		pos := int(n.Int64())	// initial position of rock on axis sa
		if a == sa { // on original axis, sa, verify we find the same thing
			if pos != psa {
				panic(fmt.Sprintf("On axis %s, computed rock position as %d instead of %d. Rock being now: %v", axisname[a], pos, psa, rock))
			}
		}
		rock.p[a] = pos
	}
	return
}

// returns id of first stone with different p&v on axis
func stoneNot(id, sa, psa, vsa int) int {
	for stones[id].p[sa] == psa || stones[id].v[sa] == vsa {
		id++
	}
	return id
}

// just performs i/j returned as a bigInt, but panics if i is not a multiple of j
// for early detection of bugs
func bigIntDivide(i, j int) *big.Int {
	if i % j != 0 {
		panic(fmt.Sprintf("intDivide; %d not a multiple of %d", i, j))
	}
	return big.NewInt(int64(i/j))
}

//////////// Common Parts code

//////////// PrettyPrinting & Debugging functions
