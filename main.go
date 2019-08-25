package main

import (
	_ "github.com/liujianping/protogen/gen"
	"github.com/x-mod/build"
	"github.com/x-mod/cmd"
)

//go:generate go-bindata -prefix tpl -nometadata -o tpl/bindata.go -ignore bindata.go -pkg tpl tpl/...
func main() {
	cmd.Version(build.String())
	cmd.PersistentFlags().StringP("input", "i", ".", "input directory")
	cmd.PersistentFlags().StringP("output", "o", ".", "output directory")
	cmd.PersistentFlags().StringP("protobuf-suffix", "p", ".proto", "protobuf suffix")
	cmd.PersistentFlags().StringP("template-suffix", "t", ".gogo", "template suffix")
	cmd.Exit(cmd.Execute())
}
