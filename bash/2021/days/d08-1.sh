#!/bin/bash
# https://adventofcode.com/days/day/8 puzzle #1
# See README.md in the parent directory
in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1

#TEST: example 26
#TEST: input 504

# simple brute force: first emit the output values one per line,
# and count the ones with the proper length via a regexp.
# Possible lengths: 2 4 3 7 for the digits 1 4 7 8

grep -oP '\| +\K.*' <"$in" |    # keep only the values after the "|"
    tr ' ' '\n' |               # emit one value per line
    grep -Ec '^(..|....|...|.......)$' # count the ones of the proper lengths
