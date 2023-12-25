// Package combinations provides an iterator that goes through
// combinations in lexicographic order.
//
// This code is significantly indebted to
//
// * https://compprog.wordpress.com/2007/10/17/generating-combinations-1/
//   (Note: This blog post has an array out of bounds error bug!)
// * https://docs.python.org/3/library/itertools.html#itertools.combinations

// Code at https://github.com/carlmjohnson/go-utils/tree/master/combinations
// by Carlana Johnson https://carlana.net/

package main

// Iterator iterates through the length K combinations of
// an N sized set. E.g. for K = 2 and N = 3, it sets Comb to {0, 1},
// then {0, 2}, and finally {1, 2}.
type Iterator struct {
	N, K int
	Comb []int
}

// Init initializes Comb as []int{0, 1, ..., k-1}. Automatically called
// on first call to Next().
func (c *Iterator) Init() {
	c.Comb = make([]int, c.K)
	for i := 0; i < c.K; i++ {
		c.Comb[i] = i
	}
}

// Next sets Comb to the next combination in lexicographic order.
// Returns false if Comb is already the last combination in
// lexicographic order. Initializes Comb if it doesn't exist.
func (c *Iterator) Next() bool {
	var (
		i int
	)

	if len(c.Comb) != c.K {
		c.Init()
		return true
	}

	// Combination (n-k, n-k+1, ..., n) reached
	// No more combinations can be generated
	if c.Comb[0] == c.N-c.K {
		return false
	}

	for i = c.K - 1; i >= 0; i-- {
		c.Comb[i]++
		if c.Comb[i] < c.N-c.K+1+i {
			break
		}
	}

	// c.Comb now looks like (..., x, n, n, n, ..., n).
	// Turn it into (..., x, x + 1, x + 2, ...)
	for i = i + 1; i < c.K; i++ {
		c.Comb[i] = c.Comb[i-1] + 1
	}

	return true
}

// StringN iterates through sub-string combinations of length K.
//
// E.g. For "abc", 2; Comb is set to "ab", "ac", and "bc" succesively.
type StringK struct {
	Comb string
	src  string
	// This is private because it violates the law of Demeter
	i    Iterator
}

func NewStringK(src string, k int) StringK {
	return StringK{
		src: src,
		i: Iterator{
			N: len(src),
			K: k,
		},
	}
}

func (s *StringK) Next() bool {
	if !s.i.Next() {
		return false
	}

	runeSrc := []rune(s.src)
	runeDest := make([]rune, len(s.i.Comb))
	for i, v := range s.i.Comb {
		runeDest[i] = runeSrc[v]
	}
	s.Comb = string(runeDest)

	return true
}

// String iterates through all sub-string combinations of its source
// from shortest to longest.
type String struct {
	sk   StringK
	Comb string
}

func NewString(src string) String {
	return String{
		sk: NewStringK(src, 1),
	}
}

func (s *String) Next() bool {
	if !s.sk.Next() {
		if s.sk.i.K+1 > len(s.sk.src) {
			return false
		}
		s.sk.i.K++
		s.sk.Next()
	}
	// This is cheap because a string is just a pointer
	s.Comb = s.sk.Comb
	return true
}
