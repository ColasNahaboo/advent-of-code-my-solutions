#!/bin/bash
# https://adventofcode.com/2021/day/3 puzzle #1
# See README.md in the parent directory
in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1

#TEST: example 198
#TEST: input 2954600

# Uses commands: rev

# The problem is ambiguous if there are as many set and unset bits in a
# column, so we raise an error in this case.
err(){ echo "***ERROR: $*" >&2; exit 1;}

# For convenience and efficiency,
# we split and reverse numbers into a shell array "digits" to read each bit:
# 001110100011 ==> 1 1 0 0 0 1 0 1 1 1 0 0
# The "set" array holds the per-column sums of digits

declare -a digits set
lines=0
cols=
while read -r -a digits; do
    if [[ -z $cols ]]; then     # inits
        cols=${#digits[@]}
        for col in "${!digits[@]}"; do
            (( set[col] = 0 ))
        done
    fi
    for col in "${!digits[@]}"; do
        (( set[col] += digits[col] ))
    done
    (( ++lines ))
done < <( rev < "$in" | sed -e 's/./\0 /g' )

(( half = lines / 2 ))

gamma=0
epsilon=0
col=0
colweight=1                     # 2^^col
while (( col < cols )); do
    if (( set[col] > half )); then
        (( gamma += colweight ))
    elif (( set[col] < half )); then
        (( epsilon += colweight ))
    elif [[ ${set[$col]} == half ]]; then
        # we cannot have the number of set bits equal to the unset ones
        err "half count set[$col] on col #$col"
    fi
    (( colweight *= 2 ))
    (( ++col ))
done

(( power = gamma * epsilon ))

echo "gamma = $gamma, epsilon = $epsilon, power = $power"
echo "$power"
