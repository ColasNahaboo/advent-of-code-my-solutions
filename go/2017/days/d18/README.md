# Advent of code challenge 2017, in GO, day d18

The part 1 is simple.

For the part2, I modelized the tablets as objects, and run them "by hand" one a t a time, the tablets comunicating by FIFO queues. But since I wanted to practice the use of Go goroutines and channels, I also made an implementation via a goroutine for each tablet, communicating by channels, that is called via the command line flag `-3`.

We use 3 goroutines: the two tablet `T0` and `T1`, and an `observer` of their status whose role is to detect when both tablets are waiting for input on a `rcv` operation.

The observer maintains an `alive status` (`1` if alive, `0` if waiting) for each tablet. On executing a `rcv` a tablet sends a clear status message to the observer channel, and once a value is read, it sends a set status.

Thus the observer can thus tell that the two tablets are in deadlock when both status are 0, and can terminate the whole program.

```
 +----+                      +----+
 | T0 |-----channel-0-1----->| T1 |
 |    |<----channel-1-0----- |    |
 +----+                      +----+
    |                          |     +----------+
    +-------observer-channel---+---> | observer |
                                     +----------+
```
