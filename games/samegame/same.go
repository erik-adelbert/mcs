// Copyright 2018 Erik Adelbert. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package samegame

import (
	"strings"

	"mcs/pkg/chaingame"
)

// SameBoard maintains an histogram alongside a  board. It's a convenience type which only
// purpose is performance. By incrementally computing histograms of boards in inner loops,
// it provides more time for problem solving.
type SameBoard struct {
	chaingame.Board
	chaingame.Histogram
}

// NewSameBoard returns an initialized 0-value ClickBoard.
// Storage is allocated but board is empty.
func NewSameBoard(h, w int) SameBoard {
	board := chaingame.NewBoard(h, w)
	histo := make(chaingame.Histogram, int(chaingame.AllColors))

	return SameBoard{board, histo}
}

// Clone returns an independent copy of a board.
func (sb SameBoard) Clone() SameBoard {
	histo := make(chaingame.Histogram, int(chaingame.AllColors))
	for k, v := range sb.Histogram {
		histo[k] = v
	}

	board := sb.Board.Clone()

	return SameBoard{board, histo}
}

// Load a SameBoard from strings and initializes its histogram.
func (sb SameBoard) Load(buf []string) {
	sb.Board.Load(buf)

	for k, v := range sb.Board.Histogram() {
		sb.Histogram[k] = v
	}
}

// Randomize the SameBoard  with n colors :
//  - repeating a color modifies the distribution accordingly.
//  - passing 'AllColors' adds one of every color to the list.
//  - a one color list produces a board filled with a unique tile.
//
// c.Randomize(AllColors, Red, Indigo)
func (sb SameBoard) Randomize(list ...chaingame.Color) {
	sb.Board.Randomize(list...)

	for k, v := range sb.Board.Histogram() {
		sb.Histogram[k] = v
	}
}

// Remove a tile from the board. The resulting histogram is built in Θ(m)
// with m < 7, the number of colors initially present on the board.
// This is a desired behaviour for this object.
func (sb SameBoard) Remove(t chaingame.Tile) SameBoard {
	color := sb.TileColor(t)

	n := sb.Histogram[color]

	if n = n - float64(len(t)); n == 0 {
		delete(sb.Histogram, color)
	} else {
		sb.Histogram[color] = n
	}

	sb.Board = sb.Board.Remove(t)

	return sb
}

// The stringer pretty prints an ansi-colorized board, the memory address of
// the SameBoard and the associated histogram.
func (sb SameBoard) String() string {
	var s strings.Builder

	s.WriteString(sb.Board.String())
	s.WriteByte('\n')
	s.WriteString(sb.Histogram.String())

	return s.String()
}

// TileColor returns the color of a given tile.
func (sb SameBoard) TileColor(t chaingame.Tile) chaingame.Color {
	// The color of a tile is the color of its first block
	if len(t) == 0 {
		return chaingame.NoColor
	}

	block := t[0]
	return sb.Board[block.Row()][block.Column()]
}
