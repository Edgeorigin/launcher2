package utils

import (
	"fmt"
)

const (
	Lt _ordering = -1
	Eq _ordering = 0
	Gt _ordering = 1
)

func _toOrdering(x int) _ordering {
	if x == 1 {
		return Gt
	}

	if x == 0 {
		return Eq
	}

	if x == -1 {
		return Lt
	}

	panic("error: invalid ordering int")
}

type Ordering interface {
	fmt.Stringer

	Lt() bool
	Gt() bool
	Eq() bool
	LtEq() bool
	GtEq() bool
}

type _ordering int

func (o _ordering) Lt() bool {
	return o == Lt
}

func (o _ordering) Eq() bool {
	return o == Eq
}

func (o _ordering) Gt() bool {
	return o == Gt
}

func (o _ordering) LtEq() bool {
	return o.Lt() || o.Eq()
}

func (o _ordering) GtEq() bool {
	return o.Gt() || o.Eq()
}

func (o _ordering) String() string {
	if o.Eq() {
		return "[==]"
	}

	if o.Gt() {
		return "[>]"
	}

	if o.Lt() {
		return "[<]"
	}

	return "[?]"
}
