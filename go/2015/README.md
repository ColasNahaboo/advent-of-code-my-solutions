# Advent of code challenge 2015, in GO

Here are my solutions to the "Advent of code" challenge of 2015 implemented in bash.
See https://adventofcode.com/2015

I am doing this to learn GO, so this must be considered as "student code". I am coding it to try my hand at various GO idioms, not seeking efficiency, scalability nor optimality. But feedback is very welcome.

The code is in standard GO, with some housekeeping scripts in bash.

**Usage:**

- Everything below is in the `days/` sub-directory, with a sub-directory `dayNN/` per day.
- Day `NN` solutions are in program of source code `days/dayNN/dayNN.go`, and executable `days/dayNN/dayNN`
- Run them with input data file (with suffix `.txt`) as argument (defaults to `input.txt`).
- The result will always be a number, alone on the last line of the output.
- They run the algorithm for Part 2 of the daily problems, unless you give the option `-1` where they will run the Part 1.
- All source code is standalone, I will try to use only standard GO and standard library functions. I will copy common code from templates in the `TEMPLATES/` directory rather than making proper packages, so you can just download a day directory and run it, it is self-contained.
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
- **Day01** Starting simple, deciding on the directopry and file structure, experimenting with testing and debugging via dlv
- **Day02** Working with regexps.
- **Day03** Working with maps.
- **Day04** Using package md5.
- **Day05** Compensing the limits of GO RE2 regexp standard.
- **Day06** Working with 2-dimensional arrays.
- **Day07** Working with function pointers and kind of closures.
- **Day08** Simple.
- **Day09** A brute-force Traveling Salesman Problem. I used a nice hack to generate all the routes (permutations), found at https://golangbyexample.com/all-permutations-string-golang/
- **Day10** Interesting discovery of the "Look and say" analysis by John Conway, with a [video](https://www.youtube.com/watch?v=ea7lJkEhytA). A naive strings-based implementation was much too slow, but using arrays of integers proved very fast.
- **Day11** Converted strings to work on arrays of integers.
- **Day12** Working with JSON, interfaces and type switches.
- **Day13** 
- **Day14** 
- **Day15** 
- **Day16** 
- **Day17** 
- **Day18** 
- **Day19** 
- **Day20** 
- **Day21** 
- **Day22** 
- **Day23** 
- **Day24** 
- **Day25** 

