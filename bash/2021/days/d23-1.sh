#!/bin/bash
# https://adventofcode.com/days/day/23 puzzle #1
# See README.md in the parent directory
in="${1:-${0%-[0-9].*}.input}"
err(){ echo "***ERROR: $*" >&2; exit 1;}
#export tmp=tmp.$$; clean(){ rm -f "$tmp" "$tmp".*;}; trap clean 0

#TEST: example 12521
#TEST: input 15111

# Starting state can be given as argument instead of a input file, e.g:
# d23-1.sh A1,A5,B8,B12,C9,C13,D10,D14

# the map coordinates are integers.
# We do not represent spaces where pods cannot stop (just outside a room)
# .. . . . ..  0 1   2   3   4   5 6
#   B C B D        7   8   9  10
#   A D C A       11  12  13  14

shopt -s extglob                # we will extended shell globbing

############ Map-specific code START
# states are sorted lists of {pod}{pos} comma-separated
goal='A7,A11,B8,B12,C9,C13,D10,D14'
hallmax=6                       # maximum pos in hallway

# finding neighbors and heuristic distance to the goal are how we encode
# the specific of the map
# map next: for each pos, space-separated list of unweighted-cost:path
# hallway positions cannot end in hallway
# room cannot end in same room
mapnexts=(                      # position
    '4:1,7,11 3:1,7 6:1,2,8,12 5:1,2,8 8:1,2,3,9,13 7:1,2,3,9 10:1,2,3,4,10,14 9:1,2,3,4,10' #0
    '3:7,11 2:7 5:2,8,12 4:2,8 7:2,3,9,13 6:2,3,9 9:2,3,4,10,14 8:2,3,4,10' #1
    '3:7,11 2:7 3:8,12 2:8 5:3,9,13 4:3,9 7:3,4,10,14 6:3,4,10' #2
    '5:2,7,11 4:2,7 3:8,12 2:8 3:9,13 2:9 5:4,10,14 4:4,10' #3
    '7:3,2,7,11 6:3,2,7 5:3,8,12 4:3,8 3:9,13 2:9 3:10,14 2:10' #4
    '9:4,3,2,7,11 8:4,3,2,7 7:4,3,8,12 6:4,3,8 5:4,9,13 4:4,9 3:10,14 2:10' #5
    '10:5,4,3,2,7,11 9:5,4,3,2,7 8:5,4,3,8,12 7:5,4,3,8 6:5,4,9,13 5:5,4,9 4:5,10,14 3:5,10' #6
    '3:1,0 2:1 9:2,3,4,5,6 8:2,3,4,5 6:2,3,4 4:2,3 2:2 9:2,3,4,10,14 8:2,3,4,10 7:2,3,9,13 6:2,3,9 5:2,8,12 4:2,8' #7
    '5:2,1,0 4:2,1 2:2 7:3,4,5,6 6:3,4,5 4:3,4 2:3 5:2,7,11 4:2,7 7:3,4,10,14 5:3,9,13 4:3,9 6:3,4,10' #8
    '7:3,2,1,0 6:3,2,1 4:3,2 2:3 5:4,5,6 4:4,5 2:4 7:3,2,7,11 6:3,2,7 5:3,8,12 4:3,8 5:4,10,14 4:4,10' # 9
    '9:4,3,2,1,0 8:4,3,2,1 6:4,3,2 4:4,3 2:4 2:5 3:5,6 9:4,3,2,7,11 8:4,3,2,7 7:4,3,8,12 6:4,3,8 5:4,9,13 4:4,9' #10
    '4:7,1,0 3:7,1 10:7,2,3,4,5,6 9:7,2,3,4,5 7:7,2,3,4 5:7,2,3 3:7,2 10:7,2,3,4,10,14 9:7,2,3,4,10 8:7,2,3,9,13 7:7,2,3,9 6:7,2,8,12 5:7,2,8' #11
    '6:8,2,1,0 5:8,2,1 3:8,2 8:8,3,4,5,6 7:8,3,4,5 5:8,3,4 3:8,3 6:8,2,7,11 5:8,2,7 6:8,3,9,13 5:8,3,9 8:8,3,4,10,14 7:8,3,4,10' #12
    '8:9,3,2,1,0 7:9,3,2,1 5:9,3,2 3:9,3 6:9,4,5,6 5:9,4,5 3:9,4 8:9,3,2,7,11 7:9,3,2,7 6:9,3,8,12 5:9,3,8 6:9,4,10,14 5:9,4,10' # 13
    '10:10,4,3,2,1,0 9:10,4,3,2,1 7:10,4,3,2 5:10,4,3 3:10,4 3:10,5 4:10,5,6 10:10,4,3,2,7,11 9:10,4,3,2,7 8:10,4,3,8,12 7:10,4,3,8 6:10,4,9,13 5:10,4,9' #14
)
# the only room pod positions allowed
declare -A ownroom=(
    [A7]=1 [A11]=1 [B8]=1 [B12]=1 [C9]=1 [C13]=1 [D10]=1 [D14]=1
)
# for each room pos, are they deeper ones?
declare -A deeper_rooms=( [A7]=11 [B8]=12 [C9]=13 [D10]=14 )
declare -A energy=( [A]=1 [B]=10 [C]=100 [D]=1000 )         # pod classes weights

