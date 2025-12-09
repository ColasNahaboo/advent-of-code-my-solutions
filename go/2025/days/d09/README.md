# Advent of code challenge 2025, in GO, day d09

[Problem statement](https://adventofcode.com/2025/day/9)

The coordinates in x and y of the red tiles are between 1900 and 99000 exclusive.
So a map of the area would have more than 6 billion places, too big for naive implementation., 
However, the x take only less than 500 distinct values, and the y 250,
so we can normalize into a 500x250 grid of 125000 cells max

For the example:
```
7,1
11,1
11,7
9,7
9,5
2,5
2,3
7,3
```
We see that x can only take the four values `2, 7, 9, 11`, and y `1, 3, 5, 7`.
We then map these actual values to "normalized" ones, their index into the array of possible values.
Thus a point `{1, 2}` in our normalized coordinates would actually be `{7, 5}`

Note that we keep the empty spaces between the normalized values: if we have x taking only values `1, 3, 4, 7`, we actually take the normalized values to be `1, 2, 3, 4, 5, 7`, the values 2 and 5 avoid collapsing the empty regions.

The example map:
```
..............
.......#...#..
..............
..#....#......
..............
..#......#....
..............
.........#.#..
..............
```
thus normalizes to:
```
..#...#
.......
#.#....
.......
#...#..
.......
....#.#
```
And, with the edges traced:
```
..#####
..#...#
###...#
#.....#
#####.#
....#.#
....###
```
And, filled:
```
..#####
..#ooo#
###ooo#
#ooooo#
#####o#
....#o#
....###
```
reducing it to a size managed simply

We also use the heuristic to find an inside point start filling the inside by propagation by looking for the first empty sapce after an edge on the second row.
