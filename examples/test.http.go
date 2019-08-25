package demo

import (
	"context"
	"net/http"

	"github.com/x-mod/httpclient"
	client "github.com/x-mod/httpclient/grpc"
	server "github.com/x-mod/httpclient/grpc"
	"github.com/x-mod/options"
	"google.golang.org/grpc"
)

//httpserver integrating
func RegisterDemoHTTPServer(s *server.HTTPServer, srv DemoServer) error {
	return s.RegisterService(&_Demo_HTTP_serviceDesc, srv)
}

var _Demo_HTTP_serviceDesc = server.ServiceDescription{
	PackageName:  "demo",
	ServiceName:  "Demo",
	Implemention: (*DemoServer)(nil),
	Option: &options.ServiceOption{
		Version: "v1",
	},
	Methods: []server.MethodDescription{
		{
			MethodName: "Hello",
			Handler:    _Demo_Hello_HTTP_Handler,
			Option: &options.HttpOption{
				Method: "post",
				Uri:    "/v1/example/hello",
			},
		},
	},
}

func _Demo_Hello_HTTP_Handler(srv interface{}, ctx context.Context, wr http.ResponseWriter, req *http.Request) {
	in := new(HelloReq)
	err := server.PBRequest(req, in)
	if err != nil {
		server.PBResponse(wr, nil, err)
		return
	}
	out, err := srv.(DemoServer).Hello(server.PBContext(req, ctx), in)
	if err != nil {
		server.PBResponse(wr, nil, err)
		return
	}
	server.PBResponse(wr, out, err)
}

//httpclient integrating
type HTTPDemoClient struct {
	*client.HTTPClient
}

func NewHTTPDemoClient(opts ...client.HTTPClientOpt) DemoClient {
	gopts := []client.HTTPClientOpt{}
	gopts = append(gopts, client.Version("v1"))
	gopts = append(gopts, client.PackageName("demo"))
	gopts = append(gopts, client.ServiceName("Demo"))
	gopts = append(gopts, client.Version("v1"))
	gopts = append(gopts, opts...)
	return &HTTPDemoClient{
		HTTPClient: client.NewHTTPClient(gopts...),
	}
}

func (c *HTTPDemoClient) Hello(ctx context.Context, in *HelloReq, opts ...grpc.CallOption) (*HelloResp, error) {
	req, err := c.MakeRequest(
		"Hello",
		httpclient.Method("post"),
		httpclient.URL(httpclient.URI("/v1/example/hello")),
		httpclient.Content(httpclient.PBJSON(in)),
	)
	if err != nil {
		return nil, err
	}
	out := new(HelloResp)
	if err := c.Execute(ctx, req, client.PBResponse(out)); err != nil {
		return nil, err
	}
	return out, nil
}
