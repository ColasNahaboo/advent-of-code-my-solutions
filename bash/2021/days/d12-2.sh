#!/bin/bash
# https://adventofcode.com/days/day/12 puzzle #2
# See README.md in the parent directory
in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1
err(){ echo "***ERROR: $*" >&2; exit 1;}

# with a second arg, display all found paths
display_paths=false; (( $# == 2 )) && display_paths=true

#TEST: example1 36
#TEST: example2 103
#TEST: example3 3509
#TEST: input 133360

# we store all the caves in a "caves" array, and then we only use their index
declare -a caves                # their "bigness": 0 or 1
declare -A names                # for each cave name, its index in caves
names['end']=0                  # we force "end" to be at index 0
caves[0]=0                      # and it is a small one
names['start']=1                # and we predeclare start at 1
caves[1]=0
# for each cave, "tunnels" is the space-delimited list of connected caves
declare -a tunnels

############ Parsing

# from a name, sets its cave index in global "cave" and auto-register it
# we use globals for efficiency, to avoid having to copy arrays
declare-cave(){
    local name="$1"
    cave="${names[$name]}"
    if [[ -z $cave ]]; then     # append cave to the end of caves[]
        cave=${#caves[@]}
        names[$name]="$cave"
        if [[ $name =~ ^[[:lower:]]+$ ]]; then caves[cave]=0
        else caves[cave]=1
        fi
    fi
}

# we register the connections from the declared cave1 to cave2
declare-tunnel(){
    local cave1="$1" cave2="$2"
    # shellcheck disable=SC2076 # yes, we want to quote $cave2
    if ! [[ ${tunnels[cave1]} =~ " $cave2 " ]]; then
        [[ -z ${tunnels[cave1]} ]] && tunnels[cave1]=" "
        tunnels[cave1]+="$cave2 "
    fi
}

# parse the input file into our caves model
while read -r line; do
    [[ $line =~ ^[[:space:]]*$ ]] && continue # ignore empty lines
    [[ $line =~ ^([[:alnum:]]+)-([[:alnum:]]+)$ ]] ||
        err "Syntax error at line: \"$line\""
    name1="${BASH_REMATCH[1]}"
    name2="${BASH_REMATCH[2]}"
    declare-cave "$name1"
    cave1="$cave"
    declare-cave "$name2"
    cave2="$cave"
    [[ "$cave1" == "$cave2" ]] && err "loop detected: $line"
    declare-tunnel "$cave1" "$cave2"
    declare-tunnel "$cave2" "$cave1"
done <"$in"

############ Debug utils, unused normally

# display a path for trace/debug
display-path(){
    local path="$1" cave sep="  " name
    for cave in $path; do
        name-of-cave "$cave"
        echo -n "$sep$name"
        sep=','
    done
    echo
}
# return the cave name in global name
name-of-cave(){
    local cave="$1"
    for name in "${!names[@]}"; do
        (( ${names[$name]} == cave )) && return
    done
    name='?'
}
# dumps the names and indexes of caves, in order
display-all-caves(){
    local name cave
    echo -n "Caves:"
    for ((cave=0; cave < ${#caves[@]}; cave++)); do
        name-of-cave "$cave"
        echo -n " $cave=$name"
    done
    echo
}

############ Pathfinding

# now we bruteforce our way to find all paths

# $1: the cave to explore one step from, and recurse
# $2: the already travelled path, space-delimited
# $3: does this path already contains a small cave twice? (true/false)
explore(){
    local cave="$1" path="$2" twice="$3"
    local next
    # echo "Explore: $1 $2 $3"    # DEBUG
    for next in ${tunnels[cave]}; do
        if ! ((next)); then     # end: found a path
            "$display_paths" && display-path "$path$next "
            ((++paths))
        elif ((caves[next])); then # big cave, OK, go through it
            explore "$next" "$path$next " "$twice"
        elif [[ $next == 1 ]]; then : # dont go twice through 'start'
        elif "$twice"; then
            # small cave, but already went twice in a small, or "start"
            # shellcheck disable=SC2076 # yes, we want to quote $next
            [[ $path =~ " $next " ]] || explore "$next" "$path$next " "$twice"
        else                    # small cave, but we still can visit twice
            # shellcheck disable=SC2076 # yes, we want to quote $next
            if [[ $path =~ " $next " ]]; then
                explore "$next" "$path$next " true # we went twice in "next"
            else
                explore "$next" "$path$next " false
            fi
        fi
    done
}

paths=0

#display-all-caves # DEBUG
explore 1 " 1 " false

echo "$paths"



