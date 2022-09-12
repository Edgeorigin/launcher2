package main

// import (
// 	"fmt"

// 	"github.com/hydrati/plugin-loader/execute"
// 	v "github.com/hydrati/plugin-loader/utils/container"

// )

import (
	_ "fmt"

	"github.com/hydrati/plugin-loader/discovery"
	_ "github.com/hydrati/plugin-loader/env"
	"github.com/hydrati/plugin-loader/loader"
)

func main() {
	system := discovery.GetSystem()
	loader.LoadSystem(system).Value()
}
