package archive

import (
	"fmt"
	"os"
	"path/filepath"

	u "github.com/hydrati/plugin-loader/utils"
	v "github.com/hydrati/plugin-loader/utils/container"
)

type Extractor struct {
	i    u.Iterator[File]
	dst  string
	done bool
	d    map[string]bool
}

func NewExtractor(i u.Iterator[File], dst string) *Extractor {
	return &Extractor{
		i:    i,
		dst:  dst,
		done: false,
		d:    make(map[string]bool),
	}
}

func (e *Extractor) Execute(p v.Option[int]) v.Future[v.Option[[]error]] {
	return v.Async(func() v.Option[[]error] {
		errs := make([]error, 0)
		async := make([]v.Future[v.Option[error]], 0)

		if e.done {
			return v.None[[]error]()
		}

		for i := range e.i {
			if p.Has() && p.Value() <= len(async) {
				for _, s := range async {
					v := s.Await()
					if v.Has() {
						errs = append(errs, v.Value())
					}
				}

				async = make([]v.Future[v.Option[error]], 0)
			}

			if _, ok := e.d[i.Path()]; !i.IsDir() && !ok {
				// if !i.IsDir() {
				e.d[i.Path()] = true
				async = append(async, e.extract(i))
			}
		}

		for _, s := range async {
			v := s.Await()
			if v.Has() {
				errs = append(errs, v.Value())
			}
		}

		if len(errs) == 0 {
			return v.None[[]error]()
		}

		return v.Some(errs)
	})
}

func (e *Extractor) Done() bool {
	return e.done
}

func (e *Extractor) extract(f File) v.Future[v.Option[error]] {
	fmt.Println("Post", f, e.dst, f.Path())

	return v.Async(func() v.Option[error] {
		p := filepath.Join(e.dst, f.Path())
		d := filepath.Dir(p)
		if d != "." {
			os.MkdirAll(d, 0777)
		}

		return f.Extract(p)
	})
}
