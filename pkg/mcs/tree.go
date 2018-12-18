// Copyright 2018 Erik Adelbert. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file implements Monte-Carlo trees for single player games.
// More information are available in [2008 Schadd et al.]
// see http://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.167.3355&rep=rep1&type=pdf

package mcs

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
)

var nodeCounter struct {
	Spinlock
	value int
}

func nodeCountInc() {
	nodeCounter.Lock()
	{
		nodeCounter.value++
	}
	nodeCounter.Unlock()
}

// NodeCount returns the grand total of created nodes.
func NodeCount() int {
	var count int
	nodeCounter.Lock()
	{
		count = nodeCounter.value
	}
	nodeCounter.Unlock()
	return count
}

// NodeRate returns the nodes creation rate on a 10s window.
var NodeRate = nodometer()

func nodometer() func() float64 {
	lock := NewSpinlock()

	last := 0
	t0 := time.Now()
	ticker := time.NewTicker(10 * time.Second)

	go func() {
		for t := range ticker.C {
			lock.Lock()
			{
				t0 = t
				last = NodeCount()
			}
			lock.Unlock()
		}
	}()

	return func() float64 {
		var count, elapsed float64

		current := NodeCount()

		lock.Lock()
		{
			count = float64(current - last)
			elapsed = float64(time.Since(t0).Seconds())
		}
		lock.Unlock()

		return count / elapsed
	}
}

// NodeStatus reflects the cycle of nodes during Monte-Carlo Searches.
// A node is either up to date (Idle), or went through selection/expansion step (Walked),
// has been sent to simulation (Sampling) or out of simulation (Sampled). The status is
// eventually reset to Idle during the update step.
type NodeStatus int

const (
	idle NodeStatus = iota
	walked
	sampling
	sampled
)

func (ns NodeStatus) String() string {
	switch ns {
	case idle:
		return "idle"
	case walked:
		return "walked"
	case sampled:
		return "sampled"
	case sampling:
		return "sampling"
	default:
		return "unknown"
	}
}

// Node composes a multi-branching tree asymptotically akin to mini-max tree.
type Node struct {
	*Spinlock

	edge   Move
	depth  int
	status NodeStatus

	up   *Node
	down []*Node

	hand  MoveSet
	state GameState

	best Decision

	solved float64

	value float64

	mean     float64
	visits   float64
	variance float64

	ε float64
	c float64
	w float64
}

// NewNode allocates a Monte-Carlo tree node.
func NewNode(up *Node, edge Move, state GameState, hand MoveSet, ε, c, w float64) *Node {
	nodeCountInc()

	depth := 0
	if up != nil {
		depth = up.depth + 1
	}

	var node = Node{
		edge:   edge,
		up:     up,
		depth:  depth,
		status: idle,

		state: state,
		hand:  hand,

		down: make([]*Node, 0, 24), // average branching factor is 20.7

		ε: ε,
		c: c,
		w: w,
	}

	node.Spinlock = NewSpinlock()

	return &node
}

// NewRoot initializes a node with an initial position and various constants used
// during the search.
func NewRoot(initial GameState, ε, c, w float64) *Node {
	return NewNode(nil, nil, initial, initial.Moves(), ε, c, w)
}

// NewTree expands a root node in order to bootstrap a search.
func NewTree(root *Node) *Node {

	if root == nil || root.hand.Len() == 0 {
		// TODO: error  handling
		panic("no moves")
	}

	root.ExpandAll(math.Inf(1))

	return root
}

// CloneRoot returns a memory independent copy of the calling node.
func CloneRoot(root *Node) *Node {
	initial, ε, c, w := root.State().Clone(), root.ε, root.c, root.w

	clone := NewRoot(initial, ε, c, w)
	clone.best = root.Best().Clone()

	return clone
}

// Best returns the best sequence found so far.
func (n *Node) Best() Decision {

	n.Lock()
	defer n.Unlock()
	{
		return n.best
	}
}

// Depth returns the distance from the calling node to the root node.
func (n *Node) Depth() int {
	n.Lock()
	defer n.Unlock()
	{
		return n.depth
	}
}

// Down returns a slice containing references to the children
// of the calling node.
func (n *Node) Down() []*Node {
	return n.down
}

