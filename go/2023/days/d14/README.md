# Advent of code challenge 2023, in GO, day d14

In part2, obviously it was supposed we did not perform the billion cycles, but detect when the cycles were looping. I did it by brute force, storing all the board positions and comparing with all of them on each cycle. Go is fast enough to perform this in less than 1/10th of a second anyways.
