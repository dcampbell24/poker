package bayes

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"poker/cards"
)

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
	for i := range cards.Suits {
		xs[i] = string([]byte{this.Dist[0], cards.Suits[i]})
		ys[i] = string([]byte{this.Dist[1], cards.Suits[i]})
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
		hands[i] = cards.StoI(shands[i])
	}
	return hands
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
			for _, seen := range cards.StoI(scards) {
				if card == seen {
					elim++
					goto nextHand
				}
			}
		}
		nextHand:
	}
	allHands := new(big.Int)
	allHands.Binomial(int64(52 - len(scards)), 2)
	return float64(len(holes) - elim) / float64(allHands.Int64())
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
