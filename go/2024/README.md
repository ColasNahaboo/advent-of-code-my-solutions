# Advent of code challenge 2024, in GO

Here are my solutions to the "Advent of code" challenge of 2024 implemented in GO (aka Golang).
See https://adventofcode.com/2024

I am doing this to keep learning GO, so this must be considered as "amateur code". I am coding it to try my hand at various GO idioms, not seeking efficiency, scalability nor optimality. But feedback is very welcome.

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
- https://github.com/topics/advent-of-code-2024 The GitHub projects tagged with `advent-of-code-2024`
- https://www.reddit.com/r/adventofcode/ The reddit to discuss AOC 

## License
Author: (c)2024 Colas Nahaboo, https://colas.nahaboo.net
License: free of use via the [MIT License](https://en.wikipedia.org/wiki/MIT_License).

## Day highlights
- **d04** Tried to use methods on generics, and immediately closed the Pandora can of worms I opened!
- **d06** Different implementation to test speeds for part2: 4 arrays bools or array of ints or bytes used as 4-bit bitfields. The latter is faster.
- **d14** How to recognize an xmas tree in the output? I suppose it is inside a box, and I look for horizontal lines
- **d23** I used the Bron Kerbosh algorithm to find the biggest clique in a graph. I coded after the deadline also alternative using the bitset package, and also the Bron-Kerbosch implementation in the graph/topo package, out of curiosity
- **d24** I solved part2 by hand to get to the finish line, but I plan to design a proper solution later.
