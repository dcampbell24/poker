// Package diff provides the Diff engine for updating a game.
//
// The job of a diff engine is to set up and maintain the connection with the
// dealer if one is needed, to send updates from the dealer, and to recieve
// Actions from players. When a game has ended, the engine will close the channel
// on which updates are sent. All updates should be an Action, Cards, or
// Players.
package diff

type Action  string   // f, c, r.
type Cards   string   // AsKd
// The names of all the players in the current hand ordered by their
// position relative to the dealer button. If the players' names are not
// known, then Names will be nil.
type Players struct {
	Names  []string // The names of all the players.
	Viewer int      // The offset into Names of the viewer.
}
