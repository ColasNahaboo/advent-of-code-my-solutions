# Advent of code challenge 2021, in bash
Here are my solutions to the "Advent of code" challenge of 2021 implemented in bash.
See https://adventofcode.com/2021

I tried to implement "smart" bash solutions (see at the end, "Algorithmic tricks"), that rely on fast GNU/Linux utilities like grep, sed, sort... (the way I use bash in real life), but most of the time my goal was readability and compliance with shellcheck, not terseness or efficiency. I am an old (retired) proficient bash programmer, so my goal here is not to learn bash, but rather un-learn some bad dirty habits accumulated over the years and force myself to code in a "shellcheck-friendly" modern way.

Although I invented, designed and implemented a commercial programming language (The SML - System Management Language - in the Bull ISM Network Management platform), and worked professionnally with "real" languages I grew in love with bash because the intellectual challenges it poses to write efficient code, making mundane tasks exciting, and that it (nearly) never breaks backwards compatibility, a code that runs now wioll run in 20 years and more... And since I retired in 2021 I now have time to play with things, and discovered the "Advent of code" challenge.

**Usage:**
- Everything is in the `days/` subdirectory, with no subdirectories. All files are in this directory, no subdirs. I hate navigating a maze of small directories, all alike.
- Day `NN` solutions are scripts `dNN-1.sh` and `dNN-2.sh` for problem 1 and 2. E.g. `d05-2.sh`
- Run them without input data file as argument (defaults to `dNN.input`). E.g. `d14.input`
- Input data is in `dNN.input`. Note that the data may be different for different accounts on the web site.
- The small sample input in the problem text is in `dNN.example`
- So running all the codes can be done by: `for i in d*.sh; do ./$i; done` (or via `EXECALL`)
- There is nearly no error checking, to keep the code small and readable.
  We only check for non-obvious errors.
- Alternate solutions are sometime left for reference in files `dNN-otherM.sh`, e.g: `d06-other2.sh`

**Notes:**
- These scripts need bash version 4.4+, and GNU linux utilities, so they will not run on non-GNU unix systems such as vanilla MacOS or FreeBSD.
- Scripts should pass shellcheck: run `CHECKALL`, or:
  `for i in d*.sh; do shellcheck -f gcc $i; done`
- When temporary files are used, they are cleaned with a clean function called on any exit with trap 0
- If we detect errors, they are listed on stderr, with the function `err`
- To generate files for test number `N`: `MAKEDAY N`

**Style:**
- We use `(( ... ))` and `[[ ... ]]` rather than let and `[ ... ]`
- We put all variable evaluations in double quotes: `"$foo"`
- We use `$( ... )` and never backquotes
- We avoid forking subprocess, by using `bar < <(foo)` instead of `foo | bar`. This way `bar` is not run in a separate variable scope from `foo`.

## Some utility code commonly used
At the start of the scripts, some convenient code can be found if needed.

Set variable `in` to the input file argument, with `dNN.input` as default for script `dNN.sh` and exit in error if the file does not exists.
`in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1`

Raise a fatal error:
`err(){ echo "***ERROR: $*" >&2; exit 1;}`

Use temporary files `$tmp` or any prefixed by `$tmp.` and delete them on any exit
`tmp=tmp.$$; clean(){ rm -f "$tmp" "$tmp".*;}; trap clean 0`

## Commands used
- The classic ones: `grep`, `sed`, `sort`, `uniq`
  Note that the perl-compatible `grep -oP` is only in GNU grep, and very useful in bash scripts.
- `rev` reverses the parameters: `echo a1 b2 c3 | rev` prints `3b 2c 1a`

## More info
- https://adventofcode.com/ The "advent of code" (aka AOC) website
- https://github.com/Bogdanp/awesome-advent-of-code Bogdan Popa list of AOC-related resources and solutions
- https://www.reddit.com/r/adventofcode/ The reddit to discuss AOC 
  - [Big inputs for AOC 2021](https://www.reddit.com/r/adventofcode/comments/r9s5pz/2021_big_inputs_for_advent_of_code_2021_puzzles/) A collection of huge inputs to stress your solutions

## License
Author: (c)2021 Colas Nahaboo, https://colas.nahaboo.net
License: free of use via the [MIT License](https://en.wikipedia.org/wiki/MIT_License).

## Algorithmic tricks
For a lot of problems, solving the problem the straightforward way is too slow in Bash. So I have used algorithmic tricks in some solutions (or seen some on the web labelled "Also:"). I have commented them in the scripts, but here they are, collected for reference:

### Day 3
To split a binary number into an array `digits` of its bits, in reverse (little endian) order, I use: `rev < "$in" | sed -e 's/./\0 /g | read -a digits'`
- `rev` to reverse the characters of the input string which is the binary number
- `sed` to split the number into a space-separated list of bits
- `read -a` to create an array of these bits

### Day 4
In a textural representation of a board as lines of space-separated numbers, detecting an empty row is just grep-ing for empty lines. I thus used files containing the board followed by its "inverted" version, each line being a column instead of a row, so that looking for empty rows or columns is done by a grep of empty lines.

I also pad the lines of numbers by spaces before and after, so I can perform a sed to remove the drawn number `N` in all files by replacing `{space}N{space}` by `{space}` without taking into account the special case of the first or last number in the line.

The code is generic and has variables to adjust the size of boards.

*Also:* I have seen on the Reddit post [Big inputs for AOC 2021](https://www.reddit.com/r/adventofcode/comments/r9s5pz/2021_big_inputs_for_advent_of_code_2021_puzzles/) a smart hack by "p_tseng":

> For each board, figure out what time it wins at. By building a map of called number -> time it's called, you can determine this very quickly. For each board, take maximum across each row/column, take minimum of those maxima.

## Day 5
Creating a (huge) 2-dimensional aray and plotting the  lines in it would be too slow in bash, so we just list the lines points one by one, each on a line in a sequential file. The points that are thus parts of more than one line are simply the coordinates that appears multiple times in the file! So a simple classic `sort|uniq -d` can find them quickly. This is a huge gain in efficiency, both in time and spacem, for data with few lines sparse in a big space as the input files are.

## Day 6
There is no way we can solve this problem with an exponential algorithm, especially with the 256 number of steps of the second exercise. So we use a linear approach where we do not actually generate a representation of individual fishes, but just keep track of the number of newborns that will be born each day in the future. When we add a fish, either at the start or because it is born on the day, we just increment the count of future newborns at the days its timer will reach zero.

We can then just iterate on the days and add the its newborn fishes this way. details in the comments of `d06-1.sh`.

I have kept some attempts that were correct but too slow for reference:
- `d06-other1.sh` the straightforward version implementin naively the explanations.
- `d06-other2.sh` a version trying to capitalize on the speed of grep with a full representation of the fishes in a file, but inverted.
These attempts were useful as it allowed me an agild coding. I could have quickly some prototypes, unusable for the "production" use of tackling the 256 steps in reasonable time, but simple enough to give correct answers to test the validity of the final solution.