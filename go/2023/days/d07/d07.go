// Adventofcode 2023, d07, in go. Arguments:
// -1: use solve part 1 of the problem, default is the second one
// the input file name: default: input,RESULT1,RESULT2.test
// TEST: -1 example 6440
// TEST: example 5905
// And any file named input-DESCRIPTION,RESULT1,RESULT2.test containing an input

package main

import (
	"flag"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

var verbose bool

// This is over-engineered, as I was anticipating part2 that would have scaled
// hugely. In retrospect, just storing cards by their score (field "hand")
// and computing other fields (e.g cards) on demand would have been fast enough.

type Bid struct {
	id int						// input position, starts at 1 (for debug)
	cards string				// the 5-letter representation of the hand
	values [5]int				// their values (1..14)
	hand int					// a number, its score
	bid int						// the bid value itself
	rank int					// the rank in the set of hands of the input
	// the following fields are only used in part2, not part1
	jhand int					// part2: max possible score with joker
	jrank int					// part2: the ranking when J is a joker
	mockcards string			// part2: the mock hand
}
// where is the type digit in a hand score
const typeOffset = 10000000000	// 10 to power twice the cards in a hand (5x2=10)
const joker = 11

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

func part1(lines []string) (winnings int) {
	bids := parse(lines)
	VP(bids)
	for _, b := range bids {
		winnings += b.bid * b.rank
		VPf("  Hand %s: %d * %d = %d\n", b.cards, b.bid, b.rank, b.rank * b.bid)
	}
	return
}

//////////// Part 2
func part2(lines []string) (winnings int) {
	bids := parse(lines)
	jokerize(bids)
	VP(bids)
	for _, b := range bids {
		winnings += b.bid * b.jrank
		VPf("  Hand %s (%s): %d * %d = %d\n", b.cards, b.mockcards, b.bid, b.jrank, b.jrank * b.bid)
	}
	return
}

//////////// Common Parts code

func parse(lines []string) (bids []Bid) {
	re := regexp.MustCompile("[0-9TJQKA]+")
	for id, line := range lines {
		tokens := re.FindAllString(line, -1)
		if tokens == nil {		// skip empty lines
			continue
		}
		hand := a2h(tokens[0])
		bid := atoi(tokens[1])
		bids = append(bids, Bid{id: id+1, cards: tokens[0], hand: hand, bid: bid, values: a2hi(tokens[0])})
	}
	// rank the hands. Sort them by increasing score
	sort.Slice(bids, func(i, j int) bool {
		return bids[i].hand < bids[j].hand
	})
	// and then assign ranks
	for i, _ := range bids {
		bids[i].rank = i + 1
	}
	return
}

// 5-letter (Ascii) hand representation to its 11-digits score number. See README
func a2h(s string) (hand int) {
	numof := make([]int,15,15)	// number of occurences for each value
	for _, d := range s {
		switch d {
		case 'A': hand = hand * 100 + 14; numof[14]++
		case 'K': hand = hand * 100 + 13; numof[13]++
		case 'Q': hand = hand * 100 + 12; numof[12]++
		case 'J': hand = hand * 100 + 11; numof[11]++
		case 'T': hand = hand * 100 + 10; numof[10]++
		default: hand = hand * 100 + int(d - '0'); numof[d - '0']++
		}
	}
	hand += handType(numof) * typeOffset
	return
}

// computes the score, but with joker cards counting as mock for type
// and as lowest (1) for strength
func scoreJoker(b Bid, mock int) (hand int) {
	numof := make([]int,15,15)	// number of occurences for each value
	hasJoker := false
	for _, v := range b.values {
		if v == joker {
			hand = hand * 100 + 1; numof[mock]++
			hasJoker = true
		} else {
			hand = hand * 100 + v; numof[v]++
		}
	}
	if hasJoker {
		hand += handType(numof) * typeOffset
	} else {
		hand = 0
	}
	return
}
	
// 5 nums hand representation
func a2hi(s string) (values [5]int) {
	for i, d := range s {
		switch d {
		case 'A': values[i] = 14
		case 'K': values[i] = 13
		case 'Q': values[i] = 12
		case 'J': values[i] = 11
		case 'T': values[i] = 10
		default: values[i] = int(d - '0')
		}
	}
	return
}

// displays a hand as its external 5-letter representation
func h2a(hand int) (s string) {
	// discard type, then process each 5 "digits"
	for n := hand % typeOffset; n > 0; n /= 100 {
		switch d := n % 100; d {
		case 10: s = "T" + s
		case 11: s = "J" + s
		case 12: s = "Q" + s
		case 13: s = "K" + s
		case 14: s = "A" + s
		default: s = itoa(d) + s
		}
	}
	return
}

// displays a hand as its external 5-letter representation
func i2card(i int) (s string) {
	switch i {
	case 10: return "T"
	case 11: return "J"
	case 12: return "Q"
	case 13: return "K"
	case 14: return "A"
	default: return itoa(i)
	}
	return
}

func handType(numof []int) int {
	// sort occurences greatest first
	sort.Slice(numof, func(i, j int) bool {
		return numof[i] > numof[j]
	})
	// keep only non-zero value and collate them
	t := 0
	for _, i := range numof {
		if i > 0 {
			t = t * 10 + i
		}
	}
	// now we can easily determine the type
	switch t {
	case 5: return 7
	case 41: return 6
	case 32: return 5
	case 311: return 4
	case 221: return 3
	case 2111: return 2
	case 11111: return 1
	}
	panic("Invalid occurences: " + itoa(t))
}

//////////// Part1 functions

//////////// Part2 functions

// process a set of hands to fill the joker-mode specific fields jhand & jrank
func jokerize(bids []Bid) {
	// compute jhands. Warning, b is a copy, read it but write in bids[i]
	for i, b := range bids {
		jokerValue := 0
		if b.hand == 71111111111 { // JJJJJ all jokers hand, special case
			bids[i].jhand = 70101010101
			bids[i].mockcards = "AAAAA"
			continue
		}
		if ! strings.Contains(b.cards, "J") { // No joker, special case
			bids[i].jhand = b.hand
			bids[i].mockcards = b.cards
			continue
		}			
		cards := cardsOf(b.hand) // the other cards that J can imitate
		jhand := 0
		for _, mock := range cards { // test score if J mocks card "mock"
			score := scoreJoker(b, mock)
			if score > jhand {
				jhand = score
				jokerValue = mock
			}
		}
		bids[i].jhand = jhand
		bids[i].mockcards = strings.ReplaceAll(b.cards, "J", i2card(jokerValue))
		VPf("  For %s, joker = %s, 	jhand=%d\n", b.cards, i2card(jokerValue), jhand)
	}
	//  then, jrank them
	sort.Slice(bids, func(i, j int) bool {
		return bids[i].jhand < bids[j].jhand
	})
	// and then assign ranks
	for i, _ := range bids {
		bids[i].jrank = i + 1
	}
}

// list all unique non-joker cards in hand: T55J5 ==> [10 5]
func cardsOf(hand int) (cards []int) {
	for i := typeOffset / 100; i > 0; i /= 100 {
		card := (hand / i) % 100
		if card != joker && indexOfInt(cards, card) == -1 {
			cards = append(cards, card)
		}
	}
	return
}

