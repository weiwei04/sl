package stack

type node struct {
	value interface{}
	next  *node
}

type slistStack struct {
	node
	count int
}

func NewSListStack() Stack {
	return &slistStack{count: 0}
}

func (s *slistStack) Push(value interface{}) {
	node := &node{
		value: value,
		next:  s.next,
	}
	s.next = node
	s.count++
}

func (s *slistStack) Pop() interface{} {
	p := s.next
	s.next = p.next
	s.count--
	return p.value
}

func (s *slistStack) Top() interface{} {
	return s.next.value
}

func (s *slistStack) Size() int {
	return s.count
}
