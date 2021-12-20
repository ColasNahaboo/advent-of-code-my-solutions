#!/bin/bash
# https://adventofcode.com/days/day/18 puzzle #2
# See README.md in the parent directory
in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1
err(){ echo "***ERROR: $*" >&2; exit 1;}
#export tmp=tmp.$$; clean(){ rm -f "$tmp" "$tmp".*;}; trap clean 0

#TEST: example3 4140
#TEST: input 4731


# We will work with a dual representation of SFNs (SnailFish Numbers)
# [1] as a hierarchical tree, with data in a heap:
# SFNs are pairs X+Y, where X and Y can be:
# - a Natural Number: 435, 645, 87...
# -a pointer to another pair: @345 where 345 is the index of the SFN in the heap
# [2] SFS (SnailFish Strings) that are the textual representations of SFNs:
# "[[1,2],3]"
#
# E.g if sfs is "[[1,2],3]", sfn is: @12, with heap[12]=@27+3 heap[27]=1+2

declare -a heap                 # the heap of SFNs
declare -A heapidx              # reverse index: string ==> its index in heap[]

############ SFN (tree+heap) operations

# read and parse a line into a sfn number
# returns in global variables "sfn" the read number, and str the remaining of
# the string to read. This is to avoid forking and preventing accessing heap
readsfn(){
    # global sfn str
    local s="$1"                # the string to read
    local sfn1 sfn2
    if [[ ${s:0:1} == '[' ]]; then
        readsfn "${s:1}"; sfn1="$sfn"; s="$str"
        [[ ${s:0:1} == ',' ]] || err "Syntax error: no comma: \"$s\""
        readsfn "${s:1}"; sfn2="$sfn"; s="$str"
        [[ ${s:0:1} == ']' ]] || err "Syntax error: no closing ]: \"$s\""
        heapadd "$sfn1+$sfn2"
        str="${s:1}"            # skips closing ']'
    elif [[ $s =~ ^([[:digit:]]+)(.*) ]]; then
        sfn="${BASH_REMATCH[1]}"
        str="${BASH_REMATCH[2]}"
    else
        return 1
    fi
    return 0
}

# allocate a new heap entry, or reuse one
heapadd(){
    # global sfn
    local s="$1" i               # the string to put on heap
    i="${heapidx[$s]:-}"
    if [[ -z $i ]]; then 
        i=${#heap[@]}
        heap+=("$s")
        heapidx["$s"]="$i"
    fi
    sfn="@$i"
}

compute-magnitude(){
    n="$1"
    # global magnitude
    magnitude=0
    if [[ ${n:0:1} == @ ]]; then # pair
        local sfn="${heap[${n:1}]}"
        if [[ $sfn =~ ^([^+]+)[+]([^+]+)$ ]]; then
            local lm rm lms rms
            lms="${BASH_REMATCH[1]}"
            rms="${BASH_REMATCH[2]}"
            compute-magnitude "$lms"; lm="$magnitude"
            compute-magnitude "$rms"; rm="$magnitude"
            ((magnitude = 3*lm +2*rm))
        else
            err "compute-magnitude: Heap[${n:1}] is not a pair: $sfn"
        fi
    else
        magnitude="$n"
    fi
}

############ SFS (flat strings) operations

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

############ DEBUG
DEBUG-HEAP(){
    local i
    echo -n "HEAP:"
    for((i=0; i<${#heap[@]}; i++)); do
        ((i % 10)) || echo
        echo -n "  @$i=${heap[i]}"
    done
    echo
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
        [[ $i == $j ]] && continue #  numbers must be different
        sfs="[${nums[i]},${nums[j]}]"
        reduce
        readsfn "$sfs"          # convert to sfn
        compute-magnitude "$sfn"
        ((magnitude > max)) && max="$magnitude"
    done
done
echo "Max of Magnitudes:"
echo "$max"

