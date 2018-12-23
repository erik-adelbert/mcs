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
	b game.Board
	h game.Histogram
}

// NewClickBoard returns an initialized 0-value ClickBoard.
// Storage is allocated but board is empty.
func NewClickBoard(h, w int) ClickBoard {
	board := game.NewBoard(h, w)
	histo := make(game.Histogram, int(game.AllColors))

	return ClickBoard{board, histo}
}

// Cap returns the maximum number of blocks a board can have.
func (cb ClickBoard) Cap() int {
	return cb.b.Cap()
}

// Caps returns the maximum lines and columns a board can have.
func (cb ClickBoard) Caps() (int, int) {
	return cb.b.Caps()
}

// Clone returns an independent copy of a board.
func (cb ClickBoard) Clone() ClickBoard {
	histo := make(game.Histogram, int(game.AllColors))
	for k, v := range cb.h {
		histo[k] = v
	}

	board := cb.b.Clone()

	return ClickBoard{board, histo}
}

// ColorTiles extract and lists all tiles of the calling board grouped by color.
func (cb ClickBoard) ColorTiles() game.ColorTiles {
	return cb.b.ColorTiles()
}

// Dims returns the current size of a board.
func (cb ClickBoard) Dims() (int, int) {
	return cb.b.Dims()
}

// Histogram is returned in Θ(1). This is a desired behaviour for this object.
func (cb ClickBoard) Histogram() game.Histogram {
	return cb.h
}

// Len is the underlying board's len
func (cb ClickBoard) Len() int {
	return cb.b.Len()
}

// Load a ClickBoard from strings and initializes its histogram.
func (cb ClickBoard) Load(buf []string) {
	cb.b.Load(buf)

	for k, v := range cb.b.Histogram() {
		cb.h[k] = v
	}
}

// Randomize the ClickBoard  with n colors :
//  - repeating a color modifies the distribution accordingly.
//  - passing 'AllColors' adds one of every color to the list.
//  - a one color list produces a board filled with a unique tile.
//
// c.Randomize(AllColors, Red, Indigo)
func (cb ClickBoard) Randomize(list ...game.Color) {
	cb.b.Randomize(list...)

	for k, v := range cb.b.Histogram() {
		cb.h[k] = v
	}
}

// Remove a tile from the board. The resulting histogram is built in Θ(m)
// with m < 7, the number of colors initially present on the board. This is
// a desired behaviour for this object.
func (cb ClickBoard) Remove(t game.Tile) ClickBoard {
	color := cb.TileColor(t)

	n := cb.h[color]

	if n = n - float64(len(t)); n == 0 {
		delete(cb.h, color)
	} else {
		cb.h[color] = n
	}

	cb.b = cb.b.Remove(t)

	return cb
}

// The stringer pretty prints an ansi-colorized board, the memory address of
// the ClickBoard and the associated histogram.
func (cb ClickBoard) String() string {
	var s strings.Builder

	s.WriteString(cb.b.String())
	s.WriteByte('\n')
	s.WriteString(cb.h.String())

	return s.String()
}

// TileColor returns the color of a given tile.
func (cb ClickBoard) TileColor(t game.Tile) game.Color {
	// The color of a tile is the color of its first block
	if len(t) == 0 {
		return game.NoColor
	}

	block := t[0]
	return cb.b[block.Row()][block.Column()]
}

// Tiles extract and lists all tiles.
func (cb ClickBoard) Tiles() game.Tiles {
	return cb.b.Tiles()
}
