#!/bin/bash
# https://adventofcode.com/days/day/12 puzzle #1
# See README.md in the parent directory
in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1
err(){ echo "***ERROR: $*" >&2; exit 1;}

#TEST: example1 10
#TEST: example2 19
#TEST: example3 226
#TEST: input 4792

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
        [[ ${names[$name]} == "$cave" ]] && return
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

# explore one step from the cave
# 2nd arg: the already travelled path, space-delimited
explore(){
    local cave="$1" path="$2"
    local next
    for next in ${tunnels[cave]}; do
        if ! ((next)); then     # end: found a path
            # display-path "$path$next " # DEBUG
            ((++paths))
        elif ((caves[next])); then # big cave, OK, go through it
            explore "$next" "$path$next "
        else                    # small cave, do go there twice
            # shellcheck disable=SC2076 # yes, we want to quote $next
            [[ $path =~ " $next " ]] || explore "$next" "$path$next "
        fi
    done
}

paths=0

#display-all-caves # DEBUG
explore 1 " 1 "

echo "$paths"



