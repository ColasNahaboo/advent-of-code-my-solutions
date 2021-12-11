#!/bin/bash
# https://adventofcode.com/days/day/11 puzzle #1
# See README.md in the parent directory
in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1
err(){ echo "***ERROR: $*" >&2; exit 1;}

#TEST: example 1656
#TEST: input 1617

steps="${2:-100}"               # number of steps in 2nd arg, default 100

# linearize the input into a single contiguous array, row by row
read -r -a levels < <(sed -r -e 's/./ \0/g' <"$in" | tr '\n' ' ')
# its dimensions
((cols = $(wc -c < <(head -1 "$in")) - 1))
((rows = $(wc -l <"$in") ))
levels_len=$((cols * rows))
[[ ${#levels[@]} != "$levels_len" ]] && err "input not a rectangle"

total=0                         # global: total number of flashes

# performs a step. Updates "total"
# If an octopus flashes, it splashes +1 on neighbors, that can overload
# and are then queued for next pass in "new"
step(){
    local i flashed=0 toflash new
    # [1] increase by 1 and mark the overloaded for [2]
    for ((i=0; i<levels_len; i++)); do
        (( (levels[i] += 1 ) == 10 )) && new="$new $i"
    done
    # [2] flash all marked, splash around them, mark overloaded ones
    # repeat pass [2] for the marked
    while [[ -n $new ]]; do
        toflash="$new"
        new=
        for i in $toflash; do
            ((levels[i])) || continue # already flashed ==> nothing
            flash "$i"
            (( i >= cols)) &&   # previous row, not for first one
                splash-row $((i - cols))
            splash-row "$i"                  # center row. i is marked so immune
            (( i < (levels_len - cols) )) && # next row, except for last
                splash-row $((i + cols))
        done
    done
    ((total += flashed))
}

# splash: add 1 and if over 9, mark for flash in later pass
splash(){
    local i="$1"
    ((levels[i])) || return   # already flashed, abort
    # only mark for next flash pass on the first time we go over 9
    (( (levels[i] += 1 ) == 10 )) && new="$new $i"
}

# flashes i and left and right neigbors if they exists (not on border)
# Updates globals flashed and new
splash-row(){
    local i="$1"
    (( (i % cols) > 0 )) && ((levels[i-1])) && splash $((i - 1))
    ((levels[i])) && splash "$i"
    (( (i % cols) < (cols - 1) )) && ((levels[i+1])) && splash $((i + 1))
}

# flashes i, and spashes the 8 adjacent position around it
flash(){
    local i="$1"
    levels[i]=0                           # flashed!
    (( ++flashed ))                       # increment flashed count
}

# functions to show a levels map, for tracing and debugging
display-levels(){
    local header="$1"
    local pause="$2"
    local zerochar="0"
    local i
    [[ -n $header ]] && echo "After step $step:"
    for ((i=0; i < levels_len; i++)); do
        if [[ ${levels[i]} == 0 ]]; then echo -n "$zerochar"
        else printf %x "${levels[i]}"
        fi
        (( (i % cols) == (cols - 1) )) && echo
    done
    if [[ -n "$pause" ]]; then
        read -r -p "[Return to continue]" rep # Pause
        [[ $rep =~ [qQnN] ]] && echo "$total" && exit 0
    fi
}    

for ((step=0; step<steps; step++)); do
    step
    #display-levels header pause   # DEBUG
done

display-levels header
echo "$total"

