package bst

import (
	"reflect"
	"testing"
)

func TestBST(t *testing.T) {
	cases := []struct {
		Name   string
		Input  []int
		Expect []int
		Min    int
		Max    int
	}{
		{
			Name:   "Empty",
			Input:  []int{},
			Expect: []int{},
		},
		{
			Name:   "1Item",
			Input:  []int{1},
			Expect: []int{1},
			Min:    1,
			Max:    1,
		},
		{
			Name:   "OrderedItems",
			Input:  []int{1, 2, 3, 4, 5},
			Expect: []int{1, 2, 3, 4, 5},
			Min:    1,
			Max:    5,
		},
		{
			Name:   "UnorderedItems",
			Input:  []int{5, 4, 2, 3, 1},
			Expect: []int{1, 2, 3, 4, 5},
			Min:    1,
			Max:    5,
		},
		{
			Name:   "DuplicateItems",
			Input:  []int{5, 4, 3, 2, 3, 1},
			Expect: []int{1, 2, 3, 3, 4, 5},
			Min:    1,
			Max:    5,
		},
	}

	for _, impl := range []struct {
		Name string
		Func func(CompareFunc) BST
	}{
		{
			Name: "SimpleBST",
			Func: NewSimpleBST,
		},
	} {
		t.Run(impl.Name, func(t *testing.T) {
			for _, c := range cases {
				bst := impl.Func(func(left interface{}, right interface{}) CompareResultType {
					if left.(int) < right.(int) {
						return Less
					} else if left.(int) == right.(int) {
						return Equal
					} else {
						return Greate
					}
				})
				var count int
				for _, i := range c.Input {
					bst.Insert(i, i)
					count++
					if count != bst.Size() {
						t.Error("case:", c.Name, ", BST:", impl.Name)
						t.Error("expect size:", count, ", got:", bst.Size())
					}
				}
				if len(c.Input) != 0 {
					if c.Min != bst.Min().(int) {
						t.Error("case:", c.Name, ", BST:", impl.Name)
						t.Error("expect min:", c.Min, ", got:", bst.Min().(int))
					}
					if c.Max != bst.Max().(int) {
						t.Error("case:", c.Name, ", BST:", impl.Name)
						t.Error("expect max:", c.Max, ", got:", bst.Max().(int))
					}
				}
				output := []int{}
				bst.InorderTraverse(func(k interface{}, v interface{}) {
					output = append(output, v.(int))
				})
				if !reflect.DeepEqual(c.Expect, output) {
					t.Error("case:", c.Name, ", BST:", impl.Name)
					t.Error("expect max:", c.Expect, ", got:", output)
				}
			}
		})
	}
}
