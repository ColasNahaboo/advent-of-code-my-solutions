#!/bin/bash
# https://adventofcode.com/days/day/18 puzzle #2
# See README.md in the parent directory
in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1
err(){ echo "***ERROR: $*" >&2; exit 1;}
#export tmp=tmp.$$; clean(){ rm -f "$tmp" "$tmp".*;}; trap clean 0

#TEST: example3 3993
#TEST: input 4731


# We use only# [2] SFS (SnailFish Strings) that are the textual representations
# of SFNs: "[[1,2],3]"

############ Magnitude

declare -i magnitude
sfs=''

# Compute magnitude in the first "number" (pair or natural number) found in sfs
magnitude(){
    # global sfs magnitude
    local -i lm
    magnitude=0
    if [[ ${sfs:0:1} == '[' ]]; then # pair
        sfs="${sfs:1}"
        magnitude; ((lm=magnitude))
        sfs="${sfs:1}"
        magnitude
        sfs="${sfs:1}"
        ((magnitude = 3*lm +2*magnitude))
    elif [[ $sfs =~ ^([[:digit:]]+)(.*) ]]; then # natural number
        magnitude="${BASH_REMATCH[1]}"
        sfs="${BASH_REMATCH[2]}"
    else
        err "Syntax error: not a SFN: \"$sfs\""
    fi
}

############ reduce

# reduce global "sfs" in place: apply the 2 rules till they dont fire
reduce(){
    while try-explode || try-split; do :; done
}

# If any pair is nested inside four pairs, the leftmost such pair explodes.
# To explode a pair, the pair's left value is added to the first regular number
# to the left of the exploding pair (if any), and the pair's right value is
# added to the first regular number to the right of the exploding pair (if
# any). Exploding pairs will always consist of two regular numbers. Then, the
# entire exploding pair is replaced with the regular number 0.
try-explode(){
    local s="$sfs" c left right
    local -i lv rv lrn rrn i open=0
    for((i=0; i<${#sfs}-8; i++)); do
        c="${s:i:1}"
        if [[ $c == '[' ]]; then
            if ((++open > 4)); then # explode
                [[ ${s:i} =~ ^\[([[:digit:]]+),([[:digit:]]+)\](.*)$ ]] ||
                    err "Explodee not a pair of regnums"
                lv="${BASH_REMATCH[1]}"
                rv="${BASH_REMATCH[2]}"
                right="${BASH_REMATCH[3]}"
                left="${s:0:i}"
                if [[ $left =~ ^(.*[\[,])([[:digit:]]+)([^[:digit:]]*)$ ]]; then
                    lrn=$((BASH_REMATCH[2] + lv))
                    left="${BASH_REMATCH[1]}$lrn${BASH_REMATCH[3]}"
                fi
                if [[ $right =~ ^([^[:digit:]]*)([[:digit:]]+)(.*)$ ]]; then
                    rrn=$((BASH_REMATCH[2] + rv))
                    right="${BASH_REMATCH[1]}$rrn${BASH_REMATCH[3]}"
                fi
                sfs="${left}0${right}"
                return 0        # yes, exploded!
            fi
        elif [[ $c == ']' ]]; then
            ((--open))
        fi
    done
    return 1                    # no explosives found.
}

# If any regular number is 10 or greater, the leftmost such regular number
# splits.
# To split a regular number, replace it with a pair; the left element of the
# pair should be the regular number divided by two and rounded down, while the
# right element of the pair should be the regular number divided by two and
# rounded up. For example, 10 becomes [5,5], 11 becomes [5,6], 12 becomes
# [6,6], and so on.

try-split(){
    local s="$sfs" c left right
    local -i n
    if [[ $s =~ ([[:digit:]][[:digit:]]+)(.*)$ ]]; then
        n="${BASH_REMATCH[1]}"
        right="${BASH_REMATCH[2]}"
        left="${s:0:$((${#s} - ${#n} - ${#right}))}"
        lv=$((n / 2))
        rv=$(((n+1)/2))
        sfs="${left}[$lv,$rv]${right}"
        return 0
    fi
    return 1
}

############ Main

# read all numbers and store the sfs into a nums[] array
declare -a nums
while read -r sfs; do
    [[ -z $sfs ]] && continue
    reduce
    nums+=("$sfs")
done <"$in"

size=${#nums[@]}
max=0                           # the largest magnitude

# now, compute the magnitude of all the possible sums of 2 of these numbers
for((i=0; i<size; i++)); do
    for((j=0; j<size; j++)); do
        [[ $i == "$j" ]] && continue #  numbers must be different
        sfs="[${nums[i]},${nums[j]}]"
        reduce
        magnitude
        ((magnitude > max)) && max="$magnitude"
    done
done
echo "Max of Magnitudes:"
echo "$max"

