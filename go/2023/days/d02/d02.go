// Adventofcode 2023, d02, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: d02-input,RESULT1,RESULT2.test
// TEST: -1 example 8
// TEST: example 2286
// And any file named d02-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"flag"
	"fmt"
	"regexp"
)

// token types & index in Draw
const red = 0
const green = 1
const blue = 2
const eog = -1			// newline = end of game
const eod = -2			// ";" = end of draw
var colornames = []string{"red", "green", "blue"}

type Game []Draw
type Draw [3]int		// red, green, blue values
type Token struct {
	color, n int
}


var verbose bool

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	flag.Parse()
	verbose = *verboseFlag
	var infile string
	if flag.NArg() > 0 {
		infile = flag.Arg(0)
	} else {
		infile = fileMatch("input,[0-9]*,[0-9]*.test")
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
	games := parseGames(lines)
	gameNum := 1
	for _, game := range games {
		for _, draw := range game {
			if ! isValidDraw(draw) {
				VPf("Game #%d impossible: %v\n", gameNum, draw)
				goto IMPOSSIBLE
			}
		}
		VPf("Game #%d OK %v\n", gameNum, game)
		result += gameNum
	IMPOSSIBLE:
		gameNum++
	}
	return
}

func isValidDraw(d Draw) bool {
	return d[red] <= 12 && d[green] <= 13 && d[blue] <= 14
}

//////////// Part 2

func part2(lines []string) (result int) {
	games := parseGames(lines)
	gameNum := 1
	for _, game := range games {
		must := [3]int{}
		for _, draw := range game {
			for color := 0; color < 3; color++ {
				if draw[color] > must[color] {
					must[color] = draw[color]
				}
			}
		}
		power := must[red] * must[green] * must[blue]
		VPf("Game #%d has power %v\n", gameNum, power)
		result += power
		gameNum++
	}
	return result
}

//////////// Common Parts code

func parseGames(lines []string) (games []Game) {
	tokens := tokens(lines)
	game := make(Game, 0)
	draw := Draw{} 
	for _, token := range tokens {
		if token.color == eog {
			game = append(game, draw)
			draw = Draw{}
			games = append(games, game)
			game = make(Game, 0)
		} else if token.color == eod {
			game = append(game, draw)
			draw = Draw{}
		} else {
			draw[token.color] = token.n
		}
	}
	return
}

func tokens(lines []string) (tokens []Token) {
	for _, line := range lines {
		token_re := regexp.MustCompile("((;)|([[:digit:]]+)[[:space:]]+(red|green|blue))")
		matches := token_re.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			if len(match[2]) > 0 {
				tokens = append(tokens, Token{eod, 1})
			} else {
				tokens = append(tokens, Token{c2i(match[4]), atoi(match[3])})
			}
		}
		tokens = append(tokens, Token{eog, 1})
	}
	return
}

func c2i(name string) int {
	if name == "red" { return red
	} else if name == "green" { return green
	} else if name == "blue" { return blue
	} else {
		panic("Unknown color: " + name)
	}
}

func i2c(i int) string {
	return colornames[i]
}
	
		
	
//////////// Part1 functions

//////////// Part2 functions
