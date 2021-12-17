#!/bin/bash
# https://adventofcode.com/days/day/17 puzzle #1
# See README.md in the parent directory
in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1
err(){ echo "***ERROR: $*" >&2; exit 1;}
#export tmp=tmp.$$; clean(){ rm -f "$tmp" "$tmp".*;}; trap clean 0

#TEST: example 45
#TEST: input 4278

############ Parse input
# shellcheck disable=SC2034 # unused values for reads
{
    read -r input <"$in"
    [[ $input =~ 'target area: 'x=([-[:digit:]]+)[.]+([-[:digit:]]+),[[:space:]]*y=([-[:digit:]]+)[.]+([-[:digit:]]+) ]] || err "Input syntax error: $input"
    x1="${BASH_REMATCH[1]}"
    x2="${BASH_REMATCH[2]}"
    y1="${BASH_REMATCH[3]}"
    y2="${BASH_REMATCH[4]}"
}
    
############ Main
# We know the initial velocity must be the lower Y of the target area
# absolute value -1, as we cross back altitude 0 at the same vertical velocity
# and the next step will be +1

((vy0 = -y1 -1))

echo "vy0 = $vy0"

# then, just compute the max height reached, step by step.
y=0
for((vy=vy0; vy>0; vy--)); do
    ((y += vy))
done

echo "$y"
