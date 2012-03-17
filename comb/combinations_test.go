package comb

import (
	"testing"
	"math/big"
)

func TestGenerator(test *testing.T) {
	nums := []uint32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	expTotal := Count(big.NewInt(int64(len(nums))), big.NewInt(5))
	actTotal := new(big.Int)

	c := Generator(nums, 5)
	loop := true
	vals := make([]uint32, 5)
	for loop {
		loop = c(vals)
		actTotal.Add(actTotal, ONE)
	}
	if expTotal.Cmp(actTotal) != 0 {
		test.Fatalf("Expected %v combinations, but saw %v\n", expTotal, actTotal)
	}
}
