package util

import (
	"testing"
	"math/big"
)

func TestComb(test *testing.T) {
	nums := []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	expTotal := new(big.Int)
	expTotal.Binomial(int64(len(nums)), 5)

	actTotal := new(big.Int)
	c := Comb(nums, 5)
	loop := true
	vals := make([]int32, 5)
	for loop {
		loop = c(vals)
		actTotal.Add(actTotal, big.NewInt(1))
	}
	if expTotal.Cmp(actTotal) != 0 {
		test.Fatalf("Expected %v combinations, but saw %v\n", expTotal, actTotal)
	}
}
