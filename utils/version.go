package utils

import (
	"fmt"

	"github.com/mcuadros/go-version"
)

type Version interface {
	fmt.Stringer

	Compare(x Version) Ordering
	Normalize() Version
}

type _version string

func (v _version) Compare(x Version) Ordering {
	return _toOrdering(version.CompareSimple(v.String(), x.String()))
}

func (v _version) Normalize() Version {
	return _version(version.Normalize(v.String()))
}

func (v _version) String() string {
	return string(v)
}

func NewVersion(v string) Version {
	return _version(v)
}
