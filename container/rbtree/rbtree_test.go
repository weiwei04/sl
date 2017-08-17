package rbtree

import (
	"reflect"
	"sort"
	"testing"
)

type element int

func (e element) Key() interface{} { return int(e) }

func (e element) Value() interface{} { return int(e) }

func compareFunc(left interface{}, right interface{}) int {
	lv := left.(int)
	rv := right.(int)
	if lv < rv {
		return Less
	} else if lv > rv {
		return Greate
	}
	return Equal
}

var insertCases = []struct {
	Name             string
	Input            []int
	Min              []int
	Max              []int
	Ok               []bool
	ExpectExist      int
	ExpectNotExist   int
	Count            int
	DistinctCount    int
	Traverse         []int
	DistinctTraverse []int
}{
	{
		Name:             "1Item",
		Input:            []int{1},
		Min:              []int{1},
		Max:              []int{1},
		Ok:               []bool{true},
		ExpectExist:      1,
		ExpectNotExist:   0,
		Count:            1,
		DistinctCount:    1,
		Traverse:         []int{1},
		DistinctTraverse: []int{1},
	},
	{
		Name:             "OrderedItems",
		Input:            []int{1, 2, 3, 4, 5, 6},
		Min:              []int{1, 1, 1, 1, 1, 1},
		Max:              []int{1, 2, 3, 4, 5, 6},
		Ok:               []bool{true, true, true, true, true, true},
		ExpectExist:      2,
		ExpectNotExist:   99,
		Count:            6,
		DistinctCount:    6,
		Traverse:         []int{1, 2, 3, 4, 5, 6},
		DistinctTraverse: []int{1, 2, 3, 4, 5, 6},
	},
	{
		Name:             "ReverseOrderedItems",
		Input:            []int{6, 5, 4, 3, 2, 1},
		Min:              []int{6, 5, 4, 3, 2, 1},
		Max:              []int{6, 6, 6, 6, 6, 6},
		Ok:               []bool{true, true, true, true, true, true},
		ExpectExist:      2,
		ExpectNotExist:   100,
		Count:            6,
		DistinctCount:    6,
		Traverse:         []int{1, 2, 3, 4, 5, 6},
		DistinctTraverse: []int{1, 2, 3, 4, 5, 6},
	},
	{
		Name:             "UnorderedItems",
		Input:            []int{3, 5, 1, 6, 4, 2},
		Min:              []int{3, 3, 1, 1, 1, 1},
		Max:              []int{3, 5, 5, 6, 6, 6},
		Ok:               []bool{true, true, true, true, true, true},
		ExpectExist:      2,
		ExpectNotExist:   17,
		Count:            6,
		DistinctCount:    6,
		Traverse:         []int{1, 2, 3, 4, 5, 6},
		DistinctTraverse: []int{1, 2, 3, 4, 5, 6},
	},
	{
		Name:             "OrderedDuplicateItems",
		Input:            []int{1, 2, 3, 3, 3, 4, 5, 5, 5, 6},
		Min:              []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		Max:              []int{1, 2, 3, 3, 3, 4, 5, 5, 5, 6},
		Ok:               []bool{true, true, true, false, false, true, true, false, false, true},
		ExpectExist:      2,
		ExpectNotExist:   99,
		Count:            10,
		DistinctCount:    6,
		Traverse:         []int{1, 2, 3, 3, 3, 4, 5, 5, 5, 6},
		DistinctTraverse: []int{1, 2, 3, 4, 5, 6},
	},
	{
		Name:             "UnorderedDuplicateItems",
		Input:            []int{3, 5, 1, 5, 6, 3, 4, 6, 2, 1},
		Min:              []int{3, 3, 1, 1, 1, 1, 1, 1, 1, 1},
		Max:              []int{3, 5, 5, 5, 6, 6, 6, 6, 6, 6},
		Ok:               []bool{true, true, true, false, true, false, true, false, true, false},
		ExpectExist:      2,
		ExpectNotExist:   17,
		Count:            10,
		DistinctCount:    6,
		Traverse:         []int{1, 1, 2, 3, 3, 4, 5, 5, 6, 6},
		DistinctTraverse: []int{1, 2, 3, 4, 5, 6},
	},
}

