#!/bin/bash
# https://adventofcode.com/days/day/9 puzzle #2
# See README.md in the parent directory
in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1
nl=$'\n'
#TEST: example 1134
#TEST: input 1113424
#TEST: small 9

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

# compute the size of a basin from a point in it
# sets the size of the basin computed from the point into the global basin_size
# keep track of already mapped points in the global array basin
map_basin(){
    #global map basin basin_size
    local i="$1"
    local j="$2"
    local altitude="$3"
    basin["$i,$j"]=1              # seen, do not count twice
    (( basin_size++ ))
    local next
    if [[ -z "${basin[$((i-1)),$j]}" ]]; then
        next="${map[$((i-1)),$j]}"
        ((next < 9)) && ((next > altitude)) &&
            map_basin $((i-1)) "$j" "$next"
    fi
    if [[ -z "${basin[$((i+1)),$j]}" ]]; then
        next="${map[$((i+1)),$j]}"
        ((next < 9)) && ((next > altitude)) &&
            map_basin $((i+1)) "$j" "$next"
    fi
    if [[ -z "${basin[$i,$((j-1))]}" ]]; then
        next="${map[$i,$((j-1))]}"
        ((next < 9)) && ((next > altitude)) &&
            map_basin "$i" $((j-1)) "$next"
    fi
    if [[ -z "${basin[$i,$((j+1))]}" ]]; then
        next="${map[$i,$((j+1))]}"
        ((next < 9)) && ((next > altitude)) &&
            map_basin "$i" $((j+1)) "$next"
    fi
}

# brute force
# Iterate on all the map, find the minimums, and recursively "flow up"
# the slopes of the basin
sizes=                          # the basins sizes found, one per line
for ((i=0; i<cols; i++)); do
    for ((j=0; j<rows; j++)); do
        altitude="${map[$i,$j]}"
        (( altitude >= 9 )) && continue # 9 cannot be a low point
        if (( map[$((i-1)),$j] > altitude)) &&
               (( map[$((i+1)),$j] > altitude)) &&
               (( map[$i,$((j-1))] > altitude)) &&
               (( map[$i,$((j+1))] > altitude)); then
            declare -A basin    # keep track of already mapped places
            basin_size=0
            map_basin "$i" "$j" "$altitude"
            unset basin
            sizes="${sizes}${sizes:+$nl}$basin_size"
        fi
    done
done

# keep the 3 biggest sizes
sizes3=$(echo "$sizes" |sort -nr | head -3)
# multiply them
total=1
for size in $sizes3; do (( total *= size )); done
echo "$total"
