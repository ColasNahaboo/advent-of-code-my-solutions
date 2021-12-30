#!/bin/bash
# https://adventofcode.com/days/day/22 puzzle #2
# See README.md in the parent directory
in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1
err(){ echo "***ERROR: $*" >&2; exit 1;}
#export tmp=tmp.$$; clean(){ rm -f "$tmp" "$tmp".*;}; trap clean 0

#TEST: small 39
#TEST: example 39769202357779
#TEST: input 1201259791805392

# This is based on d22-1-alt1.sh, where we do not explicitely
# construct the X Y Z coordinate system, and do not filter the inputs
# by the interval [-50,50]
# It works, but alas it is too slow for d22.input, as it would take a full day
# process the 569034720 "points" of the XYZ 828x830x828 space and 420 layers


nl=$'\n'
declare -ai Xx Yy Zz            # conversion tables
declare -Ai xX yY zZ            # index can be negative, so associtaive array
declare -ai Xsize Ysize Zsize   # the distance till next coord 
# layers: type (on=1, off=0), x y z
declare -ai layt layx1 layx2 layy1 layy2 layz1 layz2
declare -ai elt elx1 elx2 ely1 ely2 elz1 elz2 
declare -i i j

while read -r x1 x2 y1 y2 z1 z2; do
    xread="$xread$((x1))$nl$((x2+1))$nl"
    yread="$yread$((y1))$nl$((y2+1))$nl"
    zread="$zread$((z1))$nl$((z2+1))$nl"
done < <(tr -cs '\n[0-9]-' ' ' <"$in") # tr only keep numbers in input

readarray Xx < <(echo "${xread:0:-1}" |sort -n |uniq)
readarray Yy < <(echo "${yread:0:-1}" |sort -n |uniq)
readarray Zz < <(echo "${zread:0:-1}" |sort -n |uniq)

declare -i Xlen=$((${#Xx[@]}-1)) # last value is excluded [x1,x2[
declare -i Ylen=$((${#Yy[@]}-1))
declare -i Zlen=$((${#Zz[@]}-1))

for((i=0; i<Xlen+1; i++)); do xX["${Xx[i]}"]="$i"; done
for((i=0; i<Ylen+1; i++)); do yY["${Yy[i]}"]="$i"; done
for((i=0; i<Zlen+1; i++)); do zZ["${Zz[i]}"]="$i"; done

for((X=0; X<Xlen; X++)); do Xsize[X]=$((Xx[X+1] - Xx[X])); done
for((Y=0; Y<Ylen; Y++)); do Ysize[Y]=$((Yy[Y+1] - Yy[Y])); done
for((Z=0; Z<Zlen; Z++)); do Zsize[Z]=$((Zz[Z+1] - Zz[Z])); done

echo "== XYZ coords: 0..100 to 0..$Xlen 0..$Ylen 0..$Zlen, $((Xlen * Ylen * Zlen)) \"points\""

# we read the boot orders in reverse (|tac), as "paint' layers
# so that finding a XYZ point can stop the search
# shellcheck disable=SC2206,SC2020 # no need to quote integer values
while read -r type x1 x2 y1 y2 z1 z2; do
    [[ $type = on ]] && layt+=(1) || layt+=(0)
    layx1+=($x1)
    layx2+=($((x2+1)))
    layy1+=($y1)
    layy2+=($((y2+1)))
    layz1+=($z1)
    layz2+=($((z2+1)))
done < <(tr -cs '\n[onf0-9]-' ' ' <"$in" |tac)
declare -i laylen=${#layt[@]}
echo "== $laylen layers"

# find "loners" layers that do not intersect with any other
declare -ai intersections
for((i=0; i<laylen-1; i++)); do
    ((intersections[i])) && continue
    for((j=i+1; j<laylen; j++)); do
        ((intersections[j])) && continue
        ((layx1[j] >= layx2[i] || layx2[j] <= layx1[i] || layy1[j] >= layy2[i] || layy2[j] <= layy1[i] || layz1[j] >= layz2[i] || layz2[j] <= layz1[i])) && continue
        intersections[i]+=1
        intersections[j]+=1
    done
done
echo "== $((laylen - ${#intersections[@]})) loners and ${#intersections[@]} entangled layers"

# compute the "on" points of all the "on" loner layers, and build a stack
# of the remaining layers el (Entangled Layers)

on=0
for((i=0; i<laylen; i++)); do
    # shellcheck disable=SC2206,SC2020 # no need to quote integer values
    if ((intersections[i])); then
        #echo "== el[${#elt[@]}] = \"${layt[i]} ${layx1[i]} ${layx2[i]} ${layy1[i]} ${layy2[i]} ${layz1[i]} ${layz2[i]}"
        elt+=(${layt[i]})
        elx1+=(${layx1[i]})
        elx2+=(${layx2[i]})
        ely1+=(${layy1[i]})
        ely2+=(${layy2[i]})
        elz1+=(${layz1[i]})
        elz2+=(${layz2[i]})
    elif ((layt[i] == 1)); then
        ((on += (layx2[i]-layx1[i]) * (layy2[i]-layy1[i]) * (layz2[i]-layz1[i])))
    fi
done
echo "== $on on cubes from the loners"

# determine the sub-space enclosing only the non-loners, "el" layers
ellen=${#elt[@]}
declare -i MAXINT xmin xmax ymin ymax zmin zmax
MAXINT=$((2**63 - 1))
xmin=$MAXINT; xmax=0
ymin=$MAXINT; ymax=0
zmin=$MAXINT; zmax=0
    
for((i=0; i<ellen; i++)); do
    ((elx1[i]<xmin)) && ((xmin=elx1[i]))
    ((elx2[i]>xmax)) && ((xmax=elx2[i]))
    ((ely1[i]<ymin)) && ((ymin=ely1[i]))
    ((ely2[i]>ymax)) && ((ymax=ely2[i]))
    ((elz1[i]<zmin)) && ((zmin=elz1[i]))
    ((elz2[i]>zmax)) && ((zmax=elz2[i]))
done
echo "== Zone processed: x=$xmin..$xmax, y=$ymin..$ymax, z=$zmin..$zmax"
((Xmin=xX["$xmin"]))
((Xmax=xX["$xmax"]))
((Ymin=yY["$ymin"]))
((Ymax=yY["$ymax"]))
((Zmin=zZ["$zmin"]))
((Zmax=zZ["$zmax"]))

# now for all the "points" in the X Y Z space, look at the layers to see
# if they are on or off

for((X=Xmin; X<Xmax; X++)); do
    for((Y=Ymin; Y<Ymax; Y++)); do
        for((Z=Zmin; Z<Zmax; Z++)); do
            for((l=0; l<ellen; l++)); do
                if ((Xx[X] >= elx1[l])) && ((Xx[X] < elx2[l])) &&
                       ((Yy[Y] >= ely1[l])) && ((Yy[Y] < ely2[l])) &&
                       ((Zz[Z] >= elz1[l])) && ((Zz[Z] < elz2[l])); then
                    if ((elt[l])); then # sum of "points" sizes in xyz
                        ((on += Xsize[X] * Ysize[Y] * Zsize[Z]))
                        # off layers just stop the probing
                    fi
                    break
                fi
            done
        done
    done
done

echo "$on"
