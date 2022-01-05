#!/bin/bash
# https://adventofcode.com/days/day/24 puzzle #2
# See README.md in the parent directory
in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1
err(){ echo "***ERROR: $*" >&2; exit 1;}
export tmp=tmp.$$; clean(){ rm -f "$tmp" "$tmp".*;}; trap clean 0

#TEST: input 94992994195998

model="${2}"   # Optional: just check model number 

############ ALU
# We parse the NOMAD, as 14 sections (passes) of 18 instructions
# Value can be a lowercase letter (variable) literal number, or an uppercase
# letter which is the position of the input digit (A=0, B=1, ...)
# ZN is the value of Z of the previous pass #N (N={0..13})

export inp_digits='ABCDEFGHIJKLMNOPQRSTUVWXY'
declare -i inp_cur
declare -a alu_expr             # the expression for each pass
declare -a alu_eval             # Z{N+1}=alu_expr[N]
nl=$'\n'

# parsing symbolically. creates an expression alu_expr that can be evaluated
# by bash if the input number digits are set to the vars A,B, C, ...
alu-parse(){
    local -i pass i
    inp_cur=0
    for((pass=0; pass<14; pass++)); do 
        w=0; x=0; y=0; z="Z$pass";
        for((i=0; i<18; i++)); do
            read -r instr v1 v2
            "alu-$instr" "$v1" "$v2"
        done
        alu_expr+=("$z")
        alu_eval+=("Z$((pass+1))=$z")
    done <"$in"
}

# the instructions. $1 is a variable name, b a value of variable
# precedences levels: 1: * / %, 2: +, 3: ==
# shellcheck disable=SC2034 # w x y z are used by name 
alu-ini(){ inp_cur=0; w=0; x=0; y=0; z=0;}

alu-inp(){
    local -n a="$1"
    a="${inp_digits:$inp_cur:1}" # symbolic
    ((inp_cur++))
}

alu-add(){
    local -n a="$1"
    local b="$2"
    [[ $b =~ [[:lower:]] ]] && b=${!b} # expanse vars
    if [[ $b == 0 ]]; then :           # +0 => no change
    elif [[ $a == 0 ]]; then a="$b"
    elif [[ $a$b =~ [[:upper:]] ]]; then # symbolic
        if [[ $b =~ ^[-[:digit:]]+$ ]]; then
            if [[ $a =~ ^(.*)\+([[:digit:]]+)$ ]]; then
                # special case 'add ...+n m' => ...+{n+m}
                a="${BASH_REMATCH[1]}+$((BASH_REMATCH[2]+b))"
                return
            fi
        fi
        # encapsulate == terms
        if [[ $a =~ [=] ]]; then a="($a)+"; else a="$a+"; fi
        if [[ $b =~ [=] ]]; then a+="($b)"; else a+="$b"; fi
    else
        ((a+=b))
    fi
}

alu-mul(){
    local -n a="$1"
    local b="$2"
    [[ $b =~ [[:lower:]] ]] && b=${!b} # expanse vars
    if [[ $a == 0 || $b == 0 ]]; then a=0
    elif [[ $a == 1 ]]; then a="$b"
    elif [[ $b == 1 ]]; then :
    elif [[ $a$b =~ [[:upper:]] ]]; then # symbolic
        if [[ $a =~ [+=] ]]; then a="($a)*"; else a="$a*"; fi
        if [[ $b =~ [+=] ]]; then a+="($b)"; else a+="$b"; fi
    else
        ((a*=b))
    fi
}

alu-div(){
    local -n a="$1"
    local b="$2"
    [[ $b =~ [[:lower:]] ]] && b=${!b} # expanse vars
    if [[ $a == 0 || $b == 0 ]]; then a=0
    elif [[ $b == 1 ]]; then :
    elif [[ $a$b =~ [[:upper:]] ]]; then # symbolic
        if [[ $a =~ [+=] ]]; then a="($a)/"; else a="$a/"; fi
        if [[ $b =~ [+=] ]]; then a+="($b)"; else a+="$b"; fi
    else
        ((a/=b))
    fi
}

alu-mod(){
    local -n a="$1"
    local b="$2"
    [[ $b =~ [[:lower:]] ]] && b=${!b} # expanse vars
    if [[ $a$b =~ [[:upper:]] ]]; then # symbolic
        if [[ $a =~ [+=] ]]; then a="($a)%"; else a="$a%"; fi
        if [[ $b =~ [+=] ]]; then a+="($b)"; else a+="$b"; fi
    else
        ((a%=b))
    fi
}

alu-eql(){
    local -n a="$1"
    local b="$2"
    [[ $b =~ [[:lower:]] ]] && b=${!b} # expanse vars
    # special case: "eql non-negative+N D" is always false if N>9
    if [[ $b =~ ^[[:upper:]]$ ]] && [[ $a =~ ^[^-]*\+([[:digit:]]+)$ ]] &&
           ((BASH_REMATCH[1] > 9)); then
        ((a=0))
        return
    fi
    if [[ $a$b =~ [[:upper:]] ]]; then # symbolic
        a="$a==$b"              # no parentheses necessary, least precedence
     else
        ((a=(a==b)))
    fi
}

