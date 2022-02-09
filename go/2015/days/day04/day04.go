// Adventofcode 2015, day04, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 input 117946
// TEST: input   3938038
package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func main() {
	partOne := flag.Bool("1", false, "run part one code, instead of part 2 (default)")
	flag.Parse()
	infile := "input.txt"
	if flag.NArg() > 0 {
		infile = flag.Arg(0)
	}
	input, err := ioutil.ReadFile(infile)
	if err != nil {
		os.Exit(1)
	}

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

func Part1(key []byte) int {
	return lowestnum(key, "00000")
}

func Part2(key []byte) int {
	return lowestnum(key, "000000")
}

func lowestnum(key []byte, prefix string) int {
	keystr := string(key)
	for i := 0; true; i++ {
		data := keystr + strconv.Itoa(i)
		hash := md5.Sum([]byte(data))
		if strings.HasPrefix(hex.EncodeToString(hash[:]), prefix) {
			return i
		}
	}
	return 0
}
