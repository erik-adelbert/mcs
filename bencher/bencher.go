// Copyright 2018 Erik Adelbert. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bencher

import (
	"bufio"
	"io"
	"log"
	"mcs/pkg/mcs"
	"time"
)

type Benchmark struct {
	time.Duration

	name string
	game *Game
	out  *log.Logger
}

func NewBenchmark(name string, duration time.Duration, logfile io.Writer) (*Benchmark, error) {
	const KB = 1024

	writer := bufio.NewWriterSize(logfile, 1*KB)
	out := log.New(writer, name, 0)

	return &Benchmark{Duration: duration, name: name, out: out}, nil
}

func (b *Benchmark) Attach(game *Game) {
	b.game = game
}

func (b *Benchmark) Detach() *Game {
	var game *Game

	game, b.game = b.game, nil
	return game
}

func (b *Benchmark) Run() error {
	mcs.NodeCountReset()

	game := b.game

	initial := game.p.Initial()
	ε, C, W := game.s.Ε(), game.s.C(), game.s.W()

	policies := game.s.Policies()

	fun, duration := game.s.fun, b.Duration

	var result mcs.Decision
	start := time.Now()
	{
		root := mcs.NewRoot(initial, ε, C, W)
		result = fun(root, policies, duration)
	}
	elapsed := time.Since(start)

	b.out.Println(" ", mcs.Search(fun).String(), elapsed, game.Name(), mcs.NodeCount(), result)

	return nil
}

type Game struct {
	name string
	p    *Problem
	s    *Searcher
}

func NewGame() *Game {
	return &Game{}
}

func (g *Game) Name() string {
	return g.name
}

func (g *Game) P() *Problem {
	return g.p
}

func (g *Game) S() *Searcher {
	return g.s
}

func (g *Game) SetName(name string) {
	g.name = name
}

func (g *Game) SetP(p *Problem) {
	g.p = p
}

func (g *Game) SetS(s *Searcher) {
	g.s = s
}

type Problem struct {
	name    string
	initial mcs.GameState
}

func NewProblem(name string) *Problem {
	return &Problem{name: name}
}

func (p *Problem) Initial() mcs.GameState {
	return p.initial
}

func (p *Problem) Name() string {
	return p.name
}

func (p *Problem) SetInitial(initial mcs.GameState) {
	p.initial = initial
}

func (p *Problem) SetName(name string) {
	p.name = name
}

type Searcher struct {
	name string

	fun mcs.Search

	policies []mcs.GamePolicy

	ε float64
	c float64
	w float64
}

func NewSearcher(name string) *Searcher {
	return &Searcher{name: name}
}

func (s *Searcher) C() float64 {
	return s.c
}

func (s *Searcher) Ε() float64 {
	return s.ε
}

func (s *Searcher) Name() string {
	return s.name
}

func (s *Searcher) Policies() []mcs.GamePolicy {
	return s.policies
}

func (s *Searcher) SetC(c float64) {
	s.c = c
}

func (s *Searcher) SetΕ(ε float64) {
	s.ε = ε
}

func (s *Searcher) SetFun(fun mcs.Search) {
	s.fun = fun
}

func (s *Searcher) SetPolicies(policies []mcs.GamePolicy) {
	s.policies = policies
}

func (s *Searcher) SetW(w float64) {
	s.w = w
}

func (s *Searcher) W() float64 {
	return s.w
}
