package game

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

const ACPCversion = "VERSION:2.0.0\r\n"

type Cards struct {
	Holes, Board []string
}

func (this *Cards) String() string {
	return fmt.Sprintf("%v  %v", this.Holes, this.Board)
}

type ACPCString struct {
	position int    // Position relative to the dealer button.
	handNum  string // Unique identifier for each hand.
	bets     string // A String representing all betting actions.
	round    int    // 0-3: pre-flop, flop, turn, river.
	*Cards
}

// RoundActions returns a string of all the actions taken during the current betting round.
func (this *ACPCString) RoundActions() string {
	if this.bets == "" {
		return ""
	}
	s := strings.Split(this.bets, "/")
	return s[len(s)-1]
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

// The fields of the ACPC state string.
const (
	_ = iota
	POSITION
	HAND_NUM
	BETS
	CARDS
)

func NewACPCString(s string) (*ACPCString, error) {
	var err error
	acpc := new(ACPCString)
	state := strings.Split(s, ":")
	acpc.position, err = strconv.Atoi(state[POSITION])
	if err != nil {
		return nil, err
	}
	acpc.handNum = state[HAND_NUM]
	acpc.bets = state[BETS]
	acpc.round = strings.Count(state[BETS], "/")
	// Split into [holeCards boardCards]
	if s := strings.SplitN(state[CARDS], "/", 2); len(s) == 1 {
		acpc.Cards = &Cards{Holes: splitCards(s[0])}
	} else {
		acpc.Cards = &Cards{splitCards(s[0]), splitCards(s[1])}
	}
	return acpc, nil
}


type Player interface {
	Play(g *Game) (action byte)
}

type Game struct {
	*ACPCString          // Basic game state from the server string.
	Pot          float64 // The current amount of chips in the pot.
	Call         float64 // The amount to Call
	Raise        float64 // The amount to Raise
	lstAct       []byte  // The last action a player took.
	activePlayer int     // The player whose turn it is to act.
	*Rules               // The set of rules to use to play the game.
	numActive    int     // The number of players who have not folded.
}

func (this *Game) String() string {
	b := new(bytes.Buffer)
	fmt.Fprintf(b, "Pot: %.2f\n", this.Pot)
	fmt.Fprintf(b, "%v\n[", this.Cards)

	for i, action := range this.lstAct {
		if i == this.activePlayer {
			fmt.Fprintf(b, " (%s) ", string(action))
		} else {
			fmt.Fprintf(b, " %s ", string(action))
		}
	}
	b.WriteString("]\n")
	fmt.Fprintf(b, "Call: %.2f Raise: %.2f\n", this.Call, this.Raise)
	return b.String()
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
		i = (i + 1) % len(g.lstAct)
		if g.lstAct[i] != 'f' {
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
	state, err := NewACPCString(s)
	if err != nil {
		log.Fatalln("updateGame: malformed state string:", err)
	}
	newHand := g.handNum != state.handNum
	g.ACPCString = state
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

	// Update players with the last action each took.
	if newHand {
		for i := 0; i < len(g.lstAct); i++ {
			g.lstAct[i] = '#'
		}
		g.numActive = len(g.lstAct)
	} else if len(g.RoundActions()) == 0 {
		for i := 0; i < len(g.lstAct); i++ {
			if g.lstAct[i] != 'f' {
				g.lstAct[i] = '/'
			}
		}
	} else {
		g.lstAct[g.activePlayer] = g.bets[len(g.bets)-1]
		if g.lstAct[g.activePlayer] == 'f' {
			g.numActive--
		}
	}
	g.activePlayer = g.nextPlayer()
	//fmt.Println(g)
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
		log.Fatalln(err)
	}

	fmt.Printf("Starting a game of %s...\n", rules)
	conn.Write([]byte(ACPCversion))

	// FIXME: don't hard code the ruleset.
	game := new(Game)
	game.lstAct = make([]byte, Holdem2p.numPlayers)
	game.Rules = Holdem2p
	game.ACPCString = &ACPCString{handNum:"nil"}

	reader := bufio.NewReader(conn)
	msg := []byte("")
	//actions := map[byte]int{'f', 'c', 'r'}
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
			log.Fatalln(err)
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