func TestRBTree_EmptyTree(t *testing.T) {
	tree := New(compareFunc)
	if tree.Size() != 0 {
		t.Error("tree.Size expect: 0, got:", tree.Size())
	}
	if it := tree.First(); it != nil {
		t.Error("tree.First expect: nil, got: not nil")
	}
	if it := tree.Last(); it != nil {
		t.Error("tree.Last expect: nil, got: not nil")
	}
}

func TestRBTree_InsertUnqiue(t *testing.T) {
	for _, c := range insertCases {
		tree := New(compareFunc)
		if tree.Size() != 0 {
			t.Error("case:", c.Name, ", expect size:", 0, ", got:", tree.Size())
		}
		for i := range c.Input {
			if ok := tree.InsertUnique(element(c.Input[i])); ok != c.Ok[i] {
				t.Error("case:", c.Name, ", failed insertUnique, expect:", c.Input[i], ", got:", ok)
			}
			if tree.Min().Value().(int) != c.Min[i] {
				t.Error("case:", c.Name, ", expect min:", c.Min[i], ", got", tree.Min().Value().(int))
			}
			if tree.Max().Value().(int) != c.Max[i] {
				t.Error("case:", c.Name, "expect max:", c.Max[i], ", got", tree.Min().Value().(int))
			}
		}
		if tree.Size() != c.DistinctCount {
			t.Errorf("case:", c.Name, ", tree.Size expect:", c.DistinctCount, ", got:", tree.Size())
		}
		if tree.Find(c.ExpectExist) == nil {
			t.Errorf("case:", c.Name, "expect exist:", c.ExpectExist, ", got: 404")
		}
		if tree.Find(c.ExpectNotExist) != nil {
			t.Errorf("case:", c.Name, "expect not exist:", c.ExpectExist, ", got:",
				tree.Find(c.ExpectNotExist).Value().(int))
		}
	}
}

func TestRBTree_InsertEqual(t *testing.T) {
	for _, c := range insertCases {
		tree := New(compareFunc)
		if tree.Size() != 0 {
			t.Error("case:", c.Name, ", expect size:", 0, ", got:", tree.Size())
		}
		for i := range c.Input {
			tree.InsertEqual(element(c.Input[i]))
			if tree.Min().Value().(int) != c.Min[i] {
				t.Error("case:", c.Name, ", expect min:", c.Min[i], ", got", tree.Min().Value().(int))
			}
			if tree.Max().Value().(int) != c.Max[i] {
				t.Error("case:", c.Name, "expect max:", c.Max[i], ", got", tree.Min().Value().(int))
			}
		}
		if tree.Size() != len(c.Input) {
			t.Errorf("case:", c.Name, ", tree.Size expect:", c.DistinctCount, ", got:", tree.Size())
		}
		if tree.Find(c.ExpectExist) == nil {
			t.Error("case:", c.Name, "expect exist:", c.ExpectExist, ", got: 404")
		}
		if tree.Find(c.ExpectNotExist) != nil {
			t.Errorf("case:", c.Name, "expect not exist:", c.ExpectExist, ", got:",
				tree.Find(c.ExpectNotExist).Value().(int))
		}
	}
}

