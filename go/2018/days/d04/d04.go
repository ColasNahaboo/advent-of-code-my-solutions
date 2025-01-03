// Adventofcode 2018, d04, in go. https://adventofcode.com/2018/day/04
// Arguments: -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 240
// TEST: example 4455
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

// What we could remark on the input:
// - not in order
// - only one shift change per day, but they can happen the day before at 23:xx
// = if change at 23:xx, no falls or wakes happen before 00:00
// - only falls or change happen at 00:00, never wakes
// - Guard shift is the first thing of the day, or previous 23:xx
// - first line is a shift, but not the last

package main

import (
	"fmt"
	"regexp"
	"strings"
	// "flag"
	"sort"
)

type Guard struct {
	id int
	days [][]Sleep
	sleep int
}
type Sleep struct {
	from, to int
}

//////////// Options parsing & exec parts

func main() {
	ExecOptions(2, NoXtraOpts, part1, part2)
}

//////////// Part 1

func part1(lines []string) (res int) {
	guards := parse(lines)
	VP(guards)
	var	maxsleep, maxid, maxtimes, maxmn int
	for _, g := range guards {
		if g.sleep > maxsleep {
			maxsleep = g.sleep
			maxid = g.id
		}
	}
	mns := [60]int{}
	for _, day := range GetGuard(&guards, maxid).days {
		for _, sleep := range day {
			for i := sleep.from; i < sleep.to; i++ {
				mns[i]++
			}
		}
	}
	for mn, n := range mns {
		if n > maxtimes {
			maxtimes = n
			maxmn = mn
		}
	}
	VPf("max id: %d, max minute: %d\n", maxid, maxmn)
	return maxid * maxmn
}

//////////// Part 2

func part2(lines []string) (res int) {
	guards := parse(lines)
	var maxid, maxmn, maxtimes int
	for _, g := range guards {
		mns := [60]int{}
		for _, day := range g.days {
			for _, sleep := range day {
				for i := sleep.from; i < sleep.to; i++ {
					mns[i]++
				}
			}
		}
		for mn, n := range mns {
			if n > maxtimes {
				maxtimes = n
				maxid = g.id
				maxmn = mn
			}
		}
	}
	VPf("Max spent asleep: %d mn on minute %d for guard #%d\n", maxtimes, maxmn, maxid)
	return  maxid * maxmn
}

//////////// Common Parts code

type Day struct {
	day int
	gid int
	sleeps []Sleep
}

func parse(lines []string) (guards []Guard) {
	// first sort the input
	sort.Slice(lines, func(i, j int) bool { return lines[i] < lines[j] })
	
	// 1=day, 2=hour 3=mn 4=Guard,falls,wakes, 5= 6=gid
	renum := regexp.MustCompile("^[[][0-9]+-([-0-9]+) ([0-9]{2}):([0-9]{2})[]] (.{5}) (#([0-9]+))?")
	var day, hour, mn, gid, lastfalls int
	var act string
	var m []string

	days := []Day{}
	dayn := -1					// index of current day in days
	for lineno, line := range lines {
		if m = renum.FindStringSubmatch(line); m == nil {
			panicf("Syntax error line %d: \"%s\"", lineno+1, line)
		}
		day = atoi(strings.ReplaceAll(m[1], "-", ""))
		hour = atoi(m[2])
		mn = atoi(m[3])
		act = m[4]
		if act == "Guard" {
			gid = atoi(m[6])
			lastfalls = -1
			if hour == 23 {			// just set gid for the next day, 00:00
				continue
			}
		}
		if dayn < 0 || day != days[dayn].day { // start a new day
			dayn = len(days)
			days = append(days, Day{day, gid, []Sleep{}})
		}
		if act == "wakes" {
			if lastfalls < 0 {
				panicf("Error line %d, gid #%d: wakes no falls", gid, lineno+1)
			}
			CompleteSleep(&days, dayn, Sleep{lastfalls, mn})
			lastfalls = -1
		} else if act == "falls" {
			lastfalls = mn
		}
	}
	if lastfalls >= 0 { 		// was asleep
		CompleteSleep(&days, dayn, Sleep{lastfalls, 60})
	}
	for _, day := range days {
		g := GetGuard(&guards, day.gid)
		g.days = append(g.days, day.sleeps)
		g.sleep += TotalSleeps(day.sleeps)
	}
	return
}

func CompleteSleep(days *[]Day, dayn int, s Sleep) {
	(*days)[dayn].sleeps = append((*days)[dayn].sleeps, s)
}

func TotalSleeps(sleeps []Sleep) (res int) {
	for _, s := range sleeps {
		res += s.to - s.from
	}
	return
}

func GetGuard(guards *[]Guard, id int) *Guard {
	for i, g := range *guards {
		if g.id == id {
			return &((*guards)[i])
		}
	}
	g := Guard{id: id, days: [][]Sleep{}}
	*guards = append(*guards, g)
	// return a pointer to the actual guard instance member of the guards slice
	return &((*guards)[len(*guards)-1])
}
	
	
//////////// PrettyPrinting & Debugging functions. See also the VPx functions.

func DEBUG() {
	if ! debug { return }
	fmt.Println("DEBUG!")
	// ad hoc debug code here
}
