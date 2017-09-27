package rbtree

import (
	"github.com/moladb/sl/container/stack"
)

const (
	Less   int = -1
	Equal  int = 0
	Greate int = 1
)

type CompareFunc func(left interface{}, right interface{}) int

type Element interface {
	Key() interface{}
	Value() interface{}
}

type RBTree struct {
	root *node
	comp CompareFunc
	size int
}

type Iterator struct {
	*node
}

func (i *Iterator) Key() interface{} {
	return i.e.Key()
}

func (i *Iterator) Value() interface{} {
	return i.e.Value()
}

func (i *Iterator) Next() *Iterator {
	n := i.node
	if n.right != nil {
		return &Iterator{unguardedMin(n.right)}
	}
	parent := n.parent
	for parent != nil {
		if parent.left == n {
			return &Iterator{parent}
		}
		n = parent
		parent = n.parent
	}
	return nil
}

func (i *Iterator) Prev() *Iterator {
	n := i.node
	if n.left != nil {
		return &Iterator{unguardedMax(n.left)}
	}
	parent := n.parent
	for parent != nil {
		if parent.right == n {
			return &Iterator{parent}
		}
		n = parent
		parent = n.parent
	}
	return nil
}

func New(comp CompareFunc) *RBTree {
	return &RBTree{comp: comp}
}

func (t *RBTree) First() *Iterator {
	if t.root == nil {
		return nil
	}
	return &Iterator{unguardedMin(t.root)}
}

func (t *RBTree) Last() *Iterator {
	if t.root == nil {
		return nil
	}
	return &Iterator{unguardedMax(t.root)}
}

func (t *RBTree) InsertEqual(e Element) {
	var (
		parent *node
		equal  int
	)
	n := t.root
	for n != nil {
		parent = n
		equal = t.comp(e.Key(), n.e.Key())
		if equal < 0 {
			n = n.left
		} else {
			n = n.right
		}
	}
	if parent == nil {
		t.root = &node{e: e}
	} else {
		if equal < 0 {
			parent.left = &node{parent: parent, e: e}
		} else {
			parent.right = &node{parent: parent, e: e}
		}
	}
	t.size++
}

func (t *RBTree) InsertUnique(e Element) bool {
	var parent *node
	var equal int
	n := t.root
	for n != nil {
		parent = n
		equal = t.comp(e.Key(), n.e.Key())
		if equal < 0 {
			n = n.left
		} else if equal > 0 {
			n = n.right
		} else {
			return false
		}
	}
	if parent == nil {
		t.root = &node{e: e, black: true}
	} else {
		if equal < 0 {
			parent.left = &node{parent: parent, e: e}
		} else {
			parent.right = &node{parent: parent, e: e}
		}
		t.rotate(parent)
	}
	t.size++
	return true
}

func (t *RBTree) rotate(n *node) {
	for {
		parent := n.parent
		oldN := n
		if isRed(n.right) && !isRed(n.left) {
			n = rotateLeft(n)
		}
		if isRed(n.left) && isRed(n.left.left) {
			n = rotateRight(n)
		}
		if isRed(n.left) && isRed(n.right) {
			n = flipColor(n)
		}
		if parent == nil {
			t.root = n
			t.root.black = true
		} else if parent.left == oldN {
			parent.left = n
		} else {
			parent.right = n
		}
		if parent == nil {
			break
		}
		if !isRed(parent) {
			break
		}
		n = parent
	}
}

func (t *RBTree) Max() Element {
	return unguardedMax(t.root).e
}

func (t *RBTree) Min() Element {
	return unguardedMin(t.root).e
}

func (t *RBTree) Find(key interface{}) *Iterator {
	n := search(t.root, t.comp, key)
	if n != nil {
		return &Iterator{n}
	}
	return nil
}

func (t *RBTree) RemoveByKey(key interface{}) {
	n := search(t.root, t.comp, key)
	if n == nil {
		return
	}
	t.root = remove(t.root, n)
	t.size--
}

func (t *RBTree) RemoveByIterator(i *Iterator) {
	t.root = remove(t.root, i.node)
	t.size--
}

func (t *RBTree) Size() int {
	return t.size
}

func (t *RBTree) InorderTraverse(fn func(e Element)) {
	s := stack.NewSliceStack()
	n := t.root
	for s.Size() != 0 || n != nil {
		for ; n != nil; n = n.left {
			s.Push(n)
		}
		n = s.Pop().(*node)
		fn(n.e)
		n = n.right
	}
}

// blackEdgesToMin/Max are used to test the tree is
// balanced at black edges
func (t *RBTree) blackEdgesToMin() int {
	n := t.root
	var c int
	for ; n.left != nil; n = n.left {
		if n.black {
			c++
		}
	}
	return c
}

func (t *RBTree) blackEdgesToMax() int {
	n := t.root
	var c int
	for ; n.right != nil; n = n.right {
		if n.black {
			c++
		}
	}
	return c
}
