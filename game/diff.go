package game

import (
	"fmt"
	"strconv"
	"strings"
)

// Fields of the ACPC state string.
const (
	_ = iota
	_POSITION
	_HAND_NUM
	_BETS
	_CARDS
)

type GameDiff struct {
	Position  int    // Position of viewer relative to dealer button.
	Round     int    // 0-4: pre-flop, flop, turn, river, showdown.
	Action    string // f, c, r.
	Cards     string // New Cards.
	handNum   string // Unique identifier for each hand.
	offsets   [5]int // Current offsets into the game state string.
}

func (this *GameDiff) Update(s string) error {
	var err error
	state := strings.Split(s, ":")
	// New hand.
	if this.handNum != state[_HAND_NUM] {
		this.handNum = state[_HAND_NUM]
		this.Round = 0
		this.Position, err = strconv.Atoi(state[_POSITION])
		if err != nil {
			return fmt.Errorf("Recieved invalid position value %s\n", state[_POSITION])
		}
		// New hole cards.
		this.Cards = strings.Trim(state[_CARDS], "/|")
		this.Action = ""
		for i := range this.offsets {
			this.offsets[i] = 0
		}
	} else {
		// New Action.
		this.Action = strings.TrimRight(state[_BETS][this.offsets[_BETS]:], "/")
		// New round.
		if len(state[_CARDS][this.offsets[_CARDS]:]) > 0 {
			this.Round++
			// New board cards.
			if this.Round != 4 {
				this.Cards = state[_CARDS][this.offsets[_CARDS]:]
			// Hole cards revealed.
			} else {
				this.Cards = strings.SplitN(state[_CARDS], "/", 2)[0]
			}
		} else {
			this.Cards = ""
		}
	}
	// Update offsets.
	for i := range state {
		this.offsets[i] = len(state[i])
	}
	return nil
}
