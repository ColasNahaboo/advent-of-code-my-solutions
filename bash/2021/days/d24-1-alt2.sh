#!/bin/bash
# https://adventofcode.com/days/day/24 puzzle #1
# See README.md in the parent directory
in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1
err(){ echo "***ERROR: $*" >&2; exit 1;}
export tmp=tmp.$$; clean(){ rm -f "$tmp" "$tmp".*;}; trap clean 0

#TEST: small 
#TEST: example 
#TEST: input 

# Optional: stops after N steps and output computed expression
# N=0 to print the complete esxpression and stop
steps="${2}"
# Optional: prints from steps to stop, then stop.
stop="${3}"

############ ALU
# We parse the NOMAD, but simplifying and pre-computing symbolically as much
# as possible
# Value can be a lowercase letter (variable) literal number, or an uppercase
# letter which is the position of the input digit (A=0, B=1, ...)

export inp_digits='ABCDEFGHIJKLMNOPQRSTUVWXYZ'
declare -i inp_cur
declare -i wl wh xl xh yl yh zl zh # high and low possible values of vars
nl=$'\n'
tab=$'\t'
declare -A range                # "low space high" cache of expressions

# toplevel dispatch routine: read instructions xxx and call the alu-xxx routine
# parsing symbolically. creates an expression alu_expr that can be evaluated
# by bash if the input number digits are set to the vars A,B, C, ...
alu-parse(){
    local lineno=0 v2l v2h v1l v1h
    alu-ini
    while read -r instr v1 v2; do
        if [[ $v2 =~ [[:lower:]] ]]; then # b is a variable
            v2l="${v2}l"; v2l=${!v2l}
            v2h="${v2}h"; v2h=${!v2h}
        else                    # literal is its own low&high
            v2l="$v2"
            v2h="$v2"
        fi
        "alu-$instr" "$v1" "$v2" "${v1}l" "${v1}h" "$v2l" "$v2h"
        v1l="${v1}l"; v1l=${!v1l}
        v1h="${v1}h"; v1h=${!v1h}
        range["${!v1}"]="$v1l $v1h"
        ((lineno++))
        if ((steps && lineno >= steps)); then
            # debug mode: stop and print state of the compilation
            echo "[$lineno]=\"$instr $v1 $v2\"${nl}w=$w${tab}x=$x${tab}y=$y${tab}z=$z"
            echo "$wl-$wh${tab}$xl-$xh${tab}$yl-$yh${tab}$zl-$zh"
            #local e
            #for e in "${!range[@]}"; do echo -n " $e[${range[$e]}]"; done; echo
            ((stop && lineno < stop)) || exit 0
        fi
    done <"$in"
    c_expr="$z"                 # for use in the optional C program
    alu_expr="((alu = $z ))"
}

# the instructions. $1 is a variable name, b a value of variable
# precedences levels: 1: * / %, 2: +, 3: ==
# shellcheck disable=SC2034 # w x y z are used by name 
alu-ini(){
    inp_cur=0; w=0; x=0; y=0; z=0;
    wl=0; wh=0; xl=0; xh=0; yl=0; yh=0; zl=0; zh=0;
}

alu-inp(){
    local -n a="$1" al="$3" ah="$4"
    a="${inp_digits:$inp_cur:1}" # symbolic
    ((inp_cur++))
    al=1; ah=9;                 # digits are [1-9]
}

# in input: b is the only place with negative numbers: -15 -10 -4 -2 -1
alu-add(){
    local -n a="$1" al="$3" ah="$4"
    local b="$2" bl="$5" bh="$6"
    [[ $b =~ [[:lower:]] ]] && b=${!b} # expanse vars
    if [[ $b == 0 ]]; then :           # +0 => no change
    elif [[ $a == 0 ]]; then a="$b"; ((al=bl)); ((ah=bh))
    elif [[ $a$b =~ [[:upper:]] ]]; then # symbolic
        if [[ $b =~ ^[-[:digit:]]+$ ]]; then
            if [[ $a =~ ^([[:upper:]])\+([[:digit:]]+)$ ]]; then
                # special case 'add D+n m' => D+(n+m)
                a="${BASH_REMATCH[1]}+$((BASH_REMATCH[2]+b))"
                ((al+=bl)); ((ah+=bh))
                return
            fi
        fi
        # encapsulate == terms, the only with more precedence
        if [[ $a =~ [=] ]]; then a="($a)+"; else a="$a+"; fi
        if [[ $b =~ [=] ]]; then a+="($b)"; else a+="$b"; fi
        ((al+=bl)); ((ah+=bh))
    else
        ((a+=b))
        ((al=a)); ((ah=a)) 
    fi
}

