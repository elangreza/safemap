package safemap

import (
	"iter"
	"maps"
)

type (
	operation[k comparable, v any] struct {
		op        string
		key       k
		value     v
		replyChan chan<- any
	}

	SafeMap[k comparable, v any] struct {
		opChan chan operation[k, v]
	}
)

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
