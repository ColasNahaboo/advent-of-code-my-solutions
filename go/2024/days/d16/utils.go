// These are various convenient functions commonly used in my adventofcode
// solutions, that are to be copied in each day directory

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"regexp"
)

/////////// Useful constants

const MaxUint = ^uint(0)
const MinUint = 0
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1
// an easier to spot maxint in debug than 9223372036854775807 (^uint(0) >> 1)
const maxint = 8888888888888888888

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

//////////// Find the first file in current dir matching a regexp

// If not found, return the argument regexp
func fileMatch(regname string) string {
	re := regexp.MustCompile("^" + regname + "$")
	files, err := ioutil.ReadDir(".")
	if err != nil {
		panic("Cannot read current directory")
	}
	for _, file := range files {
		if re.MatchString(file.Name()) {
			return file.Name()
		}
	}
	return regname
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

func intAbs(i int) int {
	if i >= 0 {
		return i
	} else {
		return -i
	}
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

// same, mapped on lists for convenience
func atoil(ss []string) []int {
	is := make([]int, len(ss), len(ss))
	for i, s := range ss {
		is[i] = atoi(s)
	}
	return is
}
func itoal(is []int) []string {
	ss := make([]string, len(is), len(is))
	for i, n := range is {
		ss[i] = itoa(n)
	}
	return ss
}

////// Lines of text

// split lines into non-empty lines blocks that were separated by empty lines
func LineBlocks(lines []string) (blocks [][]string) {
	b := []string{}
	for _, line := range lines {
		if len(line) == 0 {
			if len(b) != 0 {
				blocks = append(blocks, b)
				b = b[:0]
			}
		} else {
			b = append(b, line)
		}
	}
	if len(b) != 0 {
		blocks = append(blocks, b)
	}
	return
}

////// Slice utils (maybe bow obsoleted by the slices package)

func prependInt(x []int, y int) []int {
    x = append(x, 0)
    copy(x[1:], x)
    x[0] = y
    return x
}

func insertInt(list []int, i, v int) []int {
	l := make([]int, len(list)+1)
    copy(l[:i], list[:i])
	l[i] = v
	if i < len(list) {
		copy(l[i+1:], list[i:])
	}
    return l
}

func indexOfInt(list []int, number int) (int) {
   for i, v := range list {
       if number == v {
           return i
       }
   }
   return -1    //not found.
}

// safe delete, keeps list order
func deleteOrderInt(list []int, number int) ([]int, bool) {
   for i, v := range list {
       if number == v {
		   l := make([]int, len(list) - 1)
		   copy(l, list[:i])
		   copy(l[i:], list[i+1:])
		   return l, true   // Truncate slice.
       }
   }
   return list, false    // not found.
}

// safe delete, does not keep list order (faster)
func deleteInt(list []int, number int) ([]int, bool) {
	l := make([]int, len(list))
	copy(l, list)
	return deleteFastInt(l, number)
}

// fastest delete, does not keep list order
// Should not be used if other slices point to list
func deleteFastInt(list []int, number int) ([]int, bool) {
   for i, v := range list {
       if number == v {
		   list[i] = list[len(list) - 1] // Copy last elt to index i.
		   return list[:len(list)-1], true   // Truncate slice.
       }
   }
   return list, false    // not found.
}

func sliceIntEquals(l1, l2 []int) bool {
	if len(l1) != len(l2) {
		return false
	}
	for i, v := range l1 {
       if v != l2[i] {
           return false
       }
   }
	return true
}


// Generic versions

func prependElt[T comparable](list []T, elt T) []T {
    list = append(list, *new(T))
    copy(list[1:], list)
    list[0] = elt
    return list
}


func insertElt[T comparable](list []T, i int, elt T) []T {
	l := make([]T, len(list)+1)
    copy(l[:i], list[:i])
	l[i] = elt
	if i < len(list) {
		copy(l[i+1:], list[i:])
	}
    return l
}

func IndexOf[T comparable](collection []T, elt T) int {
    for i, x := range collection {
        if x == elt {
            return i
        }
    }
    return -1
}

// safe delete, keeps list order
func deleteOrderElt[T comparable](list []T, elt T) ([]T, bool) {
   for i, v := range list {
       if elt == v {
		   l := make([]T, len(list) - 1)
		   copy(l, list[:i])
		   copy(l[i:], list[i+1:])
		   return l, true   // Truncate slice.
       }
   }
   return list, false    // not found.
}

// safe delete, does not keep list order (faster)
func deleteElt[T comparable](list []T, elt T) ([]T, bool) {
	l := make([]T, len(list))
	copy(l, list)
	return deleteFastElt(l, elt)
}

// fastest delete, does not keep list order
// Should not be used if other slices point to list
func deleteFastElt[T comparable](list []T, elt T) ([]T, bool) {
   for i, v := range list {
       if elt == v {
		   list[i] = list[len(list) - 1] // Copy last elt to index i.
		   return list[:len(list)-1], true   // Truncate slice.
       }
   }
   return list, false			// could not find it
}

func sliceEquals[T comparable](a []T, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
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
