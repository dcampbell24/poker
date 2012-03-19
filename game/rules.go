package game

type Rules struct {
	descr       string
	limit       bool
	numPlayers  int
	stack       []int
	blind       []int
	raiseSize   []float64
	firstPlayer []int
	maxRaises   []int
}

var (
	Holdem2p = &Rules{
		descr:       "two player limit Texas Hold'em",
		limit:       true,
		numPlayers:  2,
		blind:       []int{10, 5},
		raiseSize:   []float64{10.0, 10.0, 20.0, 20.0},
		firstPlayer: []int{2, 1, 1, 1},
		maxRaises:   []int{3, 4, 4, 4},
	}

	Holdem3p = &Rules{
		descr:       "three player limit Texas Hold'em",
		limit:       true,
		numPlayers:  3,
		blind:       []int{5, 10, 0},
		raiseSize:   []float64{10.0, 10.0, 20.0, 20.0},
		firstPlayer: []int{3, 1, 1, 1},
		maxRaises:   []int{3, 4, 4, 4},
	}

	HoldemNolimit2p = &Rules{
		descr:       "two player no limit Texas Hold'em",
		limit:       false,
		numPlayers:  2,
		stack:       []int{20000, 20000},
		blind:       []int{100, 50},
		firstPlayer: []int{2, 1, 1, 1},
	}

	HoldemNolimit3p = &Rules{
		descr:       "three player no limit Texas Hold'em",
		limit:       false,
		numPlayers:  3,
		stack:       []int{20000, 20000, 20000},
		blind:       []int{50, 100, 0},
		firstPlayer: []int{3, 1, 1, 1},
	}
)
