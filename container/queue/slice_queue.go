package queue

type element struct {
	value interface{}
}

type sliceQueue struct {
	s []element
}

func NewSliceQueue() Queue {
	return &sliceQueue{}
}

func (q *sliceQueue) EnQueue(value interface{}) {
	q.s = append(q.s, element{value})
}

func (q *sliceQueue) DeQueue() interface{} {
	value := q.s[0].value
	q.s = q.s[1:len(q.s)]
	return value
}

func (q *sliceQueue) Front() interface{} {
	return q.s[0].value
}

func (q *sliceQueue) Size() int {
	return len(q.s)
}
