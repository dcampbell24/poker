package equity

import (
	"testing"
	"fmt"
)

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

func TestNewLottery(_ *testing.T) {
	lotto := NewLottery(map[string]float64{"a": 0.4, "b": 0.1, "c": 0.5, "d": 0})
	fmt.Println(lotto)
	for i := 0; i < 10; i++ {
		fmt.Printf("%d:%s ", i, lotto.Play())
	}
	fmt.Println()
}

func checkCategory(expected int32, hand []string, test *testing.T) {
	cat, _ := SplitRank(EvalHand(hand))
	if cat != expected {
		test.Fatalf("The hand %v should be category %d, but was %d.\n", hand, expected, cat)
	}
}

func TestEvalHand(test *testing.T) {
	checkCategory(9, []string{"2c", "3c", "4c", "5c", "6c", "7c", "9c"}, test) // straight-flush
	checkCategory(8, []string{"Ac", "Ah", "Ad", "As", "7c", "8c", "9c"}, test) // four-of-a-kind
	checkCategory(7, []string{"2c", "2h", "2d", "Ts", "Tc", "Js", "Jh"}, test) // full-house
	checkCategory(6, []string{"2c", "3c", "4c", "5c", "7c", "8c", "9c"}, test) // flush
	checkCategory(5, []string{"2s", "3h", "4c", "5d", "6c", "8c", "9c"}, test) // straight
	checkCategory(4, []string{"5c", "5h", "5s", "Td", "7c", "8c", "9c"}, test) // three-of-a-kind
	checkCategory(3, []string{"2c", "2d", "4h", "4s", "7c", "8c", "9c"}, test) // two-pair
	checkCategory(2, []string{"Ac", "Ad", "4s", "5h", "7d", "8c", "9c"}, test) // pair
	checkCategory(1, []string{"Ac", "3d", "4s", "Tc", "7c", "5d", "9c"}, test) // high-card

	hp := []string{"Ac", "Ad", "4s", "5h", "7d", "8c", "9c"} // high pair
	lp := []string{"2c", "2d", "4s", "5h", "7d", "8c", "9c"} // low pair
	if EvalHand(hp) <= EvalHand(lp) {
		test.Fatalf("The high pair %v did not beat the low pair %v.\n", hp, lp)
	}
}

func TestEvalHands(_ *testing.T) {
	board := cardsToInts([]string{"4s", "5h", "7d", "8c", "9c"})
	hp := cardsToInts([]string{"Ac", "Ad"})
	lp := cardsToInts([]string{"2c", "2d"})
	fmt.Println(EvalHands(board, hp, lp))
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

func calcErr(hand, board []string, trials int, exp float64) {
	fmt.Printf("%v  %v  %10d    %+f\n", hand, board, trials, HandEquity(hand, board, trials) - exp)
}

func calcErrP(hand, board []string, trials int, exp float64) {
	fmt.Printf("%v  %v  %10d    %+f\n", hand, board, trials, HandEquityP(hand, board, trials) - exp)
}
func testMCHE(hand, board []string) {
	equity := HandEquity(hand, board, 0)
	fmt.Printf("%-8s %-14s %-10s %-10s\n", "Hand", "Board", "Trials", "Error")
	calcErr(hand, board, 1000, equity)
	calcErr(hand, board, 5000, equity)
	calcErr(hand, board, 10000, equity)
	calcErr(hand, board, 50000, equity)
	fmt.Println()
	calcErrP(hand, board, 1000, equity)
	calcErrP(hand, board, 5000, equity)
	calcErrP(hand, board, 10000, equity)
	calcErrP(hand, board, 50000, equity)
	fmt.Println()
}

func printHE(hand []string, trials int) {
	fmt.Println(hand, trials, HandEquity(hand, []string{}, trials))
}

func TestHE(_ *testing.T) {
	testMCHE([]string{"7d", "6c"}, []string{"2c", "2d", "3s"})
	testMCHE([]string{"Ad", "Kd"}, []string{"2c", "2d", "3s"})
	printHE([]string{"Ad", "Ac"}, 10000)
	printHE([]string{"2d", "2c"}, 10000)
	printHE([]string{"As", "Ks"}, 10000)
}

func TestNewDeck(_ *testing.T) {
	fmt.Println(NewDeck(1, 2, 3, 4, 5))
}
