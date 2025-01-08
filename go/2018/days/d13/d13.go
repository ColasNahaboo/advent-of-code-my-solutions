// Adventofcode 2018, d13, in go. https://adventofcode.com/2018/day/13
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example "7_3"
// TEST: example
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"fmt"
	"strings"
	// "flag"
	"sort"
)

const (							// Track symbols on the board
	TN = 0						// nothing
	TV = 1						// |
	TH = 2						// -
	TI = 3						// + Intersection = bitwise-or of: TV | TH
	TS = 4						// /
	TB = 8						// \
	TT = 15						// the track itself, masking cart
	TC = 16						// Cart present! (bitwise-ored)
)

type Cart struct {
	dead bool					// has it been destroyed?
	p Point 					// position
	d int						// direction
	turn int					// 0=left, 1=straight, 2=right and loops
	o int						// order: scalar coordinate: p.x + p.y * b.w
}

//////////// Options parsing & exec parts

func main() {
	ExecOptionsString(2, NoXtraOpts, part1, part2)
}

//////////// Part 1

func part1(lines []string) string {
	b, carts := parse(lines)
	VPf("grid %d x %d with %d carts\n", b.w, b.h, len(carts))
	VPtracks(b,carts)
	for {
		if p, crashed := DoTick(b, carts, false); crashed {
			return itoa(p.x) + "_" + itoa(p.y)
		}
		carts = UpdateCarts(b, carts)
		VPtracks(b,carts)
	}
}

func DoTick(b Board[byte], carts []Cart, cont bool) (crash Point, crashed bool) {
	for ci, cart := range carts {
		if cart.dead { continue }
		p := cart.p.StepOrtho(cart.d)
		track := b.a[p.x][p.y]
		if track & TC != 0 {			// there is a cart, crash!
			if cont {					// part2: mark carts dead & continue
				carts[ci].dead = true
				c := CartIndex(carts, p)
				carts[c].dead = true
				b.a[cart.p.x][cart.p.y] &= TT // remove carts markers
				b.a[p.x][p.y] &= TT
				crashed = true
				continue
			} else {
				return p, true
			}
		}
		// now we are sure there is no cart, no need to mask by TT
		b.a[cart.p.x][cart.p.y] &= TT		// remove cart marker a prev pos
		b.a[p.x][p.y] |= TC					// put it at new
		carts[ci].p = p						// move cart there
		if track == TI {
			switch cart.turn {
			case 0: carts[ci].d = RotateDirOrtho(cart.d, -1)
			case 2: carts[ci].d = RotateDirOrtho(cart.d, 1)
			}
			carts[ci].turn = (carts[ci].turn + 1) % 3
		} else if track == TS {
			switch cart.d {
			case DirsOrthoN: carts[ci].d = DirsOrthoE
			case DirsOrthoE: carts[ci].d = DirsOrthoN
			case DirsOrthoS: carts[ci].d = DirsOrthoW
			case DirsOrthoW: carts[ci].d = DirsOrthoS
			}
		} else if track == TB {
			switch cart.d {
			case DirsOrthoN: carts[ci].d = DirsOrthoW
			case DirsOrthoE: carts[ci].d = DirsOrthoS
			case DirsOrthoS: carts[ci].d = DirsOrthoE
			case DirsOrthoW: carts[ci].d = DirsOrthoN
			}
		} else if track == TN {
			panic("No track!")
		}
	}
	return
}

func SortCarts(carts []Cart) {
	sort.Slice(carts, func(i, j int) bool { return carts[i].o < carts[j].o})
}

func UpdateCarts(b Board[byte], carts []Cart) []Cart {
	for i := range carts {
		carts[i].o = carts[i].p.x + carts[i].p.y * b.w
	}
	sort.Slice(carts, func(i, j int) bool { return carts[i].o < carts[j].o})
	return carts
}

//////////// Part 2

func part2(lines []string) string {
	b, carts := parse(lines)
	VPf("grid %d x %d with %d carts\n", b.w, b.h, len(carts))
	VPtracks(b,carts)
	for {
		if _, crashed := DoTick(b, carts, true); crashed {
			carts = CleanCarts(carts)
			if len(carts) == 1 {
				return itoa(carts[0].p.x) + "_" + itoa(carts[0].p.y)
			} else if len(carts) < 1 {
				panic("No more carts!")
			}
			VPf("Crash, %d carts remaining\n", len(carts))
		}
		carts = UpdateCarts(b, carts)
		VPtracks(b,carts)
	}
}

func CleanCarts(carts []Cart) (n []Cart) {
	for _, cart := range carts {
		if ! cart.dead {
			n = append(n, cart)
		}
	}
	return
}

func CartIndex(carts []Cart, p Point) int {
	for i, cart := range carts {
		if cart.p == p {
			return i
		}
	}
	return -1
}
	
//////////// Common Parts code

func parse(lines []string) (b Board[byte], carts []Cart) {
	b = MakeBoard[byte](len(lines[0]), len(lines))
	for y, line := range lines {
		for x, r := range line {
			switch r {
			case '|','^','v': b.a[x][y] = TV
			case '-','>','<': b.a[x][y] = TH
			case '+': b.a[x][y] = TI
			case '/': b.a[x][y] = TS
			case '\\': b.a[x][y] = TB
			case ' ':
			default: panicf("Bad char: %v", r)
			}
			if strings.ContainsRune(CartDirGlyphs, r) { // cart
				cart := Cart{p: Point{x,y}}
				switch r {
				case '^': cart.d = DirsOrthoN
				case '>': cart.d = DirsOrthoE
				case 'v': cart.d = DirsOrthoS
				case '<': cart.d = DirsOrthoW
				}
				cart.o = x + y * b.w
				carts = append(carts, cart)
				b.a[x][y] |= TC
			}
		}
	}
	SortCarts(carts)
	return
}

//////////// PrettyPrinting & Debugging functions. See also the VPx functions.

func DEBUG() {
	if ! debug { return }
	fmt.Println("DEBUG!")
	// ad hoc debug code here
}

func VPtracks(b Board[byte], carts []Cart) {
	if ! verbose {return}
	for y := range b.h {
		for x  := range b.w {
			s := " "
			switch b.a[x][y] {
			case TV: s = "|"
			case TH: s = "-"
			case TI: s = "+"
			case TS: s = "/"
			case TB: s = "\\"
			}
			if b.a[x][y] & TC != 0 {
				s = CartDirGlyph(carts, x, y)
			}
			fmt.Print(s)
		}
		fmt.Print("\n")
	}
}

var CartDirGlyphs = "^>v<"
func CartDirGlyph(carts []Cart, x, y int) string {
	for _, cart := range carts {
		if cart.p.x == x && cart.p.y == y {
			return CartDirGlyphs[cart.d:cart.d+1]
		}
	}
	return "?"
}

