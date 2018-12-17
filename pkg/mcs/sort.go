// Copyright 2018 Erik Adelbert. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This sort interface is adapted from the SortKeys example of the sort package.
// see https://golang.org/pkg/sort/

package mcs

import "sort"

// Value is a "less" function that defines nodes order by their UCB values.
var Value = func(n1, n2 **Node) bool { // Sort by value
	return (*n1).Value() < (*n2).Value()
}

// By is the type of a "less" function that defines the ordering of nodes
// in a Monte-Carlo tree.
type By func(n1, n2 **Node) bool

// SortDescending sorts nodes in reverse order.
func (by By) SortDescending(nodes []*Node) {
	ns := &nodeSorter{nodes: nodes, by: by}
	sort.Sort(sort.Reverse(ns))
}

type nodeSorter struct {
	nodes []*Node
	by    func(n1, n2 **Node) bool
}

// Len is part of sort.Interface.
func (s *nodeSorter) Len() int {
	return len(s.nodes)
}

// Swap is part of sort.Interface.
func (s *nodeSorter) Swap(i, j int) {
	s.nodes[i], s.nodes[j] = s.nodes[j], s.nodes[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *nodeSorter) Less(i, j int) bool {
	return s.by(&s.nodes[i], &s.nodes[j])
}
