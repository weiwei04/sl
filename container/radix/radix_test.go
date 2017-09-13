package radix

import (
	crand "crypto/rand"
	"fmt"
	"testing"
)

func TestRadixTree(t *testing.T) {
	rTree, mMap := prepareDataSet()
	if rTree.Size() != len(mMap) {
		t.Fatalf("invalid tree size expect: %d, actual: %d",
			len(mMap), rTree.Size())
	}
	for k, v := range mMap {
		rV, ok := rTree.Get(k)
		if !ok {
			t.Fatalf("invalid get %s result expect: %d, actual: not exist", k, v.(int))
		}
		if rV.(int) != v.(int) {
			t.Fatalf("invalid get result expect: %d, actual: %d", v.(int), rV.(int))
		}
	}
	size := rTree.Size()
	for k, v := range mMap {
		rV, ok := rTree.Get(k)
		if !ok {
			t.Fatalf("invalid get result expect: %d, actual: not exist", v.(int))
		}
		if rV.(int) != v.(int) {
			t.Fatalf("invalid get result expect: %d, actual: %d", v.(int), rV.(int))
		}
		rTree.Delete(k)
		ok = rTree.Contains(k)
		if ok {
			t.Fatalf("invalid get result expect: not exist")
		}
		size--
		if size != rTree.Size() {
			t.Fatalf("invalid tree size expect: %d, actual: %d", size, rTree.Size())
		}
	}
}

func TestRadixTreeInsert(t *testing.T) {
	rTree, mMap := prepareDataSet()
	size := len(mMap)
	for k, v := range mMap {
		old, ok := rTree.Insert(k, v.(int)+1)
		if !ok {
			t.Fatalf("invalid get result expect: %d, actual: not exist", v.(int))
		}
		if old.(int) != v.(int) {
			t.Fatalf("invalid get result expect: %d, actual: %d", v.(int), old.(int))
		}
		if rTree.Size() != size {
			t.Fatalf("invalid tree size expect: %d, actual: %d", size, rTree.Size())
		}
	}
	for k, v := range mMap {
		rV, ok := rTree.Get(k)
		if !ok {
			t.Fatalf("invalid get result expect: %d, actual: not exist", v.(int))
		}
		if rV.(int) != v.(int)+1 {
			t.Fatalf("invalid get result expect: %d, actual: %d", v.(int)+1, rV.(int))
		}
	}
}

func TestRadixTreeRoot(t *testing.T) {
	rTree := NewRadixTree()
	var (
		ok  bool
		val interface{}
	)
	if rTree.Size() != 0 {
		t.Fatalf("invalid tree size expect: 0, actual: %d", rTree.Size())
	}
	ok = rTree.Contains("")
	if ok {
		t.Fatalf("invalid get result expect: nil, actual: not found")
	}
	_, ok = rTree.Insert("", true)
	if ok {
		t.Fatalf("invalid insert result expect false")
	}
	if rTree.Size() != 1 {
		t.Fatalf("invalid tree size expect: 1, actual: %d", rTree.Size())
	}
	val, ok = rTree.Get("")
	if !ok {
		t.Fatalf("invalid get result expect: true, actual: not found")
	}
	if !val.(bool) {
		t.Fatalf("invalid get result expect: true, actual: false")
	}
	rTree.Delete("")
	ok = rTree.Contains("")
	if ok {
		t.Fatalf("invalid get result expect: not exist")
	}
	if rTree.Size() != 0 {
		t.Fatalf("invalid tree size expect: 0, actual: %d", rTree.Size())
	}
}

func TestRadixTreeDelete(t *testing.T) {
	rTree := NewRadixTree()
	keys := []string{"a", "aa", "aaa"}
	for _, key := range keys {
		rTree.Insert(key, true)
	}
	for _, key := range keys {
		val, ok := rTree.Get(key)
		if !ok {
			t.Fatalf("invalid get result expect: true, actual: not found")
		}
		if !val.(bool) {
			t.Fatalf("invalid get result expect: true, actual: false")
		}
		rTree.Delete(key)
		ok = rTree.Contains(key)
		if ok {
			t.Fatalf("invalid get result expect: not found")
		}
	}
}

// generateUUID is used to generate a random UUID
func generateUUID() string {
	buf := make([]byte, 16)
	if _, err := crand.Read(buf); err != nil {
		panic(fmt.Errorf("failed to read random bytes: %v", err))
	}

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%12x",
		buf[0:4],
		buf[4:6],
		buf[6:8],
		buf[8:10],
		buf[10:16])
}

func prepareDataSet() (*RadixTree, map[string]interface{}) {
	mMap := make(map[string]interface{})
	for i := 0; i < 1000; i++ {
		k := generateUUID()
		mMap[k] = i
	}

	rTree := NewRadixTreeFromMap(mMap)

	return rTree, mMap
}
