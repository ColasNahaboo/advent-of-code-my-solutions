// Adventofcode 2023, d18, in go. https://adventofcode.com/2023/day/18
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 62
// TEST: -1 example0 9
// TEST: -1 example1 303
// TEST: example 952408144115
// TEST: example0 9
// TEST: example1 771
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"flag"
	"fmt"
	"regexp"
	"strings"
	"sort"
)

type Dig struct {
	dir int 					// direction 0 1 2 3 = U R D L (or N E S W)
	len int						// duration of the dig
	turn int					// turn from previous dig: -1 left, +1 right
	x, y int					// starting point
	x2, y2 int					// ending point, excluded (only used in part2)
}

var plan = []Dig{}				// the dig plan, set of instructions
var area Scalarray[bool]		// the dig aread map, true = hole "#"
var dirx = [4]int{0, 1, 0, -1}		// X-component of directions
var diry = [4]int{-1, 0, 1, 0}		// Y-component of directions

var verbose bool

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
	rotation := parse(lines)
	VPScallarrayBool("Trenches", area)
	digLagoon(rotation)
	VPScallarrayBool("Lagoon", area)
	return dug(area)
}

// return the rotation, number of quarter turns
func parse(lines []string) (rotation int) {
	re := regexp.MustCompile("([URDL]) ([0-9]+) [(]#([0-9a-f]{6})[)]")
	var odir, x, y, turn, maxx, minx, maxy, miny int
	for lineno, line := range lines {
		m := re.FindStringSubmatch(line); if  m == nil {
			panic(fmt.Sprintf("Syntax error line %d: \"%s\"\n", lineno+1, line))
		}
		dir := strings.Index("URDL", m[1])
		len := atoi(m[2])
		if lineno != 0 {
			if dir == odir {
				panic(fmt.Sprintf("Aligned dig [%d]: \"%s\"\n", lineno+1, line))
			} else if (dir + odir) % 2 == 0 {
				panic(fmt.Sprintf("U-turn dig [%d]: \"%s\"\n", lineno+1, line))
			}
			if (odir + 1) % 4 == dir {
				turn = 1
			} else {
				turn = -1
			}
		}
		plan = append(plan, Dig{dir: dir, len: len, turn: turn, x: x, y: y})
		odir = dir
		x += len * dirx[dir]
		y += len * diry[dir]
		if x > maxx { maxx = x }
		if x < minx { minx = x }
		if y > maxy { maxy = y }
		if y < miny { miny = y }
		rotation += turn
	}
	// now create the area map, leaving a border around it
	VPf("dig area %d x %d, [%d, %d] x [%d, %d]\n", maxx - minx +  1, maxy - miny + 1, minx, maxx, miny, maxy)
	area = makeScalarrayB[bool](maxx - minx +  1, maxy - miny + 1, 1)
	// we translate the dig coordinates to fit inside the area and its border
	for i := range plan {
		plan[i].x -= minx
		plan[i].y -= miny
	}
	// dig the trenches
	for _, dig := range plan {
		x, y = dig.x, dig.y
		p := area.PosB(x, y)
		for d := 0; d < dig.len; d++ {
			area.a[p] = true
			p += dirx[dig.dir] + diry[dig.dir]*area.w
		}
	}
	if rotation == 0 { panic("No rotation detected!")}
	return
}

// We suppose we can find a point "inpos" in the interior by looking on the
// "inside side" of the first dig, as given by the total rotation
func digLagoon(rotation int) {
	inpos := area.PosB(plan[0].x, plan[0].y) // start of digs
	inpos += dirx[plan[0].dir] + diry[plan[0].dir]*area.w // move 1 step along trench
	turndir := 1				// turn one step inside, to its right
	if rotation < 0 { turndir = -1}	// oops, no to its left
	turndir = (plan[0].dir + turndir) % 4 // actual step dir
	inpos += dirx[turndir] + diry[turndir]*area.w
	if area.a[inpos] {
		panic("Cannot find soil inside!")
	}
	// fill from inpos
	digLagoonInside(inpos)
}

func digLagoonInside(p int) {
	if area.a[p] {
		return
	}
	area.a[p] = true
	for _, step := range area.Dirs() {
		digLagoonInside(p + step)
	}
}

func dug(area Scalarray[bool]) (sum int) {
	for _, isdug := range area.a {
		if isdug {
			sum++
		}
	}
	return
}		

//////////// Part 2

// we manage tiles, rectangles instead of individual points
// the list of tile dims and position to the right of them (excluded)
// we create also special tiles of width 1 to hold the trenches, thus around
// the "wide" free soil tiles
//                                                                             
//     original map 23x8             tiles map 5x5                            
// xlist=0,7,21   ylist=0,3,7                                                  
//                                                                             
// 0111111233333333333334                                                      
//                                          all '#' are single-cell tiles      
// ########                0      ###                                          
// #      #                1      #1#       1 = 6x2                            
// #      #                1      #-###                                        
// #      ###############  2      #2|3#     2 = 6x3                            
// #                    #  3      #####                                        
// #                    #  3                3 = 13x3                           
// #                    #  3                                                   
// ######################  4                                                   

