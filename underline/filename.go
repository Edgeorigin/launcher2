package underline

import (
	"fmt"
	"path"
	"strings"

	v "github.com/hydrati/plugin-loader/utils/container"
)

type Filename interface {
	Underline
	Extname() Extname
	String() string
}

type filename struct {
	Underline
	extname Extname
}

func (f *filename) Extname() Extname {
	return f.extname
}

func (f *filename) String() string {
	return fmt.Sprintf("Filename{%s,%s}", f.Underline.String(), f.extname.String())
}

var (
	BaseExtname = ".7z"
)

func NewFilename(p string, category v.Option[string]) v.Result[Filename, error] {
	s := path.Base(p)
	if s == "." || s == "/" {
		return v.Err[Filename](ErrUnderlineParse)
	}

	e := path.Ext(s)
	n := strings.TrimSuffix(s, e)

	a := NewUnderline(n, category)
	if !a.Has() {
		return v.Err[Filename](a.Error().Value())
	}

	b := NewExtname(BaseExtname, e)
	if !a.Has() {
		return v.Err[Filename](b.Error().Value())
	}

	return v.Ok[Filename, error](&filename{a.Value(), b.Value()})
}
