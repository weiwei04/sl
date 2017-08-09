package queue

type node struct {
	value interface{}
	next  *node
}

type slistQueue struct {
	head  node
	tail  *node
	count int
}

func NewSListQueue() Queue {
	q := &slistQueue{}
	q.tail = &q.head
	return q
}

func (q *slistQueue) EnQueue(value interface{}) {
	q.tail.next = &node{
		value: value,
	}
	q.tail = q.tail.next
	q.count++
}

func (q *slistQueue) DeQueue() interface{} {
	p := q.head.next
	q.head.next = p.next
	q.count--
	if q.count == 0 {
		q.tail = &q.head
	}
	return p.value
}

func (q *slistQueue) Front() interface{} {
	return q.head.next.value
}

func (q *slistQueue) Size() int {
	return q.count
}
