# Advent of code challenge 2023, in GO, day d18

In part 1, we just dig the trenches, and record each tiome we turn at the end of a segment. By adding 1 for a turn to the right and -1 for a turn to the left, the total sum at the end of the loop tells us if the loop buckled onto the right or left. This gives us where to find a point in its inside, if we loop right, the point one step along the first segment and one step to its right. We then thus fill the adjacent cells to find all the cells inside the trenches.

In part 2, the map is too big to use the naive approach of managing the map as cells of 1 square meter. We thus look at all the possible values that x and y coordinates take at the start (and thus ends, since the plan of trenches is a closed loop), and divide the plane into an grid of tiles of various sizes. All the cells in the trenches are 1x1 tiles, and tiles in between can stretch for long distances between "interesting" coordinates.

If we take a simple exemple:

```
     original map 23x8            19 tiles on a map 5x5:
   xlist=0,7,21   ylist=0,3,7     all 16 '#' 1x1 tiles, 
                                  plus these 3 "big" tiles:
   0111111233333333333334         
   |      |             |         
   ########              --0      ### 
   #      #                1      #1#       1 = 6x2 
   #      #                1      #-### 
   #      ###############--2      #2|3#     2 = 6x3 
   #                    #  3      ##### 
   #                    #  3                3 = 13x3 
   #                    #  3 
   ######################--4 
```

We thus apply the same method as in part 1, but for this vitual map, knowing that the capacity of each tile is not always 1.
