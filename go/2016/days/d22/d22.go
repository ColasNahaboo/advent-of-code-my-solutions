// Adventofcode 2016, d22, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,,.test
// TEST: -1 example 5
// TEST: example
package main

import (
	"flag"
	"fmt"
	"regexp"
	//"time"
)

type Node struct {				// the initial state
	x, y, size, used, avail int
	t int						// taquin type of the node
}
var grid []Node					// linear of the 2d grid: pos = x + y*gw
var gw, gh, garea int			// its dims
var ghole = -1					// hole position in grid

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
	parse(lines[2:])
	
	var result int
	if *partOne {
		VP("Running Part1")
		result = part1()
	} else {
		VP("Running Part2")
		result = part2()
	}
	fmt.Println(result)
}

//////////// Part 1

// Part1 is totally simplistic
func part1() (viable int) {
	for a := 0; a < garea; a++ {
		for b := 0; b < garea; b++ {
			if grid[a].used != 0 &&
				a != b &&
				grid[a].used < grid[b].avail {
				viable++
			}
		}
	}
	return			
}

//////////// Common Parts code

func parse(lines []string) {
	re := regexp.MustCompile("/dev/grid/node-x([[:digit:]]+)-y([[:digit:]]+)[[:space:]]+([[:digit:]]+)T[[:space:]]+([[:digit:]]+)T[[:space:]]+([[:digit:]]+)T[[:space:]]+")
	nodes := []Node{}				// list, not ordered in a 2D grid
	gpos := 0						// current pos in grid
	for lineno := 0; lineno < len(lines); lineno++ {
		line := lines[lineno]
		if m := re.FindStringSubmatch(line); m != nil {
			node := Node{
				x: atoi(m[1]),
				y: atoi(m[2]),
				size: atoi(m[3]),
				used: atoi(m[4]),
				avail: atoi(m[5]),
			}
			if node.used == 0 {
				if ghole != -1 {
					panic("grid has 2 holes")
				}
				ghole = gpos
				node.t = HOLE
				
			} else if node.used > 100 {
				node.t = WALL
			} else {
				node.t = TILE
			}
			if node.x > gw -1 { gw = node.x +1; }
			if node.y > gh -1 { gh = node.y +1; }
			nodes = append(nodes, node)
		} else {
			panic(fmt.Sprintf("Syntax error line %d: %s\n", lineno, line))
		}
		gpos++
	}
	// build the 2D grid
	garea = gw * gh
	grid = make([]Node, garea, garea)
	for _, node := range nodes {
		grid[node.x + node.y * gw] = node
	}
	// our goal is at x=gw-1, y=0
	grid[gw-1].t = GOAL
	// checks: missing nodes
	for i := 0; i < garea; i++ {
		if grid[i].size == 0 {
			panic(fmt.Sprintf("grid missing a node at (%d, %d)\n", i%gw, i/gw))
		}
	}
	if ghole == -1 {
		panic("grid has no hole!")
	}
	VPf("Grid %d x %d\n", gw, gh)
}

//////////// Part 2

// For part 2, we solve this as a taquin, with A*
// We cannot use github.com/fzipp/astar as for instance many states can be the
// destination, not just one: all states where the goal is in pos (0, 0)
// There is one hole, and movable tiles or immovable ealls
// We must bring the goal tile to (0, 0) from its (gw-1, 0) position
// We use a grid augmented by borders made of walls. So, if the input
// is a grid of 31x30, our board will be 33x32, and (0, 0) in the input grid
// will be in our board a (1, 1), thus a position of 34 in the "grid" int array
// and (0, 30) is position 64 in board

type State struct {
	hole int					// position of the hole on board
	goal int					// position of the goal on board
	board string				// the board as a string, of size bw*bh
}

// taquin types (bytes). Constraint: WALL < TILE && GOAL > TILE
const (
	WALL = '#'
	TILE = '.'
	GOAL = 'X'
	HOLE = '@'
	SEP = ' '					// DEBUG: for readability of boards
)

type States []State
var bboard []byte				// used as pre-allocated space to build strings
var bw, bh, barea, bgoal int
var graph States
var stateid map[string]int
var dirs []int					// offsets for up, down, left, right

