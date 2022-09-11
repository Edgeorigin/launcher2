package env

import (
	"github.com/shirou/gopsutil/v3/disk"

	u "github.com/hydrati/plugin-loader/utils"
	v "github.com/hydrati/plugin-loader/utils/container"
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
