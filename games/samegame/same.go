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
	b game.Board
	h game.Histogram
}

// NewSameBoard returns an initialized 0-value ClickBoard.
// Storage is allocated but board is empty.
func NewSameBoard(h, w int) SameBoard {
	board := game.NewBoard(h, w)
	histo := make(game.Histogram, int(game.AllColors))

	return SameBoard{board, histo}
}

// Cap returns the maximum number of blocks a board can have.
func (sb SameBoard) Cap() int {
	return sb.b.Cap()
}

// Caps returns the maximum lines and columns a board can have.
func (sb SameBoard) Caps() (int, int) {
	return sb.b.Caps()
}

// Clone returns an independent copy of a board.
func (sb SameBoard) Clone() SameBoard {
	histo := make(game.Histogram, int(game.AllColors))
	for k, v := range sb.h {
		histo[k] = v
	}

	board := sb.b.Clone()

	return SameBoard{board, histo}
}

// ColorTiles extract and lists all tiles of the calling board grouped by color.
func (sb SameBoard) ColorTiles() game.ColorTiles {
	return sb.b.ColorTiles()
}

// Dims returns the current size of a board.
func (sb SameBoard) Dims() (int, int) {
	return sb.b.Dims()
}

// Histogram is returned in Θ(1). This is a desired behaviour for this object.
func (sb SameBoard) Histogram() game.Histogram {
	return sb.h
}

// Len is the underlying board's len
func (sb SameBoard) Len() int {
	return sb.b.Len()
}

// Load a SameBoard from strings and initializes its histogram.
func (sb SameBoard) Load(buf []string) {
	sb.b.Load(buf)

	for k, v := range sb.b.Histogram() {
		sb.h[k] = v
	}
}

// Randomize the SameBoard  with n colors :
//  - repeating a color modifies the distribution accordingly.
//  - passing 'AllColors' adds one of every color to the list.
//  - a one color list produces a board filled with a unique tile.
//
// c.Randomize(AllColors, Red, Indigo)
func (sb SameBoard) Randomize(list ...game.Color) {
	sb.b.Randomize(list...)

	for k, v := range sb.b.Histogram() {
		sb.h[k] = v
	}
}

// Remove a tile from the board. The resulting histogram is built in Θ(m)
// with m < 7, the number of colors initially present on the board.
// This is a desired behaviour for this object.
func (sb SameBoard) Remove(t game.Tile) SameBoard {
	color := sb.TileColor(t)

	n := sb.h[color]

	if n = n - float64(len(t)); n == 0 {
		delete(sb.h, color)
	} else {
		sb.h[color] = n
	}

	sb.b = sb.b.Remove(t)

	return sb
}

// The stringer pretty prints an ansi-colorized board, the memory address of
// the SameBoard and the associated histogram.
func (sb SameBoard) String() string {
	var s strings.Builder

	s.WriteString(sb.b.String())
	s.WriteByte('\n')
	s.WriteString(sb.h.String())

	return s.String()
}

// TileColor returns the color of a given tile.
func (sb SameBoard) TileColor(t game.Tile) game.Color {
	// The color of a tile is the color of its first block
	if len(t) == 0 {
		return game.NoColor
	}

	block := t[0]
	return sb.b[block.Row()][block.Column()]
}

// Tiles extract and lists all tiles.
func (sb SameBoard) Tiles() game.Tiles {
	return sb.b.Tiles()
}
