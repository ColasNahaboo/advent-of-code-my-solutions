#!/bin/bash
# https://adventofcode.com/days/day/16 puzzle #2
# See README.md in the parent directory
in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1
err(){ echo "***ERROR: $*" >&2; exit 1;}
export tmp=tmp.$$; clean(){ rm -f "$tmp" "$tmp".*;}; trap clean 0

#TEST: example01 2021 
#TEST: example21 3
#TEST: example22 54
#TEST: example23 7
#TEST: example24 9
#TEST: example25 1
#TEST: example26 0
#TEST: example27 0
#TEST: example28 1
#TEST: input 1725277876501

# bashism: letter-as-decimal-to-binary array: echo ${D2B[16#A]} ==> 1010
D2B=({0..1}{0..1}{0..1}{0..1})

# defines&parse global S, biunary input string
read -r hexstring <"$in"
S=
for((h=0; h<${#hexstring}; h++)); do
    c=$((16#${hexstring:h:1}))
    S="$S${D2B[$c]}"
done
operators=(op_sum op_product op_minimum op_maximum '' op_greater op_less op_equal)

# To work around the bash limitation of not having rw globals across
# subprocesses forked by $(...) for the accessing state of the input, and
# we just use our readpacket as aread, readin on stdin
# And we fake a multiple return value to return both the number of bits
# read and the value as a string  with both of them concatenated around a colon

# reads the binary string on stdin
# returns number of bits read + ":" + value of the packet read (recursively)
readpacket(){
    local s res
    local -i n i len type lentid value vtype l=0
    read -r -N 3 s; ((l+=3))
    # ((version = 2#$s))
    read -r -N 3 s; ((l+=3))
    ((type = 2#$s))
    if ((type == 4)); then      # value packets list
        local valuestring
        while true; do 
            read -r -N 1 vtype; ((l++))
            read -r -N 4 s; ((l+=4))
            valuestring+="$s"
            [[ $vtype == 0 ]] && break
        done
        ((value = 2#$valuestring))
    else                        # operator packet
        local subvalues=()      # array of sub-values
        read -r -N 1 lentid; ((l++))
        if ((lentid == 1)); then # N sub packets
            read -r -N 11 s; ((l+=11))
            ((n = 2#$s)) # 11-bit number
            for ((i=0; i<n; i++)); do
                res="$(readpacket)"
                ((l+="${res%%:*}"))
                subvalues+=("${res#*:}")
            done
        else                     # sub packets in S substrings
            read -r -N 15 s; ((l+=15))
            ((len = 2#$s)) # 15-bit number
            while ((len >0)); do
                res="$(readpacket)"
                ((l+="${res%%:*}"))
                ((len-="${res%%:*}"))
                subvalues+=("${res#*:}")
            done
        fi
        [[ -z ${operators[type]} ]] && err "Invalid operator ID: $type"
        value=$("${operators[type]}" "${subvalues[@]}")
    fi
    echo "$l:$value"
    return 0
}

op_sum(){
    local -i i v=0
    for i in "$@"; do ((v += i)); done
    echo "$v"
}

op_product(){
    local -i i v=1
    for i in "$@"; do ((v *= i)); done
    echo "$v"
}

op_minimum(){
    local -i i v=8888888888888888888
    for i in "$@"; do ((i < v)) && ((v=i)); done
    echo "$v"
}

op_maximum(){
    local -i i v=-8888888888888888888
    for i in "$@"; do ((i > v)) && ((v=i)); done
    echo "$v"
}

op_greater(){
    (("$1" > "$2")) && echo 1 || echo 0
}

op_less(){
    (("$1" < "$2")) && echo 1 || echo 0
}

op_equal(){
    (("$1" == "$2")) && echo 1 || echo 0
}

value=$(readpacket <<<"$S")
echo "${value#*:}"
