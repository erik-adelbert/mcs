// Copyright 2018 Erik Adelbert. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chaingame

import (
	"fmt"
	"math/rand"
)

type (
	// A Tile is a game piece, it's a list of blocks.
	Tile []Block

	// Tiles is a set of tiles.
	Tiles []Tile
)

func (t Tile) String() string {
	if t == nil {
		return "(){}"

	}

	return fmt.Sprintf("%v{%d}", t[0], len(t))
}

// ColorTiles is set of tiles grouped by color.
type ColorTiles map[Color]Tiles

// Colors lists all the colors present in a set of tiles.
func (c ColorTiles) Colors() []Color {
	colors := make([]Color, 0, len(c))
	for color := range c {
		if color != NoColor {
			colors = append(colors, color)
		}
	}
	return colors
}

// Histogram computes counts of blocks by color.
func (c ColorTiles) Histogram() Histogram {
	colors := make(Histogram, int(AllColors))
	for color := range c {
		tlen := 0
		for _, tile := range c.Tiles(color) {
			tlen += len(tile)
		}
		colors[color] = float64(tlen)
	}
	return colors
}

// Len is the number of tiles in a color group. It gives the total number
// of tiles when called with 'AllColors'.
func (c ColorTiles) Len(color Color) int {

	if color == AllColors {
		tlen := 0
		for _, tiles := range c {
			tlen += len(tiles)
		}
		return tlen
	}

	return len(c[color])
}

// PickTile removes a random tile from the set when called with 'NoColor'
// or from the corresponding subset when a color is specified.
func (c ColorTiles) PickTile(taboo Color) Tile {

	if c == nil {
		return nil
	}

	var color Color
	for color = range c { // ColorTiles are randomized by construction
		if color != taboo || taboo == NoColor {
			break
		}
	}

	var tile Tile
	if tiles := c[color]; len(tiles) > 0 { // Tiles are randomized by construction
		tile, tiles[0] = tiles[0], nil

		if len(tiles) == 1 {
			delete(c, color)
		} else {
			c[color] = tiles[1:]
		}
	}

	return tile
}

// RandomTile chooses, if possible,  a tile at random which isn't of taboo color.
// Returns a random color tile if called with 'NoColor'
func (c ColorTiles) RandomTile(taboo Color) Tile {

	var color Color
	for color = range c { // Colors are randomized by construction
		if taboo != color || taboo == NoColor {
			break
		}
	}
	tiles := c[color]

	return tiles[rand.Intn(len(tiles))]
}

// Tiles lists tiles of a color group. It lists all tiles
// when called with 'AllColors'
func (c ColorTiles) Tiles(color Color) Tiles {

	if color != AllColors {
		return c[color]
	}

	all := make(Tiles, 0, c.Len(AllColors))
	for _, tiles := range c {
		all = append(all, tiles...)
	}
	return all

}

// TaggedTiles are used to extract examples tiles (connected components)
type TaggedTiles map[Tag]Tile

// Colormap groups TaggedTiles by color
func (t TaggedTiles) Colormap(b Board) ColorTiles {

	if b.Len() == 0 {
		return nil
	}

	t.build(b)

	m := make(ColorTiles, int(AllColors))
	for tag, tile := range t { // Randomized by construction
		if len(tile) > 1 {
			color := tag.Color()

			tiles, ok := m[color]
			if !ok {
				tiles = make(Tiles, 0, 4)
			}
			m[color] = append(tiles, tile)
		}
		delete(t, tag)
	}

	return m
}

// List returns all the tiles.
func (t TaggedTiles) List(b Board) Tiles {

	if b.Len() == 0 {
		return nil
	}

	t.build(b)

	tiles := make(Tiles, 0, len(t))
	for tag, tile := range t { // Randomized by construction
		if len(tile) > 1 {
			tiles = append(tiles, tile)
		}
		delete(t, tag)
	}
	return tiles
}

// Board tiling is done efficiently in a Hoshen–Kopelman manner:
// 1) The board is copied into a larger buffer that supports labeling
//    while simplifying corners and borders handling.
// 2) a) Tags are flooded from north and west and b) conflicting ones are unified into sets.
// 3) Tiles of tag sets are merged.
// see:
// https://en.wikipedia.org/wiki/Hoshen%E2%80%93Kopelman_algorithm
func (t TaggedTiles) build(b Board) {
	if b.Len() == 0 {
		return
	}

	// 1 - Copy
	h, w := b.Dims()
	array := make([][]Tag, h+2) // add cells all around the board: no more overflows!
	cells := make([]Tag, (h+2)*(w+2))
	for i := range array {
		array[i], cells = cells[:w+2], cells[w+2:]
	}

	for i, row := range b {
		for j, block := range row {
			array[i+1][j+1] = Tag(block)
		}
	}

	t.extract(array)
}

func (t TaggedTiles) extract(array [][]Tag) {
	var h, w int
	if h = len(array); h > 0 {
		w = len(array[0])
	}

	// 2a - Flooding
	// Finding tiles (ie. tagging disjoint sets).
	tags := NewTags(1 + h*w/8)
	for i := 1; i < h-1; i++ {
		for j := 1; j < w-1; j++ {
			color := array[i][j].Color()

			if color == NoColor {
				continue
			}

			// Flooding tags from north and west
			north, west := array[i-1][j], array[i][j-1]
			switch {
			case north.Color() == color && west.Color() == color:
				array[i][j] = tags.Union(north, west) // 2b - Two tags: unify tag sets!
			case north.Color() == color:
				array[i][j] = north
			case west.Color() == color:
				array[i][j] = west
			default:
				array[i][j] = tags.NewID(color)
			}

			// Append current block to its labeled tile
			var tile Tile
			if tile = t[array[i][j]]; len(tile) == 0 {
				tile = make(Tile, 0, 2)
			}
			t[array[i][j]] = append(tile, Block{i - 1, j - 1})
		}
	}

	// 3 - Merge tiles according to previous unifications.
	for tag := range tags.List() {

		root := tags.Find(tag)
		if root == tag {
			continue
		}

		tb := t[tag]
		delete(t, tag)

		ta := t[root]
		if len(ta) < len(tb) { //  Appends the smallest tile to the bigger one.
			ta, tb = tb, ta
		}
		t[root], ta, tb = append(ta, tb...), nil, nil
	}
}
