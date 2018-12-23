// Copyright 2018 Erik Adelbert. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file implements a classical UCT

package mcs

import (
	"time"
)

// ConfidentSearch implements a classical UCT as specified in [2006 Kocsis, Szepesvári]
// see http://ggp.stanford.edu/readings/uct.pdf
func ConfidentSearch(root *Node, policies []GamePolicy, duration time.Duration) Decision {

	if root == nil {
		// TODO: error handling
		panic("no tree")
	}

	tree := GrowTree(root)

	done := make(chan struct{})
	timeout := make(chan bool)

	go func() {
		time.Sleep(duration)
		close(timeout)
	}()

	for {
		select {

		case <-timeout:
			close(done)
			return tree.Best()

		default:
			if root.IsSolved() {
				close(done)
				return tree.Best()
			}

			const VisitThreshold = 8

			var score float64
			var moves MoveSequence

			node := tree

			for node.IsExpanded() {
				node = node.Downselect()

				move := node.Edge()
				moves = moves.Enqueue(move)

				score += move.Score()
			}

			if !node.IsTerminal() && node.Visits() > VisitThreshold {
				node = node.ExpandOne(node.RandomNewEdge())

				move := node.Edge()
				moves = moves.Enqueue(move)

				score += move.Score()
			}

			clone := node.State().Clone()
			sampled := clone.Sample(done, policies[0])

			sampled.moves = moves.Join(sampled.moves)
			sampled.score += score

			node.UpdateTree(sampled)
		}
	}
}

//func NestedConfidentSearch(root *Node, policies []GamePolicy, duration time.Duration) Decision {}
