package main

// import (
// 	"fmt"

// 	"github.com/hydrati/plugin-loader/execute"
// 	v "github.com/hydrati/plugin-loader/utils/container"

// )

import (
	"fmt"

	"github.com/hydrati/plugin-loader/discovery"
	_ "github.com/hydrati/plugin-loader/env"
	"github.com/hydrati/plugin-loader/loader"
)

func main() {
	system := discovery.GetSystem()
	r := loader.LoadSystem(system)
	if r.Errored() {
		fmt.Println(r.Error().Value())
	}
}
