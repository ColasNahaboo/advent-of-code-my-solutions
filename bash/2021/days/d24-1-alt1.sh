#!/bin/bash
# https://adventofcode.com/days/day/24 playground
# See README.md in the parent directory
in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1
err(){ echo "***ERROR: $*" >&2; exit 1;}
#export tmp=tmp.$$; clean(){ rm -f "$tmp" "$tmp".*;}; trap clean 0

#TEST: input 94992994195998
#TEST part 2: 21191861151161

# working, brute force interpreted language, but much too slow

# First we suppose that input is a string of numbers: inp_buf
# inp_cur is the current position in this string
declare inp_buf
declare -i inp_cur w x y z

############ ALU
# The ALU: runs the instruction in stdin onto the input number argument
alu-run(){
    alu-ini "$1"
    while read -r instr v1 v2; do
        if [[ $v2 =~ ^[-[:digit:]]+$ ]]; then "alu-$instr-num" "$v1" "$v2"
        else "alu-$instr" "$v1" "$v2"
        fi
    done <"$in"
    #echo "ALU: input=$1, w=$w, x=$x, y=$y, z=$z"
}

# the instructions
# shellcheck disable=SC2034 # w x y z are used by name 
alu-ini(){ inp_buf="$1"; inp_cur=0; w=0; x=0; y=0; z=0;}
alu-inp(){ local -n a="$1"; a="${inp_buf:$inp_cur:1}"; ((inp_cur++));}
# when b is a variable
alu-add(){ local -n a="$1" b="$2"; ((a+=b));}
alu-mul(){ local -n a="$1" b="$2"; ((a*=b));}
alu-div(){ local -n a="$1" b="$2"; ((a/=b));}
alu-mod(){ local -n a="$1" b="$2"; ((a%=b));}
alu-eql(){ local -n a="$1" b="$2"; ((a=(a == b)));}
# when b is a number
alu-add-num(){ local -n a="$1"; ((a+="$2"));}
alu-mul-num(){ local -n a="$1"; ((a*="$2"));}
alu-div-num(){ local -n a="$1"; ((a/="$2"));}
alu-mod-num(){ local -n a="$1"; ((a%="$2"));}
alu-eql-num(){ local -n a="$1"; ((a=(a == "$2")));}

# if 2nd argument, this perform a one-shot test
if [[ -n "$2" ]]; then
    alu-run "$2"
    echo "$z"
    exit 0
fi

# brute force try of numbers satisfying the NOMAD
n=99999999999999
while true; do
    if ! [[ $n =~ 0 ]]; then
        alu-run "$n"
        echo "$n $z"
        # shellcheck disable=SC2154
        ((z)) || break
    fi
    ((n--))
done

echo "$n"
exit 0

############ Reverse-engineer input
rev-ini() {  e="$1"; vars=abcdefghijklmn; varn=13;}
rev-inp(){ local v=${vars:$varn:1}; e="${e//$1/$v}"; ((varn--));}
# when b is a variable
rev-add(){ [[ $2 != 0 ]] && e="${e//$1/($1+$2)}";}
rev-mul(){
    if [[ $2 != 0 ]]; then e="${e//$1/($1*$2)}"
    elif [[ $2 != 1 ]]; then
        e="${e//$1/0}"; e="${e//(0+/(}"; e="${e//+0)/)}"
        e="${e//(w)/w}"; e="${e//(x)/x}"; e="${e//(y)/y}"; e="${e//(z)/z}"
    fi
}
rev-div(){ e="${e//$1/($1/$2)}";}
rev-mod(){ e="${e//$1/($1%$2)}";}
rev-eql(){ e="${e//$1/($1==$2)}";}

# reverse engineer an expression
rev-run(){
    local instr v1 v2
    rev-ini "$1"
    while read -r instr v1 v2; do
        "rev-$instr" "$v1" "$v2"
        echo "$instr $v1 $v2: $e"; echo
    done < <(tac "$in")
    echo "$e"
}
