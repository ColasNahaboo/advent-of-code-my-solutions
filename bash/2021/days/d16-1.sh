#!/bin/bash
# https://adventofcode.com/days/day/16 puzzle #1
# See README.md in the parent directory
in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1
err(){ echo "***ERROR: $*" >&2; exit 1;}
#export tmp=tmp.$$; clean(){ rm -f "$tmp" "$tmp".*;}; trap clean 0

# Warning: this code is quite crude, d16-2 has much cleaner logic.

#TEST: example01 6
#TEST: example02 9
#TEST: example03 14
#TEST: example1 16
#TEST: example2 12
#TEST: example3 23
#TEST: example4 31
#TEST: input 927

# bashism: letter-as-decimal-to-binary array: echo ${D2B[16#A]} ==> 1010
D2B=({0..1}{0..1}{0..1}{0..1})

read -r hexstring <"$in"
s=
for((h=0; h<${#hexstring}; h++)); do
    c=$((16#${hexstring:h:1}))
    s="$s${D2B[$c]}"
done
echo "$hexstring = $s (${#s})"  # DEBUG
versionsum=0

# uses global s
# returns in globals: $1 version type
# increments global versionsum
readpacket(){
    local p="$1"               # position in s, returns new one in "$readpacket"
    local -i n i len end p0="$p" version type lentid
    ((p > (${#s} - 5) )) && echo "  @[$p0]: EXIT" && return 1 # at the end (+4bits?)
    ((version = 2#${s:$p:3})); ((p+=3))
    ((versionsum += version))
    ((type = 2#${s:$p:3})); ((p+=3))
    echo "  @[$p0]: v=$version, type=$type"
    if ((type == 4)); then      # value packets list
        while true; do
            [[ ${s:p:1} == 0 ]] && break
            ((p+=5))
        done
        ((p+=5))
    else                        # operator packet
        ((lentid = 2#${s:$p:1})); ((p++))
        if ((lentid == 1)); then # N sub packets
            ((n = 2#${s:$p:11})); ((p+=11)) # 11-bit number
            echo "       OP, $n subpackets"
            for ((i=0; i<n; i++)); do
                readpacket "$p"
                p="$readpacket" || return 1
            done
        else                     # sub packets in S substrings
            ((len = 2#${s:$p:15})); ((p+=15)) # 15-bit number
            ((end = p+len))
            echo "       OP, subpackets in next $len bits till $end"
            while ((p < end)); do
                readpacket "$p"
                p="$readpacket" || return 1
            done
        fi
    fi
    readpacket="$p"
    return 0
}

readpacket 0
#    readpacket="$(((p+4)/4*4))" # skip trailing garbage till next 4-bit pos
echo "$versionsum"
