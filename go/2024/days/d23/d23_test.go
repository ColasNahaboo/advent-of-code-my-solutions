package main

import (
	"testing"
	"reflect"
)

// https://en.wikipedia.org/wiki/Bron%E2%80%93Kerbosch_algorithm#Example

func EdgesInit() BKSet {
	bks := BKNew()
	bks[0] = uint64(63)
	bkedges = [][]int{[]int{1,4}, []int{0,2,4}, []int{1,3}, []int{2,4,5}, []int{0,1,3}, []int{3}}
	return bks
}

// set of unit tests
func Test(t *testing.T) {
	// setup code
	s := EdgesInit()
	e := BKNew()

	e[0] = uint64(63)
	t_add("AddNode 1", t, s, 1, e)
	e = BKNew()
	e[0] = uint64(63 + 1<<32)
	t_add("AddNode 32", t, s, 32, e)
	e = BKNew()
	e[0] = uint64(63)
	e[1] = uint64(1 << (80-64))
	t_add("AddNode 80", t, s, 80, e)

	t_res("Restrict 1", t, s, BKMake(0, 3, 5), BKMake(1, 2, 4))
}

// test function
func t_add(label string, t *testing.T, s BKSet, i int, expected BKSet) {
	t.Run(label, func(t *testing.T) {
		got := s.AddNode(i)
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}})}

// test function
func t_res(label string, t *testing.T, s, ss BKSet, expected BKSet) {
	t.Run(label, func(t *testing.T) {
		got := s.Restrict(ss)
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("expected '%v' but got '%v'", expected, got)
		}})}
