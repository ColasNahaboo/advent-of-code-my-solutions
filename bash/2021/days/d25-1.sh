#!/bin/bash
# https://adventofcode.com/days/day/25 puzzle #1
# See README.md in the parent directory
in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1
err(){ echo "***ERROR: $*" >&2; exit 1;}
#export tmp=tmp.$$; clean(){ rm -f "$tmp" "$tmp".*;}; trap clean 0

#TEST: example 58
#TEST: input 534

verbose="${2}"                  # 2nd argument prints maps

# we use as map values .=0, >=1, v=cols so that to move, we just add its value
# Note: the idea was to advance by just adding the value
# But I ended up not using it as it was not handling properluy the wraparound
# so the code below just considers the values 0, 1, and anything else.
declare -i cols rows size
declare -ai map                 # [col, row] == map[row*cols+col]

declare -i i j r c                # temp vars

# read initial state
r=0
while read -r line; do
    cols=${#line}
    for((c=0;c<cols;c++)); do
        if [[ ${line:c:1} == '>' ]]; then map+=(1)
        elif [[ ${line:c:1} == 'v' ]]; then map+=("$cols")
        else map+=(0)
        fi
    done
    ((++r))
done <"$in"
rows="$r"
((size=rows*cols))

# just used in verbose mode for debugging
print-map(){
    local -i r c i
    for((r=0; r<rows; r++)); do
        for((c=0; c<cols; c++)); do
            ((i=r*cols+c))
            if ((map[i]==0)); then echo -n '.'
            elif ((map[i]==1)); then echo -n '>'
            else echo -n 'v'
            fi
        done
        echo
    done
}

[[ -n $verbose ]] && { echo "Initial:"; print-map;}
step=0
moved=true
# each step consists of 2 passes, one for > the other for v
# on each padd we register the moves to do in new[], but only apply them
# at the end of each pass. new[] is then a sparse array.
while "$moved"; do
    moved=false
    # first, all the >
    new=()
    for((i=0; i<size; i++)); do
        if ((map[i]==1)); then
            # j is the position to move to, with wraparound
            ((j=(i/cols)*cols+((i%cols)+1)%cols))
            if ((map[j]==0)); then
                ((new[i]=0))    # register move for later
                ((new[j]=1))
                moved=true
            fi
        fi
    done
    # shellcheck disable=SC2068 # yes, elements are space-splitted
    for i in ${!new[@]}; do ((map[i]=new[i])); done # apply moves
    # then, all the v
    new=()
    for((i=0; i<size; i++)); do
        if ((map[i]==cols)); then
            # j is the position to move to, with wraparound
            ((j=((i/cols)+1)%rows*cols+i%cols))
            if ((map[j]==0)); then
                ((new[i]=0))
                ((new[j]=cols))
                moved=true
            fi
        fi
    done
    # shellcheck disable=SC2068 # yes, elements are space-splitted
    for i in ${!new[@]}; do ((map[i]=new[i])); done # apply moves

    ((++step))
    [[ -n $verbose ]] && { echo "==== Step $step:" print-map;}
done

echo "$step"
