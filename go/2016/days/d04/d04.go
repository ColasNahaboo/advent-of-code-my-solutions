// Adventofcode 2016, d04, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 input 137896
// TEST: input 501
package main

import (
	"flag"
	"fmt"
	"log"
	"sort"
	"regexp"
)

type Room struct {
	name string
	sector int
	cksum string
}

// The data type & method to be able to sort letters via the sort package

type LetterCount struct {
	letter string
	count int
}

var verbose bool
var reroom = regexp.MustCompile("^([-[:lower:]]+)-([[:digit:]]+)[[]([[:lower:]]{5})[]]") 
var rechar = regexp.MustCompile("[[:lower:]]")
var rechardash = regexp.MustCompile("[-[:lower:]]")

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	flag.Parse()
	verbose = *verboseFlag
	infile := "input.txt"
	if flag.NArg() > 0 {
		infile = flag.Arg(0)
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

func part1(lines []string) (result int) {
	rooms := parseRooms(lines)
	for _, room := range rooms {
		result += room.sector
	}
	return result
}

//////////// Part 2

const storage = "northpole object storage" // room name to find

func part2(lines []string) int {
	rooms := parseRooms(lines)
	for _, room := range rooms {
		dname := nameDecrypt(room)
		VP(room.name, "==>", dname)
		if dname == storage {
			return room.sector
		}
	}
	return 0
}

//////////// Common Parts code

// return the list of real rooms
func parseRooms(lines []string) (rooms []Room) {
	var i int
	for _, line := range lines {
		room := reroom.FindStringSubmatch(line)
		if room == nil {
			log.Fatalln("Room syntax error: " + line)
		}
		name := room[1]
		sector := room[2]
		roomchk := room[3]
		// map letters and how many of them
		counts := map[string]int{} 
		for _, letter := range rechar.FindAllString(name, -1) {
			counts[letter]++
		}
		// we create a slice of LetterCounts to sort it
		lcs := make([]LetterCount, len(counts))
		i = 0
		for char, count := range counts {
			lcs[i] = LetterCount{char, count}
			i++
		}
		sort.Slice(lcs, func(i, j int) bool {
			if lcs[i].count > lcs[j].count {
				return true
			} else if (lcs[i].count == lcs[j].count) && (lcs[i].letter < lcs[j].letter) {
				return true
			} else {
				return false
			}
		})
		cksum := ""
		for i = 0; i < 5; i++ {
			cksum += lcs[i].letter
		}
		if cksum == roomchk {
			VP("Room:", name, "is real")
			rooms = append(rooms, Room{name, atoi(sector), roomchk})
		} else {
			VP("Room:", name, "has real checksum", cksum, "instead of", roomchk)
		}
	}
	return
}

//////////// Part1 functions

//////////// Part2 functions

func nameDecrypt(room Room) (dname string) {
	for _, letter := range rechardash.FindAllString(room.name, -1) {
		if letter == "-" {
			dname += " "
		} else {
			dname += rotate(letter, room.sector)
		}
	}
	return
}

func rotate(letter string, inc int) string {
	// 26 letters starting at "a" = 97
	n := rune(97 + (int(letter[0]) - 97 + inc ) % 26)
	return string([]rune{n})
}
