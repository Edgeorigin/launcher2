package container

import "fmt"

type Option[T any] interface {
	Value[T]
	None()
	Or(T) T
	OrWith(Value[T]) T
	Map(Pipe[T, T]) Option[T]
	And(Option[T]) Option[T]
	Has() bool
}

type option[T any] struct {
	p *T
}

func (o option[T]) Has() bool {
	return o.p != nil
}

func (o option[T]) Value() T {
	if o.Has() {
		return *o.p
	}

	panic(ErrValue)
}

func (o option[T]) None() {
	if o.Has() {
		panic(*o.p)
	}
}

func (o option[T]) Or(x T) T {
	if o.Has() {
		return *o.p
	}

	return x
}

func (o option[T]) OrWith(x Value[T]) T {
	if o.Has() {
		return *o.p
	}

	return x.Value()
}

func (o option[T]) Map(x Pipe[T, T]) Option[T] {
	if o.Has() {
		return Some(x.Pipe(*o.p))
	}

	return o
}

func (o option[T]) And(x Option[T]) Option[T] {
	if o.Has() {
		return o
	}

	return x
}

func (o option[T]) String() string {
	if o.Has() {
		return fmt.Sprintf("Option.Some{%+v}", *o.p)
	}

	return "Option.None{}"
}

func Some[T any](x T) Option[T] {
	return option[T]{&x}
}

func None[T any]() Option[T] {
	return option[T]{nil}
}

func ToErrorOption(x error) Option[error] {
	if x == nil {
		return None[error]()
	}

	return Some(x)
}
