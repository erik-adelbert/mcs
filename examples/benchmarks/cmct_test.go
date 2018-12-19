// Copyright 2018 Erik Adelbert. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This benchmark is to be launched by 'go test -timeout 0'

package benchmark

import (
	"bufio"
	"fmt"
	"math/rand"
	"mcs/bencher"
	"mcs/games/samegame"
	"mcs/pkg/mcs"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

const (
	StandardSetPath = "../../assets/www.js-games.de/"
	KB              = 1024
)

func TestSameGameStandardSet(t *testing.T) {
	// 3...
	problems := make([]*bencher.Problem, 0, 20)
	for i := 1; i <= 20; i++ {

		name := fmt.Sprintf("problem%02d", i)
		problem := bencher.NewProblem(name)
		problem.SetInitial(load(StandardSetPath + name + ".txt"))

		problems = append(problems, problem)
	}

	// 2...
	searchers := []*bencher.Searcher{
		bencher.NewSearcher("Concurrent"),
		bencher.NewSearcher("Meta"),
		bencher.NewSearcher("Confident"),
	}
	searchers[0].SetFun(mcs.ConcurrentSearch)
	searchers[1].SetFun(mcs.MetaSearch)
	searchers[2].SetFun(mcs.ConfidentSearch)

	// 1...
	policies := []mcs.GamePolicy{
		samegame.TabooColor,
		samegame.TabooColor,
	}

	constants := []struct{ ε, C, W float64 }{
		{0.03, 1000, 0},
		{0.03, 40, 0},
		//{0.03, 40, 0.2},
	}

	durations := []time.Duration{
		10 * time.Minute,
		//20 * time.Minute,
		//40 * time.Minute,
		//60 * time.Minute,
	}

	for _, searcher := range searchers {
		searcher.SetPolicies(policies)
		for _, set := range constants {
			searcher.SetΕ(set.ε)
			searcher.SetC(set.C)
			searcher.SetW(set.W)
			for _, duration := range durations {
				logname := "StandardSet" + searcher.Name() + duration.String() + "C" + fmt.Sprint(set.C)
				logfile, err := os.Create(logname + ".log")
				if err != nil {
					t.Fatal(err)
				}
				logger := bufio.NewWriterSize(logfile, 1*KB)

				for _, problem := range problems {

					game := bencher.NewGame()
					game.SetP(problem)
					game.SetS(searcher)

					name := logname + "_" + problem.Name()
					benchmark, _ := bencher.NewBenchmark(name, duration, logger)
					benchmark.Duration = duration
					benchmark.Attach(game)

					if err := benchmark.Run(); err != nil {
						t.Fatal(err)
					}

					// Force memory recuperation before next iteration
					runtime.GC()
				}
				close(logfile)
			}

		}
	}

	return
}

func atoi(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return n
}

func readln(reader *bufio.Reader) string {
	s, _, err := reader.ReadLine()
	if err != nil {
		panic(err)
	}
	return string(s)
}

func close(file *os.File) {
	if err := file.Close(); err != nil {
		panic(err)
	}
}

func load(name string) mcs.GameState {
	file, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer close(file)

	const KB = 1024

	reader := bufio.NewReaderSize(file, 1*KB)

	fields := strings.Split(readln(reader), " ")
	h, w := atoi(fields[0]), atoi(fields[1])

	board := make([]string, 0, h)
	for i := 0; i < h; i++ {
		board = append(board, readln(reader))
	}

	rand.Seed(time.Now().UnixNano())

	b := samegame.NewSameBoard(h, w)
	b.Load(board)

	return mcs.GameState(samegame.State(b))
}
