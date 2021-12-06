#!/bin/bash
# https://adventofcode.com/days/day/6 puzzle #1
# See README.md
in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1
err(){ echo "***ERROR: $*" >&2; exit 1;}
tmp=tmp.$$; clean(){ rm -f "$tmp" "$tmp".*;}; trap clean 0

steps="${2:-80}"                # 2nd argument is the number of steps (def, 80)

: COMMENT <<ENDCOMMENT
The naive version (d06-other1.sh) takes too long (1h), so we use another
approach: We build the table of generations sideways: the days are columns, e.g.
instead of:

Initial state: 3,4,3,1,2
After  1 day:  2,3,2,0,1
After  2 days: 1,2,1,6,0,8
After  3 days: 0,1,0,5,6,7,8
After  4 days: 6,0,6,4,5,6,7,8,8

We use:

  --Days-->
|  3 2 1 0 6
F  4 3 2 1 6
i  3 2 1 0 6
s  1 0 6 5 4
h  2 1 0 6 5
e  . . 8 7 6
s  . . . 8 7
|  . . . . 8
V  . . . . 8

This way, since the total number of days is known (80), we can create
the first line for the 80 days at once, which is a simple repetition of a 
countdown. And repeat for all initial fished to obtain the base matrix, the
first 5 lines here.

We then can add newborns of each day as new lines based on the
number of zeros in the column of the current day.

The number of fishes is thus the number of lines.

We build the table of generations into the $tmp file

This allows us to run in 1 mn instead of 1 hour for the naive version.

But it is not yet satisfactory, as it is much faster, but still exponential.

ENDCOMMENT

initial_fishes=$(tr ',' ' ' <"$in")

# we build the filler sequence for lines, starts with 87654321065432106...
sequence=87; n=6; len=$((steps + 7))
for ((i=0; i<=len; i++)); do
    sequence="$sequence$n"
    (( --n < 0 )) && n=6
done

# for a fish (line), complete its timers until step
# globals: day steps
fish_timers(){
    local timer="$1"
    # fills with the sequence of timers up to steps
    echo "${sequence:$((8 - timer)):$((steps - day))}"
}

# First, the initial fishes
day=0
for timer in $initial_fishes; do
    fish_timers "$timer" >>$tmp
done

start=$(date +%s)
# then for each day, add newborns of the day as new lines: "0" in column "day"
while (( day < steps)); do
    newborns=$(grep -Pc "^.{$day}\K0" "$tmp")
    now=$(date +%s)
    echo "Day #$day took $((now - start)) seconds, $newborns newborns"
    (( ++day ))
    for ((newborn=0; newborn<newborns; newborn++)); do
        for ((i=1; i<=day; i++)); do echo -n "."; done # leading dots
        fish_timers 8
    done >>$tmp
done

wc -l <$tmp

