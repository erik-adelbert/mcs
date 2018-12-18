// Copyright 2018 Erik Adelbert. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package game

import (
	"fmt"
	"math/rand"
	"strings"
)

// Board is a row ordered 2D matrix of colored blocks.
// It's a directly addressable regular slice of slices.
type Board [][]Color

// NewBoard allocates storage for a h*w board with row first structure.
func NewBoard(h, w int) Board {

	if h*w == 0 {
		return nil
	}

	board := make([][]Color, h)
	cells := make([]Color, h*w)

	for i := range board {
		board[i], cells = cells[:w], cells[w:]
	}

	h, w = Board(board).Caps()

	return board
}

// Caps returns rows and columns capacities
func (b Board) Caps() (h, w int) {
	d, h, w := len(b), cap(b), 0

	if h == 0 {
		return
	}

	if d == 0 {
		b = b[:1]
	}

	w = cap(b[0]) / h

	return
}

// Cap returns the absolute capacity: a 10x15 board has the same cap as
// a 15x10 board which is 150.
func (b Board) Cap() int {
	h, w := b.Caps()
	return h * w
}

// Clone produces an independent copy.
func (b Board) Clone() Board {
	r, c := b.Caps()

	board := NewBoard(r, c)

	for i, row := range b {
		board = board[:i+1]
		copy(board[i], row)
	}

	return board
}

// ColorTiles lists tiles grouped by colors.
func (b Board) ColorTiles() ColorTiles {
	tiles := make(TaggedTiles, b.Len())
	return tiles.Colormap(b)
}

// Dims returns board's height and width
func (b Board) Dims() (h, w int) {
	h, w = len(b), 0
	if h > 0 {
		w = len(b[0])
	}
	return
}

// Histogram counts board blocks grouped by color.
func (b Board) Histogram() Histogram {
	if b.Len() == 0 {
		return nil
	}

	m := make(Histogram, int(AllColors))
	for i, row := range b {
		for j := range row {
			m[b[i][j]] = m[b[i][j]] + 1.0
		}
	}
	delete(m, NoColor) // Remove empty cells count

	return m
}

// Len is the board number of blocks.
func (b Board) Len() int {
	h, w := b.Dims()
	return h * w
}

// Load a board from a slice of strings.
func (b Board) Load(buf []string) {

	for i, row := range buf {
		for j, c := range row {
			b[i][j] = NewColor(c)
		}
	}

	return
}

// Randomize a board with n colors :
//  - repeating a color modifies the distribution accordingly.
//  - listing 'AllColors' adds one of every color to the list.
//  - a one color list produces a board filled with a unique tile.
func (b Board) Randomize(list ...Color) {
	//rand.Seed(time.Now().UnixNano())
	colors := make([]Color, 0, AllColors)

	for _, color := range list {
		switch color {
		case NoColor:
			// Do nothing
		case AllColors:
			colors = append(colors, Red, Green, Yellow, Blue, Violet, Indigo, Orange)
		default:
			colors = append(colors, color)
		}
	}

	for i, row := range b {
		for j := range row {
			b[i][j] = colors[rand.Intn(len(colors))]
		}
	}
}

// Tiles returns a list of all game pieces.
func (b Board) Tiles() Tiles {
	tiles := make(TaggedTiles, b.Len())
	return tiles.List(b)
}

// shrink is an helper: it empties a board without releasing its storage:
// If it's filled again, the storage is reused.
func shrink(b Board) Board {

	if b != nil {
		for i := range b {
			b[i] = b[i][:0]
		}
		b = b[:0]
	}

	return b
}

// Remove a tile from the board.
// The resulting board is built in three phases:
//   1) board blocks composing the tile are marked deleted;
//   2) a transposed matrix is built with no deleted block nor shrink column;
//	 3) this matrix is transposed back with no shrink row.
// This emulates perfectly the physics of clickomania.
func (b Board) Remove(t Tile) Board {
	h, w := b.Dims()

	if h*w == 0 {
		return nil
	}

	// 1 - Mark tile's blocks as deleted
	for _, block := range t {
		b[block.r][block.c] = NoColor
	}

	// 2 - Deleted blocks and shrink columns are removed while transposing.
	trans := shrink(NewBoard(w, h))
	for last, j := 0, 0; j < w; j++ {
		trans = trans[:last+1] // ExpandOne by one row

		empty := true
		for i := 0; i < h; i++ { // Dump deleted blocks
			if b[i][j] != NoColor {
				trans[last] = append(trans[last], b[i][j])
				empty = false
			}
		}

		if empty { // Last column was shrink: reuse storage!
			trans[last] = trans[last][:0] // Flush row
			trans = trans[:last]          // Retract row
			continue
		}

		// If needed, columns are padded.
		for len(trans[last]) < h {
			trans[last] = append([]Color{NoColor}, trans[last]...)
		}

		last++
	}

	// 3 - Transpose back while deleting shrink rows
	w, h = trans.Dims()

	b = shrink(b) // Reuse already allocated storage
	for last, i := 0, 0; i < h; i++ {
		b = b[:last+1] // ExpandOne by one row

		empty := true
		for j := 0; j < w; j++ {
			empty = empty && trans[j][i] == NoColor
			b[last] = append(b[last], trans[j][i])
		}

		if empty { // Last row was shrink: reuse storage!
			b[last] = b[last][:0]
			b = b[:last]
			continue
		}

		last++
	}

	return b
}

// String also prints the board's memory address.
func (b Board) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Board@%p\n", b))

	if b.Len() == 0 {
		sb.WriteString("Empty\n")
	}

	for i, row := range b {
		if _, err := fmt.Fprintf(&sb, "%3d: ", i); err != nil {
			panic("board: can't string")
		}
		for _, block := range row {
			sb.WriteString(block.String())
		}
		sb.WriteByte('\n')
	}

	return strings.TrimSuffix(sb.String(), "\n")
}
