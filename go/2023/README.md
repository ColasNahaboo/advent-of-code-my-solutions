# Advent of code challenge 2023, in GO

Here are my solutions to the "Advent of code" challenge of 2023 implemented in GO (aka Golang).
See https://adventofcode.com/2023

I am doing this to keep learning GO, so this must be considered as "student code". I am coding it to try my hand at various GO idioms, not seeking efficiency, scalability nor optimality. But feedback is very welcome.

The code is in standard GO, with some housekeeping scripts in bash.

## Usage

- Everything below is in the `days/` sub-directory, with a sub-directory `dNN/` per day.
- Day `NN` solutions are in program of source code `days/dNN/dNN.go`, and executable `days/dNN/dNN`
- Run them with input data file (with suffix `.txt` or `.test`) as argument (defaults to `input.txt`).
- The result will always be a number, alone on the last line of the output.
- They run the algorithm for Part 2 of the daily problems, unless you give the option `-1` where they will run the Part 1.
- The `-v` argument ("verbose") prints more tracing / debugging info.
- All source code is standalone, I will try to use only standard GO and standard library functions. I keep common convenience functions in the utils.go file in the same package "main" (copied from the `TEMPLATES/` directory) rather than making a proper separate packages, as I am learning as I go, and can evolve it without fear of breaking backwards compatibility with the ones used in previous days.
- Basically, all the `dNN.go` solutions consist of:
  - a `main` function
  - that read and preparse the input file via the function `ReadInput`
  - and depending on the presenc eof the command line option `-1`, calls either the `Part1` of `Part2` function to perform the calculations and return a number, the solution, that it prints.

## Testing

- Unit tests are performed via the standard GO testing system, in the optional source file `days/dNN/dNN_test.go`
- Integration tests are done by looking at the comments `// TEST: [option] input-file result` in source files and running the code with the option and input and checking the last printed line is the result. The `days/TESTALL` bash script runs all the unit and integration tests, see it for technical details. I tend to use more these integration test than the GO unit test above, as they are impervious to the major refactoring that can happen as I try many different approaches in this style of experimental programming. I rely rather on many small input samples in files `ex1.tx`, `ex2.txt`...
- **New for 2023:** input files for tests can also be specified in a standalone way, as files containing the problem input but named `input`*-DESCRIPTION,RESULT1,RESULT2*`.test`. (E.g. `input-negative-values,345,.test`). The optional description is freeform, and *result1* and *result2*, if present, are the expected results of the test. This avoids cluttering the Go source code, and if you ignore `*.test` files in your version control system (e.g. `.gitignore`), avoid publishing your personal input, as requested by the author of the adventofcode. 
- The examples given in the problem descriptions are used in GO unit tests `dNN_test.go` , whereas the input file is used for the high-level integration tests of `TESTALL`.

## Misc

- `days/MAKEDAY NN` is what I used to create a new day directory.
- `days/CLEANALL` prepares for a git commit: cleans dir, check missing info
- Notes per day may be found in each day directory, `days/dNN` as `README.md` files, if it warrants some explanations.
- All solutions run under one second, unless mentioned.

## More info
- https://adventofcode.com/ The "advent of code" (aka AOC) website
- https://github.com/Bogdanp/awesome-advent-of-code Bogdan Popa list of AOC-related resources and solutions
- https://github.com/topics/advent-of-code-2023 The GitHub projects tagged with `advent-of-code-2023`

- https://www.reddit.com/r/adventofcode/ The reddit to discuss AOC 


## License
Author: (c)2023 Colas Nahaboo, https://colas.nahaboo.net
License: free of use via the [MIT License](https://en.wikipedia.org/wiki/MIT_License).

## Day highlights

- **d01** Contained a subtle trick: it needed to find words in the strings that can overlap: E.g. find `eight` and `three` in the string `eighthree`
- **d04** Scaling for part2 needed to factorize cards so not to process each one individually
- **d05** Scaling for part2 needed to regroup number into ranges of numbers being mapped in the same way, and only work on the (few) ranges, and not the (many) numbers. The difficulty is that the ranges are different on each mapping, and you must use the intersections of them.
- **d06** Modelizing the problem shows that it amounts to finding the integer that are abscisses of the part of a parabola above a threshold. Thus solving a 2nd order polynomial equation.
- **d07** My solution was to map poker hands into numbers that I could then easily sort and compare to score the hands
- **d08** A problem that can only be solved easily by strong hypotheses on the input. In this case that values loop in a clean way.
- **d10** A topological problem, where you must find points inside a loop, my solution is tracing a ray and counting the crossings of the loop. Point is inside if this number is odd.
- **d11** Simple problem, but too huge to be solved by handling space as a grid, we must use the mere list of galaxies coordinates instead.
- **d12** I had a very hard time to debug this. The algorithm seemed simple enough to me, but I kept having bugs that I could not identify.
- **d14** A classic: detecting a loop in a series of results to avoid computing all of them. I also made a library `Scalarray.go` to help managing these 2D grids that are so often used in these kind of problems.
- **d16** The trick was to avoid looping when following the path of the light beam.
- **d17** A typical problem, where we must find the shortest path in states of "things" on a 2D board. We combine for the solution the use of `scalarray.go` (for the 2D board) and `astar.go` (for the shortest path). We also provide an alternate implementation of part2, named part3 and callable with the `-3` command line option that use a 3D scalar array instead of a hashtable map to get IDs of states.
- **d18** asks to count cells on a grid, but its becoms so huge in part2 that we must aggregate the cells in virtual mega tiles and work on them. The fun in the text was that the 1st part mentioned colors but with played no issue at all ultimately, and where a way to "hide" into the common input file what amounted to a different input for part 2.
- **d20** part2 could be only solved by looking at the possible paths in reverse from the goal and seeing that it needed that 4 nodes be in low state at some point, but these 3 nodes were entering low cyclycally with different period. The solution was thus the Least Common Multiple (LCM) of these cycles.
- **d21** was very hard, but not in the programming sense. All the difficulty was to understand how the input was a very specific case and designing and ad hoc solution.
- **d22**, on the opposite, was a very pleasant exercise, a 3D tetris. 
- **d23** a straightforward path exploration
