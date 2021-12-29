#!/bin/bash
# https://adventofcode.com/days/day/21 puzzle #2
# See README.md in the parent directory
in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1
err(){ echo "***ERROR: $*" >&2; exit 1;}
#export tmp=tmp.$$; clean(){ rm -f "$tmp" "$tmp".*;}; trap clean 0

#TEST: example 272847859601291
#TEST: input 444356092776315

# possible totals of 3 rolls with their frequencies (number of universes)
#rollfreq=([3]=1 [4]=3 [5]=6 [6]=7 [7]=6 [8]=3 [9]=1)
#rollfreqs=$'3 1\n4 3\n5 6\n6 7\n7 6\n8 3\n9 1'
rollfreqs='3,1 4,3 5,6 6,7 7,6 8,3 9,1'

# we play, but update the arrays of possible place,score
declare -Ai states               # place1,score1;place2,score2 => numof-worlds
declare -Ai statesn              # next-turn version of states[]
declare -i p1wins=0              # number of wins player1
declare -i p2wins=0              # number of wins player2
declare -i turns=0               # number of rounds (used only for tracing)
declare -i r1 f1 r2 f2 p1 p2 s1 s2 p1n s1n p2n s2n

# we number the places 0 to 9 with scores place+1 (1 to 10)
# shellcheck disable=SC2034 # s unused
{ read -r d d d d p1; read -r d d d d p2;} <"$in"
((p1--))
((p2--))

# starting positions
states["$p1,0;$p2,0"]=1          # starting world
echo "Start positions: $p1 $p2"

# p1, s1, p2, s2 are the starting positions and scores of players 1 & 2
# p1n, s1n, p2n, s2n the next ones after envisioned turn
# worlds is the number of worlds at the start of the turn
# p1worlds and p2worlds its value after the play of player 1 & 2

while ((${#states[@]} != 0)); do # turns
    for state in "${!states[@]}"; do
        # for each states, play one round: p1 then p2
        # note: this is 3 times faster than a [[ =~ ]] regexp match
        p1="${state%%,*}"
        t="${state%;*}"; s1="${t#*,}"
        t="${state#*;}"; p2="${t%,*}"
        s2="${state##*,}"
        worlds=${states["$state"]}
        # Player 1
        for rf1 in $rollfreqs; do
            r1="${rf1%,*}"; f1="${rf1#*,}"
            # worlds for this roll are # of worlds in this state * # freq
            ((p1worlds = worlds * f1))
            ((p1n = (p1+r1)%10))
            ((s1n = s1+p1n+1))
            if ((s1n >= 21)); then # player1 win, stop there
                ((p1wins += p1worlds))
            else
                # Player 2
                for rf2 in $rollfreqs; do
                    r2="${rf2%,*}"; f2="${rf2#*,}"
                    ((p2worlds = p1worlds * f2))
                    ((p2n = (p2+r2)%10))
                    ((s2n = s2+p2n+1))
                    if ((s2n >= 21)); then
                        ((p2wins += p2worlds))
                    else        # no win, register world for next round
                        ((statesn["$p1n,$s1n;$p2n,$s2n"] += p2worlds))
                    fi
                done
            fi
        done
    done
    ((turns++))
    # copy statesn into states
    statesn_def=$(declare -p statesn) && declare -A states="${statesn_def#*=}"
    statesn=()

    # Some traces
    nw=0; for i in "${states[@]}"; do ((nw += i)); done
    echo "==[$turns] states=${#states[@]} / $nw, p1wins=$p1wins, p2wins=$p2wins"
done

echo "p1wins=$p1wins, p2wins=$p2wins"
if ((p1wins > p2wins)); then echo "$p1wins"
else echo "$p2wins"
fi
