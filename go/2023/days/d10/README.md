# Advent of code challenge 2023, in GO, day d10

If part1 was trivial, part2 was not overly complex but have some tricky points:
- Topological courses tells you that you find that a point is inside a line, if traceing a ray from this point cuts the line an odd number of times. Some people went for the general solution of the [Shoelace Formula](https://en.wikipedia.org/wiki/Shoelace_formula), but it was overkill in this case, especially since we were using discrete units.
- Before trying to cross the loop, the square `S` should be replaced by the actual pipe symbol that was hidden under it.
- What is "crossing" is hard to define when we trace a ray that can be parallel to sections of the pipe. I chose the trick of using a vertical ray and not counting the vertical pipe sections, `|`, and only counting half of the "elbow" ones, to avoid counting the line crossing twice. I chose `-`, `F` and `L`, the ones with an horzontal part to the right. The proper solution may have been to use a ray going at an oblique angle (e.g. 45 degrees), to avoid these border cases by never having the ray parralle to a pipe section.

