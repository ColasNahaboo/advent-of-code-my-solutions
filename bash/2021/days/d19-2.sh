#!/bin/bash
# shellcheck disable=SC2206,SC2086,SC2046 # we rely on params expansion
# https://adventofcode.com/days/day/19 puzzle #2
# See README.md in the parent directory
in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1
err(){ echo "***ERROR: $*" >&2; exit 1;}
#export tmp=tmp.$$; clean(){ rm -f "$tmp" "$tmp".*;}; trap clean 0
D(){ echo "$*";}

#TEST: small 0
#TEST: example 3621
#TEST: input 10569

# optional second arguments: miniminum common distances to start a match
min_common=${2:-12}

# read scanners into arrays scanner0, scanner1, ... scanner$n_scan

n_scan=0
while true; do
    n=0    
    while read -r line; do
        [[ $line =~ ^--- ]] && continue
        [[ $line =~ ^[[:space:]]*$ ]] && break
        # shellcheck disable=SC1087
        declare "scanner$n_scan[$n]=${line//,/ }"
        ((++n))
    done
    ((n)) || break
    ((++n_scan))
done <"$in"

# compute an unique distance into global "d"
# this distance should uniquely identify a vector, but be invariant
# via rotation or offset
distance(){
    local -i d1="$1" d2="$2" d3="$3"
    if (($4 > d1)); then ((d1=$4-d1)); else ((d1-=$4)); fi
    if (($5 > d2)); then ((d2=$5-d2)); else ((d2-=$5)); fi
    if (($6 > d3)); then ((d3=$6-d3)); else ((d3-=$6)); fi
    d="$((d1+d2+d3)).$((d1**2+d2**2+d3**2)).$((d1*d2*d3))"
}

# only the manhattan distance, returned in global "d"
manhattan-distance(){
    local -i d1="$1" d2="$2" d3="$3"
    if (($4 > d1)); then ((d1=$4-d1)); else ((d1-=$4)); fi
    if (($5 > d2)); then ((d2=$5-d2)); else ((d2-=$5)); fi
    if (($6 > d3)); then ((d3=$6-d3)); else ((d3-=$6)); fi
    d="$((d1+d2+d3))"
}

