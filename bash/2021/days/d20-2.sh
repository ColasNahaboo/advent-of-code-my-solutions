#!/bin/bash
# shellcheck disable=SC2206
# https://adventofcode.com/days/day/20 puzzle #2
# See README.md in the parent directory
in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1
err(){ echo "***ERROR: $*" >&2; exit 1;}
#export tmp=tmp.$$; clean(){ rm -f "$tmp" "$tmp".*;}; trap clean 0

#TEST: example 
#TEST: input 

steps="${2:-50}"                # number of steps. Was 2 fo exercise 1
pad="${3}"                      # number of padding rows/cols around the image
((pad)) || ((pad=steps*2))      # should be number of steps*2

in_image=false
declare -ai iea image0 image image2
declare -i cols0 rows0 cols rows i j k x y area pad
rows0=0
while read -r line; do
    if "$in_image"; then
        image0+=($line)
        ((rows0++))
    else
        if [[ -z $line ]]; then
            in_image=true
        else
            iea+=($line)
        fi
    fi
done < <(sed -e 's/[.]/0 /g' -e 's/#/1 /g' "$in")
((${#iea[@]} != 512)) && err "IEA size is not 512 but: ${#iea[@]}"
((cols0 = ${#image0[@]} / rows0))

# image is the read image, image0, padded with pad rows of 0
((cols = cols0 + 2*pad))
((rows = rows0 + 2*pad))
((area = rows*cols))
for((i=0; i < pad*cols; i++)); do image+=(0); done
for((i=0; i<rows0; i++)); do
    for((k=0; k<pad; k++)); do image+=(0); done
    for((j=0; j<cols0; j++)); do
        image+=(${image0[j + i*cols0]})
    done
    for((k=0; k<pad; k++)); do image+=(0); done
done
for((i=0; i < pad*cols; i++)); do image+=(0); done

DI(){
    local -n I="$1"
    local c r
    for((r=0; r<rows; r++)); do
        for((c=0; c<cols; c++)); do
            echo -n "${I[c +r*cols]}"
        done
        echo
    done
    echo
}
#DI image

echo "Processing $cols x $rows image in $steps steps"

# we reduce the computed area by 1 on each step, outer rings become irrelevant
for((step=1; step <= steps; step++)); do
    for((i=0; i<area; i++)); do image2[i]=0; done # clear copy
    for((x=step; x<(cols-step); x++)); do
        for((y=step; y<(rows-step); y++)); do
            ((i = image[x+1+(y+1)*cols] + 2*image[x+(y+1)*cols] + 4*image[x-1+(y+1)*cols] +
              8*image[x+1+y*cols] + 16*image[x+y*cols] + 32*image[x-1+y*cols] +
              64*image[x+1+(y-1)*cols] + 128*image[x+(y-1)*cols] + 256*image[x-1+(y-1)*cols]
             ))
            ((iea[i])) && image2[x+y*cols]=1
        done
    done
    image=("${image2[@]}")
done
#DI image2

echo "Done"

# omit outer (# of steps) rows in computation
declare -i lit=0
for((x=steps; x<(cols-steps); x++)); do
    for((y=steps; y<(rows-steps); y++)); do
        ((image2[x+y*cols])) && ((++lit))
    done
done

echo "$lit"

