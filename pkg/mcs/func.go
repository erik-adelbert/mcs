// Copyright 2018 Erik Adelbert. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mcs

import (
	"reflect"
	"runtime"
	"time"
)

// Search is a function that implements a Monte-Carlo technique.
type Search func(*Node, []GamePolicy, time.Duration) Decision

func (s Search) String() string {
	return runtime.FuncForPC(reflect.ValueOf(s).Pointer()).Name()
}
