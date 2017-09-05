// port from https://github.com/google/leveldb/blob/master/db/skiplist_test.cc

package skiplist

import (
	"math/rand"
	"sort"
	"testing"
	"time"
)

type element int

func (e element) Key() interface{}   { return int(e) }
func (e element) Value() interface{} { return int(e) }

func intComparator(left interface{}, right interface{}) int {
	return left.(int) - right.(int)
}

func TestEmpty(t *testing.T) {
	l := NewSkipList(intComparator)

	if l.Contains(10) {
		t.Errorf("SkipList expect 10 not exist")
	}

	it := NewIterator(l)
	if it.Valid() {
		t.Errorf("SkipList Iterator expect invalid")
	}
	it.SeekToFirst()
	if it.Valid() {
		t.Errorf("SkipList Iterator expect invalid")
	}
	it.Seek(100)
	if it.Valid() {
		t.Errorf("SkipList Iterator expect invalid")
	}
	it.SeekToLast()
	if it.Valid() {
		t.Errorf("SkipList Iterator expect invalid")
	}
}

func testLookup(t *testing.T, keyM map[int]struct{}, list *SkipList) bool {
	keys := make([]int, 0, len(keyM))
	for k := range keyM {
		keys = append(keys, k)
	}
	sort.Sort(sort.IntSlice(keys))

	for _, k := range keys {
		if !list.Contains(k) {
			t.Errorf("SkipList.Contains(%d) expect:true, actual:false", k)
			return false
		}
		it := list.Find(k)
		if it == nil {
			t.Errorf("SkipList.Find(%d) expect: %d, got: nil", k, k)
			return false
		} else if it.Element().Key().(int) != k {
			t.Errorf("SkipList.Find(%d) expect: %d, got: %d", k, k, it.Element().Key().(int))
			return false
		}
	}

	minKey := keys[0]
	maxKey := keys[len(keys)-1]

	// simple iterator tests
	var it *Iterator
	it = NewIterator(list)
	it.Seek(0)
	if !it.Valid() {
		t.Errorf("Iterator.Seek(0) expect valid")
		return false
	}
	if it.Element().Key().(int) != minKey {
		t.Errorf("Iterator.Seek(0) expect:%d, got:%d", minKey, it.Element().Key().(int))
		return false
	}
	it.SeekToFirst()
	if !it.Valid() {
		t.Errorf("Iterator.SeekToFirst() expect valid")
		return false
	}
	if it.Element().Key().(int) != minKey {
		t.Errorf("Iterator.SeekToFirst() expect:%d, got:%d", minKey, it.Element().Key().(int))
		return false
	}
	it.SeekToLast()
	if !it.Valid() {
		t.Errorf("Iterator.SeekToLiast() expect valid")
		return false
	}
	if it.Element().Key().(int) != maxKey {
		t.Errorf("Iterator.SeekToLast() expect:%d, got:%d", maxKey, it.Element().Key().(int))
		return false
	}

	// forward iteration test
	it = NewIterator(list)
	it.SeekToFirst()
	for _, k := range keys {
		if !it.Valid() {
			t.Errorf("Forward iteration expect: %d, got: invalid", k)
			return false
		}
		if it.Element().Key().(int) != k {
			t.Errorf("Forward iteration expect: %d, got: %d", k, it.Element().Key().(int))
			return false
		}
		it.Next()
	}
	if it.Valid() {
		t.Errorf("Forward iteration expect iterator invalid, got: %d", it.Element().Key().(int))
		return false
	}
	t.Logf("Forward iteration checked %d items", len(keys))

	// backward iteration test
	it = NewIterator(list)
	it.SeekToLast()
	for i := len(keys) - 1; i >= 0; i-- {
		if !it.Valid() {
			t.Errorf("Backward iteration expect: %d, got: invalid", keys[i])
			return false
		}
		if it.Element().Key().(int) != keys[i] {
			t.Errorf("Backward iteration expect: %d, got: %d", keys[i], it.Element().Key().(int))
			return false
		}
		it.Prev()
	}
	if it.Valid() {
		t.Errorf("Backward iteration expect iterator invalid, got: %d", it.Element().Key().(int))
		return false
	}
	t.Logf("Backward iteration checked %d items", len(keys))
	return true
}

func TestInsertRemoveAndLookup(t *testing.T) {
	var (
		n int = 2000
		r int = 5000
	)
	keyM := make(map[int]struct{})
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	list := NewSkipList(intComparator)

	// prepare data
	for i := 0; i < n; i++ {
		k := rand.Int() % r
		if _, ok := keyM[k]; !ok {
			list.Insert(element(k))
			keyM[k] = struct{}{}
		}
	}
	if !testLookup(t, keyM, list) {
		t.Errorf("TestInsert failed")
	}

	// remove some random items
	var (
		count int = len(keyM) / 2
		i     int
	)
	for k := range keyM {
		if i == count {
			break
		}
		i++
		delete(keyM, k)
		list.Remove(k)
		if !testLookup(t, keyM, list) {
			t.Errorf("TestRemove failed, removed: %d, key: %d", i, k)
		}
	}
}
