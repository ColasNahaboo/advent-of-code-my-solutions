// Adventofcode 2015, day22, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input.txt
// TEST: -1 -e exemple1 226
// TEST: -1 -e exemple2 641
// TEST: -1 input 1824
// TEST: input
package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
)

type (
	// ID is the path of the instance through all possible choices (spells)
	ID int
	// World are "instances" in which fights happen, forked on each player choice
	World struct {
		id ID // the ID of this instance, series of digits, 9=root (debug)
		// 923 = instance where hero casted spell #2 then #3, printed /23
		time    int // the time (global)
		mana    int // consumed mana for this instance
		hero    Hero
		boss    Boss
		effects [SpellCount]Effect // the currently active effects
		curmana *int               // Optimisation: Current minimum mana consumed currently found.
		// fightRound will stop exploring a branch as soon as reached
		// must point to a variable in the caller space of fightRound
	}
	// Universe is things that never change... constants
	Universe struct {
		book [SpellCount]Spell // spellbook: the possible spells
		hard bool              // hard mode: -1 hp before each player turn
	}
	// Boss is a simple mob, without ac or spells
	Boss struct {
		hp  int
		dmg int
	}
	// Hero is the player instance, with its fields being the current values
	Hero struct {
		hp   int
		mana int
		ac   int
	}
	// Effect is a dynamic effect applied to hero. Instancied from a Spell
	Effect struct {
		spellid int // index in the spellbook
		timer   int // active while timer >0. decreased each turn
		dot     int // Damage Over Time
		aot     int // Armor Over Time
		mot     int // Mana Over Time
	}
	// Spell is a Template/Prototype/Class to create Effects from a book
	Spell struct {
		name     string
		cost     int
		dmg      int // instant actions
		ac       int
		heal     int
		duration int // "OT" (Over Time) actions
		dot      int
		aot      int
		mot      int
		// hot    int // not used here, but could be used in the general case
	}
)

const (
	MagicMissile = iota
	Drain
	Shield
	Poison
	Recharge
	SpellCount // handy hack to have the size of arrays for declarations
)

var verbose bool

// an easier to spot maxint in debug than 9223372036854775807 (^uint(0) >> 1)
const maxint = 8888888888888888888

func main() {
	partOne := flag.Bool("1", false, "run exercise part1, (default: part2)")
	verboseFlag := flag.Bool("v", false, "verbose: print extra info")
	exampleFlag := flag.Bool("e", false, "examples: Use 10 hp & 250 mana for hero instead of 50/100")
	smallFlag := flag.Bool("s", false, "small hero: Use 1 hp & 250 mana for hero instead of 50/100")
	flag.Parse()
	verbose = *verboseFlag
	infile := "input.txt"
	u := Universe{[SpellCount]Spell{
		Spell{"MagicMissile", 53, 4, 0, 0, 0, 0, 0, 0}, // 0
		Spell{"Drain", 73, 2, 0, 2, 0, 0, 0, 0},        // 1
		Spell{"Shield", 113, 0, 0, 0, 6, 0, 7, 0},      // 2
		Spell{"Poison", 173, 0, 0, 0, 6, 3, 0, 0},      // 3
		Spell{"Recharge", 229, 0, 0, 0, 5, 0, 0, 101},  // 4
	},
		false,
	}
	w := World{id: NewID(), effects: [SpellCount]Effect{
		Effect{MagicMissile, 0, 0, 0, 0}, // 0
		Effect{Drain, 0, 0, 0, 0},        // 1
		Effect{Shield, 0, 0, 7, 0},       // 2
		Effect{Poison, 0, 3, 0, 0},       // 3
		Effect{Recharge, 0, 0, 0, 101},   // 4
	},
	}

	switch flag.NArg() {
	case 1: // built-in hero + input file describing a boss
		infile = flag.Arg(0)
		fallthrough
	case 0: // built-in hero + boss file = input.txt
		lines := fileToLines(infile)
		w.boss = readBoss(lines)
	case 4: // hero HP & mana + boss HP & dmg as arguments
		w.hero.hp = atoi(flag.Arg(0))
		w.hero.mana = atoi(flag.Arg(1))
		w.boss.hp = atoi(flag.Arg(2))
		w.boss.dmg = atoi(flag.Arg(3))
	case 2: // boss HP & dmg as arguments
		w.boss.hp = atoi(flag.Arg(0))
		w.boss.dmg = atoi(flag.Arg(1))
	default:
		log.Fatal("Bad number of arguments")
	}
	if *exampleFlag {
		w.hero = Hero{hp: 10, mana: 250}
	} else if *smallFlag {
		w.hero = Hero{hp: 1, mana: 125}
	} else {
		w.hero = Hero{hp: 50, mana: 500}
	}

	var result int
	if *partOne {
		VP("Running Part1")
		result = part1(w, u)
	} else {
		VP("Running Part2")
		result = part2(w, u)
	}
	fmt.Println(result)
}

