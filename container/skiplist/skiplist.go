// port from https://github.com/google/leveldb/blob/master/db/skiplist.h
package skiplist

import (
	"math/rand"
	"time"
)

type Element interface {
	Key() interface{}
	Value() interface{}
}

type node struct {
	e     Element
	level []*node
}

func newNode(e Element, height int) *node {
	return &node{
		e:     e,
		level: make([]*node, height, height),
	}
}

func (n *node) next(level int) *node {
	return n.level[level]
}

func (n *node) setNext(level int, x *node) {
	n.level[level] = x
}

const kMaxHeight int = 12

// left < right => < 0
// left == right => == 0
// left > right => > 0
type CompareFunc func(left interface{}, right interface{}) int

type SkipList struct {
	compare   CompareFunc
	head      *node
	maxHeight int
	rand      *rand.Rand
}

func NewSkipList(compare CompareFunc) *SkipList {
	return &SkipList{
		compare:   compare,
		head:      newNode(nil, kMaxHeight),
		maxHeight: 1,
		rand:      rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (l *SkipList) Insert(e Element) {
	var prev [kMaxHeight]*node
	x := l.findGreaterOrEqual(e.Key(), &prev)
	height := l.randomHeight()
	if height > l.maxHeight {
		for i := l.maxHeight; i < height; i++ {
			// 此处设置成 head
			// 对应下面 prev[i].setNext(i, x)
			// 将 head 指向高出来的
			prev[i] = l.head
		}
		l.maxHeight = height
	}
	x = newNode(e, height)
	for i := 0; i < height; i++ {
		x.setNext(i, prev[i].next(i))
		prev[i].setNext(i, x)
	}
}

func (l *SkipList) Remove(key interface{}) {
	var prev [kMaxHeight]*node
	x := l.findGreaterOrEqual(key, &prev)
	if x != nil && l.compare(key, x.e.Key()) == 0 {
		height := len(x.level)
		// delete x at each level linked list
		for i := 0; i < height; i++ {
			prev[i].setNext(i, x.next(i))
		}
		// update maxHeight
		for l.maxHeight > 1 && l.head.next(l.maxHeight-1) == nil {
			l.maxHeight--
		}
	}
}

func (l *SkipList) Contains(key interface{}) bool {
	n := l.searchGreaterOrEqual(key)
	return n != nil && l.compare(key, n.e.Key()) == 0
}

func (l *SkipList) Find(key interface{}) *Iterator {
	n := l.searchGreaterOrEqual(key)
	if n != nil && l.compare(key, n.e.Key()) == 0 {
		return &Iterator{
			node: n,
			list: l,
		}
	}
	return nil
}

func (l *SkipList) randomHeight() int {
	var height int = 1
	// 1/4的概率增加 这层的高度
	for height < kMaxHeight && l.rand.Int()%4 == 0 {
		height++
	}
	return height
}

func (l *SkipList) keyIsAfterNode(key interface{}, n *node) bool {
	return n != nil && l.compare(n.e.Key(), key) < 0
}

//func (l *SkipList) keyIsEqualNode(key interface{}, n *node) bool {
//	return n != nil && l.compare(n.e.Key(), key) == 0
//}

func (l *SkipList) searchGreaterOrEqual(key interface{}) *node {
	x := l.head
	// 从最高的一层开始
	level := l.maxHeight - 1
	for {
		n := x.next(level)
		if l.keyIsAfterNode(key, n) {
			x = n
			//} else if l.keyIsEqualNode(key, n) {
			//	return n
		} else {
			if level == 0 {
				return n
			} else {
				level--
			}
		}
	}
}

// prev 用于记录插入操作时需要更新的各层 next
func (l *SkipList) findGreaterOrEqual(key interface{}, prev *[kMaxHeight]*node) *node {
	x := l.head
	level := l.maxHeight - 1
	for {
		n := x.next(level)
		if l.keyIsAfterNode(key, n) {
			x = n
		} else {
			// save the path
			prev[level] = x
			if level == 0 {
				return n
			} else {
				level--
			}
		}
	}
}

// return the latest node with a key < key
func (l *SkipList) findLessThan(key interface{}) *node {
	x := l.head
	level := l.maxHeight - 1
	for {
		// x 相当于 n 的 prev
		n := x.next(level)
		if n == nil || l.compare(n.e.Key(), key) >= 0 {
			if level == 0 {
				return x
			} else {
				level--
			}
		} else {
			x = n
		}
	}
}

func (l *SkipList) findLast() *node {
	x := l.head
	level := l.maxHeight - 1
	for {
		n := x.next(level)
		if n == nil {
			if level == 0 {
				return x
			} else {
				level--
			}
		} else {
			x = n
		}
	}
}

type Iterator struct {
	list *SkipList
	node *node
}

func NewIterator(l *SkipList) *Iterator {
	return &Iterator{
		list: l,
		node: nil,
	}
}

func (i *Iterator) Valid() bool {
	return i.node != nil
}

func (i *Iterator) Element() Element {
	return i.node.e
}

func (i *Iterator) Next() {
	// 第0层是包含所有元素的排好序的单向链表
	i.node = i.node.next(0)
}

func (i *Iterator) Prev() {
	i.node = i.list.findLessThan(i.node.e.Key())
	if i.node == i.list.head {
		i.node = nil
	}
}

func (i *Iterator) Seek(key interface{}) {
	i.node = i.list.searchGreaterOrEqual(key)
}

func (i *Iterator) SeekToFirst() {
	i.node = i.list.head.next(0)
}

func (i *Iterator) SeekToLast() {
	i.node = i.list.findLast()
	if i.node == i.list.head {
		i.node = nil
	}
}
