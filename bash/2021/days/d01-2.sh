#!/bin/bash
# https://adventofcode.com/2021/day/1 puzzle #2
# See README.md in the parent directory
in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1
export tmp=tmp.$$; clean(){ rm -f "$tmp" "$tmp".*;}; trap clean 0

# converts input to sliding sums of n1,n2,n3 
# and feed it to 1-1.sh to count the increases
n1=
n2=
n3=
while read -r n3; do
    [[ -n $n1 && -n $n2 ]] && echo $(( n1 + n2 + n3 ))
    n1="$n2"
    n2="$n3"
done <"$in" >"$tmp"

./d01-1.sh "$tmp"
