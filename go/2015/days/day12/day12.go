// Adventofcode 2015, day12, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 input 191164
// TEST: input 87842
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	// "regexp"
)

func main() {
	partOne := flag.Bool("1", false, "run part one code, instead of part 2 (default)")
	flag.Parse()
	infile := "input.txt"
	if flag.NArg() > 0 {
		infile = flag.Arg(0)
	}
	input := FileToBytes(infile)

	var result int
	if *partOne {
		fmt.Println("Running Part1")
		result = Part1(input)
	} else {
		fmt.Println("Running Part2")
		result = Part2(input)
	}
	fmt.Println(result)
}

//////////// Part 1

func Part1(input []byte) int {
	var root interface{}
	err := json.Unmarshal(input, &root)
	if err != nil {
		log.Fatal(err)
	}
	return CollectJsonNumbers(root)
}

func CollectJsonNumbers(inode interface{}) int {
	switch node := inode.(type) {
	case float64:
		return int(node)
	case []interface{}:
		n := 0
		for _, elt := range node {
			n = n + CollectJsonNumbers(elt)
		}
		return n
	case map[string]interface{}:
		n := 0
		for _, elt := range node {
			n = n + CollectJsonNumbers(elt)
		}
		return n
	default:
		return 0
	}
}

//////////// Part 2
func Part2(input []byte) int {
	var root interface{}
	err := json.Unmarshal(input, &root)
	if err != nil {
		log.Fatal(err)
	}
	return CollectJsonNumbersNonRed(root)
}

func CollectJsonNumbersNonRed(inode interface{}) int {
	switch node := inode.(type) {
	case float64:
		return int(node)
	case []interface{}:
		n := 0
		for _, elt := range node {
			n = n + CollectJsonNumbersNonRed(elt)
		}
		return n
	case map[string]interface{}:
		n := 0
		for _, elt := range node {
			switch s := elt.(type) {
			case string:
				if s == "red" { // skip whole object if one field value is "red"
					return 0
				}
			default:
				n = n + CollectJsonNumbersNonRed(elt)
			}
		}
		return n
	default:
		return 0
	}
}

//////////// Common Parts code

//////////// Generic code

// read the whole input file into a byte array
func FileToBytes(filename string) []byte {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Print(err)
	}
	return bytes
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
