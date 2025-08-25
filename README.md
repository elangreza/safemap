# SafeMap

[![Go Reference](https://pkg.go.dev/badge/github.com/elangreza/safemap.svg)](https://pkg.go.dev/github.com/elangreza/safemap)
[![Go Report Card](https://goreportcard.com/badge/github.com/elangreza/safemap)](https://goreportcard.com/report/github.com/elangreza/safemap)

A thread-safe map implementation in Go using goroutines and channels, providing concurrent access without explicit locking.

## Features

- **Thread-safe**: Concurrent access and modification without explicit locking
- **Generic**: Support for any comparable key type and any value type
- **Channel-based**: Uses Go's channels for internal communication
- **Iterator support**: Provides range iteration over keys and key-value pairs
- **Memory safe**: Prevents data races and ensures safe concurrent operations

## Installation

```bash
go get github.com/elangreza/safemap
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/elangreza/safemap"
)

func main() {
    // Create a new SafeMap
    m := safemap.NewSafeMap[string, int]()

    // Set values
    m.Set("apple", 5)
    m.Set("banana", 3)
    m.Set("orange", 8)

    // Get values
    value := m.Get("apple")
    fmt.Println("apple:", value) // Output: apple: 5

    // Check if key exists
    exists := m.Exist("banana")
    fmt.Println("banana exists:", exists) // Output: banana exists: true

    // Get map length
    length := m.Length()
    fmt.Println("length:", length) // Output: length: 3

    // Iterate over keys
    for key := range m.Keys() {
        fmt.Println("key:", key)
    }

    // Iterate over key-value pairs
    for key, value := range m.All() {
        fmt.Printf("%s: %d\n", key, value)
    }

    // Delete a key
    m.Delete("banana")

    // Get a copy of the internal map
    mapCopy := m.GetMap()
    fmt.Printf("Map copy: %+v\n", mapCopy)
}
```

## API Reference

### Types

#### SafeMap[K comparable, V any]

A thread-safe map that supports concurrent access and modification.

**Important**: SafeMap must be initialized using `NewSafeMap()`. Direct instantiation will cause panics.

### Functions

#### NewSafeMap[K comparable, V any]() \*SafeMap[K, V]

Creates and returns a new instance of SafeMap. This function initializes the internal goroutine that processes operations on the map.

```go
m := safemap.NewSafeMap[string, int]()
```

### Methods

#### Set(key K, val V)

Sets the value for the given key in the SafeMap.

```go
m.Set("key", "value")
```

#### Get(key K) V

Retrieves the value for the given key from the SafeMap. Returns the zero value of type V if the key doesn't exist.

```go
value := m.Get("key")
```

#### Delete(key K)

Removes the key-value pair for the given key from the SafeMap.

```go
m.Delete("key")
```

#### Exist(key K) bool

Checks if the given key exists in the SafeMap.

```go
exists := m.Exist("key")
```

#### Keys() iter.Seq[K]

Returns an iterator over all keys in the SafeMap. Can be used with range loops.

```go
for key := range m.Keys() {
    fmt.Println(key)
}
```

#### All() iter.Seq2[K, V]

Returns an iterator over all key-value pairs in the SafeMap. Can be used with range loops.

```go
for key, value := range m.All() {
    fmt.Printf("%v: %v\n", key, value)
}
```

#### Length() int

Returns the number of key-value pairs in the SafeMap.

```go
length := m.Length()
```

#### GetMap() map[K]V

Returns a copy of the internal map of the SafeMap. This is useful when you need to work with a standard Go map or pass the data to functions expecting a regular map.

```go
mapCopy := m.GetMap()
```

## Concurrent Usage

SafeMap is designed for concurrent use. Here's an example of multiple goroutines safely accessing the same SafeMap:

```go
package main

import (
    "fmt"
    "sync"
    "github.com/elangreza/safemap"
)

func main() {
    m := safemap.NewSafeMap[int, string]()
    var wg sync.WaitGroup

    // Multiple goroutines writing to the map
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            m.Set(id, fmt.Sprintf("value-%d", id))
        }(i)
    }

    // Multiple goroutines reading from the map
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            if m.Exist(id) {
                value := m.Get(id)
                fmt.Printf("Read %d: %s\n", id, value)
            }
        }(i)
    }

    wg.Wait()
    fmt.Printf("Final map length: %d\n", m.Length())
}
```

## Performance Considerations

- SafeMap uses channels for internal communication, which provides safety but may have different performance characteristics compared to mutex-based implementations
- Each operation involves channel communication, so for high-frequency operations, consider batching when possible
- The internal goroutine processes operations sequentially, ensuring consistency but potentially limiting parallelism for read operations

## Error Handling

SafeMap will panic in the following cases:

- Attempting to use SafeMap methods on an instance not created with `NewSafeMap()`
- This design ensures that SafeMap is always properly initialized and prevents undefined behavior

## Comparison with Standard Map + Mutex

| Feature           | SafeMap          | sync.Map       | Map + Mutex    |
| ----------------- | ---------------- | -------------- | -------------- |
| Type Safety       | ✅ Generic       | ❌ interface{} | ✅ Generic     |
| Memory Safety     | ✅ Channel-based | ✅ Lock-free   | ✅ Mutex-based |
| Read Performance  | ⚠️ Sequential    | ✅ Concurrent  | ⚠️ Blocking    |
| Write Performance | ⚠️ Sequential    | ✅ Optimized   | ⚠️ Blocking    |
| API Simplicity    | ✅ Simple        | ⚠️ Complex     | ✅ Simple      |
| Iterator Support  | ✅ Native        | ❌ Manual      | ❌ Manual      |

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