//////////// Part 1

// Optimisation: Current minimum mana consumed currently found.
// fightRound will stop exploring a branch as soon as reached

func part1(w World, u Universe) int {
	curmana := maxint
	w.curmana = &curmana
	return fightRound(&w, u)
}

// return the minimum mana used for a winning fight by processing
// all possible actions in a round (hero turn + boss turn) and recursing
// tick is the turn number (a round runs for 2 ticks)

func fightRound(w *World, u Universe) int {
	// then examine all possible spell to cast
	minmana := maxint
	for spellid, spell := range u.book {
		if w.hero.mana >= spell.cost && w.effects[spellid].timer <= 1 && (w.mana+spell.cost) < *w.curmana {
			// OK, we can cast this one. Fork/Clone an new instance to explore it
			nw := *w
			nw.id = nw.id*10 + ID(spellid) // add spellid as branch

			// PLAYER TURN
			nw.time++
			VPf("Instancing %v at time %v\n", nw.id, nw.time)
			// hero turn in hard mode: -1 hp
			if u.hard {
				nw.hero.hp--
				if nw.hero.hp <= 0 {
					VPf("Loses hard: Hero: %v, Boss: %v\n", nw.hero, nw.boss)
				}
			}
			// apply existing effects
			runEffects(&nw)
			if nw.boss.hp <= 0 { // Win!
				VPf("#Win Heffe [%v](%v), mana: %v\n", nw.id, nw.time, nw.mana)
				checkMana(&nw, &minmana, nw.mana)
				continue
			}
			nw.hero.mana -= spell.cost // consume the spell mana
			nw.mana += spell.cost
			runInstant(&nw, spell)         // run instant actions of the spell
			addEffect(&nw, spellid, spell) // schedule OverTime actions of the spell
			if nw.boss.hp <= 0 {           // Win!
				VPf("#Win Hinst [%v](%v), mana: %v\n", nw.id, nw.time, nw.mana)
				checkMana(&nw, &minmana, nw.mana)
				continue
			} else {

				// BOSS TURN
				nw.time++
				runEffects(&nw)
				if nw.boss.hp <= 0 { // Win!
					VPf("#Win Beffe [%v](%v), mana: %v\n", nw.id, nw.time, nw.mana)
					checkMana(&nw, &minmana, nw.mana)
					continue
				}
				hit := nw.boss.dmg - nw.hero.ac
				if hit < 1 {
					hit = 1
				}
				nw.hero.hp -= hit
				if nw.hero.hp > 0 {
					// still alive? next round!
					mana := fightRound(&nw, u)
					checkMana(&nw, &minmana, mana)
				} else {
					VPf("Loses: Hero: %v, Boss: %v\n", nw.hero, nw.boss)
				}
			}
		}
	}
	if minmana == maxint {
		VPf("No solution found for %v\n", w.id)
	} else {
		VPf("###Best solution for %v: %v mana\n", w.id, minmana)
	}
	return minmana
}

func checkMana(nw *World, mmp *int, mana int) {
	if mana < *mmp {
		*mmp = mana
		if *mmp < *nw.curmana {
			*nw.curmana = *mmp
		}
	}
}

func runInstant(nw *World, spell Spell) {
	nw.hero.hp += spell.heal
	nw.boss.hp -= spell.dmg
	nw.hero.ac += spell.ac
}

func runEffects(w *World) {
	for i := 0; i < SpellCount; i++ {
		if w.effects[i].timer > 0 {
			w.boss.hp -= w.effects[i].dot
			w.hero.mana += w.effects[i].mot
			if w.effects[i].aot != 0 && w.effects[i].timer == 1 {
				w.hero.ac -= w.effects[i].aot // remove ac on last tick
			}
			w.effects[i].timer--
		}
	}
}

func addEffect(w *World, spellid int, spell Spell) {
	if spell.duration > 0 && w.effects[spellid].timer == 0 {
		w.effects[spellid].timer = spell.duration
		w.hero.ac += spell.aot // will be removed at end of timer
	}
}

//////////// Part 2

func part2(w World, u Universe) int {
	u.hard = true
	curmana := maxint
	w.curmana = &curmana
	return fightRound(&w, u)
}

//////////// Common Parts code

func readBoss(lines []string) (m Boss) {
	s := strings.Fields(lines[0])
	m.hp = atoi(s[2])
	s = strings.Fields(lines[1])
	m.dmg = atoi(s[1])
	return
}

// IDs start with 9, we replace by / for nicer printing.
func NewID() ID {
	return ID(9)
}
func (id ID) String() string {
	return "/" + itoa(int(id))[1:]
}

//////////// Part1 functions

//////////// Part2 functions
