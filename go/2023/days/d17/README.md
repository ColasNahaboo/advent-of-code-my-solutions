# Advent of code challenge 2023, in GO, day d17

A typical problem, where we must find the shortest path in states of "things" on a 2D board. We combine for the solution the use of `scalarray.go` (for the 2D board) and `astar.go` (for the shortest path). 

## Part 3

We also provide an alternate implementation of part2, named part3 and callable with the `-3` command line option that use a 3D scalar array instead of a hashtable map to get IDs of states. 
A state ID is thus its position (number) in the array of Scalarray3D
We do not have an actual State type, we map the fields virtually on the 3D coords: pos is x, dir is y, steps is z

Note that we do not need to actually instanciate the scalar array field of id3d, we only use its coordinates/position conversion methods, as we do not have to store additional data to Nodes.

The Part3 implementation is thus 25% faster than part2, and takes less space as there is no actual data structure to do the mapping state -> ID.

## Part 4

Another alternative implementation of Part2, "Part 4" called with `-4`, also avoids the mapping state -> ID via a Go map, and instead use for Node the State type itself (a struct of 3 ints). The code is then simpler and more readable, with only a speed penaly of less than 4% compared to Part 2.
