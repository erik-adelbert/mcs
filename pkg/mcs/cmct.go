// Copyright 2018 Erik Adelbert. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file implements a Concurrent Monte-Carlo search for Trees (CMCT)

package mcs

import (
	"log"
	"runtime"
	"sync"
	"time"
)

// VisitThreshold is the minimal number of simulations a position has to go through
// before being registered in the tree as a node.
const VisitThreshold = 8

// The minimum number of walkers is 2 while there are always as many updaters and
// twice more samplers.
func numGoRoutines() int {
	switch n := runtime.NumCPU() / 8; n {
	case 0:
		return 2
	default:
		return 2 * n
	}
}

var walkers = numGoRoutines()
var updaters = walkers
var samplers = 2 * walkers // samplers are the slowest

// Jobs convey nodes and best moves between mcts steps (ie. walkers, samplers and updaters).
type job struct {
	node     *Node
	decision Decision
}

// ConcurrentSearch has common roots with tree parallelization in [2008 Chaslot, Winands et al.].
// Instead of spawning multiple threads that runs all the steps of mcts, CMCT spawns walkers,
// samplers and updaters: walkers realize selections and expansions (steps 1&2) while samplers
// play full random games (step 3) and updaters back-propagate results and compute UCB values.
// We state that it isn't mandatory to synchronize mcts steps tightly.
// see:
// high scores are on http://www.js-games.de/eng/highscores/samegame/lx (results registered as cmct)
// http://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.159.4373&rep=rep1&type=pdf
func ConcurrentSearch(root *Node, policies []GamePolicy, duration time.Duration) Decision {

	if root == nil {
		// TODO: error handling
		panic("no root")
	}

	// All possible first moves are expanded
	tree := GrowTree(root)

	done := make(chan struct{})
	timeout := make(chan bool)

	// Start the countdown asap: the search must return a decision
	// in the given time frame. A late decision is as good as an
	// illegal move, it's disqualifying.
	go func() {
		time.Sleep(duration)
		close(timeout)
	}()

	// Prepare pipelines (channels and goroutines launchers).
	positions := make(chan job, samplers)
	walk := func(count int) {
		var wg sync.WaitGroup

		wg.Add(count)
		for i := 0; i < count; i++ {
			go func() {
				walker(done, tree, positions)
				wg.Done()
			}()
		}

		go func() {
			wg.Wait()
			close(positions)
		}()
	}

	outcomes := make(chan job, samplers)
	sample := func(count int) {
		var wg sync.WaitGroup

		wg.Add(count)
		for i := 0; i < count; i++ {
			go func() {
				sampler(done, policies, positions, outcomes)
				wg.Done()
			}()
		}

		go func() {
			wg.Wait()
			close(outcomes)
		}()
	}

	update := func(count int) {
		for i := 0; i < count; i++ {
			go func() {
				updater(done, outcomes)
			}()
		}
	}

	// Launch!
	go update(updaters)
	go sample(samplers)
	go walk(walkers)

	// Wait for either timeout or solution
	for {
		select {
		case <-timeout:
			goto conclusion
		default:
			if root.IsSolved() {
				goto conclusion
			}
			runtime.Gosched()
		}
	}

	// Broadcast termination message (done) to all goroutines
	// and return the best sequence found so far.
conclusion:
	close(done)
	return tree.Best()
}

// A sampler is the slowest performer of the asynchronous pipeline. This is why there are twice
// more samplers than other kinds of goroutine: the assumption is that loading up the pipeline
// with simulation will eventually reduce dead time in walkers and updaters.
func sampler(done <-chan struct{}, policies []GamePolicy, position <-chan job, outcome chan<- job) {
	log.Println("sampler up")

	for task := range position {
		node, decision := task.node, task.decision

		if node == nil {
			continue
		}

		node.Lock()
		state := node.state.Clone()
		node.Unlock()

		switch node.Status() {
		case walked:
			node.SetStatus(simulating)
		default:
			//fmt.Printf("simulating %v node %d\n", node.Status(), node.Status())
			//continue
		}

		sampled := decision.Join(state.Sample(done, policies[0]))

		select {
		case <-done:
			log.Println("sampler done")
			return
		case outcome <- job{node, sampled}:
			node.SetStatus(simulated)
		}
	}
}

// An updater is asynchronously back propagating scores received from simulating.
// It computes UCB values along the way.
func updater(done <-chan struct{}, outcome <-chan job) {
	log.Println("updater up")

	for outcome := range outcome {
		select {
		case <-done:
			log.Println("updater done")
			return
		default:
			node, decision := outcome.node, outcome.decision
			//if node != nil && (node.Status() != simulated && node.Status() != simulating) {
			if node != nil {
				node.UpdateTree(decision)
				node.SetStatus(idle)
			}
		}
	}
}

// A walker share the very same logic as UCT: it realizes selections and expansions of nodes.
// It chooses moves to address the dilemma between exploration or exploitation.
func walker(done <-chan struct{}, root *Node, position chan<- job) {
	log.Println("walker up")

	var out chan<- job

	for {
		var score float64
		var moves MoveSequence

		out = nil
		node := root

		for node.IsExpanded() {
			node = node.Downselect()

			move := node.Edge()
			moves = moves.Enqueue(move)
			score += move.Score()
		}

		if !node.IsTerminal() && node.Visits() > VisitThreshold {
			node = node.ExpandOne(node.RandomNewEdge())

			move := node.Edge()
			moves = moves.Enqueue(move)
			score += move.Score()
		}

		if node != nil {
			node.SetStatus(walked)
			out = position // activate channel
		}

		select {
		case <-done:
			log.Println("walker done")
			return
		case out <- job{node, Decision{score: score, moves: moves}}:
			// pass along if channel is activated
		default:
			// continue
		}
	}
}
