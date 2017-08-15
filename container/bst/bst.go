package bst

type KeyCompareFunc func(leftKey interface{}, rightKey interface{}) bool

type KeyFunc func(value interface{}) interface{}

type TraverseFunc func(key interface{}, value interface{})

type BST interface {
	InsertUnique(value interface{}) bool
	InsertEqual(value interface{})
	Remove(key interface{})
	Have(key interface{}) bool
	Find(key interface{}) interface{}
	Min() interface{}
	Max() interface{}
	Size() int
	InorderTraverse(TraverseFunc)
}
