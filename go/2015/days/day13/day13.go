// Adventofcode 2015, day13, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 input 618
// TEST: input 601

// we reuse the logic and the code of Day09, but with the distance of A to B
// being arc(A,B).weight + arc(B,A).weight. The gain of happiness.
// verbose (print routes) with option -v

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	// "regexp"
)

type Graph struct {
	names   []string
	ids     map[string]int
	arcs    []Arc // only used to compute the weights
	weights [][]int
	min     int
	max     int
}
type Arc struct {
	from   int
	to     int
	weight int
}

var verbose bool

func main() {
	partOne := flag.Bool("1", false, "run part one code, instead of part 2 (default)")
	verboseFlag := flag.Bool("v", false, "verbose: print routes")
	flag.Parse()
	verbose = *verboseFlag
	infile := "input.txt"
	if flag.NArg() > 0 {
		infile = flag.Arg(0)
	}
	lines := FileToLines(infile)

	var result int
	if *partOne {
		fmt.Println("Running Part1")
		result = Part1(lines)
	} else {
		fmt.Println("Running Part2")
		result = Part2(lines)
	}
	fmt.Println(result)
}

//////////// Part 1

func Part1(lines []string) int {
	graph := NewGraph(lines, false)
	return LongestRoute(&graph)
}

//////////// Part 2
func Part2(lines []string) int {
	graph := NewGraph(lines, true)
	return LongestRoute(&graph)
}

//////////// Common Parts code

// if additional parameter "me": add oneself to the graph
func NewGraph(lines []string, me bool) Graph {
	var graph Graph
	graph.names = make([]string, 0)
	graph.ids = make(map[string]int)
	graph.arcs = make([]Arc, 0)
	var weight int
	var meid int
	if me {
		meid = GraphNodeId(&graph, "I")
	}
	for _, line := range lines {
		tokens := strings.Split(line, " ")
		fromname := tokens[0]
		toname := tokens[10][:len(tokens[10])-1] // remove trailing .
		if tokens[2] == "gain" {
			weight = Atoi(tokens[3])
		} else if tokens[2] == "lose" {
			weight = -1 * Atoi(tokens[3])
		} else {
			log.Fatalf("Parse error (no gain/lose) on line: %v\n", line)
		}
		from := GraphNodeId(&graph, fromname)
		to := GraphNodeId(&graph, toname)
		arc := Arc{from, to, weight}
		graph.arcs = append(graph.arcs, arc)
	}
	size := len(graph.names)
	graph.weights = make([][]int, size)
	for i := 0; i < size; i++ {
		graph.weights[i] = make([]int, size)
	}
	for _, arc := range graph.arcs {
		graph.weights[arc.from][arc.to] = arc.weight
	}
	// for me, we just need to set the weights, not the arcs
	if me {
		for i := 0; i < size; i++ {
			graph.weights[i][meid] = 0
			graph.weights[meid][i] = 0
		}
	}
	return graph
}

func GraphNodeId(graph *Graph, name string) int {
	id, ok := graph.ids[name]
	if !ok {
		id = len(graph.names)
		graph.names = append(graph.names, name)
		graph.ids[name] = id
	}
	return id
}

func RouteLen(graph *Graph, route []int) int {
	l := 0
	size := len(route)
	// middle cases
	for i := 1; i < size-1; i++ {
		l += graph.weights[route[i]][route[i-1]] + graph.weights[route[i]][route[i+1]]
	}
	// at ends (0 & size-1), wrap around
	l += graph.weights[route[0]][route[size-1]] + graph.weights[route[0]][route[1]]
	l += graph.weights[route[size-1]][route[size-2]] + graph.weights[route[size-1]][route[0]]
	if verbose {
		PrintRoute(graph, route, l) // trace
	}
	return l
}

func PrintRoute(graph *Graph, route []int, l int) {
	fmt.Printf("  %v", graph.names[route[0]])
	for i := 1; i < len(route); i++ {
		fmt.Printf(" -> %v", graph.names[route[i]])
	}
	fmt.Printf(" = %v\n", l)
}

func ShortestRoute(graph *Graph) int {
	size := len(graph.names)
	graph.min = 8888888888888888888 // nicer "maxint"
	route := make([]int, size)
	for i := 0; i < size; i++ {
		route[i] = i
	}
	FindShortestRoute(graph, route, 0, size-1)
	return graph.min
}

// from https://golangbyexample.com/all-permutations-string-golang/
func FindShortestRoute(graph *Graph, route []int, from, to int) {
	if from == to {
		l := RouteLen(graph, route)
		if l < graph.min {
			graph.min = l
		}
	} else {
		for i := from; i <= to; i++ {
			route[from], route[i] = route[i], route[from]
			FindShortestRoute(graph, route, from+1, to)
			route[from], route[i] = route[i], route[from]
		}
	}
}

func LongestRoute(graph *Graph) int {
	size := len(graph.names)
	graph.max = 0
	route := make([]int, size)
	for i := 0; i < size; i++ {
		route[i] = i
	}
	FindLongestRoute(graph, route, 0, size-1)
	return graph.max
}

// from https://golangbyexample.com/all-permutations-string-golang/
func FindLongestRoute(graph *Graph, route []int, from, to int) {
	if from == to {
		l := RouteLen(graph, route)
		if l > graph.max {
			graph.max = l
		}
	} else {
		for i := from; i <= to; i++ {
			route[from], route[i] = route[i], route[from]
			FindLongestRoute(graph, route, from+1, to)
			route[from], route[i] = route[i], route[from]
		}
	}
}

//////////// Generic code

// useful in tests to feed Part1 & Part2 with a simple string (with newlines)
func StringToLines(s string) (lines []string) {
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	err := scanner.Err()
	if err != nil {
		log.Fatal(err)
	}
	return
}

// read the input file into a string array for feeding Parts
func FileToLines(filePath string) (lines []string) {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	// optionally, resize scanner's capacity for lines over 64K (65536)
	const maxCapacity = 1000000 // your required line length
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)
	// end optional
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	err = scanner.Err()
	if err != nil {
		log.Fatal(err)
	}

	return
}

// simplified functions to not bother with error handling. Just abort.

func Atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal(err)
	}
	return i
}

// for completeness
func Itoa(i int) string {
	return strconv.Itoa(i)
}
