# Advent of code challenge 2025, in GO, day d06

[Problem statement](https://adventofcode.com/2025/day/6)

For part 2, we just rotate the input char matrix of the numbers (not the operators row) by a quarter turn to the left, and then parse naturally the numbers.

E.g, on the example input:

```
123 328  51 64 
 45 64  387 23 
  6 98  215 314
*   +   *   +
```

We rotate it into:

```
  4
431
623
   
175
581
 32
   
8
248
369
   
356
24 
1
```

We thus get the numbers of each problems, separated by a blank line, and starting by the end problem to the first one.

Note that the operators (plus or mult) are commutative, so the order of numbers in each problem is not significant.
