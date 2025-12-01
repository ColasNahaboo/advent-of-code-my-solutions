# Advent of code challenge 2025, in GO, day d01

[Problem statement](https://adventofcode.com/2025/day/1)

Basically, the position of the dial is the one we get modulo 100 after performing the move.
So, a simple formula would be:

`new-pos = (old-pos + move) % 100`

E.g, for R42 while the dial is at 76

`18 = (76 + 42) % 100`

The difficulty is for negative moves, e.g the first example move of `L68`. Trying to find the correct arithmetic formula from scratch is error prone, as modulo keeps the sign:

`-42 % 100` yields `-42` where we would have liked to get `58`.

The sanest solution is thus not to use the classic modulo operator (`%`), but a **Positive Modulo** function (also called true modulo or Euclidian modulo), that we call `pmod`.

Thus `new-pos = pmod(old-pos + move, 100)` works in all cases.


## Part2

Positive rotations, are simple, for negative (counterclockwise), we flip the picture horizontally:

```
      --- h flip --> 
 
     0             0 
  7     1       1     7 
 6       2     2       6 
  5     3       3     5 
     4   *     *   4 
```

Thus a L42 becomes a R42 on the mirrored dial, where numbers on the dial are `pmod(100 - n)`, and we can use the simple computations for positive rots.

