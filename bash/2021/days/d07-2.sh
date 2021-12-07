#!/bin/bash
# https://adventofcode.com/days/day/7 puzzle #2
# See README.md in the parent directory
in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1
err(){ echo "***ERROR: $*" >&2; exit 1;}
tmp=tmp.$$; clean(){ rm -f "$tmp" "$tmp".*;}; trap clean 0
#TEST: example 168
#TEST: input 92676646

# brute force: we compute the rule for each position inside the crabs
# the only optimization is to precompute the fuel costs for each possible move
# Runs fast enough in 7s

# all positions, one per line, in increasing order
crabs="$(tr ',' '\n' <"$in")"
numof_crabs=$(tr ',' '\n' <"$in" |wc -l)
positions="$(tr ',' '\n' <"$in" |sort -n |uniq)"
minpos=$(echo "$positions" | head -1)
maxpos=$(echo "$positions" | tail -1)
# precompute the fuel costs for a distance
cost[0]=0
for ((steps=1; steps <= (maxpos - minpos); steps++)); do
    ((cost[steps] = cost[steps-1] + steps))
done

# start with valid, but non-optimal values
optimalpos="$minpos"
optimalfuel=$((2 ** 62))        # MAXINT

for ((pos=minpos; pos <= maxpos; pos++)); do
    fuel=0
    for crab in $crabs; do
        ((crab >= pos)) && ((move = crab - pos)) || ((move = pos - crab))
        ((fuel += cost[move]))
    done
    if ((fuel < optimalfuel)); then
        optimalpos="$pos"
        optimalfuel="$fuel"
    fi
done

echo "Position: $optimalpos, fuel:
$optimalfuel"
