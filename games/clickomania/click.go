// Copyright 2018 Erik Adelbert. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package clickgame

import (
	"strings"

	"mcs/pkg/chaingame"
)

// ClickBoard maintains an histogram alongside a  board. It's a convenience type which only
// purpose is performance. By incrementally computing histograms of boards in inner loops,
// it provides more time for problem solving.
type ClickBoard struct {
	chaingame.Board
	chaingame.Histogram
}

// NewClickBoard returns an initialized 0-value ClickBoard.
// Storage is allocated but board is empty.
func NewClickBoard(h, w int) ClickBoard {
	board := chaingame.NewBoard(h, w)
	histo := make(chaingame.Histogram, int(chaingame.AllColors))

	return ClickBoard{board, histo}
}

// Clone returns an independent copy of a board.
func (cb ClickBoard) Clone() ClickBoard {
	histo := make(chaingame.Histogram, int(chaingame.AllColors))
	for k, v := range cb.Histogram {
		histo[k] = v
	}

	board := cb.Board.Clone()

	return ClickBoard{board, histo}
}

// Load a ClickBoard from strings and initializes its histogram.
func (cb ClickBoard) Load(buf []string) {
	cb.Board.Load(buf)

	for k, v := range cb.Board.Histogram() {
		cb.Histogram[k] = v
	}
}

// Randomize the ClickBoard  with n colors :
//  - repeating a color modifies the distribution accordingly.
//  - passing 'AllColors' adds one of every color to the list.
//  - a one color list produces a board filled with a unique tile.
//
// c.Randomize(AllColors, Red, Indigo)
func (cb ClickBoard) Randomize(list ...chaingame.Color) {
	cb.Board.Randomize(list...)

	for k, v := range cb.Board.Histogram() {
		cb.Histogram[k] = v
	}
}

// Remove a tile from the board. The resulting histogram is built in Θ(m)
// with m < 7, the number of colors initially present on the board. This is
// a desired behaviour for this object.
func (cb ClickBoard) Remove(t chaingame.Tile) ClickBoard {
	color := cb.TileColor(t)

	n := cb.Histogram[color]

	if n = n - float64(len(t)); n == 0 {
		delete(cb.Histogram, color)
	} else {
		cb.Histogram[color] = n
	}

	cb.Board = cb.Board.Remove(t)

	return cb
}

// The stringer pretty prints an ansi-colorized board, the memory address of
// the ClickBoard and the associated histogram.
func (cb ClickBoard) String() string {
	var s strings.Builder

	s.WriteString(cb.Board.String())
	s.WriteByte('\n')
	s.WriteString(cb.Histogram.String())

	return s.String()
}

// TileColor returns the color of a given tile.
func (cb ClickBoard) TileColor(t chaingame.Tile) chaingame.Color {
	// The color of a tile is the color of its first block
	if len(t) == 0 {
		return chaingame.NoColor
	}

	block := t[0]
	return cb.Board[block.Row()][block.Column()]
}
