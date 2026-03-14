package main

import "fmt"

type Deque[T any] struct {
	items []T
}

func (d *Deque[T]) append(v T) {
	d.items = append(d.items, v)
}

func (d *Deque[T]) appendleft(v T) {
	d.items = append([]T{v}, d.items...)
}

func (d *Deque[T]) pop() T {
	n := len(d.items) - 1
	v := d.items[n]
	d.items = d.items[:n]
	return v
}

func (d *Deque[T]) popleft() T {
	v := d.items[0]
	d.items = d.items[1:]
	return v
}

func (d *Deque[T]) len() int {
	return len(d.items)
}

func (d *Deque[T]) peekback() T {
	n := len(d.items) - 1
	return d.items[n]
}

func (d *Deque[T]) peekfront() T {
	return d.items[0]

}

func (d *Deque[T]) isempty() bool {
	return len(d.items) == 0
}

func main() {
	d := &Deque[int]{}
	d.append(1)     // [1]
	d.append(2)     // [1 2]
	d.appendleft(0) // [0 1 2]

	fmt.Println(d.popleft())
	fmt.Println(d.pop())
	fmt.Println(d.len())
}
