// Copyright 2018 Erik Adelbert. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package samegame

import (
	"strings"

	"mcs/pkg/game"
)

// SameBoard maintains an histogram alongside a  board. It's a convenience type which only
// purpose is performance. By incrementally computing histograms of boards in inner loops,
// it provides more time for examples solving.
type SameBoard struct {
	b game.Board
	h game.Histogram
}

// Returns an initialized 0-value SameBoard. Storage is allocated but
// board is empty.
func NewSameBoard(h, w int) SameBoard {
	board := game.NewBoard(h, w)
	histo := make(game.Histogram, int(game.AllColors))

	return SameBoard{board, histo}
}

func (cb SameBoard) Cap() int {
	return cb.b.Cap()
}

func (cb SameBoard) Caps() (int, int) {
	return cb.b.Caps()
}

// Returns a full copy with its own storage.
func (sb SameBoard) Clone() SameBoard {
	histo := make(game.Histogram, int(game.AllColors))
	for k, v := range sb.h {
		histo[k] = v
	}

	board := sb.b.Clone()

	return SameBoard{board, histo}
}

// Builds and lists all tiles grouped by color.
func (sb SameBoard) ColorTiles() game.ColorTiles {
	return sb.b.ColorTiles()
}

// A SameBoard's size is the underlying board's size
func (sb SameBoard) Dims() (int, int) {
	return sb.b.Dims()
}

// By construction, the histogram is returned in Θ(1).
// This is a desired behaviour for this object.
func (sb SameBoard) Histogram() game.Histogram {
	return sb.h
}

// A SameBoard's len is the underlying board's len
func (sb SameBoard) Len() int {
	return sb.b.Len()
}

// Loads a SameBoard from strings and initializes its histogram.
func (sb SameBoard) Load(buf []string) {
	sb.b.Load(buf)

	for k, v := range sb.b.Histogram() {
		sb.h[k] = v
	}
}

// Randomizes the SameBoard  with n colors :
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

// Removes a tile from the board. The resulting histogram is built in Θ(m)
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

// Returns the color of a given tile.
func (sb SameBoard) TileColor(t game.Tile) game.Color {
	// The color of a tile is the color of its first block
	if len(t) == 0 {
		return game.NoColor
	}

	block := t[0]
	return sb.b[block.Row()][block.Column()]
}

// Builds and lists all tiles.
func (sb SameBoard) Tiles() game.Tiles {
	return sb.b.Tiles()
}