// Downselect chooses next edge using linear ε-greedy algorithm:
// It chooses a random node with probability 1-ε/n
// see https://arxiv.org/pdf/1402.6028.pdf
func (n *Node) Downselect() *Node {

	p := 1.0
	if v := n.Visits(); v > 0 {
		//p = n.ε / v // ε-greedy linear decay
		p = n.ε
	}

	var node *Node

	n.Lock()
	{
		if rand.Float64() > p { // ε-greedy
			By(Value).SortDescending(n.down)
		} else {
			swap := func(i, j int) { n.down[i], n.down[j] = n.down[j], n.down[i] }
			rand.Shuffle(len(n.down), swap)
		}

		for _, node = range n.down {
			if node.GetLock() {
				status := node.status
				//solved := node.IsSolvedUnsafe()
				node.Unlock()
				//if status == idle && !solved {
				if status == idle {
					break
				}
			}
		}

		// next time, if possible, expand another one :
		// status will eventually be reset by updater
		n.value = math.Inf(-1)
		// n.status = walked
	}
	n.Unlock()

	return node
}

// Edge returns the move which leads to the calling node.
func (n *Node) Edge() Move {

	if n == nil {
		return nil
	}

	n.Lock()
	defer n.Unlock()
	{
		return n.edge
	}
}

// Evaluate set the UCB value of the calling node.
func (n *Node) Evaluate() float64 {
	value := n.UCB()
	n.Lock()
	{
		n.value = value
	}
	n.Unlock()

	return value
}

// ExpandOne creates and links a new children to the calling node.
func (n *Node) ExpandOne(move Move) *Node {

	state := n.State().Clone().Play(move)
	moves := state.Moves()

	node := NewNode(n, move, state, moves, n.ε, n.c, n.w)

	n.Lock()
	{
		n.down = append(n.down, node)
	}
	n.Unlock()

	return node
}

// ExpandAll creates and links children for every legal move from
// the position of the calling node.
func (n *Node) ExpandAll(value float64) {

	for _, move := range n.Hand().List() {
		node := n.ExpandOne(move)
		node.SetValue(value)
	}

	n.Lock()
	{
		n.hand = MoveSet{}
	}
	n.Unlock()

}

// Hand lists all the legal moves for the calling node.
func (n *Node) Hand() MoveSet {

	n.Lock()
	defer n.Unlock()
	{
		return n.hand
	}
}

// IsExpanded states if whether or not all the legal moves
// of the calling node have been expanded in the tree.
func (n *Node) IsExpanded() bool {

	n.Lock()
	defer n.Unlock()
	{
		return n.hand.Len() == 0 && len(n.down) > 0
	}
}

// IsLeaf states if the calling node has no child.
func (n *Node) IsLeaf() bool {

	n.Lock()
	defer n.Unlock()
	{
		return len(n.down) == 0 && n.hand.Len() > 0
	}
}

// IsSolved isn't used.
func (n *Node) IsSolved() bool {
	n.Lock()
	defer n.Unlock()
	{
		return 1 == n.solved/float64(n.hand.Len()+len(n.down))
	}
}

// IsSolvedUnsafe isn't used.
func (n *Node) IsSolvedUnsafe() bool {
	return 1 == n.solved/float64(n.hand.Len()+len(n.down))

}

// IsTerminal is true if the calling node is the second to
// last move of a game.
func (n *Node) IsTerminal() bool {

	n.Lock()
	defer n.Unlock()
	{
		return len(n.down) == 0 && n.hand.Len() == 1
	}
}

// Mean is the running mean score of the calling node.
func (n *Node) Mean() float64 {

	n.Lock()
	defer n.Unlock()
	{
		return n.mean
	}
}

// RandomNewEdge removes and return a move from the calling node's hand.
func (n *Node) RandomNewEdge() (move Move) {

	n.Lock()
	{
		move, n.hand = n.hand.Draw()
	}
	n.Unlock()

	return
}

// SampleVariance is guaranteed to be numerically stable.
func (n *Node) SampleVariance() float64 {
	n.Lock()
	defer n.Unlock()
	{
		return n.variance / (n.visits - 1)
	}
}

// SetStatus is a safe setter.
func (n *Node) SetStatus(status NodeStatus) {
	n.Lock()
	{
		n.status = status
	}
	n.Unlock()
}

// SetValue is a safe setter.
func (n *Node) SetValue(value float64) {
	n.Lock()
	{
		n.value = value
	}
	n.Unlock()
}

