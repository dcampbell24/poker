// Package equity implements functions and structures for calculating the equity of poker hands.
// The hand evaluation code is based on
// http://www.codingthewheel.com/archives/poker-hand-evaluator-roundup#2p2
//
//	Example hands that can be parsed
//	*********************************
//	String  Combinations  Description
//
//	AJs                4  Any Ace with a Jack of the same suit.
//	77                 6  Any pair of Sevens.
//	T9o               12  Any Ten and Nine of different suits.
//	54                16  Any Five and Four, suited or unsuited.
//
//	AJs+              12  Any Ace with a (Jack through King) of the same suit.
//	77+               48  Any pair greater than or equal to Sevens.
//	T9o-65o           12  Any unsuited connector between 65o and T9o.
//
//	QQ+,AQs+,AK       38  Any pair of Queen or better, any AQs, and any AK
//	                      whether suited or not.
//	AhKh,7h7d          2  Ace-King of Hearts or a pair of red Sevens.
package equity

import (
	"fmt"
	"math/rand"
	"os"
	"io"
	"log"
	"poker/cards"
	"runtime"
)

var hr [32487834]int32
// How many cpus to use for the equity calculations.
var NCPU int
var RANDS []*rand.Rand

func init() {
	fmt.Print("Loading HandRanks.dat... ")
	// Initialize hr
	buf := make([]byte, len(hr)*4)
	fp, err := os.Open("HandRanks.dat")
	if err != nil {
		log.Fatalln(err)
	}
	defer fp.Close()
	_, err = io.ReadFull(fp, buf)
	if err != nil {
		log.Fatalln(err)
	}
	for i := 0; i < len(buf); i += 4 {
		hr[i/4] = int32(buf[i+3])<<24 |
			int32(buf[i+2])<<16 |
			int32(buf[i+1])<<8 |
			int32(buf[i])
	}
	fmt.Println("Done")

	// Initialize the PRNGs
	// BUG(David): The PRNGs are not being used in a way that guarentees their
	// independence and their state space is much smaller than that of a deck of
	// cards.
	NCPU = runtime.NumCPU()
	for i := 0; i < NCPU; i++ {
		RANDS = append(RANDS, rand.New(rand.NewSource(rand.Int63())))
	}
	runtime.GOMAXPROCS(NCPU)
	fmt.Printf("Using %d CPUs\n", NCPU)
}

func evalBoard(cards []int32) int32 {
	v := hr[53+cards[0]]
	v = hr[v+cards[1]]
	v = hr[v+cards[2]]
	v = hr[v+cards[3]]
	return hr[v+cards[4]]
}

func evalHand(b int32, cards []int32) int32 {
	b = hr[b+cards[0]]
	return hr[b+cards[1]]
}

func EvalHand(hand []string) int32 {
	h := cards.StoI(hand)
	return evalHand(evalBoard(h[:5]), h[5:])
}

// Split a hand rank into two values: category and rank-within-category.
func SplitRank(rank int32) (int32, int32) {
	return rank >> 12, rank & 0xFFF
}

// Calculate the percent of the pot the first player has won.
func EvalHands(board []int32, hands ...[]int32) float64 {
	b := evalBoard(board)
	// Optimize case where there are only two hands.
	if len(hands) == 2 {
		s1 := evalHand(b, hands[0])
		s2 := evalHand(b, hands[1])
		if s1 > s2 {
			return 1.0
		}
		if s1 < s2 {
			return 0.0
		}
		return 0.5
	}
	// More than two hands.
	vals := make([]int32, len(hands))
	for i, hand := range hands {
		vals[i] = evalHand(b, hand)
	}
	// Determine the number of winners and their hand.
	winners := 1
	max := vals[0]
	for i := 1; i < len(vals); i++ {
		if v := vals[i]; v > max {
			max = v
			winners = 1
		} else if v == max {
			winners++
		}
	}
	// Calculate how much the first player has won.
	if vals[0] == max {
		return 1.0 / float64(winners)
	}
	return 0.0
}

// Choose k random items from p and put them in the first k positions of p.
func sample(p []int32, k int, r int) []int32 {
	for i := 0; i < k; i++ {
		j := RANDS[r].Intn(len(p) - i)
		p[i], p[i+j] = p[i+j], p[i]
	}
	return p[:k]
}

// Get ready to do the hand equity calculations. Returns hand, board, deck.
func handEquityInit(sHand, sBoard []string) ([]int32, []int32, []int32) {
	hole := cards.StoI(sHand)
	board := cards.StoI(sBoard)
	deck := cards.NewDeck(append(hole, board...)...)
	return hole, board, deck
}

// Exhaustive hand equity calculation.
func handEquityE(hole, board, deck []int32) float64 {
	var sum, count float64
	bLen := int32(len(board))
	board = append(board, make([]int32, 5-bLen)...)
	oHole := make([]int32, 2)
	c1 := Comb(deck, 2)
	for loop1 := true; loop1; {
		loop1 = c1(oHole)
		c2 := Comb(cards.Minus(deck, oHole), 5-bLen)
		for loop2 := true; loop2; {
			loop2 = c2(board[bLen:])
			sum += EvalHands(board, hole, oHole)
			count++
		}
	}
	return sum / count
}

// Monte-Carlo hand equity calculation.
func handEquityMC(hole, board, deck []int32, trials, r int) float64 {
	var sum float64
	bLen := len(board)
	board = append(board, make([]int32, 5-bLen)...)
	for i := 0; i < trials; i++ {
		s := sample(deck, 7-bLen, r)
		copy(board[bLen:], s[2:])
		sum += EvalHands(board, hole, s[:2])
	}
	return sum / float64(trials)
}

// Parallel Monte-Carlo hand equity calculation.
func handEquityMCP(hole, board, deck []int32, trials, r int, c chan float64) {
	c <- handEquityMC(hole, board, deck, trials, r)
}

// HandEquity returns the equity of a player's hand based on the current
// board.  trials is the number of Monte-Carlo simulations to do.  If trials
// is 0, then exhaustive enumeration will be used instead.
func HandEquity(sHand, sBoard []string, trials int) float64 {
	hole, board, deck := handEquityInit(sHand, sBoard)
	if trials == 0 {
		return handEquityE(hole, board, deck)
	}
	return handEquityMC(hole, board, deck, trials, 0)
}

// Parallel version of HandEquity.
func HandEquityP(sHand, sBoard []string, trials int) float64 {
	trials += trials % NCPU // Round to a multiple of the number of CPUs.
	c := make(chan float64) // Not buffering
	for i := 0; i < NCPU; i++ {
		hole, board, deck := handEquityInit(sHand, sBoard)
		go handEquityMCP(hole, board, deck, trials/NCPU, i, c)
	}
	sum := 0.0
	for i := 0; i < NCPU; i++ {
		sum += <-c
	}
	return sum / float64(NCPU)
}
