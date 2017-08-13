package bst

import "github.com/moladb/sl/container/stack"

type node struct {
	key   interface{}
	value interface{}
	left  *node
	right *node
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

// remove the min item in n
// input: root
// output: new-root, n (min)
func removeMin(n *node) (*node, *node) {
	root := n
	var pre *node
	for n.left != nil {
		pre = n
		n = n.left
	}
	if pre == nil {
		root = n.right
		// not nessary
		//n.right = nil
	} else {
		pre.left = n.right
	}
	return root, n
}

func removeMax(n *node) (*node, *node) {
	root := n
	var pre *node
	for n.right != nil {
		pre = n
		n = n.right
	}
	if pre == nil {
		root = n.left
		// n.left = nil
	} else {
		pre.right = n.left
	}
	return root, n
}

func find(n *node, compare CompareFunc, value interface{}) *node {
	for n != nil {
		res := compare(value, n.value)
		if res == Equal {
			return n
		} else if res == Less {
			n = n.left
		} else {
			n = n.right
		}
	}
	return nil
}

// below functions just for fun
func recursiveRemoveMin(n *node) *node {
	if n.left == nil {
		return n.right
	}
	n.left = recursiveRemoveMin(n.left)
	return n
}

func remove(t *node, compare CompareFunc, key interface{}) (*node, *node) {
	n := t
	var parent *node
	leftChild := false
	for n != nil {
		comp := compare(key, n.key)
		if comp == Equal {
			break
		} else if comp == Less {
			parent = n
			n = n.left
		} else {
			parent = n
			n = n.right
		}
	}
	if n == nil {
		return t, nil
	}
	if n.left == nil {
		if parent == nil {
			t = n.right
		} else if leftChild {
			parent.left = n.right
		} else {
			parent.right = n.right
		}
	} else if n.right == nil {
		if parent == nil {
			t = n.left
		} else if leftChild {
			parent.left = n.left
		} else {
			parent.right = n.left
		}
	} else {
		right, x := removeMin(n.right)
		x.right = right
		x.left = n.left
		if parent == nil {
			t = x
		} else if leftChild {
			parent.left = x
		} else {
			parent.right = x
		}
	}
	return t, n
}

func inorderTraverse(n *node, fn TraverseFunc) {
	s := stack.NewSliceStack()
	for s.Size() != 0 || n != nil {
		for n != nil {
			s.Push(n)
			n = n.left
		}
		n = s.Pop().(*node)
		fn(n.key, n.value)
		n = n.right
	}
}

// below functions just for fun

func recursiveRemove(n *node, compare CompareFunc, key interface{}) (*node, *node) {
	var delNode *node
	if n == nil {
		return n, nil
	}
	comp := compare(key, n.key)
	if comp == Less {
		n.left, delNode = recursiveRemove(n.left, compare, key)
	} else if comp == Equal {
		if n.left == nil {
			return n.right, n
		} else if n.right == nil {
			return n.left, n
		} else {
			delNode = n
			right, n := removeMin(n.right)
			n.right = right
			n.left = delNode.left
		}
	} else {
		n.right, delNode = recursiveRemove(n.right, compare, key)
	}
	return n, delNode
}

func recursiveInorderTraverse(n *node, fn TraverseFunc) {
	if n == nil {
		return
	}
	recursiveInorderTraverse(n.left, fn)
	fn(n.key, n.value)
	recursiveInorderTraverse(n.right, fn)
}
