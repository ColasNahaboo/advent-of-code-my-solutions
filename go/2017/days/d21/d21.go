// Adventofcode 2017, d21, in go. https://adventofcode.com/2017/day/21
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
	"strings"
	"os"
	// "golang.org/x/exp/slices"
)

var verbose, debug bool

type Rule struct {
	size int					// applies to squares of this size
	lit int						// with lit "#" pixels
	conds []string				// all variants of pattern match by rot or flip
	action string				// replacement pattern
	alit int					// lit pixels in action
}

type Image struct {
	size int
	lit int
	data string					// in the raw input format: e.g: "..#/.#./..#"
}

type Ruleset [][][]Rule			// [size][lit] -> list of rules

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	debugFlag := flag.Bool("d", false, "debug: even more verbose")
	flag.Parse()
	verbose = *verboseFlag
	debug = *debugFlag
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
	} else {
		VP("Running Part2")
		fmt.Println(part2(lines))
	}
}

//////////// Part 1

func part1(lines []string) int {
	return iterate(5, lines)
}

//////////// Part 2
func part2(lines []string) int {
	return iterate(18, lines)
}

//////////// Common Parts code

func iterate(n int, lines []string) int {
	rs := parse(lines)
	DEBUG()
	image := MakeImage(".#./..#/###")
	VPimage("Init", image)
	for iteration := 1; iteration <= n; iteration++ {
		images := []*Image{}
		for _, simage := range image.Split() {
			images = append(images, rs.Apply(simage))
		}
		image.Assemble(images)
		VPimage("Iteration #" + itoa(iteration), image)
	}
	return image.lit
}

func parse(lines []string) Ruleset {
	re := regexp.MustCompile("^([.#/]+) => ([.#/]+)")
	rs := make(Ruleset, 4, 4)
	for i := 2; i <= 3; i++ {
		rs[i] = make([][]Rule, 10, 10) // lit is max 9, as max size in input is 3
	}
	for _, line := range lines {
		m := MustFindStringSubmatch(re, line)
		rule := Rule{size: PatternSize(m[1]), lit: PatternLit(m[1]), conds: PatternVariants(m[1]), action: m[2], alit: PatternLit(m[2])}
		rs[rule.size][rule.lit] = append(rs[rule.size][rule.lit], rule)
	}
	return rs
}

func MakeImage(data string) (i Image) {
	i.data = data
	i.size = strings.Index(data, "/")
	i.lit = PatternLit(data)
	return
}

func PatternSize(s string) int {
	if len(s) < 9 {
		return 2
	} else {
		return 3
	}
}

func PatternLit(s string) int {
	return strings.Count(s, "#")
}

func PatternVariants(s string) ([]string) {
	if len(s) < 9 {
		return PatternVariants2(s)
	} else {
		return PatternVariants3(s)
	}
}

func (i *Image) Split() (l []*Image) {
	tile := 3
	if i.size % 2 == 0 {
		tile = 2
	}
	ni := i.size / tile
	for y := 0; y < ni; y++ {
		for x := 0; x < ni; x++ {
			lit, data := i.Extract(tile, x * tile, y * tile)
			image :=  &Image{tile, lit, string(data)}
			l = append(l, image)
		}
	}
	return
}

// extract sub-image of side size starting at x,y
func (i *Image) Extract(size, x, y int) (lit int, data []byte) {
	data = make([]byte, (size+1)*size-1, (size+1)*size-1)
	for yi := 0; yi < size; yi++ {
		for xi := 0; xi < size; xi++ {
			b := i.data[IPos(i.size, x + xi, y + yi)]
			if b == '#' {
				lit++
			}
			data[IPos(size, xi, yi)] = b
		}
	}
	SetSeparator(size, data)
	return
}

func (i *Image) Assemble(l []*Image) {
	ni := IntSquareRoot(len(l))
	isize := l[0].size
	i.size = ni * isize
	i.lit = 0
	data := make([]byte, (i.size+1)*i.size-1, (i.size+1)*i.size-1)
	for y := 0; y < ni; y++ {
		for x := 0; x < ni; x++ {
			ii := l[x + y*ni]
			i.lit += ii.lit
			for yi := 0; yi < isize; yi++ {
				for xi := 0; xi < isize; xi++ {
					data[IPos(i.size, x*isize + xi, y*isize + yi)] = ii.data[IPos(isize, xi, yi)]
				}
			}
		}
	}
	SetSeparator(i.size, data)
	i.data = string(data)
}

// add /-separators
func SetSeparator(size int, data []byte) {
	for y := 0; y < size - 1; y++ {
		data[IPos(size, size, y)] = '/'
	}
}	

// index in image.data string of the coordinates [x, y] in 2D image
func IPos(size, x, y int) int {
	return x + y * (size + 1)	// take the trailing / into account
}

