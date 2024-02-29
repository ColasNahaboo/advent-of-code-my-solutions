// Adventofcode 2017, d18, in go. https://adventofcode.com/2017/day/18

// Part3: alternative implementation of part2 with go goroutines
// We use 3 goroutines: the two tablet and an observer of their status
// On entering waiting for rcv on the channel, it sends a "clear status" to obs
// On reading it sends a "set status" to obs
// Thus obs can tell that the two tablets are in deadlock when both status are 0

// +----+                      +----+
// | T0 |-----channel-0-1----->| T1 |
// |    |<----channel-1-0----- |    |
// +----+                      +----+
//    |                          |     +----------+
//    +-------observer-channel---+---> | observer |
//                                     +----------+

package main

//////////// Part 3 = part2 using Goroutines and channels

const CHANBUFSIZE = 100			// raise if execution fails in deadlock

func part3(lines []string) int {
	t1 := parse3(0, lines)
	t2 := parse3(1, lines)
	cobs := make(chan int, CHANBUFSIZE)		// t1, t2 -> obs
	t1.cobs = cobs
	t2.cobs = cobs
	c12 := make(chan int, CHANBUFSIZE)		// t1 -> t2
	t1.cout = c12
	t2.cin = c12
	c21 := make(chan int, CHANBUFSIZE)		// t2 -> t1
	t1.cin = c21
	t2.cout = c21
	go t1.Run()
	go t2.Run()
    observe(cobs)				// wait for all tablets to complete
	return t2.nsends
}

func parse3(id int, lines []string) (t *Tablet) {
	t = parse(lines)
	t.id = id
	p := regOf(parseParam("p", t))
	t.regs[p] = t.id
	t.opexec[parseOp("snd", t)] = snd3Exec
	t.opexec[parseOp("rcv", t)] = rcv3Exec
	return
}

var statusbit = []int{1, 2}  		// 2^t.id

func observe(cobs chan int) {
	status := statusbit[0] + statusbit[1] // both are active
	for {
		VPf("Status: %d\n", status)
		if status == 0 && len(cobs) == 0 {
			return
		}
		status += <-cobs
	}
}

// part3 codes

func snd3Exec(t *Tablet) bool {
	value := valueOf(t.prog[t.p].x, t)
	t.cout <- value
	t.nsends++
	VPf("%sTablet [%d] SEND(#%d) %d\n", tindent[t.id], t.id, t.nsends, value)
	t.p++
	return true
}

func rcv3Exec(t *Tablet) bool {
	i := t.prog[t.p]
	if len(t.cin) > 0 {			// there is something to read, just do it
		t.regs[regOf(i.x)] = <-t.cin
		VPf("%sTablet [%d] READ %d -> %s\n", tindent[t.id], t.id, t.regs[regOf(i.x)], string(byte(regOf(i.x)) + 'a'))
		t.p++
		return true
	}
	t.cobs <- - statusbit[t.id]	// tell  observer we are entering wait state
	t.regs[regOf(i.x)] = <-t.cin
	VPf("%sTablet [%d] READ %d -> %s\n", tindent[t.id], t.id, t.regs[regOf(i.x)], string(byte(regOf(i.x)) + 'a'))
	t.cobs <- statusbit[t.id]
	t.p++
	return true
}

//////////// PrettyPrinting & Debugging functions
