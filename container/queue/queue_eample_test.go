package queue

import (
	"fmt"
)

// ExampleQueue use queue to solve Josephus problem,
// from Algorithm 4th 1.3.37
func ExampleQueue_josephusProblem() {
	var (
		m = 7
		n = 2
	)

	q := NewSListQueue()

	for i := 0; i < m; i++ {
		q.EnQueue(i)
	}

	for q.Size() != 0 {
		for i := 1; i < n; i++ {
			q.EnQueue(q.DeQueue())
		}
		fmt.Printf("%d ", q.DeQueue().(int))
	}
	// Output: 1 3 5 0 4 2 6
}
