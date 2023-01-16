
// Enumerate all the subsets of size k from a set of size n
// Fast algortihm via bitsets, but only for n < 64
// see http://math0.wvstateu.edu/~baker/cs405/code/Subset.html

func subsetsList(n, k int) []int {
	x := intPower(2, k) - 1
	max := intPower(2, n) - intPower(2, n-k)
	s := make([]int, 1, max-x)
	s[0] = x
	for x < max {
		u := x & (-x)
		v := x + u
		x = v + (((v ^ x) / u) >> 2)
		s = append(s, x)
	}
	return s
}
