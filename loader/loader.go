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
	"github.com/sirupsen/logrus"
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
	logrus.Debugf("error array %t", a)
	return fmt.Sprintf("errors: %#v", a)
}

func NewPackageLoader(path string, filename underline.Filename, system *discovery.System) *PackageLoader {
	profile := system.Profile.Value()
	logrus.Debugf("create package loader path = %v, filename = %v, system = %#v", path, filename, system)
	return &PackageLoader{path, filename, system, profile}
}

func LoadPackage(path string, filename underline.Filename, system *discovery.System) v.Future[v.Result[v.Unit, error]] {
	return NewPackageLoader(path, filename, system).Load()
}

func LoadSystem(system *discovery.System) v.Result[v.Unit, error] {
	profile := system.Profile.Value()

	logrus.Debugf("load system profile %#v", system)
	logrus.Debugf("use profile %#v", profile)

	for k := range profile.Packages {
		logrus.Debugf("load package %s", k)
		loader := NewPackageLoader(k, system.Profile.Value().Packages[k], system)
		fmt.Println(loader)
		ret := loader.Load().Await()
		logrus.Debugf("loaded package %s, %+v", k, ret)
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
	logrus.Debugf("[%s] load normal", loader.path)

	logrus.Debugf("[%s] preload clean", loader.path)
	a := loader.preloadClean()
	logrus.Debugf("[%s] preload clean done, %#v", loader.path, a)

	if a.Errored() {
		return v.ToErr[v.Unit](a)
	}

	logrus.Debugf("[%s] extract files normal", loader.path)
	b := loader.extractNormal()
	logrus.Debugf("[%s] extract file normal done, %#v", loader.path, b)
	if b.Errored() {
		return v.ToErr[v.Unit](b)
	}

	logrus.Debugf("[%s] execute scripts", loader.path)
	c := loader.executeScripts()
	logrus.Debugf("[%s] execute scripts done, %#v", loader.path, c)
	if c.Errored() {
		return v.ToErr[v.Unit](c)
	}

	logrus.Debugf("[%s] post clean", loader.path)
	d := loader.postClean(c.Value().Lookup)
	logrus.Debugf("[%s] post clean done, %#v", loader.path, d)
	if d.Errored() {
		return v.ToErr[v.Unit](d)
	}

	logrus.Debugf("[%s] load normal done", loader.path)

	return v.Ok[v.Unit, error](v.UnitVal)
}

func (loader *PackageLoader) extractNormal() v.Result[v.Unit, error] {
	logrus.Debugf("[%s] open archive", loader.path)
	file := archive.OpenArchive(loader.path, v.None[string]())
	logrus.Debugf("[%s] open archive done, %#v", loader.path, file)
	if file.Errored() {
		return v.ToErr[v.Unit](file)
	}

	logrus.Debugf("[%s] new extractor", loader.path)
	extractor := archive.NewExtractor(file.Value().All(), loader.system.Dir)
	logrus.Debugf("[%s] extract +p[%v]", loader.path, defaultPTask)
	res := extractor.Execute(defaultPTask).Await()
	logrus.Debugf("[%s] extract done, %#v", loader.path, res)
	if res.Has() {
		return v.Err[v.Unit](error(ErrorArray(res.Value())))
	}

	return v.Ok[v.Unit, error](v.UnitVal)
}

func (loader *PackageLoader) preloadClean() v.Result[v.Unit, error] {
	logrus.Debugf("[%s] preload clean", loader.path)
	logrus.Debugf("[%s] lookup scripts", loader.path)
	scripts := execute.Lookup(loader.system.Dir)
	logrus.Debugf("[%s] preload will clean scripts, %#v", loader.path, scripts)
	if scripts.Errored() {
		return v.ToErr[v.Unit](scripts)
	}

	for _, file := range scripts.Value() {
		logrus.Debugf("[%s] preload clean, remove file, %s", loader.path, file)
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
	logrus.Debugf("[%s] execute scripts", loader.path)
	scripts := execute.Lookup(loader.system.Dir)
	logrus.Debugf("[%s] execute scripts lookup, %#v", loader.path, scripts)
	if scripts.Errored() {
		return v.ToErr[ScriptResults](scripts)
	}

	outputs := make(map[string]execute.ExecuteOutput)
	errs := make(map[string]*exec.ExitError)

	for _, file := range scripts.Value() {
		script := env.PathJoin(loader.system.Dir, file)
		logrus.Debugf("[%s] executing %s", loader.path, script)
		res := execute.Execute(script).Await()
		logrus.Debugf("[%s] executed %s, %#v", loader.path, script, res)
		if res.Errored() {
			errs[file] = res.Error().Value()
		} else {
			outputs[file] = res.Value()
		}
	}

	ret := v.Ok[ScriptResults, error](
		ScriptResults{Lookup: scripts.Value(), Output: outputs, Errors: errs},
	)

	logrus.Debugf("[%s] execute scripts done, %#v", loader.path, ret)

	return ret
}

func (loader *PackageLoader) postClean(scripts execute.LookupResult) v.Result[v.Unit, error] {
	logrus.Debugf("[%s] post clean", loader.path)
	target := env.PathJoin(
		loader.system.Dir,
		scriptBackupFolderName,
		loader.filename.ToUnderline().String3(),
	)
	logrus.Debugf("[%s] post clean move target = %s", loader.path, target)
	err := os.MkdirAll(target, 0777)
	logrus.Debugf("[%s] post clean create target, %#v", loader.path, err)
	if err != nil {
		return v.Err[v.Unit](err)
	}

	for _, file := range scripts {
		logrus.Debugf("[%s] post clean, move %s", loader.path, file)
		err := os.Rename(file, env.PathJoin(target, filepath.Base(file)))
		logrus.Debugf("[%s] post clean, moved, %#v", loader.path, err)
		if err != nil {
			return v.Err[v.Unit](err)
		}
	}
	logrus.Debugf("[%s] post clean done", loader.path)
	return v.Ok[v.Unit, error](v.UnitVal)
}

func (loader *PackageLoader) Load() v.Future[v.Result[v.Unit, error]] {
	return v.Async(func() v.Result[v.Unit, error] {
		fmt.Println("Flags", loader.filename.Extname().Full(), loader.filename.Extname())
		// Normal
		if loader.filename.Extname().Full() == ".7z" {
			logrus.Debugf("[%s] load (normal mode), filename = %+v", loader.path, loader.filename)
			return loader.LoadNormal()
		}

		// Localboost, Unimplemented
		if loader.filename.Extname().Flags().Has('l') {
			logrus.Debugf("[%s] load (localboost), filename = %+v", loader.path, loader.filename)
			logrus.Errorf("[%s] load (localboost) unimplemented", loader.path)
			return v.Err[v.Unit](ErrModeUnimplemented)
		}

		// Disabled, Skip
		if loader.filename.Extname().Flags().Has('f') {
			logrus.Debugf("[%s] disabled, filename = %+v", loader.path, loader.filename)
			return v.Ok[v.Unit, error](v.UnitVal)
		}

		// Others, Skip
		logrus.Debugf("[%s] skipped, filename = %+v", loader.path, loader.filename)
		return v.Ok[v.Unit, error](v.UnitVal)
	})
}
