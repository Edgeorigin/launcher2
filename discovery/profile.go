package discovery

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/hydrati/plugin-loader/env"
	"github.com/hydrati/plugin-loader/underline"
	v "github.com/hydrati/plugin-loader/utils/container"
)

type Profile struct {
	Partition env.Partition
	Rootfs    fs.FS
	Root      string
	Version   string
	Packages  map[string]underline.Filename
}

func GetProfiles(excludePartitions ...string) v.Result[[]Profile, error] {
	partitions := env.GetPartitionNoSystem(excludePartitions...)
	if partitions.Errored() {
		return v.ToErr[[]Profile](partitions)
	}

	profiles := make([]Profile, 0)
	for _, part := range partitions.Value() {
		root := env.PathJoin(part.Mountpoint+"/", env.VendorName)
		fmt.Println(env.PathJoin(root, "version.txt"))
		if env.FileExist(root) {
			rootfs := os.DirFS(root)
			if _, err := fs.Stat(rootfs, "version.txt"); !os.IsNotExist(err) {
				versionFile := v.Resuify(rootfs.Open("version.txt"))
				if versionFile.Errored() {
					continue
				}
				versionByte := v.Resuify(io.ReadAll(versionFile.Value()))
				if versionByte.Errored() {
					continue
				}
				versionText0 := env.DecodeBytesGB18030(versionByte.Value())
				if versionText0.Errored() {
					continue
				}
				versionText := versionText0.Value()
				if s, err := fs.Stat(rootfs, "Resource"); !os.IsNotExist(err) && s.IsDir() {
					packages := ScanPackage(rootfs, root, "Resource")
					if packages.Errored() {
						continue
					}

					profiles = append(profiles, Profile{
						Partition: part,
						Rootfs:    rootfs,
						Root:      root,
						Version:   versionText,
						Packages:  packages.Value(),
					})
				}
			}
		}
	}

	return v.Ok[[]Profile, error](profiles)
}

func ScanPackage(rootfs fs.FS, rootdir string, dir string) v.Result[map[string]underline.Filename, error] {
	packages := make(map[string]underline.Filename)
	return v.Resuify(packages, fs.WalkDir(rootfs, dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fs.SkipDir
		}

		if !d.IsDir() {
			if strings.HasPrefix(filepath.Ext(d.Name()), underline.BaseExtname) {
				f := underline.NewFilename(d.Name(), v.None[string]())
				if f.Errored() {
					return f.Error().Value()
				}

				packages[env.PathJoin(rootdir, path)] = f.Value()

			}
		}

		return nil
	}))
}
