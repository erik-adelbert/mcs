// Copyright 2018 Erik Adelbert. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mcs

import "mcs/games/samegame"

// (╯°□°）╯︵ ┻━┻ poor mans's generic:

// GameState can be anything that describes accurately the state of a game.
// In samegame it's a board.
type GameState samegame.State

// Clone returns a memory-independent copy.
func (g GameState) Clone() GameState {
	return GameState(samegame.State(g).Clone())
}

// Moves returns a list of legal moves from the calling state.
func (g GameState) Moves() MoveSet {
	return MoveSet(samegame.State(g).Moves())
}

// Play returns the game state after the given move has been played in the
// calling state.
func (g GameState) Play(m Move) GameState {
	return GameState(samegame.State(g).Play(m.(samegame.Move)))
}

// Sample simulates a game to its end by applying a move selection policy. The policy usually
// embeds randomness.
func (g GameState) Sample(done <-chan struct{}, policy GamePolicy) Decision {
	score, moves := samegame.State(g).Sample(done, samegame.ColorPolicy(policy))

	return Decision{moves: MoveSequence(moves), score: score}
}

// Score returns a statically computed score of the calling state.
func (g GameState) Score() float64 {
	return samegame.State(g).Score()
}

func (g GameState) String() string {
	return samegame.State(g).String()
}

// Move is scorable, has a length and is printable.
type Move interface {
	Score() float64
	Len() int
	String() string
}

// MoveSequence is a FIFO structure.
type MoveSequence samegame.Sequence

// Clone returns an independent copy of the calling sequence.
func (s MoveSequence) Clone() MoveSequence {
	return MoveSequence(samegame.Sequence(s).Clone())
}

// Dequeue is customary for FIFO structures.
func (s MoveSequence) Dequeue() (Move, MoveSequence) {
	move, seq := samegame.Sequence(s).Dequeue()
	return Move(move), MoveSequence(seq)
}

// Enqueue is customary for FIFO structures.
func (s MoveSequence) Enqueue(m Move) MoveSequence {
	return MoveSequence(samegame.Sequence(s).Enqueue(m.(samegame.Move)))
}

// Join returns an aggregated sequence.
func (s MoveSequence) Join(t MoveSequence) MoveSequence {
	return MoveSequence(samegame.Sequence(s).Join(samegame.Sequence(t)))
}

// Len returns the number of moves in the sequence.
func (s MoveSequence) Len() int {
	return samegame.Sequence(s).Len()
}

// MoveSet is a collection of legal moves.
type MoveSet samegame.Hand

// Draw randomly removes a move from the set.
func (m MoveSet) Draw() (Move, MoveSet) {
	move, set := samegame.Hand(m).Draw()
	return Move(move), MoveSet(set)
}

// Len returns the number of legal moves.
func (m MoveSet) Len() int {
	return samegame.Hand(m).Len()
}

// List returns all the moves present in the set.
func (m MoveSet) List() []Move {
	list := samegame.Hand(m).List()

	moves := make([]Move, 0, len(list))
	for _, move := range list {
		moves = append(moves, Move(move))
	}
	return moves
}

// GamePolicy is a game policy used during the simulation step.
// It is a reference passed back to the game sampler.
type GamePolicy samegame.ColorPolicy

// (╯°□°）╯︵ ┻━┻ here is precisely what interfaces aren't for:
/*
type GameState interface {
	Clone() GameState
	Moves() MoveSet
	Play(Edge) GameState
	Sample(<-chan struct{}, float64, GamePolicy) (float64, MoveSequence)
	Score() float64
	String() string
}

type Edge interface {
	Score() float64
}

type MoveSequence interface {
	Enqueue(Edge) MoveSequence
	Join(MoveSequence) MoveSequence
	Len() int
}

type MoveSet interface {
	Draw() (Edge, MoveSet)
	Len() int
	List() []Edge
}

type GamePolicy interface{}
*/
