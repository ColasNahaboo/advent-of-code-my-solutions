#!/bin/bash
# https://adventofcode.com/days/day/6 puzzle #1
# See README.md in the parent directory
in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1
err(){ echo "***ERROR: $*" >&2; exit 1;}
tmp=tmp.$$; clean(){ rm -f "$tmp" "$tmp".*;}; trap clean 0

# This is the naive version, but it is too slow to tackle the input
# It took more than one hour for the 80 days

steps="${2:-80}"                # 2nd argument is the number of steps (def, 80)

fishes=$(tr ',' ' ' <"$in")
next=(6 0 1 2 3 4 5 6 7)

# advance once the existing fishes (olds) create the newborns (news)
# modifies in place the global space-separated list "fishes"
nextgen(){
    local olds='' news=''
    for f in $fishes; do
        olds="$olds ${next[f]}"
        [[ $f == 0 ]] && news="$news 8"
    done
    # add new ones
    fishes="$olds$news"
}

start=$(date +%s)
for ((s=0; s<steps; s++)); do
    nextgen
    now=$(date +%s)
    echo "Generation #$s took $((now - start)) seconds for $(( (${#fishes} + 1) / 2)) fishes"
    start="$now"
done

actualfishes="${fishes// /}"    # works only if internal timers is single digit
echo ${#actualfishes}
