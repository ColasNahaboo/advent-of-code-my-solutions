#!/bin/bash
# https://adventofcode.com/days/day/17 puzzle #2
# See README.md in the parent directory
in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1
err(){ echo "***ERROR: $*" >&2; exit 1;}
#export tmp=tmp.$$; clean(){ rm -f "$tmp" "$tmp".*;}; trap clean 0

#TEST: example 112
#TEST: input 1994

############ Parse input
read -r input <"$in"
[[ $input =~ 'target area: 'x=([-[:digit:]]+)[.]+([-[:digit:]]+),[[:space:]]*y=([-[:digit:]]+)[.]+([-[:digit:]]+) ]] || err "Input syntax error: $input"
x1="${BASH_REMATCH[1]}"
x2="${BASH_REMATCH[2]}"
y1="${BASH_REMATCH[3]}"
y2="${BASH_REMATCH[4]}"
echo "area: [$x1 $x2] x [$y1 $y2]"

############ Compute bounds
# for x, we know that vx0 must be at least 1 (otherwise we do not move at all)
# and at most x2 (otherwise we overshoot)
# for y, we know from part1 that vy0 must be at most (-y1 -1)
# and the minimum is y1, otherwise we overshoot.

# hack: compute the minimal vx by looking at a probe just reaching x1
# and retracing back its steps
for((x=x1, vx=0; x>0; x+=vx, vx--)); do :; done
((minvx0 = -vx -1))
echo "minvx0 = $minvx0"
rm -f /tmp/R

vok=0                           # number of OK velocities

for((vx0=minvx0; vx0 <= x2; vx0++)); do
    for((vy0=y1; vy0 <= -y1 - 1; vy0++)); do
        x=0; y=0; vx="$vx0"; vy="$vy0" # starts a new launch
        while true; do                 # trajectory steps
            ((x+=vx))
            ((y+=vy))
            ((vx > 0)) && ((vx--)) # velocity X reached 0, stopped
            ((vy--))
            (((x < x1) || (y > y2))) && continue # step is not yet in target
            (((x > x2) || (y < y1))) && break    # overshoot
            ((++vok))                                        # in target
            break
        done
    done
done

echo "$vok"
