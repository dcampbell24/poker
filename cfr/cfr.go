// Package cfr implements functions and structures for calculating epsilon Nash
// equilibrium of an abstracted poker game using counterfactual regret
// minimization.
//
// References:
//
//	Martin Zinkevich, Michael Johanson, Michael Bowling, and Carmelo Piccione.
//	Regret mini-mization in games with incomplete information. In NIPS07, 2008.
//
//	Annotation: Technical proofs of minimization bounds and horrible pseudo code.
//
//	M.B. Johanson, Robust strategies and counter-strategies: Building a champion
//	level computer poker player, Master’s thesis, University of Alberta, 2007
//
//	Annotation: Good general overview and optimization details.
//
//	M. Osborne and A. Rubenstein. A Course in Game Theory. The MIT Press, Cambridge,
//	Massachusetts, 1994.
//
//	Annotation: Explanation of a finite, extensive game with imperfect information.
//
//	J. Rubin, I. Watson. Computer poker: A review. In Artificial Intelligence 175
//	(2011) 958–987.
//
//	Annotation: General overview and mention of several bucketing strategies.
package cfr

import (
	"poker/game"
	"poker/game/diff"
)

// CFR returns the counterfactual regret for each action in a node given the
// expected value for each action, the node's current strategy, and the
// probability of the opponent reaching the node.
func CFR(actEV, strat []float64, p float64) []float64 {
	cfr := make([]float64, len(actEV))
	nodeEV := 0.0
	for i := range actEV {
		nodeEV += actEV[i] * strat[i]
	}
	for i := range cfr {
		cfr[i] = p * (actEV[i] - nodeEV)
	}
	return cfr
}

// NewStrategy calculates a regret minimizing new strategy given a set of
// cumulative counterfactual regret.
func NewStrategy(cfr []float64) []float64 {
	psum := 0.0
	strat := make([]float64, len(cfr))
	for _, r := range cfr {
		if r > 0.0 {
			psum += r
		}
	}
	if psum > 0 {
		for i, r := range cfr {
			if r > 0.0 {
				strat[i] = r/psum
			}
		}
	}
	return strat
}

/*
func NewBucketSeq(size int) {
	bs := make([]float64, 3)
	for i := range bs {
		bs[i] := float64(rand.Intn(size)
	}
	return bs
}
*/

func CalcNash(trials int) {
	g1, err := game.NewGame("2p-l")
	if err != nil {
		panic(err)
	}
	g2 := g1.Copy()
	g1.Update(&diff.Players{Viewer: 0})
	g2.Update(&diff.Players{Viewer: 1})
	r1 := new(Bucket)
	r2 := new(Bucket)
	addNodes(r1, g1)
	addNodes(r2, g2)
	/*
	for i := 0; i < trials; i++ {
		bs := NewBucketSeq(1)
		walkTrees(r1, r2, bs, 1.0, 1.0)
	}
	*/
}
/*
func walkTree(node interface{}, p int, p1, p2 []float64) {
	switch n := node.(type) {
	case Terminal:
		return utility
	case *PublicBucket:
		a := sample(node.Classes)
		return WalkTree(node.Classes[a], p1, p2)
	case *PrivateBucket:
		return WalkTree(node, p1, p2)
	case *Player:




}
*/
/*
// Walk the two trees, r1 and r2,  updating their strategies using cfr.
// Why use two trees? Wouldn't one work just as well and take up less space?
func walkTrees(r1, r2 node, b bucketSeq, p1, p2 float64) {
	switch r1.(type) {
	case *Player:
		// compute player one's strategy sigma(I(r1)) // eq. 8
		strat := newStrategy(r1)
		for _, a := range A(I(r1)) {
			c1, c2 := r1.child(a), r2.child(a) // I think...?
			u1, u2 := walkTrees(c1, c2, b, p1 * ?, p2)
		}
		// ...

	case *Opponent:
	case *Bucket:
		// find c1 and c2 from b
		u1, u2 := walkTrees(c1, c2, b, p1, p2)
		// ...
	case Terminal:
		return 45.6 67.8 // utilities u1, u2
	}
}
*/
