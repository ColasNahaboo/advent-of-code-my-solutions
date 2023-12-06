# Advent of code challenge 2023, in GO, day d05

The part1 was quite simple, however the naive approch to examine each seed does not scale for the part2.

Since for each mapping, the resulting values are always bigger than the one for the start of the range interval, we will thus examine only the start of the range, and discard the rest as it cannot have a lower location. However, as we go through the maps, we must split the range if we fall in a smaller range for the n+1 map.
