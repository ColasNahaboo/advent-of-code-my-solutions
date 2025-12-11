# Advent of code challenge 2025, in GO, day d11

[Problem statement](https://adventofcode.com/2025/day/11)

The input dscribe a graph of less than 640 nodes with max 25 directed edges connected to each of them, with no loops.

For part 2, to find all paths going through `A -> B -> C -> D` we find all paths:
- going from A to B without going through C or D, 
- then going from B to C without going through D. No need to check for A since we do not have loops
- then going from C to D, without checking for A or B
And we add the number of paths found by the same process for `A -> C -> B -> D`

We cache the results of the pathfinding functions to avoid a combinatorial explosion.
