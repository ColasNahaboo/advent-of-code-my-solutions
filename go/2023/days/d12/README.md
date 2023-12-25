# Advent of code challenge 2023, in GO, day d12

Part1 was quite simple. I solved it by just generating all the possible arragements with all the combinations of value s'.' or '#' for the '?', and then checking all of them to see which ones were valid. But this would not scale for part2.

I thus tried an iterative approach: moving up the string of the springs records, each time I encountered a '?', explored on the 2 possible branches. Nice, but not enough. I kept the code (functions named `explore`) for this approach as an alternate to part2, callable with the `-a` command line option

I then took the same approach, but skipping the '.' of the record string, and trying to fit the whole first span of contiguous '#' at the first '#' position. And then recursing not by character, but by spans.

And finally, I added a cache of calls to the recurive function, to gain a bit more speed.... but it was too slow!

I ended up copying the logic of a [python solution](https://gist.github.com/sanyi/96ccaf6d3c0a67536b4fe3e99bc53bb3) by "sanyi", that was simple and fast. 
