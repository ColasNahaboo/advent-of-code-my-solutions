# Advent of code challenge 2023, in GO, day d07

To simplify things, we will map hands to numbers so that we can easily sort them:
- each card is a 2-digit number: 2=02 ... 9=09, T=10, J=11, Q=12K=13, A=14
- we collate the five 2-digit strings to obtain a 10-digit number. E.g for hand `32T3K` we get `0302100313`
- we prepend the type of the hand: 
     - Five of a kind = 7
     - Four of a kind = 6
     - Full house = 5
     - Three of a kind = 4
     - Two pair = 3
     - One pair = 2
     - High card = 1

  So, for the hand  `32T3K` we get `20302100313`
- this is a bijective mapping, so we will just store internally the hands by their number, their "score".

To determine the type of the hand, we count the occurences of each card value and sort these in decreasing order, then collate them into digits of a number T, so we can easily detemine the type. For instance T=32 means a full house, T=221 means two pairs, T=5 five of a kind...

For Part2, we look at all the cards present in the hand along the joker(s), and test the score with the jocker imitating each of these cards, using the mocked value for computing the type, but the low value (1) of the joker for computing the strength. Be careful to take into account the special cases of hands with no or only jokers.
