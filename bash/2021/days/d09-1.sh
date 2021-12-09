#!/bin/bash
# https://adventofcode.com/days/day/9 puzzle #1
# See README.md in the parent directory
in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1
err(){ echo "***ERROR: $*" >&2; exit 1;}
#TEST: example 15
#TEST: input 591
#TEST: small 1

# Bash does not have multi-dimensional arrays. We simulate them with an
# associative array "map": map[$i,$j] queries the map hashtable for the key
# string "$i,$j"
# Also, in order to avoid having special cases for borders, we add a border
# around the table with altitudes higher than any in the table: 10
# This is fast enough: 100ms for the full input

# Another solution could have been to have just an array of rows,
# each row being a string, and getting the element i,j would be:
# ${${map[$j]}:$i:1}
# But note that then we could not add a border at index -1

# get map dimensions. We will use indexes i for cols, j for rows
((cols = $(wc -c < <(head -1 "$in")) - 1))
((rows = $(wc -l <"$in") ))

# fill the map, and put a high border (10) around it
declare -A map
for j in -1 $rows; do           # top and bottom borders
    for ((i=-1; i<=cols; i++)); do map[$i,$j]=10; done
done
for i in -1 $cols; do    # left and right borders
    for ((j=0; j<rows; j++)); do map[$i,$j]=10; done
done
# fill the inside map
j=0
while read -r line; do
    for ((i=0; i<cols; i++)); do map[$i,$j]="${line:$i:1}"; done
    (( j++ ))
done <"$in"

# brute force
# Iterate on all the map, find the minimums, and compute the cumulated risk
lows=
numof_low=0
risk=0
for ((i=0; i<cols; i++)); do
    for ((j=0; j<rows; j++)); do
        altitude="${map[$i,$j]}"
        (( altitude >= 9 )) && continue # 9 cannot be a low point
        if (( map[$((i-1)),$j] > altitude)) &&
               (( map[$((i+1)),$j] > altitude)) &&
               (( map[$i,$((j-1))] > altitude)) &&
               (( map[$i,$((j+1))] > altitude)); then
            (( risk += altitude + 1 ))
            lows="$lows $altitude"
            (( numof_low ++ ))
            # the following code is not p;art of the solution per se:
            # the problem text is ambiguous, as it does not specify what to do
            # when we are equal to a neighbour. We thus raise an error, so that
            # we are sure we are not in this case and we do not have to worry
        elif (( map[$((i-1)),$j] >= altitude)) &&
                 (( map[$((i+1)),$j] >= altitude)) &&
                 (( map[$i,$((j-1))] >= altitude)) &&
                 (( map[$i,$((j+1))] >= altitude)); then
            if (( map[$((i-1)),$j] == altitude)) ||
               (( map[$((i+1)),$j] == altitude)) ||
               (( map[$i,$((j-1))] == altitude)) ||
               (( map[$i,$((j+1))] == altitude)); then
                echo "At $i,$j, altitude $altitude is the same than a neighbor
 ${map[$i,$((j-1))]}
${map[$((i-1)),$j]}$altitude${map[$((i+1)),$j]}
 ${map[$i,$((j+1))]}"
            fi
        fi
    done
done

echo "$numof_low low points:$lows"
echo "$risk"
