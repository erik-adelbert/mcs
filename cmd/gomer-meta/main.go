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
	"mcs/pkg/chaingame"
	"mcs/pkg/mcs"
)

const (
	KB = 1024

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

	if len(*input) == 0 {
		if err := fmt.Errorf("no input file"); err != nil {
			panic(err)
		}
		return
	}

	writer := bufio.NewWriterSize(os.Stdout, 1*KB)

	rand.Seed(time.Now().UnixNano())

	var timeout = defaultTimeout
	if len(*duration) > 0 {
		duration, err := time.ParseDuration(*duration)
		if err != nil {
			timeout = duration
		}
	}

	h, w, board := load(input)
	b := samegame.NewSameBoard(h, w)
	b.Load(board)

	writeln(writer, b.String())
	flush(writer)

	{ // Everything has been taken care of, this is the core code:
		gs := mcs.GameState(b)

		policies := []mcs.GamePolicy{
			mcs.GamePolicy(samegame.TabooColor),
		}

		root := mcs.NewRoot(gs, ε, C, W)

		start := time.Now()
		result := mcs.MetaSearch(root, policies, timeout)
		elapsed := time.Since(start)

		replay(writer, b, result)

		writeln(writer, fmt.Sprintf("Meta took %v value: %v (%d nodes)", elapsed, result.Score(), mcs.NodeCount()))
		flush(writer)

		if *interactive {
			mcs.Cli(root)
		}
	}

	//panic("Show stack")
}

func load(fname *string) (h, w int, board []string) {
	file, err := os.Open(*fname)
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReaderSize(file, 1*KB)

	fields := strings.Split(readln(reader), " ")
	h, w = atoi(fields[0]), atoi(fields[1])

	board = make([]string, 0, h)
	for i := 0; i < h; i++ {
		board = append(board, readln(reader))
	}

	if err := file.Close(); err != nil {
		panic(err)
	}

	return
}

func replay(writer *bufio.Writer, b samegame.SameBoard, solution mcs.Decision) {
	moves := solution.Moves()
	for i, tile := range moves {
		color := b.TileColor(chaingame.Tile(tile))
		b = b.Remove(chaingame.Tile(tile))

		writeln(writer, fmt.Sprintf("\n#%d Removed: %s", i+1, color.AnsiColoredString(tile.String())))
		writeln(writer, b.String())
	}
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

func flush(w *bufio.Writer) {
	if err := w.Flush(); err != nil {
		panic(err)
	}
}
