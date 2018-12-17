// Copyright 2018 Erik Adelbert. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file implements various Upper Confidence Bound formulas.

package mcs

import (
	"math"
)

// SelectedUCB is the current in use formula. It isn't currently possible
// to have more than one in use.
var SelectedUCB = UCBTunedSinglePlayer

// An UCB function is an effective implementation of a formula.
type UCB func(*Node) float64

// SelectUCB chooses the current active formula.
func SelectUCB(fun UCB) {
	SelectedUCB = fun
}

// UCB1 is from [2002 Auer et Al]
// see https://homes.di.unimi.it/~cesabian/Pubblicazioni/ml-02.pdf
func UCB1(n *Node) float64 {
	var np, ni, μι, C float64

	np, ni = n.up.Visits(), n.Visits()

	n.Lock()
	{
		μι = n.mean
		C = n.c
	}
	n.Unlock()

	χi := math.Sqrt(math.Log(np) / ni)

	value := μι + C*χi

	return value
}

// UCBTunedSinglePlayer is from [2012 Beyer, Winands]
// see https://dke.maastrichtuniversity.nl/m.winands/documents/ecai2012.pdf
func UCBTunedSinglePlayer(n *Node) float64 {
	var np, ni, βi, μι, σι, C, W float64

	np, ni = n.up.Visits(), n.Visits()

	n.Lock()
	{
		βi = n.best.Score()
		μι = n.mean
		σι = n.variance
		C = n.c
		W = n.w // weighted best value
	}
	n.Unlock()

	// running average value

	χi := math.Sqrt(2 * math.Log(np) / ni)

	if σι = σι + χi; σι > 1.0/4.0 {
		σι = 1.0 / 4.0
	}
	σι = math.Sqrt(χi * σι)

	value := μι + C*σι + W*βi

	return value
}

// UCBV is from [2007 Audibert et al.]
// see http://certis.enpc.fr/~audibert/ucb_alt.pdf
func UCBV(n *Node) float64 {
	var np, ni, βi, μι, σι, C, W float64

	np, ni = n.up.Visits(), n.Visits()

	n.Lock()
	{
		βi = n.best.Score()
		μι = n.mean
		σι = n.variance
		C = n.c
		W = n.w // weighted best value
	}
	n.Unlock()

	χi := C * math.Log(np) / ni

	σι = math.Sqrt(2 * σι * χi)

	βi = 3 * W * βi * χi

	value := μι + σι + βi

	return value
}

// TODO: ADA-UCB [2018 Lattimore]
// see http://www.jmlr.org/papers/volume19/17-513/17-513.pdf
