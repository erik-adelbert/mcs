package game

import (
	"strings"
	"testing"
)

func TestHistogram_String(t *testing.T) {
	h := Histogram{Red: 1, Green: 2}
	red, green := Red, Green

	expected1 := "map[" + red.AnsiString("R") + ":1 " + green.AnsiString("G") + ":2]"
	expected2 := "map[" + green.AnsiString("G") + ":2 " + red.AnsiString("R") + ":1]"

	if s := h.String(); strings.Compare(s, expected1) != 0 && strings.Compare(s, expected2) != 0 {
		t.Errorf("ansistring: wanted %s or %s, got %s", expected1, expected2, s)
	}

	t.Log(h)
}
