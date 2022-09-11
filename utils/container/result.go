package container

import "fmt"

type Result[T any, E error] interface {
	Value[T]
	Error[E]
	Errored() bool
	Or(T) T
	OrWith(Value[T]) T
	MapErr(Pipe[T, E]) T
	Map(Pipe[T, T]) Result[T, E]
	Has() bool
	Raw() (*T, *E)
	Option() (Option[T], Option[E])
}

type result[T any, E error] struct {
	p *T
	e *E
}

func (r result[T, E]) Raw() (*T, *E) {
	return r.p, r.e
}

func (r result[T, E]) Option() (Option[T], Option[E]) {
	if r.Has() {
		return Some(r.Value()), r.Error()
	}

	return None[T](), r.Error()
}

func (r result[T, E]) Has() bool {
	return r.p != nil && r.e == nil
}

func (r result[T, E]) Errored() bool {
	return !r.Has()
}

func (r result[T, E]) Error() Option[E] {
	if r.Has() {
		return None[E]()
	}

	return Some(*r.e)
}

func (r result[T, E]) Or(x T) T {
	if r.Has() {
		return *r.p
	}

	return x
}

func (r result[T, E]) Map(x Pipe[T, T]) Result[T, E] {
	if r.Has() {
		return Ok[T, E](x.Pipe(*r.p))
	}

	return r
}

func (r result[T, E]) OrWith(x Value[T]) T {
	if r.Has() {
		return *r.p
	}

	return x.Value()
}

func (r result[T, E]) MapErr(x Pipe[T, E]) T {
	if r.Has() {
		return *r.p
	}

	return x.Pipe(*r.e)
}

func (r result[T, E]) Value() T {
	if r.Has() {
		return *r.p
	}

	panic(*r.e)
}

func (r result[T, E]) String() string {
	if r.Has() {
		return fmt.Sprintf("Result.Ok{%+v}", *r.p)
	}

	return fmt.Sprintf("Result.Err{%+v}", *r.e)
}

func Ok[T any, E error](x T) Result[T, E] {
	return result[T, E]{&x, nil}
}

func Err[T any, E error](x E) Result[T, E] {
	return result[T, E]{nil, &x}
}

func Resuify[T any](x T, y error) Result[T, error] {
	if y != nil {
		return Err[T](y)
	}

	return Ok[T, error](x)
}

func ToErr[T any, E error, U any](x Result[U, E]) Result[T, E] {
	return Err[T](x.Error().Value())
}
