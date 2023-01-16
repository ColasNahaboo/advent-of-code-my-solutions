const MaxUint = ^uint(0)
const MinUint = 0
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1

// an easier to spot maxint in debug than 9223372036854775807 (^uint(0) >> 1)
const maxint = 8888888888888888888

// Get all prime factors of a given number n
// from: https://siongui.github.io/2017/05/09/go-find-all-prime-factors-of-integer-number/

func primeFactors(n int) (pfs []int) {
	// Get the number of 2s that divide n
	for n%2 == 0 {
		pfs = append(pfs, 2)
		n = n / 2
	}
	// n must be odd at this point. so we can skip one element
	// (note i = i + 2)
	for i := 3; i*i <= n; i = i + 2 {
		// while i divides n, append i and divide n
		for n%i == 0 {
			pfs = append(pfs, i)
			n = n / i
		}
	}
	// This condition is to handle the case when n is a prime number
	// greater than 2
	if n > 2 {
		pfs = append(pfs, n)
	}
	return
}

// Get all divisors of a given number n

func allDivisors(n int) (divs []int) {
	divs = []int{1, n}
	i := 2
	for {
		r := n % i
		if r == 0 {
			j := n / i
			if j < i {
				break
			} else if j == i {
				divs = append(divs, i)
				break
			} else {
				divs = append(divs, i, j)
			}
		}
		i++
	}
	return
}

// greatest common divisor (GCD) via Euclidean algorithm
func greatestCommonDivisor(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func leastCommonMultiple(a, b int, integers ...int) int {
	result := a * b / greatestCommonDivisor(a, b)
	for i := 0; i < len(integers); i++ {
		result = leastCommonMultiple(result, integers[i])
	}
	return result
}
