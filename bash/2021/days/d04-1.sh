#!/bin/bash
# https://adventofcode.com/2021/day/4 puzzle #1
# See README.md
in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1
err(){ echo "***ERROR: $*" >&2; exit 1;}
tmp=tmp.$$; clean(){ rm -f "$tmp" "$tmp".*;}; trap clean 0

# This can be slow, around 30s.

# Globals
rows=5; cols=5                  # size of board rows & columns
(( cols < rows )) && min="$cols" || min="$rows"

# Parses the board number N, store it in a $tmp.board.N file,
# containing the rows and then the columns (inverted matrix)
# We pad the values with spaces so sed and grep do not need special cases for
# the first and last values
readboard(){
    local boardnum="$1"
    local row                   # array of values in the row
    local col                   # array of columns, space-separated values
    local r c value             # current row and col index, and the value
    # shellcheck disable=SC2034 # we do not use the empty var
    read -r empty || return 1   # error on EOF
    for ((c=0; c<cols; c++)); do col[c]=' '; done
    # first, copy the rows into the file
    for ((r=0; r<rows; r++)); do
        # short read? abort
        read -r -a row || { rm -f "$tmp.board.$boardnum"; return 1;}
        echo " ${row[*]} " >>"$tmp.board.$boardnum"
        c=0
        for value in "${row[@]}"; do
            col[c]="${col[c]}$value "
            (( ++c ))
        done
    done
    # then, append the columns as lines of space-separated values
    for ((c=0; c<cols; c++)); do
        echo "${col[c]}" >>"$tmp.board.$boardnum"
    done
}

# Parse the full input.
{
    # Numbers to draw are the first line, that we read in a space-separated list
    read -r drawline
    draws="${drawline//,/ }"
    # We parse all boards
    boardnum=0
    while readboard "$boardnum"; do (( ++boardnum )); done
} <"$in"                        # we read the input sequentially for the parsing

# Now, draw all the numbers, remove them from the boards in place by sed
# Detecting a cleared row or column is then just grep-ing for an empty line
# in the board files
ndraw=0
for draw in $draws; do
    (( ++ndraw ))
    for file in "$tmp".board.*; do # remove drawn number from all boards
        grep -q " $draw " "$file" && sed -r -i -e "s/ $draw / /g" "$file"
    done
    # now check if some boards have won
    (( draw < min )) && continue # no board can win yet
    # grep for cleared row or col: empty lines
    winners=$(grep -l '^[[:space:]]*$' "$tmp".board.* |
                  grep -oP '[.]board[.]\K[[:digit:]]+')
    [[ -z $winners ]] && continue # non winners yet
    # we have winner(s)! Normally only one, but just in case, display all
    # we sum the numbers only in the first half of the file, the rows
    for board in $winners; do
        sum=0
        # shellcheck disable=SC2013 # yes, we want to iterate on the words
        for unmarked in $(head -"$rows" "$tmp.board.$board"); do
            (( sum += unmarked ))
        done
        score=$(( sum * draw ))
        echo "Winner: Board #$board, sum=$sum, draw #$ndraw=$draw, score=$score"
    done
    exit 0
done

err "No winning board found!"

