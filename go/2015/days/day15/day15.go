// Adventofcode 2015, day15, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 input 18965440
// TEST: input 15862900
package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"
)

var verbose bool

type (
	// Ingredient holds the caracteristics of an ingredient type (Cinnamon, Sugar, ...)
	Ingredient struct {
		name       string
		id         int // index in the Stock.list
		capacity   int // the fields...
		durability int
		flavor     int
		texture    int
		calories   int
	}

	// Stock lists the type of Ingredients we have in stock
	Stock struct {
		ids   map[string]int // name -> index (id)
		names []string       // id -> name
		list  []Ingredient   // list of available Ingredients
	}
)

func main() {
	partOne := flag.Bool("1", false, "run part one code, instead of part 2 (default)")
	verboseFlag := flag.Bool("v", false, "verbose: print routes")
	flag.Parse()
	verbose = *verboseFlag
	infile := "input.txt"
	if flag.NArg() > 0 {
		infile = flag.Arg(0)
	}
	lines := fileToLines(infile)
	stock := parseIngredients(lines)

	var result int
	if *partOne {
		VP("Running Part1")
		result = part1(stock)
	} else {
		VP("Running Part2")
		result = part2(stock)
	}
	fmt.Println(result)
}

//////////// Part 1

func part1(stock *Stock) int {
	return maxRecipeScore(stock, 0, 100)
}

//////////// Part 2
func part2(stock *Stock) int {
	return maxRecipeScoreCalories(stock, 0, 100, 500)
}

//////////// Common Parts code

func parseIngredients(lines []string) *Stock {
	rn := `[[:space:]](-?[[:digit:]]+),?[[:space:]]*`
	re := regexp.MustCompile(`^([[:alpha:]]+):[[:space:]]*capacity` + rn + `durability` + rn + `flavor` + rn + `texture` + rn + `calories` + rn + `$`)
	stock := new(Stock)
	stock.ids = make(map[string]int, 0)
	stock.names = make([]string, 0)
	stock.list = make([]Ingredient, 0)
	for _, line := range lines {
		c := re.FindStringSubmatch(line)
		if c != nil {
			ing := Ingredient{c[1], len(stock.list), atoi(c[2]), atoi(c[3]), atoi(c[4]), atoi(c[5]), atoi(c[6])}
			stock.names = append(stock.names, ing.name)
			stock.ids[ing.name] = ing.id
			stock.list = append(stock.list, ing)
		} else if line != "" {
			log.Fatalf("Parse error on line: %v", line)
		}
	}

	return stock
}

func recipeScore(stock *Stock, recipe []int) int {
	capacity := 0
	for i := 0; i < len(stock.list); i++ {
		capacity += stock.list[i].capacity * recipe[i]
	}
	if capacity <= 0 {
		return 0
	}
	durability := 0
	for i := 0; i < len(stock.list); i++ {
		durability += stock.list[i].durability * recipe[i]
	}
	if durability <= 0 {
		return 0
	}
	flavor := 0
	for i := 0; i < len(stock.list); i++ {
		flavor += stock.list[i].flavor * recipe[i]
	}
	if flavor <= 0 {
		return 0
	}
	texture := 0
	for i := 0; i < len(stock.list); i++ {
		texture += stock.list[i].texture * recipe[i]
	}
	if texture <= 0 {
		return 0
	}
	return capacity * durability * flavor * texture
}

//////////// Part1 functions

// max score possible for recipes with only ingredients of id >= from
// and max total amount amount
func maxRecipeScore(stock *Stock, from int, total int) int {
	size := len(stock.list)
	recipe := make([]int, size) // quantity of ingredients. must sum to total
	return maxRecipesScore(stock, recipe, size, 0, total)
}

// We take the first ingredient in the stock.list, and for each possible values of
// teaspoons (0 to "total" = 100) , we list this value followed recursively by all
// the possible was to dispatch the rest of possible teaspoons ("total" minus the
// ones on previous ingredients) among the reamianing ingredients (the ones
// starting at the "from" index)

func maxRecipesScore(stock *Stock, r []int, size, from, total int) (score int) {
	score = 0
	if from == size-1 {
		for i := 0; i <= total; i++ {
			r[from] = i
			rs := recipeScore(stock, r)
			if rs > score {
				score = rs
			}
		}
	} else {
		for i := 0; i <= total; i++ {
			r[from] = i
			rs := maxRecipesScore(stock, r, size, from+1, total-i)
			if rs > score {
				score = rs
			}
		}
	}
	return
}

//////////// Part2 functions

// the same functions as for Part1, but only considering recipes with a fixed
// total amount of calories

func maxRecipeScoreCalories(stock *Stock, from, total, calories int) int {
	size := len(stock.list)
	recipe := make([]int, size) // quantity of ingredients. must sum to total
	return maxRecipesScoreCalories(stock, recipe, size, 0, total, calories)
}

func recipeCalories(stock *Stock, recipe []int) (calories int) {
	calories = 0
	for i := 0; i < len(stock.list); i++ {
		calories += stock.list[i].calories * recipe[i]
	}
	return
}

func maxRecipesScoreCalories(stock *Stock, r []int, size, from, total, calories int) (score int) {
	score = 0
	if from == size-1 {
		for i := 0; i <= total; i++ {
			r[from] = i
			if recipeCalories(stock, r) == calories {
				rs := recipeScore(stock, r)
				if rs > score {
					score = rs
				}
			}
		}
	} else {
		for i := 0; i <= total; i++ {
			r[from] = i
			rs := maxRecipesScoreCalories(stock, r, size, from+1, total-i, calories)
			if rs > score {
				score = rs
			}
		}
	}
	return
}
