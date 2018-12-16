/* policy.go
erik adelbert - 2018 - erik _ adelbert _ fr
*/
package clickgame

import "mcs/pkg/game"

type Mode int

const (
	PerPlayout Mode = iota
	PerMove
)

// ColorPolicy is used for selecting nodes and moves during playouts.
type ColorPolicy func(ClickBoard) (game.Color, Mode)

// NoTaboo deactivate taboo selection.
func NoTaboo(board ClickBoard) (game.Color, Mode) {
	_ = board
	return game.NoColor, PerPlayout
}

// A taboo color is not to be played unless it is the only available move.
func TabooColor(board ClickBoard) (game.Color, Mode) {

	taboo, max := game.NoColor, 0.0
	for c, n := range board.Histogram() {
		if n > max {
			taboo, max = c, n
		}
	}

	return taboo, PerPlayout
}
