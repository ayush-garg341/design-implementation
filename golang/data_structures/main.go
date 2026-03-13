package main

import (
	"fmt"
)

var nums [3]int // [0 0 0]

// Struct — custom data type, Go's "class"
// Group related fields together
// No classes in Go - structs + methods + interfaces replace them.
// Struct is a value type — copying a struct copies all fields. Use pointer (*User) to avoid copying or to mutate.

type User struct {
	Name  string
	Age   int
	Email string
}

// Methods on structs (pointer receiver = can modify)
func (u *User) Greet() string {
	u.Name = "Ayush"
	return "Hi, I'm " + u.Name
}

// Interface — Go's polymorphism
// Implicit implementation — no "implements" keyword
// Any type that has the required methods satisfies the interface automatically.

type Animal interface {
	Sound() string
}

type Dog struct{}

func (d Dog) Sound() string { return "Woo!" }

type Cat struct{}

func (c Cat) Sound() string { return "Meow" }

func main() {
	// Array — fixed-size, rarely used directly
	// Fixed length, part of the type
	// [3]int and [5]int are different types. Size must be known at compile time.
	// Arrays are values — passing to a function copies the whole array.

	primes := [3]int{2, 3, 5}
	auto := [...]int{1, 2, 3} // compiler counts
	fmt.Println(primes[0])    // 2
	fmt.Println(auto[0])      // 1

	// Slice — dynamic array, used everywhere
	// Header = (pointer, length, capacity)
	// A slice references an underlying array. Multiple slices can share the same array.

	s := []int{1, 2, 3}
	s = append(s, 4) // [1 2 3 4]
	fmt.Println(s)
	fmt.Println(len(s), cap(s))

	// make(type, len, cap)
	s2 := make([]int, 5, 10) // len=5, cap=10
	fmt.Println(s2)

	// Slicing a slice (shares memory!)
	a := []int{10, 20, 30, 40}
	b := a[1:3] // [20 30] — same array!
	b[0] = 99   // a is now [10 99 30 40]
	fmt.Println(a)

	// slices are reference types. Modifying a sub-slice modifies the original.
	// append() may create a NEW underlying array when capacity is exceeded — original slice is unaffected.

	// Map — key-value store, like HashMap:: NOT thread safe
	// Unordered. Key can be any comparable type (string, int, struct). Value can be anything.

	ages := map[string]int{
		"Alice": 30, "Bob": 25,
	}
	ages["Charlie"] = 35 // add
	delete(ages, "Bob")  // remove

	// Check if key exists (always do this!)
	_, ok := ages["Dave"]
	if !ok {
		fmt.Println("not found")
	}

	// Iterate (order is random!)
	for k, v := range ages {
		fmt.Println("Key, value == ", k, v)
	}

	// Never read/write a map from multiple goroutines — use sync.Map or protect with sync.Mutex

	u := User{Name: "Alice", Age: 30}
	u.Email = "alice@example.com"

	fmt.Println("Before modification user struct:: ", u)

	name := u.Greet()
	fmt.Println("Name:: ", name)
	fmt.Println("After modification user struct:: ", u)

	// Channel & Pointer — Go-specific essentials

	// Channel - Goroutine Comms
	// Typed pipe for passing data between goroutines safely.

	// ch := make(chan int) // Unbuffered

	// Buffered (non-blocking up to cap)
	ch2 := make(chan int, 5)

	ch2 <- 42  // send
	v := <-ch2 // receive
	close(ch2) // signal done
	fmt.Println("Received from channel = ", v)

	// Pointer Memory address
	// Avoid copying large structs; allow mutation inside functions.
	x := 42
	p := &x  // pointer to x
	*p = 100 // modify via ptr
	fmt.Println("The value of p = ", *p)

	// new() allocates zero value
	p2 := new(int) // *int, value=0
	*p2 = 7
	fmt.Println("The value of p2 = ", *p2)

	// Both satisfy Animal — no declaration needed!
	// var dog = Dog{} is also correct
	var dog Animal = Dog{}
	fmt.Println(dog.Sound())

}
