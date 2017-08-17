package rbtree

type node struct {
	parent      *node
	left, right *node
	e           Element
}

func insertUnique(root **node, comp CompareFunc, e Element) (*node, bool) {
	var (
		n      = root
		parent *node
	)
	for *n != nil {
		equal := comp(e.Key(), (*n).e.Key())
		parent = *n
		if equal < 0 {
			n = &((*n).left)
		} else if equal > 0 {
			n = &((*n).right)
		} else {
			return nil, false
		}
	}
	*n = &node{
		parent: parent,
		e:      e,
	}
	// TODO: rebalance
	return *n, true
}

func insertEqual(root **node, comp CompareFunc, e Element) *node {
	var (
		n      = root
		parent *node
	)
	for *n != nil {
		parent = *n
		if comp(e.Key(), (*n).e.Key()) < 0 {
			n = &((*n).left)
		} else {
			n = &((*n).right)
		}
	}
	*n = &node{
		parent: parent,
		e:      e,
	}
	// TODO: rebalance
	return *n
}

func search(n *node, comp CompareFunc, key interface{}) *node {
	for n != nil {
		equal := comp(key, n.e.Key())
		if equal < 0 {
			n = n.left
		} else if equal > 0 {
			n = n.right
		} else {
			return n
		}
	}
	return nil
}

func unguardedMin(n *node) *node {
	for ; n.left != nil; n = n.left {
	}
	return n
}

func unguardedMax(n *node) *node {
	for ; n.right != nil; n = n.right {
	}
	return n
}

func removeMin(root *node) (*node, *node) {
	n := root
	for ; n.left != nil; n = n.left {
	}
	if n == root {
		return n.right, n
	}
	n.parent.left = n.right
	return root, n
}

func remove(root *node, n *node) *node {
	var subTree *node

	if n.left == nil {
		subTree = n.right
	} else if n.right == nil {
		subTree = n.left
	} else {
		var right *node
		right, subTree = removeMin(n.right)
		subTree.right = right
		subTree.left = n.left
	}
	if root == n {
		return subTree
	}
	if n.parent.left == n {
		n.parent.left = subTree
	} else {
		n.parent.right = subTree
	}
	return root
}