func part2() int {
	bw, bh = gw + 2, gh + 2
	bw += 1						// DEBUG: space between rows
	barea = bw*bh
	dirs = []int{-bw, bw, -1, +1}
	bboard = make([]byte, barea, barea)
	bhole := grid2boardPos(ghole)
	bgoal = bw + 1				// the final destination of G
	VPf("Grid %d (%d, %d), board %d (%d, %d)\n", gw*gh, gw, gh, barea, bw, bh)
	for i, node := range grid {
		//VPf("  grid pos %d ==> board pos %d\n", i, grid2boardPos(i))
		bboard[grid2boardPos(i)] = byte(node.t)
	}
	// walls
	for i := 0; i < bw; i++ {	// top & bottom border wall
		if i % bw == bw -1 { continue } // DEBUG
		bboard[i] = WALL
		bboard[i + bw * (bh-1)] = WALL
	}
	for j := 1; j < bh-1; j++ {	// left & right border wall
		bboard[j * bw] = WALL
		bboard[j * bw + gw+1] = WALL
	}
	for j := 0; j < bh; j++ {	// DEBUG: space between rows
		bboard[j * bw + bw - 1] = SEP
	}
	// goal pseudo-state, id 0, used only in nodeDist comparison of goal pos
	goal_s := State{goal: bgoal}
	// start state, id 1
	start_s := State{hole: bhole, goal: bw + gw, board: string(bboard)}
	// The node object used in astar is the index of a state in the graph array
	graph = []State{goal_s, start_s}
	//VPf("Init board: [%3d] %v\n", 1, start_s)
	//VPnode(1)

	// Reverse map to find state id from its board string representation
	stateid = make(map[string]int, 0)
	stateid[goal_s.board] = 0
	stateid[start_s.board] = 1
	
	// Find the shortest path
	path := FindPath[int](1, 0, nodesConnected, nodesDistance, nodesHeuristic)
	if verbose {
		VPf("### Path found with %d steps:\n", len(path) - 1)
		for i, n :=range path {
			VPf("Step %d: state [%d]:\n", i, n)
			VPnode(n)
		}
	}
	printBoard(start_s.board)
	copy(bboard, []byte(start_s.board))
	for _, n := range path {
		bboard[graph[n].hole] = byte('+')
		bboard[graph[n].goal] = byte('x')
	}
	printBoard(string(bboard))
	return len(path) - 1
}

// return a state ID, auto-creating and declarating it if needed
func stateId(hole, goal int, board string) (id int) {
	var ok bool
	if id, ok = stateid[board]; ok == false {
		state := State{hole: hole, goal: goal, board: board}
		id = len(graph)
		graph = append(graph, state)
		VPf("  Creating new state of id %d: %v\n", id, state) 
		stateid[board] = id
	}
	return
}

func grid2boardPos(g int) int {
	// gx := g % gw, gy := g / gw
	// by := gx+1, by := gy+1
	// b = by*bw + bx
	return (g/gw+1)*bw + g%gw + 1
}

func board2gridPos(b int) int {
	// bx := b % bw, by := b / bw
	// gx := bx-1, gy := by-1
	// g = gy*gw + gx
	gx := b%bw - 1
	gy := b/bw - 1
	if gx < 0 || gx >= gw || gy < 0 || gy >= gh {
		panic(fmt.Sprintf("board2gridPos: board position %d is outside (%d, %d) of grid!", b, gx, gy))
	}
	return (b/bw+1)*gw + b%bw + 1
}

func VPnode(n int) {
	if ! verbose {
		return
	}
	printBoard(graph[n].board)
}	

func printBoard(board string) {
	for i := 0; i < bh; i++ {
		fmt.Println(board[i*bw:(i+1)*bw])
	}
}
	

// here we implement all the constraints on tile moves
// we use the global value graph, as the parameter is a copy

func nodesConnected(n int) (nexts []int) {
	state := graph[n]
	h := state.hole
	// look if we can swap hole with up, down, left, right adjacent tile
	for _, dir := range dirs {	
		if state.board[h + dir] >= TILE { // we can swap hole and tile
			copy(bboard, state.board)
			bboard[h] = state.board[h + dir]
			bboard[h + dir] = HOLE
			newgoal := state.goal
			if bboard[h] == GOAL { // we have moved the goal
				newgoal = h
			}
			newstate := stateId(h+dir, newgoal, string(bboard))
			//VPf(" |  [%2d] %v\n", n, state.board)
			//VPf(" L_ [%2d] %v\n", newstate, string(bboard))
			nexts = append(nexts, newstate)
		}
	}
	VPf("  Neighbours of state %d: %v\n", n, nexts)
	return
}

// in a path, all neigbours are at a distance of one move (a tile/hole swap)
func nodesDistance(i, j int) float64 {
	return 1
}

// heuristics: the "badness" (cost) of a node is its goal manhattan distance
func nodesHeuristic(n, dest int) float64 {
	if graph[n].goal == bgoal {	// Goal reached!
		return -1
	}
	// bx := b % bw, by := b / bw
	i := graph[n].goal
	return float64(dPos(i, bgoal) + dPos(i, graph[n].hole))
}

// the manhattan distance between 2 positions
func dPos(i, j int) int {
	return intAbs(j%bw - i%bw) + intAbs(j/bw - i/bw)
}
