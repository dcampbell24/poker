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

func NewACPC(addr string) (chan interface{}, chan string, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, nil, err
	}
	// Tell the dealer we are ready to start playing.
	conn.Write([]byte(ACPCversion))
	var gstring string
	in := make(chan string)
	go func() {
		for {
			a := <-in
			fmt.Fprintf(conn, "%s:%s\r\n", gstring, a)
		}
	}()
	ch := make(chan interface{}, 3)
	go func() {
		var offsets [5]int
		var state []string
		var handNum string
		var position, round int
		bufin := bufio.NewReader(conn)
		defer conn.Close()
		defer close(ch)
		for {
			line, err := bufin.ReadString('\n')
			if err == io.EOF {
				fmt.Println("Received EOF, shutting down.")
				return
			} else if err != nil {
				fmt.Println("Error during game update:", err)
				return
			// ";" and "#" are comment lines.
			} else if len(line) < 1 || line[0] == ';' || line[0] == '#' {
				continue
			}
			// Only update the gstring if there game has advanced.
			gstring = strings.TrimRight(line, "\r\n")
			state = strings.Split(gstring, ":")
			// New hand.
			if handNum != state[_HAND_NUM] {
				handNum = state[_HAND_NUM]
				round = 0
				position, err = strconv.Atoi(state[_POSITION])
				if err != nil {
					fmt.Printf("Recieved invalid position value %s\n", state[_POSITION])
					return
				}
				ch <- &Players{Viewer: position}
				// New hole cards.
				ch <- Cards(strings.Trim(state[_CARDS], "/|"))
				for i := range offsets {
					offsets[i] = 0
				}
			} else {
				// New action.
				ch <- Action(state[_BETS][offsets[_BETS]])
				// New round.
				if len(state[_CARDS][offsets[_CARDS]:]) > 0 {
					round++
					// New board cards.
					if round != 4 {
						ch <- Cards(state[_CARDS][offsets[_CARDS]:])
					// Hole cards revealed.
					} else {
						ch <- Cards(strings.SplitN(state[_CARDS], "/", 2)[0])
					}
				}
			}
			// Update offsets.
			for i := range state {
				offsets[i] = len(state[i])
			}
		}
	}()
	return ch, in, nil
}
