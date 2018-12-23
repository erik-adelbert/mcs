// Copyright 2018 Erik Adelbert. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chaingame

import "fmt"

// A block is a board cell.
type block struct {
	r, c int
}

// Column is an obvious getter.
func (b block) Column() int {
	return b.c
}

// Row is an obvious getter.
func (b block) Row() int {
	return b.r
}

func (b block) String() string {
	return fmt.Sprintf("(%d,%d)", b.r, b.c)
}
