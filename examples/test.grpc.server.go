package demo

import (
	"context"
	"log"
	"net"

	pb "code.subscriber.one/subscriber/protogen/examples"
	"github.com/x-mod/errors"
	"google.golang.org/grpc"
)

type DemoGrpcService struct {
	addr string
	gsrv *grpc.Server
}

func NewDemoGrpcService(addr string, opts ...grpc.ServerOption) *DemoGrpcService {
	service := &DemoGrpcService{addr: addr}
	gserver := grpc.NewServer(opts...)
	service.gsrv = gserver
	return service
}

func (svc *DemoGrpcService) Serve(ctx context.Context) error {
	pb.RegisterDemoServer(svc.gsrv, svc)
	ln, err := net.Listen("tcp", svc.addr)
	if err != nil {
		return err
	}
	log.Println("service (Demo) serving at ", svc.addr)
	return svc.gsrv.Serve(ln)
}

func (svc *DemoGrpcService) Stop() {
	svc.gsrv.Stop()
}

func (svc *DemoGrpcService) Hello(ctx context.Context, req *pb.HelloReq) (*pb.HelloResp, error) {
	return nil, errors.New("Hello unimplemented")
}
func (svc *DemoGrpcService) Greet(ctx context.Context, req *pb.HelloReq) (*pb.HelloResp, error) {
	return nil, errors.New("Greet unimplemented")
}
