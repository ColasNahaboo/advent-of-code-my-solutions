#!/bin/bash
# https://adventofcode.com/days/day/7 puzzle #1
# See README.md in the parent directory
in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1
err(){ echo "***ERROR: $*" >&2; exit 1;}
tmp=tmp.$$; clean(){ rm -f "$tmp" "$tmp".*;}; trap clean 0
#TEST: example 37
#TEST: input 336131

# brute force: we compute the rule for each position inside the crabs

# all positions, one per line, in increasing order
crabs="$(tr ',' '\n' <"$in")"
numof_crabs=$(tr ',' '\n' <"$in" |wc -l)
positions="$(tr ',' '\n' <"$in" |sort -n |uniq)"
minpos=$(echo "$positions" | head -1)
maxpos=$(echo "$positions" | tail -1)

# start with valid, but maybe non-optimal values
optimalpos="$minpos"
optimalfuel=$(( (maxpos - minpos) * numof_crabs ))

for ((pos=minpos; pos <= maxpos; pos++)); do
    fuel=0
    for crab in $crabs; do
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
