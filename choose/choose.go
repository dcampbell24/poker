package main

import (
	"math/big"
	"fmt"
	"os"
	"poker/comb"
)

func main() {
	n := new(big.Int)
	k := new(big.Int)
	n.SetString(os.Args[1], 10)
	k.SetString(os.Args[2], 10)
	fmt.Println(comb.Count(n, k))
}
