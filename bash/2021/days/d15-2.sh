#!/bin/bash
# https://adventofcode.com/days/day/15 puzzle #2
# See README.md in the parent directory
in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1
err(){ echo "***ERROR: $*" >&2; exit 1;}
#export tmp=tmp.$$; clean(){ rm -f "$tmp" "$tmp".*;}; trap clean 0

#TEST: example 315
#TEST: input 2817

# The only diff with d15-1.sh is that we duplicate the input map

# we will call level a value of a point on the map,
# and risk the sum of risk levels on a path

############ Parse map
declare -a map                  # rectangle X x Y, in a single array
                                # level(x,y) => map[x + y*X]
declare -a map1                 # The parrsed map that we will enlarge 5x
Y1=0
while read -r line; do
    # shellcheck disable=SC2206
    map1+=($line)
    ((++Y1))
done < <(sed -r -e 's/./ \0/g' <"$in")
((area1 = ${#map1[@]}))
((X1 = area1 / Y1))
((MAXINT=8888888888888888888))  # easier on the eyes for debug than 2**63-1
((area = area1 * 25))
((X = X1 * 5))

# generate the duplicates
repeatmap(){
    local -i i ox="$1" oy="$2"
    local -i inc=$((ox + oy)) x1 y1 x y
    for((i=0; i<area1; i++)); do
        ((newval=map1[i]+inc))
        ((newval > 9)) && ((newval -= 9))
        ((x1 = i%X1))
        ((y1 = i/X1))
        ((x = ox*X1 + x1))
        ((y = oy*Y1 + y1))
        ((map[x +y*X] = newval))
    done
}
for y in 0 1 2 3 4; do
    for x in 0 1 2 3 4; do
        repeatmap "$x" "$y"
    done
done

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

echo "${risks[area-1]}"
