package gen

import (
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/x-mod/cmd"
	"github.com/x-mod/errors"
)

func HttpMain(c *cmd.Command, args []string) error {
	in, err := filepath.Abs(viper.GetString("input"))
	if err != nil {
		return errors.Annotate(err, "input")
	}
	out, err := filepath.Abs(viper.GetString("output"))
	if err != nil {
		return errors.Annotate(err, "output")
	}
	ProtoSuffix = viper.GetString("protobuf-suffix")
	suffix := viper.GetString("template-suffix")

	return generate(in, out, "http", suffix, "http")
}
