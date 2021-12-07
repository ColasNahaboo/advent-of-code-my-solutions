#!/bin/bash
# https://adventofcode.com/days/day/7 puzzle #1
# See README.md in the parent directory
in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1

#TEST: example 37
#TEST: input 336131

# brute force: we compute the rule for each position inside the crabs

# array of positions of each crab
read -r -a crabs < <(tr ',' ' ' <"$in")
# compute the range of possible positions, min&max of the occupied ones
positions="$(tr ',' '\n' <"$in" |sort -n |uniq)"
minpos=$(echo "$positions" | head -1)
maxpos=$(echo "$positions" | tail -1)

# start with valid, but maybe non-optimal values
optimalpos="$minpos"
optimalfuel=$(( (maxpos - minpos) * ${#crabs[@]} ))

# iterate by computing fuel costs for all possible positions
for ((pos=minpos; pos <= maxpos; pos++)); do
    fuel=0
    for crab in "${crabs[@]}"; do
        ((crab >= pos)) && ((move = crab - pos)) || ((move = pos - crab))
        ((fuel += move))
    done
    if ((fuel < optimalfuel)); then
        optimalpos="$pos"
        optimalfuel="$fuel"
    fi
done

echo "Position: $optimalpos, fuel:
$optimalfuel"
