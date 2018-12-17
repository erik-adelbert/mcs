// Copyright 2018 Erik Adelbert. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file implements a CLI which enable navigation of Monte-Carlo Trees

package mcs

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	path2 "path"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

var (
	zeroValue = reflect.ValueOf(nil)
	once      sync.Once

	node *Node
	wdir = "/"
)

// Cli launches an interactive console which has basic Monte-Carlo trees navigational capacities.
// Available commands are listed when invoking 'help'. Node's fields are accessible trough their
// getters: to get a field's value, it's sufficient to type its name.
// Path are supported : '/0/1/state' will display the board associated with the second child
// of the first child of root.
func Cli(root *Node) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			Cli(node)
		}
	}()

	once.Do(func() { node = root })

	fmt.Println("Type 'help' to start")
	prompt()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		varargs := strings.Split(strings.TrimSpace(strings.Title(scanner.Text())), " ")
		if len(varargs) == 0 {
			continue // empty line
		}

		// extract (possibly) a command and an argument list
		command := cmd(varargs[0])
		if len(varargs) > 1 {
			varargs = varargs[1:]
		} else {
			varargs = nil
		}

		switch command {
		case "Cd":
			if len(varargs) == 0 {
				for node.up != nil {
					command.ascend()
				}
			} else {
				if err := command.Chdir(varargs[0]); err != nil {
					fmt.Println(err)
				}
			}
		case "Exit", "Quit":
			return

		case "?", "Help":
			command.Help()

		case "Up":
			command.Up()

		case "Down":
			command.Down()

		default:

			dir, base := path2.Split(path2.Clean(string(command)))
			command = cmd(base)

			currNode, currWdir := node, wdir

			if err := command.Chdir(dir); err != nil {
				fmt.Println(err)
			}

			args := make([]reflect.Value, 0, len(varargs))
			for _, vararg := range varargs {
				if arg := reflect.ValueOf(vararg); arg != zeroValue {
					args = append(args, arg)
				}
			}

			if fun := reflect.ValueOf(command).MethodByName(string(command)); fun != zeroValue {
				fun.Call(args)
			} else {
				command.Dump(string(command))
			}

			node, wdir = currNode, currWdir
		}

		prompt()
	}
}

type cmd string

func (c cmd) Chdir(path string) error {
	path = path2.Clean(path)

	if path[0] == '/' {
		for node.up != nil {
			c.ascend()
		}
		wdir = "/"

		if len(path) == 1 {
			return nil
		}

		path = path[1:]
	}

	for _, x := range strings.Split(path, "/") {

		switch x {
		case ".":
		case "..":
			c.ascend()
		default:
			if err := c.descend(x); err != nil {
				return errors.New("chdir: invalid path")
			}
		}
	}

	return nil
}

func (c cmd) Count() {
	fmt.Println(walk(node))
}

func (c cmd) Down() {
	var sb strings.Builder
	for i := 0; i < len(node.down); i++ {
		if _, err := fmt.Fprintf(&sb, "%2d: %v\t@%p\t", i, node.down[i].edge, node.down[i]); err != nil {
			panic(err)
		}

		if (i+1)%3 == 0 {
			sb.WriteByte('\n')
		}
	}
	sb.WriteByte('\n')
	fmt.Print(sb.String())
}

func (c cmd) Dump(fields ...string) {

	if len(fields) == 0 {
		fmt.Println(node)
		return
	}

	for _, field := range fields {
		switch {
		case field == "node":
			fmt.Println(node)
		default:
			if getter := reflect.ValueOf(node).MethodByName(field); getter != zeroValue {
				field := getter.Call(nil)[0]

				if field != zeroValue {
					fmt.Println(field)
					continue
				}

				if stringer := field.MethodByName("String"); stringer != zeroValue {
					fmt.Println(reflect.ValueOf(stringer.Call([]reflect.Value{})[0]))
				} else {
					fmt.Println(field)
				}
			}
		}
	}
	return
}

func (c cmd) Help() {
	command := reflect.ValueOf(c)

	names := make([]string, 0, command.NumMethod())
	for i := 0; i < command.NumMethod(); i++ {
		name := strings.ToLower(command.Type().Method(i).Name)

		switch name {
		case "chdir":
			name = "cd"
		}

		names = append(names, name)
	}

	var sb strings.Builder
	for i, name := range names {
		sb.WriteString(name + "\t")
		if i > 0 && i%7 == 6 {
			sb.WriteByte('\n')
		}
	}
	fmt.Println(sb.String())
}

func (c cmd) Ls() {
	node := reflect.ValueOf(*node)

	fields := make([]string, 0, node.NumField())
	for i := 0; i < node.NumField(); i++ {
		fields = append(fields, node.Type().Field(i).Name)
	}
	fields = fields[1:] // Loose the lock

	var sb strings.Builder
	for i, field := range fields {
		sb.WriteString(field + "\t")
		if i > 0 && i%8 == 7 {
			sb.WriteByte('\n')
		}
	}
	fmt.Println(sb.String())
}

func (c cmd) Up() {
	if _, err := fmt.Printf("@%p\n", node.up); err != nil {
		panic(err)
	}
}

func (c cmd) ascend() {
	if node.up != nil {
		node = node.up
		wdir = path2.Dir(wdir)
	}
}

func (c cmd) descend(idx string) error {
	var (
		i   int
		err error
	)

	if i, err = strconv.Atoi(idx); err == nil {
		if i < len(node.down) {
			node = node.down[i]
			wdir = path2.Join(wdir, idx)
		} else {
			return errors.New("godown: invalid node")
		}
	}

	return err
}

func prompt() {
	fmt.Printf("node@%p "+wdir+"> ", node)
}

func walk(n *Node) int {
	count := 1

	for _, child := range n.Down() {
		count += walk(child)
	}

	return count
}
