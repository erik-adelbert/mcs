// Copyright 2018 Erik Adelbert. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package samegame

import (
	"mcs/pkg/game"
)

// Mode indicates how the policy is to be used.
type Mode int

const (
	// PerSampling indicates that the policy has to be called once
	// at the start of the simulation.
	PerSampling Mode = iota

	// PerMove indicates that the policy has to be called for each
	// move during the simulation.
	PerMove
)

// ColorPolicy is used for selecting nodes and moves during playouts.
type ColorPolicy func(SameBoard) (game.Color, Mode)

// NoTaboo deactivate taboo selection.
func NoTaboo(board SameBoard) (game.Color, Mode) {
	_ = board
	return game.NoColor, PerSampling
}

// TabooColor returns a color ot to be played unless it is the only available move.
func TabooColor(board SameBoard) (game.Color, Mode) {

	taboo, max := game.NoColor, 0.0
	for c, n := range board.Histogram() {
		if n > max {
			taboo, max = c, n
		}
	}

	return taboo, PerSampling
}
