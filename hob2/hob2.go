package main

import (
	"bytes"
	"flag"
	"fmt"
	"math/rand"
	"log"
	"os"
	"runtime/pprof"

	"poker/equity"
	"poker/game"
)


// The randPlayer chooses an action uniform randomly from all legal actions each
// turn.
type randPlayer struct {
	Name string
}

func (_ *randPlayer) Play(g *game.Game) byte {
	a := g.LegalActions()
	return a[rand.Intn(len(a))]
}

// stratPlayer chooses the action which has the greatest EV based on the 7cHS,
// pot odds, and an implied call. An an opponent call is included in the EV,
// because in a two player game if the oppenent does not call, the player will
// outright win the pot.
type stratPlayer struct {
	Name   string
	equity float64
}

func (p *stratPlayer) Play(g *game.Game) byte {
	if len(g.RoundActions()) == 0 {
		p.equity = equity.HandEquity(g.Cards.Holes[:2], g.Cards.Board, 1000)
	}

	max := 0.0 // Folding has EV = 0
	action := byte('f')
	if ev := (p.equity * (g.Pot + g.Call)) - g.Call; ev >= max {
		action = 'c'
		max = ev
	}
	if (p.equity*(g.Pot+2*g.Raise-g.Call))-g.Raise >= max {
		action = 'r'
	}
	if bytes.IndexByte(g.LegalActions(), action) != -1 {
		return action
	}
	return 'c'
}

func chooseStrat(name, strat string) (game.Player, error) {
	switch strat {
	case "random":
		return &randPlayer{name}, nil
	case "7cHS":
		return &stratPlayer{Name: name}, nil
	}
	return nil, fmt.Errorf("The strategy %s was not found.", strat)
}

func main() {
	prof := flag.Bool("prof", false, "Create a pprof profile.")
	rules := flag.String("rules", "2p-l", "What rules to use.")
	strat := flag.String("strat", "7cHS", "What strategy to use.")
	flag.Parse()
	if *prof {
		f, err := os.Create("hob2.prof")
		if err != nil {
			log.Fatalln("Failed to create profile:", err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	player, err := chooseStrat("Hob2", *strat)
	if err != nil {
		log.Fatalln("Failed to create player:", err)
	}
	game.Play(*rules, player, flag.Arg(0), flag.Arg(1))
}
