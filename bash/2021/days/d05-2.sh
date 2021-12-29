#!/bin/bash
# https://adventofcode.com/2021/day/5 puzzle #2
# See README.md in the parent directory
in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1
err(){ echo "***ERROR: $*" >&2; exit 1;}
#export tmp=tmp.$$; clean(){ rm -f "$tmp" "$tmp".*;}; trap clean 0

#TEST: example 12
#TEST: input 17741

# As the naive method of creating a 2-dimensional aray and plotting the
# lines could be too slow in bash, we try the following hack:
# We will just store lines as the set of their points in X,Y form
# from A to B one per line, and we will just have to look for duplicate points
# It is thus very fast for huge sparse data

# as a variant from d05-1.sh, we use an array instead of a file
declare -A lines

while read -r ax ay arrow bx by; do
    : "$arrow"                  # keep shellcheck happy
    if (( ax == bx )); then   # vertical line
        if (( ay < by )); then min="$ay"; max="$by"
        else min="$by"; max="$ay"
        fi
        for ((y=min; y <= max; y++)); do lines["$ax,$y"]+=1; done
    elif (( ay == by )); then   # horizontal line
        if (( ax < bx )); then min="$ax"; max="$bx"
        else min="$bx"; max="$ax"
        fi
        for ((x=min; x <= max; x++)); do lines["$x,$ay"]+=1; done
    else
        # We square as a cheap way to get the absolute value of the differences
        dx=$(( (ax - bx) ** 2 ))
        dy=$(( (ay - by) ** 2 ))
        if (( dx == dy )); then # same ==> 45 degrees line
            ((ax < bx)) && dx=1 || dx=-1
            ((ay < by)) && dy=1 || dy=-1
            x="$ax"; y="$ay"
            while true; do
                lines["$x,$y"]+=1
                (( x == bx )) && break
                (( x += dx ))
                (( y += dy ))
            done
        fi                      #  ignore oblique lines
    fi
done < <(tr ',' ' ' <"$in")     # X,Y becomes X Y in input

# the result is just the number of duplicated values
duplicates=0
for line in "${lines[@]}"; do
    ((line > 1)) && ((duplicates++))
done
echo "$duplicates"

