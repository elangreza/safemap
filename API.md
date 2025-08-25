# SafeMap API Documentation

This document provides detailed API reference for the SafeMap package.

## Package Overview

```go
package safemap
```

The `safemap` package provides a thread-safe map implementation using goroutines and channels. It supports concurrent access and modification without the need for explicit locking mechanisms.

## Types

### SafeMap[K comparable, V any]

```go
type SafeMap[k comparable, v any] struct {
    // unexported fields
}
```

SafeMap is a thread-safe map implementation that uses channels for internal communication. It supports generic key-value pairs where the key must be comparable.

**Important Notes:**

- SafeMap must be initialized using `NewSafeMap()` function
- Direct instantiation (e.g., `SafeMap{}`) will cause panics when methods are called
- All operations are thread-safe and can be called from multiple goroutines

## Functions

### NewSafeMap

```go
func NewSafeMap[k comparable, v any]() *SafeMap[k, v]
```

NewSafeMap creates and returns a new instance of SafeMap. It initializes the internal goroutine that processes operations on the map.

**Parameters:**

- None

**Returns:**

- `*SafeMap[k, v]`: A pointer to a new SafeMap instance

**Example:**

```go
// String keys, int values
stringIntMap := safemap.NewSafeMap[string, int]()

// Int keys, string values
intStringMap := safemap.NewSafeMap[int, string]()

// Custom struct as value
type User struct {
    Name string
    Age  int
}
userMap := safemap.NewSafeMap[string, User]()
```

## Methods

### Set

```go
func (s *SafeMap[k, v]) Set(key k, val v)
```

Set stores a key-value pair in the SafeMap. If the key already exists, its value is updated.

**Parameters:**

- `key k`: The key to store (must be comparable)
- `val v`: The value to associate with the key

**Returns:**

- None

**Panics:**

- If SafeMap was not initialized with `NewSafeMap()`

**Example:**

```go
m := safemap.NewSafeMap[string, int]()
m.Set("apple", 5)
m.Set("banana", 3)
m.Set("apple", 10) // Updates existing key
```

### Get

```go
func (s *SafeMap[k, v]) Get(key k) (val v)
```

Get retrieves the value associated with the given key. If the key doesn't exist, returns the zero value of type V.

**Parameters:**

- `key k`: The key to retrieve

**Returns:**

- `val v`: The value associated with the key, or zero value if key doesn't exist

**Panics:**

- If SafeMap was not initialized with `NewSafeMap()`

**Example:**

```go
m := safemap.NewSafeMap[string, int]()
m.Set("apple", 5)

value := m.Get("apple")  // Returns 5
missing := m.Get("kiwi") // Returns 0 (zero value for int)
```

### Delete

```go
func (s *SafeMap[k, v]) Delete(key k)
```

Delete removes the key-value pair for the given key from the SafeMap. If the key doesn't exist, this operation is a no-op.

**Parameters:**

- `key k`: The key to remove

**Returns:**

- None

**Panics:**

- If SafeMap was not initialized with `NewSafeMap()`

**Example:**

```go
m := safemap.NewSafeMap[string, int]()
m.Set("apple", 5)
m.Delete("apple")
// Key "apple" is now removed from the map
```

### Exist

```go
func (s *SafeMap[k, v]) Exist(key k) bool
```

Exist checks whether the given key exists in the SafeMap.

**Parameters:**

- `key k`: The key to check

**Returns:**

- `bool`: true if the key exists, false otherwise

**Panics:**

- If SafeMap was not initialized with `NewSafeMap()`

**Example:**

```go
m := safemap.NewSafeMap[string, int]()
m.Set("apple", 5)

exists := m.Exist("apple")  // Returns true
missing := m.Exist("kiwi")  // Returns false
```

### Keys

```go
func (s *SafeMap[k, v]) Keys() iter.Seq[k]
```

Keys returns an iterator over all keys in the SafeMap. The iterator can be used with range loops.

**Parameters:**

- None

**Returns:**

- `iter.Seq[k]`: An iterator over the keys

**Panics:**

- If SafeMap was not initialized with `NewSafeMap()`

**Example:**

```go
m := safemap.NewSafeMap[string, int]()
m.Set("apple", 5)
m.Set("banana", 3)

for key := range m.Keys() {
    fmt.Println("Key:", key)
}
```

