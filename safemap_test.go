package safemap

import (
	"fmt"
	"testing"
)

func TestSafeMapRace(t *testing.T) {
	m := NewSafeMap[int, int]()
	done := make(chan struct{})
	go func() {
		for i := range 1000 {
			m.Set(i, i)
		}
		done <- struct{}{}
	}()
	done2 := make(chan struct{})
	go func() {
		for i := range 1000 {
			res := m.Get(i)
			fmt.Println("res", res)
		}
		done2 <- struct{}{}
	}()
	_, _ = <-done, <-done2
}
