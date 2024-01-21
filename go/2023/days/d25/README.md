# Advent of code challenge 2023, in GO, day d25

The problem is cutting a graph into two sub-graphs.

A a brute-force approch (in the function `BruteForceCut` works for the example since all combinations of 3 links (or graph edges) out of its 33 is only 32736.
But it cannot work for the 3450 edges of my input that would require looking at more than 40 billions of combinations.

We then must find an algorithm do do this, and I used the Kargen-Stein one, detailed in 
https://stanford-cs161.github.io/winter2021/assets/files/lecture16-slides.pdf 

The slide 95 giving a pseudo-code summary of the algorithm itself:

KargerStein of a graph with initial number of nodes $n$:
- if $n$ < 4:
  - find a min-cut by brute force
- Run Karger’s algorithm on G with independent repetitions until $n/\sqrt(2)$ nodes remain
- Make 2 copies copies of what’s left of G: G1 and G2
- Recurse on them, cutting a different edge at random on each one
  - S1 = KargerStein(G1)
  - S2 = KargerStein(G2)
- return whichever of S1 or S2 is the smaller cut.

Karger is the simpler version:
- if $n$ < 4:
  - find a min-cut by brute force
- Run Karger’s algorithm on G with independent repetitions
I coded it in the KargerMinCut function, the input is not big enough to warrant a Karger-Stein algorithm
