package bayes

import (
	"fmt"
	"testing"
)

func ExampleNewLottery() {
	lotto := NewLottery(map[string]float64{"a": 0.4, "b": 0.1, "c": 0.5, "d": 0})
	fmt.Println(lotto)
	for i := 0; i < 100; i++ {
		lotto.Play()
	}
	// Output:
	// [ a:0.40 c:0.90 b:1.00 ]
}

func TestParseDist (test *testing.T) {
	ako := &HandDist{"AKo"}
	if l := len(ako.Strs()); l != 12 {
		test.Fatalf("AKo should produce 12 hands, but produced %d\n", l)
	}
	aa := &HandDist{"AA"}
	if l := len(aa.Strs()); l != 6 {
		test.Fatalf("AA should produce 6 hands, but produced %d\n", l)
	}
	aks := &HandDist{"AKs"}
	if l := len(aks.Strs()); l != 4 {
		test.Fatalf("AKs should produce 4 hands, but produced %d\n", l)
	}
}

func TestPHole(test *testing.T) {
	p1 := 0.004524886877828055
	p2 := 0.0024489795918367346
	if p := PHole(&HandDist{"AA"}, []string{}); p != p1 {
		test.Fatalf("P(AA) should have been %f, but was %f\n", p1, p)
	}
	if p := PHole(&HandDist{"AA"}, []string{"As", "Ks"}); p != p2  {
		test.Fatal("P(AA | AKs) should have been %f, but was %f.\n", p2, p)
	}
}
