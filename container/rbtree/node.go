package rbtree

type node struct {
	parent      *node
	left, right *node
	e           Element
	// for new node, the color will always be red
	black bool
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

func isRed(n *node) bool {
	if n == nil {
		return false
	}
	return !n.black
}

func rotateLeft(n *node) *node {
	x := n.right
	n.right = x.left
	if x.left != nil {
		x.left.parent = n
	}
	x.left = n
	x.parent = n.parent
	n.parent = x
	x.black = n.black
	n.black = false
	return x
}

func rotateRight(n *node) *node {
	x := n.left
	n.left = x.right
	if x.right != nil {
		x.right.parent = n
	}
	x.right = n
	x.parent = n.parent
	n.parent = x
	x.black = n.black
	n.black = false
	return x
}

func flipColor(n *node) *node {
	n.left.black = true
	n.right.black = true
	n.black = false
	return n
}
