// Copyright 2018 Erik Adelbert. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chaingame

// Color is the color of a block.
type Color byte

// Possible colors.
const (
	NoColor Color = iota
	Red
	Green
	Yellow
	Blue
	Violet
	Indigo
	Orange
	AllColors
)

// NewColor encodes characters into colors
func NewColor(b int32) Color {
	switch b {
	case 'R':
		return Red
	case 'G':
		return Green
	case 'Y':
		return Yellow
	case 'B':
		return Blue
	case 'V':
		return Violet
	case 'I':
		return Indigo
	case 'O':
		return Orange
	default:
		return NoColor
	}
}

type ansiColorCode string

var lut = map[Color]ansiColorCode{
	Red:    "1",
	Green:  "28",
	Yellow: "3",
	Blue:   "6",
	Violet: "128",
	Indigo: "20",
	Orange: "202",
}

func ansiColorEscape(color ansiColorCode) string {
	return "\033[38;5;" + string(color) + "m"
}

const ansiColorReset = "\033[0m"

// AnsiColoredString returns a colorized string by embedding ansi color escape codes.
func (c Color) AnsiColoredString(s string) string {
	var code ansiColorCode

	code, ok := lut[c]
	if !ok { // Color doesn't exist in the lookup table.
		return s
	}

	return ansiColorEscape(code) + s + ansiColorReset
}

// This stringer prints an ansi-colorized letter.
func (c Color) String() string {
	var s string

	switch c {
	default:
		return "-"
	case Red:
		s = "R"
	case Green:
		s = "G"
	case Yellow:
		s = "Y"
	case Blue:
		s = "B"
	case Violet:
		s = "V"
	case Indigo:
		s = "I"
	case Orange:
		s = "O"
	}

	return c.AnsiColoredString(s)
}
