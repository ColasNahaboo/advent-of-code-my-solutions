#!/usr/bin/python3
# https://github.com/AshGriffiths/advent_of_code/blob/main/2023/day_nine/day9.py

# Comparison with my code:
# n=$(wc -l input,,.test|cut -f 1 -d ' '); i=0; while ((i<n)); do head -$i input,,.test |tail -1 >INPUT.txt; s=$(./sol1.py|tail -1); m=$(./d?? -1 INPUT.txt|tail -1); echo "[$i]: s=$s, m=$m"; ((m == s)) || { echo "**DIFFER!"; cat INPUT.txt; break;}; ((i++)); done

from itertools import pairwise


def next_reading(reading_set):
    if len(set(reading_set)) == 1:
        return reading_set[0]
    return reading_set[-1] + next_reading([x[1] - x[0] for x in pairwise(reading_set)])


with open("INPUT.txt", "r") as input:
    readings = [[int(val) for val in line.split(" ")] for line in input.readlines()]
    p1_results = []
    p2_results = []
    for reading_set in readings:
        p1_results.append(next_reading(reading_set))
        p2_results.append(next_reading(reading_set[::-1]))
    #print(f"Part One : {sum(p1_results)}")
    #print(f"Part Two : {sum(p2_results)}")
    print(sum(p1_results))
