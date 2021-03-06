package equity

import (
	"math"
	"poker/cards"
	"testing"
	"fmt"
)

// Examples

func ExampleEvalHands() {
	board := cards.StoI([]string{"4s", "5h", "7d", "8c", "9c"})
	hp := cards.StoI([]string{"Ac", "Ad"})
	lp := cards.StoI([]string{"2c", "2d"})
	fmt.Println(EvalHands(board, hp, lp))
	// Output: 1
}

// Tests

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

func calcErr(hand, board []string, trials int, exp float64) {
	fmt.Printf("%v  %v  %10d    %+f %+f\n", hand, board, trials,
		exp - HandEquity(hand, board, trials),
		exp - HandEquityP(hand, board, trials))
}

func testMCHE(hand, board []string) {
	equity := HandEquity(hand, board, 0)
	fmt.Printf("%-8s %-14s %-10s %-10s\n", "Hand", "Board", "Trials", "Error")
	calcErr(hand, board,       1, equity)
	calcErr(hand, board,      10, equity)
	calcErr(hand, board,     100, equity)
	calcErr(hand, board,    1000, equity)
	calcErr(hand, board,   10000, equity)
	calcErr(hand, board,  100000, equity)
	fmt.Println()
}

func printHE(hand []string, trials int) {
	fmt.Println(hand, trials, HandEquity(hand, []string{}, trials))
}

var hole  = []string{"Ad", "Ac"}
var flopB = []string{"2c", "2d", "3s"}
var turnB = append(flopB, "7c")
var rivB  = append(turnB, "9c")

func TestHEerr(_ *testing.T) {
	error := 0.0
	perror := 0.0
	for i := 0; i < 1000; i++ {
		d := cards.NewDeck()
		sample(d, 7, 0)
		df := cards.ItoS(d)
		exp := HandEquity(df[:2], df[2:7], 0)
		act := HandEquity(df[:2], df[2:7], 1000)
		pact := HandEquityP(df[:2], df[2:7], 1000)
		error  += math.Abs(exp - act)
		perror += math.Abs(exp - pact)
	}
	fmt.Println(error/1000)
	fmt.Println(perror/1000)
}

func TestHE(_ *testing.T) {
	testMCHE([]string{"7d", "6c"}, flopB)
	testMCHE([]string{"Ad", "Kd"}, flopB)
	printHE([]string{"Ad", "Ac"}, 10000)
	printHE([]string{"2d", "2c"}, 10000)
	printHE([]string{"As", "Ks"}, 10000)
}

// Benchmarks

func BenchmarkHS_F_0(b *testing.B) {
	for i := 0; i < b.N; i++ {
		HandEquity(hole, flopB, 0)
	}
}

func BenchmarkHE_T_0(b *testing.B) {
	for i := 0; i < b.N; i++ {
		HandEquity(hole, turnB, 0)
	}
}

func BenchmarkHE_R_0(b *testing.B) {
	for i := 0; i < b.N; i++ {
		HandEquity(hole, rivB, 0)
	}
}

func BenchmarkHE_F_1000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		HandEquity(hole, flopB, 1000)
	}
}

func BenchmarkHE_T_1000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		HandEquity(hole, turnB, 1000)
	}
}

func BenchmarkHE_R_1000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		HandEquity(hole, rivB, 1000)
	}
}

func BenchmarkHEP_F_1000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		HandEquityP(hole, flopB, 1000)
	}
}

func BenchmarkHEP_T_1000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		HandEquityP(hole, turnB, 1000)
	}
}

func BenchmarkHEP_R_1000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		HandEquityP(hole, rivB, 1000)
	}
}
