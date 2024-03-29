# Advent of code challenge 2017, in GO

Here are my solutions to the "Advent of code" challenge of 2017 implemented in GO (aka Golang).
See https://adventofcode.com/2017

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
- Input files for tests can also be specified in a standalone way, as files containing the problem input but named `input`*-DESCRIPTION,RESULT1,RESULT2*`.test`. (E.g. `input-negative-values,345,.test`). The optional description is freeform, and *result1* and *result2*, if present, are the expected results of the test. This avoids cluttering the Go source code, and if you ignore `*.test` files in your version control system (e.g. `.gitignore`), avoid publishing your personal input, as requested by the author of the adventofcode. 
- The examples given in the problem descriptions are used in GO unit tests `dNN_test.go` , whereas the input file is used for the high-level integration tests of `TESTALL`.

## Misc

- `days/MAKEDAY NN` is what I used to create a new day directory.
- `days/CLEANALL` prepares for a git commit: cleans dir, check missing info
- Notes per day may be found in each day directory, `days/dNN` as `README.md` files, if it warrants some explanations.
- All solutions run under one second, unless mentioned.

## More info
- https://adventofcode.com/ The "advent of code" (aka AOC) website
- https://github.com/Bogdanp/awesome-advent-of-code Bogdan Popa list of AOC-related resources and solutions
- https://github.com/topics/advent-of-code-2017 The GitHub projects tagged with `advent-of-code-2017`
- https://www.reddit.com/r/adventofcode/ The reddit to discuss AOC 

## License
Author: (c)2024 Colas Nahaboo, https://colas.nahaboo.net
License: free of use via the [MIT License](https://en.wikipedia.org/wiki/MIT_License).

## Day highlights
- **d01** to **d10** are very simple
- **d11** is interesting as it made me discover ways to work with hexagonal grids, see <https://www.redblobgames.com/grids/hexagons/>
- **d14** reuses the code of the solution in d10 to define a hash function
- **d18** uses Go goroutines and channels to run two proram instances in parallel
- **d19** From now on I tried to shed my habits gained by using slow shell scripts, and avoid using too smart data structures, instead relying on lots of small methods trying to explain the logic of the solution.
- **d23** was tricky, as it required to reverse-engineer the code of a toy assembly language, to see it was computing prime numbers, and thus guessing the expected results in a faster way. I could not see it, but ChatGPT did!
