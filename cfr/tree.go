package cfr

import (
	"fmt"
	"bytes"
	"poker/game"
	"poker/game/diff"
)

var (
	itoa = [3]rune{'f', 'c', 'r'}
	atoi = map[rune]int{'f': 0, 'c': 1, 'r': 2}
)

// A Bucket node nodes represents where information about the cards is observed.
// It has a child node (an opponent or player node) for each different class
// that could be observed at that point.
type Bucket struct {
	Classes [1]interface{}
}

// An Opponent nodes represents where the opponent takes an action. It has a
// child node for each action.
type Opponent struct {
	Actions [3]interface{}
}

//  A Player node represents where the current player takes an action. It
//  contains the average regret with respect to each action, the total
//  probability for each action until this point, and a child node for each
//  action (either an Opponent, Bucket, or Terminal node). There is an implicit
//  information set associated with this node, which we will write as I(n).
type Player struct {
	Actions [3]interface{}
	Regret  [3]float64 // Accumulative regret wrt each action.
	Strat   [3]float64 // The probability of taking each action.
}

// A Terminal node is a node where the game ends due to someone folding or a
// showdown. Given the probability of a win, loss, and tie, it has sufficient
// information to compute an expected utility for the node that was reached.
//
// In other words, the Terminal node holds the size of the pot when the game
// ends.
type Terminal float64


func newNode(g *game.Game) interface{} {
	if g.Round == 4 || g.NumActive() < 2 {
		return Terminal(g.Pot())
	}
	if g.Actor == -1 {
		return new(Bucket)
	}
	if g.Actor == g.Viewer {
		return &Player{}
	}
	return new(Opponent)
}

// addNodes recursively adds nodes to a game tree until all possible
// playouts have been added.
func addNodes(node interface{}, g *game.Game) {
	switch n := node.(type) {
	case *Bucket:
		g.Update(diff.Cards(""))
		for i := range n.Classes {
			n.Classes[i] = newNode(g)
			addNodes(n.Classes[i], g)
		}
	case *Player:
		la := g.LegalActions()
		s := 1 / float64(len(la))
		for _, a := range la {
			i := atoi[a]
			n.Strat[i] = s
			g1 := g.Copy()
			g1.Update(diff.Action(itoa[i]))
			n.Actions[i] = newNode(g1)
			addNodes(n.Actions[i], g1)
		}
	case *Opponent:
		for _, a := range g.LegalActions() {
			i := atoi[a]
			g1 := g.Copy()
			g1.Update(diff.Action(itoa[i]))
			n.Actions[i] = newNode(g1)
			addNodes(n.Actions[i], g1)
		}
	case Terminal:
	default:
		panic(fmt.Sprintln("Invalid type of node encountered.", n))
	}
}

// bfWalk returns a string representing a breadth-first walk of a game tree.
func bfWalk(root interface{}) string {
	var depth int
	buf := new(bytes.Buffer)
	q := make([]interface{}, 1)
	q[0] = root
	for len(q) > 0 {
		depth++
		qt := make([]interface{}, len(q))
		copy(qt, q)
		q = q[:0]
		for _, node := range qt {
			buf.WriteString("(")
			switch n := node.(type) {
			case *Bucket:
				for i := range n.Classes {
					buf.WriteString(fmt.Sprintf("B%d", i))
					q = append(q, n.Classes[i])
				}
			case *Player:
				buf.WriteString(fmt.Sprintf("%v", n.Strat))
				for i, action := range n.Actions {
					if action != nil {
						buf.WriteRune(itoa[i])
						q = append(q, n.Actions[i])
					}
				}
			case *Opponent:
				for i, action := range n.Actions {
					if action != nil {
						buf.WriteRune(itoa[i])
						q = append(q, n.Actions[i])
					}
				}
			case Terminal:
				buf.WriteString(fmt.Sprintf("%.0f", float64(n)))
			default:
				panic(fmt.Sprintln("Invalid type of node encountered.", n))
			}
			buf.WriteString(")")
		}
		buf.WriteString("\n")
	}
	buf.WriteString(fmt.Sprintf("Depth: %d\n", depth))
	return buf.String()
}