# computes distances between all the readings of a scannerN[] into distancesN[]
# key=distance, value="i j" indexes of beacons in scannerN
scanner-distances(){
    local -n scanner="$1" distances="$2"
    local -i i j
    for((i=0; i<${#scanner[@]}; i++)); do
        for((j=i+1; j<${#scanner[@]}; j++)); do
            distance ${scanner[i]} ${scanner[j]}
            # shellcheck disable=SC2034
            distances["$d"]="$i $j"
        done
    done
}

# number of common distances between two scanners, returned in scanners_common
scanners_common=
scanners-common(){
    local -n dist1="$1" dist2="$2"
    local d
    scanners_common=0
    for d in "${!dist2[@]}"; do
        [[ ${dist1[$d]} ]] && ((++scanners_common))
    done
}

# returns in scanner_bestmatch the N of known scanner with most common pairs
declare -i scanner_bestmatch
scanner-bestmatch(){
    local -i i="$1" j n=0 max=0
    local -n dist1="$2" dist2="$3"
    local d
    scanner_bestmatch=
    for((j=0; j<n_scan; j++)); do
        ((i == j)) && continue
        ((scanners_todo[j])) && continue # skip yet-unfixed ones
        scanners-common "distances$i" "distances$j"
        if ((scanners_common>max)); then
            ((max=scanners_common))
            ((scanner_bestmatch=j))
        fi
    done
    # Less than 12 common distances, skip
    ((max >= min_common))
}

# position scanner n relative to scanner n0
scanner-match(){
    local -i n="$1" n0="$4" i j 
    local -n scan="$2" dist="$3" scan0="$5" dist0="$6"
    local d b1 b2 b01 b02 o b
    local rots rot maxrot=0 pairs

    # take the most distances matches for each possible rotation,
    # to avoid false positives
    for d in "${!dist[@]}"; do
        if [[ ${dist0[$d]} ]]; then
            if rotated $(deltas "scan0" ${dist0[$d]}) $(deltas "scan" ${dist[$d]}); then
                ((rots[rotated]++))
                if ((rots[rotated]>maxrot)); then
                    # these are distance+pairs we are sure are the good ones
                    ((maxrot=rots[rotated]))
                    ((rot=rotated))
                    pairs=(${dist0[$d]} ${dist[$d]})
                fi
            fi
        fi
    done
    ((maxrot)) || return 1    # no matches
    
    # offset: take the 2 possible ones, and apply to the other to choose one
    b1=$(rotate ${scan[${pairs[2]}]} $rot)
    b2=$(rotate ${scan[${pairs[3]}]} $rot)
    b01=${scan0[${pairs[0]}]}
    b02=${scan0[${pairs[1]}]}
    o=$(sub3 $b1 $b01)
    b=$(sub3 $b2 $o)
    if ! equ3 $b $b02; then
        o=$(sub3 $b1 $b02)
        b=$(sub3 $b2 $o)
        equ3 $b $b01 || err "Could not determine offset"
    fi

    echo "  ==>  Rotation+Offset of scanner$n from scanner$n0: $rot / $o"
    # normalize the coordinates in scannerN to coords of scanner0, mark as done
    scanner-fix scan $rot $o
    scanpos[n]="$o"
    return 0
}

# fix all beacons coordinates on place in scanner with rotation and offset
scanner-fix(){
    local -n scantofix="$1"
    local -i rot="$2" x="$3" y="$4" z="$5" i
    local b
    local scanold=("${scantofix[@]}")
    for((i=0; i<${#scanold[@]}; i++)); do
        b=$(rotate ${scanold[i]} $rot)
        scantofix[i]=$(sub3 $b $x $y $z)
    done
}

# triplets arithmetic
add3(){ echo $(($1+$4)) $(($2+$5)) $(($3+$6));}
sub3(){ echo $(($1-$4)) $(($2-$5)) $(($3-$6));}
equ3(){ (($1==$4)) && (($2==$5)) && (($3==$6));}
    
# print the difference between two beacons of a scanner
deltas(){
    local -n scandelta="$1"
    local -i i="$2" j="$3"
    local -i x1 y1 z1 x2 y2 z2
    read -r x1 y1 z1 <<<"${scandelta[i]}"
    read -r x2 y2 z2 <<<"${scandelta[j]}"
    echo "$((x2-x1)) $((y2-y1)) $((z2-z1))"
}

# rotate x y z by N: (0..23) possible rotations
rotate(){
    local -i x="$1" y="$2" z="$3" n="$4"
    case "$n" in
        0) echo $x $y $z;;      # X points to old X, + the 3 X-rotations
        1) echo $x $((-z)) $y;;
        2) echo $x $((-y)) $((-z));;
        3) echo $x $z $((-y));;
        
        4) echo $((-x)) $y $((-z));; # X to -X
        5) echo $((-x)) $z $y;;
        6) echo $((-x)) $((-y)) $z;;
        7) echo $((-x)) $((-z)) $((-y));;
         
        8) echo $((-y)) $x $z;; # X to Y
        9) echo $z $x $y;;
        10) echo $y $x $((-z));;
        11) echo $((-z)) $x $((-y));;
         
        12) echo $y $((-x)) $z;; # X to -Y
        13) echo $z $((-x)) $((-y));;
        14) echo $((-y)) $((-x)) $((-z));;
        15) echo $((-z)) $((-x)) $y;;

        16) echo $y $z $x;;      # X to Z
        17) echo $((-z)) $y $x;;
        18) echo $((-y)) $((-z)) $x;;
        19) echo $z $((-y)) $x;;

        20) echo $z $y $((-x));;      # X to -Z
        21) echo $y $((-z)) $((-x));;
        22) echo $((-z)) $((-y)) $((-x));;
        23) echo $((-y)) $z $((-x));;

        *) err "bad rotation: $n";;
    esac
}

# if $x1 $y1 $z1 are $x2 $y2 $z2 rotated by r, returns true
# stores found r value in global "rotated"
rotated(){
    local -i x1="$1" y1="$2" z1="$3" x2="$4" y2="$5" z2="$6" r
    for((r=0; r<24; r++)); do
        [[ "$x1 $y1 $z1" != $(rotate $x2 $y2 $z2 $r) ]] && continue
        rotated="$r"
        return 0
    done
    rotated=-1
    return 1
}

############ Main

scanpos=('0 0 0')               # positions of scanners

# compute per-scanner data: scannerN[] and distancesN[]
for((i=0; i<n_scan; i++)); do
    declare -a "scanner$i"
    declare -A "distances$i"
    scanner-distances "scanner$i" "distances$i" "$i"
done

# The list of to-be-positioned scanners. 0 is the base, so known.
scanners_todo[0]=0
for((i=1; i<n_scan; i++)); do scanners_todo[i]=1; done

# Now, for all todo scanners, find the known ones with most common pairs
# Deduce then the rotation, and then the offset.
# fix them in place
fixed=true
while [[ ${scanners_todo[*]} =~ 1 ]]; do
    "$fixed" || err "could not fix remaining scanners!"
    fixed=false
    for((i=1; i<n_scan; i++)); do
        ((scanners_todo[i])) || continue
        scanner-bestmatch "$i" "scanner$i" "distances$i" || continue
        j="$scanner_bestmatch"
        scanner-match "$i" "scanner$i" "distances$i" \
                      "$j" "scanner$j" "distances$j" || continue
        scanners_todo[i]=0
        fixed=true
    done
done

# find the biggest manhattan distances
max=0
for((i=0; i<n_scan; i++)); do
    for((j=i+1; j<n_scan; j++)); do
        manhattan-distance ${scanpos[i]} ${scanpos[j]}
        ((d>max)) && ((max=d))
    done
done
echo $max
