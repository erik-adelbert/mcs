// Copyright 2018 Erik Adelbert. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package samegame

import (
	"mcs/pkg/chaingame"
)

// A Hand stores legal moves.
type Hand chaingame.ColorTiles

// Draw randomly removes a move from the hand.
func (h Hand) Draw() (Move, Hand) {
	tiles := chaingame.ColorTiles(h)
	tile := tiles.PickTile(chaingame.NoColor)

	return Move(tile), Hand(tiles)
}

// Len returns the number of available legal moves.
func (h Hand) Len() int {
	return chaingame.ColorTiles(h).Len(chaingame.AllColors)
}

// List returns a list containing all legal moves.
func (h Hand) List() []Move {
	moves := make([]Move, 0, h.Len())
	tiles := chaingame.ColorTiles(h).Tiles(chaingame.AllColors)

	for _, tile := range tiles {
		moves = append(moves, Move(tile))
	}

	return moves
}

// A Move is a tile that can be removed from a board.
type Move chaingame.Tile

// Len returns the number of blocks of the calling tile.
func (m Move) Len() int {
	return len(chaingame.Tile(m))
}

// Score computes the samegame score of a move.
func (m Move) Score() float64 {
	if m.Len() == 0 {
		return 0
	}

	n := float64(m.Len())
	return (n - 2) * (n - 2)
}

func (m Move) String() string {
	return chaingame.Tile(m).String()
}

// A Sequence of moves is a FIFO structure.
type Sequence []Move

// Enqueue adds a move in a sequence.
func (moves Sequence) Enqueue(m Move) Sequence {
	if m.Len() == 0 {
		return moves
	}
	return append(moves, m)
}

// Clone returns an independent copy of a sequence.
func (moves Sequence) Clone() Sequence {
	clone := make(Sequence, len(moves))
	copy(clone, moves)
	return clone
}

// Join aggregates two sequences.
func (moves Sequence) Join(seq Sequence) Sequence {
	return append(moves, seq...)
}

// Len is the number of moves in a sequence.
func (moves Sequence) Len() int {
	return len(moves)
}

// Dequeue returns the next move in a sequence.
func (moves Sequence) Dequeue() (Move, Sequence) {
	move := moves[0]
	return move, moves[1:]
}

// A State describes the board.
type State SameBoard

// Clone returns an independent copy.
func (sg State) Clone() State {
	return State(SameBoard(sg).Clone())
}

// Moves returns the legal moves.
func (sg State) Moves() Hand {
	return Hand(SameBoard(sg).ColorTiles())
}

// Play returns the state following a ply.
func (sg State) Play(m Move) State {
	return State(SameBoard(sg).Remove(chaingame.Tile(m)))
}

// Sample simulates a game to its end by applying a move selection policy. The policy usually
// embeds randomness.
func (sg State) Sample(done <-chan struct{}, policy ColorPolicy) (float64, Sequence) {

	board := SameBoard(sg)
	tiles := board.ColorTiles()

	taboo := chaingame.NoColor
	if c, mode := policy(board); mode == PerSampling {
		taboo = c
	}

	var seq Sequence
	var score float64

	for len(tiles) > 0 {
		select {
		case <-done:
			return score, seq
		default:
			if c, mode := policy(board); mode == PerMove {
				taboo = c
			}
			tile := tiles.RandomTile(taboo)

			board = board.Remove(tile)
			tiles = board.ColorTiles()

			move := Move(tile)
			seq = seq.Enqueue(move)
			score += move.Score()
		}
	}
	score += State(board).Score()

	return score, seq
}

// Score returns a statically computed score of the calling state.
func (sg State) Score() float64 {
	penalty, bonus := 0.0, 0.0
	for _, n := range SameBoard(sg).Histogram {
		penalty += n * n
	}
	if penalty == 0 {
		bonus = 1000
	}
	return bonus - penalty
}

func (sg State) String() string {
	return SameBoard(sg).String()
}