### All

```go
func (s *SafeMap[k, v]) All() iter.Seq2[k, v]
```

All returns an iterator over all key-value pairs in the SafeMap. The iterator can be used with range loops.

**Parameters:**

- None

**Returns:**

- `iter.Seq2[k, v]`: An iterator over key-value pairs

**Panics:**

- If SafeMap was not initialized with `NewSafeMap()`

**Example:**

```go
m := safemap.NewSafeMap[string, int]()
m.Set("apple", 5)
m.Set("banana", 3)

for key, value := range m.All() {
    fmt.Printf("%s: %d\n", key, value)
}
```

### Length

```go
func (s *SafeMap[k, v]) Length() int
```

Length returns the number of key-value pairs currently stored in the SafeMap.

**Parameters:**

- None

**Returns:**

- `int`: The number of key-value pairs

**Panics:**

- If SafeMap was not initialized with `NewSafeMap()`

**Example:**

```go
m := safemap.NewSafeMap[string, int]()
fmt.Println(m.Length()) // Prints: 0

m.Set("apple", 5)
m.Set("banana", 3)
fmt.Println(m.Length()) // Prints: 2
```

### GetMap

```go
func (s *SafeMap[k, v]) GetMap() map[k]v
```

GetMap returns a copy of the internal map. This method is useful when you need to work with a standard Go map or pass the data to functions expecting a regular map.

**Parameters:**

- None

**Returns:**

- `map[k]v`: A copy of the internal map

**Panics:**

- If SafeMap was not initialized with `NewSafeMap()`

**Important Notes:**

- The returned map is a copy, so modifications to it won't affect the SafeMap
- This operation creates a new map and copies all entries, which may be expensive for large maps

**Example:**

```go
m := safemap.NewSafeMap[string, int]()
m.Set("apple", 5)
m.Set("banana", 3)

mapCopy := m.GetMap()
fmt.Printf("Copy: %+v\n", mapCopy) // Prints: Copy: map[apple:5 banana:3]

// Modifying the copy doesn't affect the original SafeMap
mapCopy["orange"] = 8
fmt.Println(m.Length()) // Still prints: 2
```

## Thread Safety

All SafeMap methods are thread-safe and can be called concurrently from multiple goroutines without additional synchronization. The implementation uses a single internal goroutine that processes all operations sequentially through channels, ensuring:

1. **Data Race Prevention**: No data races can occur
2. **Consistency**: All operations are atomic from the caller's perspective
3. **Deadlock Freedom**: No possibility of deadlocks in the implementation

## Performance Characteristics

- **Time Complexity**: All operations are O(1) on average (same as regular Go maps)
- **Space Complexity**: O(n) where n is the number of key-value pairs
- **Concurrency**: Operations are processed sequentially by the internal goroutine
- **Channel Overhead**: Each operation involves channel communication

## Error Handling

SafeMap uses panics for error conditions:

1. **Uninitialized SafeMap**: All methods panic with message "safemap can be only accessed with NewSafeMap" if called on an uninitialized SafeMap
2. **No Other Errors**: Normal operations (Get on missing key, Delete on missing key) do not panic but return appropriate zero values or no-op behavior

## Best Practices

1. **Always use NewSafeMap()**: Never create SafeMap instances directly
2. **Check Exist() before Get()**: If you need to distinguish between zero values and missing keys
3. **Use iterators efficiently**: The Keys() and All() methods create snapshots, so use them when you need a consistent view
4. **Consider GetMap() for bulk operations**: If you need to perform many read operations, consider getting a copy first
5. **Handle zero values**: Remember that Get() returns zero values for missing keys

## Migration from sync.Map

If you're migrating from `sync.Map`, here are the key differences:

| Operation   | sync.Map                                   | SafeMap                     |
| ----------- | ------------------------------------------ | --------------------------- |
| Store       | `Store(key, value)`                        | `Set(key, value)`           |
| Load        | `Load(key) (value, ok)`                    | `Get(key)` + `Exist(key)`   |
| Delete      | `Delete(key)`                              | `Delete(key)`               |
| Range       | `Range(func(key, value interface{}) bool)` | `for k, v := range m.All()` |
| Type Safety | ❌ interface{}                             | ✅ Generic types            |
