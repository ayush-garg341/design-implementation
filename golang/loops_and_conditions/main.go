package main

import (
	"fmt"
	"strconv"
	"time"
)

// Type switch — used with interfaces
func describe(i interface{}) {
	switch v := i.(type) {
	case int:
		fmt.Println("int:", v)
	case string:
		fmt.Println("string:", v)
	default:
		fmt.Println("other")
	}
}

func worker(jobs <-chan int, done <-chan bool) {
	for { // Infinite loop
		select {
		case job := <-jobs: // received a job
			fmt.Println("processing job:", job)
		case <-done: // received stop signal
			fmt.Println("worker shutting down")
			return // exits goroutine
		}
	}
}

func ticker() {
	// This pattern is used in real services for health checks, cache refresh, metrics flush, and polling.

	ticker := time.NewTicker(500 * time.Millisecond) // fires every 500ms
	stop := time.After(2 * time.Second)              // stop after 2s
	count := 0

	for {
		select {
		case t := <-ticker.C:
			count++
			fmt.Printf("tick #%d at %s\n", count, t.Format("15:04:05.000"))
		case <-stop:
			ticker.Stop()
			fmt.Println("stopped after", count, "ticks")
			return
		}
	}
}

func main() {

	// if / else if / else No parentheses needed
	// Braces {} are mandatory in Go — no single-line if without them.
	// Unlike C/Java — no parentheses around the condition. Braces are always required.

	age := 20

	if age < 18 {
		fmt.Println("minor")
	} else if age < 60 {
		fmt.Println("adult")
	} else {
		fmt.Println("senior")
	}

	// if with init statement Very idiomatic in Go
	// Declare and check in one line. The variable lives only inside the if block.
	// Pattern: if init; condition { }

	if val, err := strconv.Atoi("42"); err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("value:", val) // val only exists here
	}

	myMap := map[string]string{
		"name": "Ayush",
	}

	// Also used to check map key existence
	if v, ok := myMap["name"]; ok {
		fmt.Println(v)
	}

	// switch No fallthrough by default
	// Go switch does NOT fall through automatically - no break needed between cases.

	day := "Mon"
	switch day {
	case "Mon", "Tue":
		fmt.Println("Weekday")
	case "Sat", "Sun":
		fmt.Println("Weekend")
	default:
		fmt.Println("Unknown")
	}

	// Switch with no expression
	// (acts like if-else chain)

	x := 42
	switch {
	case x < 0:
		fmt.Println("neg")
	case x == 0:
		fmt.Println("zero")
	default:
		fmt.Println("pos")
	}

	// LOOPS — Go has only ONE loop keyword: for
	for i := 0; i < 5; i++ {
		fmt.Println(i)
	}

	// while-style for Go has no while keyword
	n := 1
	for n < 100 {
		n *= 2
	}
	fmt.Println(n)

	// for range — iterate collections Most used loop
	// Works on slices, arrays, maps, strings, and channels.

	// Slice — index, value
	nums := []int{10, 20, 30}
	for i, v := range nums {
		fmt.Println(i, v)
	}

	// Skip index with _
	for _, v := range nums {
		fmt.Println(v)
	}

	// Index Only
	for i := range nums {
		fmt.Println(i)
	}

	// Map — key, value
	m := map[string]int{
		"a": 1, "b": 2,
	}

	for k, v := range m {
		fmt.Println(k, v)
	}

	// String — index, rune
	for i, ch := range "hello" {
		fmt.Println(i, ch)
	} // ch is rune, not byte

	// the loop variable v in range is a COPY — modifying v does NOT change the original slice element. Use nums[i] to modify in place.

	// break / continue / labels Control flow
	for i := 0; i < 5; i++ {
		if i == 5 {
			break
		}
		fmt.Println(i)
	}

	// continue — skip iteration
	for i := 0; i < 5; i++ {
		if i%2 == 0 {
			continue
		}
		fmt.Println(i)
	}

	// Labeled break — break outer loop
	// from inside inner loop
outer:
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if j == 1 {
				break outer
			}
			fmt.Println(i, j)
		}
	}

	// Labeled break/continue is Go's clean way to exit nested loops — preferred over using a flag variable.

	jobs := make(chan int, 5)
	done := make(chan bool)
	go worker(jobs, done)

	for i := 0; i < 3; i++ {
		jobs <- i
	}

	time.Sleep(100 * time.Millisecond) // let worker process
	done <- true                       // send stop signal
	time.Sleep(50 * time.Millisecond)  // wait for shutdown

	ticker()
}