var wtiles, htiles, xtiles, ytiles, xlist, ylist []int
var colordir = [4]int{1, 2 ,3, 0}	// 0 means R, 1 means D, 2 means L, 3 means U

func part2(lines []string) int {
	rotation := parse2(lines)
	digLagoon(rotation)
	return dugTiles()
}

func parse2(lines []string) (rotation int) {
	re := regexp.MustCompile("([URDL]) ([0-9]+) [(]#([0-9a-f]{6})[)]")
	var l, odir, x, y, x2, y2, turn, maxx, minx, maxy, miny int
	xlist = []int{}				// list of x values found
	ylist = []int{}				// list of y values found
	for lineno, line := range lines {
		m := re.FindStringSubmatch(line); if  m == nil {
			panic(fmt.Sprintf("Syntax error line %d: \"%s\"\n", lineno+1, line))
		}
		hexlen := m[3][0:5]
		fmt.Sscanf(hexlen, "%x", &l)
		dir := colordir[int(m[3][5] - '0')]
		if lineno != 0 {
			if dir == odir {
				panic(fmt.Sprintf("Aligned dig [%d]: \"%s\"\n", lineno+1, line))
			} else if (dir + odir) % 2 == 0 {
				panic(fmt.Sprintf("U-turn dig [%d]: \"%s\"\n", lineno+1, line))
			}
			if (odir + 1) % 4 == dir {
				turn = 1
			} else {
				turn = -1
			}
			rotation += turn
		}
		plan = append(plan, Dig{dir: dir, len: l, turn: turn, x: x, y: y})
		odir = dir
		x += l * dirx[dir]
		y += l * diry[dir]
		if x > maxx { maxx = x }
		if x < minx { minx = x }
		if y > maxy { maxy = y }
		if y < miny { miny = y }
		if indexOfInt(xlist, x) == -1 {
			xlist = append(xlist, x)
		}
		if indexOfInt(ylist, y) == -1 {
			ylist = append(ylist, y)
		}
	}
	// find the tiles dimensions and positions
	sort.Slice(xlist, func(i, j int) bool { return xlist[i] < xlist[j] })
	sort.Slice(ylist, func(i, j int) bool { return ylist[i] < ylist[j] })
	// tiles on X axis = len(xlist) 1-wide tiles + (len(xlist)-1) inbetween
	xtiles = make([]int, 2*len(xlist)-1, 2*len(xlist)-1)
	wtiles = make([]int, 2*len(xlist)-1, 2*len(xlist)-1)
	for t,i := 0,0; i < len(xlist); i++ {
		wtiles[t] = 1			// at t, the 1-wide tile
		xtiles[t] = xlist[i]
		t++
		if t >= len(wtiles) { break }
		wtiles[t] = xlist[i+1] - xlist[i] - 1
		xtiles[t] = xlist[i] + 1
		t++
	}
	ytiles = make([]int, 2*len(ylist)-1, 2*len(ylist)-1)
	htiles = make([]int, 2*len(ylist)-1, 2*len(ylist)-1)
	for t,i := 0,0; i < len(ylist); i++ {
		htiles[t] = 1			// at t, the 1-wide tile
		ytiles[t] = ylist[i]
		t++
		if t >= len(htiles) { break }
		htiles[t] = ylist[i+1] - ylist[i] - 1
		ytiles[t] = ylist[i] + 1
		t++
	}

	// now create the area map of the tiles
	// a tile (x,y) has upper left corner at (xtile[x], ytile[y]) and
	// has dimensions wtile[x] x wtile[y]
	area = makeScalarray[bool](len(wtiles), len(htiles))
	// we translate the dig coordinates to virtual ones in the map of tiles
	for i := range plan {
		x = plan[i].x
		y = plan[i].y
		x2 = x + dirx[plan[i].dir] * plan[i].len
		y2 = y + diry[plan[i].dir] * plan[i].len
		plan[i].x = tileOf(xtiles, wtiles, x)
		plan[i].y = tileOf(ytiles, htiles, y)
		plan[i].x2 = tileOf(xtiles, wtiles, x2)
		plan[i].y2 = tileOf(ytiles, htiles, y2)
	}
	// dig the trenches tile by tile, from x,y to x2,y2
	for _, dig := range plan {
		for x, y = dig.x, dig.y; x != dig.x2 || y != dig.y2; x, y = x+dirx[dig.dir], y+diry[dig.dir] {
			area.a[area.Pos(x, y)] = true
		}
	}
	return
}

// return the index of tile we are in
func tileOf(xt, wt []int, x int) int {
	for i := range xt {
		if x >= xt[i] && x < xt[i] + wt[i] {
			return i
		}
	}
	panic(fmt.Sprintf("tileOf: %d not found\n", x))
}

// sum the variously sized tiles
// a tile (x,y) has surface wtile[x]*wtile[y]
func dugTiles() (sum int) {
	for p, isdug := range area.a {
		if isdug {
			x, y := area.Coords(p)
			tilesurface := wtiles[x] * htiles[y]
			sum += tilesurface
		}
	}
	return
}		


//////////// Common Parts code

//////////// PrettyPrinting & Debugging functions
