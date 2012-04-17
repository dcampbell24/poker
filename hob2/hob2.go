package main

import (
	"flag"
	"fmt"
	"math/rand"
	"log"
	"os"
	"runtime/pprof"
	"strings"

	"poker/equity"
	"poker/game"
)


// The randPlayer chooses an action uniform randomly from all legal actions each
// turn.
type randPlayer struct {
	Name string
}

func (_ *randPlayer) Observe(_ *game.Game) {}

func (_ *randPlayer) Play(g *game.Game) string {
	a := g.LegalActions()
	return string(a[rand.Intn(len(a))])
}


// stratPlayer chooses the action which has the greatest EV based on the 7cHS,
// pot odds, and an implied call. An an opponent call is included in the EV,
// because in a two player game if the oppenent does not call, the player will
// outright win the pot.
type stratPlayer struct {
	Name   string
	equity float64
}

func (this *stratPlayer) Observe(g *game.Game) {
	if (g.Cards() != "") && (g.Round() != 4) {
		this.equity = equity.HandEquity(g.Holes, g.Board, 1000)
	}
}

func (this *stratPlayer) Play(g *game.Game) string {
	if g.Cards() != "" {
		this.equity = equity.HandEquity(g.Holes, g.Board, 1000)
	}

	max := 0.0 // Folding has EV = 0
	action := "f"
	c := g.CallAmt()
	r := g.RaiseAmt()
	pot := g.Pot()
	if ev := (this.equity * (pot + c)) - c; ev >= max {
		action = "c"
		max = ev
	}
	if (this.equity*(pot + 2*r - c)) - r >= max {
		action = "r"
	}
	if strings.Contains(g.LegalActions(), action) {
		return action
	}
	return "c"
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
