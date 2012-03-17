package comb

import (
	"math/big"
)

var (
	ONE = big.NewInt(1)
	TWO = big.NewInt(2)
)

// Comb takes a slice and a number of elements to choose from that slice and
// returns a function that given a slice will fill it in with a combination
// until none are left, at which point it will return false.
//
// This algorithm is based on one from TAOCP Vol. 4 by Donald Knuth.
func Generator(a []uint32, k uint32) func([]uint32) bool {
	// There is one way to choose 0 -- the empty set.
	if k == 0 {
		return func(_ []uint32) bool {
			return false
		}
	}
	var i, j, x uint32
	c := make([]uint32, k+3, k+3)
	for i = 1; i <= k; i++ {
		c[i] = i
	}
	c[k+1] = uint32(len(a)) + 1
	c[k+2] = 0
	j = k
	return func(v []uint32) bool {
		for i = k; i > 0; i-- {
			v[k-i] = a[c[i]-1]
		}

		if j > 0 {
			x = j + 1
			goto incr
		}
		if c[1]+1 < c[2] {
			c[1] += 1
			return true
		}
		j = 2
do_more:
		c[j-1] = j - 1
		x = c[j] + 1
		if x == c[j+1] {
			j++
			goto do_more
		}
		// If true, the algorithm is done.
		if j > k {
			return false
		}
incr:
		c[j] = x
		j--
		return true
	}
}

// Fact returns n! (factorial).
// Returns 1 for negative values.
func Fact(n *big.Int) *big.Int {
	if n.Cmp(TWO) == -1 {
		return ONE
	}

	sum := big.NewInt(1)
	for n.Cmp(ONE) != -1 {
		sum.Mul(sum, n)
		n.Sub(n, ONE)
	}
	return sum
}

// div returns m! / n!.
// This is an optimization to avoid calculating large factorials.
// Note: n must be less than m.
func div(m, n *big.Int) *big.Int {
	sum := new(big.Int).Add(n, ONE)
	n.Add(n, TWO)
	for n.Cmp(m) != 1  {
		sum.Mul(sum, n)
		n.Add(n, ONE)
	}
	return sum
}

// Choose returns the number of combinations for n choose k.
func Count(n, k *big.Int) *big.Int {
	if n.Cmp(k) == 0 {
		return ONE
	}
	s := new(big.Int).Sub(n, k)
	if k.Cmp(s) == 1 {
		return s.Quo(div(n, k), Fact(s))
	}
	return s.Quo(div(n, s), Fact(k))
}
