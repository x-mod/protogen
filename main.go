package main

import (
	_ "code.subscriber.one/subscriber/protogen/gen"
	"github.com/x-mod/build"
	"github.com/x-mod/cmd"
)

//go:generate go-bindata -prefix tpl -nometadata -o tpl/bindata.go -ignore bindata.go -pkg tpl tpl/...
func main() {
	cmd.Version(build.String())
	cmd.PersistentFlags().StringP("input", "i", ".", "input directory")
	cmd.PersistentFlags().StringP("output", "o", ".", "output directory")
	cmd.PersistentFlags().StringP("protobuf-suffix", "x", ".proto", "protobuf suffix")
	cmd.PersistentFlags().StringP("protobuf-files", "p", "", "protobuf gen files import path")
	cmd.PersistentFlags().StringP("template-suffix", "t", ".gogo", "template suffix")
	cmd.Execute()
}
