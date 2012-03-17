package game

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

const ACPCversion = "VERSION:2.0.0\r\n"

// The fields of the ACPC state string.
const (
	_ = iota
	POSITION
	HAND_NUM
	BETS
	CARDS
)

type Player interface {
	Play(g *Game) (action byte)
}

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

type Cards struct {
	Holes, Board []string
}

type Game struct {
	position     int     // Position relative to the dealer button.
	handNumber   string  // Unique identifier for each hand.
	bets         string  // A String representing all betting actions.
	*Cards               // A vector of cards of the form [[holes][board]].
	Pot          float64 // The current amount of chips in the pot.
	Call         float64
	Raise        float64
	folded       []bool // Whether a player has folded or not.
	activePlayer int    // The player whose turn it is to act.
	*Rules              // The set of rules to use to play the game.
	round        int    // The current round of betting.
	numActive    int    // The number of players who have not folded.
}

func newGame(r *Rules) *Game {
	return &Game{folded: make([]bool, r.numPlayers), Rules: r}
}

func splitCards(s string) []string {
	a := make([]string, len(s)/2)
	na := 0
	for i, c := range s {
		if c == 'c' || c == 'd' || c == 'h' || c == 's' {
			a[na] = s[i-1 : i+1]
			na++
		}
	}
	return a[:na]
}

// Convert an ACPC cards string into Cards.
func parseCards(s string) *Cards {
	s1 := strings.SplitN(s, "/", 2) // Split into [holeCards boardCards]
	if len(s1) == 1 {
		return &Cards{Holes: splitCards(s1[0])}
	}
	return &Cards{splitCards(s1[0]), splitCards(s1[1])}
}

// RoundActions returns a string of all the actions taken during the current betting round.
func (g *Game) RoundActions() string {
	if g.bets == "" {
		return ""
	}
	s := strings.Split(g.bets, "/")
	return s[len(s)-1]
}

// nextPlayer returns the seat of the next player who can still take an action.
func (g *Game) nextPlayer() int {
	switch {
	// Showdown or all but one folded
	case (len(g.Cards.Holes) > 2) || (g.numActive < 2):
		return -1
	// Start of a new betting round
	case g.RoundActions() == "":
		return g.Rules.firstPlayer[g.round] - 1
		fmt.Println("New Round")
	}

	i := g.activePlayer
	for {
		i = (i + 1) % len(g.folded)
		if !g.folded[i] {
			return i
		}
	}
	panic("nextPlayer: logic error")
}

// betsByRound returns the bets in the pot by round for a two player game.
func betsByRound(actions string) []float64 {
	pot, call, accum := 1.5, 0.5, make([]float64, 0, 4)
	for _, j := range actions {
		switch j {
		case '/':
			pot, call, accum = 0.0, 0.0, append(accum, pot)
		case 'c':
			pot, call = (pot + call), 0.0
		case 'r':
			pot, call = (pot + call + 1.0), 1.0
		}
	}
	return append(accum, pot)
}

// calcPot returns the amount of chips in the pot.
func (g *Game) calcPot() float64 {
	sum := 0.0
	for i, b := range betsByRound(g.bets) {
		sum += b * g.Rules.raiseSize[i]
	}
	return sum
}

// Update the card game using the match-state string."
func (g *Game) updateGame(s string) {
	var err error
	state := strings.Split(s, ":")
	g.position, err = strconv.Atoi(state[POSITION])
	if err != nil {
		panic(err.Error())
	}
	g.bets = state[BETS]
	g.Cards = parseCards(state[CARDS])
	g.round = strings.Count(state[CARDS], "/")
	g.Pot = g.calcPot()

	// Determine Call and Raise sizes.
	if len(g.bets) == 0 {
		g.Call = float64(g.Rules.raiseSize[g.round]) * 0.5
	} else if g.bets[len(g.bets)-1] == 'r' {
		g.Call = float64(g.Rules.raiseSize[g.round])
	} else {
		g.Call = 0.0
	}
	g.Raise = g.Call + g.Rules.raiseSize[g.round]

	// Determine who has folded.
	if state[HAND_NUM] != g.handNumber {
		for i := 0; i < len(g.folded); i++ {
			g.folded[i] = false
		}
		g.numActive = len(g.folded)
	} else if g.bets[len(g.bets)-1] == 'f' {
		g.folded[g.activePlayer] = true
		g.numActive--
	}
	g.activePlayer = g.nextPlayer()
	g.handNumber = state[HAND_NUM]
}

// Start playing a game.
//	rules -- a String naming the game to play
//	p     -- an object that implements the Player interface.
//	host  -- the InetAddress of the dealer passed as a String
//	port  -- the port the dealer is listening on for the client passed as a String
func Play(rules string, p Player, host, port string) {
	addr := net.JoinHostPort(host, port)
	fmt.Printf("Connecting to dealer at %s...\n", addr)
	conn, err := net.Dial("tcp", addr)
	defer conn.Close()
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Starting a game of %s...\n", rules)
	conn.Write([]byte(ACPCversion))
	game := newGame(Holdem2p)
	reader := bufio.NewReader(conn)
	msg := []byte("")
	for {
		piece, frag, err := reader.ReadLine()
		if frag {
			msg = append(msg, piece...)
			continue
		}
		msg = append(msg, piece...)
		switch {
		case err == io.EOF:
			fmt.Println("Shutting down...")
			return
		case err != nil:
			panic(err.Error())
		// ";" and "#" are comment lines.
		case len(msg) < 1 || msg[0] == ';' || msg[0] == '#':
			continue
		}
		game.updateGame(string(msg))
		if game.activePlayer == game.position {
			conn.Write(append(msg, ':', p.Play(game), '\r', '\n'))
		}
		msg = []byte("")
	}
}

// LegalActions returns a byte slice containg the currently legal actions.
func (g *Game) LegalActions() []byte {
	actions := []byte("c")
	if g.Call > 0 {
		actions = append(actions, 'f')
	}
	if strings.Count(g.RoundActions(), "r") < g.Rules.maxRaises[g.round] {
		actions = append(actions, 'r')
	}
	return actions
}
