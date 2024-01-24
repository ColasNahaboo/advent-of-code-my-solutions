# Advent of code challenge 2017, in GO, day d03

By expanding the spiral:

```
   -5  -4  -3  -2  -1   0   1   2   3   4   5
-5 101                                     91 
-4     65  --  --  --  --  --  --  --  57 
-3      |  37  36  35  34  33  32  31   |
-2      |  38  17  16  15  14  13  30   |
-1      |  39  18   5   4   3  12  29   |
0       |  40  19   6   1   2  11  28   |
1       |  41  20   7   8   9  10  27   |
2       |  42  21  22  23  24  25  26   |
3       |  43  44  45  46  47  48  49  50 
4      73  --  --  --  --  --  --  --  --  82 
5  111                              ----> ...
```

We can see that for the diagonal axises, the next corner $n_{i+1}$ is $n_i+(n_i - n_{i-1}) + 8$. The difference with the previous corner is increased by 8.

Thus the SE start of a $i$-th circle is at $2n_{i+1}-n_i+8$, starting with 0 and 2, (for i=2, $n_i$=10 ) and, if the start position of 1 is at coord (0,0), andm $n_i$ to $n_i+2*i-1$ included (for i=2: from 10 to 13) coords x = i, y = 1-i to i (for i=2: 2,-1 to 2,2), etc...

For part1, we used this approach, as to only consider corners to speed things up.

For part2, we computed step by step with a "turtle" moving one step at a time and turning left once there was not an already walked trough position on its immediate left.


