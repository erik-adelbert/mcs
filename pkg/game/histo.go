// Copyright 2018 Erik Adelbert. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package game

import "fmt"

// An Histogram counts blocks by color.
type Histogram map[Color]float64

func (h Histogram) String() string {
	return fmt.Sprint(map[Color]float64(h))
}
