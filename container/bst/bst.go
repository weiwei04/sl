package bst

type CompareResultType int

const (
	Less   CompareResultType = -1
	Equal  CompareResultType = 0
	Greate CompareResultType = 1
)

type CompareFunc func(left interface{}, right interface{}) CompareResultType

type TraverseFunc func(key interface{}, value interface{})

type BST interface {
	Insert(key interface{}, value interface{})
	Remove(key interface{})
	Have(key interface{}) bool
	Find(key interface{}) interface{}
	Min() interface{}
	Max() interface{}
	Size() int
	InorderTraverse(TraverseFunc)
}