# in input: b is always 0 or x or y
alu-mul(){
    local -n a="$1" al="$3" ah="$4"
    local b="$2" bl="$5" bh="$6"
    [[ $b =~ [[:lower:]] ]] && b=${!b} # expanse vars
    if [[ $a == 0 || $b == 0 ]]; then a=0; al=0; ah=0
    elif [[ $a == 1 ]]; then a="$b"; ((al=bl)); ((ah=bh))
    elif [[ $b == 1 ]]; then :
    elif [[ $a$b =~ [[:upper:]] ]]; then # symbolic
        if [[ $b =~ ^[-[:digit:]]+$ ]]; then
            if [[ $a =~ ^([[:upper:]])\*([[:digit:]]+)$ ]]; then
                # special case 'add D*n m' => D*(n*m)
                a="${BASH_REMATCH[1]}*$((BASH_REMATCH[2]*b))"
                ((al*=bl)); ((ah*=bh))
                return
            fi
        fi
        if [[ $a =~ [+=] ]]; then a="($a)*"; else a="$a*"; fi
        if [[ $b =~ [+=] ]]; then a+="($b)"; else a+="$b"; fi
        ((al*=bl)); ((ah*=bh))  # problem with negative nums?
    else
        ((a*=b))
        ((al=a)); ((ah=a)) 
    fi
}

# in input: b is always the literal 26, or 1 (ignored)
alu-div(){
    local -n a="$1" al="$3" ah="$4"
    local b="$2" bl="$5" bh="$6" reducible=true
    [[ $b =~ [[:lower:]] ]] && b=${!b} # expanse vars
    if [[ $a == 0 ]]; then a=0; al=0; ah=0
    elif [[ $b == 1 ]]; then :
    elif [[ $a$b =~ [[:upper:]] ]]; then # symbolic, and b==26
        # reduce
        while "$reducible"; do
            reducible=false
            # [-b b] / b ==> 0
            if ((al > -b && ah < b)); then a=0; al=0; ah=0; return; fi
        done
        if [[ $a =~ [+=] ]]; then a="($a)/"; else a="$a/"; fi
        if [[ $b =~ [+=] ]]; then a+="($b)"; else a+="$b"; fi
        ((al/=bh)); ((ah/=bl))      # note the switch bl/bh
                                    # problem with negative nums?
    else
        ((a/=b))
        ((al=a)); ((ah=a)) 
    fi
}

