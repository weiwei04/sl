package stack

type element struct {
	value interface{}
}

type sliceStack struct {
	s []element
}

func NewSliceStack() Stack {
	return &sliceStack{}
}

func (s *sliceStack) Push(value interface{}) {
	s.s = append(s.s, element{value: value})
}

func (s *sliceStack) Pop() interface{} {
	value := s.s[len(s.s)-1].value
	s.s = s.s[0 : len(s.s)-1]
	return value
}

func (s *sliceStack) Top() interface{} {
	return s.s[len(s.s)-1].value
}

func (s *sliceStack) Size() int { return len(s.s) }
