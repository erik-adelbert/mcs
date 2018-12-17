package game

import (
	"strings"
	"testing"
)

func TestNewColor(t *testing.T) {
	s := "RGYBVIO-"
	r := []Color{Red, Green, Yellow, Blue, Violet, Indigo, Orange, NoColor}

	colors := make([]Color, 0, int(AllColors))
	for _, c := range s {
		colors = append(colors, NewColor(c))
	}

	for i, color := range colors {
		if r[i] != color {
			t.Errorf("NewColor wanted %v, got %v", r[i], color)
		}
	}
}

func TestColor_AnsiString(t *testing.T) {
	r := []byte{
		033, '[', '3', '8', ';', '5', ';', '1', 'm', 'R', 'e', 'd', 033, '[', '0', 'm',
	}

	g := []byte{
		033, '[', '3', '8', ';', '5', ';', '2', '8', 'm', 'G', 'r', 'e', 'e', 'n', 033, '[', '0', 'm',
	}

	red := Red
	if s := red.AnsiString("Red"); strings.Compare(s, string(r)) != 0 {
		t.Errorf("ansistring: wanted %s, got %s", string(r), s)
	}

	green := Green
	if s := green.AnsiString("Green"); strings.Compare(s, string(g)) != 0 {
		t.Errorf("ansistring: wanted %s, got %s", string(g), s)
	}

	t.Log(red.AnsiString("Red"), green.AnsiString("Green"))
}

func TestColor_String(t *testing.T) {
	b := []byte{
		033, '[', '3', '8', ';', '5', ';', '1', 'm', 'R', 033, '[', '0', 'm',
		033, '[', '3', '8', ';', '5', ';', '2', '8', 'm', 'G', 033, '[', '0', 'm',
		033, '[', '3', '8', ';', '5', ';', '3', 'm', 'Y', 033, '[', '0', 'm',
		033, '[', '3', '8', ';', '5', ';', '6', 'm', 'B', 033, '[', '0', 'm',
		033, '[', '3', '8', ';', '5', ';', '1', '2', '8', 'm', 'V', 033, '[', '0', 'm',
		033, '[', '3', '8', ';', '5', ';', '2', '0', 'm', 'I', 033, '[', '0', 'm',
		033, '[', '3', '8', ';', '5', ';', '2', '0', '2', 'm', 'O', 033, '[', '0', 'm',
		'-',
	}

	r := []Color{Red, Green, Yellow, Blue, Violet, Indigo, Orange, NoColor}

	var sb strings.Builder
	for _, color := range r {
		sb.WriteString(color.String())
	}

	for i, c := range []byte(sb.String()) {
		if b[i] != c {
			t.Errorf("string: wanted %s, got %s", string(b), sb.String())
		}
	}

	t.Logf(sb.String())
}
