#!/bin/bash
# https://adventofcode.com/days/day/24 puzzle #2
# See README.md in the parent directory
in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1
err(){ echo "***ERROR: $*" >&2; exit 1;}
export tmp=tmp.$$; clean(){ rm -f "$tmp" "$tmp".*;}; trap clean 0

#TEST: input 21191861151161

# run d24-1.sh to generate the C program
d24-1.sh d24.input 21191861151161

# run it
gcc -O -o d24.bin d24.c; ./d24.bin 2; rm -f d24.bin d24.c
