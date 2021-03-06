#!/bin/bash
# https://adventofcode.com/2021/day/4 puzzle #1
# See README.md in the parent directory
in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1
err(){ echo "***ERROR: $*" >&2; exit 1;}

#TEST: example 4512
#TEST: input 8136

# This is a rewrite without using files with arraus

# Globals
size=5                          # size of square boards: 5x5
size2=$((size*2))
declare -a rows

# We store a board as size2 size strings (10 5-numbers strings), first the size rows
# then the size columns
# into a big array rows, so that board N occupies indexes N*size2 to (N+1)*size2-1
# A row is the string of numbers delimited by spaces and inside <RN: ... > where RN is
# the row number
readboard(){
    local row                   # array of values in the row
    local col                   # array of columns, space-separated values
    local i c                   # current row and col index
    # shellcheck disable=SC2034 # we do not use the empty var
    read -r empty || return 1   # error on EOF
    # first, read the size rows, pad with space, and store in rows[]
    for ((i=0; i<size; i++)); do
        read -r -a row
        rows+=("<${#rows[@]}: ${row[*]} >")
        # accumulate row per row on each of the columns
        for((c=0; c<size; c++)); do
            col[c]="${col[c]}${row[c]} "
        done
    done
    # then, append the completed columns as lines of space-separated values
    for ((c=0; c<size; c++)); do rows+=("<${#rows[@]}: ${col[c]}>"); done
}

# Parse the full input.
{
    # Numbers to draw are the first line, that we read in a space-separated list
    read -r drawline
    draws="${drawline//,/ }"
    # We parse all boards
    while readboard; do :; done
} <"$in"                        # we read the input sequentially for the parsing

# Now, draw all the numbers, remove them from the boards in place by string replaces
# Detecting a cleared row or column is then just finding an empty line in rows[]
# When a board wins, we stop
ndraw=0
for draw in $draws; do
    (( ++ndraw ))
    # we 's/ N / /g' in all rows in one command. This is the time saver!
    rows=("${rows[@]// $draw / }")
    # now check if some boards have won
    (( ndraw < size )) && continue # no board can win yet
    # shellcheck disable=SC2199 # yes, we want to match the whole array
    if [[ "${rows[@]}" =~ \<([[:digit:]]+):\ *\> ]]; then
        # we have a winner! and only one per the assignemnt
        rownum="${BASH_REMATCH[1]}"
        winner=$(( (rownum / size2) * size2)) # 1rst row of board
        sum=0
        # shellcheck disable=SC2013 # yes, we want to iterate on the words
        for((i=winner; i<(winner+size);i++)); do
            for number in ${rows[i]//[<>]/}; do
                [[ $number =~ : ]] || (( sum += number )) # sum all except RN:
            done
        done
        score=$(( sum * draw ))
        echo "Winner: Board #$((winner/size2)), sum=$sum, draw #$ndraw=$draw, score=$score"
        echo "$score"
        exit 0
    fi
done

err "No winning board found!"

