package utils

import (
	v "github.com/hydrati/plugin-loader/utils/container"
)

type Iterator[T any] <-chan T

// func (it Iterator[T]) Buffered(n int) Iterator[T] {
// 	return NewBufferedIterator(it, n)
// }

type ResultIterator[T any, E error] <-chan v.Result[T, E]

func NewIterator[T any]() (r Iterator[T], w chan<- T) {
	c := make(chan T)
	return c, c
}

func NewResultIterator[T any, E error]() (r ResultIterator[T, E], w chan<- v.Result[T, E]) {
	c := make(chan v.Result[T, E])
	return c, c
}

func NewBufferedIterator[T any](it Iterator[T], buffer int) Iterator[T] {
	c := make(chan T, buffer)
	go func() {
		for i := range it {
			c <- i
		}
	}()
	return c
}
