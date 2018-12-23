// Copyright 2018 Erik Adelbert. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chaingame

// A tag is composed of an ID and a color. It's used when tiling
// boards. It uniquely identifies blocks of a same group.
type tag int

// newTag returns an encoded tag.
func newTag(id int, c Color) tag {
	return tag((id << 4) | int(c))
}

// Color decodes the color from a tag.
func (t tag) Color() Color {
	return Color(int(t) & 0xf)
}

// ID decodes the unique ID of a tag.
func (t tag) ID() int {
	return int(t) >> 4
}

// tags is a collection of unique tags.
type tags map[tag]tag

// newTags returns a newly allocated tag collector.
func newTags(cap int) tags {
	return make(tags, cap)
}

// List all tags.
func (t tags) List() tags {
	delete(t, 0) // Remove id facility
	return t
}

// NewID returns a unique incrementing ID.
// There's no tags #0: tags[0] stores the current set number.
func (t tags) NewID(c Color) tag {
	id := t[0] + 1
	t[0] = id

	return newTag(int(id), c)
}

// Find implements union-find with path splitting. It's used during board
// tiling (labeling) in support of an Hoshen–Kopelman like algorithm.
func (t tags) Find(x tag) tag {

	for t[x] != x {
		x, t[x] = t[x], t[t[x]] // path splitting
	}

	return x
}

// Union implements union-find. It's used during board tiling (labeling) in
// support of an Hoshen–Kopelman like algorithm.
func (t tags) Union(x, y tag) tag {

	if _, ok := t[x]; !ok {
		t[x] = x
	}

	if _, ok := t[y]; !ok {
		t[y] = y
	}

	order := func(x, y tag) (tag, tag) {
		if x < y {
			return x, y
		}
		return y, x
	}

	x, y = order(t.Find(x), t.Find(y))
	t[y] = x

	return x
}
