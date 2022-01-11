#!/bin/bash
# https://adventofcode.com/days/day/23 puzzle #2
# See README.md in the parent directory
in="${1:-${0%-[0-9].*}.input}"
err(){ echo "***ERROR: $*" >&2; exit 1;}
#export tmp=tmp.$$; clean(){ rm -f "$tmp" "$tmp".*;}; trap clean 0

#TEST: example 12521
#TEST: input 15111

# with a 2nd argument, verbose mode
[[ -n $2 ]] && verbose=true || verbose=false

# Note: I left has a lot of debugging/tracing code commented out, labelled DDD

# with signal 1, ouput the current number of states  processed and min score
lapinfo(){
    echo "At $(date +%T), $states states processed, min $frontier1, current=$current"
    show-state "$current"
}
trap lapinfo 1

# the map coordinates are integers.
# We do not represent spaces where pods cannot stop (just outside a room)
# .. . . . ..  0 1   2   3   4   5 6
#   B C B D        7   8   9  10
#   D C B A       11  12  13  14
#   D C B A       15  16  17  18 
#   A D C A       19  20  21  22

shopt -s extglob                # we will extended shell globbing

# the following code supposes a configuration of 4 rooms, whith any depth
hallmax=6                       # maximum pos in hallway, built-in.
declare -a podname=(A B C D)  # 4 rooms, this is fixed.
declare -A podindex=([A]=0 [B]=1 [C]=2 [D]=3)

############ read input, this determines room_depth
# we only look at pods letters that can be given in compact form
# e.g: BCBDADCA for the exemple (#, . space, newlines ignored)
# DCBADCBA is inserted after the 1rst room row when reading a file,
# as we solve the 2nd exercise by default

