#!/bin/bash
# https://adventofcode.com/days/day/14 puzzle #2
# See README.md in the parent directory
in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1
err(){ echo "***ERROR: $*" >&2; exit 1;}
#export tmp=tmp.$$; clean(){ rm -f "$tmp" "$tmp".*;}; trap clean 0

#TEST: example 2188189693529
#TEST: input 3459174981021

steps="${2:-40}"

# we linearize the problem by considering only the possible pairs,
# Simple straightforward implementation, no other optimizations
# The poly is thus just an array of the numbers of each pair ID in it
declare -A poly                 # the polymer as pairs: {pair, count-in-poly}
declare -A letters              # the count of letters in poly {letter, count}
declare -A rules                # {pair, inserted-letter}

# shellcheck disable=SC2034     # yes, we do not use some read variables
{                               # parsing
    read -r template
    # count the pairs & letters in the initial poly description
    for((i=0; i < ${#template}-1; i++)); do
        ((poly[${template:i:2}]+=1))
        ((letters[${template:i:1}]+=1))
    done
    ((letters[${template:i:1}]+=1)) # we stopped before the last letter
    read -r empty
    while read -r ruleif arrow letter; do
        rules["$ruleif"]="$letter"
    done
} <"$in"

# then, in each step, "expanse" the rules, but do not store the polymer letters
# just increment the counts of created/deleted pairs in a "todo" list
# The letters, we can increase immediately, as it has no impact on computations
for((step=0; step < steps; step++)); do
    unset todo; declare -A todo
    for pair in "${!poly[@]}"; do
        ((number=poly[$pair]))      # number of this pair in poly[]
        ((todo["$pair"]-=number))   # remove the pairs as matches rule
        newletter="${rules[$pair]}"         # insert the 2 new ones
        ((todo[${pair:0:1}$newletter]+=number))
        ((todo[$newletter${pair:1:1}]+=number))
        ((letters[$newletter]+=number)) # add the new one to the counts
    done
    # apply the todo list
    for pair in "${!todo[@]}"; do
        ((poly[$pair]+=todo[$pair]))
        ((poly[$pair])) || unset poly["$pair"] # remove entries at count 0
    done
done

# compute score
((min=2**62))
max=0
for n in "${letters[@]}"; do
    ((n < min)) && ((min=n))
    ((n > max)) && ((max=n))
done
echo $((max - min))
