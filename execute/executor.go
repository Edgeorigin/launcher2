package execute

import (
	"errors"
	"os/exec"
	"path/filepath"
	"sync"

	v "github.com/hydrati/plugin-loader/utils/container"
)

var (
	executors           = &sync.Map{} // sync.Map[string, ScriptExecutor]
	ErrNotFoundExecutor = errors.New("error: not found executor")
)

func init() {
	RegisterExecutor(".cmd", FunctionalScriptExecutor(cmdExecutor))
	RegisterExecutor(".wcs", FunctionalScriptExecutor(wcsExecutor))
}

type ScriptExecutor interface {
	ExecuteScript(path string) v.Future[v.Result[[]byte, *exec.ExitError]]
}

type FunctionalScriptExecutor func(path string) v.Result[[]byte, *exec.ExitError]

func (f FunctionalScriptExecutor) ExecuteScript(path string) v.Future[v.Result[[]byte, *exec.ExitError]] {
	return v.Async(func() v.Result[[]byte, *exec.ExitError] {
		return f(path)
	})
}

func RegisterExecutor(ext string, f ScriptExecutor) bool {
	_, ok := executors.Load(ext)

	executors.Store(ext, f)

	return ok
}

func RemoveExecutor(ext string) bool {
	_, ok := executors.LoadAndDelete(ext)

	return ok
}

func getExecutor(ext string) v.Option[ScriptExecutor] {
	f, ok := executors.Load(ext)

	if !ok {
		return v.None[ScriptExecutor]()
	}

	return v.Some(f.(ScriptExecutor))
}

type ExecuteOutput struct {
	Output  []byte
	Ignored bool
}

func Execute(path string) v.Future[v.Result[ExecuteOutput, *exec.ExitError]] {
	return v.Async(func() v.Result[ExecuteOutput, *exec.ExitError] {
		ext := filepath.Ext(path)
		executor := getExecutor(ext)

		if !executor.Has() {
			return v.Ok[ExecuteOutput, *exec.ExitError](ExecuteOutput{Output: nil, Ignored: true})
		}

		r := executor.Value().ExecuteScript(PathResolve(path)).Await()
		if r.Errored() {
			return v.ToErr[ExecuteOutput](r)
		}

		return v.Ok[ExecuteOutput, *exec.ExitError](ExecuteOutput{Output: r.Value(), Ignored: false})
	})
}

func FindBin(name string) string {
	return v.Must(exec.LookPath(name))
}

func wcsExecutor(path string) v.Result[[]byte, *exec.ExitError] {
	o, err := exec.Command(FindBin("PECMD"), "LOAD", path).Output()
	if err != nil {
		return v.Err[[]byte](err.(*exec.ExitError))
	}

	return v.Ok[[]byte, *exec.ExitError](o)
}

func cmdExecutor(path string) v.Result[[]byte, *exec.ExitError] {
	o, err := exec.Command(FindBin("cmd"), "/c", path).Output()
	if err != nil {
		return v.Err[[]byte](err.(*exec.ExitError))
	}

	return v.Ok[[]byte, *exec.ExitError](o)
}
