package chaingame

import "testing"

func TestNewBoard(t *testing.T) {

}

func TestBoard_Capacities(t *testing.T) {

}

func TestBoard_Clone(t *testing.T) {

}

func TestBoard_ColorTiles(t *testing.T) {

}

func TestBoard_Dimensions(t *testing.T) {

}

func TestBoard_Histogram(t *testing.T) {

}

func TestBoard_Len(t *testing.T) {

}

func TestBoard_Load(t *testing.T) {

}

func TestBoard_Randomize(t *testing.T) {

}

func TestBoard_Remove(t *testing.T) {
	h, w := 20, 10
	board := NewBoard(h, w)
	board.Randomize(AllColors)
	t.Log(board)

	tiles := board.Tiles()
	t.Log(tiles)

	board = board.Remove(tiles[0])

	if board == nil {
		t.Errorf("Removing one tile should not empty the board")
	}

	t.Log("Removing", tiles[0])
	t.Log(board)
}

func TestBoard_RemoveAll(t *testing.T) {
	h, w := 20, 10
	board := NewBoard(h, w)
	board.Randomize(Red)
	t.Log(board)

	tiles := board.Tiles()
	t.Log(tiles)

	board = board.Remove(tiles[0])

	if board.Len() != 0 {
		t.Errorf("Removing all tiles should result in an empty board")
	}

	t.Log("Removing", tiles[0])
	t.Log(board)
}

func TestBoard_String(t *testing.T) {

}

func TestBoard_Tiles(t *testing.T) {

}
