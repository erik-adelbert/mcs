// Copyright 2018 Erik Adelbert. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package game

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

// NewColor encodes string characters into colors
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

// AnsiString returns a colorized string by embedding ansi color escape codes.
func (c Color) AnsiString(s string) string {
	var code string

	switch c {
	default:
		return s
	case Red:
		code = "1"
	case Green:
		code = "28"
	case Yellow:
		code = "3"
	case Blue:
		code = "6"
	case Violet:
		code = "128"
	case Indigo:
		code = "20"
	case Orange:
		code = "202"
	}
	ansiColor := "\033[38;5;" + code + "m"
	ansiReset := "\033[0m"

	return ansiColor + s + ansiReset
}

// This stringer prints an ansi-colorized letter.
func (c Color) String() string {
	var code, ch string

	switch c {
	default:
		return "-"
	case Red:
		ch, code = "R", "1"
	case Green:
		ch, code = "G", "28"
	case Yellow:
		ch, code = "Y", "3"
	case Blue:
		ch, code = "B", "6"
	case Violet:
		ch, code = "V", "128"
	case Indigo:
		ch, code = "I", "20"
	case Orange:
		ch, code = "O", "202"
	}

	ansiColor := "\033[38;5;" + code + "m"
	ansiReset := "\033[0m"

	return ansiColor + ch + ansiReset
}
