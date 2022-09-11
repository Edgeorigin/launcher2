package container

import "sync"

const (
	StatePending = state(0)
	StateOk      = state(1)
)

type FutureState interface {
	Pending() bool
	Ok() bool
	String() string
}

type state uint8

func (s state) Pending() bool {
	return s == StatePending
}

func (s state) Ok() bool {
	return s == StateOk
}

func (s state) String() string {
	if s.Pending() {
		return "[PENDING]"
	}

	return "[OK]"
}

type FutureResult[T any] interface {
	Value[T]
	State() FutureState
}

type Await[T any] interface {
	Await() T
}

type Future[T any] interface {
	Poll() FutureResult[T]
	Await() T
}

type fresult[T any] struct {
	s FutureState
	p *T
}

func (f fresult[T]) Value() T {
	if f.p == nil {
		panic(ErrValue)
	}

	return *f.p
}

func (f fresult[T]) State() FutureState {
	return f.s
}

type future[T any] struct {
	mu *sync.Mutex
	w  chan T
	v  *T
	d  bool
}

func (f future[T]) Poll() FutureResult[T] {
	if f.d {
		return fresult[T]{StateOk, f.v}
	}

	if f.mu.TryLock() {
		defer f.mu.Unlock()

		v, ok := <-f.w

		if ok {
			f.v = &v
			return fresult[T]{StateOk, f.v}
		}
	}

	return fresult[T]{StatePending, nil}
}

func (f future[T]) Await() T {
	if !f.d {
		f.mu.Lock()
		defer f.mu.Unlock()
		if !f.d && f.v == nil {
			f.d = true
			v := (<-f.w)
			f.v = &v
		}
	}

	return *f.v
}

func newFuture[T any]() (chan<- T, future[T]) {
	ch := make(chan T, 1)
	f := future[T]{&sync.Mutex{}, ch, nil, false}
	return ch, f
}

func Async[T any](fn func() T) Future[T] {
	ch, f := newFuture[T]()

	go func() {
		v := fn()
		ch <- v
		close(ch)
	}()

	return f
}

func AsyncVal[T any](v T) Future[T] {
	ch, f := newFuture[T]()
	go func() {
		ch <- v
		close(ch)
	}()
	return f
}

func AsyncUnit(fn func()) Future[Unit] {
	return Async(func() Unit {
		fn()
		return UnitVal
	})
}

func Async2[T any](fn func(chan<- T)) Future[T] {
	ch, f := newFuture[T]()

	go func() {
		fn(ch)
		close(ch)
	}()

	return f
}