# in input: b is always the literal 26
alu-mod(){
    local -n a="$1" al="$3" ah="$4"
    local b="$2" bl="$5" bh="$6" reducible=true ahl
    [[ $b =~ [[:lower:]] ]] && b=${!b} # expanse vars
    if [[ $a$b =~ [[:upper:]] ]]; then # symbolic
        # reduce
        while "$reducible"; do
            reducible=false
            # modulo do not change things, omit
            ((al >= 0 && ah < b)) && return
            # reduce: (((~)~)~)*26+X %26 => X%26
            if [[ $a =~ ^\(\(\(\(\(\(\(\(\([^\(\)]+\)[^\(\)]+\)[^\(\)]+\)[^\(\)]+\)[^\(\)]+\)[^\(\)]+\)[^\(\)]+\)[^\(\)]+\)[^\(\)]+\)\*26\+([^\(\)]+)$ ]] || 
                   [[ $a =~ ^\(\(\(\(\(\(\(\([^\(\)]+\)[^\(\)]+\)[^\(\)]+\)[^\(\)]+\)[^\(\)]+\)[^\(\)]+\)[^\(\)]+\)[^\(\)]+\)\*26\+([^\(\)]+)$ ]] || 
                   [[ $a =~ ^\(\(\(\(\(\(\([^\(\)]+\)[^\(\)]+\)[^\(\)]+\)[^\(\)]+\)[^\(\)]+\)[^\(\)]+\)[^\(\)]+\)\*26\+([^\(\)]+)$ ]] || 
                   [[ $a =~ ^\(\(\(\(\(\([^\(\)]+\)[^\(\)]+\)[^\(\)]+\)[^\(\)]+\)[^\(\)]+\)[^\(\)]+\)\*26\+([^\(\)]+)$ ]] || 
                   [[ $a =~ ^\(\(\(\(\([^\(\)]+\)[^\(\)]+\)[^\(\)]+\)[^\(\)]+\)[^\(\)]+\)\*26\+([^\(\)]+)$ ]] || 
                   [[ $a =~ ^\(\(\(\([^\(\)]+\)[^\(\)]+\)[^\(\)]+\)[^\(\)]+\)\*26\+([^\(\)]+)$ ]] || 
                   [[ $a =~ ^\(\(\([^\(\)]+\)[^\(\)]+\)[^\(\)]+\)\*26\+([^\(\)]+)$ ]] ||
                   [[ $a =~ ^\(\([^\(\)]+\)[^\(\)]+\)\*26\+([^\(\)]+)$ ]] ||
                   [[ $a =~ ^\([^\(\)]+\)\*26\+([^\(\)]+)$ ]] ||
                   [[ $a =~ ^[^\(\)]+\*26\+([^\(\)]+)$ ]]
            then
                a="${BASH_REMATCH[1]}"
                ahl="${range[$a]}"
                [[ -z "$ahl" ]] && err "Range not cached for \"$a\""
                al="${ahl% *}"; ah="${ahl#* }"
                reducible=true
            fi
            # reduce: (cached-expr)*26+X %26 => X%26
            if [[ $a =~ ^\(\((.+)\)\*26\)\+([^\(\)]+)$ ]] ||
                   [[ $a =~ ^\((.+)\)\*26\+([^\(\)]+)$ ]] ||
                   [[ $a =~ ^\((.+)\*26\)\+([^\(\)]+)$ ]] ||
                   [[ $a =~ ^(.+)*26\+([^\(\)]+)$ ]] 
            then
                local e
                e="${BASH_REMATCH[1]}"
                ahl="${range[$e]}" # $1 was cached, so it is a full expression
                if [[ -n "$ahl" ]]; then
                    a="${BASH_REMATCH[2]}"
                    ahl="${range[$a]}"
                    al="${ahl% *}"; ah="${ahl#* }"
                    reducible=true
                fi
            fi
           # reduce: cached-expr*26+X %26 => X%26
            if [[ $a =~ ^(.+)\)\*26\+([^\(\)]+)$ ]]; then
                local e
                e="${BASH_REMATCH[1]}"
                ahl="${range[$e]}" # $1 was cached, so it is a full expression
                if [[ -n "$ahl" ]]; then
                    a="${BASH_REMATCH[2]}"
                    ahl="${range[$a]}"
                    al="${ahl% *}"; ah="${ahl#* }"
                    reducible=true
                fi
            fi
             # reduce: Y*26+X %26 => X%26
            if [[ $a =~ ^[^\(\)]+\*26\+([^\(\)]+)$ ]]; then
                a="${BASH_REMATCH[1]}"
                ahl="${range[$a]}"
                [[ -z "$ahl" ]] && err "Range not cached for \"$a\""
                al="${ahl% *}"; ah="${ahl#* }"
                reducible=true
            fi
        done
        if [[ $a =~ [+=] ]]; then a="($a)%"; else a="$a%"; fi
        if [[ $b =~ [+=] ]]; then a+="($b)"; else a+="$b"; fi
        al=0; ah=$((b-1))
    else
        ((a%=b))
        ((al=a)); ((ah=a)) 
    fi
}

