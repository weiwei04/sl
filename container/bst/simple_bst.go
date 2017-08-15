package bst

type simpleBST struct {
	root       *node
	count      int
	keyOfValue KeyFunc
	less       KeyCompareFunc
}

func NewSimpleBST(keyOfValue KeyFunc, less KeyCompareFunc) BST {
	return &simpleBST{keyOfValue: keyOfValue, less: less}
}

func (t *simpleBST) InsertUnique(value interface{}) bool {
	if t.root != nil {
		n := t.root
		for true {
			if t.less(t.keyOfValue(value), t.keyOfValue(n.value)) {
				if n.left == nil {
					n.left = &node{value: value}
					t.count++
					return true
				} else {
					n = n.left
				}
			} else if t.less(t.keyOfValue(n.value), t.keyOfValue(value)) {
				if n.right == nil {
					n.right = &node{value: value}
					t.count++
					return true
				} else {
					n = n.right
				}
			} else {
				return false
			}
		}
	} else {
		t.root = &node{value: value}
		t.count++
		return true
	}
}

func (t *simpleBST) InsertEqual(value interface{}) {
	if t.root != nil {
		n := t.root
		for true {
			if t.less(t.keyOfValue(value), t.keyOfValue(n.value)) {
				if n.left == nil {
					n.left = &node{value: value}
					break
				} else {
					n = n.left
				}
			} else {
				if n.right == nil {
					n.right = &node{value: value}
					break
				} else {
					n = n.right
				}
			}
		}
	} else {
		t.root = &node{value: value}
	}
	t.count++
}

func (t *simpleBST) Remove(key interface{}) {
	var n *node
	t.root, n = remove(t.root, t.keyOfValue, t.less, key)
	if n != nil {
		t.count--
	}
}

func (t *simpleBST) Have(key interface{}) bool {
	return find(t.root, t.keyOfValue, t.less, key) != nil
}

func (t *simpleBST) Find(key interface{}) interface{} {
	node := find(t.root, t.keyOfValue, t.less, key)
	if node != nil {
		return node.value
	}
	return nil
}

func (t *simpleBST) Min() interface{} {
	return unguardedMin(t.root).value
}

func (t *simpleBST) Max() interface{} {
	return unguardedMax(t.root).value
}

func (t *simpleBST) Size() int {
	return t.count
}

func (t *simpleBST) InorderTraverse(fn TraverseFunc) {
	inorderTraverse(t.root, fn)
	//recursiveInorderTraverse(t.root, fn)
}
