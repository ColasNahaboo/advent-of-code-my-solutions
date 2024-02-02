// A simple generic linked list usable as LIFO/FIFO or stacks/pipes
// with randomly adressable elements, a chain
// as we added a direct pointer to the last element
// Create a LinkedList just by normal initialization:
// newlist := LinkedList[T]{}

// There is no explicit iterator function, just use for a linked list ll:
// 		for iterator := ll.head; iterator != nil; iterator = iterator.next {
//          element := iterator.val

package main

type LinkedList[T comparable] struct {
	head *LinkedCell[T]			// classic entry point, nil on empty lists
	tail *LinkedCell[T]			// optimization for FIFO use
}

type LinkedCell[T comparable] struct {
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
	if ll.IsEmpty() {
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
	if ll.IsEmpty() {
		ll.head = &cell
		ll.tail = &cell
	} else {
		ll.tail.next = &cell
	}
	ll.tail = &cell
}

//////////// Convenience common method

// simple test
func (ll *LinkedList[T])IsEmpty() bool {
	return ll.head == nil
}

// convert to list
func (ll *LinkedList[T])ToList() (l []T) {
	for c := ll.head; c != nil; c = c.next {
		l = append(l, c.val)
	}
	return 
}

// create from list, same order
func MakeLinkedList[T comparable](l []T) (ll *LinkedList[T]) {
	for _, v := range l {
		ll.Put(v)
	}
	return 
}

// create from list, reverse order. A tad faster.
func MakeLinkedListRev[T comparable](l []T) (ll *LinkedList[T]) {
	for _, v := range l {
		ll.Push(v)
	}
	return 
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

// ONLY_FOR_COMPARABLE_START
//////////// Chained lists, where you can find, insert and delete links
// Note that these are the only functions requiring a [T comparable]
// all the rest could work with a [T any]
// To obtain a packkage working with any type:
//   - replace all [T comparable] by [T any]
//   - remove everything below

// is a value already in the queue?
func (ll *LinkedList[T])Has(v T) bool {
	for c := ll.head; c != nil; c = c.next {
		if c.val == v {
			return true
		}
	}
	return false
}

// for a value, returns its cell and the previous one, to cut&paste
// 3rd returned value indicate that the value was actually found
// if v found at head, prev will be nil
func (ll *LinkedList[T])Index(v T) (prev, cell *LinkedCell[T], ok bool) {
	for c := ll.head; c != nil; c = c.next {
		if c.val == v {
			cell = c
			ok = true
			return 
		}
		prev = c
	}
	return
}

// delete value. Return true if value was present, false otherwise
func (ll *LinkedList[T])Delete(v T) bool {
	var prev *LinkedCell[T]	
	for c := ll.head; c != nil; c = c.next {
		if c.val == v {
			if c == ll.head {
				ll.head = c.next
				if ll.head == nil {
					ll.tail = nil
				}
			} else if c == ll.tail {
				ll.tail = prev
				prev.next = nil
			} else {
				prev.next = c.next
			}
			return true
		}
		prev = c
	}
	return false
}

// insert value val after value v.
func (ll *LinkedList[T])InsertAfter(v, val T) bool {
	for c := ll.head; c != nil; c = c.next {
		if c.val == v {
			cell := LinkedCell[T]{val: val, next: c.next}
			c.next = &cell
			if c == ll.tail {
				ll.tail = &cell
			}
			return true
		}
	}
	return false
}

// insert value val before value v.
func (ll *LinkedList[T])InsertBefore(v, val T) bool {
	var prev *LinkedCell[T]	
	for c := ll.head; c != nil; c = c.next {
		if c.val == v {
			cell := LinkedCell[T]{val: val, next: c}
			prev.next = &cell
			if c == ll.head {
				ll.head = &cell
			}
			return true
		}
		prev = c
	}
	return false
}

// ONLY_FOR_COMPARABLE_END
