# Advent of code challenge 2015, in GO

Here are my solutions to the "Advent of code" challenge of 2015 implemented in GO (aka Golang).
See https://adventofcode.com/2015

I am doing this to learn GO, so this must be considered as "student code". I am coding it to try my hand at various GO idioms, not seeking efficiency, scalability nor optimality. But feedback is very welcome.

The code is in standard GO, with some housekeeping scripts in bash.

**Usage:**

- Everything below is in the `days/` sub-directory, with a sub-directory `dayNN/` per day.
- Day `NN` solutions are in program of source code `days/dayNN/dayNN.go`, and executable `days/dayNN/dayNN`
- Run them with input data file (with suffix `.txt`) as argument (defaults to `input.txt`).
- The result will always be a number, alone on the last line of the output.
- They run the algorithm for Part 2 of the daily problems, unless you give the option `-1` where they will run the Part 1.
- The `-v` argument ("verbose") prints more tracing / debugging info.
- All source code is standalone, I will try to use only standard GO and standard library functions. I keep common convenience functions in the utils.go file in the same package "main" (copied from the `TEMPLATES/` directory) rather than making a proper separate packages, as I am learning as I go, and can evolve it without fear of breaking backwards compatibility with the ones used in previous days.
- Basically, all the `dayNN.go` solutions consist of:
  - a `main` function
  - that read and preparse the input file via the function `ReadInput`
  - and depending on the presenc eof the command line option `-1`, calls either the `Part1` of `Part2` function to perform the calculations and return a number, the solution, that it prints.

**Testing:**

- Unit tests are performed via the standard GO testing system, in the source file `days/dayNN/dayNN_test.go`
- Integration tests are done by looking at the comments `// TEST: [option] input-file result` in source files and running the code with the option and input and checking the last printed line is the result. The `days/TESTALL` bash script runs all the unit and integration tests, see it for technical details.
- The examples given in the problem descriptions are used in GO unit tests `dayNN_test.go` , whereas the input file is used for the high-level integration tests of `TESTALL`.

**Misc:**

- `days/MAKEDAY NN` is what I used to create a new day directory.
- `days/CLEANALL` prepares for a git commit: cleans dir, check missing info

## Notes per day
Note: all solutions run under one second, unless mentioned.

For debugging, I used sometings solutions by [devjobe](https://github.com/devjobe/advent-of-code-2015-golang) and [schwern](https://github.com/schwern/adventofcode.go) to generate more test data. And once my solutions worked, compared with them (and others in the [reddit megathread](https://www.reddit.com/r/adventofcode/wiki/solution_megathreads#wiki_december_2015) for inspiration.

- **Day01** Starting simple, deciding on the directory and file structure, experimenting with testing and debugging via dlv
- **Day02** Working with regexps.
- **Day03** Working with maps.
- **Day04** Using package md5.
- **Day05** Compensing the limits of GO RE2 regexp standard.
- **Day06** Working with 2-dimensional arrays.
- **Day07** Working with function pointers and kind of closures.
- **Day08** Simple.
- **Day09** A brute-force Traveling Salesman Problem. I used a nice hack to generate all the routes (permutations), found at https://golangbyexample.com/all-permutations-string-golang/
- **Day10** Interesting discovery of the ["Look and say"](https://en.wikipedia.org/wiki/Look-and-say_sequence) analysis by John Conway, with a nice [video](https://www.youtube.com/watch?v=ea7lJkEhytA) of him detailing it. A naive strings-based implementation was much too slow, but using arrays of integers proved very fast.
- **Day11** Converted strings to work on arrays of integers.
- **Day12** Working with JSON, interfaces and type switches.
- **Day13** Same type of problem (and solution) than Day09, so I reused the code.
- **Day14** Use of index + pointers to be able to modify structs in ranges. See https://stackoverflow.com/questions/20185511/range-references-instead-values
  ```
  for i := range array {
    e := &array[i]
    e.field = "foo"
  }
  ```
- **Day15** Generate all ways of dispatching items in boxes. Put the convenience functions in a separate file, utils.go, but still in the main package. Passes `golint`. Use of variadic functions.
- **Day16** Working with successive sub-matches of a regexp in a line
- **Day17** We use a bitset (bits sets of a uint64 number) to represent a combination of containers.
- **Day18** Working with 2-dimensional arrays.
- **Day19** Brute force could not solve the second part. So we implemented an heuristics by keeping at each step only the permutations that were the closest to the expect result. Here, since we were looking to obtain the string "e", we just took the 20 shortest strings at each step.
- **Day20** Started by for each house, decompose in divisors: backwards chaining. But it did not scale, had to just forward-simulate elves running deliveries, then looking into houses.
- **Day21** Simple.
- **Day22** An interesting problem. I explored the branches of possibilities by recursing a function, separating data that were specific to the branch from the invariants in tow different structures. It made me understand well how to copy or pass by reference in Go, The code has a lot of tracing functions, as the debug was quite hairy.
- **Day23** A refreshingly simple problem. I tried to use types as much as possible
- **Day24** Quite simple, so I decided to use the math/big package to experiment. I also used a smart way to generate all the subsets of size k from a set of size n using bitsets.
- **Day25** Simple.

