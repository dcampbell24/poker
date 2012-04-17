package diff

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

const ACPCversion = "VERSION:2.0.0\r\n"

// Fields of the ACPC state string.
const (
	_ = iota
	_POSITION
	_HAND_NUM
	_BETS
	_CARDS
)

type ACPC struct {
	position  int      // Position of viewer relative to dealer button.
	round     int      // 0-4: pre-flop, flop, turn, river, showdown.
	action    string   // f, c, r.
	cards     string   // New cards.
	handNum   string   // Unique identifier for each hand.
	gstring   string   // The latest game string from the dealer.
	offsets   [5]int   // Current offsets into the game state string.
	conn      net.Conn // Connection with the dealer server.
	bufin     *bufio.Reader // Buffer to read updates from.
}

func NewACPC(addr string) (*ACPC, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	// Tell the dealer we are ready to start playing.
	conn.Write([]byte(ACPCversion))
	return &ACPC{conn: conn, bufin: bufio.NewReader(conn)}, nil
}

func (this *ACPC) Play(action string) error {
	_, err := fmt.Fprintf(this.conn, "%s:%s\r\n", this.gstring, action)
	return err
}

func (this *ACPC) Position() int {
	return this.position
}

func (this *ACPC) Round() int {
	return this.round
}

func (this *ACPC) Action() string {
	return this.action
}

func (this *ACPC) Cards() string {
	return this.cards
}

func (this *ACPC) Update() error {
	line, err := this.bufin.ReadString('\n')
	if err == io.EOF {
		this.conn.Close() //FIXME Should we close at any other times as well?
		return new(GameOver)
	} else if err != nil {
		return err
	// ";" and "#" are comment lines.
	} else if len(line) < 1 || line[0] == ';' || line[0] == '#' {
		return this.Update()
	}
	// Only update the gstring if there game has advanced.
	this.gstring = strings.TrimRight(line, "\r\n")
	state := strings.Split(this.gstring, ":")
	// New hand.
	if this.handNum != state[_HAND_NUM] {
		this.handNum = state[_HAND_NUM]
		this.round = 0
		this.position, err = strconv.Atoi(state[_POSITION])
		if err != nil {
			return fmt.Errorf("Recieved invalid position value %s\n", state[_POSITION])
		}
		// New hole cards.
		this.cards = strings.Trim(state[_CARDS], "/|")
		this.action = ""
		for i := range this.offsets {
			this.offsets[i] = 0
		}
	} else {
		// New action.
		this.action = strings.TrimRight(state[_BETS][this.offsets[_BETS]:], "/")
		// New round.
		if len(state[_CARDS][this.offsets[_CARDS]:]) > 0 {
			this.round++
			// New board cards.
			if this.round != 4 {
				this.cards = state[_CARDS][this.offsets[_CARDS]:]
			// Hole cards revealed.
			} else {
				this.cards = strings.SplitN(state[_CARDS], "/", 2)[0]
			}
		} else {
			this.cards = ""
		}
	}
	// Update offsets.
	for i := range state {
		this.offsets[i] = len(state[i])
	}
	return nil
}
