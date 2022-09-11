package execute

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/hydrati/plugin-loader/env"
	v "github.com/hydrati/plugin-loader/utils/container"
)

type LookupResult []string

var (
	execExt = &sync.Map{} // sync.Map[string, bool]
)

func init() {
	MakeExecutable(".cmd", ".wcs")
}

func MakeExecutable(ext ...string) {
	for _, i := range ext {
		SetExtnameExecutable(i, true)
	}
}

func SetExtnameExecutable(ext string, executable bool) {
	execExt.Store(ext, executable)
}

func IsExecExt(ext string) bool {
	s, ok := execExt.Load(ext)
	return ok && s.(bool)
}

func Lookup(rootPath string) v.Result[LookupResult, error] {
	x := v.Resuify(os.ReadDir(rootPath))
	if x.Errored() {
		return v.ToErr[LookupResult](x)
	}

	f := make(LookupResult, 0)

	for _, i := range x.Value() {
		if !i.IsDir() {
			if IsExecExt(filepath.Ext(i.Name())) {
				f = append(
					f,
					env.PathResolve(env.PathJoin(rootPath, i.Name())),
				)
			}
		}
	}

	return v.Ok[LookupResult, error](f)
}
