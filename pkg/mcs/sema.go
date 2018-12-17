// Copyright 2018 Erik Adelbert. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mcs

// Semaphore is used to limit the number of active goroutines
type Semaphore chan int

// NewSemaphore returns a new allocated semaphore which can deliver 'size' tickets.
func NewSemaphore(size int) Semaphore {
	return make(Semaphore, size)
}

// Acquire will lock and keep a ticket in the calling goroutine.
func (s Semaphore) Acquire() {
	s <- 1
}

// Release will release back its ticket.
func (s Semaphore) Release() {
	<-s
}
