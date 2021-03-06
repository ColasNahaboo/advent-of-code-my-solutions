#!/bin/bash
# https://adventofcode.com/2021/day/2 puzzle #2
# See README.md in the parent directory
in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1

#TEST: example 900
#TEST: input 1880593125

# Brute force. Ignore invalid lines.
horiz=0
depth=0
aim=0
while read -r command value; do
    case "$command" in
        forward)
            (( horiz += value ))
            (( depth += ( aim * value) ))
            ;;
        down) (( aim += value ));;
        up) (( aim -= value ));;
    esac
done <"$in"
echo $(( horiz * depth ))
