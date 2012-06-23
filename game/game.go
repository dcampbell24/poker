// Package game keeps track of basic poker game information. Any Player may use
// the game.Play function in order to play a game.
package game

import (
	"fmt"
	"log"
	"net"

	"poker/game/diff"
)

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
	Round     int          // 0-4: pre-flop, flop, turn, river, showdown.
	Bets      [4][]float64 // The chips put in for each player for each round.
	Holes     []string     // All of the viewable hole cards.
	Board     []string     // All of the board cards.
	Raises    int          // The number of raises this round.
	Actions   []byte       // The last action taken by each player.
	Actor     int          // The player whose turn it is to act.
	*Rules                 // The set of rules to use to play the game.
	Event     interface{}  // The most recent event
	*diff.Players
}

// Makes a partial copy of the Game that is specifically suited for creating
// game trees.
func (this *Game) Copy() *Game {
	g := new(Game)
	g.Round = this.Round
	for i := range g.Bets {
		g.Bets[i] = make([]float64, len(this.Bets[0]))
		copy(g.Bets[i], this.Bets[i])
	}
	g.Holes = make([]string, len(this.Holes))
	copy(g.Holes, this.Holes)
	g.Board = make([]string, len(this.Board))
	copy(g.Board, this.Board)
	g.Raises = this.Raises
	g.Actions = make([]byte, len(this.Actions))
	copy(g.Actions, this.Actions)
	g.Actor = this.Actor
	g.Rules = this.Rules
	//Event     interface{}
	g.Players = this.Players
	return g
}

func NewGame(rules string) (*Game, error) {
	r, err := ChooseRules(rules)
	if err != nil {
		return nil, err
	}
	g := new(Game)
	g.Rules = r
	g.Actions = make([]byte, r.numPlayers)
	for i := range g.Bets {
		g.Bets[i] = make([]float64, r.numPlayers)
	}
	return g, nil
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
	for _, action := range this.Actions {
		if action != 'f' {
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
	for _, chips := range this.Bets[this.Round] {
		if chips > max {
			max = chips
		}
	}
	return max - this.Bets[this.Round][this.Actor]
}

func (this *Game) RaiseAmt() float64 {
	return this.CallAmt() + this.raiseSize[this.Round]
}

func (this *Game) Pot() float64 {
	var sum float64
	for _, bets := range this.Bets {
		for _, chips := range bets {
			sum += chips
		}
	}
	return sum
}

// Is there anyone in the hand who has not acted yet? If not, then does everyone
// who has not folded have the same amount in the pot?
// Pre-condition: two players have not folded.
func (this *Game) evenBets() bool {
	var j int
	var b0 float64
	for i, a := range this.Actions {
		if a == 0 {
			return false
		}
		if a != 'f' {
			b0 = this.Bets[this.Round][i]
			j = i
			break
		}
	}
	for i, a := range this.Actions[j+1:] {
		if a == 0 || (a != 'f' && this.Bets[this.Round][i+j+1] != b0) {
			return false
		}
	}
	return true
}

func (this *Game) Update(event interface{}) {
	this.Event = event
	switch e := event.(type) {
	case *diff.Players:
		this.Actor = -1
		this.Players = e
		this.Round = -1
		for _, bets := range this.Bets {
			for i := range bets {
				bets[i] = 0
			}
		}
		copy(this.Bets[0], this.blind)
		this.Board = nil
		for i := range this.Actions {
			this.Actions[i] = 0
		}
	case diff.Cards:
		for i, a := range this.Actions {
			if a != 'f' {
				this.Actions[i] = 0
			}
		}
		cards := splitCards(string(e))
		this.Round++
		this.Raises = 0
		switch this.Round {
		case PreFlop:
			this.Actor = this.firstPlayer[this.Round] - 1
			this.Holes = cards
		case Flop, Turn, River:
			this.Actor = this.firstPlayer[this.Round] - 1
			this.Board = append(this.Board, cards...)
		case Showdown:
			this.Actor = -1
			this.Holes = cards
		}
	case diff.Action:
		action := string(e)
		this.Actions[this.Actor] = action[0]
		switch action {
		case "c":
			this.Bets[this.Round][this.Actor] += this.CallAmt()
		case "r":
			this.Raises++
			this.Bets[this.Round][this.Actor] += this.RaiseAmt()
		}
		if this.NumActive() < 2  || this.evenBets() {
			this.Actor = -1
			return
		}
		i := this.Actor
		for {
			i = (i + 1) % len(this.Actions)
			if this.Actions[i] != 'f' {
				this.Actor = i
				return
			}
		}
	default:
		panic("game: Invalid event passed to Update")
	}
}

// Start playing a game.
//	rules -- a String naming the game to play
//	p     -- an object that implements the Player interface.
//	host  -- the InetAddress of the dealer passed as a String
//	port  -- the port the dealer is listening on for the client passed as a String
func Play(rules string, p Player, host, port string) {
	addr := net.JoinHostPort(host, port)
	fmt.Printf("Connecting to dealer at %s to play %s...\n", addr, rules)
	in, out, err := diff.NewACPC(addr)
	if err != nil {
		log.Fatalln(err)
	}
	game, err := NewGame(rules)
	if err != nil {
		log.Fatalln(err)
	}
	for event := range in {
		game.Update(event)
		if game.Actor == game.Viewer {
			out <- p.Play(game)
		} else {
			p.Observe(game)
		}
	}
	fmt.Println("GAME OVER")
}
