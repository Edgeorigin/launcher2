package main

// import (
// 	"fmt"

// 	"github.com/hydrati/plugin-loader/execute"
// 	v "github.com/hydrati/plugin-loader/utils/container"

// )

import (
	"fmt"

	"github.com/hydrati/plugin-loader/env"
	v "github.com/hydrati/plugin-loader/utils/container"
	"github.com/shirou/gopsutil/v3/disk"

	"github.com/hydrati/plugin-loader/discovery"
)

func main() {
	for _, p := range v.Must(disk.Partitions(true)) {
		if p.Mountpoint != env.EnvSystemDrive() {
			// fmt.Println(p)
		}
	}

	fmt.Println(discovery.GetProfiles())
}
