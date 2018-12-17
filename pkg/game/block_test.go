package game

import (
	"fmt"
	"strings"
	"testing"
)

func TestBlock_Column(t *testing.T) {
	r, c := 1, 2
	b := Block{r, c}

	if col := b.Column(); col != c {
		t.Errorf("column: %v expected %d, got %d", b, c, col)
	}
}

func TestBlock_Row(t *testing.T) {
	r, c := 1, 2
	b := Block{r, c}

	if row := b.Row(); row != r {
		t.Errorf("row: %v expected %d, got %d", b, r, row)
	}
}

func TestBlock_String(t *testing.T) {
	r, c := 1, 2
	expected := fmt.Sprintf("(%d,%d)", r, c)

	b := Block{r, c}

	if s := b.String(); strings.Compare(s, expected) != 0 {
		t.Errorf("string: expected %s, got %s", expected, s)
	}

}
