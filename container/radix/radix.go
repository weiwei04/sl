package radix

import (
	"sort"
	"strings"
)

type leafNode struct {
	key string
	val interface{}
}

type edge struct {
	// label 用于检查是否需要向这个分支继续查询
	// 存储了 node.key[0] 的内容
	label byte
	node  *node
}

type edges []edge

func (e edges) Len() int {
	return len(e)
}

func (e edges) Less(i, j int) bool {
	return e[i].label < e[j].label
}

func (e edges) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

type node struct {
	leaf   *leafNode
	prefix string
	edges  edges
}

func (n *node) isLeaf() bool {
	return n.leaf != nil
}

func (n *node) addEdge(e edge) {
	n.edges = append(n.edges, e)
	sort.Sort(n.edges)
}

func (n *node) mergeChild() {
	child := n.edges[0].node
	n.edges = child.edges
	n.prefix = n.prefix + child.prefix
	n.leaf = child.leaf
}

func (n *node) unguardedDelEdge(label byte) {
	i := 0
	for ; n.edges[i].label != label; i++ {
	}
	if i < len(n.edges)-1 {
		copy(n.edges[i:], n.edges[i+1:])
	}
	// gc can free this node
	n.edges[len(n.edges)-1] = edge{}
	n.edges = n.edges[:len(n.edges)-1]
}

func (n *node) replaceEdge(e edge) {
	// TODO: 可以使用二分查找
	for i := range n.edges {
		if n.edges[i].label == e.label {
			n.edges[i] = e
			return
		}
	}
}

func (n *node) getEdge(label byte) *node {
	// TODO: edges 是排好序的，可以使用二分查找
	for _, e := range n.edges {
		if e.label == label {
			return e.node
		} else if e.label > label {
			break
		}
	}
	return nil
}

type RadixTree struct {
	root *node // root = &node{} contains a empty k:v pair
	size int
}

func NewRadixTree() *RadixTree {
	return &RadixTree{
		root: &node{},
		size: 0,
	}
}

func NewRadixTreeFromMap(m map[string]interface{}) *RadixTree {
	t := &RadixTree{root: &node{}, size: 0}
	for k, v := range m {
		t.Insert(k, v)
	}
	return t
}

func (r *RadixTree) Size() int {
	return r.size
}

func longestPrefix(a, b string) int {
	max := len(a)
	if len(b) < max {
		max = len(b)
	}
	var i int
	for ; i < max; i++ {
		if a[i] != b[i] {
			break
		}
	}
	return i
}

//func (t *RadixTree) Insert(k string, v interface{}) {
//	var parent *node
//	n := t.root
//	search := k
//	leaf := &leafNode{k, v}
//	for {
//		commonPrefix := longestCommonPrefix(n.prefix, search)
//		if commonPrefix == n.prefix {
//			if commonPrefix == len(search) {
//				if !n.isLeaf() {
//					t.size++
//				}
//				n.leaf = leaf
//				return
//			} else {
//				// commonPrefix < len(search)
//				search = search[commonPrefix:]
//				child := n.getEdge(search[0])
//				if child == nil {
//					// just insert one
//					n.addEdge(edge{
//						label: search[0],
//						node: &node{
//							leaf:   leaf,
//							prefix: search,
//						},
//					})
//					return
//				} else {
//					parent = n
//					n = child
//				}
//			}
//		} else {
//			// commonPrefix < n.prefix
//			// create a new child between parent -> node
//			t.size++
//			child := &node{
//				prefix: search[0:commonPrefix],
//			}
//			// add new child between parent->node(do the split)
//			parent.replaceEdge(edge{
//				label: search[0],
//				node:  child,
//			})
//			// link child -> node
//			n.prefix = n.prefix[commonPrefix:]
//			child.addEdge(edge{
//				label: n.prefix[0],
//				node:  n,
//			})
//			// add the new inserted node
//			if commonPrefix == len(search) {
//				child.leaf = leaf
//			} else {
//				search = search[commonPrefix:]
//				child.addEdge(edge{
//					label: search[0],
//					node: &node{
//						prefix: search,
//						leaf:   leaf,
//					},
//				})
//			}
//		}
//	}
//}

// upsert
func (t *RadixTree) Insert(k string, v interface{}) (interface{}, bool) {
	var parent *node
	n := t.root
	search := k
	for {
		if len(search) == 0 {
			if n.isLeaf() {
				old := n.leaf.val
				n.leaf.val = v
				return old, true
			}

			n.leaf = &leafNode{k, v}
			t.size++
			return nil, false
		}

		parent = n
		n = n.getEdge(search[0])
		if n == nil {
			e := edge{
				label: search[0],
				node: &node{
					leaf:   &leafNode{k, v},
					prefix: search,
				},
			}
			parent.addEdge(e)
			t.size++
			return nil, false
		}

		commonPrefix := longestPrefix(search, n.prefix)
		if commonPrefix == len(n.prefix) {
			search = search[commonPrefix:]
			continue
		}
		// commonPrefix < len(n.prefix)
		// 在 parent 和 n 之间插入一个新的节点 child
		// parent -> child -> {n, newNode(k,v)}
		// 1. 如果 commonPrefix == len(search)
		//     parent -> {child.leaf(k, v)} & child -> n
		// 2. 如果 commonPrefix < len(search)
		//     parent -> child -> {n, newNode.leaf(k, v)}
		t.size++
		child := &node{
			prefix: search[:commonPrefix],
		}
		parent.replaceEdge(edge{
			label: search[0],
			node:  child,
		})

		n.prefix = n.prefix[commonPrefix:]
		child.addEdge(edge{
			label: n.prefix[0],
			node:  n,
		})

		// create leafNode for k, v
		// parent -> {child,}
		if commonPrefix == len(search) {
			child.leaf = &leafNode{k, v}
		} else {
			child.addEdge(edge{
				label: search[commonPrefix],
				node: &node{
					leaf:   &leafNode{k, v},
					prefix: search[commonPrefix:],
				},
			})
		}
		return nil, false
	}
}

func (t *RadixTree) Delete(k string) {
	var label byte
	var parent *node
	n := t.root
	search := k
	for {
		if len(search) == 0 {
			if !n.isLeaf() {
				// no such key
				return
			}
			n.leaf = nil
			t.size--
			if parent != nil && len(n.edges) == 0 {
				parent.unguardedDelEdge(label)
				if parent != t.root && len(parent.edges) == 1 && !parent.isLeaf() {
					// 检查 parent 删掉 n 之后是否需要和独子 merge
					parent.mergeChild()
				}
			}
			if n != t.root && len(n.edges) == 1 {
				n.mergeChild()
			}
			return
		}

		label = search[0]
		parent = n
		n = n.getEdge(label)
		if n == nil {
			return
		}

		if strings.HasPrefix(search, n.prefix) {
			search = search[len(n.prefix):]
		} else {
			return
		}
	}
}

func (t *RadixTree) Contains(k string) bool {
	_, ok := t.Get(k)
	return ok
}

func (t *RadixTree) Get(k string) (interface{}, bool) {
	n := t.root
	search := k
	for {
		if len(search) == 0 {
			if n.isLeaf() {
				return n.leaf.val, true
			}
			break
		}

		// 从 n 的子节点中寻找是否存在满足条件的子节点
		n = n.getEdge(search[0])
		if n == nil {
			break
		}
		if strings.HasPrefix(search, n.prefix) {
			search = search[len(n.prefix):]
		} else {
			break
		}
	}
	return nil, false
}
