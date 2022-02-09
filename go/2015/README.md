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

**Testing:**

- Unit tests are performed via the standard GO testing system, in the source file `days/dayNN/dayNN_test.go`
- Integration tests are done by looking at the comments `// TEST: [option] input-file result` in source files and running the code with the option and input and checking the last printed line is the result. The `days/TESTALL` bash script runs all the unit and integration tests, see it for technical details.
- The examples given in the problem descriptions are used in GO unit tests, whereas the input file is used for the integration tests.
