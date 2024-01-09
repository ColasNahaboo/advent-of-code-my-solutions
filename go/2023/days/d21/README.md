# Advent of code challenge 2023, in GO, day d21

Part1 was trivial, but Part 2 was quite hard, and not very interesting as a programming exercice. It consited in understanding howthe input was hyper-specific and design and ad hoc solution that worked only in this specific case.
- the set of reachable plots spread out in a diamond pattern undisturbed, ad there are no rocks on lines going straight out the starting place, either horizonatlly, vertically, or in diagonals
- once filled, the "tiles" ocillate between two patterns
- the number of steps is very specific, a multiple of the input "tile" side (131) plus half (65)
- the the number is expressed in a quadratic functions, that we have to solve.
