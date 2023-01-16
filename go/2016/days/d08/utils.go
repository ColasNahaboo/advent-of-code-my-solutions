package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

//////////// Reading a file in memory

// "reads" a string into an array of strings, one per line
// useful to simulate a fileToLines in tests
func stringToLines(s string) (lines []string) {
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
func fileToLines(filePath string) (lines []string) {
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

// read the whole input file into a single string
func fileToString(filename string) string {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Print(err)
	}
	return string(bytes)
}

// read the whole input file into a byte array
func fileToBytes(filename string) []byte {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Print(err)
	}
	return bytes
}

//////////// Exponentials / Power of integers. Use package math/big for >64bits

// IntPower calculates n to the mth power.
// it avoids messing with the rounding problems of float arthmeticm
// and is faster

func intPower(n int, m int) int {
	if m <= 0 {
		return 1
	}
	result := n
	for i := 2; i <= m; i++ {
		result *= n
	}
	return result
}

// IntPower64 calculates n to the mth power.
// version explicitely using uint64 for results
func intPower64(n int, m int) uint64 {
	if m <= 0 {
		return 1
	}
	result := uint64(n)
	for i := 2; i <= m; i++ {
		result *= uint64(n)
	}
	return result
}

//////////// Convenience functions

////// simplified functions to not bother with error handling. Just abort.

func atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal(err)
	}
	return i
}

// for completeness
func itoa(i int) string {
	return strconv.Itoa(i)
}

////// Slice utils

func prependInt(x []int, y int) []int {
    x = append(x, 0)
    copy(x[1:], x)
    x[0] = y
    return x
}

func indexOfInt(list []int, number int) (int) {
   for i, v := range list {
       if number == v {
           return i
       }
   }
   return -1    //not found.
}

// Generic versions

func prepend[T comparable](list []T, elt T) []T {
    list = append(list, *new(T))
    copy(list[1:], list)
    list[0] = elt
    return list
}


func IndexOf[T comparable](collection []T, el T) int {
    for i, x := range collection {
        if x == el {
            return i
        }
    }
    return -1
}

////// VP ("Verbose Print") wrappers: print only if verbose var is true

// VP = fmt.Println
func VP(v ...interface{}) {
	if verbose {
		fmt.Println(v...)
	}
}

// VPn (echo -n) = fmt.Print
func VPn(v ...interface{}) {
	if verbose {
		fmt.Print(v...)
	}
}

// VPf = fmt.Printf
func VPf(f string, v ...interface{}) {
	if verbose {
		fmt.Printf(f, v...)
	}
}

////// Convenience panic
func panicf(f string, v ...interface{}) {
	panic(fmt.Sprintf(f, v...))
}
