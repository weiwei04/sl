package stack

import (
	"reflect"
	"testing"
)

func TestStack(t *testing.T) {
	cases := []struct {
		Name     string
		PushSeqs []int
		PopSeqs  []int
	}{
		{
			Name:     "empty",
			PushSeqs: []int{},
			PopSeqs:  []int{},
		},
		{
			Name:     "One Item",
			PushSeqs: []int{1},
			PopSeqs:  []int{1},
		},
		{
			Name:     "Many Items",
			PushSeqs: []int{1, 2, 3, 4, 5},
			PopSeqs:  []int{5, 4, 3, 2, 1},
		},
	}

	for _, impl := range []struct {
		Name string
		Func func() Stack
	}{
		{
			Name: "Slice Stack",
			Func: NewSliceStack,
		},
		{
			Name: "SList Stack",
			Func: NewSListStack,
		},
	} {
		t.Run(impl.Name, func(t *testing.T) {
			for _, c := range cases {
				stack := impl.Func()
				if stack.Size() != 0 {
					t.Error(impl.Name, "case", c.Name, "expect stack.Size = 0", "got =", stack.Size())
				}
				for i, op := range c.PushSeqs {
					if stack.Size() != i {
						t.Error("case:", c.Name, ", stack:", impl.Name)
						t.Error("\texpect stack.Size() =", i, ", got", stack.Size())
					}
					stack.Push(op)
					if op != stack.Top().(int) {
						t.Error("case:", c.Name, ", stack:", impl.Name)
						t.Error("\texpect stack.Top() =", op, ", got", stack.Top().(int))
					}
				}
				popSeqs := []int{}
				for stack.Size() != 0 {
					popSeqs = append(popSeqs, stack.Top().(int))
					stack.Pop()
				}
				if !reflect.DeepEqual(c.PopSeqs, popSeqs) {
					t.Error("case:", c.Name, ", stack:", impl.Name)
					t.Error("\texpect pop seq:", c.PopSeqs, "get:", popSeqs)
				}
			}
		})
	}
}
