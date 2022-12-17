# Advent of code challenge 2022, in GO

Here are my solutions to the "Advent of code" challenge of 2022 implemented in GO (aka Golang).
See https://adventofcode.com/2022

I am doing this to keep learning GO, so this must be considered as "student code". I am coding it to try my hand at various GO idioms, not seeking efficiency, scalability nor optimality. But feedback is very welcome.

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

- Unit tests are performed via the standard GO testing system, in the optional source file `days/dNN/dNN_test.go`
- Integration tests are done by looking at the comments `// TEST: [option] input-file result` in source files and running the code with the option and input and checking the last printed line is the result. The `days/TESTALL` bash script runs all the unit and integration tests, see it for technical details. I tend to use more these integration test than the GO unit test above, as they are impervious to the major refactoring that can happen as I try many different approaches in this style of experimental programming. I rely rather on many small input samples in files `ex1.tx`, `ex2.txt`...
- The examples given in the problem descriptions are used in GO unit tests `dNN_test.go` , whereas the input file is used for the high-level integration tests of `TESTALL`.

**Misc:**

- `days/MAKEDAY NN` is what I used to create a new day directory.
- `days/CLEANALL` prepares for a git commit: cleans dir, check missing info

## Notes per day

Note: all solutions run under one second, unless mentioned.



- **d01** to **d11** Simple problems, nothing to say. I made use of the new `sort.Slice` functionality of Go 1.18, very cool!
- **d12** First difficulty, finding shortest path in a graph. I used a simple BFS search, using a FIFO by the package `github.com/gammazero/deque`, quite simple in Go after having done it in bash...
- **d15** The second part is interesting, as the huge size of the grid (4 million sides) makes looking at each position too slow. So we need to work with intervals of positions.
- **d16** We se a brute force approach, looking at all possible path, with some optimisations to quit exploring paths once we know we could not reach a better total pressure amount downthis path than the currently found best. It works for part1, but it is too slow (one hour!) for part2, I will have to do like NORMIE101 on reddit (solution copied here as `alternate-solution.py`)and determine first paths to all openable valves, and only consider paths, even multi-steps going to openable valves, skipping all the zero-flow ones.
- **d17** The tricky part is the second one, where we need to detect a loop in a seemingly chaotic output. It did it by hand with some usual shell magic (grep, sort, uniq, ...), and then coded the logic in a preprocessing function, called via `-p`, that outputs the constants to set in the go code for running the part 2
