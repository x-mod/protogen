package demo

import (
	"context"
	"net/http"

	pb "code.subscriber.one/subscriber/protogen/examples"
	"github.com/x-mod/httpserver"
	"github.com/x-mod/httpserver/grpc"
)

type DemoHTTPService struct {
	rctx    context.Context
	server  *httpserver.Server
	service pb.DemoServer
}

func NewDemoHTTPService(server *httpserver.Server, service pb.DemoServer) *DemoHTTPService {
	return &DemoHTTPService{server: server, service: service}
}

func (svc *DemoHTTPService) Serve(ctx context.Context) error {
	svc.server.Route(
		httpserver.Method("post"),
		httpserver.Pattern("/v1/example/hello"),
		httpserver.Handler(svc.Hello),
	)
	svc.server.Route(
		httpserver.Method("get"),
		httpserver.Pattern(grpc.URIFormat("v1", "demo", "Demo", "Greet")),
		httpserver.Handler(svc.Greet),
	)

	svc.rctx = ctx
	return svc.server.Serve(ctx)
}

func (svc *DemoHTTPService) Stop() {
	svc.server.Shutdown(svc.rctx)
}

func (svc *DemoHTTPService) Hello(ctx context.Context, wr http.ResponseWriter, req *http.Request) {
	in := new(pb.HelloReq)
	err := grpc.PBRequest(req, in)
	if err != nil {
		grpc.PBResponse(wr, nil, err)
		return
	}
	out, err := svc.service.Hello(grpc.PBContext(req, ctx), in)
	grpc.PBResponse(wr, out, err)
}

func (svc *DemoHTTPService) Greet(ctx context.Context, wr http.ResponseWriter, req *http.Request) {
	in := new(pb.HelloReq)
	err := grpc.PBRequest(req, in)
	if err != nil {
		grpc.PBResponse(wr, nil, err)
		return
	}
	out, err := svc.service.Greet(grpc.PBContext(req, ctx), in)
	grpc.PBResponse(wr, out, err)
}
