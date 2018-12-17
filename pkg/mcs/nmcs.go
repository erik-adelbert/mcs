// Copyright 2018 Erik Adelbert. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mcs

import (
	"fmt"
	"time"
)

const (
	maxLevel = 4
)

// cazenave & jm
// see:
// http://www.lamsade.dauphine.fr/~cazenave/papers/nested.pdf
// https://www.researchgate.net/publication/48445151_Combining_UCT_and_Nested_Monte-Carlo_Search_for_Single-Player_General_Game_Playing
func recursiveSearch(root *Node, policies []GamePolicy, duration time.Duration) Decision {

	if root == nil {
		// TODO: error handling
		panic("no root")
	}

	done := make(chan struct{})
	timeout := make(chan bool)

	go func() {
		time.Sleep(duration)
		close(timeout)
	}()

	conclude := func(root *Node) Decision {
		var final Decision

		// Tracking
		node := root
		for len(node.down) > 0 {
			node = node.down[0]

			final.score += node.edge.Score()
			final.Moves().Enqueue(node.edge)
		}

		if moves := node.Best().Moves(); moves.Len() > 0 { // Leaf has been sampled
			final.score += node.Best().Score()
			final.Moves().Join(node.Best().Moves())
		} else { // Leaf is terminal
			final.score += node.State().Score()
		}

		return final
	}

	var search func(root *Node, level int) Decision
	search = func(root *Node, level int) Decision {
		node := root

		var best Move
		var top Decision

		_ = best

		if level == 0 { // best playout with memoization
			for i := 0; i < 10; i++ {
				for _, move := range node.Hand().List() {
					start := node.State().Clone().Play(move)
					simulated := start.Sample(done, policies[0], 1.0)

					if simulated.Score() > top.Score() {
						best = move
						top = simulated
					}
				}
			}
		} else {

		}

		if top.Score() > node.Best().Score() {
			node.best = top
		}

		return conclude(root)
	}

	search(root, 0)

	fmt.Println(root.best)

	return conclude(root)
}
