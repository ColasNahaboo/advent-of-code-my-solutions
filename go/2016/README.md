# Advent of code challenge 2016, in GO

Here are my solutions to the "Advent of code" challenge of 2016 implemented in GO (aka Golang).
See https://adventofcode.com/2016

I am doing this to learn GO, so this must be considered as "student code". I am coding it to try my hand at various GO idioms, not seeking efficiency, scalability nor optimality. But feedback is very welcome.

The code is in standard GO, with some housekeeping scripts in bash.

**Usage:**

- Everything below is in the `days/` sub-directory, with a sub-directory `dNN/` per day.
- Day `NN` solutions are in program of source code `days/dNN/dNN.go`, and executable `days/dNN/dNN`
- Run them with input data file (with suffix `.txt`) as argument (defaults to `input.txt`).
- The result will always be a number, alone on the last line of the output.
- They run the algorithm for Part 2 of the daily problems, unless you give the option `-1` where they will run the Part 1.
- The `-v` argument ("verbose") prints more tracing / debugging info.
- All source code is standalone, I will try to use only standard GO and standard library functions. I keep common convenience functions in the utils.go file in the same package "main" (copied from the `TEMPLATES/` directory) rather than making a proper separate packages, as I am learning as I go, and can evolve it without fear of breaking backwards compatibility with the ones used in previous days.
- Basically, all the `dNN.go` solutions consist of:
  - a `main` function
  - that read and preparse the input file via the function `ReadInput`
  - and depending on the presenc eof the command line option `-1`, calls either the `Part1` of `Part2` function to perform the calculations and return a number, the solution, that it prints.

**Testing:**

- Unit tests are performed via the standard GO testing system, in the source file `days/dNN/dNN_test.go`
- Integration tests are done by looking at the comments `// TEST: [option] input-file result` in source files and running the code with the option and input and checking the last printed line is the result. The `days/TESTALL` bash script runs all the unit and integration tests, see it for technical details.
- The examples given in the problem descriptions are used in GO unit tests `dNN_test.go` , whereas the input file is used for the high-level integration tests of `TESTALL`.

**Misc:**

- `days/MAKEDAY NN` is what I used to create a new day directory.
- `days/CLEANALL` prepares for a git commit: cleans dir, check missing info

## Notes per day

Note: all solutions run under one second, unless mentioned.

Conce my solutions worked, I compared with them (and others in the [reddit megathread](https://www.reddit.com/r/adventofcode/wiki/solution_megathreads#wiki_december_2016) for inspiration.

- **d01** to **d10** Simple problems, nothing remarkable
- **d11** first hard problem. I used an A-star algorithm to find the shortest path in all the possible moves. First naively, archived in `d11-old1.go1` than ran in 1mn, and then 30s by adding a cache. I then implemented the optimisation to consider the actual names of metals as not important, and the only the configurations of the pairs (generator, microchip) were considered, without diffentiating by metal. I used an ID generated in a way to be the same for equivalent states. I took then the opportunity to test two A-star Golang implementation: the [fzipp/astar](https://pkg.go.dev/github.com/fzipp/astar) (default) or the one in [gonum](https://pkg.go.dev/gonum.org/v1/gonum/graph/path#AStar) (called with option `-3`). The first one is easier to use, the second is much harder to use as the doc is a awful mess of spaghetti with very few examples, and requires to build the graphe beforehand, but it seems faster once the graph is build... but much slower overall.
- **d12** is quite simple. I just tested two approaches: interpreting the code text lines, it ran in 42 seconds, but by pre-compiling it and executing as a virtual machine, it ran in less than 0.1 seconds.
- **d13** is relatively simple, except that the part2 problem text is ambiguous. I interpreted it as: try all the locations on the floor that have a shortes path to them of less than 50 steps. For my input, the shortest path to `(31,39)` was `82`, and by tring all the 1352 locations within 50 steps of `(1,1)` for the Manhattan distance (in a triangle) that had an actual shorter path of 50 or less from `(1,1)` I found 615 positions. But this was not the expected answer! I thus kept the code for reference and made it callable via a new `-3` command line option, and implemented what seemed the interpretation of the puzzle author by looking at the solutions: it was, when looking for the shortest path to `(31,39)`, all the locations examined during this specific search of the shorted path to this single location... and none other! This yield the expected result of `138` in my case. But I still think that the problem was wrongly formulated!
- **d14** to **d18** were simple.
- **d19** is easy for the part 1, but playing the given algorithm steps naively is much too slow for part 2. I thus generated the numbers for the first 5000 values of the circle length, and could find a pattern: at regular intervals, for a length `l` the result was 1, and starting from this length, incrementing the length incremented the result by 1 up to a length being `(l-1)*2`, and afterwards incremented by 2. Combined with the fact that the next length with result 1 was the previous one multiplied by 3 minus 2, it was easy to get the result extremely fast, without any calculation of the algorithm itself.
- **d22** The part2 is intended to be solved by looking at the input data as layed out on a map, and realizing that it is a very specific "taquin" (sliding tiles) game that can be solved even by hand. However, I took the opportunity (and the time) to tackle the general solution but finding the shortest path in the graph of possible sliding tiles moves with the A Star algorithm, as I thought that this approach could solve a lot of future exercices. I started with the [fzpp/astar](https://github.com/fzipp/astar) GitHub package but modified it to be both simpler to use, and flexible enough to be used in dynamically created graphs. For now it is just a `astar.go` source file, but I make publish it as a standalone package on GitHub in the future.
