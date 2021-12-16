#!/bin/bash
# https://adventofcode.com/days/day/15 puzzle #1
# See README.md in the parent directory
in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1
err(){ echo "***ERROR: $*" >&2; exit 1;}
#export tmp=tmp.$$; clean(){ rm -f "$tmp" "$tmp".*;}; trap clean 0

#TEST: example 40
#TEST: input 398

# we will call level a value of a point on the map,
# and risk the sum of risk levels on a path

############ Parse map
declare -a map                  # rectangle X x Y, in a single array
                                # level(x,y) => map[x + y*X]
Y=0
while read -r line; do
    # shellcheck disable=SC2206
    map+=($line)
    ((++Y))
done < <(sed -r -e 's/./ \0/g' <"$in")
((area = ${#map[@]}))
((X = area / Y))
((end = area - 1))
((MAXINT=8888888888888888888))  # easier on the eyes for debug than 2**63-1

############ DEBUG utils
D(){
    local i r
    for((i=0; i<area; i++)); do
        ((i % X == 0)) && echo
        ((risks[i] == MAXINT)) && r='*' || r="${risks[i]}"
        printf "%10s" "$i(${map[i]})[$r]"
    done
    echo
}
DF(){
    local i
    echo -n "Frontier:"
    for i in "${!frontier[@]}"; do
        echo -n " $i(${frontier[i]})"
    done
    echo
}
############ END DEBUG

# We perform a Dijsktra algorithm, with risk levels being distances
# https://www.geeksforgeeks.org/dijkstras-shortest-path-algorithm-greedy-algo-7/
declare -a risks                # cache map of minimal path risks from 0
((risks[0]=0))
declare -a visited              # is the node visited? blank or visited index
declare -a frontier             # unvisited nodes connected to a visited
                                # value: minimum distance to a visited
frontier[0]=0                   # initial node

# update risks. Uses globals risk, risks[], map[]. Add to frontier
update(){
    local i="$1" r
    ((visited[i])) && return
    local nextrisk
    ((nextrisk = risk + map[i]))
    r="${frontier[i]:-$MAXINT}"
    ((nextrisk < r)) && {
        ((frontier[i] = nextrisk))
    }
}  

# minimum of the frontier
minimum(){
    local i min="$MAXINT"
    for i in "${frontier[@]}"; do ((i<min)) && ((min=i)); done
    for i in "${!frontier[@]}"; do
        ((frontier[i] == min)) && echo "$i" && return
    done
    echo 0                      # nothing to do, end
}

while ((${#frontier[@]})); do
    u=$(minimum)
    ((visited[u]=1))
    ((risk = frontier[u]))
    unset "frontier[u]"
    ((risks[u] = risk))
    # Update the 4 neigbours if they are not visited
    ((u&X > 0)) && update $((u-1))
    ((u > X)) && update $((u-X))
    ((u%X < (X-1))) && update $((u+1))
    ((u < (area -X))) && update $((u+X))
done

echo "${risks[end]}"
