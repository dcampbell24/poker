// Package diff provides the Diff interface for controlling a Game.
package diff

type Diff interface {
	Play(string) error // Play an Action.
	Position() int     // Position of viewer relative to dealer button.
	Round()    int     // 0-4: pre-flop, flop, turn, river, showdown.
	Action()   string  // f, c, r.
	Cards()    string  // New Cards.
	// Update blocks waiting for the next game update. If the game has ended, a
	// *GameOver struct will be returned that contains any results.
	Update()   error
}

type GameOver struct {
	Results string
}

func (this *GameOver) Error() string {
	return "GameOver"
}
