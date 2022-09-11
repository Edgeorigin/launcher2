package underline

import (
	"errors"
	"fmt"
	"strings"

	v "github.com/hydrati/plugin-loader/utils/container"
)

var (
	ErrInvalidFullExtname = errors.New("error: invalid full extname")
)

type Extname interface {
	Flags() Flags
	Full() string
	Base() string
	String() string
}

type _extname struct {
	flags Flags
	base  string
	full  string
}

func (e *_extname) Flags() Flags {
	return e.flags
}

func (e *_extname) Base() string {
	return e.base
}

func (e *_extname) Full() string {
	return e.full
}

func (e *_extname) String() string {
	return fmt.Sprintf("Extname{%s %s}", e.Base(), e.flags.String())
}

func NewExtname(base, full string) v.Result[Extname, error] {
	full = strings.TrimSpace(full)
	base = strings.TrimSpace(base)

	if !strings.HasPrefix(full, base) {
		return v.Err[Extname](ErrInvalidFullExtname)
	}

	flags := NewFlags(strings.TrimPrefix(full, base))

	return v.Ok[Extname, error](&_extname{full: full, base: base, flags: flags})
}
