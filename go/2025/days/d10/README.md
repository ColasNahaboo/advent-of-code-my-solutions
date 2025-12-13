# Advent of code challenge 2025, in GO, day d10

[Problem statement](https://adventofcode.com/2025/day/10)

For part2, we tried the naive approach to consider the possible button pushes as a graph of, and finding the shortest path to the state with no leds lit to the state where all the leds have the desired pattern, the steps between states being a button push.

It worked on the example, but it was too slow to work on the actual input with bigger expected joltages.

I took thus the opportunity to try using a solver, by expressing the problem in term of simple liear equations, with variables being non-negative. It is coded in the `part2withZ3Solver` function, made `-3` option, and the default.

For the first machine of the example:
```
[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}
```
This gives the equation to solve, with `a` the number of presses on first button, `b` on the second, etc:
```
(3)a + (1,3)b + (2)c + (2,3)d + (0,2)e + (0,1)f =  {3,5,4,7}
```
Written in matricial form:
```
0     0     0     0     1     1     3 
0 a + 1 b + 0 c + 0 d + 0 e + 2 f = 5 
0     0     1     1     1     0     4 
1     1     0     1     0     0     7 
```
This gives the equations to solve:
```
e + f = 3
b + f = 5
c + d + e = 4
a + b + d = 7
```

We then use the Z3 solver with its [Go-Z3](https://pkg.go.dev/github.com/aclements/go-z3) bindings. Alas, the distribution does not include a [Go wrapper for the Optimizer](https://github.com/aclements/go-z3/pull/2/commits/fc200cae08e68459443ef54e6cf93f2914c2fe24), (and has really minimal docs with no examples) so we optimize by hand by adding iteratively a constraint, an inequality to assert that the sum should be smaller than the one we just found.

The iteration is thus:
- First find a solution, get the sum `s` of all variables (here `s = a+b+c+d+e+f`)
- then add the constraint `a+b+c+d+e < s`
- re-try as long as we find solutions
- Our result is the last valid `s found`

