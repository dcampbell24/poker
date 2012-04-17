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
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"io"
	"log"
	"runtime"

	"poker/comb"
)

const (
	ranks = "23456789TJQKA"
	suits = "cdhs"
)

var hr [32487834]int32
var CTOI map[string]int32
// How many cpus to use for the equity calculations.
var NCPU int
var RANDS []*rand.Rand

func init() {
	fmt.Print("Loading HandRanks.dat... ")
	// Initialize hr
	buf := make([]byte, len(hr)*4, len(hr)*4)
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

	// Initialize CTOI
	CTOI = make(map[string]int32, 52)
	var k int32 = 1
	for i := 0; i < len(ranks); i++ {
		for j := 0; j < len(suits); j++ {
			CTOI[string([]byte{ranks[i], suits[j]})] = k
			k++
		}
	}

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

func NewDeck(missing ...int32) []int32 {
	deck := make([]int32, 52, 52)
	for i := 0; i < 52; i++ {
		deck[i] = int32(i + 1)
	}
	if len(missing) > 0 {
		deck = minus(deck, missing)
	}
	return deck
}

func cardsToInts(cards []string) []int32 {
	ints := make([]int32, len(cards), len(cards))
	for i, c := range cards {
		ints[i] = CTOI[c]
	}
	return ints
}

// A hand distribution is a category of hands. Currently the only categories
// supported are those of the forms AA, AKo, and AKs.
type HandDist struct {
	Dist string
}

// Expand the HandDist into a slice of all possible hands.
func (this *HandDist) Strs() [][]string {
	hands := make([][]string, 0)
	// Expand each card into a card of each suit
	xs := make([]string, 4)
	ys := make([]string, 4)
	for i := range suits {
		xs[i] = string([]byte{this.Dist[0], suits[i]})
		ys[i] = string([]byte{this.Dist[1], suits[i]})
	}
	switch {
	case len(this.Dist) == 2:
		// pairs e.g. AA
		for i := 0; i < 3; i++ {
			for j := i+1; j < 4; j++ {
				hands = append(hands, []string{xs[i], xs[j]})
			}
		}
	case this.Dist[2] == 'o':
		// offsuit e.g. AKo
		for i := 0; i < 4; i++ {
			for j := 0; j < 4; j++ {
				if i != j {
					hands = append(hands, []string{xs[i], ys[j]})
				}
			}
		}
	default:
		// suited e. g. AKs
		for i := 0; i < 4; i++ {
			hands = append(hands, []string{xs[i], ys[i]})
		}
	}
	return hands
}

// The same as Strs, only return the hands represented by int32.
func (this *HandDist) Ints() [][]int32 {
	shands := this.Strs()
	hands := make([][]int32, len(shands))
	for i := range shands {
		hands[i] = cardsToInts(shands[i])
	}
	return hands
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

func EvalHand(cards []string) int32 {
	hand := cardsToInts(cards)
	return evalHand(evalBoard(hand[:5]), hand[5:])
}


// Split a hand rank into two values: category and rank-within-category.
func SplitRank(rank int32) (int32, int32) {
	return rank >> 12, rank & 0xFFF
}

// Calculate the percent of the pot each hand wins and return them as a slice.
func EvalHands(board []int32, hands ...[]int32) []float64 {
	b := evalBoard(board)
	// Optimize case where there are only two hands.
	if len(hands) == 2 {
		result := evalHand(b, hands[0]) - evalHand(b, hands[1])
		switch {
		case result > 0:
			return []float64{1, 0}
		case result < 0:
			return []float64{0, 1}
		default:
			return []float64{0.5, 0.5}
		}
	}
	vals := make([]int32, len(hands), len(hands))
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
	// Alot each winner his share of the pot.
	result := make([]float64, len(hands), len(hands))
	for i, v := range vals {
		if v == max {
			result[i] = 1.0 / float64(winners)
		} else {
			result[i] = 0.0
		}
	}
	return result
}

// PHole returns the probability of having a given class of hole cards given
// that scards have already been seen.
//
// Here are two example calculations:
//	          Me   Opp  Board  P(AA)
//	Pre-deal  ??   ??   ???    (4 choose 2) / (52 choose 2) ~= 0.0045
//	Pre-flop  AKs  ??   ???    (3 choose 2) / (50 choose 2) ~= 0.0024
func PHole(hd *HandDist, scards []string) float64 {
	holes := hd.Ints()
	// Count how many hands to eliminate from the holes class because a card in
	// that hand has already been seen.
	elim := 0
	for _, hand := range holes {
		for _, card := range hand {
			for _, seen := range cardsToInts(scards) {
				if card == seen {
					elim++
					goto nextHand
				}
			}
		}
		nextHand:
	}
	deck := int64(52 - len(scards))
	allHands := float64(comb.Count(big.NewInt(deck), big.NewInt(2)).Int64())
	return float64(len(holes) - elim) / allHands
}

/*
// FIXME
// CondProbs returns the PTable for P(hole | action) given the cards that have
// currently been seen, and the probabilties P(action) and P(action | hole).
// The formula for calculating the conditional probability P(hole | action):
//
//	                   P(hole) * P(action | hole)
//	P(hole | action) = --------------------------
//	                            P(action)
//
// Weisstein, Eric W. "Conditional Probability." From MathWorld--A Wolfram Web
// Resource. http://mathworld.wolfram.com/ConditionalProbability.html
//
func CondProbs(scards []string, pActHole *PTable, pAction []float64) *PTable {
	for _, vals := range actionDist {
		NewRRSDist(actionDist[:3]...) (* (PHole cards [r1 r2 s]) prob)])]
  (apply array-map (flatten values))))
*/

type Lottery struct {
	// Maybe should use ints or fixed point to make more accurate.
	probs []float64
	prizes []string
}

func (this *Lottery) String() string {
	b := bytes.NewBufferString("[ ")
	for i := 0; i < len(this.probs); i++ {
		fmt.Fprintf(b, "%s:%.2f ", this.prizes[i], this.probs[i])
	}
	b.WriteString("]")
	return b.String()
}

// Convert a discrete distribution (array-map {item prob}) into a lottery. The
// probabilities should add up to 1
func NewLottery(dist map[string] float64) *Lottery {
	sum := 0.0
	lotto := &Lottery{}
	for key, val := range dist {
		if val != 0 {
			sum += val
			lotto.probs = append(lotto.probs, sum)
			lotto.prizes = append(lotto.prizes, key)
		}
	}
	return lotto
}

// Draw a winner from a Lottery. If at least one value in the lottery is not >=
// 1, then the greatest value is effectively rounded up to 1.0"
func (this *Lottery) Play() string {
	draw := rand.Float64()
	for i, p := range this.probs {
		if p > draw {
			return this.prizes[i]
		}
	}
	return this.prizes[len(this.prizes)-1]
}

// Safe subtraction of integer sets.
func minus(a, b []int32) []int32 {
	c := make([]int32, len(a), len(a))
	var count int
	var match bool
	for _, v := range a {
		for _, w := range b {
			if v == w {
				match = true
				break
			}
		}
		if !match {
			c[count] = v
			count++
		}
		match = false
	}
	return c[:count]
}

// Fisher–Yates shuffle
// r is the rand.Rand to use.
func shuffle(a []int32, r int) {
	for i := len(a) - 1; i > 0; i-- {
		j := RANDS[r].Intn(i + 1)
		a[j], a[i] = a[i], a[j]
	}
}

// Get ready to do the hand equity calculations. Returns hand, board, bLen,
// deck.
func handEquityInit(sHand, sBoard []string) ([]int32, []int32, int32, []int32) {
	// Convert the cards from strings to ints.
	hole := cardsToInts(sHand)
	bLen := int32(len(sBoard)) // How many cards will we need to draw?
	board := make([]int32, 5, 5)
	for i, v := range sBoard {
		board[i] = CTOI[v]
	}
	// Remove the hole and board cards from the deck.
	deck := NewDeck(append(hole, board[:bLen]...)...)
	return hole, board, bLen, deck
}

// Exhaustive hand equity calculation.
func handEquityE(hole, board []int32, bLen int32, deck []int32) float64 {
	var sum, count float64
	oHole := make([]int32, 2, 2)
	c1 := comb.Generator(deck, 2)
	for loop1 := true; loop1; {
		loop1 = c1(oHole)
		c2 := comb.Generator(minus(deck, oHole), 5-bLen)
		for loop2 := true; loop2; {
			loop2 = c2(board[bLen:])
			sum += EvalHands(board, hole, oHole)[0]
			count++
		}
	}
	return sum / count
}

// Monte-Carlo hand equity calculation.
func handEquityMC(hole, board []int32, bLen int32, deck []int32, trials, r int) float64 {
	var sum float64
	for i := 0; i < trials; i++ {
		shuffle(deck, r)
		copy(board[bLen:], deck[2:8-bLen])
		sum += EvalHands(board, hole, deck[:2])[0]
	}
	return sum / float64(trials)
}

// Parallel Monte-Carlo hand equity calculation.
func handEquityMCP(hole, board []int32, bLen int32, deck []int32, trials int,
c chan float64, i int) {
	c <- handEquityMC(hole, board, bLen, deck, trials, i)
}

// HandEquity returns the equity of a player's hand based on the current
// board.  trials is the number of Monte-Carlo simulations to do.  If trials
// is 0, then exhaustive enumeration will be used instead.
func HandEquity(sHand, sBoard []string, trials int) float64 {
	hole, board, Blen, deck := handEquityInit(sHand, sBoard)
	if trials == 0 {
		return handEquityE(hole, board, Blen, deck)
	}
	return handEquityMC(hole, board, Blen, deck, trials, 0)
}

// Parallel version of HandEquity.
func HandEquityP(sHand, sBoard []string, trials int) float64 {
	trials += trials % NCPU // Round to a multiple of the number of CPUs.
	c := make(chan float64) // Not buffering
	for i := 0; i < NCPU; i++ {
		hole, board, bLen, deck := handEquityInit(sHand, sBoard)
		go handEquityMCP(hole, board, bLen, deck, trials/NCPU, c, i)
	}
	sum := 0.0
	for i := 0; i < NCPU; i++ {
		sum += <-c
	}
	return sum / float64(NCPU)
}
