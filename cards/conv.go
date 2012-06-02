package cards

const (
	Ranks = "23456789TJQKA"
	Suits = "cdhs"
)

var _STOI map[string]int32
var _ITOS [53]string

func init() {
	// Initialize STOI and ITOS
	_STOI = make(map[string]int32, 52)
	var k int32 = 1
	for i := range Ranks {
		for j := range Suits {
			card := string([]byte{Ranks[i], Suits[j]})
			_STOI[card] = k
			_ITOS[k] = card
			k++
		}
	}
}

// Safe subtraction of integer sets (cards).
func Minus(a, b []int32) []int32 {
	c := make([]int32, 0, len(a))
loop:
	for _, v := range a {
		for _, w := range b {
			if v == w {
				continue loop
			}
		}
		c = append(c, v)
	}
	return c
}

func NewDeck(missing ...int32) []int32 {
	deck := make([]int32, 52)
	for i := 0; i < 52; i++ {
		deck[i] = int32(i + 1)
	}
	if len(missing) > 0 {
		deck = Minus(deck, missing)
	}
	return deck
}

func StoI(cards []string) []int32 {
	ints := make([]int32, len(cards))
	for i, c := range cards {
		ints[i] = _STOI[c]
	}
	return ints
}

func ItoS(ints []int32) []string {
	cards := make([]string, len(ints))
	for i, c := range ints {
		cards[i] = _ITOS[c]
	}
	return cards
}