if [[ $1 =~ ^[[:upper:]]+$ ]] && ((${#1} % 4 == 0)); then
    pods="$1"                  # start state given literally e.g: BCBDADCA 
else
    # read the start state from input file
    [[ -e $in ]] || err "Input file \"$in\" not found"
    pods=$(tr -cd '[:upper:]' <"$in")
    # insert the two lines DCBA DBAC
    pods="${pods:0:4}DCBADBAC${pods:4}"
fi
((room_depth=${#pods}/4))
echo "Input: $pods (room depth = $room_depth)"
for pod in A B C D; do
    p=${pods//[^$pod]/}
    ((${#p}==room_depth)) || err "${#p} $pod instead of $room_depth"
done
goalpos=
for((r=0; r<room_depth; r++)); do
    for p in "${podname[@]}"; do
        goalpods+="$p"
    done
done
echo "Goal:  $goalpods"

# from just a list of pods in rooms, get a sorted state
pods2state(){
    local -n _var="$1"
    local pods="$2"
    local pod gridelts
    local grid=()           # sparse array, to sort results
    local -i gridsize i r pos
    ((gridsize=hallmax+4*room_depth)) # a gridsize block for each pod class
    for((i=0; i<${#pods}; i++)); do
        pod="${pods:$i:1}"
        r="${podindex[$pod]}"
        pos=$((hallmax+1+i))
        grid[r*gridsize+pos]="$pod$pos"
    done
    gridelts="${grid[*]}"
    _var="${gridelts// /,}"
}

# create the initial sorted state
pods2state start "$pods"
echo "Starting state: $start"
# And the goal state to reach: e.g: A7,A11,B8,B12,C9,C13,D10,D14
pods2state goal "$goalpods"

############ Map-specific code START
# states are sorted lists of {pod}{pos} comma-separated, eg: A11,B8,D10
# finding neighbors and heuristic distance to the goal are how we encode
# the specific of the map
# hallway positions cannot end in hallway
# room cannot end in same room
declare -A energy
for pod in {0..3}; do energy["${podname[pod]}"]=$((10**pod)); done

# hallpaths generates reachable next places for places in the hallway
# arguments: paths to 4 rooms in order
# first arg is closeness to doors: 0 for pos #0 & #6, 1 otherwise
hallpaths(){
    local room0="$hallmax" close="$1"; shift    
    local doorpath nexts i
    for doorpath in "$@"; do
        score=$((${#doorpath}+close))
        [[ -n $doorpath ]] && ((score++))
        path="$doorpath"
        ((room0++))
        for((i=0; i<room_depth; i++)); do
            ((score++))
            path="${path}${path:+,}$((room0+i*4))"
            nexts="$nexts${nexts:+ }$score:$path"
        done
    done
    echo "$nexts"
}

# roompaths generates and sets reachable next places for all places in rooms
roompaths(){
    local room0="$1"
    local -a pathroom=("$2" "$3" "$4" "$5") # path to entrance of room
    local leftpath="$6" rightpath="$7"      # . is skippable door
    local nexts subnexts pr path
    local -i score i j r
    # targets in hallway, left and right
    while [[ $leftpath =~ [[:digit:]] ]]; do
        score=$(( (${#leftpath} + 1 ) / 2))
        nexts="$nexts${nexts:+ }$score:${leftpath//.,/}"
        leftpath="${leftpath%,*}"
        leftpath="${leftpath%,.}"
    done
    while [[ $rightpath =~ [[:digit:]] ]]; do
        score=$(((${#rightpath}+1)/2))
        nexts="$nexts${nexts:+ }$score:${rightpath//.,/}"
        rightpath="${rightpath%,*}"
        rightpath="${rightpath%,.}"
    done
    # origins in this room
    r=-1
    # targets in other rooms
    for pr in "${pathroom[@]}"; do
        ((++r))
        [[ -z $pr ]] && continue
        ((score=${#pr}+2))
        path="$pr"
        # descend into next places in target room
        for((j=0; j<room_depth;j++)); do
            path="$path,$((hallmax+1+r+j*4))"
            ((score++))
            nexts="$nexts${nexts:+ }$score:$path"
        done
    done
    # nexts is done for the first place in the room. derive the others
    mapnexts[$room0]="$nexts"
    prepath="$room0"
    r="$room0"
    for((i=1; i<room_depth; i++)); do
        subnexts=
        for sp in $nexts; do
            score=$(("${sp%:*}" + i))
            path="$prepath,${sp#*:}"
            subnexts="$subnexts${subnexts:+ }$score:$path"
        done
        prepath="$r,$prepath"
        ((r+=4))
        mapnexts[$r]="$subnexts"
    done
}
    
# mapnexts: for each pos, space-separated list of unweighted-cost:path
# e.g: for place 3: '5:2,7,11 4:2,7 3:8,12 2:8 3:9,13 2:9 5:4,10,14 4:5,10'
declare -a mapnexts
# we fill it for places in the hallway:
mapnexts+=("$(hallpaths 0 1 1,2 1,2,3 1,2,3,4)") # 0
mapnexts+=("$(hallpaths 1 '' 2 2,3 2,3,4)")      # 1
mapnexts+=("$(hallpaths 1 '' '' 3 3,4)")         # 2
mapnexts+=("$(hallpaths 1 2 '' '' 4)")           # 3
mapnexts+=("$(hallpaths 1 3,2 3 '' '')")         # 4
mapnexts+=("$(hallpaths 1 4,3,2 4,3 4 '')")      # 5
mapnexts+=("$(hallpaths 0 5,4,3,2 5,4,3 5,4 5)") # 6
# and the places in rooms:
roompaths  7 '' 2 2,3 2,3,4             .,1,0 .,2,.,3,.,4,.,5,6
roompaths  8 2 '' 3 3,4             .,2,.,1,0     .,3,.,4,.,5,6
roompaths  9 3,2 3 '' 4         .,3,.,2,.,1,0         .,4,.,5,6
roompaths 10 4,3,2 4,3 4 '' .,4,.,3,.,2,.,1,0             .,5,6

"$verbose" && {
    echo 'mapnexts=('
    for((i=0; i<${#mapnexts[@]}; i++)); do echo "  '${mapnexts[i]}' #$i"; done
    echo ')'
}

# the only room pod positions allowed: [7]=A, [11]=A, [8]=B, ...
declare -a ownroom
for((d=0; d<room_depth; d++)); do
    for r in {0..3}; do
        ownroom[$((hallmax+1 + d*4 + r))]="${podname[r]}"
    done
done

# for each room pos, are they deeper ones?: [7]='11 15 19', [11]='15 19 ...
declare -a deeper_rooms
for r in {0..3}; do
    deeper=
    for((d=room_depth-1; d>0; d--)); do
        deeper="$((hallmax+1 + d*4 + r))${deeper:+ }$deeper"
        deeper_rooms[$((hallmax+1 + (d-1)*4 + r))]="$deeper"
    done
done

# heuristic function, distance of a state to goal. Return value in $heuristic
# pre-computed minimal distances to goal positions for each pod position
# weighted by pod class energy: [A0]=6 [D4]=5000 ...
declare -A heuristic_dists
for pod in {0..3}; do
    final=$((hallmax+1 + (room_depth-1)*4 + pod))
    # we use the pre-computed mapnexts paths to final
    # shellcheck disable=SC2068 # yes we space-split into params
    for p in ${!mapnexts[@]}; do
        if [[ ${mapnexts[p]} =~ (^|\ )([[:digit:]]+):([[:digit:]]+,)*"$final"(\ |$) ]]; then
            heuristic_dists["${podname[pod]}$p"]=$((cost = BASH_REMATCH[2] * (10**pod)))
        fi
    done
    # except for the pod room itself where we just measure to the bottom
    for((d=0; d<room_depth; d++)); do
        heuristic_dists["${podname[pod]}$((hallmax+1 + d*4 + pod))"]=$(((room_depth - d - 1) * (10**pod)))
    done
done

"$verbose" && {
    echo 'heuristic_dists=('
    for p in {0..3}; do
        for((i=0; i<((hallmax+1+4*room_depth)); i++)); do
            echo -n " [${podname[p]}$i]=${heuristic_dists[${podname[p]}$i]}"
        done
        echo
    done
    echo ')'
}

############ semi-generic code: can work with any board sizes

# state-next puts a list of "cost:neighbors-state" in global costnexts[]
declare -a map
declare -a costnexts
state-nexts(){
    local state="$1" podpos pod filled
    local -i i pos lower
    costnexts=()
    map=()         # occupied spaces by their letters
    # shellcheck disable=SC2086 # yes we space-split into params
    set ${state//,/ }
    for podpos in "$@"; do
        map["${podpos:1}"]="${podpos:0:1}" 
    done
    for podpos in "$@"; do
        pod="${podpos:0:1}"
        pos="${podpos:1}"
        # do not move once we are in our room and only same pods are deeper
        if [[ ${ownroom[$pos]} == "$pod" ]]; then
            filled=true
            for lower in ${deeper_rooms["$pos"]}; do
                [[ ${map[lower]} == "$pod" ]] || { filled=false; break;}
            done
            "$filled" &&
                #echo "      $podpos: already final position, do not move" && # DDD
                continue
        fi
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
        ((endpos>hallmax)) && [[ ${ownroom[$endpos]} != "$pod" ]] && continue
        free=true
        # is the path free of other pods?
        for pos in ${path//,/ }; do
            [[ -n ${map[pos]} ]] && {
                #echo "      $podpos: path $path blocked at $pos" # DDD
                free=false; break;} # path blocked
        done
        "$free" || continue
        # we cannot go into our room if an "alien" in it
        for pos in ${deeper_rooms["$endpos"]}; do
            [[ -n ${map[pos]} ]] && [[ ${map[pos]} != "$pod" ]] && {
                #echo "      $pod path: $path has alien ${BASH_REMATCH[1]}$pos" # DDD
                free=false; break;} # an alien there, abort!
        done
        "$free" || continue
        # in our room, we must go as deep as possible
        for pos in ${deeper_rooms["$endpos"]}; do
            [[ -z ${map[pos]} ]] && { 
                #echo "      endpos: $pod$endpos, has deeper at $pos" # DDD
                free=false; break;} # we can go deeper than that!
        done
        "$free" || continue
        # then create "next" by replacing the podpos by the moved one at endpos
        # ensure states are always sorted, so that representation is unique
        # insertion sort, not optimal
        next=
        # shellcheck disable=SC2086 # yes we space-split into params
        set ${state//,/ }
        inserted=false
        for pp in "$@"; do
            [[ $pp == "$podpos" ]] && continue # remove origin from next
            # [[ $pp == "$pod$endpos" ]] && err "$podpos->$pod$endpos already in $state, map[$endpos]=${map[$endpos]}" # DDD
            ! "$inserted" &&    # insert just before the first bigger one
                [[ $pod < "${pp:0:1}" || ($pod == "${pp:0:1}" && $endpos -lt ${pp:1}) ]] &&
                next+=",$pod$endpos" &&
                inserted=true
            next+=",$pp"
        done
        "$inserted" || next+=",$pod$endpos" # target greater than all in state
        # statenext=${next/,/}; s=${statenext//[^,]/}; ((${#s} == 4*room_depth-1)) || err "next state $statenext has not $((4*room_depth)) elements, from $state" # DDD
        costnexts+=("$cost${next/,/:}")
        #echo "      next: $podpos -> $pod$pos ==> $cost${next/,/:}" # DDD
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
    #((frontier1<0 || cost<frontier1)) && echo "      frontier-put[$cost]=$state (new low)" || echo "      frontier-put[$cost]=$state" # DDD
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

# Debug: display the move between 2 states
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

# extended state display
show-state(){
    local state="$1" i pod r
    for i in 0 1 . 2 . 3 . 4 . 5 6; do 
        if [[ $i == '.' ]]; then
            echo -n '.'
        else
            pod-place pod "$state" "$i"
            echo -n "$pod"
        fi
    done
    echo
    for((r=0; r<room_depth; r++)); do
        echo -n ' '
        for((i=0; i<4; i++)); do
            pod-place pod "$state" "$((hallmax+1+i+r*4))"
            echo -n " $pod"
        done
        echo
    done
}
pod-place(){
    local -n _pod="$1"
    local state="$2" p="$3"
    if [[ $state =~ ([A-D])"$p"(,|$) ]]; then
        _pod="${BASH_REMATCH[1]}"
    else
        _pod='.'
    fi
}            

############ The main generic A* code, lifted from:
# https://www.redblobgames.com/pathfinding/a-star/introduction.html

echo "Goal state:     $goal"    # marks the end of initialisation code
echo "Computing... for lap info: kill -1 $$"

declare -A came_from
declare -A cost_sofar
declare -A cost_from            # for tracing

came_from["$start"]=
cost_sofar["$start"]=0
cost_from["$start"]=0
frontier-init "$start"
states=0

export TIMEFORMAT="Computations done in %3lR"
time while true; do
    frontier-get current || break
    # s=${current//[^,]/}; ((${#s} == 4*room_depth-1)) || err "state $current has not $((4*room_depth)) elements" # DDDD
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
    cost="${cost_from[$state]}"
    "$verbose" && ((cost)) && show-state "$state" |tac
    ((cost)) || cost="cost of step" # legend for first row
    echo "  $state${tab}$(show-move "$state")  ${tab}($cost)"
    state="${came_from[$state]}"
done |tac

echo "Minimal path cost:"
echo "${cost_sofar[$current]}"

