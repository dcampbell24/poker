package cfr

import (
	"testing"
	"fmt"
)

func ExampleCFR() {
	actEV := []float64{-3, 6, 9}
	strat := []float64{1.0/3, 1.0/3, 1.0/3}
	fmt.Println(CFR(actEV, strat, 0.5))
	// Output: [-3.5 1 2.5]
}

func ExampleNewStrategy() {
	fmt.Println(NewStrategy([]float64{-3.5, 1, 2.5}))
	// Output: [0 0.2857142857142857 0.7142857142857143]
}

func TestCalcNash(_ *testing.T) {
	CalcNash(0)
}
