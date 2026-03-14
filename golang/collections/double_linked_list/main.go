package main

import "fmt"

// Node: element in the list
type Node[T any] struct {
	value T
	prev  *Node[T] // pointer to previous node (nil if head)
	next  *Node[T] // pointer to next node     (nil if tail)
}

// Doubly linked list: owns head + tail + length
type DoublyLinkedList[T any] struct {
	head   *Node[T]
	tail   *Node[T]
	length int
}

// Constructor
func NewDLL[T any]() *DoublyLinkedList[T] {
	head := &Node[T]{}
	tail := &Node[T]{}
	head.prev = nil
	tail.next = nil
	head.next = tail
	tail.prev = head
	return &DoublyLinkedList[T]{head: head, tail: tail, length: 0}
}

// Len: number of elements
func (dll *DoublyLinkedList[T]) Len() int { return dll.length }

// Is empty
func (dll *DoublyLinkedList[T]) IsEmpty() bool { return dll.length == 0 }

// PushFront — insert at head O(1)
func (dll *DoublyLinkedList[T]) PushFront(val T) *Node[T] {
	node := &Node[T]{value: val}

	node.next = dll.head.next
	node.prev = dll.head.next.prev
	dll.head.next.prev = node
	dll.head.next = node

	dll.length++

	return node
}

// PushBack - insert at tail O(1)
func (dll *DoublyLinkedList[T]) PushBack(val T) *Node[T] {
	node := &Node[T]{value: val}

	node.next = dll.tail.prev.next
	node.prev = dll.tail.prev
	dll.tail.prev.next = node
	dll.tail.prev = node

	dll.length++

	return node

}

// Insert after - insert after a given node O(1)
func (dll *DoublyLinkedList[T]) InsertAfter(node *Node[T], val T) *Node[T] {
	if node == nil {
		return nil
	}

	newNode := &Node[T]{value: val}
	newNode.next = node.next
	newNode.prev = node
	node.next.prev = newNode
	node.next = newNode

	dll.length++

	return newNode
}

// Remove - remove a specific node O(1)
func (dll *DoublyLinkedList[T]) Remove(node *Node[T]) T {
	// Stitch prev and next together, skipping this node
	value := node.value
	node.prev.next = node.next
	node.next.prev = node.prev

	node.prev = nil
	node.next = nil

	dll.length--
	return value
}

// Pop front
func (dll *DoublyLinkedList[T]) PopFront() (T, bool) {
	var zero T
	if dll.length == 0 {
		return zero, false
	}
	val := dll.Remove(dll.head.next)
	return val, true
}

// Pop Back
func (dll *DoublyLinkedList[T]) PopBack() (T, bool) {
	var zero T
	if dll.length == 0 {
		return zero, false
	}
	val := dll.Remove(dll.tail.prev)
	return val, true
}

// Forward - head to tail

func (dll *DoublyLinkedList[T]) ForEach(fn func(T)) {
	current := dll.head.next
	for current != dll.tail {
		fn(current.value)
		current = current.next // move forward
	}
}

// Backward - tail to head
func (dll *DoublyLinkedList[T]) ForEachReverse(fn func(T)) {
	current := dll.tail.prev
	for current != dll.head {
		fn(current.value)
		current = current.prev // move backward
	}
}

// ToSlice — convert to []T for easy printing
func (dll *DoublyLinkedList[T]) ToSlice() []T {
	result := make([]T, 0, dll.length)
	dll.ForEach(func(v T) { result = append(result, v) })
	return result
}

func (dll *DoublyLinkedList[T]) ToSliceReverse() []T {
	result := make([]T, 0, dll.length)
	dll.ForEachReverse(func(v T) { result = append(result, v) })
	return result
}

// Search — find first node matching condition O(n)
func (dll *DoublyLinkedList[T]) Find(fn func(T) bool) (*Node[T], bool) {
	current := dll.head.next
	for current != dll.tail {
		if fn(current.value) {
			return current, true
		}
		current = current.next
	}
	return nil, false
}

func (dll *DoublyLinkedList[T]) PeekFront() (T, bool) {
	var zero T
	if dll.length == 0 {
		return zero, false
	}
	return dll.head.next.value, true
}

func (dll *DoublyLinkedList[T]) PeekBack() (T, bool) {
	var zero T
	if dll.length == 0 {
		return zero, false
	}
	return dll.tail.prev.value, true
}

func main() {
	// Usage:
	list := NewDLL[int]()
	// Build: 10 -> 20 -> 30 -> 40
	list.PushBack(20)
	list.PushBack(30)
	list.PushFront(10)
	n40 := list.PushBack(40)

	fmt.Println("list:", list.ToSlice())           // [10 20 30 40]
	fmt.Println("reverse:", list.ToSliceReverse()) // [40 30 20 10]
	fmt.Println("len:", list.length)               // 4

	// Remove node 40 directly — O(1), no search needed
	list.Remove(n40)
	fmt.Println("after remove 40:", list.ToSlice()) // [10 20 30]

	// Pop from both ends
	v, _ := list.PopFront()
	fmt.Println("popFront:", v, "-", list.ToSlice()) // 10 -> [20 30]

	v, _ = list.PopBack()
	fmt.Println("popBack:", v, "-", list.ToSlice()) // 30 -> [20]

	node, found := list.Find(func(v int) bool { return v == 20 })
	if found {
		list.Remove(node)                                             // O(1) remove once found — no second traversal!
		fmt.Println("After remove:", node.value, "-", list.ToSlice()) // 20 -> []
	}

}
