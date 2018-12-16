/* iface.go
erik adelbert - 2018 - erik _ adelbert _ fr
*/

package clickgame

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
	return 0
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

type State ClickBoard

func (sg State) Clone() State {
	return State(ClickBoard(sg).Clone())
}

func (sg State) Moves() Hand {
	return Hand(ClickBoard(sg).ColorTiles())
}

func (sg State) Play(m Move) State {
	return State(ClickBoard(sg).Remove(game.Tile(m)))
}

func (sg State) Sample(done <-chan struct{}, p ColorPolicy) (float64, Sequence) {

	board := ClickBoard(sg)
	tiles := board.ColorTiles()

	taboo := game.NoColor
	if c, mode := p(board); mode == PerPlayout {
		taboo = c
	}

	var score float64
	var seq Sequence
	for len(tiles) > 0 {
		select {
		case <-done:
			return score, seq
		default:
			if c, mode := p(board); mode == PerMove {
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

func (sg State) Score() float64 {
	dim := float64(ClickBoard(sg).Cap())

	penalty := 0.0
	for _, n := range ClickBoard(sg).h {
		penalty += n
	}
	penalty = penalty / dim

	return 1 - penalty
}

func (sg State) String() string {
	return ClickBoard(sg).String()
}
