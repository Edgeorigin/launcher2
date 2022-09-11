package underline

import (
	"fmt"
	"strings"

	"github.com/hydrati/plugin-loader/utils"
)

type Flags interface {
	String() string
	Has(rune) bool
	Size() int
	All() utils.Iterator[rune]
}

type _flags map[rune]bool

func (f *_flags) String() string {
	b := &strings.Builder{}

	for key := range *f {
		b.WriteRune(key)
	}

	return fmt.Sprintf("Flags{%s}", b.String())
}

func (f *_flags) Has(flag rune) bool {
	_, ok := (*f)[flag]
	return ok
}

func (f *_flags) Size() int {
	return len(*f)
}

func (f *_flags) All() utils.Iterator[rune] {
	r, w := utils.NewIterator[rune]()

	go func() {
		for key := range *f {
			w <- key
		}
		close(w)
	}()

	return r
}

func NewFlags(flags string) Flags {
	r := _flags{}

	for _, v := range flags {
		r[v] = true
	}

	return &r
}