# heuristic function, distance of a state to goal. Return value in $heuristic
# pre-computed minimal distances to goal positions for each pod position
# weighted by pod class energy
declare -A heuristic_dists=(
    [A0]=4 [A1]=3 [A2]=3 [A3]=5 [A4]=7 [A5]=9 [A6]=10 [A7]=1
    [A8]=5 [A9]=7 [A10]=9 [A11]=0 [A12]=6 [A13]=8 [A14]=10
    [B0]=60 [B1]=50 [B2]=30 [B3]=30 [B4]=50 [B5]=70 [B6]=80 [B7]=50
    [B8]=10 [B9]=50 [B10]=70 [B11]=60 [B12]=0 [B13]=60 [B14]=80
    [C0]=800 [C1]=700 [C2]=500 [C3]=300 [C4]=300 [C5]=500 [C6]=600 [C7]=700
    [C8]=500 [C9]=100 [C10]=500 [C11]=800 [C12]=600 [C13]=0 [C14]=600
    [D0]=10000 [D1]=9000 [D2]=7000 [D3]=5000 [D4]=3000 [D5]=3000 [D6]=4000 [D7]=9000
    [D8]=7000 [D9]=5000 [D10]=1000 [D11]=10000 [D12]=8000 [D13]=6000 [D14]=0
)
############ map-specific code END

############ Read input

