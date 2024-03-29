// Generate all permutations
// Code by Paul Hankin copied from https://stackoverflow.com/a/30230552

package main

func nextPerm(p []int) {
    for i := len(p) - 1; i >= 0; i-- {
        if i == 0 || p[i] < len(p)-i-1 {
            p[i]++
            return
        }
        p[i] = 0
    }
}

func getPerm(orig, p []int) []int {
	result := make([]int, len(orig), len(orig))
	copy(result, orig)
    for i, v := range p {
        result[i], result[i+v] = result[i+v], result[i]
    }
    return result
}


// example of use:
// func main() {
//     orig := []int{11, 22, 33}
//     for p := make([]int, len(orig)); p[0] < len(p); nextPerm(p) {
//        fmt.Println(getPerm(orig, p))
//     }
// }

////////////////////////////////////////////////////////////////// Variant:
// Enumerate all the subsets of size k from a set of size n
// Fast algortihm via bitsets, but only for n < 64
// see http://math0.wvstateu.edu/~baker/cs405/code/Subset.html

func subsetsList(n, k int) []int {
	x := twoPower(2, k) - 1
	max := twoPower(2, n) - twoPower(2, n-k)
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

func twoPower(n int) int {
	if n <= 0 {
		return 1
	}
	result := 2
	for i := 2; i <= n; i++ {
		result *= 2
	}
	return result
}
