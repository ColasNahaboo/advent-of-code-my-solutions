#!/bin/bash
# https://adventofcode.com/days/day/14 puzzle #1
# See README.md in the parent directory
in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1
err(){ echo "***ERROR: $*" >&2; exit 1;}
#export tmp=tmp.$$; clean(){ rm -f "$tmp" "$tmp".*;}; trap clean 0

#TEST: example 1588
#TEST: input 3259

# naive approach, wont scale, as each step takes twice as much time

steps="${2:-10}"                # optional 2nd arg is the number of steps

declare -a poly                 # the polymer as array of pairs: NNCB = NN NC CB
declare -a ifs                  # the left part (conditions) of insertion rules
declare -a thens                # the right part as space separated two pairs
                                # HN -> C: if=HN, then="HC CN"
# shellcheck disable=SC2034     # yes, we do not use some read variables
{                               # parsing
    read -r template
    for((i=0; i < ${#template}-1; i++)); do
        poly+=("${template:$i:2}")
    done
    read -r empty
    while read -r ruleif arrow rulethen; do
        ifs+=("$ruleif")
        thens+=("${ruleif:0:1}$rulethen $rulethen${ruleif:1:1}")
    done
} <"$in"

# now, in each step, expanse all poly pairs
for((step=0; step<steps; step++)); do
    # apply all rules, and store in poly2 the letters to insert into poly
    unset poly2; declare -a poly2
    for((p=0; p < ${#poly[@]}; p++)); do
        for((r=0; r < ${#ifs[@]}; r++)); do
            [[ "${poly[p]}" == "${ifs[r]}" ]] && poly2[p]="${thens[r]}" && break
        done
    done
    # build the new poly into poly3, either a copy or an insertion
    unset poly3; declare -a poly3
    for((i=0; i < ${#poly[@]}; i++)); do
        if [[ -n "${poly2[i]}" ]]; then
            # shellcheck disable=SC2206 # yes, we split the space-separated pair
            poly3+=(${poly2[i]}) # adds 2 elements
        else
            poly3+=("${poly[i]}")
        fi
    done
    # poly3 is then the new poly, ready for next step
    poly=("${poly3[@]}")
done

# compute the quantity of elements
nl=$'\n'
elements="${poly[0]:0:1}"       # poly as string of one element per line
for((i=0; i < ${#poly[@]}; i++)); do
    elements="$elements${nl}${poly[i]:1:1}"
done
sorted=$(echo "$elements" | sort | uniq -c | sort -n)
# shellcheck disable=SC2034 # e is unused
mostcommon=$(echo "$sorted" | tail -1 | { read -r n e; echo "$n";})
# shellcheck disable=SC2034 # e is unused
leastcommon=$(echo "$sorted" | head -1 | { read -r n e; echo "$n";})
echo $((mostcommon - leastcommon))
