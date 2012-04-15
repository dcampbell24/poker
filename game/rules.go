package game

import "fmt"

type Rules struct {
	descr       string
	limit       bool
	numPlayers  int
	stack       []int
	blind       []float64
	raiseSize   []float64
	firstPlayer []int
	maxRaises   []int
}

func ChooseRules(rules string) (*Rules, error) {
	switch rules {
	case "2p-l":
		return &Rules{
			descr:       "two player limit Texas Hold'em",
			limit:       true,
			numPlayers:  2,
			blind:       []float64{10, 5},
			raiseSize:   []float64{10.0, 10.0, 20.0, 20.0},
			firstPlayer: []int{2, 1, 1, 1},
			maxRaises:   []int{3, 4, 4, 4},
		}, nil
	case "3p-l":
		return &Rules{
			descr:       "three player limit Texas Hold'em",
			limit:       true,
			numPlayers:  3,
			blind:       []float64{5, 10, 0},
			raiseSize:   []float64{10.0, 10.0, 20.0, 20.0},
			firstPlayer: []int{3, 1, 1, 1},
			maxRaises:   []int{3, 4, 4, 4},
		}, nil
	case "2p-nl":
		return &Rules{
			descr:       "two player no limit Texas Hold'em",
			limit:       false,
			numPlayers:  2,
			stack:       []int{20000, 20000},
			blind:       []float64{100, 50},
			firstPlayer: []int{2, 1, 1, 1},
		}, nil
	case "3p-nl":
		return &Rules{
			descr:       "three player no limit Texas Hold'em",
			limit:       false,
			numPlayers:  3,
			stack:       []int{20000, 20000, 20000},
			blind:       []float64{50, 100, 0},
			firstPlayer: []int{3, 1, 1, 1},
		}, nil
	}
	return nil, fmt.Errorf("Don't know how to play %s\n", rules)
}
