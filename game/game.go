package game

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const ACPCversion = "VERSION:2.0.0\r\n"

const (
	PreFlop = iota
	Flop
	Turn
	River
	Showdown
)

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

type Player interface {
	Play(g *Game) (action string)
	Observe(g *Game)
}

type Game struct {
	Bets      []float64 // The chips put in for each player this round.
	Holes     []string  // All of the viewable hole cards.
	Board     []string  // All of the board cards.
	Raises    int       // The number of raises this round.
	Folded    []bool    // Whether the player in the nth position has folded.
	Actor     int       // The player whose turn it is to act.
	*GameDiff           // The most recent changes in game state.
	*Rules              // The set of rules to use to play the game.
	pot       float64   // Chips in the pot from previous rounds.
}

func NewGame(rules string) (*Game, error) {
	r, err := ChooseRules(rules)
	if err != nil {
		return nil, err
	}
	return &Game{
		Rules:    r,
		Folded:   make([]bool, r.numPlayers),
		Bets:     make([]float64, r.numPlayers),
		GameDiff: new(GameDiff)}, nil
}

func (this *Game) String() string {
	s := fmt.Sprintln(this.Holes, this.Board)
	if this.Actor != -1 {
		s += fmt.Sprintln(this.Pot(), this.Bets, this.CallAmt(), this.RaiseAmt())
	}
	return s
}

// NumActive returns how many players are still in the hand.
func (this *Game) NumActive() int {
	var count int
	for _, folded := range this.Folded {
		if !folded {
			count++
		}
	}
	return count
}

// LegalActions returns a string containing the currently legal actions.
func (this *Game) LegalActions() string {
	actions := "c"
	if this.CallAmt() > 0 {
		actions += "f"
	}
	if this.Raises < this.maxRaises[this.Round] {
		actions += "r"
	}
	return actions
}

func (this *Game) CallAmt() float64 {
	var max float64
	for _, chips := range this.Bets {
		if chips > max {
			max = chips
		}
	}
	return max - this.Bets[this.Actor]
}

func (this *Game) RaiseAmt() float64 {
	return this.CallAmt() + this.raiseSize[this.Round]
}

func (this *Game) Pot() float64 {
	var sum float64
	for _, chips := range this.Bets {
		sum += chips
	}
	return this.pot + sum
}

func (this *Game) Update(s string) error {
	err := this.GameDiff.Update(s)
	if err != nil {
		return err
	}
	// Handle action updates.
	switch this.Action {
	case "f":
		this.Folded[this.Actor] = true
	case "c":
		this.Bets[this.Actor] += this.CallAmt()
	case "r":
		this.Raises++
		this.Bets[this.Actor] += this.CallAmt() + this.raiseSize[this.Round]
	}
	if this.NumActive() < 2 {
		this.Actor = -1
	} else {
		i := this.Actor
		for {
			i = (i + 1) % len(this.Folded)
			if !this.Folded[i] {
				this.Actor = i
				break
			}
		}
	}
	// Handle card updates.
	if len(this.Cards) > 0 {
		switch this.Round {
		case PreFlop:
			this.Raises = 0
			this.Actor = this.firstPlayer[this.Round] - 1
			this.pot = 0
			copy(this.Bets, this.blind)
			this.Holes = splitCards(this.Cards)
			this.Board = nil
			for i := range this.Folded {
				this.Folded[i] = false
			}
		case Flop, Turn, River:
			this.Actor = this.firstPlayer[this.Round] - 1
			this.pot = this.Pot()
			for i := range this.Bets {
				this.Bets[i] = 0
			}
			this.Board = append(this.Board, splitCards(this.Cards)...)
		case Showdown:
			this.Actor = -1
			this.Holes = splitCards(this.Cards)
		}
	}
	return nil
}

// Start playing a game.
//	rules -- a String naming the game to play
//	p     -- an object that implements the Player interface.
//	host  -- the InetAddress of the dealer passed as a String
//	port  -- the port the dealer is listening on for the client passed as a String
func Play(rules string, p Player, host, port string) {
	game, err := NewGame(rules)
	if err != nil {
		log.Fatalln(err)
	}
	// Connect to the dealer.
	addr := net.JoinHostPort(host, port)
	fmt.Printf("Connecting to dealer at %s...\n", addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	// Tell the dealer I am ready to start.
	fmt.Printf("Starting a game of %s...\n", rules)
	conn.Write([]byte(ACPCversion))

	// Read replies from dealer one line at a time.
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
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
		msg = strings.TrimRight(msg, "\r\n")
		err = game.Update(msg)
		if err != nil {
			log.Fatalln(err)
		}
		if game.Actor == game.Position {
			fmt.Fprintf(conn, "%s:%s\r\n", msg, p.Play(game))
		} else {
			p.Observe(game)
		}
	}
}
