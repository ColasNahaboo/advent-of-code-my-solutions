// Adventofcode 2016, d10, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 input 147
// TEST: example 30
// TEST: input 55637
package main

import (
	"flag"
	"fmt"
	"regexp"
	"github.com/gammazero/deque" // fast FIFO & LIFO
)

var inputs [][2]int				// {value, bot-to}
var outputs []int				// [bin#]value in output bins
var rules []Rule				// [bot#]rule for dispatch
type Rule struct {				
	low, high int				// dispatch to: <0 -> outbins+1, >0 -> bot
}
var todo *deque.Deque[int]					// bots that have 2 chips, ready to dispatch
var bots []Bot
type Bot struct {				// bot state
	low, high int				// current values in hands
}

var verbose bool

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

func part1(lines []string) int {
	parse(lines)
	return responsible(17, 61)
}

//////////// Part 2
func part2(lines []string) int {
	parse(lines)
	responsible(0, 0)
	return outputs[0] * outputs[1] * outputs[2]
}

//////////// Common Parts code

func parse(lines []string) {
	reLine := regexp.MustCompile("bot ([[:digit:]]+) gives low to (output|bot) ([[:digit:]]+) and high to (output|bot) ([[:digit:]]+)")
	maxoutputs := maxValue("value", lines) + 1
	outputs = make([]int, maxoutputs, maxoutputs) 
	maxbots := maxValue("bot", lines) + 1
	bots = make([]Bot, maxbots, maxbots)
	rules = make([]Rule, maxbots, maxbots)
	todo = deque.New[int](100)
	var value, bot int
	for ln, line := range lines {
		if n, _ := fmt.Sscanf(line, "value %d goes to bot %d", &value, &bot); n == 2 {
			inputs = append(inputs, [2]int{value, bot})
		} else if m := reLine.FindStringSubmatch(line); m != nil {
			bot := atoi(m[1])
			low := atoi(m[3])
			if m[2] == "output" { low = -1 - low;}
			high := atoi(m[5])
			if m[4] == "output" { high = -1 - high;}
			rules[bot] = Rule{low, high} // logic
			bots[bot] = Bot{}			 // state
		} else {
			panic(fmt.Sprintf("Parse error line %d: \"%s\"\n", ln, line))
		}
	}
}

// return maximum N found in all substrings "tag N"
func maxValue(tag string, lines []string) (n int) {
	reVal := regexp.MustCompile(tag + " ([[:digit:]]+)")
	for _, line := range lines {
		if matches := reVal.FindAllStringSubmatch(line, -1); matches != nil {
			for _, match := range matches {
				val := atoi(match[1])
				if val > n { n = val;}
			}
		}
	}
	return
}

func giveTo(val, bot int) {
	VPf("giving %d to %d\n", val, bot)
	if bot < 0 {				// output bin
		outputs[-1 - bot] = val
		return
	}
	if val <= 0 {
		panic(fmt.Sprintf("Tring to give value %d to bot#%d\n", val, bot))
	}
	if bots[bot].low != 0 {
		if bots[bot].low < val {
			if bots[bot].high != 0 {
				panic(fmt.Sprintf("Tring to give %d to bot#%d having %v\n", val, bot, bots[bot]))
			}
			bots[bot].high = val
		} else {
			bots[bot].low, bots[bot].high = val, bots[bot].low
		}
		todo.PushBack(bot)		//  complete! put bot on todo list to dispatch
	} else if bots[bot].high == 0 {
		bots[bot].low = val
	} else {
		panic(fmt.Sprintf("Bot #%d had high (%d) but no low!\n", bot,  bots[bot].high))
	}
}

//////////// Part1 functions

// run simulation to find bot# comparing low and high values
//  low value == 0 ==> run whole simul
func responsible(lv, hv int) int {
	// distribute from input bins
	for _, in := range inputs {
		giveTo(in[0], in[1])
		if bots[in[1]].low == lv && bots[in[1]].high == hv {
			return in[1]
		}
	}
	// run simulation: process the todo list as a FIFO
	for todo.Len() > 0 {
		bot := todo.PopFront() // de-stack next bot
		lval := bots[bot].low
		bots[bot].low = 0
		hval := bots[bot].high
		bots[bot].high = 0
		giveTo(lval, rules[bot].low)
		if rules[bot].low >= 0 && bots[rules[bot].low].low == lv && bots[rules[bot].low].high == hv {
			return rules[bot].low
		}
		giveTo(hval, rules[bot].high)
		if rules[bot].high >= 0 && bots[rules[bot].high].low == lv && bots[rules[bot].high].high == hv {
			return rules[bot].high
		}
	}
	if lv != 0 {
		panic("No bot responsible found!\n")
	}
	return 0
}

//////////// Part2 functions