############ main
# preparse the NOMAD.
alu-parse

# shellcheck disable=SC2034 # vars used by the dynamically loaded code
declare -i Z0=0 Z1 Z2 Z3 Z4 Z5 Z6 Z7 Z8 Z9 Z10 Z11 Z12 Z13 Z14
declare -i A B C D E F G H I J K L M N

# kind of a cheat: generate a C program, as the alu_expr[] are  C expression!
# Optionally, compile & run by gcc -o d24.bin -O d24.c && ./d24.bin [12]
cat >d24.c <<EOF1
#include <stdio.h>
#include <stdlib.h>
long int A, B, C, D, E, F, G, H, I, J, K ,L, M, N;
long int Z0=0, Z1, Z2, Z3, Z4, Z5, Z6, Z7, Z8, Z9, Z10, Z11, Z12, Z13, Z14;

int exercise1(){
    printf("Exercise 1:\n");
EOF1

forloop(){
    local v=${inp_digits:$2:1}
    echo "${1}for($v=9;$v>0;$v--){${nl}$1    ${alu_eval[$2]};"
    # shellcheck disable=SC2028
    (($2 < 13)) && forloop "$1    " $(($2+1)) || 
    echo "$1    if (Z14==0) {
        ${1}printf(\"%ld%ld%ld%ld%ld%ld%ld%ld%ld%ld%ld%ld%ld%ld\n\", A, B, C, D, E, F, G, H, I, J, K ,L, M, N);
        ${1}exit(0);
    $1}"
    echo "${1}}"
}
forloop '    ' 0 >>d24.c

cat >>d24.c <<EOF2
}

int exercise2(){
    printf("Exercise 2:\n");
EOF2

forloop2(){
    local v=${inp_digits:$2:1}
    echo "${1}for($v=1;$v<10;$v++){${nl}$1    ${alu_eval[$2]};"
    # shellcheck disable=SC2028
    (($2 < 13)) && forloop2 "$1    " $(($2+1)) || 
    echo "$1    if (Z14==0) {
        ${1}printf(\"%ld%ld%ld%ld%ld%ld%ld%ld%ld%ld%ld%ld%ld%ld\n\", A, B, C, D, E, F, G, H, I, J, K, L, M, N);
        ${1}exit(0);
    $1}"
    echo "${1}}"
}
forloop2 '    ' 0 >>d24.c

cat >>d24.c <<EOF
}

int main(int argc, char **argv) {
    if (argc == 1) {
        exercise1();
    } else {
        exercise2();
    }
    exit(0);
}
EOF

# one shot mode
if [[ -n $model ]]; then
    i=0
    for var in A B C D E F G H I J K L M N; do
        declare -i "$var=${model:$((i++)):1}"
    done
    for((pass=0; pass<14; pass++)); do
        echo -n "${alu_eval[$pass]}, "
        ((alu_eval[pass]))
        var="Z$((pass+1))"
        echo "$var = ${!var}"
    done
    exit
fi

# compile C code and run it.
gcc -O -o d24.bin d24.c; ./d24.bin; rm -f d24.bin d24.c

exit 0

# for reference: Interpreted bash (too slow for production)

for((A=9;A>0;A--)); do
    ((alu_eval[0]))
    for((B=9;B>0;B--)); do
        ((alu_eval[1]))
        for((C=9;C>0;C--)); do
            ((alu_eval[2]))
            for((D=9;D>0;D--)); do
                ((alu_eval[3]))
                for((E=9;E>0;E--)); do
                    ((alu_eval[4]))
                    for((F=9;F>0;F--)); do
                        ((alu_eval[5]))
                        for((G=9;G>0;G--)); do
                            ((alu_eval[6]))
                            for((H=9;H>0;H--)); do
                                ((alu_eval[7]))
                                for((I=9;I>0;I--)); do
                                    ((alu_eval[8]))
                                    for((J=9;J>0;J--)); do
                                        ((alu_eval[9]))
                                        for((K=9;K>0;K--)); do
                                            ((alu_eval[10]))
                                            for((L=9;L>0;L--)); do
                                                ((alu_eval[11]))
                                                for((M=9;M>0;M--)); do
                                                    ((alu_eval[12]))
                                                    for((N=9;N>0;N--)); do
                                                        ((alu_eval[13]))
                                                        ((Z14)) || { echo "$A$B$C$D$E$F$G$H$I$J$K$L$M$N"; exit 0;}
                                                    done
                                                done
                                            done
                                        done
                                    done
                                done
                            done
                        done
                    done
                done
            done
        done
    done
done
