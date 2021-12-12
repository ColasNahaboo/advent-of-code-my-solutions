#!/bin/bash
# https://adventofcode.com/2021/day/3 puzzle #2
# See README.md in the parent directory
in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1
err(){ echo "***ERROR: $*" >&2; exit 1;}
export tmp=tmp.$$; clean(){ rm -f "$tmp" "$tmp".*;}; trap clean 0

#TEST: example 230
#TEST: input 1662846

# Brute force script, so a bit slow, can take about 10 seconds.
# pre-compute dimensions of the bit array in input
line1=$(head -1 "$in")
numof_cols=${#line1}

# pass for a single column, uses $tmp
# creates temporary files tmp.0 and tmp.1 for the subsets of numbers having
# 0 or 1 in column col. Which one we use depends in the mode: OSR or CSR
# outputs the final number, or nothing if further passes needed.
filter_col(){
    local file="$1"             # data file, modified in place
    local col="$2"              # bit rating value 0..(numof_cols - 1)
    local mode="$3"             # can be OSR or CSR
    local set=0 unset=0         # the total number of set and unset bits in col
    local number                # found in the input line
    local numof_res             # the number of results
    local res                   # the resulting number, if alone
    rm -f "$tmp".0 "$tmp".1     # the numbers having 0 or 1 in col
    while read -r number; do
        # ${number:$col:1} is the bit at position "col"
        if [[ ${number:$col:1} == 1 ]]; then
            (( ++set ))
            echo "$number" >>"$tmp".1
        else
            (( ++unset ))
            echo "$number" >>"$tmp".0
        fi
    done <"$file"
    # what to do if set equals unset depends on mode
    if [[ $mode == OSR ]]; then
        if (( set >= unset )); then mv "$tmp".1 "$file"
        else mv "$tmp".0 "$file"
        fi
    else
        if (( set < unset )); then mv "$tmp".1 "$file"
        else mv "$tmp".0 "$file"
        fi
    fi
    numof_res=$(wc -l <"$file")
    if [[ $numof_res == 1 ]]; then
        res=$(cat "$file")
        echo "$((2#$res))"      # print res converted to decimal & stop
    elif [[ $numof_res == 0 ]]; then
        err "$mode found no number at col #$col"
    fi
}

# now call filter_col on all the columns, once for OSR, once for CSR

# OSR
cp "$in" "$tmp"
col=0
osr=
while (( col < numof_cols )) && [[ -z $osr ]]; do
    osr=$(filter_col "$tmp" "$col" OSR)
    (( ++col ))
done
[[ -z $osr ]] && err "Could not compute OSR"

# CSR
cp "$in" "$tmp"
col=0
csr=
while (( col < numof_cols )) && [[ -z $csr ]]; do
    csr=$(filter_col "$tmp" "$col" CSR)
    (( ++col ))
done
[[ -z $csr ]] && err "Could not compute CSR"

# compute and display LSR
(( lsr = osr * csr ))

echo "osr = $osr, csr = $csr, lsr = $lsr"
echo "$lsr"