// State returns the position associated to the calling node.
func (n *Node) State() GameState {

	return n.state
}

// Status is a safe getter.
func (n *Node) Status() NodeStatus {

	n.Lock()
	defer n.Unlock()
	{
		return n.status
	}
}

// StDev is a safe running variance.
func (n *Node) StDev() float64 {

	return math.Sqrt(n.Variance())
}

func (n *Node) String() string {
	if n == nil {
		return "nil node"
	}

	var sb strings.Builder
	n.Lock()
	{
		if n.edge != nil {
			sb.WriteString("edge: " + n.edge.String() + "\n")
		} else {
			sb.WriteString("edge: nil\n")
		}

		if _, err := fmt.Fprintf(&sb, "depth: %d\n", n.depth); err != nil {
			panic(err)
		}

		sb.WriteString("status: " + n.status.String() + "\n")

		if _, err := fmt.Fprintf(&sb, "solved : %g\n", n.solved); err != nil {
			panic(err)
		}

		sb.WriteByte('\n')

		if _, err := fmt.Fprintf(&sb, "up: %p\n", n.up); err != nil {
			panic(err)
		}

		if _, err := fmt.Fprintf(&sb, "down: (%d)\n", len(n.down)); err != nil {
			panic(err)
		}

		for i := 0; i < len(n.down); i++ {
			if _, err := fmt.Fprintf(&sb, "%2d: %v\t@%p\t", i, n.down[i].edge, n.down[i]); err != nil {
				panic(err)
			}

			if (i+1)%3 == 0 {
				sb.WriteByte('\n')
			}
		}

		sb.WriteByte('\n')

		if _, err := fmt.Fprint(&sb, "\nhand: ", n.hand, "\n"); err != nil {
			panic(err)
		}

		sb.WriteString("\nstate:\n" + n.state.String() + "\n")

		sb.WriteString("\nbest:\n" + n.best.String())

		sb.WriteByte('\n')

		if _, err := fmt.Fprintf(&sb, "value = %g\n", n.value); err != nil {
			panic(err)
		}

		if _, err := fmt.Fprintf(&sb, "mean = %g, visits = %g, variance = %g\n", n.mean, n.visits, n.variance); err != nil {
			panic(err)
		}

		if _, err := fmt.Fprintf(&sb, "ε = %g, c = %g, w = %g", n.ε, n.c, n.w); err != nil {
			panic(err)
		}
	}
	n.Unlock()
	return sb.String()
}

// UCB calls the package wide enabled UCB formula.
func (n *Node) UCB() float64 {

	return SelectedUCB(n)
}

// Up enables tree navigation toward tree's root.
func (n *Node) Up() *Node {
	return n.up
}

// UpdateTree guides the search: it back propagates outcomes from simulations
// in the tree path leading to the calling node. It records the best sequence
// running through this node that has been found so far. It maintains running
// mean and variance with a numerically stable technique. Finally it computes
// UCB values enabling next search iteration to select the most promising node.
func (n *Node) UpdateTree(decision Decision) {

	//if n != nil && !n.IsSolved() {
	if n != nil {
		n.Lock()
		{
			/*
				if decision.solved != 0 {
					if n.solved = n.solved + decision.solved; !n.IsSolvedUnsafe() {
						decision.solved = 0
					}
				}
			*/

			n.visits++

			score := decision.score

			// Running mean and variance from B. P. Welford.
			// This variance computation is numerically stable.
			// see:
			// D.E. Knuth TAOCP Vol 2, page 232, 3rd edition.
			old := n.mean
			cur := old + (score-old)/n.visits

			n.variance = n.variance + (score-old)*(score-cur)
			n.mean = cur

			if score > n.best.Score() {
				n.best = decision
			}
		}
		n.Unlock()

		n.up.UpdateTree(decision)

		n.Evaluate()
	}
}

// Value returns the calling node's search score.
func (n *Node) Value() float64 {

	n.Lock()
	defer n.Unlock()
	{
		return n.value
	}
}

// Variance maintains a running variance for the calling node.
func (n *Node) Variance() float64 {

	n.Lock()
	defer n.Unlock()
	{
		return n.variance / n.visits
	}
}

// Visits returns the number of simulations that run through the calling node.
func (n *Node) Visits() float64 {
	if n == nil {
		return 1.0
	}

	n.Lock()
	defer n.Unlock()
	{
		return n.visits
	}
}
