#!/bin/bash
# https://adventofcode.com/days/day/6 puzzle #1
# See README.md
in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1
err(){ echo "***ERROR: $*" >&2; exit 1;}
tmp=tmp.$$; clean(){ rm -f "$tmp" "$tmp".*;}; trap clean 0

days="${2:-80}"                # 2nd argument is the number of days (def, 80)

: COMMENT <<ENDCOMMENT
We will not actually build the generation table. We will just maintain an array
counting the number of newborns during a day, that the generation of
fish_timers will update when completing virtually a timeline of a fish

When we add a newborn fish, we virtually unroll its timer, but just increment
born where zeros would have been:

<--- day --->8765432106543210654321
             <timer->|      |
=====================+======+====== born[]

Also, we increment by the number of newborns on this day, no need to process
them individually.

We then have make the algorithm linear and no more exponential

ENDCOMMENT

initial_fishes=$(tr ',' ' ' <"$in") # list of fish timers at day 0
declare -a born                 # array of number of newborns per day
for ((i=0; i<days; i++)); do born[i]=0; done

# for a fish (column), complete its timers until days, but only update born[]
# globals: day days fishes
fish_timers(){
    local timer="$1"
    local n="${2:-1}"           # number of newborns to add
    local i
    [[ $n == 0 ]] && return
    for ((i = timer + day; i < days; i += 7)); do
        (( born[i] += n ))
    done
    (( fishes += n ))
}

# First, the initial fishes
day=0                           # current day
fishes=0                        # total number of fishes on this day
for timer in $initial_fishes; do fish_timers "$timer"; done

# then for each day (row), add newborns of the day as new lines
while (( day < days)); do
    newborns="${born[day]}"
    (( ++day ))
    fish_timers 8 "$newborns"
done

echo "$((fishes))"
