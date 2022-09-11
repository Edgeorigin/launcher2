package archive

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/hydrati/plugin-loader/archive/sevenzip"
	u "github.com/hydrati/plugin-loader/utils"
	v "github.com/hydrati/plugin-loader/utils/container"
)

var (
	ErrNotExist = errors.New("error: path does not exist")
	ErrIsNotDir = errors.New("error: path is not dir")
)

type Archive interface {
	Exist(path string) bool
	All() u.Iterator[File]
	Open(path string) v.Result[File, error]
	// ReadDir(path string) v.Result[u.Iterator[File], error]
	Path() string
}

type File interface {
	GetReader() v.Result[io.ReadCloser, error]
	Pipe(dst io.Writer) v.Option[error]
	Extract(dstPath string) v.Option[error]
	Path() string
	Size() int64
	PackedSize() v.Option[int64]
	Modified() time.Time
	Attributes() string
	CRC() string
	Encrypted() string
	Method() string
	Block() int
	HasAttr(attr ...rune) bool
	IsDir() bool
	Clone() File
}

type file struct {
	i *sevenzip.Entry
	a *sevenzip.Archive
}

func (f *file) Clone() File {
	return &file{f.i, f.a}
}

func (f *file) Path() string {
	return f.i.Path
}

func (f *file) HasAttr(a ...rune) bool {
	for _, i := range a {
		if !strings.ContainsRune(f.i.Attributes, i) {
			return false
		}
	}
	return true
}

func (f *file) Size() int64 {
	return f.i.Size
}

func (f *file) IsDir() bool {
	return f.HasAttr('D')
}

func (f *file) PackedSize() v.Option[int64] {
	if f.i.PackedSize == -1 {
		return v.None[int64]()
	}
	return v.Some(f.i.PackedSize)
}

func (f *file) Modified() time.Time {
	return f.i.Modified
}

func (f *file) Attributes() string {
	return f.i.Attributes
}

func (f *file) CRC() string {
	return f.i.CRC
}

func (f *file) Encrypted() string {
	return f.i.Encrypted
}

func (f *file) Method() string {
	return f.i.Method
}

func (f *file) Block() int {
	return f.i.Block
}

func (f *file) GetReader() v.Result[io.ReadCloser, error] {
	return v.Resuify(f.a.GetFileReader(f.i.Path))
}

func (f *file) Pipe(dst io.Writer) v.Option[error] {
	return v.ToErrorOption(f.a.ExtractToWriter(dst, f.i.Path))
}

func (f *file) Extract(dstPath string) v.Option[error] {
	return v.ToErrorOption(f.a.ExtractToFile(dstPath, f.i.Path))
}

func (f *file) String() string {
	return fmt.Sprintf("File{%s}", f.i.Path)
}

type archive struct {
	inner *sevenzip.Archive
	// f     map[string]*sevenzip.Entry
}

func (a *archive) Path() string {
	return a.inner.Path
}

func (a *archive) find(path string) v.Option[*sevenzip.Entry] {
	path = filepath.Clean(path)
	for _, i := range a.inner.Entries {
		if i.Path == path {
			return v.Some(&i)
		}
	}

	return v.None[*sevenzip.Entry]()
}

func (a *archive) Exist(path string) bool {
	return a.find(path).Has()
}

func (a *archive) Open(path string) v.Result[File, error] {
	e := a.find(path)
	if !e.Has() {
		return v.Err[File](ErrNotExist)
	}

	return v.Ok[File, error](&file{e.Value(), a.inner})
}

func (a *archive) All() u.Iterator[File] {
	r, w := u.NewIterator[File]()

	go func() {
		defer close(w)
		for _, i := range a.inner.Entries {
			w <- &file{&i, a.inner}
		}
	}()

	return r
}

func (a *archive) String() string {
	return fmt.Sprintf("File{%s}", a.inner.Path)
}

// func hasPrefix(p, prefix string) bool {
// 	if strings.HasPrefix(p, prefix) {
// 		return true
// 	}
// 	return strings.HasPrefix(strings.ToLower(p), strings.ToLower(prefix))
// }

// func (a *archive) ReadDir(path string) v.Result[u.Iterator[File], error] {
// 	path = filepath.Clean(path)
// 	f := a.Open(path)
// 	if f.Error().Has() {
// 		return v.Err[u.Iterator[File]](f.Error().Value())
// 	}

// 	if !f.Value().IsDir() {

// 		return v.Err[u.Iterator[File]](ErrIsNotDir)
// 	}

// 	r, w := u.NewIterator[File]()

// 	o := f.Value()

// 	go func() {
// 		defer close(w)
// 		for i := range a.All() {
// 			if hasPrefix(i.Path(), path) && i != o {
// 				w <- i
// 			}
// 		}
// 	}()

// 	return v.Ok[u.Iterator[File], error](r)
// }

func OpenArchive(path string, password v.Option[string]) v.Result[Archive, error] {
	var arc *sevenzip.Archive
	if password.Has() {
		_arc, err := sevenzip.NewEncryptedArchive(path, password.Value())

		if err != nil {
			return v.Err[Archive](err)
		}

		arc = _arc
	} else {
		_arc, err := sevenzip.NewArchive(path)

		if err != nil {
			return v.Err[Archive](err)
		}

		arc = _arc
	}

	// f := make(map[string]*sevenzip.Entry)

	return v.Ok[Archive, error](&archive{arc})

}
