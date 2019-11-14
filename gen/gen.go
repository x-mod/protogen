package gen

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/emicklei/proto"
	"github.com/spf13/viper"
	"github.com/x-mod/cmd"
	"github.com/x-mod/errors"
)

var ProtoSuffix = ".proto"

type AST struct {
	*proto.Proto
	Filename string
	ImportPb string
	Package  *proto.Package
	Imports  []*proto.Import
	Options  []*proto.Option
	Messages []*proto.Message
	Enums    []*Enum
	Services []*Service
}

type Service struct {
	*proto.Service
	Options []*proto.Option
	Methods []*proto.RPC
}

type Enum struct {
	*proto.Enum
	EnumFields []*proto.EnumField
}

func build(filename string, pt *proto.Proto) *AST {
	ast := &AST{
		Proto:    pt,
		Filename: filename,
		Imports:  []*proto.Import{},
		Options:  []*proto.Option{},
		Enums:    []*Enum{},
		Messages: []*proto.Message{},
		Services: []*Service{},
	}
	for _, each := range pt.Elements {
		if v, ok := each.(*proto.Package); ok {
			ast.Package = v
		}
		if v, ok := each.(*proto.Import); ok {
			ast.Imports = append(ast.Imports, v)
		}
		if v, ok := each.(*proto.Option); ok {
			ast.Options = append(ast.Options, v)
		}
		if v, ok := each.(*proto.Enum); ok {
			e := &Enum{
				Enum:       v,
				EnumFields: []*proto.EnumField{},
			}
			for _, vv := range v.Elements {
				if f, ok := vv.(*proto.EnumField); ok {
					e.EnumFields = append(e.EnumFields, f)
				}
			}
			ast.Enums = append(ast.Enums, e)
		}
		if v, ok := each.(*proto.Message); ok {
			ast.Messages = append(ast.Messages, v)
		}
		if v, ok := each.(*proto.Service); ok {
			svc := &Service{
				Service: v,
				Options: []*proto.Option{},
				Methods: []*proto.RPC{},
			}
			for _, vv := range v.Elements {
				if opt, ok := vv.(*proto.Option); ok {
					svc.Options = append(svc.Options, opt)
				}
				if rpc, ok := vv.(*proto.RPC); ok {
					svc.Methods = append(svc.Methods, rpc)
				}
			}
			ast.Services = append(ast.Services, svc)
		}
	}
	return ast
}

func getInputFilesBySuffix(dir string, suffix string) ([]string, error) {
	stat, err := os.Stat(dir)

	if err != nil && os.IsNotExist(err) {
		return nil, err
	}

	files := []string{}
	if !stat.IsDir() && strings.HasSuffix(dir, suffix) {
		files = append(files, dir)
		return files, nil
	}

	if err := filepath.Walk(dir, func(src string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(info.Name(), suffix) {
			files = append(files, src)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return files, nil
}

func generate(in string, out string, pb string, suffix string, prefix string, objects ...string) error {
	tpl, err := getTemplate(prefix)
	if err != nil {
		return errors.Annotate(err, "get templates")
	}

	protos, err := getInputFilesBySuffix(in, ProtoSuffix)
	if err != nil {
		return errors.Annotate(err, "get proto files")
	}

	for _, infile := range protos {
		_, filename := filepath.Split(infile)
		filenameTrimSuffix := strings.TrimSuffix(filename, ProtoSuffix)
		rd, err := os.Open(infile)
		if err != nil {
			return errors.Annotatef(err, "open file %s", infile)
		}
		defer rd.Close()

		parser := proto.NewParser(rd)
		pt, err := parser.Parse()
		if err != nil {
			return errors.Annotatef(err, "parse file %s", infile)
		}

		ast := build(filename, pt)
		ast.ImportPb = pb
		for _, obj := range objects {
			outfile := filepath.Join(out, fmt.Sprintf("%s.%s.go", filenameTrimSuffix, obj))
			wr, err := os.OpenFile(outfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				return errors.Annotatef(err, "write file %s", outfile)
			}
			if err := tpl.ExecuteTemplate(wr, obj, ast); err != nil {
				return errors.Annotatef(err, "template execute %s", obj)
			}
			if err := wr.Close(); err != nil {
				return errors.Annotatef(err, "close file %s", outfile)
			}
			oscmd := exec.Command("gofmt", "-w", outfile)
			if err := oscmd.Run(); err != nil {
				log.Println("gofmt: ", err)
			}
		}
	}
	return copyFilesExcludeSuffix(prefix, suffix, out, true)
}

func init() {
	cmd.Add(
		cmd.Path("/http/server"),
		cmd.Short("generate http server files from protobuf"),
		cmd.Main(HttpServer),
	)
	cmd.Add(
		cmd.Path("/http/client"),
		cmd.Short("generate http client files from protobuf"),
		cmd.Main(HttpClient),
	)
	cmd.Add(
		cmd.Path("/grpc/server"),
		cmd.Short("generate grpc server files from protobuf"),
		cmd.Main(GrpcServer),
	)
	cmd.Add(
		cmd.Path("/grpc/client"),
		cmd.Short("generate grpc client files from protobuf"),
		cmd.Main(GrpcClient),
	)
}

func HttpServer(c *cmd.Command, args []string) error {
	return Main(c, []string{"http", "server"})
}
func HttpClient(c *cmd.Command, args []string) error {
	return Main(c, []string{"http", "client"})
}
func GrpcServer(c *cmd.Command, args []string) error {
	return Main(c, []string{"grpc", "server"})
}
func GrpcClient(c *cmd.Command, args []string) error {
	return Main(c, []string{"grpc", "client"})
}

func Main(c *cmd.Command, args []string) error {
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
	pbFiles := viper.GetString("protobuf-files")
	return generate(in, out, pbFiles, suffix, args[0], strings.Join(args, "."))
}
