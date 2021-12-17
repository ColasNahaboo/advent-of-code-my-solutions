#!/bin/bash
# https://adventofcode.com/days/day/16 puzzle #2
# See README.md in the parent directory
in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1
err(){ echo "***ERROR: $*" >&2; exit 1;}
#export tmp=tmp.$$; clean(){ rm -f "$tmp" "$tmp".*;}; trap clean 0

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

############ Preprocess hexa input string to a bitstring
# bashism: letter-as-decimal-to-binary array: echo ${D2B[16#A]} ==> 1010
D2B=({0..1}{0..1}{0..1}{0..1})

# defines & parse global S, binary input string
read -r hexstring <"$in"
S=
for((h=0; h<${#hexstring}; h++)); do
    c=$((16#${hexstring:h:1}))
    S="$S${D2B[$c]}"
done

############ Reading and parsing packets
# we just use our readpacket as a read, reading on stdin
# reads the binary string on stdin
# returns the value of the packet read (recursively)

readpacket(){
    local s
    # shellcheck disable=SC2034 # unused values for reads
    local -i n i sublen type lentid value vtype version
    read-bitint 3 version || return 1  # EOF
    read-bitint 3 type
    if ((type == 4)); then      # value packets list, concatenate the chunks
        local valuestring
        while true; do 
            read-bitint 1 vtype
            read-bitstring 4 s
            valuestring+="$s"
            [[ $vtype == 0 ]] && break
        done
        ((value = 2#$valuestring))
    else                        # operator packet
        local subvalues=()      # array of sub-values
        read-bitint 1 lentid
        if ((lentid == 1)); then # n sub packets
            read-bitint 11 n
            for ((i=0; i<n; i++)); do
                # shellcheck disable=SC2207 
                subvalues+=($(readpacket))
            done
        else                     # sub packets fit in the next "sublen" bits
            read-bitint 15 sublen
            read-bitstring "$sublen" s # read from this substring till its EOF
            # shellcheck disable=SC2207 
            { while subvalues+=($(readpacket)); do :; done;} <<<"$s"
        fi
        [[ -z ${operators[type]} ]] && err "Invalid operator ID: $type"
        # map operator on subvalues
        value=$("${operators[type]}" "${subvalues[@]}")
    fi
    echo "$value"
    return 0
}

# read next $1 bits into $2 variable
read-bitint(){
    local -n _intvar="$2"
    local string
    read -r -N "$1" string || return 1 # EOF
    # shellcheck disable=SC2034 # unused values for reads
    ((_intvar = 2#$string))
    return 0
}

# same for string variable, for consistency / readability
read-bitstring(){
    local -n _strvar="$2"
    # shellcheck disable=SC2034 # unused values for reads
    read -r -N "$1" _strvar
}

############ operators
# The array of operators indexed by their IDs
operators=(op_sum op_product op_minimum op_maximum '' op_greater op_less op_equal)

#their code
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

############ Now, read!

readpacket <<<"$S"
