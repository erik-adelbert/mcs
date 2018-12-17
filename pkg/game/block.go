// Copyright 2018 Erik Adelbert. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package game

import "fmt"

// A Block is a board cell.
type Block struct {
	r, c int
}

// Column is an obvious getter.
func (b Block) Column() int {
	return b.c
}

// Row is an obvious getter.
func (b Block) Row() int {
	return b.r
}

func (b Block) String() string {
	return fmt.Sprintf("(%d,%d)", b.r, b.c)
}
