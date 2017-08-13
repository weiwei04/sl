package bst

type simpleBST struct {
	root    *node
	count   int
	compare CompareFunc
}

func NewSimpleBST(compare CompareFunc) BST {
	return &simpleBST{compare: compare}
}

func (t *simpleBST) Insert(key interface{}, value interface{}) {
	if t.root != nil {
		n := t.root
		for true {
			if t.compare(key, n.key) == Less {
				if n.left == nil {
					n.left = &node{key: key, value: value}
					break
				} else {
					n = n.left
				}
			} else {
				if n.right == nil {
					n.right = &node{key: key, value: value}
					break
				} else {
					n = n.right
				}
			}
		}
	} else {
		t.root = &node{key: key, value: value}
	}
	t.count++
}

func (t *simpleBST) Remove(key interface{}) {
	var n *node
	t.root, n = remove(t.root, t.compare, key)
	if n != nil {
		t.count--
	}
}

func (t *simpleBST) Have(key interface{}) bool {
	return find(t.root, t.compare, key) != nil
}

func (t *simpleBST) Find(key interface{}) interface{} {
	node := find(t.root, t.compare, key)
	if node != nil {
		return node.value
	}
	return nil
}

func (t *simpleBST) Min() interface{} {
	return unguardedMin(t.root).key
}

func (t *simpleBST) Max() interface{} {
	return unguardedMax(t.root).key
}

func (t *simpleBST) Size() int {
	return t.count
}

func (t *simpleBST) InorderTraverse(fn TraverseFunc) {
	inorderTraverse(t.root, fn)
	//recursiveInorderTraverse(t.root, fn)
}
