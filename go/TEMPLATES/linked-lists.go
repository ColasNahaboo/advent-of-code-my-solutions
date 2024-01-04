// A simple generic linked list usable as LIFO/FIFO or stacks/pipes
// as we added a direct pointer to the last element
// Create a LinkedList just by normal initialization:
// newlist := LinkedList[T]{}

package main

type LinkedList[T any] struct {
	head *LinkedCell[T]			// classic entry point, nil on empty lists
	tail *LinkedCell[T]			// optimization for FIFO use
}

type LinkedCell[T any] struct {
	val T						// the stored value
	next *LinkedCell[T]			// next cell or nil
}

//////////// LIFO (aka stack): push / pop

// push a value in front of the LL
func (ll *LinkedList[T])Push (val T) {
	cell := LinkedCell[T]{val, ll.head}
	ll.head = &cell
	if ll.tail == nil {
		ll.tail = &cell
	}
}

// pop a value from the front of the LL
func (ll *LinkedList[T])Pop() (val T) {
	if ll.isEmpty() {
		panic("Linkedlist Pop on empty list")
	}
	val = ll.head.val
	ll.head = ll.head.next
	return
}

//////////// FIFO (aka pipe): put / pop

// put a value at the end of the  LL
func (ll *LinkedList[T])Put (val T) {
	cell := LinkedCell[T]{val: val}
	if ll.isEmpty() {
		ll.head = &cell
		ll.tail = &cell
	} else {
		ll.tail.next = &cell
	}
	ll.tail = &cell
}

//////////// Convenience common method

// simple test
func (ll *LinkedList[T])isEmpty() bool {
	return ll.head == nil
}

//////////// Convenience (but expensive for single-linked lists) methods

// list size
func (ll *LinkedList[T])Len() (size int) {
	if ll.head == nil {
		return 0
	}
	for c := ll.head; c != nil; c = c.next {
		size++
	}
	return size
}

// Find previous cell
func (ll *LinkedList[T])Prev(cell *LinkedCell[T]) *LinkedCell[T] {
	c := ll.head
	for {
		if c.next == c {
			return c
		} else  if c == ll.tail {
			return nil
		}
		c = c.next 
	}
}

// Get a value from the end of the LL and remove it: a "Pop" from the end
func (ll *LinkedList[T])Tail() (val T) {
	val = ll.tail.val
	beforeTail := ll.Prev(ll.tail)
	ll.tail = beforeTail
	beforeTail.next = nil
	return
}
