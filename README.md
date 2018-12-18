# mcs
mcs is a package implementing CMCT, a Concurrent Monte-Carlo Search for Tree akin to a concurrent UCT.

```bash 
$ cd go/src
$ git clone https://github.com/erik-adelbert/mcs.git
$ cd mcs/examples/gomer-uct/
$ go run ./main.go -i -f ../../assets/www.js-game.de/problem01.txt
```

Benchmark is available but will take at least 20\*50mn to complete and generate a lot of log files:

```bash 
$ cd go/src
$ git clone https://github.com/erik-adelbert/mcs.git
$ cd mcs/examples/benchmarks/
$ go test -timeout 0
```

There's also a Makefile:

```bash 
$ cd go/src
$ git clone https://github.com/erik-adelbert/mcs.git
$ make check
$ make bench
```

See the [up to date documentation.](https://godoc.org/github.com/erik-adelbert/mcs/pkg/mcs)\
See also the [Clickomania/Samegame boardgame engine documentation.](https://godoc.org/github.com/erik-adelbert/mcs/pkg/game)

**_This is an experimental setup which is currently studied. Do not use for production purposes._**
