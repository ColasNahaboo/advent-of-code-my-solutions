// Adventofcode 2015, day09, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 input
// TEST: input
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Graph struct {
	names   []string
	ids     map[string]int
	arcs    []Arc
	weights [][]int
	min     int
	max     int
}
type Arc struct {
	from   int
	to     int
	weight int
}

func main() {
	partOne := flag.Bool("1", false, "run part one code, instead of part 2 (default)")
	flag.Parse()
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
	graph := NewGraph(lines)
	return ShortestRoute(&graph)
}

//////////// Part 2
func Part2(lines []string) int {
	graph := NewGraph(lines)
	return LongestRoute(&graph)
}

//////////// Common Parts code

func NewGraph(lines []string) Graph {
	var graph Graph
	graph.names = make([]string, 0)
	graph.ids = make(map[string]int)
	graph.arcs = make([]Arc, 0)
	for _, line := range lines {
		tokens := strings.Split(line, " ")
		fromname := tokens[0]
		toname := tokens[2]
		weight, err := strconv.Atoi(tokens[4])
		if err != nil {
			log.Fatal(err)
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
		graph.weights[arc.to][arc.from] = arc.weight
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
	for i := 1; i < len(route); i++ {
		l += graph.weights[route[i-1]][route[i]]
	}
	PrintRoute(graph, route, l) // trace
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
