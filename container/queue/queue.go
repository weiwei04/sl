package queue

type Queue interface {
	EnQueue(value interface{})
	DeQueue() interface{}
	Front() interface{}
	Size() int
}
