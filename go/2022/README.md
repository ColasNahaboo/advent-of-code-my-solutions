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
- **d16** We used a brute force approach, looking at all possible path, with some optimizations to quit exploring paths once we know we could not reach a better total pressure amount down this path than the currently found best. It works for part1, but it is too slow (one hour!) for part2, so I did like NORMIE101 on reddit (solution copied here as `alternate-solution.py`): determine first paths to all openable valves (closed valves with potential flow non-null), and only consider paths, even multi-steps going to openable valves, skipping all the zero-flow ones.
  Edit: I used a Floyd-Warshall algorithm to compute all paths between valves, and only explored these paths that lead to openable valves, greatly diminishing the exploration graph. I thus went from 75 minutes to 19 seconds. Still far from optimal, but not insane anymore. I kept the naïve implementation as reference as `d16.go-naive`.
  Edit#2: And then I recoded it (old code in `d16-old.go`) by not exploring step by step, but going all the way to the valve we want to open. Ended up at 8s running time.
- **d17** The tricky part is the second one, where we need to detect a loop in a seemingly chaotic output. It did it by hand with some usual shell magic (grep, sort, uniq, ...), and then coded the logic in a preprocessing function, called via `-p`, that outputs the constants to set in the go code for running the part 2
- **d19** Similar to d16, but this time I did not explore step by step, but each step was going all the way to build a robot. The speed gain was huge, runs in less than 0.5s. I kept the step-by-step naive implementation as `d19-naive.go`
- **d21** As computing an answer to part1 is so fast, and the possible operators being linear, for part2 add a monkey "delta" performing a substraction of the values of the two monkeys listened to by "root", and I just interpolate on values of the monkey "humn" to find one that get the value of "delta" down to zero. And as more than one humn value can give the same result (division on integers is not a bijection), we then find the smallest one from there.
- **d22** The second part is quite tricky. I ended up describing the specific geometry of mapping my input flattenned faces into a 3D cube in a table describing the connections between faces. I did not have the courage to code an automatic, input agnostic, cube mapping functioanlity for this. So my code only works for my input geometry.
