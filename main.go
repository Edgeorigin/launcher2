package main

import (
	"fmt"

	"github.com/hydrati/plugin-loader/execute"
	v "github.com/hydrati/plugin-loader/utils/container"
	"golang.org/x/text/encoding/simplifiedchinese"
)

func main() {
	// file := archive.OpenArchive(
	// 	"./WizTree_4.10_Hydration.7z",
	// 	v.None[string](),
	// ).Value()

	// extractor := archive.NewExtractor(file.All(), "./WizTree")
	// fmt.Println(extractor.Execute(v.Some(16)).Await())
	// fmt.Println(u.Must(exec.LookPath("cmd")))

	output, err := execute.Execute("./test.wcs").Await().Option()

	decoder := simplifiedchinese.GB18030.NewDecoder()

	if err.Has() {
		utf8 := v.Must(decoder.Bytes(err.Value().Stderr))

		fmt.Println(err, string(utf8))
	} else {
		fmt.Println(output)
		if output.Value().Output != nil {
			utf8 := v.Must(decoder.Bytes(output.Value().Output))
			fmt.Println(string(utf8))
		}
	}

}
