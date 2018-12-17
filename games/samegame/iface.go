// Copyright 2018 Erik Adelbert. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package samegame

import (
	"mcs/pkg/game"
)

type Hand game.ColorTiles

func (h Hand) Draw() (Move, Hand) {
	tiles := game.ColorTiles(h)
	tile := tiles.PickTile(game.NoColor)

	return Move(tile), Hand(tiles)
}

func (h Hand) Len() int {
	return game.ColorTiles(h).Len(game.AllColors)
}

func (h Hand) List() []Move {
	moves := make([]Move, 0, h.Len())
	tiles := game.ColorTiles(h).Tiles(game.AllColors)

	for _, tile := range tiles {
		moves = append(moves, Move(tile))
	}

	return moves
}

type Move game.Tile

func (m Move) Len() int {
	return len(game.Tile(m))
}

func (m Move) Score() float64 {
	if m.Len() == 0 {
		return 0
	}

	n := float64(m.Len())
	return (n - 2) * (n - 2)
}

func (m Move) String() string {
	return game.Tile(m).String()
}

type Sequence []Move

func (moves Sequence) Enqueue(m Move) Sequence {
	if m.Len() == 0 {
		return moves
	}
	return append(moves, m)
}

func (moves Sequence) Clone() Sequence {
	clone := make(Sequence, len(moves))
	copy(clone, moves)
	return clone
}

func (moves Sequence) Join(seq Sequence) Sequence {
	return append(moves, seq...)
}

func (moves Sequence) Len() int {
	return len(moves)
}

func (moves Sequence) Dequeue() (Move, Sequence) {
	move, moves := moves[0], moves[1:]
	return move, moves
}

type State SameBoard

func (sg State) Clone() State {
	return State(SameBoard(sg).Clone())
}

func (sg State) Moves() Hand {
	return Hand(SameBoard(sg).ColorTiles())
}

func (sg State) Play(m Move) State {
	return State(SameBoard(sg).Remove(game.Tile(m)))
}

func (sg State) Sample(done <-chan struct{}, policy ColorPolicy, solved float64) (float64, Sequence, float64) {

	board := SameBoard(sg)
	tiles := board.ColorTiles()

	taboo := game.NoColor
	if c, mode := policy(board); mode == PerPlayout {
		taboo = c
	}

	var seq Sequence
	var score float64

	if len(tiles) > 1 {
		solved = 0
	}

	for len(tiles) > 0 {
		select {
		case <-done:
			return score, seq, solved
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

	return score, seq, solved
}

func (sg State) Score() float64 {
	penalty, bonus := 0.0, 0.0
	for _, n := range SameBoard(sg).h {
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
