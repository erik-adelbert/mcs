// Copyright 2018 Erik Adelbert. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chaingame

// A Tag is composed of an ID and a color. It's used when tiling
// boards. It uniquely identifies blocks of a same group.
type Tag int

// NewTag returns an encoded tag.
func NewTag(id int, c Color) Tag {
	return Tag((id << 4) | int(c))
}

// Color decodes the color from a tag.
func (t Tag) Color() Color {
	return Color(int(t) & 0xf)
}

// ID decodes the unique ID of a tag.
func (t Tag) ID() int {
	return int(t) >> 4
}

// Tags is a collection of unique tags.
type Tags map[Tag]Tag

// NewTags returns a newly allocated tag collector.
func NewTags(cap int) Tags {
	return make(Tags, cap)
}

// List all tags.
func (t Tags) List() Tags {
	delete(t, 0) // Remove id facility
	return t
}

// NewID returns a unique incrementing ID.
// There's no Tags #0: Tags[0] stores the current set number.
func (t Tags) NewID(c Color) Tag {
	id := t[0] + 1
	t[0] = id

	return NewTag(int(id), c)
}

// Find implements union-find with path splitting. It's used during board
// tiling (labeling) in support of an Hoshen–Kopelman like algorithm.
func (t Tags) Find(x Tag) Tag {

	for t[x] != x {
		x, t[x] = t[x], t[t[x]] // path splitting
	}

	return x
}

// Union implements union-find. It's used during board tiling (labeling) in
// support of an Hoshen–Kopelman like algorithm.
func (t Tags) Union(x, y Tag) Tag {

	if _, ok := t[x]; !ok {
		t[x] = x
	}

	if _, ok := t[y]; !ok {
		t[y] = y
	}

	order := func(x, y Tag) (Tag, Tag) {
		if x < y {
			return x, y
		}
		return y, x
	}

	x, y = order(t.Find(x), t.Find(y))
	t[y] = x

	return x
}
