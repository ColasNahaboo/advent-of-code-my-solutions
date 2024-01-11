// Adventofcode 2023, d22, in go. https://adventofcode.com/2023/day/22
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 5
// TEST: example 7
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"flag"
	"fmt"
	"regexp"
	"sort"
)

var verbose bool

type Voxel struct {
	x, y, z int
}
type Brick struct {
	id int
	x1, y1, z1, x2, y2, z2 int	// in order: n1 <= n2
	supports []int				// list of brick IDs supporting it
}
var bricks []Brick				// list of bricks, by order of parsing
var roofs [][]int				// for each z, list of brick IDs of z2 == z
var destroyed = -1				// for part2 the brick ID we destroy. -1 in part1

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
	parse(lines)
	VPbricks()
	VProofs()
	freeFall()
	VProofs()
	findSupports()
	VPbricks()
	return len(bricks) - len(monoSupports())
}


//////////// Part 2
func part2(lines []string) (fell int) {
	// first, let the bricks fall & rest and archive in allbricks
	parse(lines)
	freeFall()
	findSupports()
	allbricks := make([]Brick, len(bricks), len(bricks))
	copy(allbricks, bricks)
	// then, for each fall-causing-if-desintegrated bricks, delete and let fall
	for b := range monoSupports() {
		destroyed = b
		VPf(" Destroying brick %d\n", b)
		copy(bricks, allbricks)
		// top brick does not support any other, so roofs keep the same length
		makeRoofs(len(roofs))
		fell += freeFall()		// incr by number of bricks that fell this round
	}
	return
}

//////////// Common Parts code

func parse(lines []string) {
	re := regexp.MustCompile("([0-9]+),([0-9]+),([0-9]+)~([0-9]+),([0-9]+),([0-9]+)")
	maxroof := 0
	for lineno, line := range lines {
		m := re.FindStringSubmatch(line)
		if m == nil {
			panic(fmt.Sprintf("Syntax error line %d: \"%s\"", lineno, line))
		}
		x1 := atoi(m[1])
		y1 := atoi(m[2])
		z1 := atoi(m[3])
		x2 := atoi(m[4])
		y2 := atoi(m[5])
		z2 := atoi(m[6])
		if z2 > maxroof {
			maxroof = z2
		}
		if x2 < x1 || y2 < y1 || z2 < z1 {
			panic(fmt.Sprintf("Wrong coords order in line %d: \"%s\"", lineno, line))
		}
		bricks = append(bricks, Brick{id:len(bricks), x1:x1, y1:y1, z1:z1, x2:x2, y2:y2, z2:z2})
	}
	makeRoofs(maxroof)
}

func makeRoofs(l int) {
	roofs = make([][]int, l+1, l+1)
	for _, b := range bricks {
		if b.id == destroyed {
			continue
		}
		roofs[b.z2] = append(roofs[b.z2], b.id)
	}
}	

func moveDown(bid, dz int) {
	oz2 := bricks[bid].z2
	bricks[bid].z1 -= dz
	bricks[bid].z2 -= dz
	z2 := bricks[bid].z2
	var ok bool
	roofs[oz2], ok = deleteFastInt(roofs[oz2], bid) // move roof level
	if ! ok {
		panic(fmt.Sprintf("Brick %d not found its its roof %d", bid, oz2))
	}
	roofs[z2] = append(roofs[z2], bid)
	roofstop := len(roofs)-1	// shorten roofs if needed
	if oz2 == roofstop {
		for len(roofs[roofstop]) == 0 {
			roofs = roofs[0:roofstop]
			roofstop--
		}
	}
	VPf("  moveDown by %d: %v\n", dz, bricks[bid])
}

func disjoints(b1, b2 Brick) bool {
	return b1.x2 < b2.x1 || b1.x1 > b2.x2 || b1.y2 < b2.y1 || b1.y1 > b2.y2
}

// enumeration(n) => []int{0, 1, 2, 3, ... n-1}
func enumeration(n int) (e []int) {
	for i := 0; i < n; i++ {
		e = append(e, i)
	}
	return
}

// let all the bricks fall freely, return how many fell
func freeFall() (fell int) {
	// fall down all bricks, starting by the bottom
	bottomup := enumeration(len(bricks))
	sort.Slice(bottomup, func(i, j int) bool { return bricks[bottomup[i]].z1 < bricks[bottomup[j]].z1})
BRICK:
	for _, i := range bottomup {
		if i == destroyed {
			continue
		}
		b := bricks[i]
		for roof := b.z1 - 1; roof > 0; roof-- {
			for _, under := range roofs[roof] {
				if under == destroyed {
					continue
				}
				if disjoints(bricks[under], b) {
					continue
				}
				if b.z1 - roof - 1 > 0 {
					VPf("  move %v on top of %v at roof %d\n", b, bricks[under], roof)
					moveDown(i, b.z1 - roof - 1)
					fell++
				}
				continue BRICK	// we found a support on this "roof" level
			}
		}
		if b.z1 - 1 > 0 {
			VPf("  drop %v by %d on ground\n", b, b.z1 -1)
			moveDown(i, b.z1 - 1)		// none was under it, fall to ground (roof 0)
			fell++
		}
	}
	return
}

// fill the supports field of all bricks: the bricks directly supporting them
func findSupports() {
	for roof := 2; roof < len(roofs); roof++ {
		for _, bid := range roofs[roof] {
			roofunder := bricks[bid].z1 - 1 // take brick height into account
			for _, under := range roofs[roofunder] {
				if disjoints(bricks[under], bricks[bid]) {
					continue
				}
				// we found a support for bid
				bricks[bid].supports = append(bricks[bid].supports, under)
			}
		}
	}
}

// list of all the bricks that would cause a collapse if desintegrated
func monoSupports() map[int]bool {
	supporters := make(map[int]bool, 0)
	for _, b := range bricks {
		if len(b.supports) == 1 {
			VPf("  brick %d cannot be disintegrated, supports %d\n", b.supports[0], b.id)
			supporters[b.supports[0]] = true
		}
	}
	VPf("  %d mono supporters\n", len(supporters))
	return supporters
}

//////////// PrettyPrinting & Debugging functions

func VPbricks() {
	for _, b := range bricks {
		VPf("  %v\n", b)
	}
}


func VProofs() {
	for i := len(roofs) - 1; i > 0; i-- {
		VPf("  [%d]", i)
		for _, bid := range roofs[i] {
			VPf(" %d", bid)
		}
		VP()
	}
}
