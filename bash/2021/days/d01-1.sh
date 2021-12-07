#!/bin/bash
# https://adventofcode.com/2021/day/1 puzzle #1
# See README.md in the parent directory
in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1

# brute force
previous=
count=0
while read -r number; do
    [[ -n $previous ]] && (( number > previous )) && (( ++count ))
    previous="$number"
done <"$in"
echo "$count"
