# Advent of code challenge 2017, in GO, day d21

For this problem, I tried to keep the data structures as close to the human-readable input as possible, and put the logic in the functions/met6hods handling them instead of optimizing the data structure.

I just pre-generated all the possible variants of a pattern, so that rule matching would be simple string comparisons. However, I kind of made an ad hoc hashtable in that I grouped the matching patterns (the conditions part of the rule) by string lenghth and number of lit pixels, to limit the number of matches to perform.
