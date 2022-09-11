package loader

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/hydrati/plugin-loader/archive"
	"github.com/hydrati/plugin-loader/discovery"
	"github.com/hydrati/plugin-loader/env"
	"github.com/hydrati/plugin-loader/execute"
	"github.com/hydrati/plugin-loader/underline"
	v "github.com/hydrati/plugin-loader/utils/container"
)

var (
	ErrModeUnimplemented = errors.New("error: mode unimplemented")
)

var (
	scriptBackupFolderName = "plugin_script_backup"
	defaultPTask           = v.Some(4)
)

type ErrorArray []error

func (a ErrorArray) Error() string {
	return fmt.Sprintf("errors: %#v", a)
}

func NewPackageLoader(path string, filename underline.Filename, system *discovery.System) *PackageLoader {
	profile := system.Profile.Value()
	return &PackageLoader{path, filename, system, profile}
}

func LoadPackage(path string, filename underline.Filename, system *discovery.System) v.Future[v.Result[v.Unit, error]] {
	return NewPackageLoader(path, filename, system).Load()
}

func LoadSystem(system *discovery.System) v.Result[v.Unit, error] {
	fmt.Println(system)
	fmt.Println()
	profile := system.Profile.Value()
	fmt.Println(profile)

	for k := range profile.Packages {
		loader := NewPackageLoader(k, system.Profile.Value().Packages[k], system)
		fmt.Println(loader)
		ret := loader.Load().Await()
		if ret.Errored() {
			return ret
		}
	}

	return v.Ok[v.Unit, error](v.UnitVal)
}

type PackageLoader struct {
	path     string
	filename underline.Filename
	system   *discovery.System
	profile  *discovery.Profile
}

func (loader *PackageLoader) LoadNormal() v.Result[v.Unit, error] {
	a := loader.preloadClean()
	if a.Errored() {
		return v.ToErr[v.Unit](a)
	}

	b := loader.extractNormal()
	if b.Errored() {
		return v.ToErr[v.Unit](b)
	}

	c := loader.executeScripts()
	if c.Errored() {
		return v.ToErr[v.Unit](c)
	}

	d := loader.postClean(c.Value().Lookup)
	if d.Errored() {
		return v.ToErr[v.Unit](d)
	}

	return v.Ok[v.Unit, error](v.UnitVal)
}

func (loader *PackageLoader) extractNormal() v.Result[v.Unit, error] {
	file := archive.OpenArchive(loader.path, v.None[string]())
	if file.Errored() {
		return v.ToErr[v.Unit](file)
	}

	extractor := archive.NewExtractor(file.Value().All(), loader.system.Dir)
	res := extractor.Execute(defaultPTask).Await()
	if res.Has() {
		return v.Err[v.Unit](error(ErrorArray(res.Value())))
	}

	return v.Ok[v.Unit, error](v.UnitVal)
}

func (loader *PackageLoader) preloadClean() v.Result[v.Unit, error] {
	scripts := execute.Lookup(loader.system.Dir)
	if scripts.Errored() {
		return v.ToErr[v.Unit](scripts)
	}

	for _, file := range scripts.Value() {
		os.Remove(file)
	}

	return v.Ok[v.Unit, error](v.UnitVal)
}

type ScriptResults struct {
	Lookup execute.LookupResult
	Output map[string]execute.ExecuteOutput
	Errors map[string]*exec.ExitError
}

func (loader *PackageLoader) executeScripts() v.Result[ScriptResults, error] {
	scripts := execute.Lookup(loader.system.Dir)
	if scripts.Errored() {
		return v.ToErr[ScriptResults](scripts)
	}

	outputs := make(map[string]execute.ExecuteOutput)
	errs := make(map[string]*exec.ExitError)

	for _, file := range scripts.Value() {
		script := env.PathJoin(loader.system.Dir, file)
		fmt.Println("Execute", script)
		res := execute.Execute(script).Await()
		fmt.Println(file, res)
		if res.Errored() {
			errs[file] = res.Error().Value()
		} else {
			outputs[file] = res.Value()
		}
	}

	return v.Ok[ScriptResults, error](
		ScriptResults{Lookup: scripts.Value(), Output: outputs, Errors: errs},
	)
}

func (loader *PackageLoader) postClean(scripts execute.LookupResult) v.Result[v.Unit, error] {
	target := env.PathJoin(loader.system.Dir, scriptBackupFolderName, loader.filename.ToUnderline().String3())
	err := os.MkdirAll(target, 0777)
	if err != nil {
		return v.Err[v.Unit](err)
	}

	for _, file := range scripts {
		err := os.Rename(file, env.PathJoin(target, filepath.Base(file)))
		if err != nil {
			return v.Err[v.Unit](err)
		}
	}
	return v.Ok[v.Unit, error](v.UnitVal)
}

func (loader *PackageLoader) Load() v.Future[v.Result[v.Unit, error]] {
	return v.Async(func() v.Result[v.Unit, error] {
		fmt.Println("Flags", loader.filename.Extname().Full(), loader.filename.Extname())
		// Normal
		if loader.filename.Extname().Full() == ".7z" {
			fmt.Println("normal")
			return loader.LoadNormal()
		}

		// Localboost, Unimplemented
		if loader.filename.Extname().Flags().Has('l') {
			fmt.Println("localboost2")
			return v.Err[v.Unit](ErrModeUnimplemented)
		}

		// Disabled, Skip
		if loader.filename.Extname().Flags().Has('f') {
			fmt.Println("disabled")
			return v.Ok[v.Unit, error](v.UnitVal)
		}

		// Others, Skip
		fmt.Println("skip")
		return v.Ok[v.Unit, error](v.UnitVal)
	})
}
