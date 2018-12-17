// Copyright 2018 Erik Adelbert. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file defines the Decision type and associated methods.

package mcs

import (
	"fmt"
	"strings"
)

// Decision is made of a sequence of moves from the starting position that yields
// the recorded score. Decisions are formed during selection/expansion and finalized
// during games simulations.
type Decision struct {
	moves  MoveSequence
	score  float64
	solved float64 // not used yet
}

// Clone returns an independent copy of a decision.
func (d Decision) Clone() Decision {
	var clone Decision

	clone.moves = make(MoveSequence, d.Moves().Len())
	copy(clone.moves, d.Moves())
	clone.score = d.Score()
	//clone.solved = d.solved

	return clone
}

// Join merges two decisions.
func (d Decision) Join(other Decision) Decision {
	d.moves = d.moves.Join(other.Moves())
	d.score += other.Score()
	//d.solved = d.solved + other.solved
	return d
}

// Moves is a getter.
func (d Decision) Moves() MoveSequence {
	return d.moves
}

// Score is a getter.
func (d Decision) Score() float64 {
	return d.score
}

// Solved is not yet functional.
func (d Decision) Solved() float64 {
	return d.solved
}

// SetMoves is a setter.
func (d Decision) SetMoves(m MoveSequence) {
	d.moves = m
}

// SetScore is a setter.
func (d Decision) SetScore(score float64) {
	d.score = score
}

// SetSolved isn't functional yet.
func (d Decision) SetSolved(solved float64) {
	d.solved = solved
}

func (d Decision) String() string {
	var sb strings.Builder

	if _, err := fmt.Fprintf(&sb, "score: %g, solved: %v\n", d.score, d.solved); err != nil {
		panic(err)
	}

	for i := 0; i < len(d.moves); i++ {
		if _, err := fmt.Fprintf(&sb, "%2d: %s\t", i, d.moves[i].String()); err != nil {
			panic(err)
		}

		if (i+1)%6 == 0 {
			sb.WriteByte('\n')
		}
	}
	sb.WriteByte('\n')

	return sb.String()
}
