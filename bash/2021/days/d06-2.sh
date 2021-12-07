#!/bin/bash
# https://adventofcode.com/days/day/6 puzzle #2
# See README.md in the parent directory
in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1
err(){ echo "***ERROR: $*" >&2; exit 1;}
tmp=tmp.$$; clean(){ rm -f "$tmp" "$tmp".*;}; trap clean 0
#TEST: input 1650309278600
#TEST: example 26984457539

# The second problem is just the first, with steps 256 instead of 80
steps="${2:-256}"                # 2nd argument is the number of steps (def, 80)

"${0//d06-2/d06-1}" "$in" "$steps"
