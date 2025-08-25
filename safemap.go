package safemap

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
			}
		}
	}()

	return sm
}

func (s *SafeMap[k, v]) Set(key k, val v) {
	if s.opChan == nil {
		panic("safemap can be only access with NewSafeMap")
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
		panic("safemap can be only access with NewSafeMap")
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
		panic("safemap can be only access with NewSafeMap")
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
		panic("safemap can be only access with NewSafeMap")
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

func (s *SafeMap[k, v]) Keys() func(yield func(k) bool) {
	if s.opChan == nil {
		panic("safemap can be only access with NewSafeMap")
	}

	replyChan := make(chan any)
	s.opChan <- operation[k, v]{
		op:        "getKeys",
		replyChan: replyChan,
	}

	keys := <-replyChan
	keysSlice := keys.([]k)
	return func(yield func(k) bool) {
		for _, k := range keysSlice {
			if !yield(k) {
				break
			}
		}
	}
}

func (s *SafeMap[k, v]) Seq2() func(yield func(k, v) bool) {
	return func(yield func(k, v) bool) {

	}
}
