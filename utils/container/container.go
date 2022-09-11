package container

import "errors"

type Value[T any] interface {
	Value() T
}

type Error[E error] interface {
	Error() Option[E]
}

type Pipe[T any, I any] interface {
	Pipe(I) T
}

var (
	ErrValue = errors.New("error: none or rejected value container")
)
