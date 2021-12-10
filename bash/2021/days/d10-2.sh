#!/bin/bash
# https://adventofcode.com/days/day/10 puzzle #2
# See README.md in the parent directory
in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1

#TEST: example 288957
#TEST: input 2377613374

# We implement a tree parser, pushing the current chunk type on a stack
# We cannot use associative arrays because we cannot use ([{< chars as indexes
# Instead we use the string operators # and % to fetch our data in a string
# E.g, instead of using map[key], we define map="key1value1 key2value2..."
# and find value by cuttting left and right:
# value=${map#*key}; value=${value%% *}

# global variables (exported to subprocess calls)
export openchars='([{<'         # chunks opening chars
export closechars=')]}>'        # chunks closing chars
export openclosechars='()[]{}<>' # the pairings between them
export chunkscores='(1 [2 {3 <4' # list of completion (open) chars and scores

# prints the completion score of the line, ignore invalid ones
parse-line(){
    local line="$1"
    local c i openchunk closechunk stack score=0 pointvalue invalid=false
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
                invalid=true
                break
            else
                (( --top ))
                openchunk="${stack[top]}" # char opening the curent chunk
                stack[top]=               # uneeded, but eases debugging
                closechunk="${openclosechars#*[$openchunk]}" # expected close
                if [[ "$c" != "${closechunk:0:1}" ]]; then
                    invalid=true
                    break
                fi
            fi
            # other chars are ignored
        fi
    done
    "$invalid" && return

    # now we find the completion characters on the stack
    while ((top > 0)); do
        (( top-- ))
        c="${stack[top]}"
        pointvalue="${chunkscores#*[$c]}"
        pointvalue="${pointvalue%% *}"
        (( score = score * 5 + pointvalue ))
    done
    echo "$score"
}

# get the sorted list of scores
scores=$(while read -r line; do
             parse-line "$line"
         done <"$in" | sort -n)
# outputs the middle one
numof_scores=$(wc -l <<<"$scores")
middle_pos=$((numof_scores / 2 + 1))
head -"$middle_pos" <<<"$scores" | tail -1
