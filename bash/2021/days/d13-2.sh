#!/bin/bash
# https://adventofcode.com/days/day/13 puzzle #2
# See README.md in the parent directory
in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1
err(){ echo "***ERROR: $*" >&2; exit 1;}
#export tmp=tmp.$$; clean(){ rm -f "$tmp" "$tmp".*;}; trap clean 0

#TEST: example 16
#TEST: input 97
# Note: the actual result of the run is the text code in ascii art on the paper

# We manage only the of dots in dots[] rather than the paper grid itself
# Each dot is its input coordinates (string) x,y
# The folding instructions are in the array folds[] as x=N or y=M

declare -A dots
declare -a folds
# first read the dots
{
    while read -r line; do
        if [[ $line =~ ^[[:digit:]]+,[[:digit:]]+ ]]; then
            dots["$line"]=1
        else
            break
        fi
    done
    while read -r line; do
        [[ $line =~ 'fold along'\ *([xy]=[[:digit:]]+) ]] &&
            folds+=("${BASH_REMATCH[1]}")
    done
} <"$in"

# fold: for each dot, if after fold copy them over the fold, and remove old 
fold(){
    local coord="$1"            # x or y
    local fold="$2"             # the fold value
    local x y nx ny
    for dot in "${!dots[@]}"; do
        x="${dot%,*}"
        y="${dot#*,}"
        if [[ $coord == x ]] && ((x > fold)); then
            ((nx = 2*fold - x))
            dots["$nx,$y"]=1
            unset dots["$dot"]
        elif [[ $coord == y ]] && ((y > fold)); then
            ((ny = 2*fold - y))
            dots["$x,$ny"]=1
            unset dots["$dot"]
        fi
    done
}

# apply all folds
for fold in "${folds[@]}"; do fold "${fold%=*}" "${fold#*=}"; done

# display folded paper. The letters shown are the actual result, the code
minx=1000; miny=1000; maxx=-1000; maxy=-1000
for dot in "${!dots[@]}"; do
    x="${dot%,*}"
    y="${dot#*,}"
    ((x < minx)) && minx="$x"
    ((y < miny)) && miny="$y"
    ((x > maxx)) && maxx="$x"
    ((y > maxy)) && maxy="$y"
done
for((y=miny; y<=maxy; y++)); do
    for((x=minx; x<=maxx; x++)); do
        [[ -n ${dots["$x,$y"]} ]] && echo -n '#' || echo -n '.'
    done
    echo
done

echo ${#dots[@]}
