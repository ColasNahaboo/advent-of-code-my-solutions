// from: https://stackoverflow.com/questions/8757389/reading-a-file-line-by-line-in-go/16615559#16615559
package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func StringToLines(s string) (lines []string) {
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	err := scanner.Err()
	if err != nil {
		os.Exit(1)
	}
	return
}

func FileToLines(filePath string) (lines []string) {
	f, err := os.Open(filePath)
	if err != nil {
		return
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
		os.Exit(1)
	}

	return
}

// read the whole input file into a single string
func FileToString(filename string) string {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Print(err)
	}
	return string(bytes)
}

// read the whole input file into a byte array
func FileToBytes(filename string) []byte {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Print(err)
	}
	return bytes
}
