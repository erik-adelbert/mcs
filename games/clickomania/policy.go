// Copyright 2018 Erik Adelbert. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package clickgame

import "mcs/pkg/chaingame"

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
type ColorPolicy func(ClickBoard) (chaingame.Color, Mode)

// NoTaboo deactivate taboo selection.
func NoTaboo(board ClickBoard) (chaingame.Color, Mode) {
	_ = board
	return chaingame.NoColor, PerSampling
}

// TabooColor returns a color ot to be played unless it is the only available move.
func TabooColor(board ClickBoard) (chaingame.Color, Mode) {

	taboo, max := chaingame.NoColor, 0.0
	for c, n := range board.Histogram() {
		if n > max {
			taboo, max = c, n
		}
	}

	return taboo, PerSampling
}
