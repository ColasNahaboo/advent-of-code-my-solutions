#!/bin/bash

i=4
r=1

while true; do
    while ((r*2 < i)); do
        ((i++))
        ((r++))
        echo "$i [$r]"
    done
    while ((r < i)); do
        ((i++))
        ((r+=2))
        echo "$i [$r]"
    done
    ((i++))
    r=1
    echo "$i [$r]"
done