# in input: b is always w or 0
alu-eql(){
    local -n a="$1" al="$3" ah="$4"
    local b="$2" bl="$5" bh="$6"
    [[ $b =~ [[:lower:]] ]] && b=${!b} # expanse vars
    # special case: "eql non-negative+N D" is always false if N>9
    if [[ $b =~ ^[[:upper:]]$ ]] && [[ $a =~ ^[^-]*\+([[:digit:]]+)$ ]] &&
           ((BASH_REMATCH[1] > 9)); then
        ((a=0)); ((al=0)); ((ah=0))
        return
    fi
    if [[ $a$b =~ [[:upper:]] ]]; then # symbolic
        if ((ah < bl || bh < al)); then # empty intersection
            a=0; al=0; ah=0
            return
        fi
        a="$a==$b"              # no parentheses necessary, least precedence
        al=0; ah=1
    else
        ((a=(a==b)))
        ((al=a)); ((ah=a)) 
    fi
}

############ main


# preparse the NOMAD.
alu-parse

# kind of a cheat: generate a C program, as the alu_expr is a C expression!
# Optionally, compile & run by gcc -O d24.c && ./a,out
cat >d24.c <<EOF
#include <stdio.h>
#include <stdlib.h>
long int A, B, C, D, E, F, G, H, I, J, K ,L, M, N;
long int alu;
long int n = 22876792454960; // 9#88888888888888
int main(int argc, char **argv) {
    while (n >=0) {
	N=(n%9)+1;
	M=(n/9%9)+1;
	L=(n/81%9)+1;
	K=(n/729%9)+1;
	J=(n/6561%9)+1;
	I=(n/59049%9)+1;
	H=(n/531441%9)+1;
	G=(n/4782969%9)+1;
	F=(n/43046721%9)+1;
	E=(n/387420489%9)+1;
	D=(n/3486784401%9)+1;
	C=(n/31381059609%9)+1;
	B=(n/282429536481%9)+1;
	A=(n/2541865828329%9)+1;

        alu = $c_expr;

	//printf("9#%ld %ld\n", n, alu);
	if (alu==0) {
	    printf("9#%ld\n", n);
	    exit(0);
	}
	n--;
    }
}
EOF

if [[ -n $steps ]]; then
    echo "$alu_expr"
    echo "alu_expr is ${#alu_expr} chars"
    exit 0
fi

# Create a bash function "alu-eval" which evaluate the ALU into the file $tmp,
# and source this file to define it here.
# shellcheck disable=SC2016 # yes, we do not want to expand the $
echo 'alu-eval(){
local -i A="$1" B="$2"  C="$3"  D="$4"  E="$5"  F="$6"  G="$7"  H="$8"
local -i I="$9" J="${10}" K="${11}" L="${12}" M="${13}" N="${14}"
' >$tmp
echo "$alu_expr;}" >>$tmp
# shellcheck disable=SC1090
. $tmp
cp $tmp /tmp/alu-eval           # DDD

((n=9#100000000000000))         # we work in base 9 [0-8], and +1 to get [1-9]
while ((--n)); do
    # split the base-9 digits into arguments, with 1 added
    alu-eval $((n/2541865828329%9+1)) $((n/282429536481%9+1)) $((n/31381059609%9+1)) $((n/3486784401%9+1)) $((n/387420489%9+1)) $((n/43046721%9+1)) $((n/4782969%9+1)) $((n/531441%9+1)) $((n/59049%9+1)) $((n/6561%9+1)) $((n/729%9+1)) $((n/81%9+1)) $((n/9%9+1)) $((n%9+1))
    # shellcheck disable=SC2154 # alu-eval defines alu2, sourced dynamically
    echo "$((n/2541865828329%9+1))$((n/282429536481%9+1))$((n/31381059609%9+1))$((n/3486784401%9+1))$((n/387420489%9+1))$((n/43046721%9+1))$((n/4782969%9+1))$((n/531441%9+1))$((n/59049%9+1))$((n/6561%9+1))$((n/729%9+1))$((n/81%9+1))$((n/9%9+1))$((n%9+1)) $alu"
    # shellcheck disable=SC2154
    ((alu)) || break
done

echo "$n"

