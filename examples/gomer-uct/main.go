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
	"os"
	"strconv"
	"strings"
	"time"

	"mcs/games/samegame"
	"mcs/pkg/game"
	"mcs/pkg/mcs"
)

const (
	timeout = 1 * time.Minute // This timer has in +/- 10ms accuracy due to programming structure.

	ε = 0.03 // ε-greedy

	C = 0.012 // UCB-SP
	W = 0.0   // UCB Best score
)

var profiling = flag.Bool("pprof", false, "launch a live profiling web service on port 6060")

func main() {
	flag.Parse()
	if *profiling {
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}

	const KB = 1024

	reader := bufio.NewReaderSize(os.Stdin, 1*KB)
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
		panic("flush")
	}

	start := time.Now()

	policies := []mcs.GamePolicy{samegame.TabooColor}
	root := mcs.NewRoot(gs, ε, C, W)

	result := mcs.ConfidentSearch(root, policies, timeout)
	score, tiles := result.Score(), result.Moves()

	elapsed := time.Since(start)

	fmt.Println(tiles, tiles.Len())

	for i, tile := range tiles {
		tile := game.Tile(tile)
		color := b.TileColor(tile)
		b = b.Remove(tile)
		writeln(writer, fmt.Sprintf("\n#%d Removed: %s", i+1, color.AnsiString(tile.String())))
		writeln(writer, b.String())
	}
	writeln(writer, fmt.Sprintf("TimedUCT took %v value: %v (%d nodes)", elapsed, score, mcs.NodeCount()))
	if err := writer.Flush(); err != nil {
		panic("flush")
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
