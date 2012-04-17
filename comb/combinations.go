// Package comb provides functions for calculating and interating over
// combinations.
package comb

import (
	"math/big"
)

// Comb takes a slice and a number of elements to choose from that slice and
// returns a function that given a slice will fill it in with a combination
// until none are left, at which point it will return false.
//
// This algorithm is based on one from TAOCP Vol. 4 by Donald Knuth.
func Generator(a []int32, k int32) func([]int32) bool {
	// There is one way to choose 0 -- the empty set.
	if k == 0 {
		return func(_ []int32) bool {
			return false
		}
	}
	var i, j, x int32
	c := make([]int32, k+3, k+3)
	for i = 1; i <= k; i++ {
		c[i] = i
	}
	c[k+1] = int32(len(a)) + 1
	c[k+2] = 0
	j = k
	return func(v []int32) bool {
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
	if n.Cmp(big.NewInt(2)) == -1 {
		return big.NewInt(1)
	}

	sum := big.NewInt(1)
	n0 := new(big.Int).Set(n)
	for n0.Cmp(big.NewInt(1)) != -1 {
		sum.Mul(sum, n0)
		n0.Sub(n0, big.NewInt(1))
	}
	return sum
}

// div returns m! / n!.
// This is an optimization to avoid calculating large factorials.
// Note: n must be less than m.
func div(m, n *big.Int) *big.Int {
	n0 := new(big.Int).Set(n)
	sum := new(big.Int).Add(n0, big.NewInt(1))
	n0.Add(n0, big.NewInt(2))
	for n0.Cmp(m) != 1  {
		sum.Mul(sum, n0)
		n0.Add(n0, big.NewInt(1))
	}
	return sum
}

// Choose returns the number of combinations for n choose k.
func Count(n, k *big.Int) *big.Int {
	if n.Cmp(k) == 0 {
		return big.NewInt(1)
	}
	s := new(big.Int).Sub(n, k)
	if k.Cmp(s) == 1 {
		return s.Quo(div(n, k), Fact(s))
	}
	return s.Quo(div(n, s), Fact(k))
}
