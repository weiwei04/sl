package queue

import (
	"reflect"
	"testing"
)

func TestQueue(t *testing.T) {
	cases := []struct {
		Name    string
		EnQSeqs []int
		DeQSeqs []int
	}{
		{
			Name:    "Empty",
			EnQSeqs: []int{},
			DeQSeqs: []int{},
		},
		{
			Name:    "One Item",
			EnQSeqs: []int{1},
			DeQSeqs: []int{1},
		},
		{
			Name:    "Many Items",
			EnQSeqs: []int{1, 2, 3, 4, 5},
			DeQSeqs: []int{1, 2, 3, 4, 5},
		},
	}

	for _, impl := range []struct {
		Name string
		Func func() Queue
	}{
		{
			Name: "Slice Queue",
			Func: NewSliceQueue,
		},
		{
			Name: "SList Queue",
			Func: NewSListQueue,
		},
	} {
		t.Run("batch: "+impl.Name, func(t *testing.T) {
			for _, c := range cases {
				queue := impl.Func()
				if queue.Size() != 0 {
					t.Error(impl.Name, "case", c.Name, "expect queue.Size = 0", "got =", queue.Size())
				}
				for i, op := range c.EnQSeqs {
					if queue.Size() != i {
						t.Error("case:", c.Name, ", queue:", impl.Name)
						t.Error("\texpect queue.Size() =", i, ", got", queue.Size())
					}
					queue.EnQueue(op)
				}
				deQSeqs := []int{}
				for queue.Size() != 0 {
					deQSeqs = append(deQSeqs, queue.Front().(int))
					queue.DeQueue()
				}
				if !reflect.DeepEqual(c.DeQSeqs, deQSeqs) {
					t.Error("case:", c.Name, ", stack:", impl.Name)
					t.Error("\texpect pop seq:", c.DeQSeqs, "get:", deQSeqs)
				}
			}
		})
		t.Run("1by1: "+impl.Name, func(t *testing.T) {
			for _, c := range cases {
				queue := impl.Func()
				if queue.Size() != 0 {
					t.Error(impl.Name, "case", c.Name, "expect stack.Size = 0", "got =", queue.Size())
				}
				deQSeqs := []int{}
				for _, op := range c.EnQSeqs {
					queue.EnQueue(op)
					if queue.Front().(int) != op {
						t.Error("case", c.Name, ", queue:", impl.Name)
						t.Error("\texpect queue.Front:", op, ", got", queue.Front().(int))
					}
					deQSeqs = append(deQSeqs, queue.DeQueue().(int))
					if queue.Size() != 0 {
						t.Error("case", c.Name, ", queue:", impl.Name)
						t.Error("\texpect queue.Size():", 0, "got:", queue.Size())
					}
				}
			}
		})
	}
}
