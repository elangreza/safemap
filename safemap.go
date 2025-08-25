package safemap

import (
	"iter"
	"maps"
)

type (
	// operation represents a request to perform an operation on the SafeMap.
	// It includes the operation type, key, value (if applicable), and a channel to send the result back.
	operation[k comparable, v any] struct {
		op        string
		key       k
		value     v
		replyChan chan any
	}

	// SafeMap is a thread-safe map implementation using goroutines and channels.
	// It supports concurrent access and modification of the map without the need for explicit locking.
	// for initializing must use NewSafeMap function. if initialization NewSafeMap is not used will be panic if not used.
	SafeMap[k comparable, v any] struct {
		opChan chan operation[k, v]
	}
)

// NewSafeMap creates and returns a new instance of SafeMap.
// It initializes the internal goroutine that processes operations on the map.
func NewSafeMap[k comparable, v any]() *SafeMap[k, v] {
	sm := &SafeMap[k, v]{
		opChan: make(chan operation[k, v]),
	}
	data := make(map[k]v)

	go func() {
		for op := range sm.opChan {
			switch op.op {
			case "set":
				data[op.key] = op.value
				op.replyChan <- struct{}{}
			case "get":
				op.replyChan <- data[op.key]
			case "delete":
				delete(data, op.key)
				op.replyChan <- struct{}{}
			case "exist":
				_, ok := data[op.key]
				op.replyChan <- ok
			case "getMap":
				copyMap := make(map[k]v, len(data))
				maps.Copy(copyMap, data)
				op.replyChan <- copyMap
			case "getLen":
				op.replyChan <- len(data)
			}
		}
	}()

	return sm
}

// Set sets the value for the given key in the SafeMap.
// If the SafeMap was not initialized using NewSafeMap, it panics.
func (s *SafeMap[k, v]) Set(key k, val v) {
	if s.opChan == nil {
		panic("safemap can be only accessed with NewSafeMap")
	}

	replyChan := make(chan any)
	s.opChan <- operation[k, v]{
		op:        "set",
		key:       key,
		value:     val,
		replyChan: replyChan,
	}
	<-replyChan
}

// Get retrieves the value for the given key from the SafeMap.
// If the SafeMap was not initialized using NewSafeMap, it panics.
func (s *SafeMap[k, v]) Get(key k) (val v) {
	if s.opChan == nil {
		panic("safemap can be only accessed with NewSafeMap")
	}

	replyChan := make(chan any)
	s.opChan <- operation[k, v]{
		op:        "get",
		key:       key,
		replyChan: replyChan,
	}

	reply := <-replyChan
	return reply.(v)
}

// Delete removes the key-value pair for the given key from the SafeMap.
// If the SafeMap was not initialized using NewSafeMap, it panics.
func (s *SafeMap[k, v]) Delete(key k) {
	if s.opChan == nil {
		panic("safemap can be only accessed with NewSafeMap")
	}

	replyChan := make(chan any)
	s.opChan <- operation[k, v]{
		op:        "delete",
		key:       key,
		replyChan: replyChan,
	}

	<-replyChan
}

// Exist checks if the given key exists in the SafeMap.
// If the SafeMap was not initialized using NewSafeMap, it panics.
func (s *SafeMap[k, v]) Exist(key k) bool {
	if s.opChan == nil {
		panic("safemap can be only accessed with NewSafeMap")
	}

	replyChan := make(chan any)
	s.opChan <- operation[k, v]{
		op:        "exist",
		key:       key,
		replyChan: replyChan,
	}

	exist := <-replyChan
	return exist.(bool)
}

// Keys returns a slice of all keys in the SafeMap.
// If the SafeMap was not initialized using NewSafeMap, it panics.
// example
//
//	m := NewSafeMap[int, int]()
//	m.Set(1, 2)
//	for key := range m.Keys() {
//		fmt.Println(key)
//	}
func (s *SafeMap[k, v]) Keys() iter.Seq[k] {
	if s.opChan == nil {
		panic("safemap can be only accessed with NewSafeMap")
	}

	replyChan := make(chan any)
	s.opChan <- operation[k, v]{
		op:        "getMap",
		replyChan: replyChan,
	}

	m := <-replyChan

	return maps.Keys(m.(map[k]v))
}

// All returns a slice of all key-value pairs in the SafeMap.
// If the SafeMap was not initialized using NewSafeMap, it panics.
// example
//
//	m := NewSafeMap[int, int]()
//	m.Set(1, 2)
//	for key, value := range m.All() {
//		fmt.Println(key, value)
//	}
func (s *SafeMap[k, v]) All() iter.Seq2[k, v] {
	if s.opChan == nil {
		panic("safemap can be only accessed with NewSafeMap")
	}

	replyChan := make(chan any)
	s.opChan <- operation[k, v]{
		op:        "getMap",
		replyChan: replyChan,
	}

	m := <-replyChan

	return maps.All(m.(map[k]v))
}

// Length returns the number of key-value pairs in the SafeMap.
// If the SafeMap was not initialized using NewSafeMap, it panics.
func (s *SafeMap[k, v]) Length() int {
	if s.opChan == nil {
		panic("safemap can be only accessed with NewSafeMap")
	}

	replyChan := make(chan any)
	s.opChan <- operation[k, v]{
		op:        "getLen",
		replyChan: replyChan,
	}

	length := <-replyChan
	return length.(int)
}

// GetMap returns a copy of the internal map of the SafeMap.
// If the SafeMap was not initialized using NewSafeMap, it panics.
func (s *SafeMap[k, v]) GetMap() map[k]v {
	if s.opChan == nil {
		panic("safemap can be only accessed with NewSafeMap")
	}

	replyChan := make(chan any)
	s.opChan <- operation[k, v]{
		op:        "getMap",
		replyChan: replyChan,
	}

	items := <-replyChan
	return items.(map[k]v)
}
