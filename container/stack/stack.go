package stack

type Stack interface {
	Push(val interface{})
	Pop() interface{}
	Top() interface{}
	Size() int
}
