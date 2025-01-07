// Adventofcode 2018, d10, in go. https://adventofcode.com/2018/day/10
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example
// TEST: example
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"fmt"
	"regexp"
	//"github.com/eiannone/keyboard" // to wait for a key press
	"os"
	// "flag"
	// "slices"
)

//////////// Options parsing & exec parts

func main() {
	ExecOptions(2, NoXtraOpts, part1, part2)
}

//////////// Part 1

func part1(lines []string) (res int) {
	part(lines)
	os.Exit(0)
	return
}

//////////// Part 2

func part2(lines []string) (res int) {
	return part(lines)
}
	

//////////// Common Parts code

func part(lines []string) (res int) {
	//if err := keyboard.Open(); err != nil { panic(err) } // eiannone/keyboard
    //defer keyboard.Close()
	ops, vs := parse(lines)
	ps := make([]Point, len(ops), len(ops))
	hmin := MaxInt
	tclose := CloseTime(ops, vs)	// heuristics: time close to tmin
	VPf("Lookin at times %d to %d\n", tclose, tclose + tclose / 500)
	var tmin int
	
	for t := tclose; t < tclose + 20; t++ {
		for i, p := range ops {
			ps[i] = Point{p.x + t * vs[i].x, p.y + t * vs[i].y}
		}
		_,_,_, h := BoundingBox(ps)
		if h < hmin {
			hmin = h
			tmin = t
		}
	}
	for i, p := range ops {
		ps[i] = Point{p.x + tmin * vs[i].x, p.y + tmin * vs[i].y}
	}
	x, y, w, h := BoundingBox(ps)
	PrintPic(tmin, x, y, w, h, ps)

	//fmt.Println("Press any key to continue...")
	//	char, _, _ := keyboard.GetKey() // eiannone/keyboard Wait for a key press
	//if char == 'q' { os.Exit(0) }
	return tmin
}

func PrintPic(t, xo, yo, w, h int, ps []Point) {
	b := MakeBoard[bool](w,h)
	for _, p := range ps {
		b.a[p.x-xo][p.y-yo] = true
	}
	fmt.Printf("Pic @%d from [%d %d], %d x %d\n", t, xo, yo, w, h)
	for y := range b.h {
		for x  := range b.w {
			if b.a[x][y] {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Print("\n")
	}
}

// find the probable time t where lights meet
// take 1 & 2 points, we look for t that has p1 + t * v1 =~ p2 + t * v2
// t * (v1 - v2) = p2 - p1;  t = (p2 - p1) / (v1 - v2)
// we look at the y coord only, as texts heights are less than their width
// we look for higest opposite velocities (5 , -5) for more precision
// We thus only look at the 20 times around this probable time

func CloseTime(ps, vs []Point) (t int) {
	var i1, i2 int
	for i, v := range vs {
		switch v.y {
		case 5: i1 = i
		case -5: i2 = i
		}
		if i1 != 0 && i2 != 0 {
			break
		}
	}
	t = (ps[i2].y - ps[i1].y) / (vs[i1].y - vs[i2].y)
	return t - 10				// start of a 20 seconds interval around close t
}

func parse(lines []string) (ps, vs []Point) {
	renum := regexp.MustCompile("-?[[:digit:]]+") // example code body, replace.
	for _, line := range lines {
		ns := atoil(renum.FindAllString(line, -1))
		ps = append(ps, Point{ns[0], ns[1]})
		vs = append(vs, Point{ns[2], ns[3]})
	}
	return
}

//////////// PrettyPrinting & Debugging functions. See also the VPx functions.

func DEBUG() {
	if ! debug { return }
	fmt.Println("DEBUG!")
	// ad hoc debug code here
}
