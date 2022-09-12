package env

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	u "github.com/hydrati/plugin-loader/utils"
	v "github.com/hydrati/plugin-loader/utils/container"
	"github.com/sirupsen/logrus"
)

var (
	VendorName     = "Edgeless"
	ErrUnsupported = errors.New("error: unsuported platform")
)

var (
	supportedOS   = u.MkStringSet("windows")
	supportedArch = u.MkStringSet("amd64")
)

func EnvIsSupportedOS() bool {
	_, ok := supportedOS[runtime.GOOS]
	return ok
}

func EnvIsSupportedArch() bool {
	_, ok := supportedArch[runtime.GOARCH]

	return ok
}

func EnvIsSupport() bool {
	return EnvIsSupportedOS() || EnvIsSupportedArch()
}

func EnvSystemRoot() string {
	return os.Getenv("SYSTEMROOT")
}

func EnvProgramFiles() string {
	return os.Getenv("PROGRAMFILES")
}

func EnvSystemDrive() string {
	return os.Getenv("SYSTEMDRIVE")
}

func Cwd() string {
	return v.Must(os.Getwd())
}

func FindBin(name string) string {
	bin := v.Must(exec.LookPath(name))
	logrus.Debugf("find bin %s -> %s", name, bin)
	return bin
}

func PathJoin(path ...string) string {
	return filepath.Clean(filepath.Join(path...))
}

func PathResolve(path string) string {
	if filepath.IsAbs(path) {
		return path
	}

	return PathJoin(v.Must(os.Getwd()), path)
}

func FileExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.Debugf("runtime version %s", runtime.Version())
	if !EnvIsSupport() {
		logrus.Debugf("unsupported platform")
		panic(ErrUnsupported)
	}
	logrus.Debugf("supported platform")
}
