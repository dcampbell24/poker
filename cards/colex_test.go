package cards

import (
	"testing"
	"fmt"
	"poker/util"
	"time"
)

type hand struct {
	count    int
	strength float64
}

// 2 cards =     169 hands
// 5 cards = 134,459 hands
func TestIso(_ *testing.T) {
	colex := newColex()
	canonical := newCanonical()
	var count int
	t0 := time.Now()
	hands := make([]hand, binomial(52, 5))
	hand := make([]int32, 5)
	deck := NewDeck()
	c1 := util.Comb(deck, 5)
	for loop := true; loop; {
		loop = c1(hand)
		hands[colex(canonical(hand))].count++
	}
	t1 := time.Now()
	for _, v := range hands {
		if v.count != 0 {
			count++
		}
	}
	fmt.Println(t1.Sub(t0), count)
}
