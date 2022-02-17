// Adventofcode 2015, day21, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 input 111
// TEST: input 188
package main

import (
	"flag"
	"fmt"
	"strings"
	// "regexp"
)

type (
	mob struct {
		hp  int
		dmg int
		ac  int
	}
	item struct {
		name string
		cost int
		dmg  int
		ac   int
	}
	shop struct {
		weapons []item
		armors  []item
		rings   []item
	}
)

var verbose bool

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	shopfileFlag := flag.String("s", "shop.txt", "shop file")
	flag.Parse()
	verbose = *verboseFlag
	shopfile := *shopfileFlag
	infile := "input.txt"
	if flag.NArg() > 0 {
		infile = flag.Arg(0)
	}
	lines := fileToLines(infile)
	boss := readMob(lines)
	theShop := readShop(shopfile)

	var result int
	if *partOne {
		VP("Running Part1")
		result = part1(boss, theShop)
	} else {
		VP("Running Part2")
		result = part2(boss, theShop)
	}
	fmt.Println(result)
}

//////////// Part 1

func part1(boss mob, s shop) int {
	hero := mob{100, 0, 0}
	mincost := 1 << 31
	var gear []string
	for _, weapon := range s.weapons {
		for _, armor := range s.armors {
			for _, ring1 := range s.rings {
				for _, ring2 := range s.rings {
					if ring2.name != ring1.name { // rings are different
						hero.dmg = weapon.dmg + armor.dmg + ring1.dmg + ring2.dmg
						hero.ac = weapon.ac + armor.ac + ring1.ac + ring2.ac
						if winFight(hero, boss) {
							cost := weapon.cost + armor.cost + ring1.cost + ring2.cost
							if cost < mincost {
								mincost = cost
								gear = []string{weapon.name, armor.name, ring1.name, ring2.name}
							}
						}
					}
				}
			}
		}
	}
	VP("Minimal gear:", gear)
	return mincost
}

//////////// Part 2

func part2(boss mob, s shop) int {
	hero := mob{100, 0, 0}
	maxcost := 0
	var gear []string
	for _, weapon := range s.weapons {
		for _, armor := range s.armors {
			for _, ring1 := range s.rings {
				for _, ring2 := range s.rings {
					if ring2.name != ring1.name { // rings are different
						hero.dmg = weapon.dmg + armor.dmg + ring1.dmg + ring2.dmg
						hero.ac = weapon.ac + armor.ac + ring1.ac + ring2.ac
						if !winFight(hero, boss) {
							cost := weapon.cost + armor.cost + ring1.cost + ring2.cost
							if cost > maxcost {
								maxcost = cost
								gear = []string{weapon.name, armor.name, ring1.name, ring2.name}
							}
						}
					}
				}
			}
		}
	}
	VP("Minimal gear:", gear)
	return maxcost
}

//////////// Common Parts code

func readShop(filename string) (s shop) {
	lines := fileToLines(filename)
	// "nothing" is a - virtual - item
	s.weapons = make([]item, 0) // weapon is mandatory
	s.armors = []item{item{"noArmor", 0, 0, 0}}
	s.rings = []item{item{"noRing1", 0, 0, 0}, item{"noRing2", 0, 0, 0}}
	section := ""
	for _, line := range lines {
		if line == "" {
			section = ""
			continue
		}
		switch section {
		case "Weapons:":
			l := strings.Fields(line)
			s.weapons = append(s.weapons, item{l[0], atoi(l[1]), atoi(l[2]), atoi(l[3])})
		case "Armor:":
			l := strings.Fields(line)
			s.armors = append(s.armors, item{l[0], atoi(l[1]), atoi(l[2]), atoi(l[3])})
		case "Rings:":
			l := strings.Fields(line)
			s.rings = append(s.rings, item{l[0] + l[1], atoi(l[2]), atoi(l[3]), atoi(l[4])})
		case "":
			l := strings.Fields(line)
			section = l[0]
		}
	}
	VP("theShop:", s)
	return
}

func readMob(lines []string) (m mob) {
	s := strings.Fields(lines[0])
	m.hp = atoi(s[2])
	s = strings.Fields(lines[1])
	m.dmg = atoi(s[1])
	s = strings.Fields(lines[2])
	m.ac = atoi(s[1])
	return
}

func winFight(hero, boss mob) bool {
	h := hero.hp
	b := boss.hp
	for {
		b = hit(hero, boss, b)
		if b <= 0 {
			return true
		}
		h = hit(boss, hero, h)
		if h <= 0 {
			return false
		}
	}
}

func hit(p1, p2 mob, hp int) int {
	loss := p1.dmg - p2.ac
	if loss < 1 {
		return hp - 1
	} else {
		return hp - loss
	}
}

//////////// Part1 functions

//////////// Part2 functions
