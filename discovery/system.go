package discovery

import (
	"errors"
	"io/fs"
	"os"

	"github.com/hydrati/plugin-loader/env"
	v "github.com/hydrati/plugin-loader/utils/container"
)

var (
	ErrNotFoundSystem = errors.New("error: not found system")
)

type System struct {
	Partition  env.Partition
	Fs         fs.FS
	Dir        string
	AllProfile []Profile
	Profile    v.Option[*Profile]
}

func GetSystem() *System {
	p := env.GetSystemPartition().Value()

	workdir := env.PathJoin(env.EnvProgramFiles(), env.VendorName)
	f := v.Must(os.Stat(workdir))
	if !f.IsDir() {
		panic(ErrNotFoundSystem)
	}

	workfs := os.DirFS(workdir)

	profiles := GetProfiles().Value()
	if len(profiles) < 1 {
		return &System{p, workfs, workdir, profiles, v.None[*Profile]()}
	}

	return &System{p, workfs, workdir, profiles, v.Some(&profiles[0])}
}
