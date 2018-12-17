// Copyright 2018 Erik Adelbert. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mcs

import (
	"log"
	"time"
)

const slot = 10 * time.Minute

// MetaSearch splits allowed thinking time into time slots. A new search is launched
// for each time slot (cycle) and the best result is returned.
func MetaSearch(root *Node, policies []GamePolicy, duration time.Duration) Decision {
	var best Decision
	var cycle time.Duration

	switch {
	case duration > 3*slot: // from 31mn to 50mn: 2 runs
		cycle = duration / 2
	case duration > 5*slot: // 51mn and up: evenly divide runtime among 4 runs
		cycle = duration / 4
	default: // up to 10mn : 1 run, from 11mn to 30m : 10mn runs
		cycle = slot
	}

	if first := duration % cycle; first != 0 {
		best = ConcurrentSearch(root, policies, first)
		log.Printf("[meta] first (%v) : %g\n", cycle.Seconds(), best.Score())
	}

	cycles := duration / cycle
	clone := CloneRoot(root)
	for cycles > 0 {
		best = ConcurrentSearch(clone, policies, cycle)
		clone = CloneRoot(clone)
		log.Printf("[meta] cycle #%d (%v) : %g\n", cycles, cycle, best.Score())

		cycles--
	}

	return best
}
