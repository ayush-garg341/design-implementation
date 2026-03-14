package main

import "fmt"

type Queue[T any] struct {
	items []T
}

func (q *Queue[T]) append(v T) {
	q.items = append(q.items, v)
}

func (q *Queue[T]) pop() T {
	n := len(q.items) - 1
	v := q.items[n]
	q.items = q.items[:n]
	return v
}

func (q *Queue[T]) len() int {
	return len(q.items)
}

func (q *Queue[T]) peek() T {
	n := len(q.items) - 1
	return q.items[n]
}

func main() {
	q := &Queue[int]{}

	q.append(1) // [1]
	q.append(2) // [1, 2]
	q.append(3) // [1, 2, 3]

	len := q.len()
	fmt.Println("Len = ", len)

	fmt.Println("Last Element = ", q.peek())
	fmt.Println("Last Element popped = ", q.pop())

	fmt.Println("Len = ", q.len())

}