func TestRBTree_IterateUnique(t *testing.T) {
	for _, c := range insertCases {
		var (
			tree   *RBTree
			result []int
		)

		tree = New(compareFunc)
		for i := range c.Input {
			tree.InsertUnique(element(c.Input[i]))
		}
		result = []int{}
		for it := tree.First(); it != nil; it = it.Next() {
			result = append(result, it.Value().(int))
		}
		if !reflect.DeepEqual(c.DistinctTraverse, result) {
			t.Error("case:", c.Name, ", iterate seqs")
			t.Error("\texpect:", c.DistinctTraverse)
			t.Error("\tgot:", result)
		}

		tree = New(compareFunc)
		for i := range c.Input {
			tree.InsertUnique(element(c.Input[i]))
		}
		result = []int{}
		for it := tree.Last(); it != nil; it = it.Prev() {
			result = append(result, it.Value().(int))
		}
		sort.Sort(sort.Reverse(sort.IntSlice(c.DistinctTraverse)))
		if !reflect.DeepEqual(c.DistinctTraverse, result) {
			t.Error("case:", c.Name, ", iterate seqs")
			t.Error("\texpect:", c.DistinctTraverse)
			t.Error("\tgot:", result)
		}
	}
}

func TestRBTree_IterateDuplicate(t *testing.T) {
	for _, c := range insertCases {
		var (
			tree   *RBTree
			result []int
		)

		tree = New(compareFunc)
		for i := range c.Input {
			tree.InsertEqual(element(c.Input[i]))
		}
		result = []int{}
		for it := tree.First(); it != nil; it = it.Next() {
			result = append(result, it.Value().(int))
		}
		if !reflect.DeepEqual(c.Traverse, result) {
			t.Error("case:", c.Name, ", iterate seqs")
			t.Error("\texpect:", c.Traverse)
			t.Error("\tgot:", result)
		}

		tree = New(compareFunc)
		for i := range c.Input {
			tree.InsertEqual(element(c.Input[i]))
		}
		result = []int{}
		for it := tree.Last(); it != nil; it = it.Prev() {
			result = append(result, it.Value().(int))
		}
		sort.Sort(sort.Reverse(sort.IntSlice(c.Traverse)))
		if !reflect.DeepEqual(c.Traverse, result) {
			t.Error("case:", c.Name, ", iterate seqs")
			t.Error("\texpect:", c.Traverse)
			t.Error("\tgot:", result)
		}
	}
}

var removeCases = []struct {
	Name             string
	Input            []int
	RemoveExist      int
	RemoveNotExist   int
	Traverse         []int
	DistinctTraverse []int
}{
	{
		Name:             "1Item",
		Input:            []int{1},
		RemoveExist:      1,
		RemoveNotExist:   2,
		Traverse:         []int{},
		DistinctTraverse: []int{},
	},
	{
		Name:             "Remove",
		Input:            []int{3, 1, 2, 6, 5, 4},
		RemoveExist:      3,
		RemoveNotExist:   100,
		Traverse:         []int{1, 2, 4, 5, 6},
		DistinctTraverse: []int{1, 2, 4, 5, 6},
	},
	{
		Name:             "OrderedRemoveFirst",
		Input:            []int{1, 2, 3, 4, 5, 6},
		RemoveExist:      1,
		RemoveNotExist:   100,
		Traverse:         []int{2, 3, 4, 5, 6},
		DistinctTraverse: []int{2, 3, 4, 5, 6},
	},
	{
		Name:             "ReverseOrderedRemoveFirst",
		Input:            []int{6, 5, 4, 3, 2, 1},
		RemoveExist:      6,
		RemoveNotExist:   100,
		Traverse:         []int{1, 2, 3, 4, 5},
		DistinctTraverse: []int{1, 2, 3, 4, 5},
	},
	{
		Name:             "RemoveMin",
		Input:            []int{5, 2, 3, 1, 4, 6},
		RemoveExist:      1,
		RemoveNotExist:   100,
		Traverse:         []int{2, 3, 4, 5, 6},
		DistinctTraverse: []int{2, 3, 4, 5, 6},
	},
	{
		Name:             "RemoveMax",
		Input:            []int{1, 5, 3, 6, 2, 4},
		RemoveExist:      6,
		RemoveNotExist:   100,
		Traverse:         []int{1, 2, 3, 4, 5},
		DistinctTraverse: []int{1, 2, 3, 4, 5},
	},
	{
		Name:             "RemoveDuplicate",
		Input:            []int{5, 6, 1, 3, 5, 4, 4, 2, 3},
		RemoveExist:      4,
		RemoveNotExist:   100,
		Traverse:         []int{1, 2, 3, 3, 4, 5, 5, 6},
		DistinctTraverse: []int{1, 2, 3, 5, 6},
	},
	{
		Name:             "RemoveUnqiue",
		Input:            []int{5, 6, 1, 3, 5, 4, 4, 2, 3},
		RemoveExist:      1,
		RemoveNotExist:   100,
		Traverse:         []int{2, 3, 3, 4, 4, 5, 5, 6},
		DistinctTraverse: []int{2, 3, 4, 5, 6},
	},
}

