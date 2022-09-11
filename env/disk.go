package env

import (
	"errors"

	"github.com/shirou/gopsutil/v3/disk"

	u "github.com/hydrati/plugin-loader/utils"
	v "github.com/hydrati/plugin-loader/utils/container"
)

var (
	ErrNotFoundSystemDrive = errors.New("error: not found system drive")
)

type Partition = disk.PartitionStat

func GetPartition(exclude ...string) v.Result[[]Partition, error] {
	p := make([]Partition, 0)
	ex := u.MkStringSet(exclude...)

	s, err := disk.Partitions(true)
	if err != nil {
		return v.Err[[]Partition](err)
	}

	for _, i := range s {
		if _, ok := ex[i.Device]; !ok {
			p = append(p, i)
		}
	}

	return v.Ok[[]Partition, error](p)
}

func GetPartitionNoSystem(exclude ...string) v.Result[[]Partition, error] {
	return GetPartition(append(exclude, EnvSystemDrive())...)
}

func GetSystemPartition() v.Result[Partition, error] {
	s, err := disk.Partitions(true)
	p := EnvSystemDrive()
	if err != nil {
		return v.Err[Partition](err)
	}

	for _, i := range s {
		if i.Device == p {
			return v.Ok[Partition, error](i)
		}
	}

	return v.Err[Partition](ErrNotFoundSystemDrive)
}