func IntSquareRoot(i int) (r int) {
	for r = 1; r*r < i; r++ {
	}
	if r*r == i {
		return r
	}
	panic("Not a square: "+ itoa(i))
}

func (rs Ruleset) Apply(i *Image) (ni *Image) {
	for _, rule := range rs[i.size][i.lit] {
		ni = rule.Apply(i)
		if ni != nil {
			return ni
		}
	}
	panic("No rule applies to: \"" + i.data + "\"")
}

func (r *Rule) Apply(i *Image) (ni *Image) {
	for cn, cond := range r.conds {
		if i.data == cond {
			VPf("  Match: %s [%d] %s | %s ==> %s\n", i.data, cn, cond, r.conds[0], r.action)
			return &Image{size: i.size+1, lit: r.alit, data: r.action}
		}
	}
	return nil
}

// a right-rotation transforms AB into CA
//                             CD      DB
// e.g: AB/CD into CA/DB. The index transform is thus 01234 => 14203
// as in char at position 0 moves to position 1, at pos 1 to 4, etc...
// we give thus 1, 4, 2, 0, 3 as args to PatternTransform
//                   r1 r1
//    r1 r2 r3 fh fv fh fv
// 01 20 32 13 10 23 02 31
// 23 31 10 02 32 01 13 20


func PatternVariants2(s string) []string {
	// rotations. / stays at 2
	r1 := PatternTransform(s,  1, 4, 2, 0, 3)
	r2 := PatternTransform(r1, 1, 4, 2, 0, 3)
	r3 := PatternTransform(r2, 1, 4, 2, 0, 3)
	// flips
	fh := PatternTransform(s, 1, 0, 2, 4, 3)
	fv := PatternTransform(s, 3, 4, 2, 0, 1)
	// rotation 90 + flips
	r1fh := PatternTransform(s, 0, 3, 2, 1, 4)
	r1fv := PatternTransform(s, 4, 1, 2, 3, 0)
	return []string{s, r1, r2, r3, fh, fv, r1fh, r1fv}
}

//                   r1  r1 
//      r1   fh  fv  fh  fv
// 012/ 840 210 89a 048 a62
// 456/ 951 654 456 159 951
// 89a  a62 a98 012 26a 840

func PatternVariants3(s string) []string {
	// rotations. / stays at 3 and 7
	r1 := PatternTransform(s,  2, 6,10, 3, 1, 5, 9, 7, 0, 4, 8) 
	r2 := PatternTransform(r1, 2, 6,10, 3, 1, 5, 9, 7, 0, 4, 8)
	r3 := PatternTransform(r2, 2, 6,10, 3, 1, 5, 9, 7, 0, 4, 8)
	// flips
	fh := PatternTransform(s, 2, 1, 0, 3, 6, 5, 4, 7, 10, 9, 8)
	fv := PatternTransform(s, 8, 9, 10, 3, 4, 5, 6, 7, 0, 1, 2)
	// rotation 90 + flips
	r1fh := PatternTransform(r1, 2, 1, 0, 3, 6, 5, 4, 7, 10, 9, 8)
	r1fv := PatternTransform(r1, 8, 9, 10, 3, 4, 5, 6, 7, 0, 1, 2)
	return []string{s, r1, r2, r3, fh, fv, r1fh, r1fv}
}

func PatternTransform(s string, pos ...int) string {
	t := make([]rune, len(s), len(s))
	for i, b := range s {
		t[pos[i]] = b
	}
	return string(t)
}		

func MustFindStringSubmatch(re *regexp.Regexp, s string) (m []string) {
	m = re.FindStringSubmatch(s)
	if m == nil {
		panic("No match of \"" + re.String() + "\" on: \"" + s + "\"")
	}
	return
}

//////////// PrettyPrinting & Debugging functions

func VPimage(label string, i Image) {
	if ! verbose {
		return
	}
	fmt.Printf("%s: Image size=%d, lit=%d\n", label, i.size, i.lit)
	for _, b := range i.data {
		if b == '/' {
			fmt.Println()
		} else {
			fmt.Print(string(b))
		}
	}
	fmt.Println()
}

func VPruleset(rs Ruleset) {
	if ! verbose {
		return
	}
	fmt.Printf("Ruleset:\n")
	for size := 2; size <= 3; size++ {
		fmt.Printf("  Size %d\n", size)
		for lit, rules := range rs[size] {
			if len(rules) > 0 {
				fmt.Printf("    lit %d:\n", lit)
				for _, rule := range rules {
					for _, cond := range rule.conds {
						fmt.Printf("      %s\n", cond)
					}
					fmt.Printf("      ==> %s\n", rule.action)
				}
			}
		}
	}
}

func DEBUG() {
	if ! debug { return }
	for i, s := range PatternVariants3("###/..#/#..") {
		fmt.Printf("  [%d] %s\n", i, s)
	}
	os.Exit(0)
}
