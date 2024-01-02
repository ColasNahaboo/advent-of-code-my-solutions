# Advent of code challenge 2023, in GO, day d19

The part 2 difficulty is scaling up. Trying all the combinations would mean trying 4000^4, more than 10^18.

At firts I tried the approach I used in the day before, to divide the 4D solution space into 4D tiles, but it was still big (4 billion tiles), running in more than 2 minutes. (I saved this approach as a "part 3"). So I just explored the tree, auto-creating new tiles on the go, splitting at each condition comparison. I Then only had to look at 500 tiles altogether...
