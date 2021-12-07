#!/bin/bash
# https://adventofcode.com/2021/day/2 puzzle #1
# See README.md in the parent directory
in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1

# Brute force. Ignore invalid lines.
horiz=0
depth=0
while read -r command value; do
    case "$command" in
        forward) (( horiz += value ));;
        down) (( depth += value ));;
        up) (( depth -= value ));;
    esac
done <"$in"
echo $(( horiz * depth ))