func fillUniqueTree(input []int) *RBTree {
	t := New(compareFunc)
	for _, i := range input {
		t.InsertUnique(element(i))
	}
	return t
}

func removeByKey(t *RBTree, key interface{}) {
	t.RemoveByKey(key)
}

func removeByIterator(t *RBTree, key interface{}) {
	if it := t.Find(key); it != nil {
		t.RemoveByIterator(it)
	}
}

func visitFunc(result *[]int) func(e Element) {
	return func(e Element) {
		*result = append(*result, e.Key().(int))
	}
}

func TestRBTree_UniqueRemove(t *testing.T) {
	for _, c := range removeCases {
		for _, fn := range []struct {
			Name string
			Func func(t *RBTree, key interface{})
		}{
			{
				Name: "RemoveByKey",
				Func: removeByKey,
			},
			{
				Name: "RemoveByIterator",
				Func: removeByIterator,
			},
		} {
			t.Run(fn.Name, func(t *testing.T) {
				tree := fillUniqueTree(c.Input)
				fn.Func(tree, c.RemoveExist)
				result := []int{}
				tree.InorderTraverse(visitFunc(&result))
				if !reflect.DeepEqual(c.DistinctTraverse, result) {
					t.Error("case:", c.Name, ", removeFunc:", fn.Name, ", after remove")
					t.Error("\texpect:", c.DistinctTraverse)
					t.Error("\tgot:", result)
				}
				fn.Func(tree, c.RemoveNotExist)
				result = []int{}
				tree.InorderTraverse(visitFunc(&result))
				if !reflect.DeepEqual(c.DistinctTraverse, result) {
					t.Error("case:", c.Name, ", removeFunc:", fn.Name, ", after remove not exist")
					t.Error("\texpect:", c.DistinctTraverse)
					t.Error("\tgot:", result)
				}
			})
		}
	}
}

func fillEqualTree(input []int) *RBTree {
	t := New(compareFunc)
	for _, i := range input {
		t.InsertEqual(element(i))
	}
	return t
}

func TestRBTree_EqualRemove(t *testing.T) {
	for _, c := range removeCases {
		for _, fn := range []struct {
			Name string
			Func func(t *RBTree, key interface{})
		}{
			{
				Name: "RemoveByKey",
				Func: removeByKey,
			},
			{
				Name: "RemoveByIterator",
				Func: removeByIterator,
			},
		} {
			t.Run(fn.Name, func(t *testing.T) {
				tree := fillEqualTree(c.Input)
				fn.Func(tree, c.RemoveExist)
				result := []int{}
				tree.InorderTraverse(visitFunc(&result))
				if !reflect.DeepEqual(c.Traverse, result) {
					t.Error("case:", c.Name, ", removeFunc:", fn.Name, ", after remove")
					t.Error("\texpect:", c.Traverse)
					t.Error("\tgot:", result)
				}
				fn.Func(tree, c.RemoveNotExist)
				result = []int{}
				tree.InorderTraverse(visitFunc(&result))
				if !reflect.DeepEqual(c.Traverse, result) {
					t.Error("case:", c.Name, ", removeFunc:", fn.Name, ", after remove not exist")
					t.Error("\texpect:", c.Traverse)
					t.Error("\tgot:", result)
				}
			})
		}
	}
}
