# Advent of code challenge 2023, in GO, day d06

Instead of trying all the values, we just solve the problem mathematically.

For a race of duration $d$, keeping the button pressed for a time $t$ makes the boat goes to length $l$ by the formula:

$l=speed *(d-t)$

That is, since $speed$ is given by the time $t$ itself:

$l=t(d-t) \implies -t^2+dt$

So, the winning times $w$ are the ones that can generate a distance $r$, equal or greater to the current record + 1.

$-t^2+dt=r \implies -t^2+dt-r=0$

This means that the winning times are on a parabola, in the "bump" between the two solutions of this second degree equation, which are:

$(-d ± \sqrt({d^2-4r})/-2 \implies (d ± \sqrt({d^2-4r})/2$

For instance, for the first race of the example, with $d=7$ and $r=10$ (9+1), the solutions are $(7 +- \sqrt{49 - 40}) / 2$, that is $(7-3)/2=2$ and $(7+3)/2=5$. So the solutions are the integer between 2 and 5, inclusive.

For the second race, we get the solutions 3.5948... and 11.4051... we round these values up and down, to 4 and 11, to stay inside the acceptable values range.
