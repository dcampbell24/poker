package diff

import (
	"fmt"
	"strings"
	"io/ioutil"
)

// The structure of record string is
// STATE:handNum:actions:cards:result:playerPositions
type record struct {
	actions []string // Actions split by round.
	holes   []string // Hole cards split by player.
	board   []string // Board cards split by round.
	players []string // Players indexed by position.
}

func newRecord(line string) *record {
	s := strings.Split(line, ":")
	cards := strings.SplitN(s[3], "/", 2)
	return &record{
		actions: strings.Split(s[2], "/"),
		holes:   strings.Split(cards[0], "|"),
		board:   strings.Split(cards[1], "/"),
		players: strings.Split(s[5], "|")}
}

type ACPCLog struct {
	hands    []string // All of the hands in the log file.
	action   string   // The most recent action.
	cards    string   // The most recently revealed cards.
	round    int      // The current round.
	// The players at the table ordered by their position relative to the
	// dealer button.
	players  []string
}

func (this *ACPCLog) String() string {
	return fmt.Sprintln(this.round, this.action, this.cards)
}

func NewACPCLog(file string) (*ACPCLog, error) {
	log, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return &ACPCLog{hands: strings.Split(string(log), "\n")}, nil
}

func (this *ACPCLog) Replay(c chan *ACPCLog) {
	for _, hand := range this.hands {
		if len(hand) > 4 && hand[:5] == "STATE" {
			record := newRecord(hand)
			this.cards = strings.Join(record.holes, " ")
			this.action = ""
			c <- this
			for this.round < 4 {
				for _, action := range record.actions[this.round] {
					this.action = string(action)
					this.cards = ""
					c <- this
				}
				this.round++
				if this.round < 4 {
					this.cards = record.board[this.round-1]
					this.action = ""
					c <- this
				}
			}
		}
	}
}
