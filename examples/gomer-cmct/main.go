// Copyright 2018 Erik Adelbert. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
	"strings"
	"time"

	"mcs/games/samegame"
	"mcs/pkg/game"
	"mcs/pkg/mcs"
)

const (
	defaultTimeout = 1 * time.Minute // This program has in +/- 10ms accuracy due to its structure.

	ε = 0.03 // ε-greedy

	C = 40  // UCB-SP
	W = 0.0 // UCB Best score
)

var (
	input       = flag.String("f", "", "problem file")
	duration    = flag.String("t", "", "timeout")
	interactive = flag.Bool("i", false, "launch an inspection console at the end of the search")
	profiling   = flag.Bool("pprof", false, "launch a live profiling web service on port 6060")
)

func main() {
	flag.Parse()

	if *profiling {
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}

	var timeout = defaultTimeout
	if len(*duration) > 0 {
		duration, err := time.ParseDuration(*duration)
		if err != nil {
			timeout = duration
		}
	}

	if len(*input) == 0 {
		if err := fmt.Errorf("no input file"); err != nil {
			panic(err)
		}
		return
	}

	file, err := os.Open(*input)
	if err != nil {
		panic(err)
	}

	const KB = 1024

	reader := bufio.NewReaderSize(file, 1*KB)
	writer := bufio.NewWriterSize(os.Stdout, 1*KB)

	fields := strings.Split(readln(reader), " ")
	h, w := atoi(fields[0]), atoi(fields[1])

	board := make([]string, 0, h)
	for i := 0; i < h; i++ {
		board = append(board, readln(reader))
	}

	rand.Seed(time.Now().UnixNano())

	b := samegame.NewSameBoard(h, w)
	b.Load(board)

	gs := mcs.GameState(samegame.State(b))

	writeln(writer, b.String())
	if err := writer.Flush(); err != nil {
		panic(err)
	}

	policies := []mcs.GamePolicy{
		mcs.GamePolicy(samegame.TabooColor),
		//mcs.GamePolicy(samegame.TabooColor),
		//mcs.GamePolicy(samegame.NoTaboo),
	}

	root := mcs.NewRoot(gs, ε, C, W)

	start := time.Now()
	result := mcs.ConcurrentSearch(root, policies, timeout)
	elapsed := time.Since(start)

	score, tiles := result.Score(), result.Moves()

	clone := b.Clone()
	for i, tile := range tiles {
		color := clone.TileColor(game.Tile(tile))
		clone = clone.Remove(game.Tile(tile))
		writeln(writer, fmt.Sprintf("\n#%d Removed: %s", i+1, color.AnsiString(tile.String())))
		writeln(writer, clone.String())
	}
	writeln(writer, fmt.Sprintf("TimedCMCT took %v value: %v (%d nodes)", elapsed, score, mcs.NodeCount()))

	if err := writer.Flush(); err != nil {
		panic(err)
	}

	if *interactive {
		mcs.Cli(root)
	}

	if err := file.Close(); err != nil {
		panic(err)
	}

	//panic("Show stack")
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

func writeln(w *bufio.Writer, s string) {
	if _, err := w.WriteString(s); err != nil {
		panic(err)
	}
	if err := w.WriteByte('\n'); err != nil {
		panic(err)
	}
}
