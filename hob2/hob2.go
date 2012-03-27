package main

import (
	"bytes"
	"math/rand"
	"os"
	"runtime/pprof"

	"poker/equity"
	"poker/game"
)

const CPU_PROF = true

type randPlayer struct {
	Name string
}

type stratPlayer struct {
	Name   string
	equity float64
}

func (_ *randPlayer) Play(g *game.Game) byte {
	a := g.LegalActions()
	return a[rand.Intn(len(a))]
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

func main() {
	if CPU_PROF {
		f, err := os.Create("hob2.prof")
		if err != nil {
			panic(err)
		}

		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	game.Play("Holdem2p", &stratPlayer{Name: "Hob2"}, os.Args[1], os.Args[2])
}
