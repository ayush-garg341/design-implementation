# Go: Interface vs Generics

> A complete guide to understanding the difference, when to use each, and how to combine them.

---

## Table of Contents

- [The Core Idea](#the-core-idea)
- [The One Question That Decides](#the-one-question-that-decides)
- [Side-by-Side Comparison](#side-by-side-comparison)
- [How Each One Works](#how-each-one-works)
  - [Interface - Runtime Polymorphism](#interface--runtime-polymorphism)
  - [Generics - Compile-time Polymorphism](#generics--compile-time-polymorphism)
- [When to Use Interface](#when-to-use-interface)
  - [Plugin / Strategy Pattern](#1-plugin--strategy-pattern--swap-implementations)
  - [Mixed Types in One Collection](#2-mixed-types-in-one-collection)
  - [Testability - Mock the Dependency](#3-testability--mock-the-dependency)
- [When to Use Generics](#when-to-use-generics)
  - [Data Structures](#1-data-structures-that-work-for-any-type)
  - [Utility Functions](#2-utility-functions--filter-map-reduce)
  - [Numeric Operations](#3-numeric-operations--operators-need-constraints)
  - [Map Keys with comparable](#4-map-keys--comparable-constraint)
- [Common Mistakes](#common-mistakes)
- [The Decision Rule](#the-decision-rule)
- [Using Both Together](#using-both-together)
- [One-Line Summary](#one-line-summary)

---

## The Core Idea

Both interfaces and generics solve the same high-level problem - *"write code that works for multiple types"* - but they solve it in fundamentally different ways.

| | Interface | Generics `[T]` |
|---|---|---|
| **When type is resolved** | Runtime | Compile time |
| **Mechanism** | Method dispatch (vtable) | Type substitution |
| **Types in one collection** | Yes - different types allowed | No - one type per instance |
| **Overhead** | Slight (vtable lookup) | Zero (compiler generates typed code) |
| **Best for** | Behaviour abstraction | Data structures & algorithms |

---

## The One Question That Decides

> **Are you abstracting *behaviour* (what it can do) or abstracting a *container/algorithm* (what it holds)?**

- **Behaviour** → Interface
- **Container / Algorithm** → Generics

---

## Side-by-Side Comparison

```go
// INTERFACE - "anything that can speak"
// Different types, different behaviour, live together in one slice

type Speaker interface {
    Speak() string
}

type Dog struct{}
func (d Dog) Speak() string { return "woof" }

type Cat struct{}
func (c Cat) Speak() string { return "meow" }

// Dog and Cat are DIFFERENT types in ONE slice
animals := []Speaker{Dog{}, Cat{}}
for _, a := range animals {
    fmt.Println(a.Speak()) // woof / meow - each does its own thing
}
```

```go
// GENERICS - "a stack that holds anything"
// Same logic, different type per instance, fully type-safe

type Stack[T any] struct {
    items []T
}

func (s *Stack[T]) Push(v T)  { s.items = append(s.items, v) }
func (s *Stack[T]) Pop() T    {
    n := len(s.items) - 1
    v := s.items[n]
    s.items = s.items[:n]
    return v
}

// Each stack holds ONE specific type - compiler enforces it
intStack := &Stack[int]{}
strStack := &Stack[string]{}

intStack.Push(42)
intStack.Push("oops") // COMPILE ERROR - caught immediately
```

**Key insight:** In the interface example, `Dog` and `Cat` are *different types* living in the *same slice*. In the generics example, `Stack[int]` and `Stack[string]` are the *same code shape* applied to *one type each*.

---

## How Each One Works

### Interface - Runtime Polymorphism

The type is resolved **at runtime**. One compiled function handles all types via method dispatch.

```go
// What actually happens at runtime:
// animals[0] → looks up Dog's method table → calls Dog.Speak()
// animals[1] → looks up Cat's method table → calls Cat.Speak()

// The compiler does NOT know which Speak() will be called at compile time
// It figures it out at runtime by looking at the actual type stored in the interface value
```

### Generics - Compile-time Polymorphism

The type is resolved **at compile time**. The compiler generates separate typed versions of the function/struct.

```go
// What actually happens at compile time:
// Stack[int]    → compiler generates a Stack specifically for int
// Stack[string] → compiler generates a Stack specifically for string

// No runtime lookup - the compiler bakes the type in directly
// This is why generics have zero overhead
```

---

## When to Use Interface

### 1. Plugin / Strategy Pattern - Swap Implementations

Use interface when you want to swap implementations without changing the code that uses them.

```go
// Define the behaviour contract
type Storage interface {
    Save(key, value string) error
    Get(key string) (string, error)
}

// Two completely different implementations
type RedisStorage  struct{ client *redis.Client }
type MemoryStorage struct{ data map[string]string }

func (r RedisStorage)  Save(k, v string) error { return r.client.Set(...) }
func (m MemoryStorage) Save(k, v string) error { m.data[k] = v; return nil }

// Service takes the interface - doesn't care which implementation
type LinkService struct {
    store Storage
}

// Swap without changing LinkService at all
svc := LinkService{store: RedisStorage{...}}  // production
svc  = LinkService{store: MemoryStorage{...}} // tests
```

### 2. Mixed Types in One Collection

Use interface when you need to hold *different types* with *different behaviour* in the same slice.

```go
type Notifier interface {
    Notify(msg string) error
}

type EmailNotifier struct{ addr    string }
type SMSNotifier   struct{ phone   string }
type SlackNotifier struct{ webhook string }

// ALL THREE different types in one slice - only possible via interface
notifiers := []Notifier{
    EmailNotifier{"a@b.com"},
    SMSNotifier{"+91999"},
    SlackNotifier{"https://hooks..."},
}

for _, n := range notifiers {
    n.Notify("server down!") // each type handles it differently
}
```

### 3. Testability - Mock the Dependency

The most important pattern in real Go codebases. Define a DB interface so you can swap real DB with a mock in tests.

```go
type DB interface {
    GetUser(id int) (User, error)
}

// Real implementation - hits actual Postgres
type PostgresDB struct{ db *sql.DB }
func (p PostgresDB) GetUser(id int) (User, error) { /* real query */ }

// Mock for tests - no real DB needed, no network, instant
type MockDB struct{}
func (m MockDB) GetUser(id int) (User, error) {
    return User{Name: "Alice"}, nil // fake data
}

// Handler works with both - testable without a running database
func NewHandler(db DB) *Handler { return &Handler{db: db} }

// In production:
handler := NewHandler(PostgresDB{db: realDB})

// In tests:
handler := NewHandler(MockDB{})
```

---

## When to Use Generics

### 1. Data Structures That Work for Any Type

Use generics when the logic is identical regardless of the type being stored.

```go
type Stack[T any] struct {
    items []T
}

func (s *Stack[T]) Push(v T)        { s.items = append(s.items, v) }
func (s *Stack[T]) Pop() T          {
    n := len(s.items) - 1
    v := s.items[n]
    s.items = s.items[:n]
    return v
}
func (s *Stack[T]) Len() int        { return len(s.items) }
func (s *Stack[T]) IsEmpty() bool   { return len(s.items) == 0 }

// Fully type-safe - no type assertion, no panic risk
intStack := &Stack[int]{}
intStack.Push(42)
val := intStack.Pop() // val is int - compiler knows

strStack := &Stack[string]{}
strStack.Push("hello")
s := strStack.Pop() // s is string - compiler knows
```

### 2. Utility Functions - Filter, Map, Reduce

Use generics for functions where the algorithm is the same but the type varies.

```go
// Works for []int, []string, []User - any type
func Filter[T any](s []T, fn func(T) bool) []T {
    var out []T
    for _, v := range s {
        if fn(v) {
            out = append(out, v)
        }
    }
    return out
}

// Works for any two types - T input, R output
func Map[T, R any](slice []T, fn func(T) R) []R {
    result := make([]R, len(slice))
    for i, v := range slice {
        result[i] = fn(v)
    }
    return result
}

// Usage
users  := []User{{"Alice", 30}, {"Bob", 17}, {"Charlie", 25}}
adults := Filter(users, func(u User) bool { return u.Age >= 18 })
// result type is []User - not []interface{}, no type assertion needed

nums := []int{1, 2, 3}
strs := Map(nums, strconv.Itoa)
// result type is []string
```

### 3. Numeric Operations - Operators Need Constraints

Use generics when you need `+`, `-`, `<`, `>` to work across multiple numeric types. Interfaces cannot do this - you cannot call operators on `interface{}`.

```go
// Interface can't do this - you can't call + on interface{}
// Generics with a union constraint makes it possible

func Sum[T int | int64 | float64](nums []T) T {
    var total T
    for _, n := range nums {
        total += n
    }
    return total
}

func Min[T int | float64 | string](a, b T) T {
    if a < b { return a }
    return b
}

fmt.Println(Sum([]int{1, 2, 3}))              // 6
fmt.Println(Sum([]float64{1.1, 2.2, 3.3}))   // 6.6
fmt.Println(Min("apple", "mango"))            // "apple"
fmt.Println(Min(3, 7))                        // 3
```

> **Tip:** Use `golang.org/x/exp/constraints` for ready-made constraints:
> - `constraints.Ordered` → int | float | string (supports `<` `>`)
> - `constraints.Integer` → all int types
> - `constraints.Float`   → float32 | float64

### 4. Map Keys - comparable Constraint

Use `comparable` constraint when T needs to be used as a map key.

```go
// Generic Set - works for any hashable type
type Set[T comparable] struct {
    m map[T]struct{}
}

func NewSet[T comparable]() *Set[T] {
    return &Set[T]{m: make(map[T]struct{})}
}

func (s *Set[T]) Add(v T)         { s.m[v] = struct{}{} }
func (s *Set[T]) Has(v T) bool    { _, ok := s.m[v]; return ok }
func (s *Set[T]) Remove(v T)      { delete(s.m, v) }
func (s *Set[T]) Len() int        { return len(s.m) }

intSet := NewSet[int]()
strSet := NewSet[string]()
// comparable constraint ensures T can be used as a map key
```

---

## Common Mistakes

### Mistake 1 - Using generics when you need interface (mixed types)

```go
// WRONG - trying to hold mixed types with generics
// T cannot be int AND string at the same time
items := []any{1, "hello", true}
PrintAll(items) // works but T = any, all type info is lost

// BETTER - use interface directly for truly mixed types
func PrintAll(items []fmt.Stringer) {
    for _, v := range items {
        fmt.Println(v.String())
    }
}
```

### Mistake 2 - Using interface{} when you need generics (type safety lost)

```go
// WRONG - using interface{}, losing all type safety
type OldStack struct{ items []interface{} }
func (s *OldStack) Pop() interface{} { /* ... */ }

s := &OldStack{}
s.items = append(s.items, 42)
val := s.Pop().(int) // type assertion - can PANIC at runtime!

// RIGHT - generics gives compile-time safety, zero runtime risk
type Stack[T any] struct{ items []T }
func (s *Stack[T]) Pop() T { /* ... */ }

gs := &Stack[int]{}
gs.items = append(gs.items, 42)
val := gs.Pop() // already int - no assertion, no panic possible
```

### Mistake 3 - Calling methods on T with constraint `any`

```go
// WRONG - trying to call methods on an unconstrained T
func Process[T any](item T) {
    item.Process() // COMPILE ERROR - T is "any", has no methods
}

// RIGHT - if you need method dispatch, use an interface
type Processor interface {
    Process()
}
func Process(item Processor) {
    item.Process() // works - interface guarantees Process() exists
}
```

> **Rule:** If you need to call methods on `T`, use an interface - not generics with `any`.

---

## The Decision Rule

Ask yourself these questions in order:

```
Q1: Do different types need DIFFERENT BEHAVIOUR for the same operation?
    e.g. Dog.Speak() → "woof", Cat.Speak() → "meow"
    → USE INTERFACE

Q2: Do you need to store MIXED TYPES in one collection?
    e.g. []Shape holding Circle, Square, Triangle simultaneously
    → USE INTERFACE

Q3: Are you building a CONTAINER or ALGORITHM that works identically for any type?
    e.g. Stack, Queue, Set, Filter, Map, Sort - same logic, different type
    → USE GENERICS

Q4: Do you need OPERATORS like +, -, <, > across multiple types?
    e.g. Sum of ints AND floats, Min/Max across ints AND strings
    → USE GENERICS (with union constraint)

Q5: Do you want to SWAP IMPLEMENTATIONS (prod vs test, Redis vs Memory)?
    e.g. LinkService works with any Storage that has Save/Get methods
    → USE INTERFACE

Q6: Do you want ZERO RUNTIME OVERHEAD and the type is known at call time?
    e.g. Compiler generates specific code per type, no boxing/unboxing
    → USE GENERICS
```

---

## Using Both Together

You can and should combine them. A generic struct with an interface field is a powerful pattern.

```go
// Generic cache that stores any type (T),
// but the loader is an interface so it can be swapped

type Loader[T any] interface {
    Load(key string) (T, error)
}

type Cache[T any] struct {
    store  map[string]T
    loader Loader[T] // interface inside a generic struct
}

func (c *Cache[T]) Get(key string) (T, error) {
    if v, ok := c.store[key]; ok {
        return v, nil // cache hit
    }
    v, err := c.loader.Load(key) // interface call - swappable
    if err == nil {
        c.store[key] = v
    }
    return v, err
}

// T = User  → Cache[User]    with any UserLoader implementation
// T = Product → Cache[Product] with any ProductLoader implementation
// Generics handles type safety of what's stored
// Interface handles the swappable loading strategy
```

---

## One-Line Summary

| | Summary |
|---|---|
| **Interface** | *"I don't care what type it is, as long as it **can do** this."* → about capability, resolved at runtime, swappable implementations |
| **Generics `[T]`** | *"Same algorithm, but let me **choose the type**."* → about structure, resolved at compile time, type-safe without assertions |

---

### Quick Reference Card

| Scenario | Use |
|---|---|
| Service layer / repository | Interface |
| Mock in tests | Interface |
| Mixed types in one slice | Interface |
| Swap Redis ↔ Memory ↔ Postgres | Interface |
| Stack, Queue, Set, LinkedList | Generics |
| Filter, Map, Reduce, Contains | Generics |
| Sum, Min, Max across numeric types | Generics |
| Type-safe wrapper (no assertions) | Generics |
| Generic struct + swappable loader | Both |

---

*Part of the Go learning series - covering interfaces, generics, data structures, concurrency, and more.*

