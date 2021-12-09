#!/bin/bash
# https://adventofcode.com/days/day/8 puzzle #2
# See README.md in the parent directory
in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1
err(){ echo "***ERROR: $*" >&2; exit 1;}
tmp=tmp.$$; clean(){ rm -f "$tmp" "$tmp".*;}; trap clean 0
#TEST: example-1 5353
#TEST: example 61229
#TEST: input 

max=0                           # debug: max number of lines to process

# We suppose the input always includes all digits, that means:
# one word with 2, 3, 4, and 7 chars and three words with 5 and 6 chars.
# We then first identify digits 1 and 4 since the only ones with 2 and 4 chars,
# we also identify the unique 7 and 8 this way.
# And to differentiate the 5 and 6 chars, we compare the signature of their
# intersections with 1 and 4, which allows to differentiate them
# In bash, the characters in s1 that are not in s2 are: ${s1//[$s2]/}

# we preprocess the input, yielding a word length and a word per line, in order
prepro(){
    local i c r
    for i in "$@"; do
        r=
        for c in a b c d e f g; do [[ $i =~ $c ]] && r="$r$c"; done
        echo "${#r} $r"
    done |sort -n |uniq
}

# reminder: the segment map for the display digits
# i2 = size of intersection with 2-chars (digit 1), i4 with 4-chars (digit 4)
#                      len i2 i4 
#  segments[0]=abc_efg  6  4  3 
#  segments[1]=__c__f_  2                                                    
#  segments[2]=a_cde_g  5  4  3 
#  segments[3]=a_cd_fg  5  3  2 
#  segments[4]=_bcd_f_  4       
#  segments[5]=ab_d_fg  5  4  2 
#  segments[6]=ab_defg  6  5  3 
#  segments[7]=a_c__f_  3       
#  segments[8]=abcdefg  7       
#  segments[9]=abcd_fg  6  4  2 

# we now unroll an ad-hoc algorithm for our special input
# we set the global digit array to translate chars, and we use globals
build_digits_table(){
    local i len w2 w4 word delta
    read -r len word            # 2-chars word
    [[ $len != 2 ]] && err "Expected 2-chars word, not '$word' on $line"
    digits[1]="$word"
    w2="$word"
    read -r len word            # 3-chars word
    [[ $len != 3 ]] && err "Expected 3-chars word, not '$word' on $line"
    digits[7]="$word"
    read -r len word            # 4-chars word
    [[ $len != 4 ]] && err "Expected 4-chars word, not '$word' on $line"
    digits[4]="$word"
    w4="$word"
    for i in 0 1 2; do          # the three 5-chars words
        read -r len word
        [[ $len != 5 ]] && err "Expected 5-chars word, not '$word' on $line"
        # find the actual digits by comparing intersections with digits 2 & 4
        delta=",${word//[$w2]/},${word//[$w4]/},"
        if [[ $delta =~ ,....,..., ]]; then digits[2]="$word"
        elif [[ $delta =~ ,...,.., ]]; then digits[3]="$word"
        else digits[5]="$word"
        fi
    done
    for i in 0 1 2; do          # the three 6-chars words
        read -r len word
        [[ $len != 6 ]] && err "Expected 6-chars word, not '$word' on $line"
        # find the actual digits by comparing intersections with digits 2 & 4
        delta=",${word//[$w2]/},${word//[$w4]/},"
        if [[ $delta =~ ,....,..., ]]; then digits[0]="$word"
        elif [[ $delta =~ ,.....,..., ]]; then digits[6]="$word"
        else digits[9]="$word"
        fi
    done
    read -r len word
    [[ $len != 7 ]] && err "Expected 7-chars word, not '$word' on $line"
    digits[8]="$word"
    [[ ${#digits[@]} != 10 ]] && err "Could find only ${#digits[@]} digits in $line"
}

# return the total number from its display signals into the global var total
parse_display(){
    local signal digit digitsig display
    display=
    for signal in "$@"; do
        for digit in {0..9}; do
            digitsig="${digits[$digit]}"
            if [[ -z "${signal//[$digitsig]/}" ]] &&
                [[ -z "${digitsig//[$signal]/}" ]]; then
                display="$display$digit"
                break
            fi
        done
    done
    (( total += "1$display" - 10000 )) # avoid leading 0 considered octal
}

declare -A digits                  # for each digit, its input chars
n=0
while read -r line; do
    # shellcheck disable=SC2086 # yes, we want to iterate on the words
    build_digits_table < <(prepro ${line//[^[:alnum:][:space:]]/} )
    # shellcheck disable=SC2086 # yes, we want to iterate on the words
    parse_display ${line##*[|]}
    (( max > 0 )) && (( ++n >= max )) && break # debug
done <"$in"

echo "$total"
