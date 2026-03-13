package main

import (
	"fmt"
	"math/rand/v2"
)

// var with explicit type
// Explicit - use when zero value matters or type clarity is needed.
// Works at package level (outside functions). Preferred for package-level vars.
var name string
var count int = 0
var isReady bool

// var with inferred type
// Like := but using var keyword — less idiomatic inside functions.
// Same as := but can be used at package scope. Inside functions, prefer :=
var lastName = "Garg"
var counter = 100

// const — compile-time constants
// Value must be known at compile time. Can't be changed later.
const Pi = 3.14
const MaxRetries int = 3

// iota — auto-incrementing in const blocks
const (
	Sunday  = iota // 0
	Monday         // 1
	Tuesday        // 2
)

// var block — grouped declarations
// Clean way to declare multiple package-level variables together.
var (
	host    string = "localhost"
	port    int    = 8080
	timeout float64
)

func random() (int, error) {
	return rand.IntN(100), nil
}

func main() {

	// Declare first
	name := "Alice"

	// Then reassign with =
	name = "Bob"

	// Key rule: := declares + assigns. = only assigns (variable must already exist).

	// Multiple at once
	x, y := 10, 20

	// Used heavily with function returns
	val, err := random()

	if err == nil {
		fmt.Println("x, y = ", x, y)
	}
	fmt.Println("name, random num = ", name, val)
}
