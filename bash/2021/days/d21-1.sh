#!/bin/bash
# https://adventofcode.com/days/day/21 puzzle #1
# See README.md in the parent directory
in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1
err(){ echo "***ERROR: $*" >&2; exit 1;}
#export tmp=tmp.$$; clean(){ rm -f "$tmp" "$tmp".*;}; trap clean 0

#TEST: example 739785
#TEST: input 551901

# we number the places 0 to 9 with scores place+1 (1 to 10)
declare -i p1 p2 s1=0 s2=0
# shellcheck disable=SC2034 # s unused
{ read -r d d d d p1; read -r d d d d p2;} <"$in"

((--p1))
((--p2))

declare -i dicerolls=0
declare -i diceval=0
dice_roll(){
    local i
    dice_roll=0
    for((i=0; i<3; i++)); do
        ((++diceval > 100)) && ((diceval-=100))
        ((dice_roll += diceval))
    done
    ((dicerolls+=3))
}

declare -i sloser=0
while true; do                  # turns
    # Player 1
    dice_roll
    ((p1+=dice_roll))
    ((p1 = p1 % 10))
    ((s1 += (p1 + 1)))
    ((s1 >= 1000)) && { sloser="$s2"; break;}
    # Player 2
    dice_roll
    ((p2+=dice_roll))
    ((p2 = p2 % 10))
    ((s2 += (p2 + 1)))
    ((s2 >= 1000)) && { sloser="$s1"; break;}
done
echo "$sloser * $dicerolls ="
echo $((sloser * dicerolls))
