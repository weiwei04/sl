package stack

type Stack interface {
	Push(val interface{})
	Pop()
	Top() interface{}
	Size() int
}
