package safemap

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSafeMap_Exist(t *testing.T) {
	m := NewSafeMap[int, int]()

	for i := range 10 {
		m.Set(i, i)
	}

	for i := range 10 {
		ok := m.Exist(i)
		assert.True(t, ok)
	}
}

func TestSafeMap_Keys(t *testing.T) {
	m := NewSafeMap[int, int]()

	for i := range 10 {
		m.Set(i, i)
	}

	keys := m.Keys()
	var keysSlice []int
	for key := range keys {
		keysSlice = append(keysSlice, key)
	}
	for i := range 10 {
		assert.Contains(t, keysSlice, i)
	}
}

func TestSafeMap_All(t *testing.T) {
	m := NewSafeMap[int, int]()

	for i := range 10 {
		m.Set(i+1, i)
	}

	items := m.All()
	var itemsSlice []int
	var keysSlice []int
	for key, item := range items {
		itemsSlice = append(itemsSlice, item)
		keysSlice = append(keysSlice, key)
	}
	for i := range 10 {
		assert.Contains(t, itemsSlice, i)
		assert.Contains(t, keysSlice, i+1)
	}
}

func TestSafeMap_Get(t *testing.T) {
	m := NewSafeMap[int, int]()

	for i := range 10 {
		m.Set(i, i)
	}

	for i := range 10 {
		v := m.Get(i)
		assert.Equal(t, i, v)
	}
}

func TestSafeMap_Delete(t *testing.T) {
	m := NewSafeMap[int, int]()

	for i := range 10 {
		m.Set(i, i)
	}

	for i := range 10 {
		m.Delete(i)
	}

	for i := range 10 {
		ok := m.Exist(i)
		assert.False(t, ok)
	}
}

func TestSafeMap_Length(t *testing.T) {
	m := NewSafeMap[int, int]()

	for i := range 10 {
		m.Set(i, i)
	}

	assert.Equal(t, 10, m.Length())
}

func TestSafeMap_GetMap(t *testing.T) {
	m := NewSafeMap[int, int]()

	for i := range 3 {
		m.Set(i, i)
	}

	assert.Equal(t, map[int]int{0: 0, 1: 1, 2: 2}, m.GetMap())
}

func TestSafeMap_Panic(t *testing.T) {

	m := &SafeMap[int, int]{}
	assert.Panics(t, func() { m.Get(1) })
	assert.Panics(t, func() { m.Set(1, 1) })
	assert.Panics(t, func() { m.Delete(1) })
	assert.Panics(t, func() { m.Exist(1) })
	assert.Panics(t, func() { m.Keys() })
	assert.Panics(t, func() { m.All() })
	assert.Panics(t, func() { m.Length() })
	assert.Panics(t, func() { m.GetMap() })

}

func TestSafeMapRace(t *testing.T) {
	m := NewSafeMap[int, int]()

	var wg sync.WaitGroup

	wg.Add(8)

	go func() {
		for i := range 1000 {
			m.Set(i, i)
		}
		wg.Done()
	}()

	go func() {
		for i := range 1000 {
			m.Get(i)
		}
		wg.Done()
	}()

	go func() {
		for i := range 1000 {
			m.Delete(i)
		}
		wg.Done()
	}()

	go func() {
		for i := range 1000 {
			m.Exist(i)
		}
		wg.Done()
	}()

	go func() {
		m.Keys()
		wg.Done()
	}()

	go func() {
		m.All()
		wg.Done()
	}()

	go func() {
		m.Length()
		wg.Done()
	}()

	go func() {
		m.GetMap()
		wg.Done()
	}()

	wg.Wait()
}
