#!/bin/bash
# https://adventofcode.com/days/day/10 puzzle #1
# See README.md in the parent directory
in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1

#TEST: example 26397
#TEST: input 436497

# We implement a tree parser, pushing the current chunk type on a stack
# We cannot use associative arrays because we cannot use ([{< chars as indexes
# Instead we use the string operators # and % to fetch our data in a string
# E.g, instead of using map[key], we define map="key1value1 key2value2..."
# and find value by cuttting left and right:
# value=${map#*key}; value=${value%% *}

# global variables (exported to subprocess calls)
export openchars='([{<'         # chunks opening chars
export closechars=')]}>'        # chunks closing chars
export openclosechars='()[]{}<>'         # the pairings between them
export chunkscores=')3 ]57 }1197 >25137' # list of illegal chars and scores

# prints the syntax error score of the line
# can be 0,3 57 1197 25137
parse-line(){
    local line="$1"
    local c i openchunk closechunk stack score=0
    declare -a stack               # the stack of chars
    local top=0                    # index of next available stack position

    for ((i=0; i<${#line}; i++)); do
        c="${line:i:1}"
        # shellcheck disable=SC2076 # yes, we want to quote $c
        if [[ $openchars =~ "$c" ]]; then # push on stack
            stack[top]="$c"
            (( ++top ))
        elif [[ $closechars =~ "$c" ]]; then # pull from stack and check
            if (( top <= 0 )); then          # empty stack, should not happen
                score=$(get-score "$c")
                break
            else
                (( --top ))
                openchunk="${stack[top]}" # char opening the curent chunk
                closechunk="${openclosechars#*[$openchunk]}" # expected close
                if [[ "$c" != "${closechunk:0:1}" ]]; then
                    score=$(get-score "$c")
                    break
                fi
            fi
            # other chars are ignored
        fi
    done
    echo "$score"
}

get-score(){
    local c="$1" score
    score="${chunkscores#*[$c]}"
    score="${score%% *}"
    echo "$score"
}

total=0
while read -r line; do
    score=$(parse-line "$line")
    ((total += score))
done <"$in"

echo "$total"
