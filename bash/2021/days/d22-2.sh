#!/bin/bash
# https://adventofcode.com/days/day/22 puzzle #2
# See README.md in the parent directory
in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1
err(){ echo "***ERROR: $*" >&2; exit 1;}
#export tmp=tmp.$$; clean(){ rm -f "$tmp" "$tmp".*;}; trap clean 0

#TEST: small 39
#TEST: example 39769202357779
#TEST: input 

# This implements the solution of the reddit megathread, one described by aexl:

# Algorithm: Parse the input cuboids. Then start with an empty list (clist) of
# cuboids. For each cuboid from the input, calculate the intersections with the
# cuboids in the clist. If the intersection is non-empty add it to the clist
# with inverted sign w.r.t to the current cuboid in clist. If that cuboid from
# the input is set on, add it to the clist as well. Then add together all the
# volumes (cuboids with a negative sign will have a negative volume).

# A cuboid is described by "type x1 x2 y1 y2 z1 z2", type: 1=on, -1=off

declare -a clist                # list of processed cuboids
declare -i clistvol             # total volume of clist

# merges an input cuboid to the clist
merge-cuboid(){
    local -i i clistLen=${#clist[@]}
    for((i=0; i<clistLen; i++)); do
        # shellcheck disable=SC2086 # yes we split the cuboids by param passing
        add-cuboid-intersection ${clist[i]} "$@"
    done
    (($1 > 0)) && append-cuboid "$@"
}

# compute cuboids intersection, and append with inverted sign wrt first
add-cuboid-intersection(){
    local -i sign1="$1" x11="$2" x12="$3" y11="$4" y12="$5" z11="$6" z12="$7"
    local -i x21="$9" x22="${10}" y21="${11}" y22="${12}" z21="${13}" z22="${14}"
    local -i x1 x2 y1 y2 z1 z2
    if ((x21 <= x12)); then
        ((x22 < x11)) && return
        if ((x21 > x11)); then ((x1=x21)); else ((x1=x11)); fi
        if ((x22 < x12)); then ((x2=x22)); else ((x2=x12)); fi
    else return
    fi
    if ((y21 <= y12)); then
        ((y22 < y11)) && return
        if ((y21 > y11)); then ((y1=y21)); else ((y1=y11)); fi
        if ((y22 < y12)); then ((y2=y22)); else ((y2=y12)); fi
    else return
    fi
    if ((z21 <= z12)); then
        ((z22 < z11)) && return
        if ((z21 > z11)); then ((z1=z21)); else ((z1=z11)); fi
        if ((z22 < z12)); then ((z2=z22)); else ((z2=z12)); fi
    else return
    fi
    append-cuboid "$((-1 * sign1))" "$x1" "$x2" "$y1" "$y2" "$z1" "$z2"
}

# just append to clist
append-cuboid(){
    local -i sign="$1" x1="$2" x2="$3" y1="$4" y2="$5" z1="$6" z2="$7"
    local -i vol
    ((vol = sign * (x2 - x1 + 1) * (y2 - y1 + 1) * (z2 -z1 +1)))
    clist+=("$*")
    ((clistvol += vol))
}
    
############ Main
# we read the cuboids, and apply them to the clist
n=0
tot=$(wc -l <"$in")
# shellcheck disable=SC2020 # tr use is legit!
while read -r line; do          # sign x1 x2 y1 y2 z1 z2
    # shellcheck disable=SC2086 # yes we split the cuboids by param passing
    merge-cuboid $line
    echo "==[$((++n))/$tot] \"$line\", clist: ${#clist[@]}"
done < <(tr -cs '\n[onf0-9]-' ' ' <"$in" |sed -e 's/on/1/' -e 's/off/-1/')

echo "== clist has ${#clist[@]} cuboids, for a total volume of:"
echo "$clistvol"
