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
	root  *node
	comp  CompareFunc
	count int
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
	insertEqual(&t.root, t.comp, e)
	t.count++
}

func (t *RBTree) InsertUnique(e Element) bool {
	_, ok := insertUnique(&t.root, t.comp, e)
	if ok {
		t.count++
	}
	return ok
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
	t.count--
}

func (t *RBTree) RemoveByIterator(i *Iterator) {
	t.root = remove(t.root, i.node)
	t.count--
}

func (t *RBTree) Size() int {
	return t.count
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