if [[ $1 =~ ^[[:upper:]]+$ ]] && ((${#1} % 4 == 0)); then
    pods="$1"                  # start state given literally e.g: BCBDADCA 
else
    # read the start state from input file
    [[ -e $in ]] || err "Input file \"$in\" not found"
    pods=$(tr -cd '[:upper:]' <"$in")
fi
echo "Input: $pods (room depth fixed to 2)"
echo "Goal:  ABCDABCD"
start=$(for((i=0; i<${#pods}; i++)); do
            echo "${pods:$i:1}""$(printf %02d $((hallmax+i+1)))"
        done |sort |sed -r -e 's/([[:upper:]])0([[:digit:]])/\1\2/g' |tr '\n' ,)
start="${start%,}"
echo "Starting state: $start"

############ semi-generic code: can work with any board sizes
declare -A startpos                        # did we move yet?
for pp in ${start//,/ }; do startpos["$pp"]=1; done

# state-next puts a list of "cost:neighbors-state" in global costnexts[]
declare -a map
declare -a costnexts
state-nexts(){
    local state="$1" podpos
    local -i i
    costnexts=()
    map=()
    for i in ${state//[^0-9]/ }; do map[$i]=1; done # occupied spaces
    # shellcheck disable=SC2086 # yes we space-split into params
    set ${state//,/ }
    for podpos in "$@"; do
        ((ownroom[$podpos] && ! startpos[$podpos])) && continue # arrived
        # shellcheck disable=SC2086 # yes we space-split into params
        podpos-next "$podpos" ${mapnexts[${podpos:1}]}
    done
}

podpos-next(){
    local podpos="$1" pod="${1:0:1}"
    local costpath cost path pos free endpos weight next pp inserted
    shift
    weight=${energy["$pod"]}
    #echo "      podpos paths of $podpos: $*" # DDD
    for costpath in "$@"; do
        ((cost = ${costpath%:*} * weight))
        path="${costpath#*:}"
        endpos="${path##*,}"
        # we cannot end into a room other than our own
        ((endpos>hallmax && ! ownroom[$pod$endpos])) && continue
        # in our room, we must go as deep as possible
        for pos in ${deeper_rooms["$pod$endpos"]}; do
            ((map[pos])) || { 
                #echo "      endpos: $pod$endpos, has deeper at $pos" # DDD
                free=false; break;} # we can go deeper than that!
        done
        free=true               # is the path free of other pods?
        for pos in ${path//,/ }; do
            ((map[pos])) && {
                #echo "      $podpos: path $path blocked at $pos" # DDD
                free=false; break;} # path blocked
        done
        "$free" || continue
        # ensure states are always sorted, so that representation is unique
        # insertion sort, not optimal
        next=
        # shellcheck disable=SC2086 # yes we space-split into params
        set ${state//,/ }
        inserted=false
        for pp in "$@"; do
            [[ $pp == "$podpos" ]] && continue
            ! "$inserted" &&
                [[ $pod < "${pp:0:1}" || ($pod == "${pp:0:1}" &&$pos -lt ${pp:1}) ]] &&
                next+=",$pod$pos" &&
                inserted=true
            next+=",$pp"
        done
        "$inserted" || next+=",$pod$pos" # new is greater than all
        costnexts+=("$cost${next/,/:}")
        #echo "      next: $podpos -> $pod$pos ==> $cost${next/,/:}" # DDD
    done
}

state-sort(){
    for podpos in "$@"; do
        # shellcheck disable=SC2086 # yes we space-split into params
        ${podpos:1}+="$podpos"
    done
}

############ Generic code: can work on any style of map
# the priority queue frontier. Ops: put(state, cost) get()
declare -a frontier           # array of colon-separated states indexed by cost
declare -i frontier1          # index of 1rt (lowest) element. negative if empty
# sort-insert a state with cost
frontier-put(){
    local state="$1"
    local -i cost="$2"
    frontier[$cost]+="$state:"
    ((frontier1<0 || cost<frontier1)) && ((frontier1=cost))
    #((frontier1<0 || cost<frontier1)) && echo "      frontier-put new low: $cost: $state" || echo "      frontier-put: $cost: $state" # DDD
}
# remove the least cost state and return it, Return 1 if empty
frontier-get(){
    local -n _state="$1"        # returns the lowest cost state into it
    ((frontier1<0)) && { _state=; return 1;}
    local list="${frontier[frontier1]}"
    # shellcheck disable=SC2034
    _state="${list%%:*}"
    list="${list#*:}"
    #echo "      frontier-get[$frontier1]: $_state" # DDD
    if [[ -n $list ]]; then
        frontier[frontier1]="$list"
    else
        local -i i
        unset "frontier[frontier1]"
        if ((${#frontier[@]})); then
            # shellcheck disable=SC2068 # yes we space-split into params
            set ${!frontier[@]}
            frontier1="$1"
        else
            frontier1=-1            # queue is empty
        fi
    fi
    return 0
}
# initialize frontier with start state
frontier-init(){
    local state="$1"
    frontier[0]="$state:"
    frontier1=0
}

# heuristic is the approximate distance to the goal
declare -i heuristic
heuristic(){
    # shellcheck disable=SC2086 # yes we space-split into params
    set ${2//,/ }            # split the state into pod positions
    local podpos
    heuristic=0
    for podpos in "$@"; do ((heuristic+=heuristic_dists[$podpos])); done
}

# Debug: print the state of the frontier queue
frontier-print(){
    local -i i
    echo -n "==f1=$frontier1$1"
    # shellcheck disable=SC2068 # yes we space-split into params
    for i in ${!frontier[@]}; do
        [[ -n ${frontier[i]} ]] && echo -n " [$i]=${frontier[i]}"
    done    
    echo
}

# displays the move between 2 states
show-move(){
    local state="$1"
    local old="${came_from[$state]}"
    # shellcheck disable=SC2155,SC2086
    if [[ -n $old ]]; then
        local o=$({ tr , '\n' <<<$old; tr , '\n' <<<$old; tr , '\n' <<<$state;}|sort|uniq -c|grep -oP '2 \K.*')
        local n=$({ tr , '\n' <<<$old; tr , '\n' <<<$old; tr , '\n' <<<$state;}|sort|uniq -c|grep -oP '1 \K.*')
        echo "$o->$n"
    fi
}

############ The main generic A* code, lifted from:
# https://www.redblobgames.com/pathfinding/a-star/introduction.html

echo "Goal state:     $goal"    # marks the end of initialisation code
echo "Computing..."

declare -A came_from
declare -A cost_sofar
declare -A cost_from            # for tracing

came_from["$start"]=
cost_sofar["$start"]=0
cost_from["$start"]=0
frontier-init "$start"
states=0

export TIMEFORMAT="Minimal path found in %3lR"
time while true; do
    frontier-get current || break
    ((states++))
    # shellcheck disable=SC2154 # current is assigned by the above
    [[ $current == "$goal" ]] && break # found!
    #echo "== current: [$current]" # DDD
    state-nexts "$current"
    # shellcheck disable=SC2154 # costnexts[] is assigned by the above
    for costnext in "${costnexts[@]}"; do
        cost="${costnext%:*}"
        next="${costnext#*:}"
        ((new_cost = cost_sofar[$current] + cost))
        if [[ -z ${cost_sofar[$next]} ]] || ((new_cost<cost_sofar[$next])); then
            cost_sofar[$next]="$new_cost"
            heuristic "$goal" "$next"
            ((priority = new_cost + heuristic))
            frontier-put "$next" "$priority"
            came_from["$next"]="$current"
            cost_from["$next"]="$cost"
        fi
    done
done

echo "States processed: $states"
[[ -z $current ]] && err "No path found!"

# display path
tab=$'\t'
echo "Minimal path:"
state="$current"
while [[ -n $state ]]; do
    echo "  $state${tab}$(show-move "$state")${tab}(${cost_from[$state]})"
    state="${came_from[$state]}"
done |tac

echo "Minimal path cost:"
echo "${cost_sofar[$current]}"

