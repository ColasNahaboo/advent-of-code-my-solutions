# Advent of code challenge 2017, in GO, day d11

An hexa grid can be simulated on a sparse normal grid, by using one column out of 4, 
```
 ....1.......2
 3.......4....
 ....5.......6
 7.......8....
 ....9.......a
 b.......c....
 ```
 With adjacent directions being:
 - n  = [0, -2]
 - ne = [4, -1]
 - se = [4, 1]
 - s  = [0, 2]
 - sw = [-4, 1]
 - nw = [-4, -1]

Or we can use axial coordinates, see https://www.redblobgames.com/grids/hexagons/
```
     1       2
     | \     |
 3   |   4   |           +
 | \ |   | \ |           |\
 |   5   |   6           | \
 |   | \ |   |           |  \
 7   |   8   |           |   \
 | \ |   | \ |           v    J
 |   9   |   a           Y     X
 |     \ |
 b       c
```
```
   2
 146
358a
79c
b
 ```
 With adjacent directions being:
 - n  = [0, -1]
 - ne = [1, -1]
 - se = [1, 0]
 - s  = [0, 1]
 - sw = [-1, 1]
 - nw = [-1, 0]

Or, as in a 3D cube

```
        +
       /|\
      / | \
     /  |  \
    /   |   \
   L    v    J
  X     Z     Y
```

With adjacent directions being:
 - n  = [0, -1, 1]
 - ne = [1, -1, 0]
 - se = [1, 0, -1]
 - s  = [0, 1, -1]
 - sw = [-1, 1, 0]
 - nw = [-1, 0, 1]
