# Advent of code challenge 2017, in GO, day d17

Part2 can be coded quite efficiently if you note that the first position has always value 0 as it is never inserted, since we always insert after. So, unlike part1, we do not need toa ctually build the circular buffer! We just have to register the value when the position wraps up at 0 (the first position).

If we call:
- `s` the step, our program input
- `v` the iteration value, which is also the length of the circular buffer as we always insert one new element at each iteration
- `p` the start position

We get:
```
   p-s->q
1234567890123
<-----v----->

   p-s->q
123456789v0123
<-----v+1---->
```
So we only have to compute the successives values of `p`, registering the value `v` when `q` wraps to 0, never actually building the buffer.
