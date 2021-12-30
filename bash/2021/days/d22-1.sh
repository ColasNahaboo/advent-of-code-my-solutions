#!/bin/bash
# https://adventofcode.com/days/day/22 puzzle #1
# See README.md in the parent directory
in="${1:-${0%-[0-9].*}.input}"; [[ -e $in ]] || exit 1
err(){ echo "***ERROR: $*" >&2; exit 1;}
#export tmp=tmp.$$; clean(){ rm -f "$tmp" "$tmp".*;}; trap clean 0

#TEST: small 39
#TEST: example 590784
#TEST: input 611176

# we will work on another referential: positive integer coordinates X Y Z
# were integers represent all the frontiers between boot cubes
# x y z are input coords plus 50 so they are >0 and usable as array indexes
# also the upper bounds are exclusive: in the input x=x1..x2 means all x1 to x2,
# but internally our interval x1 to x2 means x1 to x2-1
# in other words, we use intervals [x1,x2[ where the input is in [x1,x2]

# debugging functions
# print a layer
p_lay(){
    local l="$1"
    local t
    ((layt[l])) && t=ON || t=Off
    echo "${2}Layer [$l] $t: x=$((${Xx[${layx1[l]}]}-50))..$((${Xx[${layx2[l]}]}-50)),y=$((${Yy[${layy1[l]}]}-50))..$((${Yy[${layy2[l]}]}-50)),z=$((${Zz[${layz1[l]}]}-50))..$((${Zz[${layz2[l]}]}-50))"
}
# print a point
p_point(){
    echo "${4}Point X=$1,Y=$2,Z=$3 ==> x=$((${Xx[$1]}-50)),y=$((${Yy[$2]}-50)),z=$((${Zz[$3]}-50))"
}

# debugging end

nl=$'\n'
declare -ai Xx Yy Zz xX yY zZ   # conversion tables
declare -ai Xsize Ysize Zsize   # the size of a XYZ point in xyz cubes
# layers: type (on=1, off=0), X Y Z
declare -ai layt layx1 layx2 layy1 layy2 layz1 layz2 
declare -i i

while read -r x1 x2 y1 y2 z1 z2; do
    if ((x1 >= -50)) && ((x2 <= 50)) &&
           ((y1 >= -50)) && ((y2 <= 50)) &&
           ((z1 >= -50)) && ((z2 <= 50)); then
        xread="$xread$((x1+50))$nl$((x2+51))$nl"
        yread="$yread$((y1+50))$nl$((y2+51))$nl"
        zread="$zread$((z1+50))$nl$((z2+51))$nl"
    fi
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
while read -r type x1 x2 y1 y2 z1 z2; do
    if ((x1 >= -50)) && ((x2 <= 50)) &&
           ((y1 >= -50)) && ((y2 <= 50)) &&
           ((z1 >= -50)) && ((z2 <= 50)); then
        [[ $type = on ]] && layt+=(1) || layt+=(0)
        layx1+=(${xX[$((x1+50))]})
        layx2+=(${xX[$((x2+51))]})
        layy1+=(${yY[$((y1+50))]})
        layy2+=(${yY[$((y2+51))]})
        layz1+=(${zZ[$((z1+50))]})
        layz2+=(${zZ[$((z2+51))]})
    fi
done < <(tr -cs '\n[onf0-9]-' ' ' <"$in" |tac)
declare -i laylen=${#layt[@]}
echo "== $laylen layers"

# now for all the "points" in the X Y Z space, look at the layers to see
# if they are on or off

on=0
for((X=0; X<Xlen; X++)); do
    for((Y=0; Y<Ylen; Y++)); do
        for((Z=0; Z<Zlen; Z++)); do
            for((l=0; l<laylen; l++)); do
                if ((X >= layx1[l])) && ((X < layx2[l])) &&
                       ((Y >= layy1[l])) && ((Y < layy2[l])) &&
                       ((Z >= layz1[l])) && ((Z < layz2[l])); then
                    if ((layt[l])); then # sum of "points" sizes in xyz
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
