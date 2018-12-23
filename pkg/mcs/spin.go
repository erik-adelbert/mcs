// Copyright 2018 Erik Adelbert. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mcs

import (
	"fmt"
	"runtime"
	"sync/atomic"
)

const (
	unlocked = iota
	locked
)

// spinlock are locks suited for very short period of time.
// see https://en.wikipedia.org/wiki/Spinlock
type spinlock struct {
	f uint32
}

// newSpinlock allocates and returns an unlocked spinlock.
func newSpinlock() *spinlock {
	return &spinlock{}
}

// GetLock will try to lock sl and return whether it succeed or not without blocking.
func (s *spinlock) GetLock() bool {
	return atomic.CompareAndSwapUint32(&s.f, unlocked, locked)
}

// Lock will simply wait in the loop until it can acquire the lock.
// This busy wait scheme won't hog the system.
func (s *spinlock) Lock() {
	for !s.GetLock() {
		runtime.Gosched()
	}
}

func (s *spinlock) String() string {
	if atomic.LoadUint32(&s.f) == locked {
		return fmt.Sprintf("Locked@%p", s)
	}
	return fmt.Sprintf("Unlocked@%p", s)
}

// Unlock cause no harm when called on an already loose lock.
func (s *spinlock) Unlock() {
	atomic.StoreUint32(&s.f, unlocked)
}
