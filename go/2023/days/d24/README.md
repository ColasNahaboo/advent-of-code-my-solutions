# Advent of code challenge 2023, in GO, day d24

Part1 was trivial, but part2 was quite challenging.

A general solution seems to see the problem as an optimization one, obtain a set of linear equations and reduce them.

Another can be to place ourselves in the referential of a hailstone, which is thus at position 0,0,0 and velocity 0,0,0 and try to work from there.

However, it seems that on every input there are two hailstones that have the same position and velocity on one of the XYZ axis. So, if we call these two $ha$ and $hb$ having the same position and velocity (aka $px$ and $vx$) on the X axis, and $r$ the rock, we can say:

- The rock must also have the same psx and vx, otherwise it could never meet ha and hb, so we then know their values: $r.px = ha.px$ and $r.vx = ha.vx$
- Let's take two other hailstones at random: $h0$ and $h1$
- At time $t$, the X coord of any hailstone $h$ (or the rock) being $h.px + t * h.vx$, we have at the time $t0$ where the rock hits $h0$ the equation: $r.px + t0 * r.vx = h0.px + t0 * h0.vx$
- And the same for h1: $r.px + t1 * r.vx = h1.px + t1 * h1.vx$
- We can thus deduce $t0$ and $t1$
- Now we can compute the positions of $r$ at $t0$ and $t1$ by factorizing its velocities in the two equations above. We get: $r.px = (t0* t1 * h1.vx - t0* t1 * h0.vx + t0 * h1.px - t1 * h0.px) / (t0 - t1)$
- And we repeat the computations on the other 2 axis to get the rock position!

Note that most of the input coordinates being 50-bit numbers, calculations must be done with big integers as multiplaying them will overflow the 64 bit integers of Go.

